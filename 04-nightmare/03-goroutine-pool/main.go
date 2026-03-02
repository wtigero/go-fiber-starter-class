package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ============ Basic Worker Pool ============

type Job func()

type WorkerPool struct {
	workers    int
	jobQueue   chan Job
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc

	// Stats
	processed  int64
	processing int64
}

func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &WorkerPool{
		workers:  workers,
		jobQueue: make(chan Job, queueSize),
		ctx:      ctx,
		cancel:   cancel,
	}
	pool.start()
	return pool
}

func (p *WorkerPool) start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case job, ok := <-p.jobQueue:
			if !ok {
				return
			}
			atomic.AddInt64(&p.processing, 1)
			job()
			atomic.AddInt64(&p.processing, -1)
			atomic.AddInt64(&p.processed, 1)
		}
	}
}

func (p *WorkerPool) Submit(job Job) error {
	select {
	case <-p.ctx.Done():
		return fmt.Errorf("pool is closed")
	case p.jobQueue <- job:
		return nil
	}
}

func (p *WorkerPool) SubmitWait(job Job) error {
	done := make(chan struct{})
	wrappedJob := func() {
		job()
		close(done)
	}

	if err := p.Submit(wrappedJob); err != nil {
		return err
	}

	<-done
	return nil
}

func (p *WorkerPool) Stop() {
	p.cancel()
	close(p.jobQueue)
	p.wg.Wait()
}

func (p *WorkerPool) Stats() (processed, processing int64, queueLen int) {
	return atomic.LoadInt64(&p.processed),
		atomic.LoadInt64(&p.processing),
		len(p.jobQueue)
}

// ============ Advanced Pool with Results ============

type Task[T any] struct {
	Fn     func() T
	Result chan T
}

type ResultPool[T any] struct {
	workers  int
	tasks    chan Task[T]
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewResultPool[T any](workers int, queueSize int) *ResultPool[T] {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &ResultPool[T]{
		workers: workers,
		tasks:   make(chan Task[T], queueSize),
		ctx:     ctx,
		cancel:  cancel,
	}
	pool.start()
	return pool
}

func (p *ResultPool[T]) start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-p.ctx.Done():
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}
					result := task.Fn()
					task.Result <- result
				}
			}
		}()
	}
}

func (p *ResultPool[T]) Submit(fn func() T) <-chan T {
	result := make(chan T, 1)
	task := Task[T]{Fn: fn, Result: result}

	select {
	case <-p.ctx.Done():
		close(result)
	case p.tasks <- task:
	}

	return result
}

func (p *ResultPool[T]) Stop() {
	p.cancel()
	close(p.tasks)
	p.wg.Wait()
}

// ============ Dynamic Scaling Pool ============

type DynamicPool struct {
	minWorkers  int
	maxWorkers  int
	workers     int64
	jobQueue    chan Job
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.Mutex

	processed   int64
	queueHighWater int64
}

func NewDynamicPool(minWorkers, maxWorkers, queueSize int) *DynamicPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &DynamicPool{
		minWorkers: minWorkers,
		maxWorkers: maxWorkers,
		jobQueue:   make(chan Job, queueSize),
		ctx:        ctx,
		cancel:     cancel,
	}

	// Start minimum workers
	for i := 0; i < minWorkers; i++ {
		pool.addWorker()
	}

	// Start scaler
	go pool.scaler()

	return pool
}

func (p *DynamicPool) addWorker() {
	atomic.AddInt64(&p.workers, 1)
	go func() {
		defer atomic.AddInt64(&p.workers, -1)

		idleTimer := time.NewTimer(30 * time.Second)
		defer idleTimer.Stop()

		for {
			select {
			case <-p.ctx.Done():
				return
			case job, ok := <-p.jobQueue:
				if !ok {
					return
				}
				idleTimer.Reset(30 * time.Second)
				job()
				atomic.AddInt64(&p.processed, 1)
			case <-idleTimer.C:
				// Idle timeout, remove worker if above minimum
				if atomic.LoadInt64(&p.workers) > int64(p.minWorkers) {
					return
				}
				idleTimer.Reset(30 * time.Second)
			}
		}
	}()
}

func (p *DynamicPool) scaler() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			queueLen := int64(len(p.jobQueue))
			workers := atomic.LoadInt64(&p.workers)

			// Update high water mark
			if queueLen > atomic.LoadInt64(&p.queueHighWater) {
				atomic.StoreInt64(&p.queueHighWater, queueLen)
			}

			// Scale up if queue is building up
			if queueLen > workers && workers < int64(p.maxWorkers) {
				p.addWorker()
			}
		}
	}
}

func (p *DynamicPool) Submit(job Job) {
	select {
	case <-p.ctx.Done():
		return
	case p.jobQueue <- job:
	}
}

func (p *DynamicPool) Stop() {
	p.cancel()
	close(p.jobQueue)
}

func (p *DynamicPool) Stats() (workers, processed, queueLen, highWater int64) {
	return atomic.LoadInt64(&p.workers),
		atomic.LoadInt64(&p.processed),
		int64(len(p.jobQueue)),
		atomic.LoadInt64(&p.queueHighWater)
}

// ============ Global Pools ============

var (
	basicPool   *WorkerPool
	dynamicPool *DynamicPool
	resultPool  *ResultPool[int]
)

func init() {
	basicPool = NewWorkerPool(runtime.NumCPU(), 1000)
	dynamicPool = NewDynamicPool(2, runtime.NumCPU()*2, 1000)
	resultPool = NewResultPool[int](runtime.NumCPU(), 100)
}

// ============ Handlers ============

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(logger.New())

	// Routes
	app.Get("/", homeHandler)
	app.Get("/stats", statsHandler)

	// Comparison endpoints
	app.Get("/no-pool", noPoolHandler)
	app.Get("/basic-pool", basicPoolHandler)
	app.Get("/dynamic-pool", dynamicPoolHandler)
	app.Get("/result-pool", resultPoolHandler)

	// Benchmark
	app.Get("/bench", benchHandler)

	log.Println("🚀 Goroutine Pool Demo running on http://localhost:3000")
	log.Println("")
	log.Println("Endpoints:")
	log.Println("  /no-pool      - Creates new goroutine each request")
	log.Println("  /basic-pool   - Uses fixed worker pool")
	log.Println("  /dynamic-pool - Uses auto-scaling pool")
	log.Println("  /result-pool  - Pool with result channel")
	log.Println("  /stats        - Pool statistics")
	log.Println("  /bench        - Run benchmark comparison")

	log.Fatal(app.Listen(":3000"))
}

func homeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Goroutine Pool Demo",
		"pools": fiber.Map{
			"basic":   fmt.Sprintf("%d workers", runtime.NumCPU()),
			"dynamic": fmt.Sprintf("%d-%d workers", 2, runtime.NumCPU()*2),
			"result":  fmt.Sprintf("%d workers", runtime.NumCPU()),
		},
	})
}

func statsHandler(c *fiber.Ctx) error {
	bp, bProcessing, bQueue := basicPool.Stats()
	dWorkers, dProcessed, dQueue, dHighWater := dynamicPool.Stats()

	return c.JSON(fiber.Map{
		"basic_pool": fiber.Map{
			"processed":  bp,
			"processing": bProcessing,
			"queue_len":  bQueue,
		},
		"dynamic_pool": fiber.Map{
			"workers":         dWorkers,
			"processed":       dProcessed,
			"queue_len":       dQueue,
			"queue_high_water": dHighWater,
		},
		"goroutines": runtime.NumGoroutine(),
	})
}

// ❌ NO POOL - สร้าง goroutine ใหม่ทุกครั้ง
func noPoolHandler(c *fiber.Ctx) error {
	var wg sync.WaitGroup
	results := make([]int, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// Simulate work
			time.Sleep(10 * time.Millisecond)
			results[idx] = idx * 2
		}(i)
	}

	wg.Wait()

	return c.JSON(fiber.Map{
		"method":     "no_pool",
		"results":    results,
		"goroutines": runtime.NumGoroutine(),
		"note":       "Created 10 new goroutines",
	})
}

// ✅ BASIC POOL - ใช้ fixed worker pool
func basicPoolHandler(c *fiber.Ctx) error {
	var wg sync.WaitGroup
	results := make([]int, 10)
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		idx := i
		basicPool.Submit(func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			results[idx] = idx * 2
			mu.Unlock()
		})
	}

	wg.Wait()
	processed, _, _ := basicPool.Stats()

	return c.JSON(fiber.Map{
		"method":          "basic_pool",
		"results":         results,
		"total_processed": processed,
		"goroutines":      runtime.NumGoroutine(),
		"note":            "Reused existing workers",
	})
}

// ✅ DYNAMIC POOL - auto-scaling
func dynamicPoolHandler(c *fiber.Ctx) error {
	var wg sync.WaitGroup
	results := make([]int, 10)
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		idx := i
		dynamicPool.Submit(func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			results[idx] = idx * 2
			mu.Unlock()
		})
	}

	wg.Wait()
	workers, processed, _, _ := dynamicPool.Stats()

	return c.JSON(fiber.Map{
		"method":          "dynamic_pool",
		"results":         results,
		"active_workers":  workers,
		"total_processed": processed,
		"goroutines":      runtime.NumGoroutine(),
		"note":            "Auto-scaling workers",
	})
}

// ✅ RESULT POOL - with typed results
func resultPoolHandler(c *fiber.Ctx) error {
	// Submit tasks and collect result channels
	resultChans := make([]<-chan int, 10)

	for i := 0; i < 10; i++ {
		idx := i
		resultChans[i] = resultPool.Submit(func() int {
			time.Sleep(10 * time.Millisecond)
			return idx * 2
		})
	}

	// Collect results
	results := make([]int, 10)
	for i, ch := range resultChans {
		results[i] = <-ch
	}

	return c.JSON(fiber.Map{
		"method":     "result_pool",
		"results":    results,
		"goroutines": runtime.NumGoroutine(),
		"note":       "Type-safe results with channels",
	})
}

// Benchmark comparison
func benchHandler(c *fiber.Ctx) error {
	iterations := 100
	jobCount := 100

	// Without pool
	start := time.Now()
	for i := 0; i < iterations; i++ {
		var wg sync.WaitGroup
		for j := 0; j < jobCount; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(1 * time.Millisecond)
			}()
		}
		wg.Wait()
	}
	noPoolDuration := time.Since(start)

	// With pool
	start = time.Now()
	for i := 0; i < iterations; i++ {
		var wg sync.WaitGroup
		for j := 0; j < jobCount; j++ {
			wg.Add(1)
			basicPool.Submit(func() {
				defer wg.Done()
				time.Sleep(1 * time.Millisecond)
			})
		}
		wg.Wait()
	}
	withPoolDuration := time.Since(start)

	return c.JSON(fiber.Map{
		"iterations":       iterations,
		"jobs_per_iter":    jobCount,
		"total_jobs":       iterations * jobCount,
		"without_pool":     noPoolDuration.String(),
		"with_pool":        withPoolDuration.String(),
		"speedup":          fmt.Sprintf("%.2fx", float64(noPoolDuration)/float64(withPoolDuration)),
		"goroutines_saved": fmt.Sprintf("%d goroutines", iterations*jobCount),
	})
}

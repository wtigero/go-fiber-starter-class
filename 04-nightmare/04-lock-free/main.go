package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ============ Lock-Free Stack (Treiber Stack) ============

type StackNode struct {
	value interface{}
	next  unsafe.Pointer // *StackNode
}

type LockFreeStack struct {
	head unsafe.Pointer // *StackNode
	size int64
}

func NewLockFreeStack() *LockFreeStack {
	return &LockFreeStack{}
}

func (s *LockFreeStack) Push(value interface{}) {
	node := &StackNode{value: value}

	for {
		oldHead := atomic.LoadPointer(&s.head)
		node.next = oldHead

		if atomic.CompareAndSwapPointer(&s.head, oldHead, unsafe.Pointer(node)) {
			atomic.AddInt64(&s.size, 1)
			return
		}
		// CAS failed, retry
	}
}

func (s *LockFreeStack) Pop() (interface{}, bool) {
	for {
		oldHead := atomic.LoadPointer(&s.head)
		if oldHead == nil {
			return nil, false
		}

		node := (*StackNode)(oldHead)
		newHead := node.next

		if atomic.CompareAndSwapPointer(&s.head, oldHead, newHead) {
			atomic.AddInt64(&s.size, -1)
			return node.value, true
		}
		// CAS failed, retry
	}
}

func (s *LockFreeStack) Size() int64 {
	return atomic.LoadInt64(&s.size)
}

// ============ Lock-Free Counter ============

type LockFreeCounter struct {
	value int64
}

func NewLockFreeCounter() *LockFreeCounter {
	return &LockFreeCounter{}
}

func (c *LockFreeCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

func (c *LockFreeCounter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

func (c *LockFreeCounter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

func (c *LockFreeCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func (c *LockFreeCounter) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.value, old, new)
}

// ============ Lock-Free Queue (Michael & Scott) ============

type QueueNode struct {
	value interface{}
	next  unsafe.Pointer // *QueueNode
}

type LockFreeQueue struct {
	head unsafe.Pointer // *QueueNode
	tail unsafe.Pointer // *QueueNode
	size int64
}

func NewLockFreeQueue() *LockFreeQueue {
	// Dummy node
	node := &QueueNode{}
	q := &LockFreeQueue{
		head: unsafe.Pointer(node),
		tail: unsafe.Pointer(node),
	}
	return q
}

func (q *LockFreeQueue) Enqueue(value interface{}) {
	node := &QueueNode{value: value}

	for {
		tail := atomic.LoadPointer(&q.tail)
		tailNode := (*QueueNode)(tail)
		next := atomic.LoadPointer(&tailNode.next)

		if tail == atomic.LoadPointer(&q.tail) {
			if next == nil {
				// Try to link node at end
				if atomic.CompareAndSwapPointer(&tailNode.next, nil, unsafe.Pointer(node)) {
					// Success, try to swing tail
					atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(node))
					atomic.AddInt64(&q.size, 1)
					return
				}
			} else {
				// Tail is behind, try to advance
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			}
		}
	}
}

func (q *LockFreeQueue) Dequeue() (interface{}, bool) {
	for {
		head := atomic.LoadPointer(&q.head)
		tail := atomic.LoadPointer(&q.tail)
		headNode := (*QueueNode)(head)
		next := atomic.LoadPointer(&headNode.next)

		if head == atomic.LoadPointer(&q.head) {
			if head == tail {
				if next == nil {
					return nil, false // Empty
				}
				// Tail is behind, advance it
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			} else {
				// Read value before CAS
				nextNode := (*QueueNode)(next)
				value := nextNode.value

				if atomic.CompareAndSwapPointer(&q.head, head, next) {
					atomic.AddInt64(&q.size, -1)
					return value, true
				}
			}
		}
	}
}

func (q *LockFreeQueue) Size() int64 {
	return atomic.LoadInt64(&q.size)
}

// ============ Mutex-based versions for comparison ============

type MutexStack struct {
	data []interface{}
	mu   sync.Mutex
}

func NewMutexStack() *MutexStack {
	return &MutexStack{data: make([]interface{}, 0)}
}

func (s *MutexStack) Push(value interface{}) {
	s.mu.Lock()
	s.data = append(s.data, value)
	s.mu.Unlock()
}

func (s *MutexStack) Pop() (interface{}, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.data) == 0 {
		return nil, false
	}

	value := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return value, true
}

type MutexCounter struct {
	value int64
	mu    sync.Mutex
}

func (c *MutexCounter) Increment() int64 {
	c.mu.Lock()
	c.value++
	v := c.value
	c.mu.Unlock()
	return v
}

func (c *MutexCounter) Value() int64 {
	c.mu.Lock()
	v := c.value
	c.mu.Unlock()
	return v
}

// ============ Global instances ============

var (
	lockFreeStack   = NewLockFreeStack()
	lockFreeQueue   = NewLockFreeQueue()
	lockFreeCounter = NewLockFreeCounter()
	mutexStack      = NewMutexStack()
	mutexCounter    = &MutexCounter{}
)

// ============ Handlers ============

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(logger.New())

	// Routes
	app.Get("/", homeHandler)
	app.Get("/stats", statsHandler)

	// Lock-Free operations
	app.Post("/stack/push", stackPushHandler)
	app.Get("/stack/pop", stackPopHandler)
	app.Post("/queue/enqueue", queueEnqueueHandler)
	app.Get("/queue/dequeue", queueDequeueHandler)
	app.Post("/counter/increment", counterIncrHandler)

	// Benchmarks
	app.Get("/bench/counter", benchCounterHandler)
	app.Get("/bench/stack", benchStackHandler)

	log.Println("🚀 Lock-Free Data Structures Demo running on http://localhost:3000")
	log.Println("")
	log.Println("Endpoints:")
	log.Println("  /stack/push, /stack/pop     - Lock-free stack operations")
	log.Println("  /queue/enqueue, /queue/dequeue - Lock-free queue operations")
	log.Println("  /counter/increment          - Lock-free counter")
	log.Println("  /bench/counter, /bench/stack - Benchmarks")
	log.Println("  /stats                      - Current state")

	log.Fatal(app.Listen(":3000"))
}

func homeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Lock-Free Data Structures Demo",
		"structures": fiber.Map{
			"stack":   "Treiber Stack (CAS-based)",
			"queue":   "Michael & Scott Queue",
			"counter": "Atomic counter",
		},
		"endpoints": fiber.Map{
			"stack":   "POST /stack/push, GET /stack/pop",
			"queue":   "POST /queue/enqueue, GET /queue/dequeue",
			"counter": "POST /counter/increment",
			"bench":   "GET /bench/counter, GET /bench/stack",
		},
	})
}

func statsHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"lock_free_stack_size":   lockFreeStack.Size(),
		"lock_free_queue_size":   lockFreeQueue.Size(),
		"lock_free_counter":      lockFreeCounter.Value(),
		"mutex_counter":          mutexCounter.Value(),
		"goroutines":             runtime.NumGoroutine(),
	})
}

func stackPushHandler(c *fiber.Ctx) error {
	value := c.Query("value", "item")
	lockFreeStack.Push(value)

	return c.JSON(fiber.Map{
		"action":     "push",
		"value":      value,
		"stack_size": lockFreeStack.Size(),
	})
}

func stackPopHandler(c *fiber.Ctx) error {
	value, ok := lockFreeStack.Pop()

	return c.JSON(fiber.Map{
		"action":     "pop",
		"value":      value,
		"found":      ok,
		"stack_size": lockFreeStack.Size(),
	})
}

func queueEnqueueHandler(c *fiber.Ctx) error {
	value := c.Query("value", "item")
	lockFreeQueue.Enqueue(value)

	return c.JSON(fiber.Map{
		"action":     "enqueue",
		"value":      value,
		"queue_size": lockFreeQueue.Size(),
	})
}

func queueDequeueHandler(c *fiber.Ctx) error {
	value, ok := lockFreeQueue.Dequeue()

	return c.JSON(fiber.Map{
		"action":     "dequeue",
		"value":      value,
		"found":      ok,
		"queue_size": lockFreeQueue.Size(),
	})
}

func counterIncrHandler(c *fiber.Ctx) error {
	lfValue := lockFreeCounter.Increment()
	mValue := mutexCounter.Increment()

	return c.JSON(fiber.Map{
		"lock_free_counter": lfValue,
		"mutex_counter":     mValue,
	})
}

// Benchmark: Counter
func benchCounterHandler(c *fiber.Ctx) error {
	iterations := 100000
	goroutines := runtime.NumCPU()

	// Lock-free counter
	lfCounter := NewLockFreeCounter()
	start := time.Now()
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations/goroutines; i++ {
				lfCounter.Increment()
			}
		}()
	}
	wg.Wait()
	lockFreeDuration := time.Since(start)

	// Mutex counter
	mCounter := &MutexCounter{}
	start = time.Now()
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations/goroutines; i++ {
				mCounter.Increment()
			}
		}()
	}
	wg.Wait()
	mutexDuration := time.Since(start)

	return c.JSON(fiber.Map{
		"iterations":     iterations,
		"goroutines":     goroutines,
		"lock_free": fiber.Map{
			"duration": lockFreeDuration.String(),
			"value":    lfCounter.Value(),
		},
		"mutex": fiber.Map{
			"duration": mutexDuration.String(),
			"value":    mCounter.Value(),
		},
		"speedup": fmt.Sprintf("%.2fx", float64(mutexDuration)/float64(lockFreeDuration)),
	})
}

// Benchmark: Stack
func benchStackHandler(c *fiber.Ctx) error {
	iterations := 50000
	goroutines := runtime.NumCPU()

	// Lock-free stack
	lfStack := NewLockFreeStack()
	start := time.Now()
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations/goroutines; i++ {
				lfStack.Push(i)
				lfStack.Pop()
			}
		}()
	}
	wg.Wait()
	lockFreeDuration := time.Since(start)

	// Mutex stack
	mStack := NewMutexStack()
	start = time.Now()
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations/goroutines; i++ {
				mStack.Push(i)
				mStack.Pop()
			}
		}()
	}
	wg.Wait()
	mutexDuration := time.Since(start)

	return c.JSON(fiber.Map{
		"iterations": iterations,
		"goroutines": goroutines,
		"operations": "push + pop per iteration",
		"lock_free": fiber.Map{
			"duration": lockFreeDuration.String(),
		},
		"mutex": fiber.Map{
			"duration": mutexDuration.String(),
		},
		"speedup": fmt.Sprintf("%.2fx", float64(mutexDuration)/float64(lockFreeDuration)),
	})
}

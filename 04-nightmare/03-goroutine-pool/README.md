# 🔄 Goroutine Pool & Worker Patterns

## 🎯 วัตถุประสงค์
เรียนรู้การจัดการ **Goroutine Pool** เพื่อควบคุม resource usage และเพิ่ม throughput ในระบบ high-concurrency

## 💡 ปัญหาของ Unlimited Goroutines

```go
// ❌ BAD: Goroutine explosion
func handleRequests(requests <-chan Request) {
    for req := range requests {
        go func(r Request) { // New goroutine for each request
            processRequest(r)  // Can create millions of goroutines
        }(req)
    }
}

// Problems:
// - Memory overhead (2KB+ per goroutine)
// - Context switching overhead
// - Resource exhaustion
// - Unpredictable performance
```

## 🏊‍♂️ Worker Pool Patterns

### 1. **Basic Worker Pool**
```go
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    quit       chan bool
    wg         sync.WaitGroup
}

type Job func()

func NewWorkerPool(workers int, queueSize int) *WorkerPool {
    return &WorkerPool{
        workers:  workers,
        jobQueue: make(chan Job, queueSize),
        quit:     make(chan bool),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    for {
        select {
        case job := <-wp.jobQueue:
            job() // Execute job
        case <-wp.quit:
            return
        }
    }
}

func (wp *WorkerPool) Submit(job Job) {
    wp.jobQueue <- job
}

func (wp *WorkerPool) Stop() {
    close(wp.quit)
    wp.wg.Wait()
}
```

### 2. **Advanced Worker Pool with Context**
```go
type AdvancedPool struct {
    workers     int
    jobQueue    chan JobWithContext
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
    metrics     PoolMetrics
}

type JobWithContext struct {
    Task    func(context.Context) error
    OnDone  func(error)
    Timeout time.Duration
}

type PoolMetrics struct {
    JobsProcessed   int64
    JobsFailed      int64
    AvgProcessTime  time.Duration
    QueueSize       int64
}

func NewAdvancedPool(workers, queueSize int) *AdvancedPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &AdvancedPool{
        workers:  workers,
        jobQueue: make(chan JobWithContext, queueSize),
        ctx:      ctx,
        cancel:   cancel,
    }
}

func (ap *AdvancedPool) Start() {
    for i := 0; i < ap.workers; i++ {
        ap.wg.Add(1)
        go ap.worker(i)
    }
}

func (ap *AdvancedPool) worker(id int) {
    defer ap.wg.Done()
    
    for {
        select {
        case job := <-ap.jobQueue:
            start := time.Now()
            
            // Create job context with timeout
            jobCtx := ap.ctx
            if job.Timeout > 0 {
                var cancel context.CancelFunc
                jobCtx, cancel = context.WithTimeout(ap.ctx, job.Timeout)
                defer cancel()
            }
            
            // Execute job
            err := job.Task(jobCtx)
            
            // Update metrics
            atomic.AddInt64(&ap.metrics.JobsProcessed, 1)
            if err != nil {
                atomic.AddInt64(&ap.metrics.JobsFailed, 1)
            }
            
            // Record processing time
            duration := time.Since(start)
            // Update average (simplified)
            
            // Callback
            if job.OnDone != nil {
                job.OnDone(err)
            }
            
        case <-ap.ctx.Done():
            return
        }
    }
}

func (ap *AdvancedPool) Submit(task func(context.Context) error, timeout time.Duration) error {
    job := JobWithContext{
        Task:    task,
        Timeout: timeout,
    }
    
    select {
    case ap.jobQueue <- job:
        atomic.AddInt64(&ap.metrics.QueueSize, 1)
        return nil
    case <-ap.ctx.Done():
        return ap.ctx.Err()
    default:
        return errors.New("queue full")
    }
}
```

### 3. **Burst Worker Pool (Dynamic Scaling)**
```go
type BurstPool struct {
    minWorkers    int
    maxWorkers    int
    currentWorkers int64
    jobQueue      chan Job
    workerQueue   chan chan Job
    quit          chan bool
    mu            sync.RWMutex
}

func NewBurstPool(min, max, queueSize int) *BurstPool {
    return &BurstPool{
        minWorkers:  min,
        maxWorkers:  max,
        jobQueue:    make(chan Job, queueSize),
        workerQueue: make(chan chan Job, max),
        quit:        make(chan bool),
    }
}

func (bp *BurstPool) Start() {
    // Start minimum workers
    for i := 0; i < bp.minWorkers; i++ {
        bp.startWorker()
    }
    
    // Start dispatcher
    go bp.dispatcher()
}

func (bp *BurstPool) dispatcher() {
    for {
        select {
        case job := <-bp.jobQueue:
            // Try to assign to existing worker
            select {
            case workerJobQueue := <-bp.workerQueue:
                workerJobQueue <- job
            default:
                // No available worker, try to scale up
                if bp.canScaleUp() {
                    bp.startWorker()
                    // Retry assignment
                    workerJobQueue := <-bp.workerQueue
                    workerJobQueue <- job
                } else {
                    // Block until worker available
                    workerJobQueue := <-bp.workerQueue
                    workerJobQueue <- job
                }
            }
        case <-bp.quit:
            return
        }
    }
}

func (bp *BurstPool) canScaleUp() bool {
    current := atomic.LoadInt64(&bp.currentWorkers)
    return int(current) < bp.maxWorkers
}

func (bp *BurstPool) startWorker() {
    atomic.AddInt64(&bp.currentWorkers, 1)
    
    worker := &Worker{
        ID:          int(atomic.LoadInt64(&bp.currentWorkers)),
        JobQueue:    make(chan Job),
        WorkerQueue: bp.workerQueue,
        Quit:        make(chan bool),
    }
    
    worker.Start()
}
```

## 🔄 Specialized Pool Patterns

### 1. **Pipeline Worker Pool**
```go
type PipelineStage struct {
    Input  <-chan interface{}
    Output chan<- interface{}
    Worker func(interface{}) interface{}
}

type Pipeline struct {
    stages []PipelineStage
    pools  []*WorkerPool
}

func NewPipeline(stages []PipelineStage, workersPerStage int) *Pipeline {
    pools := make([]*WorkerPool, len(stages))
    
    for i, stage := range stages {
        pool := NewWorkerPool(workersPerStage, 100)
        pools[i] = pool
        
        // Start workers for this stage
        for j := 0; j < workersPerStage; j++ {
            go func(s PipelineStage) {
                for data := range s.Input {
                    result := s.Worker(data)
                    s.Output <- result
                }
            }(stage)
        }
    }
    
    return &Pipeline{stages: stages, pools: pools}
}
```

### 2. **Priority Worker Pool**
```go
type PriorityJob struct {
    Task     Job
    Priority int
    Created  time.Time
}

type PriorityPool struct {
    workers   int
    queues    []chan PriorityJob // Multiple priority queues
    quit      chan bool
    wg        sync.WaitGroup
}

func NewPriorityPool(workers int, priorities int) *PriorityPool {
    queues := make([]chan PriorityJob, priorities)
    for i := range queues {
        queues[i] = make(chan PriorityJob, 100)
    }
    
    return &PriorityPool{
        workers: workers,
        queues:  queues,
        quit:    make(chan bool),
    }
}

func (pp *PriorityPool) worker() {
    defer pp.wg.Done()
    
    for {
        // Check higher priority queues first
        for i := len(pp.queues) - 1; i >= 0; i-- {
            select {
            case job := <-pp.queues[i]:
                job.Task()
                continue
            default:
                // Continue to next priority level
            }
        }
        
        // If no high priority jobs, wait on any queue
        select {
        case <-pp.quit:
            return
        default:
            // Block on highest priority queue
            select {
            case job := <-pp.queues[len(pp.queues)-1]:
                job.Task()
            case <-pp.quit:
                return
            }
        }
    }
}
```

## 📊 Performance Optimization

### **Channel vs Mutex Comparison**
```go
// Channel-based coordination (slower but safer)
type ChannelPool struct {
    jobQueue chan Job
    workers  int
}

// Mutex-based coordination (faster but more complex)
type MutexPool struct {
    jobs    []Job
    mu      sync.Mutex
    cond    *sync.Cond
    workers int
}

func (mp *MutexPool) worker() {
    for {
        mp.mu.Lock()
        
        // Wait for job
        for len(mp.jobs) == 0 {
            mp.cond.Wait()
        }
        
        // Get job
        job := mp.jobs[0]
        mp.jobs = mp.jobs[1:]
        
        mp.mu.Unlock()
        
        // Execute
        job()
    }
}
```

### **Lock-free Work Stealing**
```go
type WorkStealingPool struct {
    workers []WorkStealingWorker
    global  chan Job
}

type WorkStealingWorker struct {
    id     int
    local  []Job // Local queue
    mu     sync.Mutex
    pool   *WorkStealingPool
}

func (w *WorkStealingWorker) run() {
    for {
        job := w.getJob()
        if job != nil {
            job()
        } else {
            // No work available, sleep briefly
            time.Sleep(time.Microsecond)
        }
    }
}

func (w *WorkStealingWorker) getJob() Job {
    // Try local queue first
    w.mu.Lock()
    if len(w.local) > 0 {
        job := w.local[len(w.local)-1]
        w.local = w.local[:len(w.local)-1]
        w.mu.Unlock()
        return job
    }
    w.mu.Unlock()
    
    // Try global queue
    select {
    case job := <-w.pool.global:
        return job
    default:
    }
    
    // Try stealing from other workers
    return w.steal()
}

func (w *WorkStealingWorker) steal() Job {
    for _, victim := range w.pool.workers {
        if victim.id == w.id {
            continue
        }
        
        victim.mu.Lock()
        if len(victim.local) > 1 {
            // Steal from the front (FIFO for stealing)
            job := victim.local[0]
            victim.local = victim.local[1:]
            victim.mu.Unlock()
            return job
        }
        victim.mu.Unlock()
    }
    return nil
}
```

## 🔍 Pool Sizing & Tuning

### **Optimal Pool Size Formula**
```go
func calculateOptimalPoolSize() int {
    // CPU-bound tasks
    cpuBound := runtime.NumCPU()
    
    // I/O-bound tasks  
    ioBound := runtime.NumCPU() * 2 // or more
    
    // Mixed workload
    mixed := int(float64(runtime.NumCPU()) * 1.5)
    
    // Consider memory constraints
    maxMemory := getAvailableMemory()
    goroutineMemory := 2 * 1024 // ~2KB per goroutine
    memoryLimit := maxMemory / goroutineMemory
    
    return min(mixed, memoryLimit)
}

// Dynamic sizing based on metrics
func (pool *AdaptivePool) adjustSize() {
    metrics := pool.getMetrics()
    
    if metrics.QueueDepth > 100 && metrics.CPUUsage < 80 {
        // Queue building up, CPU available -> scale up
        pool.scaleUp()
    } else if metrics.QueueDepth < 10 && metrics.CPUUsage > 90 {
        // Low queue, high CPU -> scale down
        pool.scaleDown()
    }
}
```

## 📈 Real-world Examples

### **HTTP Server Pool**
```go
type HTTPWorkerPool struct {
    pool     *WorkerPool
    handlers map[string]http.HandlerFunc
}

func (hwp *HTTPWorkerPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Submit request to worker pool
    hwp.pool.Submit(func() {
        if handler, exists := hwp.handlers[r.URL.Path]; exists {
            handler(w, r)
        } else {
            http.NotFound(w, r)
        }
    })
}
```

### **Database Connection Pool Integration**
```go
type DatabaseWorkerPool struct {
    workers int
    dbPool  *sql.DB
    jobQueue chan DatabaseJob
}

type DatabaseJob struct {
    Query    string
    Args     []interface{}
    Callback func(*sql.Rows, error)
}

func (dwp *DatabaseWorkerPool) worker() {
    for job := range dwp.jobQueue {
        rows, err := dwp.dbPool.Query(job.Query, job.Args...)
        job.Callback(rows, err)
    }
}
```

## 🎯 Best Practices

### **Do's ✅**
- Size pools based on workload type
- Monitor queue depth and processing time
- Use context for cancellation
- Implement graceful shutdown
- Pool long-running goroutines, not short ones

### **Don'ts ❌**
- Don't create unlimited goroutines
- Don't ignore backpressure
- Don't forget error handling
- Don't block workers unnecessarily
- Don't pool very short tasks (< 1ms)

## ⚠️ Common Pitfalls

1. **Goroutine Leaks**: Always ensure proper cleanup
2. **Deadlocks**: Avoid circular dependencies
3. **Resource Starvation**: Monitor queue sizes
4. **Context Propagation**: Pass context through jobs
5. **Error Propagation**: Handle errors appropriately

## 📊 Performance Metrics

```
Without Pool:     10,000 goroutines, 20MB memory, 50ms P99
With Pool (50):   50 goroutines, 2MB memory, 5ms P99
Improvement:      200x less goroutines, 10x less memory, 10x faster
``` 
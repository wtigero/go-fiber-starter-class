# 🏊 Memory Pool & Object Reuse

## 🎯 วัตถุประสงค์
เรียนรู้เทคนิค **Memory Pooling** เพื่อลด allocation overhead และ GC pressure ในระบบ high-throughput

## 💡 หลักการ Memory Pool

```go
// ❌ BAD: ใหม่ทุกครั้ง
func processRequest() {
    buffer := make([]byte, 1024) // Allocates every time
    // ... process with buffer
} // Buffer gets GC'd

// ✅ GOOD: Reuse from pool
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func processRequestOptimized() {
    buffer := bufferPool.Get().([]byte)
    defer bufferPool.Put(buffer[:0]) // Return cleared buffer
    
    // ... process with buffer
} // Buffer reused, no GC
```

## 🏊‍♂️ Types of Memory Pools

### 1. **sync.Pool (Built-in)**
```go
var stringBuilderPool = sync.Pool{
    New: func() interface{} {
        return &strings.Builder{}
    },
}

func optimizedStringProcess() string {
    builder := stringBuilderPool.Get().(*strings.Builder)
    defer func() {
        builder.Reset()
        stringBuilderPool.Put(builder)
    }()
    
    builder.WriteString("Hello")
    builder.WriteString(" World")
    return builder.String()
}
```

### 2. **Custom Ring Buffer Pool**
```go
type RingPool struct {
    buffers [][]byte
    index   int64
    mask    int64
}

func NewRingPool(size, bufferSize int) *RingPool {
    // Ensure size is power of 2
    poolSize := 1
    for poolSize < size {
        poolSize <<= 1
    }
    
    buffers := make([][]byte, poolSize)
    for i := range buffers {
        buffers[i] = make([]byte, bufferSize)
    }
    
    return &RingPool{
        buffers: buffers,
        mask:    int64(poolSize - 1),
    }
}

func (p *RingPool) Get() []byte {
    idx := atomic.AddInt64(&p.index, 1) & p.mask
    return p.buffers[idx][:0] // Reset length
}
```

### 3. **Sized Pool (เลือกขนาดได้)**
```go
type SizedPool struct {
    pools map[int]*sync.Pool
    sizes []int
}

func NewSizedPool(sizes []int) *SizedPool {
    pools := make(map[int]*sync.Pool)
    for _, size := range sizes {
        size := size // Capture for closure
        pools[size] = &sync.Pool{
            New: func() interface{} {
                return make([]byte, 0, size)
            },
        }
    }
    return &SizedPool{pools: pools, sizes: sizes}
}

func (p *SizedPool) Get(size int) []byte {
    // Find smallest buffer that fits
    for _, poolSize := range p.sizes {
        if poolSize >= size {
            return p.pools[poolSize].Get().([]byte)
        }
    }
    // Fallback: create new buffer
    return make([]byte, 0, size)
}

func (p *SizedPool) Put(buf []byte) {
    capacity := cap(buf)
    if pool, exists := p.pools[capacity]; exists {
        pool.Put(buf[:0]) // Reset length
    }
    // Ignore buffers that don't match pool sizes
}
```

## 🔧 Advanced Pool Patterns

### 1. **Object Pool สำหรับ Complex Structs**
```go
type HTTPResponse struct {
    Headers map[string]string
    Body    []byte
    Status  int
}

var responsePool = sync.Pool{
    New: func() interface{} {
        return &HTTPResponse{
            Headers: make(map[string]string),
            Body:    make([]byte, 0, 1024),
        }
    },
}

func (r *HTTPResponse) Reset() {
    // Clear map but keep capacity
    for k := range r.Headers {
        delete(r.Headers, k)
    }
    r.Body = r.Body[:0] // Reset slice length
    r.Status = 0
}

func handleRequest() *HTTPResponse {
    resp := responsePool.Get().(*HTTPResponse)
    resp.Reset()
    
    // Use response...
    resp.Status = 200
    resp.Headers["Content-Type"] = "application/json"
    
    return resp
}

func releaseResponse(resp *HTTPResponse) {
    responsePool.Put(resp)
}
```

### 2. **Worker Pool with Object Reuse**
```go
type Task struct {
    ID   int
    Data []byte
    done chan bool
}

type WorkerPool struct {
    workers   int
    taskPool  sync.Pool
    taskQueue chan *Task
}

func NewWorkerPool(workers int) *WorkerPool {
    wp := &WorkerPool{
        workers:   workers,
        taskQueue: make(chan *Task, workers*2),
        taskPool: sync.Pool{
            New: func() interface{} {
                return &Task{
                    Data: make([]byte, 0, 1024),
                    done: make(chan bool, 1),
                }
            },
        },
    }
    
    // Start workers
    for i := 0; i < workers; i++ {
        go wp.worker()
    }
    
    return wp
}

func (wp *WorkerPool) worker() {
    for task := range wp.taskQueue {
        // Process task...
        processTask(task)
        
        // Reset and return to pool
        task.Data = task.Data[:0]
        task.ID = 0
        task.done <- true
        wp.taskPool.Put(task)
    }
}

func (wp *WorkerPool) Submit(id int, data []byte) {
    task := wp.taskPool.Get().(*Task)
    task.ID = id
    task.Data = append(task.Data, data...)
    
    wp.taskQueue <- task
    <-task.done // Wait for completion
}
```

## 📊 Performance Comparison

### **Memory Allocation**
```go
// Without pool: 1000+ allocs/op, 64KB+ allocated
func withoutPool() {
    for i := 0; i < 1000; i++ {
        buffer := make([]byte, 1024)
        processBuffer(buffer)
    }
}

// With pool: 1-5 allocs/op, <1KB allocated
func withPool() {
    for i := 0; i < 1000; i++ {
        buffer := bufferPool.Get().([]byte)
        processBuffer(buffer)
        bufferPool.Put(buffer[:0])
    }
}
```

### **Benchmark Results**
```
BenchmarkWithoutPool-8    1000    1050000 ns/op    1048576 B/op   1000 allocs/op
BenchmarkWithPool-8      10000     105000 ns/op       1024 B/op      1 allocs/op
```

## ⚡ Real-world Examples

### **JSON Processing Pool**
```go
var jsonPool = sync.Pool{
    New: func() interface{} {
        return &jsonProcessor{
            encoder: json.NewEncoder(nil),
            decoder: json.NewDecoder(nil),
            buffer:  bytes.NewBuffer(make([]byte, 0, 1024)),
        }
    },
}

type jsonProcessor struct {
    encoder *json.Encoder
    decoder *json.Decoder
    buffer  *bytes.Buffer
}

func (j *jsonProcessor) Reset() {
    j.buffer.Reset()
    j.encoder = json.NewEncoder(j.buffer)
    j.decoder = json.NewDecoder(j.buffer)
}

func processJSON(data interface{}) ([]byte, error) {
    proc := jsonPool.Get().(*jsonProcessor)
    defer func() {
        proc.Reset()
        jsonPool.Put(proc)
    }()
    
    if err := proc.encoder.Encode(data); err != nil {
        return nil, err
    }
    
    return proc.buffer.Bytes(), nil
}
```

### **HTTP Client Pool**
```go
var httpClientPool = sync.Pool{
    New: func() interface{} {
        return &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     90 * time.Second,
            },
        }
    },
}

func makeRequest(url string) (*http.Response, error) {
    client := httpClientPool.Get().(*http.Client)
    defer httpClientPool.Put(client)
    
    return client.Get(url)
}
```

## 🔍 Pool Design Guidelines

### **Do's ✅**
- ใช้ `sync.Pool` สำหรับ short-lived objects
- Reset objects ก่อน return to pool
- Monitor pool hit rate
- Use different pools for different sizes
- Consider pool warming (pre-populate)

### **Don'ts ❌**
- Don't pool long-lived objects
- Don't assume objects stay in pool
- Don't store pointers to pooled data
- Avoid overly complex pool logic
- Don't pool tiny objects (< 64 bytes)

## 📈 Pool Monitoring

```go
type PoolStats struct {
    Gets    int64
    Puts    int64
    News    int64 // New allocations
    Reuses  int64 // Pool hits
}

type MonitoredPool struct {
    pool  sync.Pool
    stats PoolStats
}

func (p *MonitoredPool) Get() interface{} {
    atomic.AddInt64(&p.stats.Gets, 1)
    obj := p.pool.Get()
    
    if wasNew := isNewObject(obj); wasNew {
        atomic.AddInt64(&p.stats.News, 1)
    } else {
        atomic.AddInt64(&p.stats.Reuses, 1)
    }
    
    return obj
}

func (p *MonitoredPool) Put(obj interface{}) {
    atomic.AddInt64(&p.stats.Puts, 1)
    p.pool.Put(obj)
}

func (p *MonitoredPool) HitRate() float64 {
    gets := atomic.LoadInt64(&p.stats.Gets)
    if gets == 0 {
        return 0
    }
    reuses := atomic.LoadInt64(&p.stats.Reuses)
    return float64(reuses) / float64(gets)
}
```

## 🚨 Common Pitfalls

### **1. Pool Pollution**
```go
// ❌ BAD: Returning modified object
buffer := pool.Get().([]byte)
buffer = append(buffer, data...) // Modifies capacity
pool.Put(buffer) // Returns larger buffer

// ✅ GOOD: Reset before returning
buffer := pool.Get().([]byte)
buffer = append(buffer, data...)
// ... use buffer
pool.Put(buffer[:0]) // Reset length
```

### **2. Memory Leaks**
```go
// ❌ BAD: Storing references
type BadPool struct {
    objects []MyObject // Prevents GC
}

// ✅ GOOD: Let sync.Pool handle it
var goodPool = sync.Pool{
    New: func() interface{} {
        return &MyObject{}
    },
}
```

## 🎯 เมื่อไหร่ควรใช้?

✅ **ใช้เมื่อ:**
- High allocation rate (> 1000 allocs/sec)
- Temporary objects (short lifecycle)
- Fixed-size or predictable-size objects
- GC pressure เป็นปัญหา

❌ **ไม่ควรใช้เมื่อ:**
- Objects ถูกใช้นาน
- Allocation rate ต่ำ
- Object size เล็กมาก (< 64 bytes)
- Code complexity เพิ่มมากเกินไป 
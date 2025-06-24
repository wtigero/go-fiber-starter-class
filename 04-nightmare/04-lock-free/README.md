# 🔓 Lock-Free Programming

## 🎯 วัตถุประสงค์
เรียนรู้เทคนิค **Lock-Free Programming** ด้วย atomic operations เพื่อเพิ่ม performance และลด contention

## 💡 ปัญหาของ Mutex/Lock

```go
// ❌ BAD: Lock contention
type Counter struct {
    mu    sync.Mutex
    value int64
}

func (c *Counter) Add(delta int64) {
    c.mu.Lock()         // Lock contention
    c.value += delta    // Critical section
    c.mu.Unlock()      // Context switching
}

// Problems:
// - Context switching overhead
// - Lock contention at high concurrency
// - Potential deadlocks
// - Thread blocking
```

## 🚀 Atomic Operations

### 1. **Basic Atomic Counter**
```go
// ✅ GOOD: Lock-free counter
type AtomicCounter struct {
    value int64
}

func (c *AtomicCounter) Add(delta int64) int64 {
    return atomic.AddInt64(&c.value, delta)
}

func (c *AtomicCounter) Get() int64 {
    return atomic.LoadInt64(&c.value)
}

func (c *AtomicCounter) Set(val int64) {
    atomic.StoreInt64(&c.value, val)
}

// Compare-and-Swap
func (c *AtomicCounter) CompareAndSwap(old, new int64) bool {
    return atomic.CompareAndSwapInt64(&c.value, old, new)
}
```

### 2. **Lock-Free Stack**
```go
type LockFreeStack struct {
    head unsafe.Pointer
}

type node struct {
    data interface{}
    next unsafe.Pointer
}

func (s *LockFreeStack) Push(data interface{}) {
    newNode := &node{data: data}
    
    for {
        head := atomic.LoadPointer(&s.head)
        newNode.next = head
        
        if atomic.CompareAndSwapPointer(&s.head, head, unsafe.Pointer(newNode)) {
            break // Successfully pushed
        }
        // Retry if CAS failed (another goroutine modified head)
    }
}

func (s *LockFreeStack) Pop() interface{} {
    for {
        head := atomic.LoadPointer(&s.head)
        if head == nil {
            return nil // Empty stack
        }
        
        headNode := (*node)(head)
        next := atomic.LoadPointer(&headNode.next)
        
        if atomic.CompareAndSwapPointer(&s.head, head, next) {
            return headNode.data // Successfully popped
        }
        // Retry if CAS failed
    }
}
```

### 3. **Lock-Free Queue (Michael & Scott Algorithm)**
```go
type LockFreeQueue struct {
    head unsafe.Pointer
    tail unsafe.Pointer
}

type queueNode struct {
    data interface{}
    next unsafe.Pointer
}

func NewLockFreeQueue() *LockFreeQueue {
    dummy := &queueNode{}
    return &LockFreeQueue{
        head: unsafe.Pointer(dummy),
        tail: unsafe.Pointer(dummy),
    }
}

func (q *LockFreeQueue) Enqueue(data interface{}) {
    newNode := &queueNode{data: data}
    
    for {
        tail := atomic.LoadPointer(&q.tail)
        tailNode := (*queueNode)(tail)
        next := atomic.LoadPointer(&tailNode.next)
        
        if tail == atomic.LoadPointer(&q.tail) { // Tail hasn't changed
            if next == nil {
                // Tail is pointing to last node, try to link new node
                if atomic.CompareAndSwapPointer(&tailNode.next, next, unsafe.Pointer(newNode)) {
                    // Successfully linked, now try to move tail
                    atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(newNode))
                    break
                }
            } else {
                // Tail is not pointing to last node, try to advance tail
                atomic.CompareAndSwapPointer(&q.tail, tail, next)
            }
        }
    }
}

func (q *LockFreeQueue) Dequeue() interface{} {
    for {
        head := atomic.LoadPointer(&q.head)
        tail := atomic.LoadPointer(&q.tail)
        headNode := (*queueNode)(head)
        next := atomic.LoadPointer(&headNode.next)
        
        if head == atomic.LoadPointer(&q.head) { // Head hasn't changed
            if head == tail {
                if next == nil {
                    return nil // Queue is empty
                }
                // Tail is falling behind, try to advance it
                atomic.CompareAndSwapPointer(&q.tail, tail, next)
            } else {
                // Read data before CAS (prevents use-after-free)
                nextNode := (*queueNode)(next)
                data := nextNode.data
                
                // Try to move head to next node
                if atomic.CompareAndSwapPointer(&q.head, head, next) {
                    return data
                }
            }
        }
    }
}
```

## ⚡ Advanced Lock-Free Patterns

### 1. **Lock-Free Hash Map (simplified)**
```go
type LockFreeMap struct {
    buckets []unsafe.Pointer
    size    int
}

type bucket struct {
    key   string
    value interface{}
    next  unsafe.Pointer
}

func NewLockFreeMap(size int) *LockFreeMap {
    return &LockFreeMap{
        buckets: make([]unsafe.Pointer, size),
        size:    size,
    }
}

func (m *LockFreeMap) hash(key string) int {
    h := fnv.New32a()
    h.Write([]byte(key))
    return int(h.Sum32()) % m.size
}

func (m *LockFreeMap) Put(key string, value interface{}) {
    bucketIdx := m.hash(key)
    newBucket := &bucket{key: key, value: value}
    
    for {
        head := atomic.LoadPointer(&m.buckets[bucketIdx])
        
        // Check if key already exists
        current := head
        for current != nil {
            currentBucket := (*bucket)(current)
            if currentBucket.key == key {
                // Key exists, update value atomically
                // Note: This is simplified, real implementation needs more careful handling
                currentBucket.value = value
                return
            }
            current = atomic.LoadPointer(&currentBucket.next)
        }
        
        // Key doesn't exist, add new bucket
        newBucket.next = head
        if atomic.CompareAndSwapPointer(&m.buckets[bucketIdx], head, unsafe.Pointer(newBucket)) {
            break
        }
    }
}

func (m *LockFreeMap) Get(key string) (interface{}, bool) {
    bucketIdx := m.hash(key)
    head := atomic.LoadPointer(&m.buckets[bucketIdx])
    
    current := head
    for current != nil {
        currentBucket := (*bucket)(current)
        if currentBucket.key == key {
            return currentBucket.value, true
        }
        current = atomic.LoadPointer(&currentBucket.next)
    }
    
    return nil, false
}
```

### 2. **Wait-Free Ring Buffer**
```go
type WaitFreeRingBuffer struct {
    buffer []interface{}
    mask   int64
    head   int64  // Write position
    tail   int64  // Read position
}

func NewWaitFreeRingBuffer(size int) *WaitFreeRingBuffer {
    // Ensure size is power of 2
    if size&(size-1) != 0 {
        panic("size must be power of 2")
    }
    
    return &WaitFreeRingBuffer{
        buffer: make([]interface{}, size),
        mask:   int64(size - 1),
    }
}

func (rb *WaitFreeRingBuffer) Push(item interface{}) bool {
    head := atomic.LoadInt64(&rb.head)
    tail := atomic.LoadInt64(&rb.tail)
    
    // Check if buffer is full
    if head-tail >= int64(len(rb.buffer)) {
        return false
    }
    
    // Write to buffer
    rb.buffer[head&rb.mask] = item
    
    // Advance head
    atomic.StoreInt64(&rb.head, head+1)
    return true
}

func (rb *WaitFreeRingBuffer) Pop() (interface{}, bool) {
    head := atomic.LoadInt64(&rb.head)
    tail := atomic.LoadInt64(&rb.tail)
    
    // Check if buffer is empty
    if tail >= head {
        return nil, false
    }
    
    // Read from buffer
    item := rb.buffer[tail&rb.mask]
    
    // Advance tail
    atomic.StoreInt64(&rb.tail, tail+1)
    return item, true
}
```

## 🧠 Memory Ordering

### 1. **Acquire-Release Semantics**
```go
type Flag struct {
    ready int32
    data  int64
}

// Writer
func (f *Flag) SetData(value int64) {
    atomic.StoreInt64(&f.data, value)      // Store data first
    atomic.StoreInt32(&f.ready, 1)         // Release: signal ready
}

// Reader
func (f *Flag) GetData() (int64, bool) {
    if atomic.LoadInt32(&f.ready) == 1 {   // Acquire: check ready
        return atomic.LoadInt64(&f.data), true  // Load data after
    }
    return 0, false
}
```

### 2. **Memory Barriers**
```go
// Sequential consistency example
type SeqConsistent struct {
    x, y int32
}

func (s *SeqConsistent) Write() {
    atomic.StoreInt32(&s.x, 1)
    atomic.StoreInt32(&s.y, 1)
}

func (s *SeqConsistent) Read() (int32, int32) {
    y := atomic.LoadInt32(&s.y)
    x := atomic.LoadInt32(&s.x)
    return x, y
}
```

## 🔍 ABA Problem & Solutions

### **ABA Problem Example**
```go
// ❌ Susceptible to ABA problem
type VulnerableStack struct {
    head unsafe.Pointer
}

func (s *VulnerableStack) Pop() interface{} {
    for {
        head := atomic.LoadPointer(&s.head)
        if head == nil {
            return nil
        }
        
        headNode := (*node)(head)
        next := headNode.next
        
        // ABA problem: head might have changed to same value
        // but point to different node
        if atomic.CompareAndSwapPointer(&s.head, head, next) {
            return headNode.data
        }
    }
}
```

### **Solution: Tagged Pointers**
```go
type TaggedPointer struct {
    ptr unsafe.Pointer
    tag uint64
}

type SafeStack struct {
    head TaggedPointer
    mu   sync.Mutex // Only for tag update
}

func (s *SafeStack) Pop() interface{} {
    for {
        head := s.loadTaggedPointer()
        if head.ptr == nil {
            return nil
        }
        
        headNode := (*node)(head.ptr)
        next := headNode.next
        
        newHead := TaggedPointer{
            ptr: next,
            tag: head.tag + 1, // Increment tag
        }
        
        if s.compareAndSwapTaggedPointer(head, newHead) {
            return headNode.data
        }
    }
}
```

## 📊 Performance Comparison

```go
// Benchmark: Mutex vs Atomic
func BenchmarkMutexCounter(b *testing.B) {
    var counter MutexCounter
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            counter.Add(1)
        }
    })
}

func BenchmarkAtomicCounter(b *testing.B) {
    var counter AtomicCounter
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            counter.Add(1)
        }
    })
}

// Typical results (8 cores):
// BenchmarkMutexCounter-8     10000000    150 ns/op
// BenchmarkAtomicCounter-8    50000000     30 ns/op
// Atomic is 5x faster!
```

## 🛠️ Real-world Applications

### 1. **Metrics Collection**
```go
type Metrics struct {
    requests  int64
    errors    int64
    latency   int64
    count     int64
}

func (m *Metrics) RecordRequest(latency time.Duration, isError bool) {
    atomic.AddInt64(&m.requests, 1)
    atomic.AddInt64(&m.latency, int64(latency))
    atomic.AddInt64(&m.count, 1)
    
    if isError {
        atomic.AddInt64(&m.errors, 1)
    }
}

func (m *Metrics) GetAverageLatency() time.Duration {
    totalLatency := atomic.LoadInt64(&m.latency)
    count := atomic.LoadInt64(&m.count)
    if count == 0 {
        return 0
    }
    return time.Duration(totalLatency / count)
}
```

### 2. **Connection Pool**
```go
type ConnectionPool struct {
    available int64
    total     int64
}

func (cp *ConnectionPool) Acquire() bool {
    for {
        current := atomic.LoadInt64(&cp.available)
        if current <= 0 {
            return false // No connections available
        }
        
        if atomic.CompareAndSwapInt64(&cp.available, current, current-1) {
            return true
        }
        // Retry if CAS failed
    }
}

func (cp *ConnectionPool) Release() {
    atomic.AddInt64(&cp.available, 1)
}
```

## ⚠️ Lock-Free Pitfalls

### **1. ABA Problem**
- Value changes from A→B→A
- CAS thinks nothing changed
- Solution: Use tagged pointers or hazard pointers

### **2. Memory Reclamation**
- When is it safe to free memory?
- Solutions: Hazard pointers, RCU, epochs

### **3. Livelock**
- Threads keep retrying but make no progress
- Solution: Exponential backoff

### **4. Memory Ordering**
- Different architectures have different guarantees
- Use proper acquire/release semantics

## 🎯 เมื่อไหร่ควรใช้?

### ✅ **ใช้เมื่อ:**
- High contention scenarios
- Simple operations (counters, flags)
- Performance critical paths
- Low-latency requirements

### ❌ **ไม่ควรใช้เมื่อ:**
- Complex data structures
- Rare contention
- Debug/development phase
- Team lacks expertise

## 🔧 Tools & Debugging

```bash
# Race detection
go run -race program.go

# Memory profiling
go tool pprof mem.prof

# CPU profiling  
go tool pprof cpu.prof

# Contention profiling
go test -mutexprofile=mutex.prof
```

## 📈 Performance Guidelines

1. **Start simple**: Use channels and mutexes first
2. **Measure**: Profile before optimizing
3. **Target hotspots**: Focus on high-contention areas
4. **Test thoroughly**: Race conditions are subtle
5. **Keep it simple**: Complex lock-free code is bug-prone 
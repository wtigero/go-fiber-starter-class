# 🚀 Zero Allocation Techniques

## 🎯 วัตถุประสงค์
เรียนรู้เทคนิคการเขียน Go ที่ **ไม่ allocate memory** เพื่อลด GC pressure และเพิ่ม performance

## 💡 หลักการ Zero Allocation

```go
// ❌ BAD: Allocates memory
func BadConcat(strs []string) string {
    result := ""
    for _, s := range strs {
        result += s  // Each += allocates new string
    }
    return result
}

// ✅ GOOD: Zero allocation (with pre-allocated buffer)
func GoodConcat(strs []string, buf []byte) string {
    buf = buf[:0]  // Reset without allocation
    for _, s := range strs {
        buf = append(buf, s...)
    }
    return string(buf)  // Only one allocation
}
```

## 🧪 Techniques ต่างๆ

### 1. **String Builder Optimization**
```go
// ❌ String concatenation (multiple allocations)
func slowConcat(parts []string) string {
    result := ""
    for _, part := range parts {
        result += part
    }
    return result
}

// ✅ strings.Builder (minimal allocations)
func fastConcat(parts []string) string {
    var builder strings.Builder
    builder.Grow(calculateTotalSize(parts)) // Pre-allocate
    for _, part := range parts {
        builder.WriteString(part)
    }
    return builder.String()
}
```

### 2. **Slice Pre-allocation**
```go
// ❌ Growing slice (multiple allocations)
func slowSlice() []int {
    var result []int
    for i := 0; i < 1000; i++ {
        result = append(result, i) // Reallocates when capacity exceeded
    }
    return result
}

// ✅ Pre-allocated slice (single allocation)
func fastSlice() []int {
    result := make([]int, 0, 1000) // Pre-allocate capacity
    for i := 0; i < 1000; i++ {
        result = append(result, i)
    }
    return result
}
```

### 3. **Interface{} Avoidance**
```go
// ❌ interface{} causes boxing allocation
func slowProcess(data interface{}) {
    switch v := data.(type) {
    case int:
        processInt(v)
    case string:
        processString(v)
    }
}

// ✅ Type-specific functions (zero allocation)
func fastProcessInt(data int) {
    processInt(data)
}

func fastProcessString(data string) {
    processString(data)
}
```

### 4. **Buffer Reusing**
```go
// Buffer pool for reusing
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 1024)
    },
}

func processWithReuse(data []byte) []byte {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf[:0]) // Return cleared buffer
    
    // Use buf for processing
    buf = append(buf, data...)
    return append([]byte(nil), buf...) // Copy result
}
```

## 📊 Benchmarks

```go
func BenchmarkStringConcat(b *testing.B) {
    parts := []string{"hello", " ", "world", "!"}
    
    b.Run("Bad", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = slowConcat(parts)
        }
    })
    
    b.Run("Good", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = fastConcat(parts)
        }
    })
}

// Typical results:
// BenchmarkStringConcat/Bad-8    1000000   1050 ns/op   32 B/op   4 allocs/op
// BenchmarkStringConcat/Good-8   5000000    210 ns/op    8 B/op   1 allocs/op
```

## 🔍 Escape Analysis

```bash
# ดู escape analysis
go build -gcflags="-m" main.go

# ตัวอย่าง output:
# ./main.go:10:13: inlining call to fmt.Sprintf
# ./main.go:15:12: make([]byte, size) escapes to heap
# ./main.go:20:6: moved to heap: x
```

## ⚡ Advanced Techniques

### 1. **Stack vs Heap Allocation**
```go
// ✅ Stack allocation (fast)
func stackAlloc() {
    var arr [1024]byte  // On stack
    process(arr[:])
}

// ❌ Heap allocation (slower)
func heapAlloc() {
    arr := make([]byte, 1024)  // On heap
    process(arr)
}
```

### 2. **Inlining Optimization**
```go
//go:noinline
func notInlined(x int) int {
    return x * 2
}

// Small functions are auto-inlined
func inlined(x int) int {
    return x * 2
}
```

### 3. **Zero-copy Techniques**
```go
// ❌ Copy data
func withCopy(data []byte) string {
    return string(data)  // Copies data
}

// ✅ Zero-copy (unsafe but fast)
func zeroCopy(data []byte) string {
    return *(*string)(unsafe.Pointer(&data))
}
```

## 🧪 Real-world Examples

### **JSON Processing**
```go
// ✅ Using jsoniter (faster JSON library)
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ✅ Pre-allocated decoder
type OptimizedDecoder struct {
    decoder *jsoniter.Decoder
    buffer  bytes.Buffer
}

func (d *OptimizedDecoder) Decode(data []byte, v interface{}) error {
    d.buffer.Reset()
    d.buffer.Write(data)
    d.decoder.ResetBytes(data)
    return d.decoder.Decode(v)
}
```

### **HTTP Response**
```go
// ✅ Buffer pooling for HTTP responses
var responsePool = sync.Pool{
    New: func() interface{} {
        return &bytes.Buffer{}
    },
}

func optimizedHandler(w http.ResponseWriter, r *http.Request) {
    buf := responsePool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        responsePool.Put(buf)
    }()
    
    // Build response in buffer
    buf.WriteString(`{"status":"ok"}`)
    w.Write(buf.Bytes())
}
```

## 📈 Performance Impact

### **Memory Allocation Reduction**
- Before: 1000+ allocs/op
- After: 1-5 allocs/op
- **90-99% reduction**

### **GC Pressure**
- Before: 15% CPU on GC
- After: 1-2% CPU on GC
- **90% reduction**

### **Latency**
- Before: P99 = 50ms
- After: P99 = 5ms
- **10x improvement**

## 🔧 Tools

```bash
# Memory profiling
go test -memprofile=mem.prof
go tool pprof mem.prof

# Allocation tracking
go test -bench=. -benchmem

# Escape analysis
go build -gcflags="-m=2"
```

## ⚠️ Trade-offs

### **Pros**
- Dramatically reduced GC pressure
- Better memory efficiency
- Lower latency
- Higher throughput

### **Cons**
- More complex code
- Harder to maintain
- Potential memory leaks
- Unsafe operations risk

## 🎯 เมื่อไหร่ควรใช้?

✅ **ใช้เมื่อ:**
- High-frequency functions (hot path)
- Memory-constrained environments
- Latency-critical applications
- Profiling แสดง allocation hotspots

❌ **ไม่ควรใช้เมื่อ:**
- Cold path functions
- Development/prototype phase
- Team ไม่มีประสบการณ์
- Premature optimization 
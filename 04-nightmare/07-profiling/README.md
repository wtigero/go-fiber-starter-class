# 07. Profiling & Performance Analysis

## ระดับ: Expert

## สิ่งที่จะได้เรียนรู้

1. **pprof** - Go's built-in profiler
2. **CPU Profiling** - หา CPU bottlenecks
3. **Memory Profiling** - หา memory leaks
4. **Goroutine Profiling** - ตรวจสอบ goroutine leaks
5. **Trace** - Detailed execution tracing
6. **Benchmark** - การเขียน benchmarks ที่ดี

## ทำไมต้อง Profile?

> "Premature optimization is the root of all evil" - Donald Knuth

แต่เมื่อต้อง optimize จริงๆ ต้องรู้ว่า **จะ optimize ตรงไหน**

## Tools Overview

| Tool | ใช้ทำอะไร |
|------|-----------|
| `go tool pprof` | CPU & Memory profiling |
| `go tool trace` | Execution tracing |
| `go test -bench` | Benchmarking |
| `go build -gcflags="-m"` | Escape analysis |
| `go test -race` | Race condition detection |

## การใช้งาน pprof

### 1. เพิ่ม pprof endpoint

```go
import _ "net/http/pprof"

// ใน main()
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

### 2. เก็บ Profile

```bash
# CPU profile (30 วินาที)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Block profile (blocking calls)
go tool pprof http://localhost:6060/debug/pprof/block

# Mutex profile (lock contention)
go tool pprof http://localhost:6060/debug/pprof/mutex
```

### 3. วิเคราะห์ผลลัพธ์

```bash
# Text mode
(pprof) top 10
(pprof) list functionName

# Web UI
(pprof) web

# Flame graph
go tool pprof -http=:8080 profile.pb.gz
```

## Benchmark Best Practices

```go
func BenchmarkXxx(b *testing.B) {
    // Setup (ไม่นับเวลา)
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        // โค้ดที่ต้องการวัด
    }
}

// รันด้วย
// go test -bench=. -benchmem -benchtime=5s
```

## Memory Profiling Tips

### Escape Analysis
```bash
go build -gcflags="-m -m" 2>&1 | grep escape
```

### Common Causes of Allocations
1. Interface conversions
2. Closures capturing variables
3. Slice/map growth
4. String concatenation
5. Large structs passed by value

## โครงสร้างโปรเจค

```
07-profiling/
├── README.md
├── go.mod
├── main.go              # Web server with pprof
├── handlers/
│   └── slow.go          # Intentionally slow handlers
├── profiling/
│   ├── cpu_test.go      # CPU benchmarks
│   ├── mem_test.go      # Memory benchmarks
│   └── examples.go      # Problematic code examples
└── scripts/
    ├── profile_cpu.sh
    ├── profile_mem.sh
    └── analyze.sh
```

## การรัน

```bash
# Start server
go run main.go

# ในอีก terminal - สร้าง load
hey -n 10000 -c 100 http://localhost:3000/slow

# Profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

## Common Performance Issues

### 1. CPU Bound
- ลูปที่ไม่ optimize
- Regex compilation ในลูป
- JSON marshal/unmarshal

### 2. Memory Bound
- ไม่ใช้ sync.Pool
- String concatenation
- Unbounded caches

### 3. I/O Bound
- Blocking database calls
- ไม่ใช้ connection pooling
- Synchronous external API calls

### 4. Concurrency Issues
- Lock contention
- Goroutine leaks
- Channel blocking

## Optimization Workflow

```
1. Establish baseline (benchmark)
2. Profile under realistic load
3. Identify hotspots
4. Optimize ONE thing
5. Measure again
6. Repeat
```

## Real-World Example

### Before
```
BenchmarkHandler-8    1000   1500000 ns/op   50000 B/op   100 allocs/op
```

### After (using sync.Pool + buffer reuse)
```
BenchmarkHandler-8   50000     30000 ns/op     500 B/op     5 allocs/op
```

**Result: 50x faster, 100x less memory, 20x less allocations**

package main

import (
	"fmt"
	"runtime"
	"time"
)

// ❌ BAD: Growing slice causes multiple allocations
func slowSliceGrowth(n int) []int {
	var result []int // len=0, cap=0
	for i := 0; i < n; i++ {
		result = append(result, i) // Reallocates when cap exceeded
	}
	return result
}

// ✅ GOOD: Pre-allocated slice (single allocation)
func fastSlicePrealloc(n int) []int {
	result := make([]int, 0, n) // Pre-allocate capacity
	for i := 0; i < n; i++ {
		result = append(result, i) // No reallocation needed
	}
	return result
}

// 🚀 BETTER: Direct indexing (no append overhead)
func ultraFastSliceDirect(n int) []int {
	result := make([]int, n) // Pre-allocate with length
	for i := 0; i < n; i++ {
		result[i] = i // Direct assignment
	}
	return result
}

// ✅ Slice reusing pattern
type SlicePool struct {
	pool [][]int
}

func NewSlicePool() *SlicePool {
	return &SlicePool{
		pool: make([][]int, 0, 10),
	}
}

func (p *SlicePool) Get(capacity int) []int {
	if len(p.pool) > 0 {
		// Reuse existing slice
		slice := p.pool[len(p.pool)-1]
		p.pool = p.pool[:len(p.pool)-1]

		// Ensure capacity
		if cap(slice) < capacity {
			return make([]int, 0, capacity)
		}
		return slice[:0] // Reset length but keep capacity
	}
	return make([]int, 0, capacity)
}

func (p *SlicePool) Put(slice []int) {
	if cap(slice) > 0 {
		p.pool = append(p.pool, slice[:0])
	}
}

// 📊 Memory allocation tracker
func trackAllocs(name string, fn func()) {
	var start, end runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&start)

	fn()

	runtime.GC()
	runtime.ReadMemStats(&end)

	fmt.Printf("%s: %d allocs, %d bytes\n",
		name,
		end.Mallocs-start.Mallocs,
		end.TotalAlloc-start.TotalAlloc)
}

// 🔄 Slice filtering without allocation
func filterInPlace(slice []int, predicate func(int) bool) []int {
	writeIndex := 0
	for _, value := range slice {
		if predicate(value) {
			slice[writeIndex] = value
			writeIndex++
		}
	}
	return slice[:writeIndex] // Return filtered slice without allocation
}

// ❌ BAD: Creates new slice
func filterWithAlloc(slice []int, predicate func(int) bool) []int {
	var result []int
	for _, value := range slice {
		if predicate(value) {
			result = append(result, value) // Allocates
		}
	}
	return result
}

func benchmark(name string, fn func(), iterations int) {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		fn()
	}
	duration := time.Since(start)
	fmt.Printf("%s: %v (%v per op)\n", name, duration, duration/time.Duration(iterations))
}

func main() {
	const n = 1000
	const iterations = 10000

	fmt.Println("🚀 Slice Optimization Benchmarks")
	fmt.Println("=================================")

	// Benchmark slice growth
	benchmark("❌ Slow growth", func() {
		_ = slowSliceGrowth(n)
	}, iterations)

	benchmark("✅ Pre-allocated", func() {
		_ = fastSlicePrealloc(n)
	}, iterations)

	benchmark("🚀 Direct indexing", func() {
		_ = ultraFastSliceDirect(n)
	}, iterations)

	fmt.Println("\n🔄 Slice Pool Example")
	fmt.Println("====================")

	pool := NewSlicePool()

	// Use pool
	slice1 := pool.Get(100)
	for i := 0; i < 10; i++ {
		slice1 = append(slice1, i)
	}
	fmt.Printf("Used slice: %v\n", slice1)

	// Return to pool
	pool.Put(slice1)

	// Reuse from pool
	slice2 := pool.Get(100)
	fmt.Printf("Reused slice capacity: %d\n", cap(slice2))

	fmt.Println("\n🔍 Filtering Comparison")
	fmt.Println("======================")

	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	// Filter even numbers
	isEven := func(x int) bool { return x%2 == 0 }

	// In-place filtering (no allocation)
	dataCopy := make([]int, len(data))
	copy(dataCopy, data)
	filtered1 := filterInPlace(dataCopy, isEven)
	fmt.Printf("In-place filtered: %d items\n", len(filtered1))

	// With allocation
	filtered2 := filterWithAlloc(data, isEven)
	fmt.Printf("With alloc filtered: %d items\n", len(filtered2))
}

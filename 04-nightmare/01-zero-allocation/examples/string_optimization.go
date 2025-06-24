package main

import (
	"fmt"
	"strings"
	"time"
	"unsafe"
)

// ❌ BAD: Multiple allocations
func slowStringConcat(parts []string) string {
	result := ""
	for _, part := range parts {
		result += part // Each += creates new string
	}
	return result
}

// ✅ GOOD: Single allocation with strings.Builder
func fastStringConcat(parts []string) string {
	var builder strings.Builder

	// Pre-calculate total size to avoid reallocations
	totalSize := 0
	for _, part := range parts {
		totalSize += len(part)
	}
	builder.Grow(totalSize)

	for _, part := range parts {
		builder.WriteString(part)
	}
	return builder.String()
}

// ✅ BETTER: Reuse buffer (zero allocation after first call)
func ultraFastStringConcat(parts []string, buf *strings.Builder) string {
	buf.Reset() // Clear without deallocating

	for _, part := range parts {
		buf.WriteString(part)
	}
	return buf.String()
}

// 🚀 EXTREME: Zero-copy string conversion (unsafe but fastest)
func bytesToStringZeroCopy(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ✅ SAFE: Zero allocation byte slice to string
func bytesToString(b []byte) string {
	return string(b) // Go optimizes this in many cases
}

// 📊 Benchmark helper
func benchmark(name string, fn func(), iterations int) {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		fn()
	}
	duration := time.Since(start)
	fmt.Printf("%s: %v (%v per op)\n", name, duration, duration/time.Duration(iterations))
}

func main() {
	// Test data
	parts := []string{
		"The", " ", "quick", " ", "brown", " ", "fox", " ",
		"jumps", " ", "over", " ", "the", " ", "lazy", " ", "dog",
	}

	iterations := 100000

	fmt.Println("🚀 String Concatenation Benchmarks")
	fmt.Println("==================================")

	// Benchmark slow concat
	benchmark("❌ Slow concat", func() {
		_ = slowStringConcat(parts)
	}, iterations)

	// Benchmark fast concat
	benchmark("✅ Fast concat", func() {
		_ = fastStringConcat(parts)
	}, iterations)

	// Benchmark ultra fast concat with reused buffer
	var builder strings.Builder
	benchmark("🚀 Ultra fast", func() {
		_ = ultraFastStringConcat(parts, &builder)
	}, iterations)

	fmt.Println("\n💡 Zero-copy Examples")
	fmt.Println("====================")

	// Zero-copy example
	data := []byte("Hello, World!")

	// Normal conversion (may allocate)
	normal := string(data)
	fmt.Printf("Normal: %s\n", normal)

	// Zero-copy conversion (unsafe)
	zeroCopy := bytesToStringZeroCopy(data)
	fmt.Printf("Zero-copy: %s\n", zeroCopy)

	// ⚠️ WARNING: Modifying original data affects zero-copy string
	data[0] = 'h'
	fmt.Printf("After modification - Normal: %s, Zero-copy: %s\n", normal, zeroCopy)
}

package benchmarks

import (
	"strings"
	"testing"
	"unsafe"
)

// String concatenation benchmarks
func BenchmarkStringConcat(b *testing.B) {
	parts := []string{"hello", " ", "world", " ", "from", " ", "go"}

	b.Run("PlusOperator", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := ""
			for _, part := range parts {
				result += part
			}
			_ = result
		}
	})

	b.Run("StringBuilder", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			for _, part := range parts {
				builder.WriteString(part)
			}
			_ = builder.String()
		}
	})

	b.Run("StringBuilderPrealloc", func(b *testing.B) {
		b.ReportAllocs()
		totalSize := 0
		for _, part := range parts {
			totalSize += len(part)
		}

		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			builder.Grow(totalSize)
			for _, part := range parts {
				builder.WriteString(part)
			}
			_ = builder.String()
		}
	})

	b.Run("StringBuilderReuse", func(b *testing.B) {
		b.ReportAllocs()
		var builder strings.Builder
		totalSize := 0
		for _, part := range parts {
			totalSize += len(part)
		}
		builder.Grow(totalSize)

		for i := 0; i < b.N; i++ {
			builder.Reset()
			for _, part := range parts {
				builder.WriteString(part)
			}
			_ = builder.String()
		}
	})
}

// Slice allocation benchmarks
func BenchmarkSliceAllocation(b *testing.B) {
	const size = 1000

	b.Run("GrowingSlice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var slice []int
			for j := 0; j < size; j++ {
				slice = append(slice, j)
			}
		}
	})

	b.Run("PreallocatedSlice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			slice := make([]int, 0, size)
			for j := 0; j < size; j++ {
				slice = append(slice, j)
			}
		}
	})

	b.Run("DirectIndexing", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			slice := make([]int, size)
			for j := 0; j < size; j++ {
				slice[j] = j
			}
		}
	})
}

// Interface{} vs specific types
func BenchmarkInterfaceVsSpecific(b *testing.B) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.Run("InterfaceSum", func(b *testing.B) {
		b.ReportAllocs()
		interfaceData := make([]interface{}, len(data))
		for i, v := range data {
			interfaceData[i] = v
		}

		for i := 0; i < b.N; i++ {
			sum := 0
			for _, v := range interfaceData {
				sum += v.(int)
			}
		}
	})

	b.Run("SpecificSum", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sum := 0
			for _, v := range data {
				sum += v
			}
		}
	})
}

// Memory reuse patterns
func BenchmarkMemoryReuse(b *testing.B) {
	const size = 100

	b.Run("NewAllocation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buffer := make([]byte, size)
			for j := 0; j < size; j++ {
				buffer[j] = byte(j)
			}
		}
	})

	b.Run("ReuseBuffer", func(b *testing.B) {
		b.ReportAllocs()
		buffer := make([]byte, size)
		for i := 0; i < b.N; i++ {
			for j := 0; j < size; j++ {
				buffer[j] = byte(j)
			}
		}
	})
}

// String to bytes conversion
func BenchmarkStringToBytes(b *testing.B) {
	s := "Hello, World! This is a test string for benchmarking."

	b.Run("StandardConversion", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = []byte(s)
		}
	})

	b.Run("UnsafeConversion", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = stringToBytes(s)
		}
	})
}

// Zero-copy string to bytes (unsafe)
func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		cap int
	}{s, len(s)}))
}

// Map vs switch for small sets
func BenchmarkMapVsSwitch(b *testing.B) {
	// Map approach
	statusMap := map[int]string{
		200: "OK",
		404: "Not Found",
		500: "Internal Server Error",
		400: "Bad Request",
	}

	// Switch function
	statusSwitch := func(code int) string {
		switch code {
		case 200:
			return "OK"
		case 404:
			return "Not Found"
		case 500:
			return "Internal Server Error"
		case 400:
			return "Bad Request"
		default:
			return "Unknown"
		}
	}

	codes := []int{200, 404, 500, 400, 200, 404}

	b.Run("MapLookup", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, code := range codes {
				_ = statusMap[code]
			}
		}
	})

	b.Run("SwitchStatement", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, code := range codes {
				_ = statusSwitch(code)
			}
		}
	})
}

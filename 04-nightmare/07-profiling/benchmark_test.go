package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
)

// ============ String Concatenation Benchmarks ============

func BenchmarkStringConcat_Plus(b *testing.B) {
	// BAD: Using + operator
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < 1000; j++ {
			s += "x"
		}
		_ = s
	}
}

func BenchmarkStringConcat_Builder(b *testing.B) {
	// GOOD: Using strings.Builder
	for i := 0; i < b.N; i++ {
		var builder strings.Builder
		builder.Grow(1000)
		for j := 0; j < 1000; j++ {
			builder.WriteString("x")
		}
		_ = builder.String()
	}
}

func BenchmarkStringConcat_Buffer(b *testing.B) {
	// GOOD: Using bytes.Buffer
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		buf.Grow(1000)
		for j := 0; j < 1000; j++ {
			buf.WriteString("x")
		}
		_ = buf.String()
	}
}

// ============ Slice Append Benchmarks ============

func BenchmarkSlice_NoPrealloc(b *testing.B) {
	// BAD: No pre-allocation
	for i := 0; i < b.N; i++ {
		s := []int{}
		for j := 0; j < 10000; j++ {
			s = append(s, j)
		}
		_ = s
	}
}

func BenchmarkSlice_Prealloc(b *testing.B) {
	// GOOD: Pre-allocated
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 10000)
		for j := 0; j < 10000; j++ {
			s = append(s, j)
		}
		_ = s
	}
}

// ============ Object Pool Benchmarks ============

type TestObject struct {
	ID   int
	Name string
	Data []byte
}

var objectPool = sync.Pool{
	New: func() interface{} {
		return &TestObject{
			Data: make([]byte, 1024),
		}
	},
}

func BenchmarkObject_NewEachTime(b *testing.B) {
	// BAD: Allocate new object each time
	for i := 0; i < b.N; i++ {
		obj := &TestObject{
			ID:   i,
			Name: fmt.Sprintf("obj-%d", i),
			Data: make([]byte, 1024),
		}
		_ = obj
	}
}

func BenchmarkObject_Pool(b *testing.B) {
	// GOOD: Use object pool
	for i := 0; i < b.N; i++ {
		obj := objectPool.Get().(*TestObject)
		obj.ID = i
		obj.Name = fmt.Sprintf("obj-%d", i)
		// obj.Data already allocated
		objectPool.Put(obj)
	}
}

// ============ JSON Benchmarks ============

type JSONData struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Tags    []string `json:"tags"`
	Active  bool     `json:"active"`
	Balance float64  `json:"balance"`
}

var jsonPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func BenchmarkJSON_MarshalDirect(b *testing.B) {
	data := JSONData{
		ID: 1, Name: "Test", Email: "test@test.com",
		Tags: []string{"a", "b", "c"}, Active: true, Balance: 100.50,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, _ := json.Marshal(data)
		_ = result
	}
}

func BenchmarkJSON_EncoderPooled(b *testing.B) {
	data := JSONData{
		ID: 1, Name: "Test", Email: "test@test.com",
		Tags: []string{"a", "b", "c"}, Active: true, Balance: 100.50,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := jsonPool.Get().(*bytes.Buffer)
		buf.Reset()
		enc := json.NewEncoder(buf)
		enc.Encode(data)
		_ = buf.Bytes()
		jsonPool.Put(buf)
	}
}

// ============ Map vs Struct Benchmarks ============

func BenchmarkMap_Access(b *testing.B) {
	m := map[string]interface{}{
		"id":    1,
		"name":  "Test",
		"email": "test@test.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m["id"]
		_ = m["name"]
		_ = m["email"]
	}
}

func BenchmarkStruct_Access(b *testing.B) {
	s := struct {
		ID    int
		Name  string
		Email string
	}{1, "Test", "test@test.com"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.ID
		_ = s.Name
		_ = s.Email
	}
}

// ============ Interface vs Concrete Type ============

func processInterface(v interface{}) int {
	return v.(int) + 1
}

func processConcrete(v int) int {
	return v + 1
}

func BenchmarkInterface_Call(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = processInterface(i)
	}
}

func BenchmarkConcrete_Call(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = processConcrete(i)
	}
}

// ============ Run with: go test -bench=. -benchmem ============

package main

import (
	"fmt"
	"time"
)

// ❌ BAD: interface{} causes boxing allocation
func slowProcessAny(data interface{}) string {
	switch v := data.(type) {
	case int:
		return fmt.Sprintf("int: %d", v) // Boxing allocation
	case string:
		return fmt.Sprintf("string: %s", v)
	case float64:
		return fmt.Sprintf("float: %.2f", v)
	default:
		return "unknown"
	}
}

// ✅ GOOD: Type-specific functions (zero allocation)
func fastProcessInt(data int) string {
	// Pre-allocated buffer approach or direct return
	return "int: " + itoa(data) // Custom itoa to avoid allocation
}

func fastProcessString(data string) string {
	return "string: " + data
}

func fastProcessFloat(data float64) string {
	return "float: " + ftoa(data) // Custom ftoa
}

// 🚀 EXTREME: Union type simulation (zero allocation)
type Value struct {
	typ  uint8   // Type indicator
	iVal int64   // Integer value
	fVal float64 // Float value
	sVal string  // String value
}

const (
	TypeInt = iota
	TypeFloat
	TypeString
)

func NewIntValue(v int64) Value {
	return Value{typ: TypeInt, iVal: v}
}

func NewFloatValue(v float64) Value {
	return Value{typ: TypeFloat, fVal: v}
}

func NewStringValue(v string) Value {
	return Value{typ: TypeString, sVal: v}
}

func (v Value) Process() string {
	switch v.typ {
	case TypeInt:
		return "int: " + itoa(int(v.iVal))
	case TypeFloat:
		return "float: " + ftoa(v.fVal)
	case TypeString:
		return "string: " + v.sVal
	default:
		return "unknown"
	}
}

// 🔧 Custom number to string (avoiding fmt allocations)
func itoa(i int) string {
	if i == 0 {
		return "0"
	}

	var buf [20]byte // Stack allocated buffer
	pos := len(buf)
	negative := i < 0
	if negative {
		i = -i
	}

	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}

	if negative {
		pos--
		buf[pos] = '-'
	}

	return string(buf[pos:])
}

func ftoa(f float64) string {
	// Simplified float to string (production would use strconv)
	i := int(f)
	frac := int((f - float64(i)) * 100)
	return itoa(i) + "." + itoa(frac)
}

// 📊 Generic vs Specific comparison
// ❌ Generic function with interface{}
func genericSum(values []interface{}) interface{} {
	var intSum int
	var floatSum float64
	hasFloat := false

	for _, v := range values {
		switch val := v.(type) {
		case int:
			intSum += val
		case float64:
			floatSum += val
			hasFloat = true
		}
	}

	if hasFloat {
		return float64(intSum) + floatSum
	}
	return intSum
}

// ✅ Specific functions (zero allocation)
func sumInts(values []int) int {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum
}

func sumFloats(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum
}

// 🔄 Polymorphism without interface{} - using type switch on concrete types
type ProcessorFunc func(string) string

func getProcessor(dataType string) ProcessorFunc {
	switch dataType {
	case "upper":
		return func(s string) string {
			// Inline upper case conversion (avoiding strings.ToUpper allocation)
			buf := make([]byte, len(s))
			for i, b := range []byte(s) {
				if b >= 'a' && b <= 'z' {
					buf[i] = b - 32
				} else {
					buf[i] = b
				}
			}
			return string(buf)
		}
	case "prefix":
		return func(s string) string { return "prefix_" + s }
	default:
		return func(s string) string { return s }
	}
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
	const iterations = 100000

	fmt.Println("🚀 Interface{} Avoidance Benchmarks")
	fmt.Println("===================================")

	// Test data
	testInt := 42
	testString := "hello"
	testFloat := 3.14

	// Benchmark interface{} vs specific
	benchmark("❌ With interface{}", func() {
		_ = slowProcessAny(testInt)
		_ = slowProcessAny(testString)
		_ = slowProcessAny(testFloat)
	}, iterations)

	benchmark("✅ Type specific", func() {
		_ = fastProcessInt(testInt)
		_ = fastProcessString(testString)
		_ = fastProcessFloat(testFloat)
	}, iterations)

	// Union type approach
	intVal := NewIntValue(42)
	strVal := NewStringValue("hello")
	floatVal := NewFloatValue(3.14)

	benchmark("🚀 Union type", func() {
		_ = intVal.Process()
		_ = strVal.Process()
		_ = floatVal.Process()
	}, iterations)

	fmt.Println("\n📊 Sum Comparison")
	fmt.Println("=================")

	// Generic vs specific sum
	genericData := []interface{}{1, 2, 3, 4, 5}
	specificData := []int{1, 2, 3, 4, 5}

	benchmark("❌ Generic sum", func() {
		_ = genericSum(genericData)
	}, iterations)

	benchmark("✅ Specific sum", func() {
		_ = sumInts(specificData)
	}, iterations)

	fmt.Println("\n🔄 Polymorphism Without Interface{}")
	fmt.Println("==================================")

	// Function-based polymorphism
	processor := getProcessor("upper")
	result := processor("hello world")
	fmt.Printf("Processed: %s\n", result)

	fmt.Println("\n💡 Custom Number Conversion")
	fmt.Println("===========================")

	// Compare custom vs fmt
	testNum := 12345

	benchmark("❌ fmt.Sprintf", func() {
		_ = fmt.Sprintf("%d", testNum)
	}, iterations)

	benchmark("✅ Custom itoa", func() {
		_ = itoa(testNum)
	}, iterations)

	// Show results
	fmt.Printf("fmt result: %s\n", fmt.Sprintf("%d", testNum))
	fmt.Printf("custom result: %s\n", itoa(testNum))
}

package examples

import (
	"fmt"
	"strings"
	"testing"
)

// BenchmarkStringConcatenation benchmarks string concatenation with +
func BenchmarkStringConcatenation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < 100; j++ {
			s += "a"
		}
	}
}

// BenchmarkStringBuilder benchmarks string concatenation with strings.Builder
func BenchmarkStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		for j := 0; j < 100; j++ {
			sb.WriteString("a")
		}
		_ = sb.String()
	}
}

// BenchmarkStringJoin benchmarks string concatenation with strings.Join
func BenchmarkStringJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parts := make([]string, 100)
		for j := 0; j < 100; j++ {
			parts[j] = "a"
		}
		_ = strings.Join(parts, "")
	}
}

// BenchmarkSprintf benchmarks string formatting with fmt.Sprintf
func BenchmarkSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("Hello %s, you are %d years old", "World", 42)
	}
}

// BenchmarkStringFormat benchmarks string formatting with concatenation
func BenchmarkStringFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		name := "World"
		age := 42
		_ = "Hello " + name + ", you are " + fmt.Sprint(age) + " years old"
	}
}

package examples

import (
	"testing"
)

// BenchmarkSliceAppend benchmarks appending to a slice without pre-allocation
func BenchmarkSliceAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := []int{}
		for j := 0; j < 1000; j++ {
			s = append(s, j)
		}
	}
}

// BenchmarkSliceAppendPrealloc benchmarks appending to a pre-allocated slice
func BenchmarkSliceAppendPrealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 1000)
		for j := 0; j < 1000; j++ {
			s = append(s, j)
		}
	}
}

// BenchmarkSliceCopy benchmarks copying slices
func BenchmarkSliceCopy(b *testing.B) {
	src := make([]int, 1000)
	for i := range src {
		src[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := make([]int, len(src))
		copy(dst, src)
	}
}

// BenchmarkMapAccess benchmarks map access
func BenchmarkMapAccess(b *testing.B) {
	m := make(map[int]int)
	for i := 0; i < 1000; i++ {
		m[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[500]
	}
}

// BenchmarkMapIteration benchmarks iterating over a map
func BenchmarkMapIteration(b *testing.B) {
	m := make(map[int]int)
	for i := 0; i < 1000; i++ {
		m[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for _, v := range m {
			sum += v
		}
	}
}

package benchmark

import (
	"testing"
)

func BenchmarkMakeSliceWithoutAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MakeSliceWithoutAlloc()
	}
}

func BenchmarkMakeSliceWithAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MakeSliceWithAlloc()
	}
}

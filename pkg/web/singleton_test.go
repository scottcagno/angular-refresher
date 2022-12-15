package web

import (
	"testing"
	"unsafe"
)

type T1A struct{ a int }
type T2A struct{ a int }

func TestSingleton(t *testing.T) {
	const N = 100

	v1 := GetSingleton[T1A]()
	v2 := GetSingleton[T2A]()

	t.Logf("v1=(%#v, %p) v2=(%#v %p)", v1, &v1, v2, &v2)

	if v1 == nil || v2 == nil {
		t.Fatalf("v1 or v2 should not be nil")
	}

	if unsafe.Pointer(&v1) == unsafe.Pointer(&v2) {
		t.Fatalf("v1 and v2 should not be the same")
	}

	for i := 0; i < N; i++ {
		v := GetSingleton[T1A]()

		if v != v1 {
			t.Fatalf("all singleton should be the same")
		}
	}

	for i := 0; i < N; i++ {
		v := GetSingleton[T2A]()

		if v != v2 {
			t.Fatalf("all singleton should be the same")
		}
	}
}

func BenchmarkSingleton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetSingleton[T1A]()
	}
}

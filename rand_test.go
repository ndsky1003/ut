package rand

import (
	"fmt"
	"testing"
)

func BenchmarkRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := GenID(6, func(id int) bool { return true })
		if len(fmt.Sprintf("%d", v)) != 6 {
			b.Error("dd")
		}
	}
}

func TestPick(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < 100; i++ {
		pick := i % 9
		v := Pick[int](arr, pick)
		t.Log(v)
	}
	t.Error("333")
}

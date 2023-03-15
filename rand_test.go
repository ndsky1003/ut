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

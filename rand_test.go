package rand

import (
	"testing"
)

func TestGenID(t *testing.T) {
	type args struct {
		num_width uint8
		canUse    func(id int) bool
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 int
	}{
		{
			name: "test1",
			args: func(t *testing.T) args {
				return args{
					num_width: 8,
					canUse: func(id int) bool {
						return true
					},
				}
			},
			want1: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			got1 := GenID(tArgs.num_width, tArgs.canUse)
			t.Log(got1)
		})
	}
}

func BenchmarkGenID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenID(8, func(id int) bool {
			return true
		})
	}
}

// 本来生成范围就是1->20来位，所以模糊测试一定失败
func FuzzGenID(f *testing.F) {
	testcases := []uint8{1, 1, 2, 3, 4, 5, 6}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig uint8) {
		t.Log(orig)
		got1 := GenID(orig, func(id int) bool {
			return true
		})
		t.Logf("out:%v", got1)
	})
}

package ut

import (
	"fmt"
	"math/rand/v2"
)

var metas = [max_D_width + 1]rangeMeta{}

type rangeMeta struct {
	min, max uint64
}

const (
	max_D_width uint8  = 19                  //支持的最大数字宽度
	min_D_width uint8  = 1                   //支持的最小数字宽度
	min_num     uint64 = 1000000000000000000 //保证编译环境能正常运行
	max_num     uint64 = 9999999999999999999
)

func init() {
	for i := min_D_width; i <= max_D_width; i++ {
		min, max := gen_min_max_by_width(i)
		metas[i] = rangeMeta{min: min, max: max}
	}
}

func gen_min_max_by_width(width uint8) (min uint64, max uint64) {
	left := max_D_width - width
	var base uint64 = 1
	var i uint8
	for ; i < left; i++ {
		base *= 10
	}
	min = min_num / base
	max = max_num / base
	return
}

// num_width 需要的id宽度（十进制位数，1~19）
// canUse 判断该ID是否可用（例如未被占用）；返回 0 表示重试上限内未找到可用ID
func GenID(num_width uint8, canUse func(id uint64) bool) uint64 {
	if num_width > max_D_width || num_width < min_D_width {
		panic(fmt.Errorf("num_width:%d is invalid", num_width))
	}

	meta := metas[num_width]
	if canUse == nil {
		panic("canUse func is nil")
	}

	for {
		// 在 [min, max] 上均匀取值，保证 ID 分布均匀（首位数字无偏斜）
		id := meta.min + rand.Uint64N(meta.max-meta.min+1)
		if canUse(id) {
			return id
		}
	}
}

func Pick[T any](origins []T, count int) []T {
	if count == 0 {
		return []T{}
	}
	length := len(origins)
	newOrigins := make([]T, length)
	copy(newOrigins, origins)
	origins = newOrigins
	if length <= count {
		return origins
	} else {
		results := make([]T, count)
		for i := range count {
			ri := rand.IntN(length - i)
			results[i] = origins[ri]
			origins = append(origins[:ri], origins[ri+1:]...)
		}
		return results
	}
}

// Peek 从 origins 中随机取 count 个不重复元素（不放回抽样），不修改原切片。
// 语义与 Pick 一致。
func Peek[T any](origins []T, count int) []T {
	if count == 0 {
		return []T{}
	}
	length := len(origins)
	pool := append([]T(nil), origins...)
	if length <= count {
		return pool
	}
	results := make([]T, count)
	for i := range count {
		ri := rand.IntN(length - i)
		results[i] = pool[ri]
		pool = append(pool[:ri], pool[ri+1:]...)
	}
	return results
}

func Shuffle[T any](arr []T, opts ...*shuffle_option) {
	opt := ShuffleOption().SetCount(1).Merge(opts...)
	shuffleCount := opt.GetCount()
	length := len(arr)
	for range shuffleCount {
		rand.Shuffle(length, func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
}

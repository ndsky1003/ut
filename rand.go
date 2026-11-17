package ut

import (
	"fmt"
	"math/rand/v2"
)

var metas = map[uint8]*flag_items{}

type flag_items struct {
	min, max                           uint64
	flag_arr                           []uint8 //为了随机宽度
	min_binary_width, max_binary_width uint8
}

func (this *flag_items) String() string {
	return fmt.Sprintf("min:%d,max:%d,flag_arr:%v,min_binary_width:%d,max_binary_width:%d", this.min, this.max, this.flag_arr, this.min_binary_width, this.max_binary_width)
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
		item := &flag_items{
			min:              min,
			max:              max,
			min_binary_width: uint8(len(fmt.Sprintf("%b", min))),
			max_binary_width: uint8(len(fmt.Sprintf("%b", max))),
		}
		length := item.max_binary_width - item.min_binary_width + 1
		item.flag_arr = make([]uint8, length)
		var j uint8
		for ; j < length; j++ {
			item.flag_arr[j] = j
		}
		metas[i] = item
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

func shuffle(arr []uint8) []uint8 {
	for i := len(arr) - 1; i >= 1; i-- {
		ridx := rand.IntN(i + 1)
		arr[i], arr[ridx] = arr[ridx], arr[i]
	}
	return arr
}

// num_width 需要的id宽度
// canUse 是否可以使用的这个ID
func GenID(num_width uint8, canUse func(id uint64) bool) uint64 {
	if num_width > max_D_width || num_width < min_D_width {
		panic(fmt.Errorf("num_width:%d is invalid", num_width))
	}

	meta, ok := metas[num_width]
	if !ok {
		panic("meta is not support")
	}
	if len(meta.flag_arr) == 0 {
		panic(fmt.Errorf("meta flag arr:%v is invalid", meta.flag_arr))
	}

	if canUse == nil {
		panic("canUse func is nil")
	}
	// randI := rand.Uint64()
	var randI uint64
	for randI = rand.Uint64(); randI < meta.min; randI = rand.Uint64() {
	}
	flags := make([]uint8, len(meta.flag_arr))
	copy(flags, meta.flag_arr)
	flags = shuffle(flags)
	for _, flag := range flags {
		width := meta.min_binary_width + flag
		var leftOffset uint8 = 64 - width
		var i uint8
		for ; i <= leftOffset; i++ {
			tmpwidth := width + i
			// result := (randI & (1<<tmpwidth - 1)) >> i
			mask := ^uint64(0) >> (64 - tmpwidth)
			result := (randI & mask) >> i
			if result >= meta.min && result <= meta.max && canUse(result) {
				return result
			}
		}
	}
	return GenID(num_width, canUse)
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
		for i := 0; i < count; i++ {
			ri := rand.IntN(length - i)
			results[i] = origins[ri]
			origins = append(origins[:ri], origins[ri+1:]...)
		}
		return results
	}
}

func Shuffle[T any](arr []T, opts ...*shuffle_option) {
	opt := ShuffleOption().Merge(opts...)
	shuffleCount := opt.GetCount()
	length := len(arr)
	var c uint8
	for c = 0; c < shuffleCount; c++ {
		rand.Shuffle(length, func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	}
}

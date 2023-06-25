package rand

import (
	"fmt"
	"math/rand"
	"time"
)

var r rand.Source
var metas = map[uint8]*flag_items{}

const (
	max_D_width uint8 = 18                 //支持的最大数字宽度
	min_D_width uint8 = 1                  //支持的最小数字宽度
	min_num     int64 = 100000000000000000 //保证编译环境能正常运行
	max_num     int64 = 999999999999999999
)

func init() {
	v := time.Now().UnixNano()
	rand.Seed(v)
	r = rand.NewSource(v)
	for i := min_D_width; i <= max_D_width; i++ {
		min, max := gen_min_max_by_width(i)
		item := &flag_items{
			min:              min,
			max:              max,
			min_binary_width: uint8(len(fmt.Sprintf("%b", min))),
			max_binary_width: uint8(len(fmt.Sprintf("%b", max))),
		}
		length := item.max_binary_width - item.min_binary_width + 1
		item.flag_arr = make([]uint8, length, length)
		var j uint8
		for ; j < length; j++ {
			item.flag_arr[j] = j
		}
		metas[i] = item
	}
}

func gen_min_max_by_width(width uint8) (min int64, max int64) {
	left := max_D_width - width
	var base int64 = 1
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
		var ridx = rand.Intn(i)
		arr[i], arr[ridx] = arr[ridx], arr[i]
	}
	return arr
}

type flag_items struct {
	min, max                           int64
	flag_arr                           []uint8 //为了随机宽度
	min_binary_width, max_binary_width uint8
}

func (this *flag_items) String() string {
	return fmt.Sprintf("min:%d,max:%d,flag_arr:%v,min_binary_width:%d,max_binary_width:%d", this.min, this.max, this.flag_arr, this.min_binary_width, this.max_binary_width)
}

// num_width 需要的id宽度
// canUse 是否可以使用的这个ID
func GenID(num_width uint8, canUse func(id int) bool) int {
	randI := r.Int63()
	if num_width > max_D_width || num_width < min_D_width {
		panic(fmt.Sprintf("num_width:%d is invalid", num_width))
	}

	meta, ok := metas[num_width]
	if !ok {
		panic("meta is not support")
	}
	if len(meta.flag_arr) == 0 {
		panic(fmt.Sprintf("meta flag arr:%v is invalid", meta.flag_arr))
	}

	flags := make([]uint8, len(meta.flag_arr))
	copy(flags, meta.flag_arr)
	flags = shuffle(flags)
	width := meta.min_binary_width
	for _, flag := range flags {
		width += flag
		var leftOffset uint8 = 63 - width
		var i uint8
		for ; i <= leftOffset; i++ {
			tmpwidth := width + i
			result := (randI & (1<<tmpwidth - 1)) >> i
			if result > meta.min && result < meta.max && canUse(int(result)) {
				return int(result)
			}
		}
		return GenID(num_width, canUse)
	}
	return 0
}

func Pick[T any](origins []T, count int) []T {
	if count == 0 {
		return []T{}
	}
	length := len(origins)
	newOrigins := make([]T, length, length)
	copy(newOrigins, origins)
	origins = newOrigins
	if length <= count {
		return origins
	} else {
		results := make([]T, count)
		for i := 0; i < count; i++ {
			ri := rand.Intn(length - i)
			results[i] = origins[ri]
			origins = append(origins[:ri], origins[ri+1:]...)
		}
		return results
	}
}

// 洗牌
func Shuffle[T any](arr []T, shuffleCount uint8) {
	length := len(arr)
	var j uint8
	for j = 0; j < shuffleCount; j++ {
		for i := 0; i < length; i++ {
			r := rand.Intn(length)
			arr[i], arr[r] = arr[r], arr[i]
		}
	}
}

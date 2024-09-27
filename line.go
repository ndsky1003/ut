package ut

import (
	"errors"
	"math/rand/v2"
	"sync"
	"sync/atomic"
)

type Line[T any] struct {
	l            sync.Mutex
	index        int64
	arr          []T
	length       int
	shuffle_func func([]T) error
}

func NewLine[T any](opts ...*line_option[T]) *Line[T] {
	ch := make(chan struct{})
	close(ch)
	l := &Line[T]{}
	opt := LineOption[T]().Merge(opts...)
	if opt.shuffle_func != nil {
		l.shuffle_func = opt.shuffle_func
	} else {
		l.shuffle_func = func(arr []T) error {
			rand.Shuffle(len(arr), func(i, j int) {
				arr[i], arr[j] = arr[j], arr[i]
			})
			return nil
		}
	}
	return l
}

// 在外头洗牌,避免加锁,浪费卡出时间
func (this *Line[T]) SetArr(index_arr []T) error {
	if this.arr == nil || len(this.arr) == 0 {
		return errors.New("arr is nil")
	}
	this.l.Lock()
	defer this.l.Unlock()
	atomic.SwapInt64(&this.index, 0)
	this.arr = index_arr
	this.length = len(this.arr)
	return nil
}

func (this *Line[T]) get_value(index int) ([]T, T) {
	this.l.Lock()
	defer this.l.Unlock()
	index = index % this.length
	return this.arr, this.arr[index]
}

func (this *Line[T]) Step() (T, error) {
	newIndex := atomic.AddInt64(&this.index, 1)
	oldIndex := newIndex - 1
	arr, r := this.get_value(int(oldIndex))
	if int(newIndex) == this.length {
		tmp_arr := make([]T, len(arr))
		copy(tmp_arr, arr)
		go this.shuffle(tmp_arr)
	}
	return r, nil
}

func (this *Line[T]) shuffle(arr []T) {
	if this.shuffle_func != nil {
		this.shuffle_func(arr)
		this.SetArr(arr)
	}
	atomic.SwapInt64(&this.index, 0)
}

package ut

import "time"

type shuffle_option struct {
	count *uint8
}

func ShuffleOption() *shuffle_option {
	return &shuffle_option{}
}

func (this *shuffle_option) SetCount(count uint8) {
	this.count = &count
}

func (this *shuffle_option) GetCount() uint8 {
	if this.count == nil {
		return 1
	}
	return *this.count
}

func (this *shuffle_option) merge(delta *shuffle_option) *shuffle_option {
	if delta == nil {
		return this
	}
	if delta.count != nil {
		this.count = delta.count
	}
	return this
}

func (this *shuffle_option) Merge(opts ...*shuffle_option) *shuffle_option {
	for _, option := range opts {
		this.merge(option)
	}
	return this
}

// =============================

type protect_run_option struct {
	interval *time.Duration
}

func ProtectRunOption() *protect_run_option {
	return &protect_run_option{}
}

func (this *protect_run_option) SetInterval(interval time.Duration) {
	this.interval = &interval
}

func (this *protect_run_option) GetInterval() time.Duration {
	if this.interval == nil {
		return time.Millisecond
	}
	return *this.interval
}

func (this *protect_run_option) merge(delta *protect_run_option) *protect_run_option {
	if delta == nil {
		return this
	}
	if delta.interval != nil {
		this.interval = delta.interval
	}
	return this
}

func (this *protect_run_option) Merge(opts ...*protect_run_option) *protect_run_option {
	for _, option := range opts {
		this.merge(option)
	}
	return this
}

//==================================================

type line_option[T any] struct {
	shuffle_func func([]T) error
}

func LineOption[T any]() *line_option[T] {
	return &line_option[T]{}
}

func (this *line_option[T]) merge(delta *line_option[T]) *line_option[T] {
	if delta == nil {
		return this
	}
	if delta.shuffle_func != nil {
		this.shuffle_func = delta.shuffle_func
	}
	return this
}

func (this *line_option[T]) Merge(opts ...*line_option[T]) *line_option[T] {
	for _, opt := range opts {
		this.merge(opt)
	}
	return this
}

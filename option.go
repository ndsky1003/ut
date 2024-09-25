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

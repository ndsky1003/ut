package ut

import (
	"log/slog"
	"runtime/debug"
	"time"
)

func ProtectRun(f func(), opts ...*protect_run_option) {
	option := ProtectRunOption().Merge(opts...)
	interval := option.GetInterval()
	for {
		var is_panic bool
		func() {
			defer func() {
				is_panic = false
				if err := recover(); err != nil {
					is_panic = true
					slog.Error("recover", "err", err)
					debug.PrintStack()
				}
			}()
			f()
		}()
		if !is_panic {
			break
		}
		time.Sleep(interval)
	}
}

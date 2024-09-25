package ut

import (
	"fmt"
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
					fmt.Println("recover:", err)
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

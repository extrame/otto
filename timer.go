package otto

import (
	"fmt"
	"time"
)

type _timer struct {
	timer    *time.Timer
	duration time.Duration
	interval bool
	call     FunctionCall
}

func AddTimerToOtto(vm *Otto) error {

	registry := map[*_timer]*_timer{}
	ready := make(chan *_timer)

	newTimer := func(call FunctionCall, interval bool) (*_timer, Value) {
		delay, _ := call.Argument(1).ToInteger()
		if 0 >= delay {
			delay = 1
		}

		timer := &_timer{
			duration: time.Duration(delay) * time.Millisecond,
			call:     call,
			interval: interval,
		}
		registry[timer] = timer

		timer.timer = time.AfterFunc(timer.duration, func() {
			ready <- timer
		})

		value := vm.runtime.toValue(timer)

		return timer, value
	}

	setTimeout := func(call FunctionCall) Value {
		fmt.Println("setTimeout")
		_, value := newTimer(call, false)
		return value
	}
	vm.Set("setTimeout", setTimeout)

	setInterval := func(call FunctionCall) Value {
		_, value := newTimer(call, true)
		return value
	}
	vm.Set("setInterval", setInterval)

	clearTimeout := func(call FunctionCall) Value {
		timer, _ := call.Argument(0).Export()
		if timer, ok := timer.(*_timer); ok {
			timer.timer.Stop()
			delete(registry, timer)
		}
		return UndefinedValue()
	}
	vm.Set("clearTimeout", clearTimeout)
	vm.Set("clearInterval", clearTimeout)

	go func() {
		wg := vm.newRoutine()
		for {
			select {
			case timer := <-ready:
				var arguments []interface{}
				if len(timer.call.ArgumentList) > 2 {
					tmp := timer.call.ArgumentList[2:]
					arguments = make([]interface{}, 2+len(tmp))
					for i, value := range tmp {
						arguments[i+2] = value
					}
				} else {
					arguments = make([]interface{}, 1)
				}
				arguments[0] = timer.call.ArgumentList[0]
				_, err := vm.Call(`Function.call.call`, nil, arguments...)
				if err != nil {
					for _, timer := range registry {
						timer.timer.Stop()
						delete(registry, timer)
					}
				}
				if timer.interval {
					timer.timer.Reset(timer.duration)
				} else {
					delete(registry, timer)
				}
				// default:
				// 	fmt.Println("timer.go: waiting...")
				// Escape valve!
				// If this isn't here, we deadlock...
			}
			if len(registry) == 0 {
				break
			}
		}
		wg.Done()
	}()

	return nil
}

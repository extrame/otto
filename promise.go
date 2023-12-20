package otto

type promiseStatus int

const (
	promisePending promiseStatus = iota
	promiseFulfilled
	promiseRejected
)

type _promiseResolver func(arguments ...interface{})

type promise struct {
	resolver _promiseResolver

	status promiseStatus

	resolveFn *object
	rejectFn  *object

	result []interface{}
	err    *exception

	callbacks map[promiseStatus]_promiseResolver
}

func (p *promise) start(rt *runtime, resolver Value) {
	p.resolveFn = rt.newNativeFunction("get", "internal", 0, func(fn FunctionCall) Value {
		if p.status != promisePending {
			return UndefinedValue()
		}
		p.updateStatue(promiseFulfilled, fn.ArgumentList)
		return UndefinedValue()
	})
	p.rejectFn = rt.newNativeFunction("get", "internal", 0, func(fn FunctionCall) Value {
		if p.status != promisePending {
			return UndefinedValue()
		}
		p.updateStatue(promiseRejected, fn.ArgumentList)
		return UndefinedValue()
	})
	p.callbacks = make(map[promiseStatus]_promiseResolver)
	resolver.call(rt, UndefinedValue(), p.resolveFn, p.rejectFn)
}

func (p *promise) resolveBy(rt *runtime, result Value) {
	p.updateStatue(promiseFulfilled, result)
}

func (p *promise) updateStatue(status promiseStatus, arguments ...interface{}) {
	p.status = status
	if status == promiseFulfilled {
		p.result = arguments
	} else if status == promiseRejected {
		p.err = arguments[0].(*exception)
	}

	if callback, ok := p.callbacks[status]; ok {
		callback(arguments...)
	}
}

// return a new promise
func (p *promise) callThen(rt *runtime, this Value, fn Value) Value {
	var nextPromise = new(promise)
	defer func() {
		if p.status == promisePending {
			p.addCallback(promiseFulfilled, nextPromise.resolver)
			p.addCallback(promiseRejected, func(arguments ...interface{}) {
				nextPromise.updateStatue(promiseRejected, arguments)
			})
		} else if p.status == promiseFulfilled {
			nextPromise.resolver(p.result...)
		} else {
			nextPromise.updateStatue(promiseRejected, p.err)
		}
	}()
	nextPromise.resolver = nextPromise.wrapResolver(rt, this, fn)
	nextPromise.callbacks = make(map[promiseStatus]_promiseResolver)

	o := rt.newPromiseObject(nextPromise)
	o.prototype = rt.global.PromisePrototype
	return rt.toValue(o)
}

func (p *promise) callCatch(rt *runtime, this Value, fn Value) Value {
	if p.status == promiseRejected {
		fn.call(rt, this, p.rejectFn)
		return UndefinedValue()
	}
	p.addCallback(promiseRejected, func(arguments ...interface{}) {
		fn.call(rt, this, arguments)
	})
	return UndefinedValue()
}

func (p *promise) addCallback(status promiseStatus, callback _promiseResolver) {
	if p.status == status {
		callback()
		return
	}
	p.callbacks[status] = callback
}

func (p *promise) wrapResolver(rt *runtime, this Value, fn Value) func(arguments ...interface{}) {
	return func(arguments ...interface{}) {
		defer func() {
			if err := recover(); err != nil {
				p.updateStatue(promiseRejected, err)
			}
		}()

		var result = fn.call(rt, this, arguments)
		if result.Class() == "Promise" {
			var nextPromise = result.value.(*object).value.(*promise)
			nextPromise.addCallback(promiseFulfilled, func(arguments ...interface{}) {
				p.updateStatue(promiseFulfilled, arguments)
			})
			nextPromise.addCallback(promiseRejected, func(arguments ...interface{}) {
				p.updateStatue(promiseRejected, arguments)
			})
		} else {
			p.updateStatue(promiseFulfilled, result)
		}
	}
}

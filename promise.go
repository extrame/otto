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

func (p *promise) updateStatue(status promiseStatus, arguments ...interface{}) {
	p.status = status
	if callback, ok := p.callbacks[status]; ok {
		callback(arguments...)
	}
}

// return a new promise
func (p *promise) callThen(rt *runtime, this Value, fn Value) Value {
	var nextPromise promise
	nextPromise.resolver = nextPromise.wrapResolver(rt, this, fn)
	nextPromise.callbacks = make(map[promiseStatus]_promiseResolver)
	p.addCallback(promiseFulfilled, nextPromise.resolver)
	o := rt.newPromiseObject(&nextPromise)
	o.prototype = rt.global.PromisePrototype
	return rt.toValue(o)
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

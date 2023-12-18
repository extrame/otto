package otto

func builtinNewPromise(obj *object, argumentList []Value) Value {
	resolver := UndefinedValue()

	if len(argumentList) > 0 {
		resolver = argumentList[0]
	}
	return objectValue(obj.runtime.newPromise(resolver))
}

func builtinPromiseThen(call FunctionCall) Value {
	thisObject := call.thisObject()
	var p = thisObject.value.(*promise)
	var arg = call.Argument(0)
	if !arg.IsFunction() {
		panic(call.runtime.panicTypeError("Promise.then() requires a function argument"))
	}
	return p.callThen(call.runtime, call.This, arg)
}

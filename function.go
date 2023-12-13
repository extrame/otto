package otto

import (
	"fmt"
	"strings"
)

func (f *FunctionLiteral) String() string {
	return fmt.Sprintf("function %s(%s) { [native code] }", f.name, strings.Join(f.parameterList, ", "))
}

func (f *FunctionLiteral) CallScript() string {
	return fmt.Sprintf("%s(%s)", f.name, strings.Join(f.parameterList, ", "))
}

func (f *FunctionLiteral) Name() string {
	return f.name
}

func (f *FunctionLiteral) Call(call *FunctionCall) Value {
	value := call.runtime.cmplEvaluateNodeStatement(f.body)
	if value.kind == valueResult {
		return value.value.(result).value
	}
	return value
}

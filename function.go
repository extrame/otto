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

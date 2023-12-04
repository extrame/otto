package otto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetOwnPropertyNames(t *testing.T) {
	vm := New()
	// Commented out values aren't implemented yet.
	tests := map[string][]string{
		"Object.prototype": {
			"constructor",
			// "__defineGetter__",
			// "__defineSetter__",
			"hasOwnProperty",
			// "__lookupGetter__",
			// "__lookupSetter__",
			"isPrototypeOf",
			"propertyIsEnumerable",
			"toString",
			"valueOf",
			// "__proto__",
			"toLocaleString",
		},
		"Array.prototype": {
			"length",
			"constructor",
			// "at",
			"concat",
			// "copyWithin",
			// "fill",
			// "find",
			// "findIndex",
			"lastIndexOf",
			"pop",
			"push",
			"reverse",
			"shift",
			"unshift",
			"slice",
			"sort",
			"splice",
			// "includes",
			"indexOf",
			"join",
			// "keys",
			// "entries",
			// "values",
			"forEach",
			"filter",
			// "flat",
			// "flatMap",
			"map",
			"every",
			"some",
			"reduce",
			"reduceRight",
			"toLocaleString",
			"toString",
			// "findLast",
			// "findLastIndex",
		},
		"String.prototype": {
			"length",
			"constructor",
			// "anchor",
			// "at",
			// "big",
			// "blink",
			// "bold",
			"charAt",
			"charCodeAt",
			// "codePointAt",
			"concat",
			// "endsWith",
			// "fontcolor",
			// "fontsize",
			// "fixed",
			// "includes",
			"indexOf",
			// "italics",
			"lastIndexOf",
			// "link",
			"localeCompare",
			"match",
			// "matchAll",
			// "normalize",
			// "padEnd",
			// "padStart",
			// "repeat",
			"replace",
			// "replaceAll",
			"search",
			"slice",
			// "small",
			"split",
			// "strike",
			// "sub",
			"substr",
			"substring",
			// "sup",
			// "startsWith",
			"toString",
			"trim",
			// "trimStart",
			"trimLeft",
			// "trimEnd",
			"trimRight",
			"toLocaleLowerCase",
			"toLocaleUpperCase",
			"toLowerCase",
			"toUpperCase",
			"valueOf",
		},
		"Boolean.prototype": {
			"constructor",
			"toString",
			"valueOf",
		},
		"Number.prototype": {
			"constructor",
			"toExponential",
			"toFixed",
			"toPrecision",
			"toString",
			"valueOf",
			"toLocaleString",
		},
		"Math": {
			"abs",
			"acos",
			"acosh",
			"asin",
			"asinh",
			"atan",
			"atanh",
			"atan2",
			"cbrt",
			"ceil",
			// "clz32",
			"cos",
			"cosh",
			"exp",
			"expm1",
			"floor",
			// "fround",
			// "hypot",
			// "imul",
			"log",
			"log10",
			"log1p",
			"log2",
			"max",
			"min",
			"pow",
			"random",
			"round",
			// "sign",
			"sin",
			"sinh",
			"sqrt",
			"tan",
			"tanh",
			"trunc",
			"E",
			"LN10",
			"LN2",
			"LOG10E",
			"LOG2E",
			"PI",
			"SQRT1_2",
			"SQRT2",
		},
		"Date.prototype": {
			"constructor",
			"toString",
			"toDateString",
			"toTimeString",
			"toISOString",
			"toUTCString",
			"toGMTString",
			"getDate",
			"setDate",
			"getDay",
			"getFullYear",
			"setFullYear",
			"getHours",
			"setHours",
			"getMilliseconds",
			"setMilliseconds",
			"getMinutes",
			"setMinutes",
			"getMonth",
			"setMonth",
			"getSeconds",
			"setSeconds",
			"getTime",
			"setTime",
			"getTimezoneOffset",
			"getUTCDate",
			"setUTCDate",
			"getUTCDay",
			"getUTCFullYear",
			"setUTCFullYear",
			"getUTCHours",
			"setUTCHours",
			"getUTCMilliseconds",
			"setUTCMilliseconds",
			"getUTCMinutes",
			"setUTCMinutes",
			"getUTCMonth",
			"setUTCMonth",
			"getUTCSeconds",
			"setUTCSeconds",
			"valueOf",
			"getYear",
			"setYear",
			"toJSON",
			"toLocaleString",
			"toLocaleDateString",
			"toLocaleTimeString",
		},
		"RegExp.prototype": {
			"constructor",
			"exec",
			// "dotAll",
			// "flags",
			// "global",
			// "hasIndices",
			// "ignoreCase",
			// "multiline",
			// "source",
			// "sticky",
			// "unicode",
			"compile",
			"toString",
			"test",
		},
		"Error.prototype": {
			"constructor",
			"name",
			"message",
			"toString",
		},
		"EvalError.prototype": {
			"constructor",
			"name",
			"message",
			"toString",
		},
		"TypeError.prototype": {
			"constructor",
			"name",
			"message",
			"toString",
		},
		"ReferenceError.prototype": {
			"constructor",
			"name",
			"message",
			"toString",
		},
		"SyntaxError.prototype": {
			"constructor",
			"name",
			"message",
			"toString",
		},
		"URIError.prototype": {
			"constructor",
			"name",
			"message",
			"toString",
		},
		"JSON": {
			"parse",
			"stringify",
		},
		"NaN":      {},
		"Infinity": {},
	}
	for _, v := range []string{
		"eval", "parseInt", "parseFloat", "isNaN", "isFinite",
		"decodeURI", "decodeURIComponent", "encodeURI", "encodeURIComponent",
		"escape", "unescape",
	} {
		tests[v] = []string{"length", "name"}
	}
	for name, expected := range tests {
		t.Run(name, func(t *testing.T) {
			val, err := vm.Run(fmt.Sprintf("Object.getOwnPropertyNames(%s)", name))
			require.NoErrorf(t, err, "%#v", err)

			export, err := val.Export()
			require.NoError(t, err)
			if len(expected) == 0 {
				// Zero length doesn't know type.
				require.Equal(t, []interface{}{}, export)
				return
			}

			require.Equal(t, expected, export)
		})
	}
}

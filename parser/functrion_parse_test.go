package parser

import (
	"testing"

	"github.com/extrame/otto/ast"
)

func TestParserAnonymousFunctionSimple(t *testing.T) {
	tt(t, func() {
		test := func(input string, expect interface{}) (*ast.Program, *parser) {
			parser := newParser("", input, 1, nil)
			program, err := parser.parse()
			is(firstErr(err), expect)
			return program, parser
		}

		// test("(1,2,3)", nil)
		// test("()", nil)
		test("()=>{}", nil)
		test("(first,second)=>{}", nil)
	})
}

func TestParserAnonymousFunction(t *testing.T) {
	tt(t, func() {
		test := func(input string, expect interface{}) (*ast.Program, *parser) {
			parser := newParser("", input, 1, nil)
			program, err := parser.parse()
			is(firstErr(err), expect)
			return program, parser
		}

		test("var f = (1,2,3)", nil)
		test("var f = ()", nil)
		test("var f = ()=>{}", nil)
	})
}

func TestParserFunction(t *testing.T) {
	tt(t, func() {
		test := func(input string, expect interface{}) (*ast.Program, *parser) {
			parser := newParser("", input, 1, nil)
			program, err := parser.parse()
			is(firstErr(err), expect)
			return program, parser
		}

		test("var f = function(){}", nil)
	})
}

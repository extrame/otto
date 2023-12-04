package otto

import (
	"errors"
	"fmt"
	"testing"

	"github.com/extrame/otto/file"
	"github.com/extrame/otto/parser"
	"github.com/extrame/otto/token"
)

type Fixer struct {
}

func (f *Fixer) Fix(src string, idx file.Idx, tok token.Token) (*parser.Fixture, error) {
	fmt.Println("Fixer.Fix called", idx, tok)
	if tok == token.IDENTIFIER {
		return &parser.Fixture{
			Str: " fn",
		}, nil
	}
	return nil, errors.New("TODO")
}

func RunJS(code string) error {
	vm := New()
	vm.runtime.newArgumentsObject()
	script, err := vm.Compile("", code)
	if err != nil {
		return err
	}
	script.program.body = append(script.program.body, &nodeExpressionStatement{
		expression: &nodeCallExpression{
			callee: &nodeIdentifier{
				name: "login",
			},
			argumentList: []nodeExpression{
				&nodeIdentifier{
					name: "login",
				}, &nodeIdentifier{
					name: "login",
				},
			},
		},
	})

	_, err = vm.Run(script)
	return err
}

func RunJS2(code string) error {
	vm := New()
	script, err := vm.CompileWithFixer("", code, nil)
	if err != nil {
		return err
	}
	_, err = vm.Run(script)
	return err
}

func TestRunJS1(t *testing.T) {
	// Test case 1: Valid JavaScript code
	err := RunJS2("function() {var email = '';var username = '';var phone = '';if(register()){var email = getAttribute('Email');var username = getAttribute('Username');var phone = getAttribute('Phone');if(email.trim() != '' \u0026\u0026 username.trim() != '' \u0026\u0026 phone.trim() != ''){axios.post('/api/login', {email: email, username: username, phone: phone}).then(response =\u003e {if(response.data.code == 0){router.push({path: '/system'});}else{alert('登录失败，请检查您的信息！')}});}}}}")
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	// Test case 2: Invalid JavaScript code
	err = RunJS("console.log('Hello, world!'")
	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	}
}

func TestRunJS2(t *testing.T) {
	// Test case 1: Valid JavaScript code
	err := RunJS(`
	var req = {};
	var res = {};
	function login(req, res) {  var username = req.getQuery('username');  var password = req.getQuery('password');  var service = services.Get('UserService');  service.validateLogin(username, password)  .then(function(result) {    if (result.success) {      res.send({ success: true });    } else {      res.send({ success: false, message: result.message });    }  })  .catch(function(error) {    res.send({ success: false, message: error });  });}
	login(req, res);
	`)
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	// Test case 2: Invalid JavaScript code
	err = RunJS2("console.log('Hello, world!'")
	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	}
}

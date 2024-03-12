package otto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPromiseCall(t *testing.T) {
	tt(t, func() {
		vm := New()
		script, err := vm.Compile("", `Promise()`)
		require.NoError(t, err)
		_, err = vm.Run(script)
		require.EqualError(t, err, "TypeError: Promise constructor cannot be invoked without 'new'")
	})
}

func TestPromiseNew(t *testing.T) {
	tt(t, func() {
		vm := New()
		script, err := vm.Compile("", `new Promise()`)
		require.NoError(t, err)
		_, err = vm.Run(script)
		require.EqualError(t, err, "TypeError: Promise resolver undefined is not a function")
	})
}

func TestPromiseNew1(t *testing.T) {
	tt(t, func() {
		vm := New()
		script, err := vm.Compile("", `new Promise(1)`)
		require.NoError(t, err)
		_, err = vm.Run(script)
		require.EqualError(t, err, "TypeError: Promise resolver 1 is not a function")
	})
}

func TestPromiseNewWithThen(t *testing.T) {
	tt(t, func() {
		vm := New()
		script, err := vm.Compile("", `new Promise(function(){console.log(this)});`)
		require.NoError(t, err)
		_, err = vm.Run(script)
		require.EqualError(t, err, "TypeError: Promise resolver 1 is not a function")
	})
}

func TestPromiseNewWithThenAnony(t *testing.T) {
	tt(t, func() {
		vm := New()
		script, err := vm.Compile("", `new Promise(()=>{console.log(this)});`)
		require.NoError(t, err)
		_, err = vm.Run(script)
		require.EqualError(t, err, "TypeError: Promise resolver 1 is not a function")
	})
}

func TestPromiseNewWithFuncThen(t *testing.T) {
	tt(t, func() {
		vm := New()
		AddTimerToOtto(vm)
		script, err := vm.Compile("", `new Promise((resolve,reject)=>{
			setTimeout(function(){
				resolve(1);
			},1000);
		}).then(()=>{
			console.log("second");
		});
		console.log("first")
		`)
		require.NoError(t, err)
		_, err = vm.Run(script)
		require.EqualError(t, err, "TypeError: Promise resolver 1 is not a function")
	})
}

func TestPromiseNewWithFuncThen2(t *testing.T) {
	tt(t, func() {
		vm := New()
		AddTimerToOtto(vm)
		script, err := vm.Compile("", `var p = new Promise((resolve, reject) => {
			console.log(0);
			  resolve(1);
		}).then(() => {
			console.log(1);
			return new Promise((resolve1, reject1) => {
			  console.log(2);
			  setTimeout(() => {
				console.log(3);
				resolve1(2);
			  },500);
			});
		  }).then(()=>{console.log(4);});
		  console.log(p);
		`)
		require.NoError(t, err)
		_, err = vm.Run(script)
		require.EqualError(t, err, "TypeError: Promise resolver 1 is not a function")
	})
}

func TestPromiseNewWithFuncThen3(t *testing.T) {
	tt(t, func() {
		vm := New()
		AddTimerToOtto(vm)
		script, err := vm.Compile("", `var p = new Promise((resolve, reject) => {
			console.log(0);
			  resolve(1);
		}).then(() => {
			console.log(1);
			return new Promise((resolve1, reject1) => {
			  console.log(2);
			  setTimeout(() => {
				console.log(3);
				reject1(2);
			  },500);
			});
		  }).then(()=>{console.log(4);});
		  console.log(p);
		`)
		require.NoError(t, err)
		_, err = vm.Run(script)
		require.EqualError(t, err, "TypeError: Promise resolver 1 is not a function")
	})
}

func TestPromiseNewWithFuncThen4(t *testing.T) {
	tt(t, func() {
		vm := New()
		AddTimerToOtto(vm)
		script, err := vm.Compile("", `var p = new Promise((resolve, reject) => {
			console.log(0);
			  resolve(1);
		}).then(() => {
			console.log(1);
			throw new Error("error");
		  }).then(()=>{console.log(4);}).catch((e)=>{console.log(5);});
		  console.log(p);
		`)
		require.NoError(t, err)
		_, err = vm.Run(script)
		require.EqualError(t, err, "TypeError: Promise resolver 1 is not a function")
	})
}

func TestPromiseEffectedCalculator(t *testing.T) {
	tt(t, func() {
		vm := New()
		AddTimerToOtto(vm)
		_, err := vm.Compile("", `
		console.log( (1+2)*3);
		`)
		require.NoError(t, err)
	})
}

func TestPromiseEffectedFunction(t *testing.T) {
	tt(t, func() {
		vm := New()
		AddTimerToOtto(vm)
		script, err := vm.Compile("", `
		var fn = function(num){
			console.log(num);
		}
		console.log(fn);
		fn(2);
		`)
		require.NoError(t, err)
		_, err = vm.Run(script)
		require.EqualError(t, err, "TypeError: Cannot convert function to object")
	})
}

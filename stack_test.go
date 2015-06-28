package gostack

import (
	"errors"
	"runtime"
	"testing"
)

func testStacks(t *testing.T, stacks Stacks, testPanic bool, expectedFnName string) {
	t.Logf("\n\n%s\n\n", stacks)

	buf := make([]byte, 2<<20)
	buf = buf[:runtime.Stack(buf, false)]
	t.Log(string(buf))

	if testPanic {
		panicCall := stacks.GetPanic()
		if panicCall == nil {
			t.Errorf("Panic call not found")
			return
		}
		if panicCall.Func.Name() != expectedFnName {
			t.Errorf("GetPanic() returned wrong stuff: %v", panicCall)
			return
		}
	} else {

		for _, stack := range stacks {
			if stack.Func.Name() == expectedFnName { // Ok.
				return
			}
		}
		t.Errorf("Expected function stack not found.")
	}
}

func TestGetForPanic(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fail()
			return
		}
		if e, ok := err.(error); ok {
			testStacks(t, Get(e), true, "github.com/ceram1/gostack.TestGetForPanic")
		} else {
			t.Fail()
		}
	}()

	panic(errors.New("TEST_ERROR"))
}

func TestGetForNewError(t *testing.T) {
	err := errors.New("TestGetForNewError")

	testStacks(t, Get(err), false, "github.com/ceram1/gostack.TestGetForNewError")

}

// golang issue: https://github.com/golang/go/issues/11440
//
// func TestGetInGoroutine(t *testing.T) {
// 	go func() {
// 		defer func() {
// 			err := recover()
// 			if err == nil {
// 				t.Fail()
// 				return
// 			}
// 			if e, ok := err.(error); ok {
// 				testStacks(t, Get(e), "github.com/ceram1/gostack.TestGetInGoroutine")
// 			} else {
// 				t.Fail()
// 			}
// 		}()
//
// 		panic(errors.New("TEST_ERROR"))
// 	}()
// }

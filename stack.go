package gostack

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// Stack is a line in stacktrace.
type Stack struct {
	File string
	Line int

	Func *runtime.Func
	// FuncName string

}

// Stacks is a array of Stack.
type Stacks []Stack

// Get stacktrace of an error.
// Currently, this cannnot get function created goroutine.
//
// golang issue: https://github.com/golang/go/issues/11440
func Get(err error) Stacks {
	var stacks []Stack
	for skip := 1; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)

		if !ok { // not recoverable.
			break
		}

		if strings.HasSuffix(file, "c") { // Skip if it's c file.
			continue
		}

		fn := runtime.FuncForPC(pc)

		stack := Stack{file, line, fn}
		stacks = append(stacks, stack)
	}
	return stacks
}

func (s Stack) String() string {
	return fmt.Sprintf("%s:%d %s()\n", s.File, s.Line, s.Func.Name())
}

func (s Stacks) String() string {
	var b bytes.Buffer

	for _, st := range s {
		b.WriteString(st.String())
	}

	return b.String()
}

// GetPanic returns panic(call)
// return nil if not found.
func (s Stacks) GetPanic() *Stack {
	found := false

	for _, st := range s {
		if found {
			return &st
		}

		if st.Func.Name() == "runtime.gopanic" {
			found = true
		}
	}

	return nil
}

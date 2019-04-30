package util

import (
	"fmt"
	"runtime/debug"
)

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception struct {
	Error      interface{}
	StackTrace []byte
}

func (e Exception) String() string {
	return fmt.Sprintf("%s\n%s", e.Error, string(e.StackTrace))
}

func (tcf Block) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(Exception{
					Error:      r,
					StackTrace: debug.Stack(),
				})
			}
		}()
	}
	tcf.Try()
}

package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	cl          *object.Closure // the frame's associated closure
	ip          int             // internal IP ptr for this specific frame and func
	basePointer int             // keeps track of the stack ptr pos at the call site
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	f := &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}
	return f
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}

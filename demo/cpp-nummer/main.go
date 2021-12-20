package main

// #include "num.h"
import "C"
import (
	"reflect"
	"unsafe"
)

type GoNum struct {
	num C.Num
}

func New() GoNum {
	var ret GoNum
	ret.num = C.NumInit()
	return ret
}
func (n GoNum) Free() {
	C.NumFree((C.Num)(unsafe.Pointer(n.num)))
}
func (n GoNum) Inc() {
	C.NumIncrement((C.Num)(unsafe.Pointer(n.num)))
}
func (n GoNum) GetValue() int {
	return int(C.NumGetValue((C.Num)(unsafe.Pointer(n.num))))
}

func main() {
	num := New()
	num.Inc()
	if num.GetValue() != 2 {
		panic("unexpected value received")
	}
	num.Inc()
	num.Inc()
	num.Inc()
	if num.GetValue() != 5 {
		panic("unexpected value received")
	}
	value := num.GetValue()
	num.Free()

	typ := reflect.TypeOf(value)
	if typ.Name() != "int" {
		panic("got unexpected type")
	}
}

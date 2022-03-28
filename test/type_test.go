package test

import (
	"github.com/rarnu/gojvm"
	"testing"
)

func TestTypes(t *testing.T) {
	var i0 int32 = 1
	var i1 int = 2
	var i2 int64 = 1
	var f1 float32 = 1.1
	var f2 float64 = 2.0
	var f3 byte = 200
	var s0 = 'a'
	var s1 string = "abc"
	gojvm.ArgumentsCheck(true, i0, i1, i2, f1, f2, f3, s0, s1, []int{1, 2, 3}, map[string]int{"a": 1, "b": 2})
}

package test

import (
	"fmt"
	"github.com/rarnu/gojvm"
	"testing"
)

func TestJvm(t *testing.T) {
	jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.FindClass("Hello")
	fmt.Printf("%v\n", c1)

	s0, _ := c1.CallStaticStringMethod("hello", "rarnu", 233)
	fmt.Printf("%v\n", s0)

	// c2 := env.FindClass("Hello1")
	// fmt.Printf("%v\n", c2)

	env.Detach()
	jvm.Free()
}

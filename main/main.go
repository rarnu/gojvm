package main

import (
	"fmt"
	"github.com/rarnu/gojvm"
)

func main() {
	jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.FindClass("Hello")

	s0, _ := c1.CallStaticStringMethod("hello", "rarnu", 233)
	fmt.Printf("%v\n", s0)

	// c2 := env.FindClass("Hello1")
	// fmt.Printf("%v\n", c2)

	env.Detach()
	jvm.Free()
}

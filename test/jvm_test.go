package test

import (
	"fmt"
	"github.com/rarnu/gojvm"
	"testing"
)

func TestStaticInvoke(t *testing.T) {
	jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()
	c1 := env.FindClass("Hello")
	s0, _ := c1.InvokeString("hello", "rarnu", 8)
	fmt.Printf("%v\n", s0)
	env.Detach()
	jvm.Free()
}

func TestStaticGetClass(t *testing.T) {
	jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.FindClass("Hello")
	jo := c1.GetObject("h", "H1")
	fmt.Printf("%v\n", jo)
	jo.Free()
	c1.Free()

	env.Detach()
	jvm.Free()
}

func TestNewClass(t *testing.T) {
	jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.NewObject("Hello")
	fmt.Printf("%v\n", c1)
	c1.Free()

	env.Detach()
	jvm.Free()
}

func TestClassFields(t *testing.T) {
	jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.NewObject("H1")
	fmt.Printf("%v\n", c1)

	s1 := c1.GetString("v")
	fmt.Printf("%v\n", s1)

	c1.SetString("v", "rarnu")
	s2 := c1.GetString("v")
	fmt.Printf("%v\n", s2)

	c1.Free()

	env.Detach()
	jvm.Free()
}

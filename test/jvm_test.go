package test

import (
	"github.com/rarnu/gojvm"
	"runtime"
	"testing"
)

func showMem(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	t.Logf("Alloc = %v bytes", m.Alloc)
}

func TestStaticInvoke(t *testing.T) {
	jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()
	c1 := env.FindClass("Hello")
	s0, _ := c1.InvokeString("hello", "rarnu", 8)
	t.Logf("%v\n", s0)
	env.Detach()
	jvm.Free()
}

func TestStaticGetClass(t *testing.T) {
	jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.FindClass("Hello")
	jo := c1.GetObject("h", "H1")
	t.Logf("%v\n", jo)
	jo.Free()
	c1.Free()

	env.Detach()
	jvm.Free()
}

func TestNewClass(t *testing.T) {
	jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.NewObject("Hello")
	t.Logf("%v\n", c1)
	c1.Free()

	env.Detach()
	jvm.Free()
}

func TestClassFields(t *testing.T) {
	showMem(t)
	jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	showMem(t)

	c1 := env.NewObject("H1")
	t.Logf("%v\n", c1)

	s1 := c1.GetString("v")
	t.Logf("%v\n", s1)

	c1.SetString("v", "rarnu")
	s2 := c1.GetString("v")
	t.Logf("%v\n", s2)

	showMem(t)

	c1.Free()

	showMem(t)

	env.Detach()
	jvm.Free()

	showMem(t)

	runtime.GC()
	showMem(t)
}

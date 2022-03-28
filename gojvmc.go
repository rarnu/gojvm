package gojvm

import (
	"fmt"
	"runtime"
)

func printMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("malloc: %d\n", m.Mallocs)
	fmt.Printf("free: %d\n", m.Frees)
}

func CallC() {
	jvm := NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.FindClass("Hello")
	fmt.Printf("%v\n", c1)
	c2 := env.FindClass("Hello1")
	fmt.Printf("%v\n", c2)

	env.Detach()
	jvm.Free()
}

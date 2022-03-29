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

}

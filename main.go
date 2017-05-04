package main

import h "./handlers"
import (
	"fmt"
	"time"
)

func main() {
	for {
		cpu, mem := h.Load_CPU_MEM()
		fmt.Printf("CPU: %.2f MEM: %.2f \n", cpu, mem)
		time.Sleep(time.Second * 2)
	}
}

package main

import (
	"fmt"	
)
import h "./handlers"

func main() {
	disk := h.GetDiskUsage("/")
	fmt.Printf("All:%.2f GB, Free: %.2f GB, Used: %.2f GB ", (float64(disk.All) / float64(h.GB)), (float64(disk.Free) / float64(h.GB)), (float64(disk.Used) / float64(h.GB)))

}

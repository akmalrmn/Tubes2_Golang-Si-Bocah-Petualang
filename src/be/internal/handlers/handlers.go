package handlers

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
)

func UpdateMaxGoroutine() int {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	// Convert to GB
	availableMemory := v.Available / 1024 / 1024 / 1024

	c, err := cpu.Counts(true)
	if err != nil {
		log.Fatal(err)
	}

	if availableMemory < 1 {
		return 1
	} else {
		// Return CPU Core
		return c
	}
}
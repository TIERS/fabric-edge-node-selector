package internal

import (
	"github.com/shirou/gopsutil/mem"
)

func GetMemoryUsage() (MemStats DrcMemStats) {
	tmpMem, _ := mem.VirtualMemory()
	return DrcMemStats{
		Total:     tmpMem.Total,
		Available: tmpMem.Available,
		Used:      tmpMem.UsedPercent,
	}
}

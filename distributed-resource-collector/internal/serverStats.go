package internal

import (
	"time"
)

func GetServerStats() (dcrStats DrcStats) {
	tmpTime := time.Now()
	timestamp := DrcTimestamp{
		TimeLocal:   tmpTime,
		TimeSeconds: tmpTime.Unix(),
		TimeNano:    tmpTime.UnixNano(),
	}
	return DrcStats{
		Timestamp:  timestamp,
		DrcHost:    GetHostStats(),
		CPUStats:   GetCPUUsage(),
		MemStats:   GetMemoryUsage(),
		DiskStats:  GetDiskUsage(),
		ProcStats:  GetProcStats(),
		DockerSats: GetDockerStats(),
	}
}

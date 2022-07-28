package internal

import (
	"github.com/shirou/gopsutil/load"
)

func GetProcStats() (procStats DrcProcStats) {
	tmpProcs, _ := load.Misc()
	return DrcProcStats{
		TotalProcs:   tmpProcs.ProcsTotal,
		CreatedProcs: tmpProcs.ProcsCreated,
		RunningProcs: tmpProcs.ProcsRunning,
		BlockedProcs: tmpProcs.ProcsBlocked,
	}
}

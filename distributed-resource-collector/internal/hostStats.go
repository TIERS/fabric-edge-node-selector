package internal

import (
	"github.com/shirou/gopsutil/host"
)

func GetHostStats() (drcHost DrcHost) {
	tmpHost, _ := host.Info()
	return DrcHost{
		Hostname:             tmpHost.Hostname,
		Uptime:               int64(tmpHost.Uptime),
		BootTime:             int64(tmpHost.BootTime),
		Platform:             tmpHost.Platform,
		VirtualizationSystem: tmpHost.VirtualizationSystem,
		VirtualizationRole:   tmpHost.VirtualizationRole,
		HostID:               tmpHost.HostID,
	}
}

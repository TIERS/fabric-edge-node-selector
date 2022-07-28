package internal

import (
	"github.com/shirou/gopsutil/disk"
)

func GetDiskUsage() []DrcDiskStats {
	parts, err := disk.Partitions(false)
	CheckError(err)

	var drcUsage []DrcDiskStats

	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		CheckError(err)

		if !CustomContains(u.Path, "/snap/", "/etc/") {
			tmpUsage := DrcDiskStats{
				Device: part.Device,
				//SerialNumber: disk.GetDiskSerialNumber(part.Device),
				Path:        u.Path,
				Label:       disk.GetLabel(part.Device),
				Fstype:      part.Fstype,
				Total:       u.Total,
				Used:        u.Used,
				UsedPercent: u.UsedPercent,
			}
			//fmt.Println(tmpUsage)
			drcUsage = append(drcUsage, tmpUsage)
		}
	}
	return drcUsage
}

package internal

import (
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// This function uses GO-PsUtil to get the CPU usage information
// There is a function that brings the complete average usage, however, I opted for each core performance
// Both can be used, however, the average usage is calculated in a similar fashion, and it would require
// running a second function. If we calculate the average from the data we retrieved, we have more information to work with.
func GetCPUUsage() (CPUStats DrcCPUStats) {
	tmpCPU, _ := cpu.Percent(time.Second/5, true)
	totalPercent := 0.0
	for _, percent := range tmpCPU {
		totalPercent += percent
	}

	vendorList := []string{}
	modelList := []string{}
	cpuInfo, _ := cpu.Info()
	for _, cpu := range cpuInfo {
		vendorList = append(vendorList, cpu.VendorID)
		modelList = append(modelList, cpu.ModelName)
	}

	modelList = UniqueString(modelList)
	model := ""
	if len(modelList) > 1 {
		model = strings.Join(modelList, " / ")
	} else {
		model = modelList[0]
	}
	vendorList = UniqueString(vendorList)
	vendor := ""
	if len(vendorList) > 1 {
		vendor = strings.Join(vendorList, " / ")
	} else {
		vendor = vendorList[0]
	}

	return DrcCPUStats{
		ModelName:    model,
		VendorID:     vendor,
		AverageUsage: totalPercent / float64(len(tmpCPU)),
		CoreUsage:    tmpCPU,
	}
}

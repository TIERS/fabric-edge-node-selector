package internal

import (
	"encoding/json"
	"time"

	"github.com/wI2L/jettison"
)

// -- CPU
type DrcCPUStats struct {
	ModelName    string    `json:"modelName"`
	VendorID     string    `json:"vendorId"`
	AverageUsage float64   `json:"averageUsage"`
	CoreUsage    []float64 `json:"coreUsage"`
}

func (d DrcCPUStats) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

// -- DISK
type DrcDiskStats struct {
	Device string `json:"device"`
	//SerialNumber string  `json:"serialNumber"`
	Path        string  `json:"path"`
	Label       string  `json:"label"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

func (d DrcDiskStats) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

// -- MEMORY / RAM
type DrcMemStats struct {
	Total     uint64  `json:"total"`
	Available uint64  `json:"available"`
	Used      float64 `json:"used"`
}

func (d DrcMemStats) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

// -- PROCESSES
type DrcProcStats struct {
	TotalProcs   int `json:"totalProcs"`
	CreatedProcs int `json:"createdProcs"`
	RunningProcs int `json:"runningProcs"`
	BlockedProcs int `json:"blockedProcs"`
}

func (d DrcProcStats) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

// -- DOCKER
type DrcDockerStats struct {
	ContainerID string `json:"containerID"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Status      string `json:"status"`
	State       string `json:"State"`
}

func (d DrcDockerStats) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

// -- TIMESTAMP
type DrcTimestamp struct {
	TimeLocal   time.Time `json:"timeLocal"`
	TimeSeconds int64     `json:"timeSeconds"`
	TimeNano    int64     `json:"timeNano"`
}

func (d DrcTimestamp) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

// -- HOST INFO
type DrcHost struct {
	Hostname             string `json:"hostname"`
	Uptime               int64  `json:"uptime"`
	BootTime             int64  `json:"boottime"`
	Platform             string `json:"platform"`
	VirtualizationSystem string `json:"virtualizationSystem"`
	VirtualizationRole   string `json:"virtualizationRole"`
	HostID               string `json:"hostid"`
}

func (d DrcHost) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

// -- RESPONSE OBJECT
type DrcStats struct {
	Timestamp  DrcTimestamp     `json:"timestamp"`
	DrcHost    DrcHost          `json:"host"`
	CPUStats   DrcCPUStats      `json:"cpuStats"`
	MemStats   DrcMemStats      `json:"memStats"`
	DiskStats  []DrcDiskStats   `json:"diskStats"`
	ProcStats  DrcProcStats     `json:"procStats"`
	DockerSats []DrcDockerStats `json:"dockerStats"`
}

type StoredStat struct {
	ID         string           `json:"id"`
	Hostname   string           `json:"hostname"`
	Timestamp  DrcTimestamp     `json:"timestamp"`
	DrcHost    DrcHost          `json:"host"`
	CPUStats   DrcCPUStats      `json:"cpuStats"`
	MemStats   DrcMemStats      `json:"memStats"`
	DiskStats  []DrcDiskStats   `json:"diskStats"`
	ProcStats  DrcProcStats     `json:"procStats"`
	DockerSats []DrcDockerStats `json:"dockerStats"`
}

func JsonToStoredStat(v string) (storedStat StoredStat, err error) {
	err = json.Unmarshal([]byte(v), &storedStat)
	return storedStat, err
}

func ArrayStoredStat(v string) (storedStats []StoredStat, err error) {
	err = json.Unmarshal([]byte(v), &storedStats)
	return storedStats, err
}

func (d DrcStats) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

func DrcJsonToStruct(v string) (drcStats DrcStats, err error) {
	err = json.Unmarshal([]byte(v), &drcStats)
	return drcStats, err
}

func ConvertToStorage(drcStats DrcStats) StoredStat {
	return StoredStat{
		ID:         "",
		Hostname:   "",
		Timestamp:  drcStats.Timestamp,
		DrcHost:    drcStats.DrcHost,
		CPUStats:   drcStats.CPUStats,
		MemStats:   drcStats.MemStats,
		DiskStats:  drcStats.DiskStats,
		ProcStats:  drcStats.ProcStats,
		DockerSats: drcStats.DockerSats,
	}
}

func (d StoredStat) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

////

type StatSummary struct {
	ID                  string       `json:"id"`
	Timestamp           DrcTimestamp `json:"timestamp"`
	CPUAverageUsage     float64      `json:"cpuAverageUsage"`
	MemoryUsePercentage float64      `json:"MemoryUsePercentage"`
	ContainersRunning   int          `json:"containersRunning"`
}

func SummarizeStoredStat(d StoredStat) StatSummary {
	var summary StatSummary
	summary.ID = d.ID
	summary.Timestamp = d.Timestamp
	summary.CPUAverageUsage = d.CPUStats.AverageUsage
	summary.MemoryUsePercentage = d.MemStats.Used

	var runningCount = 0
	for _, v := range d.DockerSats {
		if v.Status == "running" || v.State == "true" {
			runningCount += 1
		}
	}

	summary.ContainersRunning = runningCount

	return summary
}

func (d StatSummary) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

type StatAnalysis struct {
	Hostname            string        `json:"hostname"`
	Duration            int           `json:"duration"`
	CPUAverageUsage     float64       `json:"cpuAverageUsage"`
	MemoryUsePercentage float64       `json:"MemoryUsePercentage"`
	ContainersRunning   int           `json:"containersRunning"`
	StatSummary         []StatSummary `json:"statSummary"`
}

func AnalizeStatSummary(statAnalysis StatAnalysis) StatAnalysis {
	if len(statAnalysis.StatSummary) > 0 {
		var CPUAverageUsage float64 = 0
		var MemoryUsePercentage float64 = 0
		var ContainersRunning int = 0
		for _, summary := range statAnalysis.StatSummary {
			CPUAverageUsage += summary.CPUAverageUsage
			MemoryUsePercentage += summary.MemoryUsePercentage
			ContainersRunning += summary.ContainersRunning
		}
		statAnalysis.CPUAverageUsage = CPUAverageUsage / float64(len(statAnalysis.StatSummary))
		statAnalysis.MemoryUsePercentage = MemoryUsePercentage / float64(len(statAnalysis.StatSummary))
		statAnalysis.ContainersRunning = (ContainersRunning / len(statAnalysis.StatSummary))
	}

	return statAnalysis
}

func (d StatAnalysis) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

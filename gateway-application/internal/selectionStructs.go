package internal

import (
	"encoding/json"
	"time"

	"github.com/wI2L/jettison"
)

// Selection
type Timestamp struct {
	TimeLocal   time.Time `json:"timeLocal"`
	TimeSeconds int64     `json:"timeSeconds"`
	TimeNano    int64     `json:"timeNano"`
}

type ServerSelection struct {
	Asset               Asset   `json:"asset"`
	Target              string  `json:"target"`
	AverageLatency      float64 `json:"averageLatency"`
	CPUAverageUsage     float64 `json:"cpuAverageUsage"`
	MemoryUsePercentage float64 `json:"memoryUsePercentage"`
	ContainersRunning   int     `json:"containersRunning"`
}

func (d ServerSelection) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

func JsonToServerSelection(v string) (selection ServerSelection, err error) {
	err = json.Unmarshal([]byte(v), &selection)
	return selection, err
}

func StoreSelection(s ServerSelection) StoredSelection {
	tmpTime := time.Now()
	timestamp := Timestamp{
		TimeLocal:   tmpTime,
		TimeSeconds: tmpTime.Unix(),
		TimeNano:    tmpTime.UnixNano(),
	}
	return StoredSelection{
		ID:                  s.Target + "-" + DateFormatID(timestamp.TimeSeconds),
		AssetID:             s.Asset.ID,
		Target:              s.Target,
		Timestamp:           timestamp,
		AverageLatency:      s.AverageLatency,
		CPUAverageUsage:     s.CPUAverageUsage,
		MemoryUsePercentage: s.MemoryUsePercentage,
		ContainersRunning:   s.ContainersRunning,
	}
}

/// Selection Store
type StoredSelection struct {
	ID                  string    `json:"id"`
	AssetID             string    `json:"assetID"`
	Target              string    `json:"target"`
	Timestamp           Timestamp `json:"timestamp"`
	AverageLatency      float64   `json:"averageLatency"`
	CPUAverageUsage     float64   `json:"cpuAverageUsage"`
	MemoryUsePercentage float64   `json:"memoryUsePercentage"`
	ContainersRunning   int       `json:"containersRunning"`
}

func (d StoredSelection) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

func JsonToStoredSelection(v string) (selection StoredSelection, err error) {
	err = json.Unmarshal([]byte(v), &selection)
	return selection, err
}

func JsonToStoredSelectionArray(v string) (selection []StoredSelection, err error) {
	err = json.Unmarshal([]byte(v), &selection)
	return selection, err
}

package internal

import (
	"encoding/json"
	"time"

	"github.com/wI2L/jettison"
)

// INVENTORY ASSET
type Asset struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Owner      string     `json:"owner"`
	Type       int        `json:"type"`       //[0: Server, 1: Robot, 2: Sensor]
	State      int        `json:"state"`      //[0: Disabled, 1: Enabled]
	Properties Properties `json:"properties"` //{GPU: TRUE ...}
}

// PROPERTY ASSET
// Can be expanded to match the evolution of the PDP (Policy Decision Point) that determines how the Edge Server is selected
// Updated from being a simple map[string]string because it would be difficult to index the results in CouchDB otherwise (data integrity)
// Storing the information in plain text is not recommended due to security issues, even if the data can be saved as private in the Blockchain
// instead, servers should be assigned SSH keys
type Properties struct {
	GPU          int    `json:"gpu"` //0 = false, 1 = true
	Hostname     string `json:"hostname"`
	HostPort     string `json:"hostPort"`
	HostUser     string `json:"hostUser"`
	HostPassword string `json:"hostPassword"`
}

func (d Asset) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

func AssetArrayToJson(d []Asset) []byte {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return s
}

func JsonToAsset(v string) (asset Asset, err error) {
	err = json.Unmarshal([]byte(v), &asset)
	return asset, err
}

func JsonToProperties(v string) (properties Properties, err error) {
	err = json.Unmarshal([]byte(v), &properties)
	return properties, err
}

func JsonToAssetArray(v string) (assets []Asset, err error) {
	err = json.Unmarshal([]byte(v), &assets)
	return assets, err
}

///

type Timestamp struct {
	TimeLocal   time.Time `json:"timeLocal"`
	TimeSeconds int64     `json:"timeSeconds"`
	TimeNano    int64     `json:"timeNano"`
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

package internal

import (
	"encoding/json"
	"time"

	"github.com/wI2L/jettison"
)

// -- Latency Targets

type LatencyTargets struct {
	Source  string          `json:"source"`
	Targets []LatencyTarget `json:"targets"`
}

type LatencyTarget struct {
	Hostname     string `json:"hostname"`
	Hostport     string `json:"hostPort"`
	HostUser     string `json:"hostUser"`
	HostPassword string `json:"hostPassword"`
}

func LatencyTargetsJsonToStruct(v string) (targets LatencyTargets, err error) {
	err = json.Unmarshal([]byte(v), &targets)
	return targets, err
}

// -- TIMESTAMP
type LatencyTimestamp struct {
	TimeLocal   time.Time `json:"timeLocal"`
	TimeSeconds int64     `json:"timeSeconds"`
	TimeNano    int64     `json:"timeNano"`
}

func (d LatencyTimestamp) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

// Latency Results
type LatencyResults struct {
	Source    string           `json:"source"`
	Timestamp LatencyTimestamp `json:"timestamp"`
	Results   []LatencyResult  `json:"results"`
}

func LatencyResultsJsonToStruct(v string) (targets LatencyResults, err error) {
	err = json.Unmarshal([]byte(v), &targets)
	return targets, err
}

func (r LatencyResults) String() string {
	s, _ := json.Marshal(r)
	return string(s)
}

type LatencyResult struct {
	Hostname string `json:"hostname"`
	Latency  int64  `json:"latency"`
}

func (r LatencyResult) String() string {
	s, _ := json.Marshal(r)
	return string(s)
}

func LatencyJsonToStruct(v string) (targets LatencyTargets, err error) {
	err = json.Unmarshal([]byte(v), &targets)
	return targets, err
}

type LatencyAsset struct {
	ID        string           `json:"id"`
	Source    string           `json:"source"`
	Timestamp LatencyTimestamp `json:"timestamp"`
	Results   []LatencyResult  `json:"results"`
}

func (d LatencyAsset) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

func LatencyAssetJsonToStruct(v string) (asset LatencyAsset, err error) {
	err = json.Unmarshal([]byte(v), &asset)
	return asset, err
}

func CreateLatencyAsset(id string, latencyResults LatencyResults) LatencyAsset {
	return LatencyAsset{
		ID:        id,
		Source:    latencyResults.Source,
		Timestamp: latencyResults.Timestamp,
		Results:   latencyResults.Results,
	}
}

// We can create the ID because the Fabric App is letting us know the IP that's running this program
// Since we always ask which servers we should detect the latency anyway, we can create the ID here.
// Unlike with the resource usage collector, in which we always just blindly push our stats.
func CreateLatencyID(appType string, source string, timestamp LatencyTimestamp) string {
	if appType == "single_insert" {
		return source + "-" + DateFormatID(timestamp.TimeSeconds)
	} else if appType == "single_upsert" {
		return source
	} else {
		panic("APP_TYPE not implemented")
	}
}

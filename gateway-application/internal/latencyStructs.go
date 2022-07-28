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

func (d LatencyTargets) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

type LatencyTarget struct {
	Hostname     string `json:"hostname"`
	Hostport     string `json:"hostPort"`
	HostUser     string `json:"hostUser"`
	HostPassword string `json:"hostPassword"`
}

func LatencyTargetFromMap(properties Properties) LatencyTarget {
	return LatencyTarget{
		Hostname:     properties.Hostname,
		Hostport:     properties.HostPort,
		HostUser:     properties.HostUser,
		HostPassword: properties.HostPassword,
	}
}

func LatencyJsonToStrcut(v string) (targets LatencyTargets, err error) {
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
func JsonToLatencyAssetArray(v string) (assets []LatencyAsset, err error) {
	err = json.Unmarshal([]byte(v), &assets)
	return assets, err
}

func CreateLatencyAsset(id string, latencyResults LatencyResults) LatencyAsset {
	return LatencyAsset{
		ID:        id,
		Source:    latencyResults.Source,
		Timestamp: latencyResults.Timestamp,
		Results:   latencyResults.Results,
	}
}

func CreateLatencyID(appType string, source string, timestamp LatencyTimestamp) string {
	if appType == "single_insert" {
		return source + "-" + DateFormatID(timestamp.TimeSeconds)
	} else if appType == "single_upsert" {
		return source
	} else {
		panic("APP_TYPE not implemented")
	}
}

/////////////////////

type LatencyAnalysis struct {
	Hostname       string  `json:"hostname"`
	Target         string  `json:"target"`
	Duration       int     `json:"duration"`
	AverageLatency float64 `json:"averageLatency"`
	LatencyCount   int     `json:"latencyCount"`
	LatencySummary []int64 `json:"statSummary"`
}

func (d LatencyAnalysis) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

//func SummarizeLantecy(latencyAsset LatencyAsset) LatencyAnalysis {}

func AnalizeLatencySummary(latencyAnalysis LatencyAnalysis) LatencyAnalysis {
	if len(latencyAnalysis.LatencySummary) > 0 {
		var AverageLatency float64 = 0
		for _, summary := range latencyAnalysis.LatencySummary {
			AverageLatency += float64(summary)
		}
		latencyAnalysis.AverageLatency = AverageLatency / float64(len(latencyAnalysis.LatencySummary))
	}

	return latencyAnalysis
}

func JsonToLatencyAnalysisArray(v string) (assets []LatencyAnalysis, err error) {
	err = json.Unmarshal([]byte(v), &assets)
	return assets, err
}

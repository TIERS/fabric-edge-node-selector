package pkg

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/dmonteroh/fabric-distributed-resources/internal"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func filterServers(servers []internal.Asset, latencyAnalysis []internal.LatencyAnalysis) ([]internal.Asset, []internal.LatencyAnalysis) {
	filteredServers := []internal.Asset{}
	filteredAnalysis := []internal.LatencyAnalysis{}

	if len(servers) > 0 {
		for _, server := range servers {
			for _, analysis := range latencyAnalysis {
				if server.ID == analysis.Hostname {
					filteredServers = append(filteredServers, server)
					filteredAnalysis = append(filteredAnalysis, analysis)
				}
			}
		}
	}

	return filteredServers, filteredAnalysis
}

func combineAnalysis(target string, servers []internal.Asset, latencyAnalysis []internal.LatencyAnalysis, statAnalysis []internal.StatAnalysis) []internal.ServerSelection {
	selectionSlice := []internal.ServerSelection{}
	for _, server := range servers {
		tmpSel := internal.ServerSelection{}
		tmpSel.Asset = server
		tmpSel.Target = target

		for _, lat := range latencyAnalysis {
			if lat.Hostname == server.ID {
				tmpSel.AverageLatency = lat.AverageLatency
			}
		}
		for _, stat := range statAnalysis {
			if stat.Hostname == server.ID {
				tmpSel.CPUAverageUsage = stat.CPUAverageUsage
				tmpSel.MemoryUsePercentage = stat.MemoryUsePercentage
				tmpSel.ContainersRunning = stat.ContainersRunning
			}
		}

		selectionSlice = append(selectionSlice, tmpSel)
	}
	return selectionSlice
}

func sortSelection(items []internal.ServerSelection) {
	sort.Slice(items, func(i, j int) bool {
		var sortedByLatency, sortedByCPU bool

		sortedByLatency = items[i].AverageLatency < items[j].AverageLatency

		if items[i].AverageLatency == items[j].AverageLatency {
			sortedByCPU = items[i].CPUAverageUsage < items[j].CPUAverageUsage
			return sortedByCPU
		}
		return sortedByLatency
	})
}

func GetSelectedAssetHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	//contract := c.MustGet("selector").(*gateway.Contract)
	target := c.Param("target")
	minutes := c.Param("minutes")
	gpu, err := strconv.Atoi(c.Param("gpu"))
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Selecting SERVER for %s after %s minute analysis", target, minutes)

	// GET Servers
	var servers []internal.Asset
	if gpu == 1 {
		servers = ManualGPUServersInventoryHandler(c)
	} else {
		servers = ManualServersInventoryHandler(c)
	}

	// Check that there are servers

	// GET Latency Analysis of TARGET
	var latencyAnalysis []internal.LatencyAnalysis = ManualAnalysisTimeTarget(c)

	// FILTER Latency Analysis and Server List
	var filteredServers, filteredAnalysis = filterServers(servers, latencyAnalysis)

	// GET Resource Analysis for Server List
	var resourceAnalysis []internal.StatAnalysis = []internal.StatAnalysis{}
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(len(filteredServers))

	c1 := make(chan internal.StatAnalysis)

	go func() {
		for cr := range c1 {
			resourceAnalysis = append(resourceAnalysis, cr)
		}
	}()

	for _, server := range filteredServers {
		go func(server internal.Asset) {
			ManualSummaryAnalysisTime(c, server.ID, minutes, c1, waitGroup)
		}(server)
	}
	// Timestamp after operations
	waitGroup.Wait()
	//Without this sleep, the last result is skipped, don't know why
	time.Sleep(time.Microsecond * 15)

	// Combine Data into single Slice
	selectionObj := combineAnalysis(target, filteredServers, filteredAnalysis, resourceAnalysis)
	// Sort combined Data
	sortSelection(selectionObj)

	contract := c.MustGet("selector").(*gateway.Contract)
	storeSelection := internal.StoreSelection(selectionObj[0])

	fmt.Println(storeSelection)

	_, err = contract.SubmitTransaction("CreateAsset", storeSelection.String())
	if err != nil {
		panic(err.Error())
	}

	if len(selectionObj) == 0 {
		c.JSON(200, gin.H{"selected": nil, "options": nil})
	} else if len(selectionObj) == 1 {
		c.JSON(200, gin.H{"selected": selectionObj[0], "options": nil})
	} else {
		c.JSON(200, gin.H{"selected": selectionObj[0], "options": selectionObj[1:]})
	}

}

// CRUD
func GetAllSelectionsHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("selector").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToStoredSelectionArray(string(res))
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, readRes)
}

func GetSelectorHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("selector").(*gateway.Contract)
	asset := c.Param("id")

	res, err := contract.EvaluateTransaction("ReadAsset", asset)
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToStoredSelection(string(res))
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, readRes)
}

func GetAllSelectionTargetHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("selector").(*gateway.Contract)
	target := c.Param("target")

	res, err := contract.EvaluateTransaction("GetAllSelectionTarget", target)
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToStoredSelectionArray(string(res))
	if err != nil {
		panic(err.Error())
	} else if len(readRes) < 1 || readRes == nil {
		readRes = []internal.StoredSelection{}
	}

	c.JSON(200, readRes)
}

func GetAllSelectionServerHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("selector").(*gateway.Contract)
	asset := c.Param("asset")

	res, err := contract.EvaluateTransaction("GetAllSelectionServer", asset)
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToStoredSelectionArray(string(res))
	if err != nil {
		panic(err.Error())
	} else if len(readRes) < 1 || readRes == nil {
		readRes = []internal.StoredSelection{}
	}

	c.JSON(200, readRes)
}

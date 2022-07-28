package pkg

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/dmonteroh/fabric-distributed-resources/internal"
)

func UpsertResourceHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	appType := c.MustGet("APP_TYPE").(string)
	var clientIP string = c.ClientIP()
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	drcStats, err := internal.DrcJsonToStruct(string(jsonData))
	if err != nil {
		panic(err)
	}

	if appType == "single_insert" {
		stats := internal.ConvertToStorage(drcStats)
		stats.ID = clientIP + "-" + internal.DateFormatID(drcStats.Timestamp.TimeSeconds)
		stats.Hostname = clientIP
		createResource(c, stats)
	} else if appType == "single_upsert" {
		stats := internal.ConvertToStorage(drcStats)
		stats.ID = clientIP
		stats.Hostname = clientIP
		if resourceExists(c, stats.ID) {
			updateResource(c, stats)
		} else {
			createResource(c, stats)
		}
	} else {
		panic("APP_TYPE not implemented")
	}
}

func UpdateResourceHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)

	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	stats, err := internal.JsonToStoredStat(string(jsonData))
	if err != nil {
		panic(err)
	}

	if resourceExists(c, stats.ID) {
		updateResource(c, stats)
	} else {
		c.JSON(404, gin.H{"error": fmt.Sprintf("Error: Failed to submit, the key %s does not exist", stats.ID)})
	}
}

func GetAllResourcesHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("resources").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.ArrayStoredStat(string(res))
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, readRes)
}

func GetResourceHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("resources").(*gateway.Contract)
	asset := c.Param("asset")

	res, err := contract.EvaluateTransaction("ReadAsset", asset)
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToStoredStat(string(res))
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, readRes)
}

//GetAssetResourceListTime, , GetSummaryAnalysisTime
func ExecuteCustomQuery(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	payload, _ := ioutil.ReadAll(c.Request.Body)
	contractType := c.Param("contract")
	var response string

	if contractType != "resources" && contractType != "latency" && contractType != "inventory" {
		c.JSON(404, gin.H{"error": fmt.Sprintf("Error: Failed to submit, the contract '%s' does not exist", contractType)})
	} else {
		contract := c.MustGet(contractType).(*gateway.Contract)
		res, err := contract.EvaluateTransaction("ExecuteQuery", string(payload))
		if err != nil {
			panic(err.Error())
		}
		response = string(res)

		c.JSON(200, response)
	}

}

func GetAllAssetResourceList(c *gin.Context) {
	contract := c.MustGet("resources").(*gateway.Contract)
	device := c.Param("device")
	res, err := contract.EvaluateTransaction("GetAssetResource", device)

	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.ArrayStoredStat(string(res))
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, readRes)
}

func GetAssetResourceListTime(c *gin.Context) {
	contract := c.MustGet("resources").(*gateway.Contract)
	device := c.Param("device")
	minutes := c.Param("minutes")
	res, err := contract.EvaluateTransaction("GetAssetResourceListTime", device, minutes)

	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.ArrayStoredStat(string(res))
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, readRes)
}

func GetSummaryAnalysisTime(c *gin.Context) {
	contract := c.MustGet("resources").(*gateway.Contract)
	device := c.Param("device")
	minutes := c.Param("minutes")
	res, err := contract.EvaluateTransaction("GetSummaryAnalysisTime", device, minutes)

	if err != nil {
		panic(err.Error())
	}
	//log.Println("/n---------/n" + string(res) + "/n---------/n")
	readRes, err := internal.JsonToStatAnalysis(string(res))
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, readRes)
}

func ManualSummaryAnalysisTime(c *gin.Context, device string, minutes string, c1 chan internal.StatAnalysis, waitGroup *sync.WaitGroup) {
	contract := c.MustGet("resources").(*gateway.Contract)
	res, err := contract.EvaluateTransaction("GetSummaryAnalysisTime", device, minutes)

	if err != nil {
		panic(err.Error())
	}
	//log.Println("/n---------/n" + string(res) + "/n---------/n")
	readRes, err := internal.JsonToStatAnalysis(string(res))
	if err != nil {
		panic(err.Error())
	}
	c1 <- readRes
	defer waitGroup.Done()
}

// GetLastResourceSummary Not working due to CouchDB being flaky
// func GetLastResourceSummary(c *gin.Context) {
// 	defer internal.RecoverEndpoint(c)
// 	contract := c.MustGet("resources").(*gateway.Contract)
// 	asset := c.Param("asset")

// 	res, err := contract.EvaluateTransaction("GetLastResourceSummary", asset)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	readRes, err := internal.JsonToStatSummary(string(res))
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	c.JSON(200, readRes)
// }

// PRIVATE FUNCTIONS
func resourceExists(c *gin.Context, statIP string) bool {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("resources").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("AssetExists", statIP)
	if err != nil {
		panic(err.Error())
	}
	boolRes, err := strconv.ParseBool(string(res))
	if err != nil {
		panic(err.Error())
	}
	return boolRes
}

func updateResource(c *gin.Context, stats internal.StoredStat) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("resources").(*gateway.Contract)

	_, err := contract.SubmitTransaction("UpdateAsset", stats.ID, stats.String())
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, gin.H{"key": stats.ID})
}

func createResource(c *gin.Context, stats internal.StoredStat) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("resources").(*gateway.Contract)

	_, err := contract.SubmitTransaction("CreateAsset", stats.ID, stats.String())
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, gin.H{"key": stats.ID})
}

// FABRIC CALLS

// log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
// result, err := contract.SubmitTransaction("InitLedger")
// if err != nil {
// 	log.Fatalf("Failed to Submit transaction: %v", err)
// }
// log.Println(string(result))

// log.Println("--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")
// result, err = contract.EvaluateTransaction("GetAllAssets")
// if err != nil {
// 	log.Fatalf("Failed to evaluate transaction: %v", err)
// }
// log.Println(string(result))

// log.Println("--> Submit Transaction: CreateAsset, creates new asset with ID, color, owner, size, and appraisedValue arguments")
// result, err = contract.SubmitTransaction("CreateAsset", "localhost", "{\"timestamp\":{\"timeLocal\":\"2022-02-07T11:54:21.222970965Z\",\"timeSeconds\":1644234861,\"timeNano\":1644234861222970965},\"host\":{\"hostname\":\"426ede137da2\",\"uptime\":16061,\"boottime\":1644218801,\"platform\":\"alpine\",\"virtualizationSystem\":\"docker\",\"virtualizationRole\":\"guest\",\"hostid\":\"a72ab14c-76bf-ea11-8105-842afd4cdfcb\"},\"cpuStats\":{\"modelName\":\"Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz\",\"vendorId\":\"GenuineIntel\",\"averageUsage\":9.746240601477155,\"coreUsage\":[5.2631578947179465,9.523809523840459,9.523809523757965,9.999999999954525,4.999999999977263,5.2631578947179465,23.809523809286638,23.809523809503187,5.000000000022737,5.000000000022737,10.000000000045475,4.7619047618789825]},\"memStats\":{\"total\":67268472832,\"available\":60582289408,\"used\":5795409920},\"diskStats\":[{\"device\":\"/dev/sda5\",\"path\":\"/app\",\"label\":\"\",\"fstype\":\"ext4\",\"total\":78693273600,\"used\":63548289024,\"usedPercent\":85.1262654153894}],\"procStats\":{\"totalProcs\":2182,\"createdProcs\":176236,\"runningProcs\":1,\"blockedProcs\":0},\"dockerStats\":[{\"containerID\":\"6440e6be7720f2fd842aea47b6a87de484414561a212acd71df7a6e04a915cee\",\"name\":\"/dev-peer0.org2.example.com-basic_1.0-5f042bbcb3e3b1b4b6e6a25f30f746f263614a8b838865b6f72deb9cbd8ab981\",\"image\":\"dev-peer0.org2.example.com-basic_1.0-5f042bbcb3e3b1b4b6e6a25f30f746f263614a8b838865b6f72deb9cbd8ab981-9e2ae745b02b13626fa3a2f2d71e307d9bfe37fc81485ce810802cb0859f0872\",\"status\":\"Up 5 minutes\",\"State\":\"running\"},{\"containerID\":\"f53321eee853444cf550909a339ed5182d74e8490a6dc7e64fa16f26fdde5154\",\"name\":\"/dev-peer0.org1.example.com-basic_1.0-5f042bbcb3e3b1b4b6e6a25f30f746f263614a8b838865b6f72deb9cbd8ab981\",\"image\":\"dev-peer0.org1.example.com-basic_1.0-5f042bbcb3e3b1b4b6e6a25f30f746f263614a8b838865b6f72deb9cbd8ab981-1b5f8eb9971213e57fd41c10cac95d657f48b197f27f836e8cc606f0a4cf27fc\",\"status\":\"Up 5 minutes\",\"State\":\"running\"},{\"containerID\":\"670c610a8b5641cd060585be5d607fe1f0365974d4c6c02f915feb5ebc64b1a5\",\"name\":\"/cli\",\"image\":\"hyperledger/fabric-tools:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"8d933fb4377bb80f5eead71ccf160912e53abf54f62f96290a625cca89ea3800\",\"name\":\"/peer0.org2.example.com\",\"image\":\"hyperledger/fabric-peer:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"24069d7b1df3f28f556276af70be52d09dc3cf323581a327d7a7ad3084ca1396\",\"name\":\"/peer0.org1.example.com\",\"image\":\"hyperledger/fabric-peer:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"c9d549992b95170d7ae750c5dc6bf464b03c37188d33476dbb0e87126d52af6e\",\"name\":\"/orderer.example.com\",\"image\":\"hyperledger/fabric-orderer:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"f29776a0078c97076ca2ea5b4d3121bb4c4e307f676a6ba88981b7bb13195c26\",\"name\":\"/couchdb1\",\"image\":\"couchdb:3.1.1\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"3493a420a5df5b12efa116c279d57f46aecbbfaf0a820c0208a2462b7a068726\",\"name\":\"/couchdb0\",\"image\":\"couchdb:3.1.1\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"8ee4ad099a3e4446274444df66cfe3cb38aa257cfd5cf1b0012bec9ba71ed284\",\"name\":\"/ca_orderer\",\"image\":\"hyperledger/fabric-ca:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"d1b7e077111df8525e738388e70d99b5e75d8d611b3e6b1c8994a90f963d7832\",\"name\":\"/ca_org2\",\"image\":\"hyperledger/fabric-ca:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"674b38d89aa0d381c1b7476d773d0b346081a3c7202649c0aa81ad5b34661b7a\",\"name\":\"/ca_org1\",\"image\":\"hyperledger/fabric-ca:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"426ede137da2f30fa6aede67b757b6545a51d51e30a673e9cbf8cbebe4b4f0b2\",\"name\":\"/distributed-resource-collector\",\"image\":\"distributed-resource-collector\",\"status\":\"Up 3 seconds\",\"State\":\"running\"}]}")
// if err != nil {
// 	log.Fatalf("Failed to Submit transaction: %v", err)
// }
// log.Println(string(result))

// log.Println("--> Evaluate Transaction: ReadAsset, function returns an asset with a given assetID")
// result, err = contract.EvaluateTransaction("ReadAsset", "localhost")
// if err != nil {
// 	log.Fatalf("Failed to evaluate transaction: %v\n", err)
// }
// log.Println(string(result))

// log.Println("--> Evaluate Transaction: AssetExists, function returns 'true' if an asset with given assetID exist")
// result, err = contract.EvaluateTransaction("AssetExists", "localhost")
// if err != nil {
// 	log.Fatalf("Failed to evaluate transaction: %v\n", err)
// }
// log.Println(string(result))

// log.Println("--> Submit Transaction: TransferAsset asset1, transfer to new owner of Tom")
// _, err = contract.SubmitTransaction("TransferAsset", "asset1", "Tom")
// if err != nil {
// 	log.Fatalf("Failed to Submit transaction: %v", err)
// }

// log.Println("--> Evaluate Transaction: ReadAsset, function returns 'asset1' attributes")
// result, err = contract.EvaluateTransaction("ReadAsset", "asset1")
// if err != nil {
// 	log.Fatalf("Failed to evaluate transaction: %v", err)
// }
// log.Println(string(result))
// log.Println("============ application-golang ends ============")

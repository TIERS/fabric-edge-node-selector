package pkg

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/dmonteroh/fabric-distributed-resources/internal"
)

func GetAllInventoryHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("inventory").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToAssetArray(string(res))
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, readRes)
}

func GetInventoryHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("inventory").(*gateway.Contract)
	asset := c.Param("asset")

	res, err := contract.EvaluateTransaction("ReadAsset", asset)
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToAsset(string(res))
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, readRes)
}

func UpdateInventoryHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("inventory").(*gateway.Contract)

	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	inventory, err := internal.JsonToAsset(string(jsonData))
	if err != nil {
		panic(err)
	}

	_, err = contract.SubmitTransaction("UpdateAsset", inventory.String())
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, gin.H{"key": inventory.ID})
}

func CreateInventoryHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("inventory").(*gateway.Contract)

	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	inventory, err := internal.JsonToAsset(string(jsonData))
	if err != nil {
		panic(err)
	}

	_, err = contract.SubmitTransaction("CreateAsset", inventory.String())
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, gin.H{"key": inventory.ID})
}

func GetServersInventoryHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("inventory").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetServerAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToAssetArray(string(res))
	if err != nil {
		panic(err.Error())
	} else if len(readRes) < 1 || readRes == nil {
		readRes = []internal.Asset{}
	}

	c.JSON(200, readRes)
}

func GetGPUServersInventoryHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("inventory").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetServerGPUAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToAssetArray(string(res))
	if err != nil {
		panic(err.Error())
	} else if len(readRes) < 1 || readRes == nil {
		readRes = []internal.Asset{}
	}

	c.JSON(200, readRes)
}

func ManualServersInventoryHandler(c *gin.Context) []internal.Asset {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("inventory").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetServerAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToAssetArray(string(res))
	if err != nil {
		panic(err.Error())
	} else if len(readRes) < 1 {
		return []internal.Asset{}
	}

	return readRes
}

func ManualGPUServersInventoryHandler(c *gin.Context) []internal.Asset {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("inventory").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetServerGPUAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToAssetArray(string(res))
	if err != nil {
		panic(err.Error())
	} else if len(readRes) < 1 {
		return []internal.Asset{}
	}

	return readRes
}

func GetRobotInventoryHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("inventory").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetRobotAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToAssetArray(string(res))
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, readRes)
}

func GetSensorInventoryHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("inventory").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetSensorAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToAssetArray(string(res))
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, readRes)
}

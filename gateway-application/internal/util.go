package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		//fmt.Println(value)
		return value
	}
	return fallback
}

//// Removed as is currently unused
// func KeysInStringMap(mapObj map[string]string, keys []string) bool {
// 	for _, key := range keys {
// 		_, ok := mapObj[key]
// 		if !ok {
// 			return false
// 		}
// 	}
// 	return true
// }

func DateFormatID(d int64) string {
	t := time.Unix(d, 0)
	layout := "2006-01-02T15:04:05"
	return t.Format(layout)
}

// Adds every key and value in map to the gin context as middleware. Allows access to these variables from inside the handlers
func EnviromentMiddleware(variables map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, value := range variables {
			if key != "" && value != "" {
				c.Set(key, value)
				//c.Next()
			}
		}
	}
}

// Adds the fabric contract to the gin context as middleware.
func ContractMiddleware(contractName string, contract *gateway.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(contractName, contract)
		//c.Next()
	}
}

func RecoverEndpoint(c *gin.Context) {
	if err := recover(); err != nil {
		msg := "Error: [Recovered] "
		switch errType := err.(type) {
		case string:
			msg += err.(string)
		case *json.SyntaxError:
			msg += errType.Error()
		default:
		}
		fmt.Println(msg)
		c.JSON(400, gin.H{"error": msg})
	}
}

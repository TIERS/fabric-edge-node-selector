package internal

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shomali11/parallelizer"
)

// GetEnv is a simple function that will give you a fallback value in case the environment variable is empty. Sort of a default option.
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		//fmt.Println(value)
		return value
	}
	return fallback
}

func UrlMaker(protocol string, hostname string, endpoint string) string {
	return strings.Join([]string{protocol, strings.Join([]string{hostname, endpoint}, "/")}, "://")
}

// Transforming the clasic Unix time in seconds to human readable date
func DateFormatID(d int64) string {
	t := time.Unix(d, 0)
	layout := "2006-01-02T15:04:05"
	return t.Format(layout)
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
}

// Return a slice of unique strings. This function is not implemented in default into GO unlike other languages.
func UniqueString(slice []string) (unique []string) {
	for _, v := range slice {
		skip := false
		for _, u := range unique {
			if v == u {
				skip = true
				break
			}
		}
		if !skip {
			unique = append(unique, v)
		}
	}
	return unique
}

// Return a boolean if at least one of substrings is contained into a main string.
func CustomContains(str string, subStrings ...string) bool {
	if len(subStrings) == 0 {
		return true
	}

	for _, subString := range subStrings {
		if strings.Contains(str, subString) {
			return true
		}
	}
	return false
}

// There are multiple ways to check if we're running from inside a Docker Container, this default operation is good enough for what we need.
func InDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

// Adds the parallelization group to the gin context as middleware. Allows access to these variables from inside the handlers
// It's currently unknown if creating a new parallize group creates unnecessary overhead.
// Marked for deletion in next release -> NOT USING PARALLELIZER, GO ROUTINES AND CHANNELS ARE ENOUGH
func GroupMiddleware(group *parallelizer.Group) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("latencyGroup", group)
		//c.Next()
	}
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

// Reuse code to return ERROR message
// RESTFUL implementation is incomplete, as all results should use the same structure (error should always be present, and results should be interfaced through the same key)
func RecoverEndpoint(c *gin.Context) {
	if err := recover(); err != nil {
		msg := "Error: [Recovered] "
		switch errType := err.(type) {
		case string:
			msg += err.(string)
		case error:
			msg += errType.Error()
		default:
		}
		fmt.Println(msg)
		c.JSON(400, gin.H{"error": msg})
	}
}

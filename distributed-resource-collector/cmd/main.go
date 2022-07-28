package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"

	"github.com/dmonteroh/distributed-resource-collector/internal"
	"github.com/dmonteroh/distributed-resource-collector/pkg"
)

func main() {
	// I was using RUNTIME to determine if parallelization started correctly, only for debugging
	// fmt.Println("Version", runtime.Version())
	// fmt.Println("NumCPU", runtime.NumCPU())
	// fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))

	// Create a GINGONIC http server
	r := gin.Default()

	// ENVIROMENTAL VARIABLES
	execMode := internal.GetEnv("EXEC_MODE", "DEBUG")
	listenPort := internal.GetEnv("INTERNAL_PORT", "8081")
	appProtocol := internal.GetEnv("APP_PROTOCOL", "http")
	appIP := internal.GetEnv("APP_IP", "localhost:8080")
	collectorUrl := internal.GetEnv("APP_URL", "collector")
	targetsUrl := internal.GetEnv("TARGETS_URL", "latency/servers/targets")
	latencyUrl := internal.GetEnv("LATENCY_URL", "latency")
	collectorApp := internal.UrlMaker(appProtocol, appIP, collectorUrl)
	targetsApp := internal.UrlMaker(appProtocol, appIP, targetsUrl)
	latencyApp := internal.UrlMaker(appProtocol, appIP, latencyUrl)
	appCron, _ := strconv.Atoi(internal.GetEnv("APP_CRON", "30"))
	heartbeat, _ := strconv.ParseBool(internal.GetEnv("HEARTBEAT", "true"))

	// MAP VARIABLES INTO MAP
	variables := map[string]string{
		"EXEC_MODE":     execMode,
		"COLLECTOR_APP": collectorApp,
		"LATENCY_APP":   latencyApp,
		"TARGETS_APP":   targetsApp,
	}

	// SAVE VARIABLES INSIDE GIN CONTEXT
	r.Use(internal.EnviromentMiddleware(variables))
	//r.Use(internal.GroupMiddleware(latencyGroup))

	// HTTP SERVER ROUTES
	r.GET("/heartbeat", pkg.HeartbeatEndpoint)
	r.GET("/latency", pkg.LatencyEndpoint)
	r.POST("/latency", pkg.ManualLatencyEndpoint)

	// HEARTBEAT AND LATENCY AUTO POSTING
	if heartbeat {
		fmt.Println("INITIATE HEARTBEAT SCHEDULER")
		cron := gocron.NewScheduler(time.Local)
		pkg.HeartbeatCron(cron, appCron, collectorApp, execMode)
		pkg.LatencyCron(cron, appCron, targetsApp, latencyApp, execMode)
		cron.StartAsync()
	}

	// Start listening on the desired port (similar to ros' spin, "blocks" thread)
	r.Run(":" + listenPort)
}

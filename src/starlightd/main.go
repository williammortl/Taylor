package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/williammortl/Taylor/src/starlightd/pkg/handlers"
	telemetryLib "github.com/williammortl/Taylor/src/starlightd/pkg/telemetry"
)

const subsystemName string = "starlightd"
const prometheusEndpoint string = "/metrics"
const prometheusPort int = 8888
const restPort int = 1604

func main() {
	var telemetry *telemetryLib.Telemetry = telemetryLib.InitializeTelemetry(subsystemName)
	telemetryLib.StartTelemetryEndpoint(prometheusEndpoint, prometheusPort)
	telemetry.LogTrace("main", "test", "hi")

	// REST endpoints
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/launcher", handlers.LauncherHandlerPost(telemetry)).Methods("POST")
	router.HandleFunc("/launcher/{runID}", handlers.LauncherHandlerGet(telemetry)).Methods("GET")

	// launch REST server
	http.ListenAndServe(fmt.Sprintf(":%v", restPort), router)
}

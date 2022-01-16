package main

import (
	"bufio"
	"os"

	telemetryLib "github.com/williammortl/Taylor/src/starlightd/pkg/telemetry"
)

const subsystemName string = "starlightd"
const prometheusEndpoint string = "/metrics"
const prometheusPort int = 8888

func main() {
	var telemetry *telemetryLib.Telemetry = telemetryLib.InitializeTelemetry(subsystemName,
		prometheusEndpoint, prometheusPort)
	telemetry.LogTrace("main", "test", "hi")

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

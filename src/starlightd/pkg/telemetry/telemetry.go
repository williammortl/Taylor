package telemetry

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// namespaceNameDefault is the default namespace for this project
const namespaceNameDefault = "Taylor"

// exeuctionTimeStart base time == 0
const exeuctionTimeStart = 0

// exeuctionTimeWidth is the width of a bucket in the histogram, here it is 1m
const exeuctionTimeWidth = 60

// executionTimeBuckets is the number of buckets, here it 20 minutes worth of 1m buckets
const executionTimeBuckets = 20

var (
	statusCounter     *prometheus.CounterVec
	durationHistogram *prometheus.HistogramVec
	endpointStarted   bool = false
)

// TelemetryClient contains the functions for Telemetry
type TelemetryClient interface {
	LogDuration(componentName string, durationName string, durationInSecs float64)
	LogTrace(componentName string, typeTrace string, message string)
	LogInfo(componentName string, typeInfo string, message string)
	LogWarning(componentName string, typeWarning string, message string)
	LogError(componentName string, typeError string, err error)
}

// Telemetry contains data for the TelemetryClient
type Telemetry struct {
	NamespaceName string
	SubsystemName string
	Endpoint      string
	Port          int
}

// InitializeTelemetry initializes a TelemetryFactory client
func InitializeTelemetry(subsystemName string) *Telemetry {

	var telemetryConfig Telemetry = Telemetry{
		NamespaceName: namespaceNameDefault,
		SubsystemName: subsystemName,
	}

	// initialize global metrics if neccessary
	if statusCounter == nil {

		statusCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespaceNameDefault,
				Subsystem: telemetryConfig.SubsystemName,
				Name:      "Status",
				Help:      "Status messages",
			},
			[]string{
				"component",
				"level",
				"type",
				"message",
			},
		)
		prometheus.MustRegister(statusCounter)

		durationHistogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespaceNameDefault,
				Subsystem: telemetryConfig.SubsystemName,
				Name:      "ExecutionTime",
				Buckets: prometheus.LinearBuckets(
					exeuctionTimeStart,
					exeuctionTimeWidth,
					executionTimeBuckets),
				Help: "Time to execute",
			},
			[]string{
				"component",
				"durationName",
			},
		)
		prometheus.MustRegister(durationHistogram)
	}

	return &telemetryConfig
}

// StartTelemetryEndpoint starts the telemetry service endpoint
func StartTelemetryEndpoint(endpoint string, port int) {
	if !endpointStarted {
		endpointStarted = true
		go goStartLoggingEndpoint(endpoint, port)
	}
}

// goStartTelemetryEndpoint is a go routine for starting the telemetry service endpoint
func goStartLoggingEndpoint(endpoint string, port int) {
	http.Handle(endpoint, promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// LogDuration logs the duration of an operation in seconds
func (t *Telemetry) LogDuration(componentName string, durationName string, durationInSecs float64) {
	t.printLogging(componentName, "Duration", durationName, fmt.Sprintf("%f(s)", durationInSecs))
	durationHistogram.WithLabelValues(componentName, durationName).Observe(durationInSecs)
}

// LogTrace logs a trace
func (t *Telemetry) LogTrace(componentName string, typeTrace string, message string) {
	t.printLogging(componentName, "Trace", typeTrace, message)
}

// LogInfo logs an informational message
func (t *Telemetry) LogInfo(componentName string, typeInfo string, message string) {
	t.printLogging(componentName, "Info", typeInfo, message)
	statusCounter.WithLabelValues(componentName, "Info", typeInfo, message).Inc()
}

// LogWarning logs a warning
func (t *Telemetry) LogWarning(componentName string, typeWarning string, message string) {
	t.printLogging(componentName, "Warning", typeWarning, message)
	statusCounter.WithLabelValues(componentName, "Warning", typeWarning, message).Inc()
}

// LogErrorByInstance logs an error
func (t *Telemetry) LogError(componentName string, typeError string, err error) {
	t.printLogging(componentName, "Error", typeError, err.Error())
	statusCounter.WithLabelValues(componentName, "Error", typeError, err.Error()).Inc()
}

// printLogging prints the logging message
func (t *Telemetry) printLogging(componentName string, level string, typeOf string, message string) {
	log.Printf("%s|%s|%s|%s::%s - %s", t.NamespaceName, t.SubsystemName, componentName, level, typeOf, message)
}

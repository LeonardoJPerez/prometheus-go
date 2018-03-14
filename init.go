package prometheusTelemetry

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	appNamespace = "gawkbox_"

	currentQueries               prometheus.Gauge
	dbTransactionDurationHistVec *prometheus.HistogramVec
	httpRequestsCounterVector    *prometheus.CounterVec
	httpRequestsDurationHistVec  *prometheus.HistogramVec
)

const (
	dbCurrentQueries         = "current_queries"
	dbTransactionHistVecName = "transaction_latency"
	requestCounterName       = "http_requests_total"
	requestDurationName      = "http_request_duration_seconds"

	echoSubsystem = "echo"
	dbSubsystem   = "db"
)

func initialize() {
	currentQueries = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appNamespace,
		Subsystem: dbSubsystem,
		Name:      dbCurrentQueries,
		Help:      "The current number of database queries being executed or waiting.",
	})

	dbTransactionDurationHistVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: appNamespace,
		Subsystem: dbSubsystem,
		Name:      dbTransactionHistVecName,
		Help:      "The transaction time while querying the database.",
		Buckets:   prometheus.DefBuckets,
	},
		[]string{"operation"},
	)

	httpRequestsCounterVector = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: appNamespace,
		Subsystem: echoSubsystem,
		Name:      requestCounterName,
		Help:      "HTTP requests processed, partitioned by status code, HTTP method and handler.",
	},
		[]string{"code", "method", "handler"},
	)

	httpRequestsDurationHistVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: appNamespace,
		Subsystem: echoSubsystem,
		Name:      requestDurationName,
		Help:      "Duration of request partitioned by method and handler.",
		Buckets:   prometheus.DefBuckets,
	},
		[]string{"method", "handler"},
	)

	// Registering metrics.
	prometheus.MustRegister(currentQueries)
	prometheus.MustRegister(dbTransactionDurationHistVec)
	prometheus.MustRegister(httpRequestsCounterVector)
	prometheus.MustRegister(httpRequestsDurationHistVec)
}

// SetupTelemetry :
func SetupTelemetry(namespace string) {
	if namespace != "" {
		appNamespace = fmt.Sprintf("%s_%s", appNamespace, namespace)
	}

	initialize()
}

// IncrementCurrentDbQueries increments the current query counter by 1.
func IncrementCurrentDbQueries() {
	if currentQueries == nil {
		return
	}

	currentQueries.Inc()
}

// DecrementCurrentDbQueries decrements the current query counter by 1.
func DecrementCurrentDbQueries() {
	if currentQueries == nil {
		return
	}

	currentQueries.Dec()
}

// AddRequestMetric creates a series point a specific request.
func AddRequestMetric(code, method, handler string) {
	if httpRequestsCounterVector == nil {
		return
	}

	httpRequestsCounterVector.WithLabelValues(code, method, handler).Inc()
}

// ObserveTransaction :
func ObserveTransaction(operation string, begin time.Time) {
	if dbTransactionDurationHistVec == nil {
		return
	}

	dbTransactionDurationHistVec.
		WithLabelValues(operation).
		Observe(time.Since(begin).Seconds())
}

// ObserveRequestDuration records the time lapse from a starting time object.
func ObserveRequestDuration(method, handler string, begin time.Time) {
	if httpRequestsDurationHistVec == nil {
		return
	}

	httpRequestsDurationHistVec.
		WithLabelValues(method, handler).
		Observe(time.Since(begin).Seconds())
}

package metrics

import (
	"context"
	"time"

	httpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/authz/http"
	"github.com/prometheus/client_golang/prometheus"
	ctrlmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	httpRequestsMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "authz_http_server_validating_requests",
			Help: "can be used to track the number of http requests",
		},
		[]string{"method", "host", "path", "schema", "status"},
	)
	httpRequestsErrorMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "authz_http_server_validating_request_errors",
			Help: "can be used to track the number of http request errors",
		},
		[]string{"method", "host", "path", "schema"},
	)
	httpDurationMetric = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "authz_http_server_validating_requests_duration_seconds",
			Help: "can be used to track the latencies (in seconds) associated with the entire http request.",
		},
		[]string{"method", "host", "path", "schema", "status"},
	)
)

func init() {
	ctrlmetrics.Registry.MustRegister(httpRequestsMetric, httpDurationMetric, httpRequestsErrorMetric)
}

func RecordHTTPRequest(ctx context.Context, startTime time.Time, req httpcel.CheckRequest, res *httpcel.CheckResponse) {
	var status string
	if res.Denied != nil {
		status = "denied"
	} else {
		status = "ok"
	}

	httpRequestsMetric.WithLabelValues(
		req.Attributes.Method,
		req.Attributes.Host,
		req.Attributes.Path,
		req.Attributes.Scheme,
		status,
	).Inc()

	if httpDurationMetric != nil {
		defer func() {
			latency := float64(time.Since(startTime))
			httpDurationMetric.WithLabelValues(
				req.Attributes.Method,
				req.Attributes.Host,
				req.Attributes.Path,
				req.Attributes.Scheme,
				status,
			).Observe(latency)
		}()
	}
}

func RecordHTTPRequestError(ctx context.Context, req httpcel.CheckRequest, err error) {
	httpRequestsErrorMetric.WithLabelValues(
		req.Attributes.Method,
		req.Attributes.Host,
		req.Attributes.Path,
		req.Attributes.Scheme,
	).Inc()
}

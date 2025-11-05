package metrics

import (
	"context"
	"fmt"
	"time"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/prometheus/client_golang/prometheus"
	ctrlmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	envoyRequestsMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "authz_envoy_server_validating_requests",
			Help: "can be used to track the number of envoy requests",
		},
		[]string{"method", "host", "path", "schema", "status"},
	)
	envoyRequestsErrorMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "authz_envoy_server_validating_request_errors",
			Help: "can be used to track the number of envoy request errors",
		},
		[]string{"method", "host", "path", "schema"},
	)
	envoyDurationMetric = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "authz_envoy_server_validating_requests_duration_seconds",
			Help: "can be used to track the latencies (in seconds) associated with the entire envoy request.",
		},
		[]string{"method", "host", "path", "schema", "status"},
	)
)

func init() {
	ctrlmetrics.Registry.MustRegister(envoyRequestsMetric, envoyDurationMetric, envoyRequestsErrorMetric)
}

func RecordEnvoyRequest(ctx context.Context, startTime time.Time, req *authv3.CheckRequest, res *authv3.CheckResponse) {
	var status string
	if res.Status != nil {
		status = fmt.Sprint(res.Status.Code)
	}

	envoyRequestsMetric.WithLabelValues(
		req.Attributes.Request.Http.Method,
		req.Attributes.Request.Http.Host,
		req.Attributes.Request.Http.Path,
		req.Attributes.Request.Http.Scheme,
		status,
	).Inc()

	defer func() {
		latency := float64(time.Since(startTime))

		envoyDurationMetric.WithLabelValues(
			req.Attributes.Request.Http.Method,
			req.Attributes.Request.Http.Host,
			req.Attributes.Request.Http.Path,
			req.Attributes.Request.Http.Scheme,
			status,
		).Observe(latency)
	}()
}

func RecordEnvoyRequestError(ctx context.Context, req *authv3.CheckRequest, err error) {
	envoyRequestsErrorMetric.WithLabelValues(
		req.Attributes.Request.Http.Method,
		req.Attributes.Request.Http.Host,
		req.Attributes.Request.Http.Path,
		req.Attributes.Request.Http.Scheme,
	).Inc()
}

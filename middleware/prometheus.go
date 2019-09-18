package middleware

import (
	"context"
	"micro/meta"
	"micro/middleware/prometheus"
	"time"
)

var (
	DefaultServiceMetrics = prometheus.NewServerMetrics()
)

func PrometheusServerMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		serverMeta := meta.GetServerMeta(ctx)
		DefaultServiceMetrics.IncrRequest(ctx, serverMeta.ServiceName, serverMeta.Method)

		startTime := time.Now()
		resp, err = next(ctx, req)

		DefaultServiceMetrics.IncrCode(ctx, serverMeta.ServiceName, serverMeta.Method, err)
		DefaultServiceMetrics.Latency(ctx, serverMeta.ServiceName, serverMeta.Method, time.Since(startTime).Nanoseconds()/1000)
		return
	}
}

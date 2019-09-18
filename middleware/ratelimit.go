package middleware

import (
	"context"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Allow() bool
}

func NewRateLimitMiddleware(l Limiter) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			glog.Info("access limit")
			allow := l.Allow()
			if !allow {
				glog.Info("limit deny")
				err = status.Error(codes.ResourceExhausted, "rate limited")
				return
			}
			glog.Info("limit pass")

			return next(ctx, req)
		}
	}
}

package middleware

import (
	"context"
	"time"

	"github.com/golang/glog"
)

/*
func init() {
	Use(CostMiddleware)
}
*/

func CostMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		startTimeNano := time.Now().UnixNano()
		resp, err = next(ctx, req)
		endTimeNano := time.Now().UnixNano()
		cost := (endTimeNano - startTimeNano) / 1000
		glog.Info("cost ", cost/1000, " ms")
		return
	}
}

package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMiddleware(t *testing.T) {

	costmidware := func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			startTimeNano := time.Now().UnixNano()
			resp, err = next(ctx, req)
			endTimeNano := time.Now().UnixNano()
			cost := (endTimeNano - startTimeNano) / 1000
			fmt.Println("costmidware ", cost/1000, " ms")
			return
		}
	}

	middleware1 := func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			fmt.Printf("middleware 1 start\n")
			num := rand.Intn(2)
			if num <= 2 {
				err = fmt.Errorf("this is request is not allow")
				return
			}
			resp, err = next(ctx, req)
			if err != nil {
				return
			}
			fmt.Printf("middleware1 end\n")
			return
		}
	}

	middleware2 := func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			fmt.Printf("middleware 2 start\n")

			resp, err = next(ctx, req)
			if err != nil {
				return
			}
			fmt.Printf("middleware2 end\n")
			return
		}
	}

	outer := func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			fmt.Printf("outer  start\n")
			resp, err = next(ctx, req)
			if err != nil {
				return
			}
			fmt.Printf("outer end\n")
			return
		}
	}

	proc := func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		fmt.Printf("req process start\n")
		fmt.Printf("req process end\n")
		return
	}

	chain := Chain(outer, middleware1, middleware2, costmidware)
	chainFunc := chain(proc)

	resp, err := chainFunc(context.Background(), "test")
	fmt.Printf("resp:%#v, err:%v\n", resp, err)
}

package concurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type MyContext interface {
	Say() string
	DoThing() error
}

type emptyCtx struct {
}

func (ctx *emptyCtx) Say() string {
	return ""
}

func (ctx *emptyCtx) DoThing() error {
	return nil
}

type valueCtx struct {
	MyContext
	key, value interface{}
}

func TODO() MyContext {
	return new(emptyCtx)
}

func (ctx *valueCtx) Say() string {
	return "说一段相声"
}

func TestCtx(t *testing.T) {
	ctx := TODO()
	vc := valueCtx{
		MyContext: ctx,
	}
	fmt.Println(vc.Say())
	_ = vc.DoThing()
}

func HandleRequest(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Handle request done")
			return
		default:
			fmt.Println("Handle request running, parameter: ", ctx.Value("parameter"))
			time.Sleep(2 * time.Second)
		}
	}
}

func TestValueCtx(t *testing.T) {
	ctx := context.WithValue(context.Background(), "parameter", "1")
	go HandleRequest(ctx)
	time.Sleep(10 * time.Second)
}

func HandleComplexRequest(ctx context.Context) {
	go WriteRedis(ctx)
	go WriteDatabase(ctx)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("HandleComplexRequest done")
			return
		default:
			fmt.Println("HandleComplexRequest running")
			time.Sleep(2 * time.Second)
		}
	}
}

func WriteRedis(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("WriteRedis done")
			return
		default:
			fmt.Println("WriteRedis running")
			time.Sleep(2 * time.Second)
		}
	}
}

func WriteDatabase(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("WriteDatabase done")
			return
		default:
			fmt.Println("WriteDatabase running")
			time.Sleep(2 * time.Second)
		}
	}
}

func TestCancelCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go HandleComplexRequest(ctx)

	time.Sleep(5 * time.Second)
	fmt.Println("It's time to stop all sub goroutines!")
	cancel()

	time.Sleep(5 * time.Second)
}

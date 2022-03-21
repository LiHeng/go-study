package mapreduce

import (
	"context"
)

// GenerateFunc 数据生产func
type GenerateFunc func(source <-chan interface{})

// MapperFunc 数据加工func
// item - 生产出来的数据
// writer - 调用writer.Write()可以将加工后的向后传递至reducer
// cancel - 终止流程func
type MapperFunc func(item interface{}, writer Writer, cancel func(error))

// ReducerFunc 数据聚合func
// pipe - 加工出来的数据
// writer - 调用writer.Write()可以将聚合后的数据返回给用户
// cancel - 终止流程func
type ReducerFunc func(pipe <-chan interface{}, writer Writer, cancel func(error))

type Writer interface {
	Write(v interface{})
}

type (
	// ForEachFunc is used to do element processing, but no output.
	ForEachFunc func(item interface{})
	// MapFunc is used to do element processing and write the output to writer.
	MapFunc func(item interface{}, writer Writer)
	// VoidReducerFunc is used to reduce all the mapping output, but no output.
	// Use cancel func to cancel the processing.
	VoidReducerFunc func(pipe <-chan interface{}, cancel func(error))
	// Option defines the method to customize the mapreduce.
	Option func(opts *mapReduceOptions)

	mapperContext struct {
		ctx       context.Context
		mapper    MapFunc
		source    <-chan interface{}
		panicChan *onceChan
		collector chan<- interface{}
		doneChan  <-chan struct{}
		workers   int
	}

	mapReduceOptions struct {
		ctx     context.Context
		workers int
	}
)

type onceChan struct {
	channel chan interface{}
	wrote   int32
}

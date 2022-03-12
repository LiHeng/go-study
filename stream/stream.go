package stream

import (
	"sort"
	"sync"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

type Stream struct {
	source <-chan interface{}
}

type rxOptions struct {
	unlimitedWorkers bool
	workers          int
}

type GenerateFunc func(source <-chan interface{})

type KeyFunc func(item interface{}) interface{}

// ForEachFunc 对每个item执行操作
type ForEachFunc func(item interface{})

// FilterFunc 过滤函数
type FilterFunc func(item interface{}) bool

// Option defines the method to customize a Stream.
type Option func(opts *rxOptions)

// WalkFunc defines the method to walk through all the elements in a Stream.
type WalkFunc func(item interface{}, pipe chan<- interface{})

// LessFunc defines the method to compare the elements in a Stream.
type LessFunc func(a, b interface{}) bool

type PlaceholderType struct{}

var Placeholder PlaceholderType

// Range 创建阶段/数据获取阶段
func Range(source <-chan interface{}) Stream {
	return Stream{
		source: source,
	}
}

// Just 通过可变参数模式创建 stream
func Just(items ...interface{}) Stream {
	// 带缓冲的channel
	source := make(chan interface{}, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)
	return Range(source)
}

// From 通过函数创建 stream
func From(generate GenerateFunc) Stream {
	source := make(chan interface{})
	GoSafe(func() {
		defer close(source)
		generate(source)
	})
	return Range(source)
}

// Concat 拼接stream
// func Concat(s Stream, other ...Stream) Stream {
//
// }

func (s Stream) Distinct(keyFunc KeyFunc) Stream {
	source := make(chan interface{})
	GoSafe(func() {
		defer close(source)
		keys := make(map[interface{}]PlaceholderType)
		for item := range s.source {
			// 自定义去重逻辑
			key := keyFunc(item)
			// 如果key不存在,则将数据写入新的channel
			if _, ok := keys[key]; !ok {
				source <- item
				keys[key] = Placeholder
			}
		}
	})
	return Range(source)
}

// Filter 新的 Stream 中 channel 里面的数据顺序是随机的。
func (s Stream) Filter(filterFunc FilterFunc, opts ...Option) Stream {
	return s.Walk(func(item interface{}, pipe chan<- interface{}) {
		if filterFunc(item) {
			pipe <- item
		}
	}, opts...)
}

func (s Stream) Walk(fn WalkFunc, opts ...Option) Stream {
	option := buildOptions(opts...)
	if option.unlimitedWorkers {
		return s.walkUnlimited(fn, option)
	}

	return s.walkLimited(fn, option)
}

func (s Stream) Sort(fn LessFunc) Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return fn(items[i], items[j])
	})
	return Just(items...)
}

func (s Stream) walkUnlimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan interface{}, option.workers)

	go func() {
		var wg sync.WaitGroup
		for item := range s.source {
			// important, used in another goroutine
			val := item
			wg.Add(1)
			// better to safely run caller defined method
			GoSafe(func() {
				defer wg.Done()
				fn(val, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	}()
	return Range(pipe)
}

func (s Stream) walkLimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan interface{}, option.workers)

	go func() {
		var wg sync.WaitGroup
		// 每次只允许option.workers个goroutine工作
		pool := make(chan PlaceholderType, option.workers)

		for item := range s.source {
			// important, used in another goroutine
			val := item
			pool <- Placeholder
			wg.Add(1)

			// better to safely run caller defined method
			GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				fn(val, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

func (s Stream) ForEach(fn ForEachFunc) {
	for item := range s.source {
		fn(item)
	}
}

// buildOptions returns a rxOptions with given customizations.
func buildOptions(opts ...Option) *rxOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

// newOptions returns a default rxOptions.
func newOptions() *rxOptions {
	return &rxOptions{
		workers: defaultWorkers,
	}
}

// UnlimitedWorkers lets the caller use as many workers as the tasks.
func UnlimitedWorkers() Option {
	return func(opts *rxOptions) {
		opts.unlimitedWorkers = true
	}
}

// WithWorkers lets the caller customize the concurrent workers.
func WithWorkers(workers int) Option {
	return func(opts *rxOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

package retry

const defaultRetryTimes = 3

type retryOptions struct {
	times int
}

type RetryOption func(options *retryOptions)

func DoWithRetry(fn func() error, opts ...RetryOption) error {
	options := newRetryOptions()
	for _, opt := range opts {
		opt(options)
	}

	var berr BatchError
	for i := 0; i < options.times; i++ {
		if err := fn(); err != nil {
			berr.Add(err)
		} else {
			return nil
		}
	}

	return berr.Err()
}

// WithRetry customize a DoWithRetry call with given retry times.
func WithRetry(times int) RetryOption {
	return func(options *retryOptions) {
		options.times = times
	}
}

func newRetryOptions() *retryOptions {
	return &retryOptions{
		times: defaultRetryTimes,
	}
}

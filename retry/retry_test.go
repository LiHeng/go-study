package retry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	assert.NotNil(t, DoWithRetry(func() error {
		return errors.New("any")
	}))

	assert.Nil(t, DoWithRetry(func() error {
		return nil
	}))

	assert.NotNil(t, DoWithRetry(func() error {
		return errors.New("any")
	}, WithRetry(2)))
}

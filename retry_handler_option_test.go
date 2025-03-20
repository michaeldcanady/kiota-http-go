package nethttplibrary

import (
	"errors"
	"fmt"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithMaxRetries_ItSetsMaxRetries(t *testing.T) {
	config := RetryHandlerOptions{}

	err := WithMaxRetries(5)(&config)

	assert.Nil(t, err)
	assert.Equal(t, 5, config.MaxRetries)
}

func TestWithMaxRetries_ItErrorsWhenLessThanDefault(t *testing.T) {
	config := RetryHandlerOptions{}

	err := WithMaxRetries(defaultMaxRedirects - 1)(&config)

	assert.Equal(t, fmt.Errorf("maxRetries must be greater than %v", defaultMaxRedirects), err)
	assert.Equal(t, 0, config.MaxRetries)
}

func TestWithMaxRetries_ItErrorsWhenOptionsIsNil(t *testing.T) {
	err := WithMaxRetries(5)(nil)

	assert.Equal(t, errors.New("r is nil"), err)
}

func TestWithRetryDelay_ItSetsDelay(t *testing.T) {
	config := RetryHandlerOptions{}

	err := WithRetryDelay(5)(&config)

	assert.Nil(t, err)
	assert.Equal(t, 5, config.DelaySeconds)
}

func TestWithRetryDelay_ItErrorsWhenLessThanDefault(t *testing.T) {
	config := RetryHandlerOptions{}

	err := WithRetryDelay(defaultDelaySeconds - 1)(&config)

	assert.Equal(t, fmt.Errorf("delayInSeconds must be greater than %v", defaultDelaySeconds), err)
	assert.Equal(t, 0, config.DelaySeconds)
}

func TestWithRetryDelay_ItErrorsWhenOptionsIsNil(t *testing.T) {
	err := WithRetryDelay(5)(nil)

	assert.Equal(t, errors.New("r is nil"), err)
}

func TestWithShouldRetry_ItSetsShouldRetry(t *testing.T) {
	config := RetryHandlerOptions{}

	err := WithShouldRetry(func(delay time.Duration, executionCount int, request *nethttp.Request, response *nethttp.Response) bool {
		return true
	})(&config)

	assert.Nil(t, err)
	assert.NotNil(t, config.ShouldRetry)
}

func TestWithShouldRetry_ItErrorsWhenShouldRetryIsNil(t *testing.T) {
	config := RetryHandlerOptions{}

	err := WithShouldRetry(nil)(&config)

	assert.Equal(t, errors.New("shouldRetry is nil"), err)
	assert.Nil(t, config.ShouldRetry)
}

func TestWithShouldRetry_ItErrorsWhenOptionsIsNil(t *testing.T) {
	err := WithShouldRetry(func(delay time.Duration, executionCount int, request *nethttp.Request, response *nethttp.Response) bool {
		return true
	})(nil)

	assert.Equal(t, errors.New("r is nil"), err)
}

package nethttplibrary

import (
	nethttp "net/http"
	"time"

	abs "github.com/microsoft/kiota-abstractions-go"
)

// RetryHandlerOptions to apply when evaluating the response for retrial
type RetryHandlerOptions struct {
	// Callback to determine if the request should be retried
	ShouldRetry func(delay time.Duration, executionCount int, request *nethttp.Request, response *nethttp.Response) bool
	// The maximum number of times a request can be retried
	MaxRetries int
	// The delay in seconds between retries
	DelaySeconds int
}

// GetKey returns the key value to be used when the option is added to the request context
func (options *RetryHandlerOptions) GetKey() abs.RequestOptionKey {
	return retryKeyValue
}

// GetShouldRetry returns the should retry callback function which evaluates the response for retrial
func (options *RetryHandlerOptions) GetShouldRetry() func(delay time.Duration, executionCount int, request *nethttp.Request, response *nethttp.Response) bool {
	return options.ShouldRetry
}

// GetDelaySeconds returns the delays in seconds between retries
func (options *RetryHandlerOptions) GetDelaySeconds() int {
	if options.DelaySeconds < 1 {
		return defaultDelaySeconds
	} else if options.DelaySeconds > absoluteMaxDelaySeconds {
		return absoluteMaxDelaySeconds
	} else {
		return options.DelaySeconds
	}
}

// GetMaxRetries returns the maximum number of times a request can be retried
func (options *RetryHandlerOptions) GetMaxRetries() int {
	if options.MaxRetries < 1 {
		return defaultMaxRetries
	} else if options.MaxRetries > absoluteMaxRetries {
		return absoluteMaxRetries
	} else {
		return options.MaxRetries
	}
}

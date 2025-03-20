package nethttplibrary

import (
	"errors"
	"fmt"
	nethttp "net/http"
	"time"

	"github.com/microsoft/kiota-http-go/internal"
)

// WithMaxRetries sets the maximum number of times a request can be retried.
func WithMaxRetries(maxRetries int) internal.Option[*RetryHandlerOptions] {
	return func(r *RetryHandlerOptions) error {
		if maxRetries < defaultMaxRedirects {
			return fmt.Errorf("maxRetries must be greater than %v", defaultMaxRedirects)
		}

		if r == nil {
			return errors.New("r is nil")
		}

		r.MaxRetries = maxRetries
		return nil
	}
}

// WithRetryDelay sets the delay in seconds between retries.
func WithRetryDelay(delayInSeconds int) internal.Option[*RetryHandlerOptions] {
	return func(r *RetryHandlerOptions) error {
		if delayInSeconds < defaultDelaySeconds {
			return fmt.Errorf("delayInSeconds must be greater than %v", defaultDelaySeconds)
		}

		if r == nil {
			return errors.New("r is nil")
		}

		r.DelaySeconds = delayInSeconds
		return nil
	}
}

// WithShouldRetry sets the callback to determine if the request should be retried.
func WithShouldRetry(shouldRetry func(delay time.Duration, executionCount int, request *nethttp.Request, response *nethttp.Response) bool) internal.Option[*RetryHandlerOptions] {
	return func(r *RetryHandlerOptions) error {
		if shouldRetry == nil {
			return errors.New("shouldRetry is nil")
		}

		if r == nil {
			return errors.New("r is nil")
		}

		r.ShouldRetry = shouldRetry
		return nil
	}
}

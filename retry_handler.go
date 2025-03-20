package nethttplibrary

import (
	"context"
	"fmt"
	"io"
	"math"
	nethttp "net/http"
	"strconv"
	"time"

	abs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-http-go/internal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// RetryHandler handles transient HTTP responses and retries the request given the retry options
type RetryHandler struct {
	// default options to use when evaluating the response
	options RetryHandlerOptions
}

// NewRetryHandler2 creates a new RetryHandler with the given options
func NewRetryHandler2(opts ...internal.Option[*RetryHandlerOptions]) (*RetryHandler, error) {
	options := RetryHandlerOptions{
		ShouldRetry: func(delay time.Duration, executionCount int, request *nethttp.Request, response *nethttp.Response) bool {
			return true
		},
		DelaySeconds: defaultDelaySeconds,
		MaxRetries:   defaultMaxRetries,
	}

	if err := internal.ApplyOptions(&options, opts...); err != nil {
		return nil, err
	}

	return NewRetryHandlerWithOptions(options), nil
}

// NewRetryHandler creates a new RetryHandler with default options
func NewRetryHandler() *RetryHandler {

	handler, _ := NewRetryHandler2()

	return handler
}

// NewRetryHandlerWithOptions creates a new RetryHandler with the given options
func NewRetryHandlerWithOptions(options RetryHandlerOptions) *RetryHandler {
	return &RetryHandler{options: options}
}

const defaultMaxRetries = 3
const absoluteMaxRetries = 10
const defaultDelaySeconds = 3
const absoluteMaxDelaySeconds = 180

type retryHandlerOptionsInt interface {
	abs.RequestOption
	GetShouldRetry() func(delay time.Duration, executionCount int, request *nethttp.Request, response *nethttp.Response) bool
	GetDelaySeconds() int
	GetMaxRetries() int
}

var retryKeyValue = abs.RequestOptionKey{
	Key: "RetryHandler",
}

const retryAttemptHeader = "Retry-Attempt"
const retryAfterHeader = "Retry-After"

const tooManyRequests = 429
const serviceUnavailable = 503
const gatewayTimeout = 504

// Intercept implements the interface and evaluates whether to retry a failed request.
func (middleware RetryHandler) Intercept(pipeline Pipeline, middlewareIndex int, req *nethttp.Request) (*nethttp.Response, error) {
	obsOptions := GetObservabilityOptionsFromRequest(req)
	ctx := req.Context()
	var span trace.Span
	var observabilityName string
	if obsOptions != nil {
		observabilityName = obsOptions.GetTracerInstrumentationName()
		ctx, span = otel.GetTracerProvider().Tracer(observabilityName).Start(ctx, "RetryHandler_Intercept")
		span.SetAttributes(attribute.Bool("com.microsoft.kiota.handler.retry.enable", true))
		defer span.End()
		req = req.WithContext(ctx)
	}
	response, err := pipeline.Next(req, middlewareIndex)
	if err != nil {
		return response, err
	}
	reqOption, ok := req.Context().Value(retryKeyValue).(retryHandlerOptionsInt)
	if !ok {
		reqOption = &middleware.options
	}
	return middleware.retryRequest(ctx, pipeline, middlewareIndex, reqOption, req, response, 0, 0, observabilityName)
}

// retryRequest retries the request if the response is retriable.
func (middleware RetryHandler) retryRequest(ctx context.Context, pipeline Pipeline, middlewareIndex int, options retryHandlerOptionsInt, req *nethttp.Request, resp *nethttp.Response, executionCount int, cumulativeDelay time.Duration, observabilityName string) (*nethttp.Response, error) {
	if middleware.isRetriableErrorCode(resp.StatusCode) &&
		middleware.isRetriableRequest(req) &&
		executionCount < options.GetMaxRetries() &&
		cumulativeDelay < time.Duration(absoluteMaxDelaySeconds)*time.Second &&
		options.GetShouldRetry()(cumulativeDelay, executionCount, req, resp) {
		executionCount++
		delay := middleware.getRetryDelay(req, resp, options, executionCount)
		cumulativeDelay += delay
		req.Header.Set(retryAttemptHeader, strconv.Itoa(executionCount))
		if req.Body != nil {
			s, ok := req.Body.(io.Seeker)
			if ok {
				s.Seek(0, io.SeekStart)
			}
		}
		if observabilityName != "" {
			ctx, span := otel.GetTracerProvider().Tracer(observabilityName).Start(ctx, "RetryHandler_Intercept - attempt "+fmt.Sprint(executionCount))
			span.SetAttributes(attribute.Int("http.request.resend_count", executionCount),

				httpResponseStatusCodeAttribute.Int(resp.StatusCode),
				attribute.Float64("http.request.resend_delay", delay.Seconds()),
			)
			defer span.End()
			req = req.WithContext(ctx)
		}
		t := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			// Return without retrying if the context was cancelled.
			return nil, ctx.Err()

			// Leaving this case empty causes it to exit the switch-block.
		case <-t.C:
		}
		response, err := pipeline.Next(req, middlewareIndex)
		if err != nil {
			return response, err
		}
		return middleware.retryRequest(ctx, pipeline, middlewareIndex, options, req, response, executionCount, cumulativeDelay, observabilityName)
	}
	return resp, nil
}

// isRetriableErrorCode determines if the response status code should be retried.
func (middleware RetryHandler) isRetriableErrorCode(code int) bool {
	return code == tooManyRequests || code == serviceUnavailable || code == gatewayTimeout
}

// isRetriableRequest determines if the request should be retried.
func (middleware RetryHandler) isRetriableRequest(req *nethttp.Request) bool {
	isBodiedMethod := req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH"
	if isBodiedMethod && req.Body != nil {
		return req.ContentLength != -1
	}
	return true
}

// getRetryDelay calculates the delay between retries.
func (middleware RetryHandler) getRetryDelay(req *nethttp.Request, resp *nethttp.Response, options retryHandlerOptionsInt, executionCount int) time.Duration {
	retryAfter := resp.Header.Get(retryAfterHeader)
	if retryAfter != "" {
		retryAfterDelay, err := strconv.ParseFloat(retryAfter, 64)
		if err == nil {
			return time.Duration(retryAfterDelay) * time.Second
		}

		// parse the header if it's a date
		t, err := time.Parse(time.RFC1123, retryAfter)
		if err == nil {
			return t.Sub(time.Now())
		}
	}

	return time.Duration(math.Pow(float64(options.GetDelaySeconds()), float64(executionCount))) * time.Second
}

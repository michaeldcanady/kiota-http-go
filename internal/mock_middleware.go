package internal

import (
	nethttp "net/http"

	"github.com/stretchr/testify/mock"
)

// Pipeline contract for middleware infrastructure
type Pipeline interface {
	// Next moves the request object through middlewares in the pipeline
	Next(req *nethttp.Request, middlewareIndex int) (*nethttp.Response, error)
}

type MockMiddleware struct {
	mock.Mock
}

func NewMockMiddleware() *MockMiddleware {
	return &MockMiddleware{
		Mock: mock.Mock{},
	}
}

func (mM *MockMiddleware) Intercept(pipeline Pipeline, retry int, request *nethttp.Request) (*nethttp.Response, error) {
	args := mM.Called(pipeline, retry, request)
	return args.Get(0).(*nethttp.Response), args.Error(1)
}

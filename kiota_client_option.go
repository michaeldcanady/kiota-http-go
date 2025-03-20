package nethttplibrary

import (
	"errors"
	nethttp "net/http"
	"net/url"

	"github.com/microsoft/kiota-http-go/internal"
)

func matchAny[T any](input T, matchers ...func(T) bool) bool {
	for _, matcher := range matchers {
		if matcher(input) {
			return true
		}
	}
	return false
}

// WithTransportMiddleware sets the middleware for the HTTP client.
func WithTransportMiddleware(middleware ...Middleware) internal.Option[*nethttp.Client] {
	return func(c *nethttp.Client) error {
		if c == nil {
			return errors.New("c is nil")
		}

		// possibility that all or some middleware is nil, currently errors if any middleware is nil
		if len(middleware) == 0 || matchAny(middleware, func(middlewares []Middleware) bool {
			for _, middleware := range middlewares {
				return middleware == nil
			}
			return false
		}) {
			return errors.New("middleware is empty")
		}

		c.Transport = NewCustomTransport(middleware...)
		return nil
	}
}

// WithTransportProxyAuthenticated sets the proxy URL, username, and password for the HTTP client.
func WithTransportProxyAuthenticated(proxyUrlStr string, username string, password string, middleware ...Middleware) internal.Option[*nethttp.Client] {
	return func(c *nethttp.Client) error {
		if c == nil {
			return errors.New("c is nil")
		}

		user := url.UserPassword(username, password)
		transport, err := getTransportWithProxy(proxyUrlStr, user, middleware...)
		if err != nil {
			return err
		}
		c.Transport = transport
		return nil
	}
}

// WithTransportProxy sets the proxy URL for the HTTP client.
func WithTransportProxyUnauthenticated(proxyUrlStr string, middleware ...Middleware) internal.Option[*nethttp.Client] {
	return func(c *nethttp.Client) error {
		if c == nil {
			return errors.New("c is nil")
		}

		transport, err := getTransportWithProxy(proxyUrlStr, nil, middleware...)
		if err != nil {
			return err
		}
		c.Transport = transport
		return nil
	}
}

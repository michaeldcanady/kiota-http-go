package nethttplibrary

import (
	"errors"
	nethttp "net/http"
	"testing"

	"github.com/microsoft/kiota-http-go/internal"
	"github.com/stretchr/testify/assert"
)

func TestWithTransportMiddleware_ItSetsClientTransport(t *testing.T) {
	middleware := internal.NewMockMiddleware()

	client := nethttp.Client{}

	err := WithTransportMiddleware(middleware)(&client)

	assert.Nil(t, err)
	assert.NotNil(t, client.Transport)
}

func TestWithTransportMiddleware_ItErrorsIfNoMiddlewareIsProvided(t *testing.T) {
	client := nethttp.Client{}

	err := WithTransportMiddleware(nil)(&client)

	assert.Equal(t, errors.New("middleware is empty"), err)
	assert.Nil(t, client.Transport)
}

func TestWithTransportMiddleware_ItErrorsIfClientIsNil(t *testing.T) {
	middleware := internal.NewMockMiddleware()

	err := WithTransportMiddleware(middleware)(nil)

	assert.Equal(t, errors.New("c is nil"), err)
}

func TestWithTransportProxyAuthenticated_ItSetsClientTransport(t *testing.T) {
	middleware := internal.NewMockMiddleware()

	client := nethttp.Client{}

	err := WithTransportProxyAuthenticated("http://localhost", "username", "password", middleware)(&client)

	assert.Nil(t, err)
	assert.NotNil(t, client.Transport)
}

func TestWithTransportProxyAuthenticated_ItErrorsIfClientIsNil(t *testing.T) {
	middleware := internal.NewMockMiddleware()

	err := WithTransportProxyAuthenticated("http://localhost", "username", "password", middleware)(nil)

	assert.Equal(t, errors.New("c is nil"), err)
}

func TestWithTransportProxyUnauthenticated_ItSetsClientTransport(t *testing.T) {
	middleware := internal.NewMockMiddleware()

	client := nethttp.Client{}

	err := WithTransportProxyUnauthenticated("http://localhost", middleware)(&client)

	assert.Nil(t, err)
	assert.NotNil(t, client.Transport)
}

func TestWithTransportProxyUnauthenticated_ItErrorsIfClientIsNil(t *testing.T) {
	middleware := internal.NewMockMiddleware()

	err := WithTransportProxyUnauthenticated("http://localhost", middleware)(nil)

	assert.Equal(t, errors.New("c is nil"), err)
}

package nethttplibrary

import (
	"errors"
	"net/http"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-http-go/internal"
)

func WithParseNodeFactory(parseNodeFactory absser.ParseNodeFactory) internal.Option[*NetHttpRequestAdapter] {
	return func(adapter *NetHttpRequestAdapter) error {
		if parseNodeFactory == nil {
			return errors.New("parseNodeFactory is nil")
		}

		if adapter == nil {
			return errors.New("adapter is nil")
		}

		adapter.parseNodeFactory = parseNodeFactory

		return nil
	}
}

func WithSerializationWriterFactory(serializationWriterFactory absser.SerializationWriterFactory) internal.Option[*NetHttpRequestAdapter] {
	return func(adapter *NetHttpRequestAdapter) error {
		if serializationWriterFactory == nil {
			return errors.New("serializationWriterFactory is nil")
		}

		if adapter == nil {
			return errors.New("adapter is nil")
		}

		adapter.serializationWriterFactory = serializationWriterFactory

		return nil
	}
}

func WithHttpClient(client *http.Client) internal.Option[*NetHttpRequestAdapter] {
	return func(adapter *NetHttpRequestAdapter) error {
		if client == nil {
			return errors.New("client is nil")
		}

		if adapter == nil {
			return errors.New("adapter is nil")
		}

		adapter.httpClient = client

		return nil
	}
}

// WithObservabilityOptions sets the observability options for the adapter
func WithObservabilityOptions(options ObservabilityOptions) internal.Option[*NetHttpRequestAdapter] {
	// if there is desire to be more specific with this option it can be made into WithObservabilityOptionsXXX or it can accept Options (I'd advise against the latter)
	return func(adapter *NetHttpRequestAdapter) error {
		if adapter == nil {
			return errors.New("adapter is nil")
		}

		adapter.observabilityOptions = options

		return nil
	}
}

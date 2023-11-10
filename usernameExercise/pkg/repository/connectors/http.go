package connectors

import (
	"context"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	GetClient(ctx context.Context) (http.Client, error)
	GetEndpoint(ctx context.Context) *url.URL
}

func NewHTTPClient(ctx context.Context, endpoint string) (HTTPClient, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return httpClient{}, err
	}
	return httpClient{client: http.Client{}, endpoint: endpointURL}, nil
}

type httpClient struct {
	client   http.Client
	endpoint *url.URL
}

func (h httpClient) GetClient(ctx context.Context) (http.Client, error) {
	return h.client, nil
}

func (h httpClient) GetEndpoint(ctx context.Context) *url.URL {
	return h.endpoint
}

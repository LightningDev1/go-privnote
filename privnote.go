package privnote

import (
	http "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client/profiles"

	tls_client "github.com/bogdanfinn/tls-client"
)

// Headers are the default headers used by the Privnote client.
var Headers = http.Header{
	"accept":           {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"},
	"accept-encoding":  {""},
	"accept-language":  {"en-US,en;q=0.5"},
	"referer":          {"https://privnote.com/"},
	"origin":           {"https://privnote.com"},
	"sec-fetch-dest":   {"empty"},
	"sec-fetch-mode":   {"cors"},
	"sec-fetch-site":   {"same-origin"},
	"user-agent":       {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"},
	"x-requested-with": {"XMLHttpRequest"},
}

// TLSClientOptions are the default TLS client options used by the Privnote client.
var TLSClientOptions = []tls_client.HttpClientOption{
	tls_client.WithTimeoutSeconds(20),
	tls_client.WithClientProfile(profiles.Chrome_117),
	tls_client.WithCookieJar(tls_client.NewCookieJar()),
	tls_client.WithRandomTLSExtensionOrder(),
}

// NewClient creates a new Privnote client for creating and opening notes.
func NewClient() *Client {
	// A TLS client is required to make requests to the Privnote API.
	// This is because Cloudflare blocks requests without a valid TLS fingerprint.
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), TLSClientOptions...)

	return &Client{
		BaseURL:    "https://privnote.com",
		HTTPClient: client,
	}
}

package rbac

import (
	"bytes"
	"crypto/tls"
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

// Client for the RBAC API
type Client struct {
	resty  *resty.Client
	strict bool
}

// NewClient access the RBAC API via TLS
func NewClient(hostURL string, tlsConfig *tls.Config) *Client {
	r := resty.New()
	if tlsConfig != nil {
		r.SetTLSClientConfig(tlsConfig)
	}
	r.SetHostURL(hostURL)
	r.SetError(APIError{})
	client := Client{resty: r}
	r.JSONUnmarshal = func(data []byte, v interface{}) error {
		d := json.NewDecoder(bytes.NewReader(data))
		if client.strict {
			d.DisallowUnknownFields()
		}
		return d.Decode(v)
	}
	return &client
}

// APIError represents an error response from the RBAC API
type APIError struct {
	Kind string `json:"kind"`
	Msg  string `json:"msg"`
}

func (oe *APIError) Error() string {
	return oe.Msg
}

package sebel

import (
	"net/http"
)

// roundTripper is a wrapper for the [http.RoundTripper] interface with a
// reference to Sebel.
type roundTripper struct {
	http.RoundTripper
	sebel *Sebel
}

// RoundTrip implements the [http.RoundTripper] interface and checks the TLS
// connection after the round trip.
func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := rt.RoundTripper.RoundTrip(req)
	if resp.TLS != nil && err == nil {
		rt.sebel.tls = resp.TLS

		_, err = rt.sebel.checkTLS()
		if err != nil {
			resp = nil
		}
	}

	return resp, err
}

package sebel

import (
	"testing"

	"net/http"
)

func TestRoundTrip(t *testing.T) {
	t.Parallel()

	tr := http.DefaultTransport
	rt := &roundTripper{RoundTripper: tr, sebel: sebel}

	req, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
	if err != nil {
		t.SkipNow()
	}

	_, err = rt.RoundTrip(req)
	if err != nil {
		t.Fail()
	}
}

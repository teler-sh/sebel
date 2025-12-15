package sebel_test

import (
	"net/http"
	"time"

	"github.com/teler-sh/sebel"
)

func ExampleNew() {
	client := &http.Client{
		Transport: sebel.New().RoundTripper(http.DefaultTransport),
	}

	resp, err := client.Get("https://c2.host")
	if err != nil && sebel.IsBlacklist(err) {
		// certificate blacklisted
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	println("OK")
}

// ExampleNew_autoRefresh demonstrates creating a Sebel instance with automatic
// background SSLBL data refresh.
func ExampleNew_autoRefresh() {
	s := sebel.New(sebel.Options{
		DataRefreshInterval: 5 * time.Minute,
	})
	defer s.Close() // Stop background refresh

	client := &http.Client{
		Transport: s.RoundTripper(http.DefaultTransport),
	}

	resp, err := client.Get("https://c2.host")
	if err != nil && sebel.IsBlacklist(err) {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
}

func ExampleSebel_CheckTLS() {
	r, err := http.Get("https://c2.host")
	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()

	s := sebel.New()

	_, err = s.CheckTLS(r.TLS)
	if err != nil && sebel.IsBlacklist(err) {
		// certificate blacklisted
		panic(err)
	}
}

// To seamlessly integrate it without need to configure a new client, you can
// simply replace your current [http.DefaultClient] with sebel's RoundTripper.
func ExampleSebel_RoundTripper() {
	http.DefaultClient.Transport = sebel.New().RoundTripper(http.DefaultTransport)
}

func ExampleSebel_CheckHost() {
	s := sebel.New()

	_, err := s.CheckHost("c2.host", "443", nil)
	if err != nil && sebel.IsBlacklist(err) {
		// certificate blacklisted
		panic(err)
	}
}

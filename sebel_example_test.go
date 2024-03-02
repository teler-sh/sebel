package sebel_test

import (
	"net/http"

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
	defer resp.Body.Close()

	println("OK")
}

func ExampleSebel_CheckTLS() {
	r, err := http.Get("https://c2.host")
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	s := sebel.New()

	_, err = s.CheckTLS(r.TLS)
	if err != nil && sebel.IsBlacklist(err) {
		// certificate blacklisted
		panic(err)
	}
}

// To seamlessly integrate it without need to configure a new client, you can
// simply replace your current http.DefaultClient with sebel's RoundTripper.
func ExampleSebel_RoundTripper() {
	http.DefaultClient.Transport = sebel.New().RoundTripper(http.DefaultTransport)
}

package sebel

import (
	"testing"

	"net/http"
)

var sebel *Sebel

func init() {
	sebel = New()
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		if sebel == nil {
			t.Fail()
		}
	})

	t.Run("set", func(t *testing.T) {
		sebelWithOptions := New(Options{DisableSSLBlacklist: true})
		if !sebelWithOptions.options.DisableSSLBlacklist {
			t.Fail()
		}
	})
}

func TestRoundTripper(t *testing.T) {
	t.Parallel()

	tr := http.DefaultTransport
	if tr == sebel.RoundTripper(tr) {
		t.Fail()
	}
}

func TestCheckTLS(t *testing.T) {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		t.SkipNow()
	}

	if resp.TLS == nil {
		t.SkipNow()
	}

	t.Parallel()

	t.Run("default", func(t *testing.T) {
		_, err := sebel.CheckTLS(resp.TLS)
		if err != nil {
			t.Fail()
		}
	})

	t.Run("WithDisableSSLBlacklist", func(t *testing.T) {
		sebelWithOptions := New(Options{DisableSSLBlacklist: true})

		_, err := sebelWithOptions.CheckTLS(resp.TLS)
		if err != nil {
			t.Fail()
		}
	})
}

func BenchmarkCheckTLS(b *testing.B) {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		b.SkipNow()
	}

	if resp.TLS == nil {
		b.SkipNow()
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = sebel.CheckTLS(resp.TLS)
	}
}

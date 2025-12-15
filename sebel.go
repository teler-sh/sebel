package sebel

import (
	"fmt"
	"net"

	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/teler-sh/sebel/pkg/cert"
	"github.com/teler-sh/sebel/pkg/sslbl"
)

// Sebel holds information and [Options].
type Sebel struct {
	data    data
	options *Options
	tls     *tls.ConnectionState
}

// New creates a new instance of [Sebel] with the provided options.
func New(opt ...Options) *Sebel {
	sebel := new(Sebel)

	if len(opt) > 0 {
		sebel.options = &opt[0]
	} else {
		sebel.options = new(Options)
	}

	sebel.data.sslbl = sslbl.MustGet()

	if sebel.options.DataRefreshInterval > 0 {
		sslbl.StartBackgroundRefresh(sebel.options.DataRefreshInterval)
	}

	return sebel
}

// Close stops the background refresh goroutine if running.
//
// It is only necessary to call this if [Options].DataRefreshInterval was set.
// Safe to call even if background refresh was not started.
func (s *Sebel) Close() {
	sslbl.StopBackgroundRefresh()
}

// RoundTripper creates a new RoundTripper using the provided [http.RoundTripper]
// and [Sebel] instance.
func (s *Sebel) RoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &roundTripper{RoundTripper: rt, sebel: s}
}

// CheckTLS checks the TLS connection against the SSLBL (SSL Blacklist) and
// returns the SSLBL record.
//
// It returns [ErrSSLBlacklist] error if the certificate is blacklisted.
func (s *Sebel) CheckTLS(connState *tls.ConnectionState) (*sslbl.Record, error) {
	s.tls = connState

	return s.checkTLS()
}

// CheckHost connects to the specified host, retrieves its TLS certificate,
// and checks it against the SSLBL (SSL Blacklist).
//
// The config parameter allows customizing TLS behavior. Pass nil to use defaults.
//
// It returns [ErrSSLBlacklist] error if the certificate is blacklisted.
func (s *Sebel) CheckHost(host, port string, config *tls.Config) (*sslbl.Record, error) {
	conn, err := tls.Dial("tcp", net.JoinHostPort(host, port), config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	connState := conn.ConnectionState()
	s.tls = &connState

	return s.checkTLS()
}

// getCert retrieves the peer certificate from the TLS connection state.
func (s *Sebel) getCert() *x509.Certificate {
	if s.tls == nil {
		return nil
	}

	if cert := s.tls.PeerCertificates; len(cert) > 0 {
		return cert[0]
	}

	return nil
}

// checkTLS runs actual checks on the TLS connection and returns the SSLBL
// record and [ErrSSLBlacklist] error if blacklisted.
func (s *Sebel) checkTLS() (*sslbl.Record, error) {
	record, ok := new(sslbl.Record), false

	// return early if disabled
	if s.options.DisableSSLBlacklist {
		return record, nil
	}

	data := s.data.sslbl
	if len(data) == 0 {
		return record, ErrNoSSLBLData
	}

	certificate := s.getCert()
	if certificate == nil {
		return record, nil
	}

	fingerprint := cert.New(certificate)
	sha1sum := fingerprint.SHA1().String()

	record, ok = sslbl.Find(sha1sum, data)
	if ok {
		reason := record.Listing.Reason
		s.write(record, sha1sum)

		return record, fmt.Errorf("%w: %s detected", ErrSSLBlacklist, reason)
	}

	return record, nil
}

// write writes the blacklist detection to the configured output writer.
func (s *Sebel) write(record *sslbl.Record, fingerprint string) {
	if s.options.Output == nil {
		return
	}

	var msg string
	if s.options.Formatter != nil {
		msg = s.options.Formatter(record, fingerprint)
	} else {
		msg = fmt.Sprintf("%s: fingerprint=%q reason=%q\n", ErrSSLBlacklist, fingerprint, record.Listing.Reason)
	}

	_, _ = fmt.Fprint(s.options.Output, msg)
}

// Package cert provides utilities for working with SSL/TLS certificates,
// including fingerprint generation.
package cert

import (
	"bytes"
	"fmt"

	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
)

// Fingerprint defines methods to calculate SHA1 and SHA256 fingerprints.
//
// If the certificate is nil, an empty buffer is returned from SHA1 & SHA256.
type Fingerprint interface {
	SHA1() *bytes.Buffer
	SHA256() *bytes.Buffer
}

// fingerprint implementation of Fingerprint interface.
type fingerprint struct {
	*x509.Certificate
}

// New creates a new Fingerprint instance based on the [x509.Certificate].
func New(cert *x509.Certificate) Fingerprint {
	return &fingerprint{Certificate: cert}
}

func (f *fingerprint) SHA1() *bytes.Buffer {
	var b bytes.Buffer

	if f.Certificate == nil {
		return &b
	}

	sha1sum := sha1.Sum(f.Certificate.Raw)
	for _, v := range sha1sum {
		fmt.Fprintf(&b, "%02x", v)
	}

	return &b
}

func (f *fingerprint) SHA256() *bytes.Buffer {
	var b bytes.Buffer

	if f.Certificate == nil {
		return &b
	}

	sha256sum := sha256.Sum256(f.Certificate.Raw)
	for _, v := range sha256sum {
		fmt.Fprintf(&b, "%02x", v)
	}

	return &b
}

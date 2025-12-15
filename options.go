package sebel

import (
	"io"

	"github.com/teler-sh/sebel/pkg/sslbl"
)

// Options holds configuration settings for the [Sebel] package.
type Options struct {
	// DisableSSLBlacklist, when set to true, disables SSL/TLS certificate
	// blacklist checks.
	DisableSSLBlacklist bool

	// Output specifies the writer for logging blacklist detections.
	// If nil, no output is written.
	Output io.Writer

	// Formatter customizes the output format for blacklist detections.
	// If nil, a default format is used.
	Formatter func(record *sslbl.Record, fingerprint string) string

	// TODO(dwisiswant0): Add these fields
	// DisableHostBlacklist bool
}

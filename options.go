package sebel

import (
	"io"
	"time"

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

	// DataRefreshInterval specifies how often to refresh SSLBL data in the
	// background. SSLBL updates every 5 minutes, so values less than 5 minutes
	// are not recommended.
	//
	// If zero or negative, background refresh is disabled.
	DataRefreshInterval time.Duration

	// TODO(dwisiswant0): Add these fields
	// DisableHostBlacklist bool
}

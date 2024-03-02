package sebel

// Options holds configuration settings for the [Sebel] package.
type Options struct {
	// DisableSSLBlacklist, when set to true, disables SSL/TLS certificate
	// blacklist checks.
	DisableSSLBlacklist bool

	// TODO(dwisiswant0): Add these fields
	// DisableHostBlacklist bool
	// Output               io.Writer
}

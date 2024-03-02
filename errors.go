package sebel

import "errors"

var (
	ErrSSLBlacklist = errors.New("certificate blacklisted")
	ErrNoSSLBLData  = errors.New("no SSLBL data")
)

// IsBlacklist checks if the given error is an [ErrSSLBlacklist].
func IsBlacklist(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrSSLBlacklist)
}

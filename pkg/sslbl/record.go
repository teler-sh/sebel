package sslbl

// Record represents an entry in the SSL Blacklist (SSLBL).
type Record struct {
	Listing struct {
		// Date is the date when the SSL certificate was listed in the SSLBL.
		Date string
		// Reason provides information about why the SSL certificate was
		// blacklisted.
		Reason string
	}
	// SHA1Sum is the SHA-1 fingerprint of the SSL certificate associated with
	// the record.
	SHA1Sum string
}

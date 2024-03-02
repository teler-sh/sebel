package sslbl

import (
	"io"

	"encoding/csv"
	"net/http"

	"github.com/samber/lo"
)

// Find searches for a record with a given SHA-1 fingerprint in a slice of SSLBL
// records.
//
// It returns the found record and a boolean indicating whether the record was
// found.
func Find(sha1sum string, records []Record) (*Record, bool) {
	record, ok := lo.Find(records, func(r Record) bool {
		return r.SHA1Sum == sha1sum
	})

	return &record, ok
}

// Get retrieves SSLBL records from a `sslbl.abuse.ch`, parses the CSV data,
// and returns them.
func Get() ([]Record, error) {
	var records []Record

	resp, err := http.Get(dataURL)
	if err != nil {
		return records, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return records, err
	}

	csvData := sanitizeBody(body)
	csvReader := csv.NewReader(csvData)

	data, err := csvReader.ReadAll()
	if err != nil {
		return records, err
	}
	records = parseCSV(data)

	return records, nil
}

// MustGet is like [Get] but panics if there is an error during the retrieval
// process.
func MustGet() []Record {
	records, err := Get()
	if err != nil {
		panic(err)
	}

	return records
}

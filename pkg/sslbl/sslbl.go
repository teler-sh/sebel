package sslbl

import (
	"encoding/csv"
	"io"
	"net/http"
	"time"

	"github.com/maypok86/otter/v2"
	"github.com/samber/lo"
	"github.com/teler-sh/sebel/internal/cache"
)

// recordCache is the package-level cache for SSLBL records.
var recordCache = cache.Must(&cache.Options[string, Record]{
	Options: otter.Options[string, Record]{
		InitialCapacity: 10_000,
	},
	Filepath: cache.BuildCacheFilepath("sebel", "records.gob"),
})

// Find searches for a record with a given SHA-1 fingerprint.
//
// It first checks the cache, then falls back to the provided records slice if
// not found. Returns the found record and a boolean indicating whether the
// record was found.
func Find(sha1sum string, records []Record) (*Record, bool) {
	if record, ok := recordCache.GetIfPresent(sha1sum); ok {
		return &record, true
	}

	record, ok := lo.Find(records, func(r Record) bool {
		return r.SHA1Sum == sha1sum
	})

	return &record, ok
}

// Get retrieves SSLBL records from cache if available, otherwise fetches from
// `sslbl.abuse.ch`, parses the CSV data, caches the results, and returns them.
func Get() ([]Record, error) {
	if recordCache.EstimatedSize() > 0 {
		var records []Record

		for _, record := range recordCache.All() {
			records = append(records, record)
		}

		return records, nil
	}

	records, err := fetch()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		recordCache.Set(record.SHA1Sum, record)
	}
	_ = recordCache.Save()

	return records, nil
}

// Fetch retrieves fresh SSLBL records from `sslbl.abuse.ch`, bypassing the
// cache. It updates the cache with the fetched records.
func Fetch() ([]Record, error) {
	records, err := fetch()
	if err != nil {
		return nil, err
	}

	recordCache.InvalidateAll()
	for _, record := range records {
		recordCache.Set(record.SHA1Sum, record)
	}
	_ = recordCache.Save()

	return records, nil
}

// fetch retrieves SSLBL records from the remote URL.
func fetch() ([]Record, error) {
	resp, err := http.Get(dataURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	csvData := sanitizeBody(body)
	csvReader := csv.NewReader(csvData)

	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return parseCSV(data), nil
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

// StartBackgroundRefresh starts a goroutine that periodically fetches fresh
// SSLBL data at the specified interval. This keeps the cache up-to-date.
// SSLBL updates every 5 minutes, so intervals less than 5 minutes are not
// recommended.
//
// Call StopBackgroundRefresh to stop the background refresh.
func StartBackgroundRefresh(interval time.Duration) {
	recordCache.StartBackgroundRefresh(interval, func() error {
		_, err := Fetch()
		return err
	})
}

// StopBackgroundRefresh stops the background refresh goroutine.
func StopBackgroundRefresh() {
	recordCache.StopBackgroundRefresh()
}

package sslbl

import "bytes"

func sanitizeBody(body []byte) *bytes.Reader {
	lines := bytes.Split(body, []byte("\n"))
	var split [][]byte

	for _, line := range lines {
		if !bytes.HasPrefix(line, []byte("#")) {
			split = append(split, line)
		}
	}

	cleaned := bytes.Join(split, []byte("\n"))

	// pattern := regexp.MustCompile(`^\s*#.*$`)
	// cleaned := pattern.ReplaceAll(body, []byte(""))

	return bytes.NewReader(cleaned)
}

func parseCSV(records [][]string) []Record {
	var sslBlRecords []Record

	for _, record := range records {
		var data Record

		for i, value := range record {
			switch i {
			case 0:
				data.Listing.Date = value
			case 1:
				data.SHA1Sum = value
			case 2:
				data.Listing.Reason = value
			}
		}

		sslBlRecords = append(sslBlRecords, data)
	}

	return sslBlRecords
}

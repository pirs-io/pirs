package parsers

import (
	"encoding/csv"
	"os"
)

func ReadCsvFile(filePath string, lazyQuotes bool) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		// todo log
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	csvReader.LazyQuotes = lazyQuotes
	records, err := csvReader.ReadAll()
	if err != nil {
		// todo log
	}
	return records
}

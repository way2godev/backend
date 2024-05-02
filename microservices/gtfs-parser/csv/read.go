package csv

import (
	"encoding/csv"
	"os"
)

// Read reads the contents of a CSV file and returns a slice of maps, where each map represents a row in the CSV file.
// The keys of the map are the column names from the CSV file, and the values are the corresponding field values for each row.
// The first row of the CSV file is assumed to be the header row containing the column names.
func Read(filename string) ([]map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var data []map[string]string
	header := records[0]
	for _, record := range records[1:] {
		row := make(map[string]string)
		for i, field := range record {
			row[header[i]] = field
		}
		data = append(data, row)
	}
	return data, nil
}

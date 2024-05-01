package csv

import (
	"encoding/csv"
	"os"
)

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

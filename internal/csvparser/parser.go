package csvparser

import (
	"encoding/csv"
	"os"
)

type Comment struct {
	Text   string `csv:"comment"`
	// Add other fields from your CSV as needed
}

func ParseCSV(filePath string) ([]Comment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // or '\t' for TSV
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var comments []Comment
	count := 0
	for _, record := range records[1:] { // Skip header
		if len(record) > 0 {
			comments = append(comments, Comment{Text: record[0]})
		}
		count++
		if count >= 10 {
			break
		}
	}

	return comments, nil
}
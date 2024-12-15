package repositories

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"

	"dictionary-scraper/entities"
)

type CSVRepository struct{}

func NewCSVRepository() *CSVRepository {
	return &CSVRepository{}
}

func (r *CSVRepository) ReadWordsFromFile(filePath string, language entities.Language) ([]entities.Word, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var words []entities.Word
	for _, record := range records {
		if len(record) < 2 {
			return nil, errors.New("invalid CSV format")
		}
		words = append(words, entities.Word{
			Language:     language,
			Term:         strings.TrimSpace(record[0]),
			PartOfSpeech: strings.TrimSpace(record[1]),
		})
	}

	return words, nil
}

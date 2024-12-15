package repositories

import "dictionary-scraper/entities"

// TranslationFetcher defines an interface for fetching translations from a source
type TranslationFetcher interface {
	FetchTranslations(word entities.Word, toLang entities.Language) ([]string, string, error)
}

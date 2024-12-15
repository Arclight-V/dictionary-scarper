package usecases

import (
	"dictionary-scraper/entities"
	"dictionary-scraper/repositories"
	"errors"
	log "github.com/sirupsen/logrus"
)

// TranslationService manages fetching translations from various sources
type TranslationService struct {
	fetchers map[entities.Language]repositories.TranslationFetcher
}

// NewTranslationService creates a new TranslationService
func NewTranslationService(fetchers map[entities.Language]repositories.TranslationFetcher) *TranslationService {
	return &TranslationService{fetchers: fetchers}
}

// GetTranslations fetches translations from all available sources
func (s *TranslationService) GetTranslations(word entities.Word, toLang entities.Language) (entities.Translations, string, error) {
	if word.Term == "" {
		return nil, "", errors.New("term cannot be empty")
	}

	allTranslations := make(entities.Translations)

	translations, transcription, err := s.fetchers[word.Language].FetchTranslations(word, toLang)
	allTranslations[toLang] = append(allTranslations[toLang], translations...)

	if len(allTranslations) == 0 {
		return nil, "", errors.New("no translations found")
	}
	log.Info(translations, transcription)
	return allTranslations, transcription, err
}

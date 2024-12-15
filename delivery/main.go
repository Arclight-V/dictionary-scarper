package main

import (
	"dictionary-scraper/entities"
	"dictionary-scraper/repositories"
	"dictionary-scraper/usecases"
	"flag"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.New()

	// Flags for CSV file paths
	frenchFilePath := flag.String("french", "", "Path to the French CSV file")
	englishFilePath := flag.String("english", "", "Path to the English CSV file")
	russianFilePath := flag.String("russian", "", "Path to the Russian CSV file")
	flag.Parse()

	// Checking for mandatory flags
	if *frenchFilePath == "" || *englishFilePath == "" || *russianFilePath == "" {
		logger.Error("Error: all file paths must be provided")
		flag.Usage()
		return
	}
	// Repositories
	csvRepo := repositories.NewCSVRepository()

	// Read data from CSV files
	frenchWords, err := csvRepo.ReadWordsFromFile(*frenchFilePath, entities.French)
	if err != nil {
		logger.Fatal("Error reading French words:", err)
		return
	}

	englishWords, err := csvRepo.ReadWordsFromFile(*englishFilePath, entities.English)
	if err != nil {
		logger.Fatal("Error reading English words:", err)
		return
	}

	russianWords, err := csvRepo.ReadWordsFromFile(*russianFilePath, entities.Russian)
	if err != nil {
		logger.Fatal("Error reading Russian words:", err)
		return
	}

	// UseCases
	wordTranslator := usecases.NewWordTranslator()

	// Aggregate all words by language
	wordsByLanguage := map[entities.Language][]entities.Word{
		entities.French:  frenchWords,
		entities.English: englishWords,
		entities.Russian: russianWords,
	}

	// Translate words between languages
	translatedWords := wordTranslator.TranslateWords(wordsByLanguage)

	// Initialize parsers
	dictionaryCambridgeParser := repositories.NewDictionaryCambridgeParser()

	// Register parsers in the service
	translationService := usecases.NewTranslationService(map[entities.Language]repositories.TranslationFetcher{
		entities.French: dictionaryCambridgeParser,
	})

	// Display the translated words
	for _, word := range translatedWords {
		logger.Infof("[%s] %s (%s)\n", word.Language, word.Term, word.PartOfSpeech)
		if word.Language == entities.French {
			word.Translations, word.Transcription, err = translationService.GetTranslations(word, entities.English)
			if err != nil {
				log.Errorf("Error GetTranslations", err)
				continue
			}
		}

	}
}

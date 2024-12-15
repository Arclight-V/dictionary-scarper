package usecases

import (
	"dictionary-scraper/common"
	"dictionary-scraper/entities"
)

type WordTranslator struct {
	posParser *common.PartOfSpeechParser
}

func NewWordTranslator() *WordTranslator {
	return &WordTranslator{
		posParser: common.GetPartOfSpeechParserInstance(),
	}
}

// TranslateWords links words between languages and adds translations to the Word structure
func (t *WordTranslator) TranslateWords(wordsByLanguage map[entities.Language][]entities.Word) []entities.Word {
	translations := make([]entities.Word, 0)

	// Indexing words by term and language for quick search
	index := make(map[string]map[entities.Language]*entities.Word)
	for lang, words := range wordsByLanguage {
		for i := range words {
			words[i].PartOfSpeech = t.posParser.Parse(words[i].PartOfSpeech)
			if index[words[i].Term] == nil {
				index[words[i].Term] = make(map[entities.Language]*entities.Word)
			}
			index[words[i].Term][lang] = &words[i]
		}
	}

	// We are looking for a translation for each word
	for lang, words := range wordsByLanguage {
		for _, word := range words {
			// Adding transfers
			translationsMap := make(entities.Translations)
			for targetLang, wordMap := range index[word.Term] {
				if targetLang != lang {
					translationsMap[targetLang] = append(translationsMap[targetLang], wordMap.Term) // Adding the full word
				}
			}
			word.Translations = translationsMap
			translations = append(translations, word)
		}
	}

	return translations
}

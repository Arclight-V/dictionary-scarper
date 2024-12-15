package repositories

import (
	"dictionary-scraper/common"
	"dictionary-scraper/entities"
	"errors"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	baseUrl = "https://dictionary.cambridge.org/dictionary/"
)

type DictionaryCambridgeParser struct{}

func NewDictionaryCambridgeParser() *DictionaryCambridgeParser {
	return &DictionaryCambridgeParser{}
}

// Extract the translation from French to English
func fetchTranslationFromFrenchToEnglis(word entities.Word) ([]string, string, error) {
	c := colly.NewCollector()
	c.UserAgent = DefaultUserAgent

	var firstMatchFound bool
	parser := common.GetPartOfSpeechParserInstance()
	var translations []string
	var transcription string
	c.OnHTML("div.pr.dictionary", func(e *colly.HTMLElement) {
		if firstMatchFound {
			return
		}
		partOfSpeach := e.DOM.Find("span.pos.dpos").First().Text()
		if parser.Parse(partOfSpeach) != word.PartOfSpeech {
			return
		}

		transcription = e.DOM.Find("span.pron.dpron").First().Text()

		// TODO: When you need several translations of one word, then change here
		translation := e.DOM.Find("span.trans.dtrans").First().Text()
		translations = append(translations, translation)
		log.Infof("partOfSpeach: %s, translation: %s , transcription: %s", partOfSpeach, translation, transcription)

		// Set flag, for stopping process
		firstMatchFound = true
	})

	var builder strings.Builder
	builder.WriteString(baseUrl)
	builder.WriteString(string(entities.French))
	builder.WriteRune('-')
	builder.WriteString(string(entities.English))
	builder.WriteRune('/')
	builder.WriteString(word.Term)

	err := c.Visit(builder.String())

	return translations, transcription, err

}

func (p *DictionaryCambridgeParser) FetchTranslations(word entities.Word, toLang entities.Language) ([]string, string, error) {
	if word.Term == "" {
		return nil, "", errors.New("term cannot be empty")
	}
	if word.PartOfSpeech == "" {
		return nil, "", errors.New("part_of_speech cannot be empty")
	}

	var translations []string
	var transcription string
	var err error
	if word.Language == entities.French {
		translations, transcription, err = fetchTranslationFromFrenchToEnglis(word)
		log.Infof("translations: %s, transcription: %s", translations, transcription)
		if err != nil {
			log.Error("Failed to visit URL:", err)
		}
	}

	return translations, transcription, err
}

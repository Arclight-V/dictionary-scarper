package entities

type Language string

const (
	French  Language = "french"
	English Language = "english"
	Russian Language = "russian"
)

type Translations = map[Language][]string

type Word struct {
	Language      Language
	Term          string
	PartOfSpeech  string
	Transcription string
	Translations  Translations
}

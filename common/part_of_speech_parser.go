package common

import (
	"strings"
	"sync"
)

// PartOfSpeechParser handles parsing of part of speech names
type PartOfSpeechParser struct {
	mappings map[string]string
}

var (
	instance *PartOfSpeechParser
	once     sync.Once
)

// NewPartOfSpeechParser creates a new instance of the parser
func GetPartOfSpeechParserInstance() *PartOfSpeechParser {
	once.Do(func() {
		instance = &PartOfSpeechParser{
			mappings: map[string]string{
				"noun":         "n",
				"n":            "n",
				"verb":         "v",
				"v":            "v",
				"adjective":    "adj",
				"adj":          "adj",
				"adverb":       "adv",
				"adv":          "adv",
				"conjunction":  "conj",
				"conj":         "conj",
				"pronoun":      "pron",
				"pron":         "pron",
				"preposition":  "prep",
				"prep":         "prep",
				"interjection": "interj",
				"interj":       "interj",
				"article":      "art",
				"art":          "art",
			},
		}
	})
	return instance
}

// Parse converts a full or shortened part of speech into a normalized format
func (p *PartOfSpeechParser) Parse(partOfSpeech string) string {
	partOfSpeech = strings.ToLower(strings.TrimSpace(partOfSpeech))
	if normalized, exists := p.mappings[partOfSpeech]; exists {
		return normalized
	}
	return "Unknown" // Default for unrecognized parts of speech
}

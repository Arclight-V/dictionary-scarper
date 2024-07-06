// scraper.go
package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	// import Colly

	"github.com/atselvan/ankiconnect"
	"github.com/gocolly/colly"
)

const (
	english = "english"
	french  = "french"
	russian = "russian"
)

type Word struct {
	Word, Transcription, Language, PartOfTheSpeech string
}

type Translate struct {
	From, To, Russian Word
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkLanguages(str string) string {
	switch str {
	case "fr-en":
		return french
	case "en-fr":
		return english
	case "en-ru":
		return russian
	default:
		panic(errors.New("The language " + str + " is not supported"))
	}
}

func setLanguages(str string) string {
	switch str {
	case "fr-en":
		return english
	case "en-fr":
		return french
	case "en-ru":
		return russian
	default:
		panic(errors.New("The language " + str + " is not supported"))
	}
}

func checkPartOfTheSpeeche(str string) string {
	switch str {
	case "noun":
	case "pronoun":
	case "verb":
	case "adjective":
	case "adverb":
	case "preposition":
	case "conjunction":
	case "interjection":
	case "article":
		return str
	default:
		panic(errors.New("Unknown part of speech " + str))
	}
	return str
}

func readWordsFromFile(inputFile *string) []Translate {
	readFile, err := os.Open(*inputFile)
	check(err)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var words []Translate
	for fileScanner.Scan() {
		word_info := strings.Fields(fileScanner.Text())
		if len(word_info) > 0 {
			word := Translate{
				From: Word{
					Word:            word_info[0],
					PartOfTheSpeech: checkPartOfTheSpeeche(word_info[1]),
					Language:        checkLanguages(word_info[2]),
				},
				To: Word{
					PartOfTheSpeech: checkPartOfTheSpeeche(word_info[1]),
					Language:        setLanguages(word_info[2]),
				},
				Russian: Word{
					PartOfTheSpeech: checkPartOfTheSpeeche(word_info[1]),
					Language:        russian,
				},
			}
			words = append(words, word)
		}
	}
	readFile.Close()
	return words

}

func searchOne(translate *Translate) {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	part_of_speach := translate.From.PartOfTheSpeech
	c.OnHTML("span.link.dlink", func(e *colly.HTMLElement) {
		if translate.To.Word == "" {
			if p_f_s := e.ChildText("span.pos.dpos"); strings.Contains(p_f_s, part_of_speach) {
				e.ForEachWithBreak("span.pron.dpron", func(i int, e *colly.HTMLElement) bool {
					translate.From.Transcription = e.Text
					return false
				})
				e.ForEachWithBreak("span.trans.dtrans", func(i int, e *colly.HTMLElement) bool {
					str := strings.TrimSpace(e.Text)
					str_split := strings.Split(str, "  ")
					if len(str_split) > 1 {
						translate.To.Word = str_split[0]
						translate.To.Transcription = str_split[1]
					} else {
						translate.To.Word = str_split[0]
					}
					return false
				})
			}
		}

	})
	url := "https://dictionary.cambridge.org/dictionary/" + translate.From.Language + "-" + translate.To.Language + "/" + translate.From.Word
	c.Visit(url)
}

func searchTwo(translate *Translate) {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	part_of_speach := translate.To.PartOfTheSpeech
	c.OnHTML("span.link.dlink", func(e *colly.HTMLElement) {
		if p_f_s := e.ChildText("span.pos.dpos"); strings.Contains(p_f_s, part_of_speach) {
			e.ForEachWithBreak("span.pron-info.dpron-info", func(i int, e *colly.HTMLElement) bool {
				translate.To.Transcription = e.Text + " " + translate.To.Transcription
				return false
			})
		}

	})
	url := "https://dictionary.cambridge.org/dictionary/" + translate.To.Language + "-" + translate.From.Language + "/" + translate.To.Word
	c.Visit(url)

}

func searchThree(translate *Translate) {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	c.OnHTML("div.di-body", func(e *colly.HTMLElement) {
		if translate.Russian.Word == "" {
			if p_f_s := e.ChildText("div.posgram.dpos-g.hdib.lmr-5"); strings.Contains(p_f_s, translate.Russian.PartOfTheSpeech) {
				e.ForEachWithBreak("span.trans.dtrans.dtrans-se", func(i int, e *colly.HTMLElement) bool {
					translate.Russian.Word = e.Text
					return false
				})
			}

		}
	})

	language := translate.From.Language
	search_word := translate.From.Word
	if language == french {
		language = translate.To.Language
		search_word = translate.To.Word
	}

	if strings.Contains(search_word, "to ") {
		search_word = strings.TrimPrefix(search_word, "to ")
	}

	url := "https://dictionary.cambridge.org/dictionary/" + language + "-" + translate.Russian.Language + "/" + search_word
	c.Visit(url)

}

func main() {

	fileToRead := flag.String("input", "/Users/vladimir/Documents/words.txt", "a file with a list of words to translate")
	flag.Parse()

	words_in_file := readWordsFromFile(fileToRead)

	for _, v := range words_in_file {
		searchOne(&v)
		searchTwo(&v)
		searchThree(&v)

		if v.From.Language == french {
			tmp := v.From
			v.From = v.To
			v.To = tmp
		}

		deck_name := ""
		if v.From.Word == "" || v.To.Word == "" || v.Russian.Word == "" {
			log.Println("Need work")
			deck_name = "Need Work"
		} else {
			deck_name = "New Deck"
		}

		client := ankiconnect.NewClient()
		note := ankiconnect.Note{
			DeckName:  deck_name,
			ModelName: "Basic (three reversed card)",
			Fields: ankiconnect.Fields{
				"Front": "<h1>" + v.From.Word + "</h1>" + "<br></br>" + v.From.PartOfTheSpeech + "<br />" + v.From.Transcription,
				"Back":  "<h1>" + v.To.Word + "</h1>" + "<br></br>" + v.To.PartOfTheSpeech + "<br />" + v.To.Transcription,
				"three": v.Russian.Word,
			},
		}
		restErr := client.Notes.Add(note)
		if restErr != nil {
			log.Println(restErr)
		}
	}
}

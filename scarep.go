// scraper.go
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// import Colly

	"github.com/atselvan/ankiconnect"
	"github.com/gocolly/colly"
)

// definr some data structures
// to store the scraped data

type Language int

const (
	english_french Language = iota
	french_english
	english_russian
)

type WordInfo struct {
	WordToTranslate, PartOfTheSpeache, Transcripton string
	Lang                                            Language
}

type TranslateWord struct {
	Word                      WordInfo
	Translate, AdditionalInfo string
}

type EnglishFrenchRussian struct {
	English TranslateWord
	French  TranslateWord
	Russian TranslateWord
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkLanguages(str string) Language {
	switch str {
	case "fr-en":
		return french_english
	case "en-fr":
		return english_french
	case "en-ru":
		return english_russian
	default:
		panic(errors.New("The language " + str + " is not supported"))
	}
}

func setTranslateLanguages(lang Language) Language {
	switch lang {
	case french_english:
		return english_french
	case english_french:
		return french_english
	default:
		return english_french
	}
}

func readWordsFromFile(inputFile *string) []WordInfo {
	readFile, err := os.Open(*inputFile)
	check(err)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var words []WordInfo
	for fileScanner.Scan() {
		word_info := strings.Fields(fileScanner.Text())
		if len(word_info) > 0 {
			word := WordInfo{
				WordToTranslate:  word_info[0],
				PartOfTheSpeache: word_info[1],
				Lang:             checkLanguages(word_info[2]),
			}
			words = append(words, word)
		}

	}
	readFile.Close()
	return words

}

func deleteToIsWordVerbEnglish(english_words []TranslateWord) {
	for i, v := range english_words {
		if strings.Contains(v.Word.PartOfTheSpeache, "verb") && strings.HasPrefix(v.Translate, "to ") {
			english_words[i].Translate = strings.TrimLeft(english_words[i].Translate, "to ")
		}
	}
}

func export(word_to_add EnglishFrenchRussian) {
	// --- create export dir
	path := "export"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	// --- export to CSV ---

	// open the output CSV file
	csvFile, csvErr := os.Create(path + "/" + word_to_add.English.Word.WordToTranslate + ".csv")
	// if the file creation fails
	if csvErr != nil {
		log.Fatalln("Failed to create the output CSV file", csvErr)
	}
	// release the resource allocated to handle
	// the file before ending the execution
	defer csvFile.Close()

	// create a CSV file writer
	writer := csv.NewWriter(csvFile)
	// release the resources associated with the
	// file writer before ending the execution
	defer writer.Flush()

	// add the header row to the CSV
	headers := []string{
		"word",
		"part of speaches",
		"transcription",
		"translate",
	}
	writer.Write(headers)

	// store each Industry product in the
	// output CSV file

	record := []string{
		word_to_add.English.Word.WordToTranslate,
		word_to_add.English.Word.PartOfTheSpeache,
		word_to_add.English.Word.Transcripton,
		word_to_add.English.Translate,
		word_to_add.French.Word.WordToTranslate,
		word_to_add.French.Word.PartOfTheSpeache,
		word_to_add.French.Word.Transcripton,
		word_to_add.French.Translate,
		word_to_add.Russian.Word.WordToTranslate,
		word_to_add.Russian.Word.PartOfTheSpeache,
		word_to_add.Russian.Word.Transcripton,
		word_to_add.Russian.Translate,
	}
	// add a new CSV record
	writer.Write(record)

	// --- export to JSON ---

	// open the output JSON file
	jsonFile, jsonErr := os.Create(path + "/" + word_to_add.English.Word.WordToTranslate + ".json")
	if jsonErr != nil {
		log.Fatalln("Failed to create the output JSON file", jsonErr)
	}
	defer jsonFile.Close()
	// convert industries to an indented JSON string
	jsonString, _ := json.MarshalIndent(word_to_add, " ", " ")

	// write the JSON string to file
	jsonFile.Write(jsonString)
}

func searchWord(lookingWord WordInfo) []TranslateWord {
	// initialize the struct slices

	var words []TranslateWord

	// initialize the Collector
	c := colly.NewCollector()

	// set a valid User-Agent header
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"

	// iterating over the list of industry card
	// HTML elements

	c.OnHTML("span.link.dlink", func(e *colly.HTMLElement) {
		if len(words) == 0 {
			if part_of_speaches := e.ChildText("div.dpos-g.hdib"); strings.Contains(part_of_speaches, lookingWord.PartOfTheSpeache) {
				word := TranslateWord{
					Word: WordInfo{
						WordToTranslate:  e.ChildText("h2"),
						PartOfTheSpeache: part_of_speaches,
						Transcripton:     e.ChildText("span.pron-info.dpron-info"),
						Lang:             lookingWord.Lang,
					},
					Translate: strings.Split(e.ChildText("div.def-body.ddef_b.ddef_b-t"), "\n")[0],
				}
				words = append(words, word)

			}
		}
	})

	c.OnHTML("div.pr.entry-body__el", func(e *colly.HTMLElement) {
		if len(words) == 0 {
			if lookingWord.Lang == english_russian {
				if part_of_speaches := e.ChildText("span.pos.dpos"); strings.Contains(part_of_speaches, lookingWord.PartOfTheSpeache) {
					word := TranslateWord{
						Word: WordInfo{
							WordToTranslate:  e.ChildText("span.hw.dhw"),
							PartOfTheSpeache: part_of_speaches,
							Transcripton:     e.ChildText("span.pron-info.dpron-info"),
							Lang:             lookingWord.Lang,
						},
						Translate: strings.Split(e.ChildText("div.def-body.ddef_b"), "\n")[0],
					}
					words = append(words, word)
				}
			}
		}
	})

	// connect to the target site
	var url string
	switch lookingWord.Lang {
	case english_french:
		url = "https://dictionary.cambridge.org/dictionary/english-french/" + lookingWord.WordToTranslate
	case french_english:
		url = "https://dictionary.cambridge.org/dictionary/french-english/" + lookingWord.WordToTranslate
	case english_russian:
		url = "https://dictionary.cambridge.org/dictionary/english-russian/" + lookingWord.WordToTranslate
	default:
		url = ""
	}
	c.Visit(url)

	if lookingWord.Lang == english_russian {
		words = append(words, TranslateWord{
			Word: lookingWord,
		})
	}

	deleteToIsWordVerbEnglish(words)

	return words
}

func main() {

	fileToRead := flag.String("input", "/Users/vladimir/Documents/words.txt", "a file with a list of words to translate")
	isExport := flag.Bool("export", false, "exporting json and csv files to the export directory")
	flag.Parse()

	words_in_file := readWordsFromFile(fileToRead)

	for _, v := range words_in_file {
		word_to_add_first := searchWord(v)
		var additional_info string
		if strings.Contains(word_to_add_first[0].Translate, "/") {
			additional_info = word_to_add_first[0].Translate
			word_to_add_first[0].Translate = word_to_add_first[0].Translate[:strings.Index(word_to_add_first[0].Translate, "/")]
			word_to_add_first[0].AdditionalInfo = additional_info
		}
		var word_to_translate_for_second string
		if strings.Contains(word_to_add_first[0].Translate, "[") {
			word_to_translate_for_second = word_to_add_first[0].Translate[:strings.Index(word_to_add_first[0].Translate, "[")]
		} else {
			word_to_translate_for_second = word_to_add_first[0].Translate
		}
		word_to_add_second := searchWord(WordInfo{
			WordToTranslate:  word_to_translate_for_second,
			PartOfTheSpeache: v.PartOfTheSpeache,
			Lang:             setTranslateLanguages(word_to_add_first[0].Word.Lang),
		})

		var word_to_translate_for_russian string
		if word_to_add_first[0].Word.Lang == english_french {
			word_to_translate_for_russian = word_to_add_first[0].Word.WordToTranslate

		} else {
			word_to_translate_for_russian = word_to_add_second[0].Word.WordToTranslate
		}

		word_to_add_third := searchWord(WordInfo{
			WordToTranslate:  word_to_translate_for_russian,
			PartOfTheSpeache: v.PartOfTheSpeache,
			Lang:             english_russian,
		})

		if len(word_to_add_first) == 0 || len(word_to_add_second) == 0 || len(word_to_add_third) == 0 {
			fmt.Println("ERROR", v, " not add")
			continue
		}

		var efr EnglishFrenchRussian

		if word_to_add_first[0].Word.Lang == english_french {
			efr.English = word_to_add_first[0]
			efr.French = word_to_add_second[0]
			efr.Russian = word_to_add_third[0]
		} else {
			efr.French = word_to_add_first[0]
			efr.English = word_to_add_second[0]
			efr.Russian = word_to_add_third[0]
		}

		if *isExport {
			export(efr)
		}

		client := ankiconnect.NewClient()

		note := ankiconnect.Note{
			DeckName:  "New Deck",
			ModelName: "Basic (three reversed card)",
			Fields: ankiconnect.Fields{
				"Front": "<h1>" + efr.English.Word.WordToTranslate + "</h1>" + "<br></br>" + efr.English.Word.PartOfTheSpeache + "<br />" + efr.English.Word.Transcripton,
				"Back":  "<h1>" + efr.French.Word.WordToTranslate + "</h1>" + "<br></br>" + efr.French.Word.PartOfTheSpeache + "<br />" + efr.French.Word.Transcripton,
				"three": efr.Russian.Translate,
			},
		}

		restErr := client.Notes.Add(note)
		if restErr != nil {
			log.Fatal(restErr)
		}
		fmt.Println(v)
	}
}

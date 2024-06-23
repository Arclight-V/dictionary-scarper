// scraper.go
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	// import Colly
	"github.com/gocolly/colly"
)

// definr some data structures
// to store the scraped data
type EnglishWord struct {
	Word, PartOfTheSpeache, Transcripton string
	Translate                            []string
}

type FrenchWord struct {
	Word, PartOfTheSpeache, Transcripton string
}

type EnglishFrench struct {
	English EnglishWord
	French  FrenchWord
}

type Word struct {
	WordToTranslate, PartOfTheSpeache, Language string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readWordsFromFile(inputFile *string) []Word {
	readFile, err := os.Open(*inputFile)
	check(err)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var words []Word
	for fileScanner.Scan() {
		word_info := strings.Fields(fileScanner.Text())
		word := Word{
			WordToTranslate:  word_info[0],
			PartOfTheSpeache: word_info[1],
			Language:         word_info[2],
		}
		words = append(words, word)
	}
	readFile.Close()
	return words

}

func getWordByPartOfSpeaches(english_words []EnglishWord, part_of_speaches string) []EnglishWord {
	var word []EnglishWord
	for _, v := range english_words {
		if v.PartOfTheSpeache == part_of_speaches {
			word = append(word, v)
		}
	}
	return word
}

func searchWord(lookingWord Word) EnglishFrench {
	// initialize the struct slices

	var englishWords []EnglishWord
	var words, part_of_the_speaches, transcriptons []string
	var translate [][]string

	// initialize the Collector
	c := colly.NewCollector()

	// set a valid User-Agent header
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"

	// iterating over the list of industry card
	// HTML elements

	c.OnHTML("h2.tw-bw.dhw.dpos-h_hw.di-title", func(e *colly.HTMLElement) {
		text := e.Text
		if text != "" {
			words = append(words, text)
		}
	})

	c.OnHTML("span.pos.dpos", func(e *colly.HTMLElement) {
		text := e.Text
		if text != "" {
			part_of_the_speaches = append(part_of_the_speaches, text)
		}
	})

	c.OnHTML("span.pron-info.dpron-info", func(e *colly.HTMLElement) {
		text := e.Text
		if text != "" {
			transcriptons = append(transcriptons, text)
		}
	})

	c.OnHTML("span.trans.dtrans.lmr--5", func(e *colly.HTMLElement) {
		text := e.Text
		texts := strings.Split(text, "/")
		if text != "" {
			translate = append(translate, texts)
		}
	})

	// connect to the target site
	var url string
	if lookingWord.Language == "en" {
		url = "https://dictionary.cambridge.org/dictionary/english-french/" + lookingWord.WordToTranslate
	} else {
		url = "https://dictionary.cambridge.org/dictionary/french-english/" + lookingWord.WordToTranslate
	}
	c.Visit(url)

	for _, w := range words {
		word := EnglishWord{
			Word: w,
		}
		englishWords = append(englishWords, word)
	}

	for i := 0; i < len(part_of_the_speaches); i++ {
		englishWords[i].PartOfTheSpeache = part_of_the_speaches[i]
	}

	for i := 0; i < len(transcriptons); i++ {
		englishWords[i].Transcripton = transcriptons[i]
	}

	for i := 0; i < len(translate); i++ {
		englishWords[i].Translate = translate[i]
	}

	word_to_add := getWordByPartOfSpeaches(englishWords, lookingWord.PartOfTheSpeache)
	// --- export to CSV ---

	// open the output CSV file
	csvFile, csvErr := os.Create(word_to_add[0].Word + ".csv")
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
		"In French",
	}
	writer.Write(headers)

	// store each Industry product in the
	// output CSV file
	for _, word := range word_to_add {
		// convert the Industry instance to
		// a slice of strings
		record := []string{
			word.Word,
			word.PartOfTheSpeache,
			word.Transcripton,
			word.Translate[0],
		}
		// add a new CSV record
		writer.Write(record)
	}

	// --- export to JSON ---

	// open the output JSON file
	jsonFile, jsonErr := os.Create(word_to_add[0].Word + ".json")
	if jsonErr != nil {
		log.Fatalln("Failed to create the output JSON file", jsonErr)
	}
	defer jsonFile.Close()
	// convert industries to an indented JSON string
	jsonString, _ := json.MarshalIndent(word_to_add, " ", " ")

	// write the JSON string to file
	jsonFile.Write(jsonString)

	var en_fr EnglishFrench
	return en_fr
}

func main() {

	fileToRead := flag.String("input", "/Users/vladimir/Documents/words.txt", "a file with a list of words to translate")
	flag.Parse()

	words_in_file := readWordsFromFile(fileToRead)
	for _, v := range words_in_file {
		searchWord(v)
	}

}
package stopwords

import (
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/jeffprestes/stopwords/corpus"
)

var languages = []string{
	"danish", "dutch", "english", "finnish", "french", "german", "hungarian",
	"italian", "norwegian", "portuguese", "russian", "spanish", "swedish",
	"turkish",
}

const punctuation = "!@£$%^&*()-_+¡€#,<.>/?`~'\"[{}];:\\|"

func (sf *StopwordFilter) isStopWord(token string) (wordStatus bool) {
	token = strings.ToLower(token)
	index := sort.SearchStrings(sf.words, token)
	wordStatus = index >= 0 && index < len(sf.words) && sf.words[index] == token
	if token == "n�o" {
		for i, b := range token {
			fmt.Printf("i: %d - %+v ", i, b)
		}
		fmt.Println(" ")
		fmt.Println("Token: ", token, " indice: ", index, " status: ", wordStatus)
	}
	return
}

type StopwordFilter struct {
	source   io.Reader
	dumpMode bool
	word     []byte
	j        int
	words    []string
}

func (sf *StopwordFilter) dump(p []byte) bool {
	sf.dumpMode = true

	for ; len(sf.word) > 0 && sf.j < len(p)-1; sf.j++ {

		poppedLetter := sf.word[0]
		sf.word = sf.word[1:]
		p[sf.j] = poppedLetter
	}
	if len(sf.word) == 0 {
		sf.dumpMode = false
	}
	return sf.dumpMode
}

func (sf *StopwordFilter) filterAndDump(p []byte) bool {
	// Delete the last character (space) when running the check
	word := sf.word
	if len(word) > 0 {
		word = word[:len(word)-1]
	}
	if sf.isStopWord(string(word)) {
		sf.word = make([]byte, 0, 0)
		return false
	}
	return sf.dump(p)
}

func (sf *StopwordFilter) Read(buffer []byte) (n int, err error) {
	letter := make([]byte, 1, 1)
	sf.j = 0

	if sf.dumpMode {
		sf.dump(buffer)
		return sf.j, nil
	}

	for {
		n, err := sf.source.Read(letter)
		if err == io.EOF {
			sf.filterAndDump(buffer)
			if sf.j == 0 {
				return 0, io.EOF
			}
			return sf.j, nil
		} else if n == 1 {
			if strings.Contains(punctuation, string(letter[0])) {
				// Skip letter if it is punctuation
				continue
			}
			// Accumulate word letters
			sf.word = append(sf.word, letter[0])

			// Dump all accumulated word letters and append a space if a space is
			// read
			if letter[0] == ' ' && sf.filterAndDump(buffer) {
				return sf.j, nil
			}
		}
	}
}

// Filter uses NewReader to perform non-streamed stop words filter
func Filter(str string, language corpus.Language) (string, error) {
	filter := NewReader(strings.NewReader(str), language)
	bytes, err := ioutil.ReadAll(filter)
	return string(bytes), err
}

// NewReader takes a Reader stream and exposes the Read method which filters
// out stop words
func NewReader(reader io.Reader, language corpus.Language) *StopwordFilter {
	return &StopwordFilter{reader, false, make([]byte, 0, 0), 0, corpus.Stopwords[language]}
}

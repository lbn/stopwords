package stopwords

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"
)

var languages = []string{
	"danish", "dutch", "english", "finnish", "french", "german", "hungarian",
	"italian", "norwegian", "portuguese", "russian", "spanish", "swedish",
	"turkish",
}

const punctuation = "!@£$%^&*()-_+¡€#,<.>/?`~'\"[{}];:\\|"

var words = make([]string, 0, 0)

func isStopWord(token string) bool {
	token = strings.ToLower(token)
	index := sort.SearchStrings(words, token)
	return index >= 0 && index < len(words) && words[index] == token
}

type StopwordFilter struct {
	source   io.Reader
	dumpMode bool
	word     []byte
	j        int
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
	if isStopWord(string(word)) {
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
			return sf.j, nil
		} else if n == 1 {
			if strings.Contains(punctuation, string(letter[0])) {
				// Skip letter if it is punctuation
				continue
			} else {
				// Accumulate word letters
				sf.word = append(sf.word, letter[0])
				if letter[0] == ' ' {
					// Dump all accumulated word letters and append a space
					if sf.filterAndDump(buffer) {
						return sf.j, nil
					}
				}
			}
		}
	}
}

func NewReader(reader io.Reader) *StopwordFilter {
	return &StopwordFilter{reader, false, make([]byte, 0, 0), 0}
}

func init() {
	for _, language := range languages {
		wordFile, _ := os.Open("./corpus/" + language)

		defer wordFile.Close()

		scanner := bufio.NewScanner(wordFile)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			words = append(words, scanner.Text())
		}
	}
	sort.Strings(words)
}

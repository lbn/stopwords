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

func (sf *StopwordFilter) dump(p []byte) {
	sf.dumpMode = true

	for ; len(sf.word) > 0 && sf.j < len(p)-1; sf.j++ {

		poppedLetter := sf.word[0]
		sf.word = sf.word[1:]
		p[sf.j] = poppedLetter
	}
	if len(sf.word) == 0 {
		sf.dumpMode = false
		p[sf.j] = ' '
		sf.j++
	}
}

func (sf *StopwordFilter) Read(p []byte) (n int, err error) {
	letter := make([]byte, 1, 1)
	if sf.dumpMode {
		sf.dump(p)
		return sf.j, nil
	}

	// Reset the pointer at the end
	defer func() {
		sf.j = 0
	}()

	for {
		n, err := sf.source.Read(letter)
		if err == io.EOF {
			if isStopWord(string(sf.word)) {
				sf.word = make([]byte, 0, 0)
			} else {
				sf.dump(p)
				if sf.dumpMode {
					return sf.j, nil
				}
			}
			return sf.j, nil
		} else if n == 1 {
			// If space read
			if strings.Contains(punctuation, string(letter[0])) {
				continue
			} else if letter[0] == ' ' {
				if isStopWord(string(sf.word)) {
					sf.word = make([]byte, 0, 0)
				} else {
					sf.dump(p)
					if sf.dumpMode {
						return sf.j, nil
					}
				}
			} else {
				sf.word = append(sf.word, letter[0])
			}
		} else if n == 0 {
			// True EOF
			return 0, io.EOF
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

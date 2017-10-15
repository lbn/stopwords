package stopwords

import (
	"strings"
	"testing"

	"github.com/jeffprestes/stopwords/corpus"
)

var (
	sentence   = "In 2009 for the first time in History the Spanish was the first \"mother tongue\" language of the western world."
	expected   = "2009 first time History Spanish first mother tongue language western world"
	sentencePT = "Em 2009 pela primeira vez na história o Espanhol se tornou a maior \"lingua mãe\" do mundo ocidental e não o Inglês."
	expectedPT = "2009 primeira vez história Espanhol tornou maior lingua mãe mundo ocidental inglês"
)

func TestStream(t *testing.T) {
	sen := strings.NewReader(sentence)
	filter := NewReader(sen, corpus.English)

	outdata := make([]byte, 0, 0)
	for {
		buf := make([]byte, 10)
		n, err := filter.Read(buf)
		if err != nil || n == 0 {
			break
		}

		outdata = append(outdata, buf[:n]...)
	}

	actualOutput := string(outdata)

	if actualOutput != expected {
		t.Logf("Expected: %s | Length: %d", expected, len(expected))
		t.Logf("Actual:   %s | Length: %d", actualOutput, len(actualOutput))
		t.Fatal("Actual output does not match the expected output")
	}
}

func TestStreamPt(t *testing.T) {
	//sen := strings.NewReader(corpus.RemoveDiacriticMark(sentencePT))
	sen := strings.NewReader(sentencePT)
	filter := NewReader(sen, corpus.Portuguese)

	outdata := make([]byte, 0, 0)
	for {
		buf := make([]byte, 10)
		n, err := filter.Read(buf)
		if err != nil || n == 0 {
			break
		}

		outdata = append(outdata, buf[:n]...)
	}

	actualOutput := string(outdata)

	if actualOutput != expected {
		t.Logf("Expected: %s | Length: %d", expectedPT, len(expectedPT))
		t.Logf("Actual:   %s | Length: %d", actualOutput, len(actualOutput))
		t.Fatal("Actual output does not match the expected output")
	}
}

func TestNonStream(t *testing.T) {
	actualOutput, err := Filter(sentence, corpus.English)
	if err != nil || actualOutput != expected {
		t.Logf("Expected: %s | Length: %d", expected, len(expected))
		t.Logf("Actual:   %s | Length: %d", actualOutput, len(actualOutput))
		t.Fatal("Actual output does not match the expected output")
	}
}

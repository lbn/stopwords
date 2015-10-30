package stopwords

import (
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	sentence := "In 2009 for the first time in History the Spanish was the first \"mother tongue\" language of the western world."
	expected := "2009 first time History Spanish first mother tongue language western world"

	sen := strings.NewReader(sentence)
	filter := NewReader(sen)

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

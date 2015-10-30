package stopwords

import (
	"fmt"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	sentence := "In 2009 for the first time in History the Spanish was the first \"mother tongue\" language of the western world."
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

	fmt.Println(string(outdata))
}

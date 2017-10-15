# stopwords
stopwords is a simple Go library to remove [stop
words](https://en.wikipedia.org/wiki/Stop_words) from a body of text.

## Instructions how to use
Please, check the stopwords_test.go

## Additional languages
Just add support to import other corpuses besides English, Spanish and Portuguese. PR are welcomed.

## Dependency
In your program you will need to create a folder called corpus and add the specific corpus file to it.

## Algorithm
Stop words are sorted on start up and stored in a slice. The algorithm
accumulates letters in memory until it finds the end of the word (space
character). Once the entire word is in memory we use binary search to test its
membership in the slice of stop words. The computational complexity for this
particular approach is `O(n log m)` where `n` is the number of words in a body
of text and `m` is the total number of stop words loaded.

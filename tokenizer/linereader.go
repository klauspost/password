// Tokenizers for various formats,
// that satisfies the password.Tokenizer interface.
package tokenizer

import (
	"bufio"
	"io"

	gzip "github.com/klauspost/pgzip"
)

// LineReader will return one line when Next() is called.
type LineReader struct {
	in io.Reader     // Original supplied reader
	gr *gzip.Reader  // Set if input is gzip compressed, otherwise nil
	br *bufio.Reader // Used for reading
}

// NewLine reads one password per line until
// \0xa (newline) is encountered.
// Input is streamed.
func NewLine(in io.Reader) *LineReader {
	l := &LineReader{in: in}
	l.br = bufio.NewReader(in)
	return l
}

// NewGzLine reads one password per line until
// \0xa (newline) is encountered.
// The input is assumed to be gzip compressed.
// Input is streamed.
func NewGzLine(in io.Reader) (*LineReader, error) {
	l := &LineReader{in: in}
	var err error
	l.gr, err = gzip.NewReader(in)
	if err != nil {
		return nil, err
	}
	l.br = bufio.NewReader(l.gr)
	return l, nil
}

// Next returns the data on the next line.
// Will return io.EOF when there is no more data.
func (l *LineReader) Next() (string, error) {
	record, err := l.br.ReadBytes(10)
	if err != nil {
		return "", err
	}
	return string(record), nil
}

// Should be called when finished
func (l *LineReader) Close() error {
	if l.gr != nil {
		return l.gr.Close()
	}
	return nil
}

// Tokenizers for various formats,
// that satisfies the password.Tokenizer interface.
package tokenizer

import (
	"bufio"
	"io"

	gzip "github.com/klauspost/pgzip"
)

type LineReader struct {
	in io.Reader
	gr *gzip.Reader
	br *bufio.Reader
}

// NewLine reads one password per line until
// \0xa (newline) is encountered.
// Input is streamed.
func NewLine(in io.Reader) (*LineReader, error) {
	l := &LineReader{in: in}
	l.br = bufio.NewReader(in)
	return l, nil
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

func (l *LineReader) Next() (string, error) {
	record, err := l.br.ReadBytes(10)
	if err != nil {
		return "", err
	}
	return string(record), nil
}

func (l *LineReader) Close() error {
	if l.gr != nil {
		return l.gr.Close()
	}
	return nil
}

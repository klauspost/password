package tokenizer

import (
	"bufio"
	"io"

	gzip "github.com/klauspost/pgzip"
)

type GzLineReader struct {
	in io.Reader
	gr *gzip.Reader
	br *bufio.Reader
}

func NewGzLine(in io.Reader) (*GzLineReader, error) {
	l := &GzLineReader{in: in}
	var err error
	l.gr, err = gzip.NewReader(in)
	if err != nil {
		return nil, err
	}
	l.br = bufio.NewReader(l.gr)
	return l, nil
}

func (l *GzLineReader) Next() (string, error) {
	record, err := l.br.ReadBytes(10)
	if err != nil {
		return "", err
	}
	return string(record), nil
}

func (l *GzLineReader) Close() error {
	return l.gr.Close()
}

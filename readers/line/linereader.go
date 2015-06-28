package line

import (
	"bufio"
	"io"

	gzip "github.com/klauspost/pgzip"
)

type Reader struct {
	in io.Reader
	gr *gzip.Reader
	br *bufio.Reader
}

func New(in io.Reader) (*Reader, error) {
	l := &Reader{in: in}
	var err error
	l.gr, err = gzip.NewReader(in)
	if err != nil {
		return nil, err
	}
	l.br = bufio.NewReader(l.gr)
	return l, nil
}

func (l *Reader) Next() (string, error) {
	record, err := l.br.ReadBytes(10)
	if err != nil {
		return "", err
	}
	return string(record), nil
}

func (l *Reader) Close() error {
	return l.gr.Close()
}

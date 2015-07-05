// Copyright 2015, Klaus Post, see LICENSE for details.

// Tokenizers for various formats,
// that satisfies the password.Tokenizer interface.
package tokenizer

import (
	"bufio"
	"compress/bzip2"
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
func NewLine(r io.Reader) *LineReader {
	l := &LineReader{in: r}
	l.br = bufio.NewReader(r)
	return l
}

// NewGzLine reads one password per line until
// \0xa (newline) is encountered.
// The input is assumed to be gzip compressed.
// Input is streamed.
func NewGzLine(r io.Reader) (*LineReader, error) {
	l := &LineReader{in: r}
	var err error
	l.gr, err = gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	l.br = bufio.NewReader(l.gr)
	return l, nil
}

// NewBz2Line reads one password per line until
// \0xa (newline) is encountered.
// The input is assumed to be bzip2 compressed.
// Input is streamed. If r does not also implement io.ByteReader,
// the decompressor may read more data than necessary from in.
func NewBz2Line(r io.Reader) *LineReader {
	l := &LineReader{in: r}
	bz := bzip2.NewReader(r)
	l.br = bufio.NewReader(bz)
	return l
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

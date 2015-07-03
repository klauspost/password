// Copyright 2015, Klaus Post, see LICENSE for details.

package password

import ()

// If your DbWriter implements this, input will be sent
// in batches instead of using Add.
type BulkWriter interface {
	AddMultiple([]string) error
}

type bulkWrapper struct {
	out BulkWriter
	res chan error
	in  chan []string
	buf []string
}

// BulkMax is the maximum number of passwords sent at once to the writer.
// You can change this before starting an import.
var BulkMax = 1000

func bulkWrap(out BulkWriter) DbWriter {
	b := &bulkWrapper{
		out: out,
		res: make(chan error, 1),
		in:  make(chan []string, 0),
		buf: make([]string, 0, BulkMax),
	}
	go func() {
		b.res <- nil
		defer close(b.res)
		for {
			select {
			case x, ok := <-b.in:
				if !ok {
					return
				}
				err := b.out.AddMultiple(x)
				b.res <- err
				if err != nil {
					return
				}
			}
		}
	}()
	return b
}

func (b *bulkWrapper) Add(s string) error {
	b.buf = append(b.buf, s)
	if len(b.buf) >= BulkMax {
		// Get last result
		last := <-b.res
		if last != nil {
			return last
		}
		// Send next
		b.in <- b.buf

		// Create new
		b.buf = make([]string, 0, BulkMax)
	}
	return nil
}

func (b *bulkWrapper) Close() error {
	if len(b.buf) > 0 {
		last := <-b.res
		if last != nil {
			return last
		}
		// Send next
		b.in <- b.buf
	}
	close(b.in)
	return <-b.res
}

package password

import ()

type BulkWriter interface {
	AddMultiple([]string) error
}

type bulkWrapper struct {
	out BulkWriter
	res chan error
	in  chan []string
	buf []string
}

var BulkMax = 1000

func BulkWrap(out BulkWriter) Writer {
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
	close(b.in)
	return <-b.res
}

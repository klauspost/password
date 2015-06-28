package password

import (
	"io"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

type Writer interface {
	Add(string) error
}

type DB interface {
	Has(string) (bool, error)
}

type Sanitizer interface {
	Sanitize(string) string
}

type Tokenizer interface {
	Next() (string, error)
}

type DefaultSanitizer struct{}

// Sanitize performs the following sanitation:
// * Trim space, tab and newlines from start+end of input
// * Check that there is at least 8 runes.
// * Normalize input using Unicode Normalization Form KD
// * Convert to unicode lower case.
// If input is less than 8 characters an empty string is returned.
func (d DefaultSanitizer) Sanitize(in string) string {
	in = strings.Trim(in, "\r\n \t")
	if utf8.RuneCountInString(in) < 8 {
		return ""
	}
	// Normalize using Unicode Normalization Form KD
	in = norm.NFKD.String(in)
	in = strings.ToLower(in)
	return in
}

// This will populate the known password list with common passwords
// It is a simple line-reader reading one password per line.
// Similar to format at https://crackstation.net/buy-crackstation-wordlist-password-cracking-dictionary.htm
func Import(in Tokenizer, out Writer, san Sanitizer) error {

	bulk, ok := out.(BulkWriter)
	if ok {
		closer, ok := out.(io.Closer)
		if ok {
			// TODO: Check error
			defer closer.Close()
		}
		out = BulkWrap(bulk)
	}

	closer, ok := out.(io.Closer)
	if ok {
		// TODO: Check error
		defer closer.Close()
	}

	if san == nil {
		san = DefaultSanitizer{}
	}

	start := time.Now()
	i := 0
	added := 0
	for {
		record, err := in.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		valstring := san.Sanitize(record)
		if len(valstring) > 0 {
			err = out.Add(valstring)
			if err != nil {
				return err
			}
			added++
		}
		i++
		if i%10000 == 0 {
			elapsed := time.Since(start)
			log.Printf("Read %d, (%0.0f per sec). Added: %d (%d%%)\n", i, float64(i)/elapsed.Seconds(), added, (added*100)/i)
		}
	}
	elapsed := time.Since(start)
	log.Printf("Processing took %s, processing %d entries.\n", elapsed, i)
	log.Printf("%0.2f entries/sec.", float64(i)/elapsed.Seconds())
	return nil
}

func InDB(password string, db DB, san Sanitizer) (bool, error) {
	if san == nil {
		san = DefaultSanitizer{}
	}
	p := san.Sanitize(password)
	if p == "" {
		return false, nil
	}
	return db.Has(p)
}

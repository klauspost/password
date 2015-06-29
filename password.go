package password

import (
	"errors"
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
	Sanitize(string) (string, error)
}

type Tokenizer interface {
	Next() (string, error)
}

// DefaultSanitizer should be used for adding passwords
// to the database.
// Assumes input is UTF8.
//
// DefaultSanitizer performs the following sanitazion:
//
//  - Trim space, tab and newlines from start+end of input
//  - Check that there is at least 8 runes (returns ErrSanitizeTooShort if not).
//  - Normalize input using Unicode Normalization Form KD
//  - Convert to unicode lower case.
//
// If input is less than 8 runes ErrSanitizeTooShort is returned.
var DefaultSanitizer Sanitizer

// CheckSanitizer should be used to sanitize a password
// before hasing. Performs the same operations as DefaultSanitizer
// except it doesn't convert the password to lower case.
// Assumes input is UTF8.
//
// CheckSanitizer performs the following sanitazion:
//
//  - Trim space, tab and newlines from start+end of input
//  - Check that there is at least 8 runes (returns ErrSanitizeTooShort if not).
//  - Normalize input using Unicode Normalization Form KD
var CheckSanitizer Sanitizer

func init() {
	CheckSanitizer = &checkSanitizer{}
	DefaultSanitizer = &defaultSanitizer{}
}

// ErrSanitizeTooShort is returned by the default sanitizer,
// if the input password is less than 8 runes.
var ErrSanitizeTooShort = errors.New("password too short")

// ErrPasswordInDB is returedn by Check, if the password is in the
// database.
var ErrPasswordInDB = errors.New("password found in database")

type defaultSanitizer struct {
	checkSanitizer
}

func (d defaultSanitizer) Sanitize(in string) (string, error) {
	in, err := d.checkSanitizer.Sanitize(in)
	in = strings.ToLower(in)
	return in, err
}

type checkSanitizer struct{}

func (c checkSanitizer) Sanitize(in string) (string, error) {
	in = strings.TrimSpace(in)
	if utf8.RuneCountInString(in) < 8 {
		return "", ErrSanitizeTooShort
	}
	in = norm.NFKD.String(in)
	return in, nil
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
		san = DefaultSanitizer
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

		valstring, err := san.Sanitize(record)
		if err == nil {
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

// Check a password.
// It will return an error if:
// * Sanitazition fails.
// * DB lookup returns an error
// * Password is in database (ErrPasswordInDB)
// If nil is passed as Sanitizer, DefaultSanitizer will be used.
func Check(password string, db DB, san Sanitizer) error {
	if san == nil {
		san = DefaultSanitizer
	}
	p, err := san.Sanitize(password)
	if err != nil {
		return err
	}
	has, err := db.Has(p)
	if err != nil {
		return err
	}
	if has {
		return ErrPasswordInDB
	}
	return nil
}

func inDB(password string, db DB, san Sanitizer) (bool, error) {
	if san == nil {
		san = DefaultSanitizer
	}
	p, err := san.Sanitize(password)
	if err != nil {
		return false, nil
	}
	return db.Has(p)
}

func SanitizeOK(password string, san Sanitizer) error {
	if san == nil {
		san = DefaultSanitizer
	}
	_, err := san.Sanitize(password)
	return err
}

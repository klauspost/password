// Copyright 2015, Klaus Post, see LICENSE for details.

// Dictionary Password Validation for Go
//
// For usage and examples see: https://github.com/klauspost/password
// (or open the README.md)
//
// This library will help you import a password dictionary and will allow you
// to validate new/changed passwords against the dictionary.
//
// You are able to use your own database and password dictionary.
// Currently the package supports importing dictionaries similar to
// CrackStation's Password Cracking Dictionary: https://crackstation.net/buy-crackstation-wordlist-password-cracking-dictionary.htm
//
// It and has "drivers" for various backends, see the "drivers"  directory, where there are
// implementations and a test framework that will help you test your own drivers.
package password

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// Logger used for output during Import.
// This can be exchanged with your own.
var Logger = log.New(os.Stdout, "", log.LstdFlags)

// A DbWriter is used for adding passwords to a database.
// Items sent to Add has always been sanitized, however
// the same passwords can be sent multiple times.
type DbWriter interface {
	Add(string) error
}

// A DB should check the database for the supplied password.
// The password sent to the interface has always been sanitized.
type DB interface {
	Has(string) (bool, error)
}

// A Sanitizer should prepare a password, and check
// the basic properties that should be satisfied.
// For an example, see DefaultSanitizer
type Sanitizer interface {
	Sanitize(string) (string, error)
}

// Tokenizer delivers input tokens (passwords).
// Calling Next() should return the next password, and when
// finished io.EOF should be returned.
//
// It is ok for the Tokenizer to send empty strings and duplicate
// values.
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
//  - Check that there is at least 8 runes. Return ErrSanitizeTooShort if not.
//  - Check that the input is valid utf8. Return ErrInvalidString if not.
//  - Normalize input using Unicode Normalization Form KD
//
// If input is less than 8 runes ErrSanitizeTooShort is returned.
var DefaultSanitizer Sanitizer

func init() {
	DefaultSanitizer = &defaultSanitizer{}
}

// ErrSanitizeTooShort is returned by the default sanitizer,
// if the input password is less than 8 runes.
var ErrSanitizeTooShort = errors.New("password too short")

// ErrInvalidString is returned by the default sanitizer
// if the string contains an invalid utf8 character sequence.
var ErrInvalidString = errors.New("invalid utf8 sequence")

// ErrPasswordInDB is returedn by Check, if the password is in the
// database.
var ErrPasswordInDB = errors.New("password found in database")

// doc at DefaultSanitizer
type defaultSanitizer struct{}

// doc at DefaultSanitizer
func (d defaultSanitizer) Sanitize(in string) (string, error) {
	in = strings.TrimSpace(in)
	if utf8.RuneCountInString(in) < 8 {
		return "", ErrSanitizeTooShort
	}
	if !utf8.ValidString(in) {
		return "", ErrInvalidString
	}
	in = norm.NFKD.String(in)
	in = strings.TrimSpace(in)
	return in, nil
}

type initer interface {
	Init() error
}

// Import will populate a database with common passwords.
//
// You must supply a Tokenizer (see tokenizer package for default tokenizers)
// that will deliver the passwords,
// a DbWriter, where the passwords will be sent,
// and finally a Sanitizer to clean up the passwords -
// - if you send nil DefaultSanitizer will be used.
func Import(in Tokenizer, out DbWriter, san Sanitizer) (err error) {
	bulk, ok := out.(BulkWriter)
	if ok {
		initer, ok := out.(initer)
		if ok {
			err := initer.Init()
			if err != nil {
				return err
			}
		}
		closer, ok := out.(io.Closer)
		if ok {
			defer func() {
				e := closer.Close()
				if e != nil && err == nil {
					err = e
				}
			}()
		}
		out = bulkWrap(bulk)
	}

	initer, ok := out.(initer)
	if ok {
		err := initer.Init()
		if err != nil {
			return err
		}
	}

	closer, ok := out.(io.Closer)
	if ok {
		defer func() {
			e := closer.Close()
			if e != nil && err == nil {
				err = e
			}
		}()
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
			valstring = strings.ToLower(valstring)
			err = out.Add(valstring)
			if err != nil {
				return err
			}
			added++
		}
		i++
		if i%10000 == 0 {
			elapsed := time.Since(start)
			Logger.Printf("Read %d, (%0.0f per sec). Added: %d (%d%%)\n", i, float64(i)/elapsed.Seconds(), added, (added*100)/i)
		}
	}
	elapsed := time.Since(start)
	Logger.Printf("Processing took %s, processing %d entries.\n", elapsed, i)
	Logger.Printf("%0.2f entries/sec.", float64(i)/elapsed.Seconds())
	return nil
}

// Check a password against the database.
// It will return an error if:
//  - Sanitazition fails.
//  - DB lookup returns an error
//  - Password is in database (ErrPasswordInDB)
// If nil is passed as Sanitizer, DefaultSanitizer will be used.
func Check(password string, db DB, san Sanitizer) error {
	if san == nil {
		san = DefaultSanitizer
	}
	p, err := san.Sanitize(password)
	if err != nil {
		return err
	}
	p = strings.ToLower(p)
	has, err := db.Has(p)
	if err != nil {
		return err
	}
	if has {
		return ErrPasswordInDB
	}
	return nil
}

// Sanitize will sanitize a password, useful before hashing
// and storing it.
//
// If the sanitizer is nil, DefaultSanitizer will be used.
func Sanitize(password string, san Sanitizer) (string, error) {
	if san == nil {
		san = DefaultSanitizer
	}
	p, err := san.Sanitize(password)
	return p, err
}

// SanitizeOK can be used to check if a password passes the sanitizer.
//
// If the sanitizer is nil, DefaultSanitizer will be used.
func SanitizeOK(password string, san Sanitizer) error {
	if san == nil {
		san = DefaultSanitizer
	}
	_, err := san.Sanitize(password)
	return err
}

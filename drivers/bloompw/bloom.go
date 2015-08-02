// Copyright 2015, Mike Houston, see LICENSE for details.

// A bitset Bloom filter for a reduced memory password representation.
//
// See https://github.com/AndreasBriese/bbloom
package bloompw

import (
	"github.com/AndreasBriese/bbloom"
)

type BloomPW struct {
	Filter *bbloom.Bloom
}

// New will return a new Database interface that stores entries in
// a bloom filter
func New(filter *bbloom.Bloom) (*BloomPW, error) {
	b := &BloomPW{Filter: filter}

	return b, nil
}

// Has satisfies the password.DB interface
func (b BloomPW) Has(s string) (bool, error) {
	return b.Filter.Has([]byte(s)), nil
}

// Has satisfies the password.DbWriter interface.
// It writes a single password to the database
func (b BloomPW) Add(s string) error {
	b.Filter.Add([]byte(s))
	return nil
}

// AddMultiple satisfies the password.BulkWriter interface.
// It writes a number of passwords to the database
func (b BloomPW) AddMultiple(s []string) error {
	for _, v := range s {
		if err := b.Add(v); err != nil {
			return err
		}
	}
	return nil
}

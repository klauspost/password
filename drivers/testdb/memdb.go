// Copyright 2015, Klaus Post, see LICENSE for details.

// An in-memory database for testing
//
// This database is completely in memory
// and can be used as a reference for your
// own implementation.
//
// Since it has rather high memory use, it is
// not recommended to be used in production.
// For a good in-memory database, see the bloom
// driver.
package testdb

// This is the simplest possible database that must be supported.
// If you can mimmic this with your database you are good to go!
type MemDB map[string]struct{}

// NewMemDB creates a new MemDB instance
func NewMemDB() *MemDB {
	m := MemDB(make(map[string]struct{}))
	return &m
}

// Add a password to the MemDB.
// It must silently ignore duplicates
// If an error is returned Import is aborted.
func (m *MemDB) Add(s string) error {
	db := *m
	db[s] = struct{}{}
	return nil
}

// Has will check if the database has a specific password.
// If any error is returned, it will be forwarded to your
// "Check()" call.
func (m MemDB) Has(s string) (bool, error) {
	_, ok := m[s]
	return ok, nil
}

// MemDBBulk is the same as MemDB, but also
// satisfies the bulk interface.
type MemDBBulk map[string]struct{}

// NewMemDBBulk will return a new MemDBBulk
func NewMemDBBulk() *MemDBBulk {
	m := MemDBBulk(make(map[string]struct{}))
	return &m
}

// AddMultiple is the function that will be called
// with several items at once.
func (m *MemDBBulk) AddMultiple(s []string) error {
	db := *m
	for _, p := range s {
		db[p] = struct{}{}
	}
	return nil
}

// Add a single entry
func (m *MemDBBulk) Add(s string) error {
	db := *m
	db[s] = struct{}{}
	return nil
}

// Has returns true if the map has the string
func (m MemDBBulk) Has(s string) (bool, error) {
	_, ok := m[s]
	return ok, nil
}

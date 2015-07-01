package testdb

import (
	"testing"

	"github.com/klauspost/password/drivers"
)

// Test a MemDB database
func TestMemDB(t *testing.T) {
	err := drivers.TestDriver(NewMemDB())
	if err != nil {
		t.Fatal(err)
	}
}

// Test a MemDBBulk database
func TestMemDBBulk(t *testing.T) {
	err := drivers.TestDriver(NewMemDBBulk())
	if err != nil {
		t.Fatal(err)
	}
}

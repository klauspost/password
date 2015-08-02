// Copyright 2015, Mike Houston, see LICENSE for details.

package bloompw

import (
	"testing"

	"github.com/AndreasBriese/bbloom"
	"github.com/klauspost/password/drivers"
)

// Test a bloom database
func TestBloom(t *testing.T) {
	// create a bloom filter for 65536 items and 0.001 % wrong-positive ratio
	filter := bbloom.New(float64(1<<16), float64(0.00001))

	bloom, err := New(&filter)
	if err != nil {
		t.Fatal(err)
	}
	err = drivers.TestDriver(bloom)
	if err != nil {
		t.Fatal(err)
	}
}

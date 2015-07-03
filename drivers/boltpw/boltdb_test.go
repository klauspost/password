// Copyright 2015, Klaus Post, see LICENSE for details.

package boltpw

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/klauspost/password/drivers"
)

// tempfile returns a temporary file path.
func tempfile() string {
	f, _ := ioutil.TempFile("", "bolt-")
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

// Test a bolt database
func TestBolt(t *testing.T) {
	db, err := bolt.Open(tempfile(), 0666, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(db.Path())
	defer db.Close()

	bolt, err := New(db, "commonpwd")
	if err != nil {
		t.Fatal(err)
	}
	err = drivers.TestDriver(bolt)
	if err != nil {
		t.Fatal(err)
	}
}

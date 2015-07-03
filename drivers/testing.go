// Copyright 2015, Klaus Post, see LICENSE for details.

// Provides a standard test library for drivers.
//
// You can use the functions of this package to test
// your own driver implementations.
//
// If your driver provides both read and write functionality
// and no inbetween synchronization you can use the TestDriver
// function, otherwise read and write is split into TestImport
// and TestData
package drivers

import (
	"github.com/klauspost/password"
	"github.com/klauspost/password/testdata"
	"github.com/klauspost/password/tokenizer"

	"bytes"
	"fmt"
)

type TestDB interface {
	password.DbWriter
	password.DB
}

// TestDriver will test a driver by
// running TestImport followed by TestData
func TestDriver(db TestDB) error {
	err := TestImport(db)
	if err != nil {
		return err
	}
	err = TestData(db)
	if err != nil {
		return err
	}
	return nil
}

var single_val = "j984lop!#\"{}"

// TestImport will import about 1500 entries into your
// database. It will test "Add" and "AddMultiple" (if available).
// If any error is returned the test failed.
func TestImport(db password.DbWriter) error {
	buf, err := testdata.Asset("testdata.txt.gz")
	if err != nil {
		return err
	}
	in, err := tokenizer.NewGzLine(bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	writer, ok := db.(password.DbWriter)
	if !ok {
		return fmt.Errorf("%T is not a DbWriter", db)
	}
	err = password.Import(in, writer, nil)
	if err != nil {
		return err
	}
	// Test Add once separately
	err = writer.Add(single_val)
	if err != nil {
		return err
	}
	// .. and test that it is ok to write it again.
	err = writer.Add(single_val)
	if err != nil {
		return err
	}
	return nil
}

// TestData will test that the data imported with TestImport
// is correctly returned.
// It will test "Has" function of the driver.
// If any error is returned the test failed.
func TestData(db password.DB) error {
	for p := range testdata.TestSet {
		if password.SanitizeOK(p, nil) != nil {
			continue
		}
		err := password.Check(p, db, nil)
		if err != password.ErrPasswordInDB {
			return fmt.Errorf("%s not found in database: %v", p, err)
		}
		if err != password.ErrPasswordInDB && err != nil {
			return fmt.Errorf("check %s returned unexpected error: %v", p, err)
		}
	}
	for p := range testdata.NotInSet {
		if password.SanitizeOK(p, nil) != nil {
			continue
		}
		err := password.Check(p, db, nil)
		if err == password.ErrPasswordInDB {
			return fmt.Errorf("%s should NOT be not found in database: %v", p, err)
		} else if err != nil {
			return err
		}
	}

	has, err := db.Has(single_val)
	if !has {
		return fmt.Errorf("%s not found in database. (single insert)", single_val)
	}
	if err != nil {
		return err
	}
	has, err = db.Has(single_val + "*")
	if has {
		return fmt.Errorf("%s* WAS found in database, it shouldn't. (single insert)", single_val)
	}
	if err != nil {
		return err
	}
	return nil
}

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

// TestDriver will test a driver.
func TestDriver(db TestDB) error {
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
	// Test Add once separately
	val := "j984lop!#\"{}"
	err = writer.Add(val)
	if err != nil {
		return err
	}
	has, err := db.Has(val)
	if !has {
		return fmt.Errorf("%s not found in database. (single insert)", val)
	}
	if err != nil {
		return err
	}
	has, err = db.Has(val + "*")
	if has {
		return fmt.Errorf("%s* WAS found in database, it shouldn't. (single insert)", val)
	}
	if err != nil {
		return err
	}
	return nil
}

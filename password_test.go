package password

import (
	"bytes"
	"os"
	"testing"

	"github.com/klauspost/password/drivers/testdb"
	"github.com/klauspost/password/readers/line"
	"github.com/klauspost/password/testdata"
)

func TestImport(t *testing.T) {
	buf := testdata.MustAsset("testdata.txt.gz")
	mem := testdb.NewMemDB()
	in, err := line.New(bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal(err)
	}
	err = Import(in, mem, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestImportBig(t *testing.T) {
	r, err := os.Open("crackstation-human-only.txt.gz")
	if err != nil {
		t.Skip("Skipping big file test. 'crackstation-human-only.txt.gz' must be in current dir")
	}
	mem := testdb.NewMemDB()
	in, err := line.New(r)
	if err != nil {
		t.Fatal(err)
	}
	err = Import(in, mem, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInDB(t *testing.T) {
	buf := testdata.MustAsset("testdata.txt.gz")
	mem := testdb.NewMemDB()
	in, err := line.New(bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal(err)
	}
	err = Import(in, mem, nil)
	for p := range testdata.TestSet {
		if SanitizeOK(p, nil) != nil {
			continue
		}
		has, err := inDB(p, mem, nil)
		if err != nil {
			t.Fatal(err)
		}
		if !has {
			t.Fatalf("db should have: %s", p)
		}
		err = Check(p, mem, nil)
		if err != ErrPasswordInDB {
			t.Fatal("check failed on:", p, err)
		}
	}
	for p := range testdata.NotInSet {
		if SanitizeOK(p, nil) != nil {
			continue
		}
		has, err := inDB(p, mem, nil)
		if err != nil {
			t.Fatal(err)
		}
		if has {
			t.Fatalf("db should not have: %s", p)
		}
		err = Check(p, mem, nil)
		if err != nil {
			t.Fatal("check failed on:", p, err)
		}
	}
}

func TestDefaultSanitizer(t *testing.T) {
	san := DefaultSanitizer
	for p := range testdata.TestSet {
		s, err := san.Sanitize(p)
		expect, ok := testdata.SanitizeExpect[p]
		if !ok {
			t.Fatalf("Sanitized version of `%s` not defined.", p)
		}
		if s != expect.S {
			t.Fatalf("Sanitized difference. Expected `%s`, got `%s`", expect.S, s)
		}
		e := ""
		if err != nil {
			e = err.Error()
		}
		if e != expect.E {
			t.Fatalf("Sanitized error difference. Expected `%s`, got `%s`", expect.E, e)
		}
	}
}

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

func TestCheck(t *testing.T) {
	buf := testdata.MustAsset("testdata.txt.gz")
	mem := testdb.NewMemDB()
	in, err := line.New(bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal(err)
	}
	err = Import(in, mem, nil)
	san := DefaultSanitizer{}
	for p := range testdata.TestSet {
		pass := san.Sanitize(p)
		if pass == "" {
			continue
		}
		has, err := mem.Has(pass)
		if err != nil {
			t.Fatal(err)
		}
		if !has {
			t.Fatalf("db should have: %s", pass)
		}
	}
	for p := range testdata.NotInSet {
		pass := san.Sanitize(p)
		if pass == "" {
			continue
		}
		has, err := mem.Has(pass)
		if err != nil {
			t.Fatal(err)
		}
		if has {
			t.Fatalf("db should not have: %s", pass)
		}
	}
}

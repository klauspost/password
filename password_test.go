package password

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/klauspost/password/drivers/testdb"
	"github.com/klauspost/password/readers/line"
	"github.com/klauspost/password/testdata"
)

func TestImport(t *testing.T) {
	buf, err := testdata.Asset("testdata.txt.gz")
	if err != nil {
		t.Fatal(err)
	}
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
	mem := testdb.NewMemDBBulk()
	in, err := line.New(r)
	if err != nil {
		t.Fatal(err)
	}
	err = Import(in, mem, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestImportBulk(t *testing.T) {
	buf, err := testdata.Asset("testdata.txt.gz")
	if err != nil {
		t.Fatal(err)
	}
	mem := testdb.NewMemDBBulk()
	in, err := line.New(bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal(err)
	}
	err = Import(in, mem, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInDB(t *testing.T) {
	buf, err := testdata.Asset("testdata.txt.gz")
	if err != nil {
		t.Fatal(err)
	}
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

type CustomSanitizer struct {
	email    string
	username string
}

func (c CustomSanitizer) Sanitize(s string) (string, error) {
	s, err := DefaultSanitizer.Sanitize(s)
	if err != nil {
		return "", err
	}
	if strings.EqualFold(s, c.email) {
		return "", errors.New("password cannot be the same as email")
	}
	if strings.EqualFold(s, c.username) {
		return "", errors.New("password cannot be the same as user name")
	}
	return s, nil
}

// This example shows how to create a custom sanitizer that checks if
// the password matches the username or email.
//
// CustomSanitizer is defined as:
//  type CustomSanitizer struct {
//      email string
//      username string
//  }
//
//  func (c CustomSanitizer) Sanitize(s string) (string, error) {
//      s, err := DefaultSanitizer.Sanitize(s)
//      if err != nil {
//          return "", err
//      }
//      if strings.EqualFold(s, c.email) {
//          return "", errors.New("password cannot be the same as email")
//      }
//      if strings.EqualFold(s, c.username) {
//          return "", errors.New("password cannot be the same as user name")
//      }
//      return s, nil
//  }
func ExampleSanitizer() {
	// Create a custom sanitizer.
	san := CustomSanitizer{email: "john@doe.com", username: "johndoe73"}

	// Check some passwords
	err := SanitizeOK("john@doe.com", san)
	fmt.Println(err)

	err = SanitizeOK("JohnDoe73", san)
	fmt.Println(err)

	err = SanitizeOK("MyP/|$$W0rd", san)
	fmt.Println(err)
	// Output: password cannot be the same as email
	// password cannot be the same as user name
	// <nil>
}

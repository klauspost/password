# password
[![GoDoc][1]][2] [![Build Status][3]][4]

[1]: https://godoc.org/github.com/klauspost/password?status.svg
[2]: https://godoc.org/github.com/klauspost/password
[3]: https://travis-ci.org/klauspost/password.svg?branch=master
[4]: https://travis-ci.org/klauspost/password

Dictionary Password Validation for Go.

Protect your bcrypt/scrypt/PBKDF encrypted passwords against dictionary attacks.

Motivated by [Password Requirements Done Better](http://blog.klauspost.com/password-requirements-done-better/) - or *why password requirements help hackers*

This library will help you import a password dictionary and will allow you to validate new/changed passwords against the dictionary.

You are able to use your own database and password dictionary. Currently the package supports importing dictionaries similar to [CrackStation's Password Cracking Dictionary](https://crackstation.net/buy-crackstation-wordlist-password-cracking-dictionary.htm), and has "drivers" for [MongoDB](https://godoc.org/github.com/klauspost/password/drivers/mgopw), [BoltDB](https://godoc.org/github.com/klauspost/password/drivers/boltpw), [MySQL](https://godoc.org/github.com/klauspost/password/drivers/sqlpw) and [PostgreSQL](https://godoc.org/github.com/klauspost/password/drivers/sqlpw).

# installation

As always, the package is installed with `go get github.com/klauspost/password`.

# usage

With this library you can

1. Import a password dictionary into your database
2. Check new passords against the dictionary
3. Sanitize passwords before authenticating a user

All of the 3 functionality parts can be used or replaced as it suits your application. In particult you probably do not want to import dictionaries on your webserver, so you can separate that functionality into a separate command.

## setting up a database


```Go
import(
  "github.com/boltdb/bolt"
  "github.com/klauspost/password"
  "github.com/klauspost/password/drivers/boltpw"
)

  // Open the database using the Bolt driver
  // You probably have this elsewhere if you already use Bolt
  db, err := bolt.Open("password.db", 0666, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
```

So far pretty standard. We open the database as we always would. This is used by the driver in [`github.com/klauspost/password/drivers/boltpw`](https://godoc.org/github.com/klauspost/password/drivers/boltpw) to write and check passwords.

```Go
  // Use the driver to create in/out using the bucket "commonpwd"
	chk, err := boltpw.New(db, "commonpwd")
	if err != nil {
		t.Fatal(err)
	} 
```

The object we get back can then be used to check passwords, assuming you have imported a database.
```Go
	err = password.Check(chk, "SecretPassword", nil)
	if err != nil {
	  // Password failed sanitazion or was in database.
	  panic(err)
	}
```	

## importing a dictionary

Example that will import the crackstation into memory. Replace `testdb.NewMemDBBulk()` with a constructor to the database you want to use.
```Go
import (
	"os"
	
	"github.com/klauspost/password"
	"github.com/klauspost/password/drivers/testdb"
	"github.com/klauspost/password/tokenizer"
)

func Import() {
	r, err := os.Open("crackstation-human-only.txt.gz")
	if err != nil {
	  panic(err)
	}
	mem := testdb.NewMemDBBulk()
	in, err := tokenizer.NewGzLine(r)
	if err != nil {
		panic(err)
	}
	err = password.Import(in, mem, nil)
	if err != nil {
		panic(err)
	}
}

```
## checking a password

This is an example of checking and preparing a password to be stored in the database. Passwords allowed to be full UTF8, and are compared case insensitively.
```Go
func PreparePassword(db password.DB, toCheck string)  (string, error) {
	err := password.Check(db, toCheck, nil)
	if err != nil {
		// Password failed sanitazion or was in database.
		return "", err
	}
	
	// We use the default sanitizer to sanitize/normalize the password
	toStore, _ := password.Sanitize(toCheck, nil)
	if err != nil {
		// Shouldn't happen, since we already passed sanitaztion in the check once
		// File a bug if it does.
		panic(err)
	}

	// bcrypt the result and return it
	return bcrypt.GenerateFromPassword([]byte(toStore), 12)
}
```	

# dictionaries
* [CrackStation's Password Cracking Dictionary](https://crackstation.net/buy-crackstation-wordlist-password-cracking-dictionary.htm). Very good start.
* HTTP download provided by me, use only if you cannot download torrents:
 * [Actual Passwords only](http://5.9.40.76/static/dicts/crackstation-human-only.txt.gz) - 200MB.
 * [Full Dictionary](http://5.9.40.76/static/dicts/crackstation.full.txt.gz) - 3.5GB
  
# [SkullSecurity Passwords](https://wiki.skullsecurity.org/Passwords). 
Mostly small and varying quality.

# [g0tmi1k Dictionaries + Wordlists](https://blog.g0tmi1k.com/2011/06/dictionaries-wordlists/). 
Hard to download. `18-in-1` has a lot of `sameword1`; `sameword2`, etc. Mostly ascii passwords.

# compatibility

Unless security related issues should show up, the interfaces and functions should not change in this package. If it is impossible to remain compatible, it will always be shown by a compiler error. So if the library compiles after an update it will remain compatible.

# license

This code is published under an MIT license. See LICENSE file for more information.


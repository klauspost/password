// Copyright 2015, Klaus Post, see LICENSE for details.

// Driver for MongoDB
//
// Tested on Mongo v3.0.4 and 2.6.x
//
// Supply a session and the database and collection name
// you would like to use.
package mgopw

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo can be used for adding and checking passwords.
// Note that passwords are cut until they are below 512 bytes
// when checking against the database, since the maximum index length
// is 1024 bytes including some overhead.
type Mongo struct {
	session    *mgo.Session
	db         string
	collection string
}

// New returns a new database.
// Supply a valid (copy of a) session and the database and collection
// you would like to use.
func New(session *mgo.Session, db, collection string) *Mongo {
	m := Mongo{
		session:    session,
		db:         db,
		collection: collection,
	}
	return &m
}

// Add an entry to the password database
func (m Mongo) Add(s string) error {
	s = truncate(s)
	_, err := m.session.DB(m.db).C(m.collection).UpsertId(s, bson.M{"_id": s})
	return err
}

// Has will return true if the database has the entry.
func (m Mongo) Has(s string) (bool, error) {
	s = truncate(s)
	n, err := m.session.DB(m.db).C(m.collection).FindId(s).Count()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Cut runes off the end until the
// string is below 512 bytes.
func truncate(s string) string {
	for {
		b := []byte(s)
		if len(b) < 512 {
			return s
		}
		r := []rune(s)
		// Cut one rune
		s = string(r[:len(r)-1])
	}
}

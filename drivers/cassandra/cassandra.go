// Copyright 2015, Klaus Post, see LICENSE for details.

// Package cassandra is a driver for Apache Cassandra
//
// Supply a session and the database and collection name
// you would like to use.
package cassandra

import "github.com/gocql/gocql"

// Cassandra can be used for adding and checking passwords.
type Cassandra struct {
	session *gocql.Session
	table   string
}

// New returns a new database.
// Supply a valid (copy of a) session and the database and collection
// you would like to use.
// The correct keyspace must be set on the session.
// Before using the database, the table should be created:
//
//		create table keyspace.table(password text, PRIMARY KEY(password));
//
// Replace the keyspace and table with the keyspace and table you want to use.
func New(session *gocql.Session, table string) *Cassandra {
	m := Cassandra{
		session: session,
		table:   table,
	}
	return &m
}

// Add an entry to the password database
func (m Cassandra) Add(s string) error {
	return m.session.Query(`INSERT INTO `+m.table+` (password) VALUES (?)`, s).Exec()
}

// Has will return true if the database has the entry.
func (m Cassandra) Has(s string) (bool, error) {
	n := 0
	if err := m.session.Query(`SELECT COUNT(*) FROM `+m.table+` WHERE password = ?`, s).
		Consistency(gocql.One).Scan(&n); err != nil {
		return false, err
	}

	return n != 0, nil
}

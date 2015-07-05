// Copyright 2015, Klaus Post, see LICENSE for details.

// Wrapper for an SQL database backend
//
// This can be used to use an existing database for input
// output.
//
// There is constructors for
//
// See "mysql_test.go" and "postgres_test.go" for examples on
// how to create those.
//
// Note that passwords are truncated at 64 runes (not bytes).
package sqlpw

import (
	"database/sql"
)

// Sql can be used for adding and checking passwords.
// Insert and Query are generated for MySQL, and should very likely
// be changed for other databases. See the "postgres_test" for an example.
type Sql struct {
	db     *sql.DB
	Table  string
	Query  string // Query string, used to get a count of hits
	Insert string // Insert string,used to insert an item
	qStmt  *sql.Stmt
	iStmt  *sql.Stmt
}

// New returns a new database.
//
func New(db *sql.DB, table, query, insert string) *Sql {
	s := Sql{
		db:     db,
		Table:  table,
		Query:  query,
		Insert: insert,
	}
	return &s
}

// NewMysql returns a new database wrapper, set up for MySQL.
//
func NewMysql(db *sql.DB, table, column string) *Sql {
	s := Sql{
		db:     db,
		Table:  table,
		Query:  "SELECT COUNT(*) FROM `" + table + "` WHERE `" + column + "`=?;",
		Insert: "INSERT IGNORE INTO `" + table + "` (`" + column + "`) VALUE (?);",
	}
	return &s
}

func NewPostgresql(db *sql.DB, table, column string) *Sql {
	s := Sql{
		db:     db,
		Table:  table,
		Query:  `INSERT INTO ` + table + ` (` + column + `) VALUES ($1)`,
		Insert: `SELECT COUNT(*) FROM  ` + table + ` WHERE ` + column + `=$1`,
	}
	return &s
}

// Add an entry to the password database
func (m *Sql) Add(s string) error {
	var err error
	if m.iStmt == nil {
		m.iStmt, err = m.db.Prepare(m.Insert)
		if err != nil {
			return err
		}
	}
	_, err = m.iStmt.Exec(truncate(s))
	return err
}

// Has will return true if the database has the entry.
func (m *Sql) Has(s string) (bool, error) {
	var err error
	if m.qStmt == nil {
		m.qStmt, err = m.db.Prepare(m.Query)
		if err != nil {
			return false, err
		}
	}
	var num int
	err = m.qStmt.QueryRow(truncate(s)).Scan(&num)
	if err != nil {
		return false, err
	}
	return num > 0, nil
}

func truncate(s string) string {
	r := []rune(s)
	if len(r) <= 64 {
		return s
	}
	return string(r[:64])
}

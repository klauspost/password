// Copyright 2015, Klaus Post, see LICENSE for details.

// Wrapper for an SQL database backend
//
// This can be used to use an existing database for input
// output.
//
// Without modifying the default queries
// the library assumes a table with a unique column "Pass"
// is created already.
//
// See "mysql_test.go" and "postgres_test.go" for examples on
// how to create those
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
	Query  string // Query string, used to get a count of hits, it defaults to "SELECT COUNT(*) FROM `" + table + "` WHERE `Pass`=?;"
	Insert string // Insert string,used to insert an item, defaults to "INSERT IGNORE INTO `" + table + "` (`Pass`) VALUE (?);"
	qStmt  *sql.Stmt
	iStmt  *sql.Stmt
}

// New returns a new database.
//
func New(db *sql.DB, table string) *Sql {
	s := Sql{
		db:     db,
		Table:  table,
		Query:  "SELECT COUNT(*) FROM `" + table + "` WHERE `Pass`=?;",
		Insert: "INSERT IGNORE INTO `" + table + "` (`Pass`) VALUE (?);",
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
	_, err = m.iStmt.Exec(s)
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
	err = m.qStmt.QueryRow(s).Scan(&num)
	if err != nil {
		return false, err
	}
	return num > 0, nil
}

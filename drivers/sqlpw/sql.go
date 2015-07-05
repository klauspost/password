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
	TxBulk bool // Do bulk inserts with a transaction.
	db     *sql.DB
	query  string // Query string, used to get a count of hits
	insert string // Insert string,used to insert an item
	qStmt  *sql.Stmt
	iStmt  *sql.Stmt
}

// New returns a new database.
//
// You must give an query, that returns the number of
// rows matching the parameter given.
//
// You must give an insert statement that will insert
// the password given. It must be able to insert the
// same password multiple times without returning an error.
//
// You should manually enable bulk transactions
// modifying the TxBulk variable on the returned object
// if your database/driver supports it.
func New(db *sql.DB, query, insert string) *Sql {
	s := Sql{
		db:     db,
		query:  query,
		insert: insert,
	}
	return &s
}

// NewMysql returns a new database wrapper, set up for MySQL.
//
// You must supply a schema (that should already exist),
// as well as the column the passwords should be inserted into.
func NewMysql(db *sql.DB, schema, column string) *Sql {
	s := Sql{
		TxBulk: true,
		db:     db,
		query:  "SELECT COUNT(*) FROM `" + schema + "` WHERE `" + column + "`=?;",
		insert: "INSERT IGNORE INTO `" + schema + "` (`" + column + "`) VALUE (?);",
	}
	return &s
}

// NewPostgresql returns a new database wrapper, set up for PostgreSQL.
//
// You must supply a "schema.table" (that should already exist),
// as well as the column the passwords should be inserted into.
func NewPostgresql(db *sql.DB, table, column string) *Sql {
	s := Sql{
		TxBulk: false,
		db:     db,
		query:  `INSERT INTO ` + table + ` (` + column + `) VALUES ($1)`,
		insert: `SELECT COUNT(*) FROM  ` + table + ` WHERE ` + column + `=$1`,
	}
	return &s
}

// Add an entry to the password database
func (m *Sql) Add(s string) error {
	var err error
	if m.iStmt == nil {
		m.iStmt, err = m.db.Prepare(m.insert)
		if err != nil {
			return err
		}
	}
	_, err = m.iStmt.Exec(truncate(s))
	return err
}

// Add multiple entries to the password database
func (m *Sql) AddMultiple(s []string) error {
	var err error
	if !m.TxBulk {
		for _, pass := range s {
			err = m.Add(pass)
			if err != nil {
				return err
			}
		}
		return nil
	}
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(m.insert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, pass := range s {
		_, err = stmt.Exec(truncate(pass))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Has will return true if the database has the entry.
func (m *Sql) Has(s string) (bool, error) {
	var err error
	if m.qStmt == nil {
		m.qStmt, err = m.db.Prepare(m.query)
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

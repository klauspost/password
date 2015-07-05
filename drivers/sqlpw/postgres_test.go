// Copyright 2015, Klaus Post, see LICENSE for details.

package sqlpw

import (
	"database/sql"
	"flag"
	"testing"

	"github.com/klauspost/password/drivers"
	_ "github.com/lib/pq"
)

var postGresPwd = flag.String("pgpass", "", "Postgres password")

func init() {
	flag.Parse()
}

// Test a Postgres database
// To run locally, use the "-pgpass" password to set the postgres user password.
func TestPostgres(t *testing.T) {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password="+*postGresPwd)
	if err != nil {
		t.Skip("Postgres connect error:", err)
	}
	err = db.Ping()
	if err != nil {
		t.Skip("Postgres ping error:", err)
	}

	table := "testschema.pwtesttable"
	drop := `DROP TABLE ` + table + `;`
	schema := `CREATE SCHEMA testschema AUTHORIZATION postgres`
	create := `CREATE TABLE ` + table + ` ("pass" VARCHAR(64) PRIMARY KEY);`
	ignore_rule := `
        CREATE OR REPLACE RULE db_table_ignore_duplicate_inserts AS
            ON INSERT TO ` + table + `
            WHERE (EXISTS ( 
                SELECT 1
                FROM ` + table + `
                WHERE ` + table + `.pass = NEW.pass
            ) ) DO INSTEAD NOTHING;`

	_, _ = db.Exec(drop)
	_, err = db.Exec(schema)
	if err != nil {
		t.Log("warning:", err)
	}
	_, err = db.Exec(create)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(ignore_rule)
	if err != nil {
		t.Fatal(err)
	}

	d := NewPostgresql(db, table, "pass")

	err = drivers.TestDriver(d)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(drop)
	if err != nil {
		t.Log("DROP returned:", err)
	}
}

// Example of using a Postgres database
//
// Uses 'pwtesttable' in the 'testschema' schema,
// and reads/adds to the "pass" column.
//
// Table can be created like this:
//  `CREATE TABLE ` + table + ` ("`+ column +`" VARCHAR(64) PRIMARY KEY);`
//
// For Postgres to ignore duplicate inserts, you can use a rule
// like this:
//
//  `CREATE OR REPLACE RULE db_table_ignore_duplicate_inserts AS
//      ON INSERT TO ` + table + `
//      WHERE (EXISTS (
//          SELECT 1
//          FROM ` + table + `
//          WHERE ` + table + `.` + column + ` = NEW.` + column + `
//      )
//  ) DO INSTEAD NOTHING;`
func ExampleNewPostgresql() {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic("Postgres connect error " + err.Error())
	}
	table := "testschema.pwtesttable"

	d := NewPostgresql(db, table, "pass")

	// Test it
	err = drivers.TestDriver(d)
	if err != nil {
		panic(err)
	}
}

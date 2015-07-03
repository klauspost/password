package sqlpw

import (
	"database/sql"
	"testing"

	"github.com/klauspost/password/drivers"
	_ "github.com/lib/pq"
)

// Test a Postgres database
func TestPostgres(t *testing.T) {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		t.Skip("Postgres connect error ", err)
	}
	table := "testschema.pwtesttable"
	drop := `DROP TABLE ` + table + `;`
	schema := `CREATE SCHEMA IF NOT EXISTS testschema AUTHORIZATION postgres`
	create := `CREATE TABLE ` + table + ` ("pass" VARCHAR(128) PRIMARY KEY);`
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
		t.Fatal(err)
	}
	_, err = db.Exec(create)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(ignore_rule)
	if err != nil {
		t.Fatal(err)
	}

	d := New(db, table)
	// Override Insert/Query
	d.Insert = `INSERT INTO ` + table + ` (pass) VALUES ($1)`
	d.Query = `SELECT COUNT(*) FROM  ` + table + ` WHERE pass=$1`
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
func ExampleNew_Postgres() {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic("Postgres connect error " + err.Error())
	}
	table := "testschema.pwtesttable"

	d := New(db, table)
	// Override Insert/Query
	d.Insert = `INSERT INTO ` + table + ` (pass) VALUES ($1)`
	d.Query = `SELECT COUNT(*) FROM  ` + table + ` WHERE pass=$1`
	err = drivers.TestDriver(d)
	if err != nil {
		panic(err)
	}
}

// Copyright 2015, Klaus Post, see LICENSE for details.

package sqlpw

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/klauspost/password/drivers"
)

// Test a MySQL database
func TestMySQL(t *testing.T) {
	db, err := sql.Open("mysql", "travis:@/testdb")
	if err != nil {
		t.Skip("MySQL connect error:", err)
	}
	err = db.Ping()
	if err != nil {
		t.Skip("MySQL ping error:", err)
	}

	drop := "DROP TABLE `testdb`.`test-table`;"
	create := "CREATE TABLE `testdb`.`test-table` (`Pass` VARCHAR(64) NOT NULL COMMENT '', PRIMARY KEY (`Pass`)  COMMENT '', UNIQUE INDEX `Pass_UNIQUE` (`Pass` ASC)  COMMENT '') DEFAULT CHARACTER SET = utf8;"

	_, _ = db.Exec(drop)
	_, err = db.Exec(create)
	if err != nil {
		t.Fatal(err)
	}

	d := NewMysql(db, "test-table", "Pass")
	err = drivers.TestDriver(d)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(drop)
	if err != nil {
		t.Log("DROP returned:", err)
	}
}

// Example of using a MySQL database.
//
// This assumes "test-table" has been created in "testdb".
// It could have been created like this:
//
//	CREATE TABLE `testdb`.`test-table` (
//    `Pass` VARCHAR(64) NOT NULL,
//    PRIMARY KEY (`Pass`),
//    UNIQUE INDEX `Pass_UNIQUE` (`Pass` ASC)
//  ) DEFAULT CHARACTER SET = utf8;
func ExampleNewMysql() {
	db, err := sql.Open("mysql", "travis:@/testdb")
	if err != nil {
		panic("MySQL connect error " + err.Error())
	}

	d := NewMysql(db, "test-table", "Pass")

	// Test it
	err = drivers.TestDriver(d)
	if err != nil {
		panic(err)
	}
}

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
		t.Skip("MySQL connect error ", err)
	}
	drop := "DROP TABLE `testdb`.`test-table`;"
	create := "CREATE TABLE `testdb`.`test-table` (`Pass` VARCHAR(128) NOT NULL COMMENT '', PRIMARY KEY (`Pass`)  COMMENT '', UNIQUE INDEX `Pass_UNIQUE` (`Pass` ASC)  COMMENT '') DEFAULT CHARACTER SET = utf8;"

	_, _ = db.Exec(drop)
	_, err = db.Exec(create)
	if err != nil {
		t.Fatal(err)
	}

	d := New(db, "test-table")
	err = drivers.TestDriver(d)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(drop)
	if err != nil {
		t.Log("DROP returned:", err)
	}
}

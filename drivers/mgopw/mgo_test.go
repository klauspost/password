package mgopw

import (
	"testing"

	"github.com/klauspost/password/drivers"
	"gopkg.in/mgo.v2"
)

// Test a Mongo database
func TestMemDB(t *testing.T) {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		t.Skip("No database: ", err)
	}
	coll := session.DB("testdb").C("password-test")
	_ = coll.DropCollection()

	db := New(session, "testdb", "password-test")
	err = drivers.TestDriver(db)
	if err != nil {
		t.Fatal(err)
	}
	err = coll.DropCollection()
	if err != nil {
		t.Log("Drop returned", err, "(ignoring)")
	}
	session.Close()
}

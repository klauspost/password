// Copyright 2015, Klaus Post, see LICENSE for details.

package cassandra

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/klauspost/password/drivers"
)

// Test a Cassandra database
func TestCassandra(t *testing.T) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "testkeyspace"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		t.Fatal("createsession", err)
	}
	defer session.Close()

	err = session.Query(`CREATE KEYSPACE IF NOT EXISTS testkeyspace with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`).Exec()
	if err != nil {
		t.Fatal("createkeyspace", err)
	}
	_ = session.Query(`DROP TABLE IF EXISTS passwords`).Exec()

	err = session.Query(`CREATE TABLE testkeyspace.passwords(password text, PRIMARY KEY(password));`).Exec()
	if err != nil {
		t.Fatal("createkeyspace", err)
	}

	defer session.Query(`DROP TABLE IF EXISTS passwords`).Exec()

	db := New(session, "passwords")
	err = drivers.TestImport(db)
	if err != nil {
		t.Fatal(err)
	}

	err = drivers.TestData(db)
	if err != nil {
		t.Fatal(err)
	}
}

package mgopw

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongo struct {
	session    *mgo.Session
	db         string
	collection string
}

func New(session *mgo.Session, db, collection string) *Mongo {
	m := Mongo{
		session:    session,
		db:         db,
		collection: collection,
	}
	return &m
}

func (m Mongo) Add(s string) error {
	_, err := m.session.DB(m.db).C(m.collection).UpsertId(s, bson.M{"_id": s})
	return err
}

func (m Mongo) Has(s string) (bool, error) {
	n, err := m.session.DB(m.db).C(m.collection).FindId(s).Count()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

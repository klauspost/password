package testdb

type MemDB map[string]struct{}

func NewMemDB() *MemDB {
	m := MemDB(make(map[string]struct{}))
	return &m
}

func (m *MemDB) Add(s string) error {
	db := *m
	db[s] = struct{}{}
	return nil
}

func (m MemDB) Has(s string) (bool, error) {
	_, ok := m[s]
	return ok, nil
}

type MemDBBulk map[string]struct{}

func NewMemDBBulk() *MemDBBulk {
	m := MemDBBulk(make(map[string]struct{}))
	return &m
}

func (m *MemDBBulk) AddMultiple(s []string) error {
	db := *m
	for _, p := range s {
		db[p] = struct{}{}
	}
	return nil
}
func (m *MemDBBulk) Add(s string) error {
	db := *m
	db[s] = struct{}{}
	return nil
}

func (m MemDBBulk) Has(s string) (bool, error) {
	_, ok := m[s]
	return ok, nil
}

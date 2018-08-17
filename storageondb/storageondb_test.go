package storageondb

import (
	"testing"

	"github.com/go-pg/pg"
)

func TestUserStorageDB_setDB(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	type args struct {
		db *pg.DB
	}
	u := &UserStorageDB{}
	if u.db != nil {
		t.Errorf("Expected db == nil")
	}
	u.setDB(&pg.DB{})
	if u.db == nil {
		t.Errorf("Expected db != nil")
	}
}

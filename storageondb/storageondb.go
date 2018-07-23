package storageondb

import (
	"fmt"

	"github.com/go-pg/pg"

	"reactionbot/environment"
)

const (
	USERS_TABLE = "users"
)

type users struct {
	ID    string
	Token string
}

type temporaryUserTokens struct {
	Uuid   string
	UserID string
}

type UserStorageDB struct {
	db *pg.DB
}

func SetUp() (*UserStorageDB, error) {
	userStorage := UserStorageDB{}

	connStr := environment.GetPostgresDBURL()
	opts, err := pg.ParseURL(connStr)

	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)

	userStorage.setDB(db)
	return &userStorage, nil
}

func (u *UserStorageDB) setDB(db *pg.DB) {
	u.db = db
}

func (u *UserStorageDB) AddUserToken(id string, token string) error {
	return u.db.Insert(&users{ID: id, Token: token})
}

func (u *UserStorageDB) LookupUserToken(id string) (bool, error) {
	user := users{}
	count, err := u.db.Model(&user).Where("id = ?", id).Count()

	if err != nil {
		return false, nil
	}
	if count > 1 {
		panic(fmt.Sprintf("Wrong count value %d when looking for user %s", count, id))
	}
	return count == 1, nil
}

func (u *UserStorageDB) remove(model interface{}) error {
	return u.db.Delete(model)
}

func (u *UserStorageDB) RemoveUserToken(id string) error {
	found, err := u.LookupUserToken(id)

	if !found {
		return fmt.Errorf("No user %s in table users", id)
	} else if err != nil {
		return err
	}

	return u.remove(&users{ID: id})
}

func (u *UserStorageDB) GetUserToken(id string) (string, error) {
	user := users{ID: id}

	if err := u.db.Select(&user); err != nil {
		return "", err
	}
	return user.Token, nil
}

func (u *UserStorageDB) PopUserToken(id string) (string, error) {

	token, err := u.GetUserToken(id)
	if err != nil {
		return "", err
	}
	if err = u.RemoveUserToken(id); err != nil {
		return "", err
	}
	return token, nil
}

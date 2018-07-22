package storageondb

import (
	"fmt"

	"github.com/go-pg/pg"

	// _ "github.com/lib/pq"

	// "github.com/go-pg/pg"
	// "github.com/go-pg/pg/orm"

	"reactionbot/environment"
)

const (
	USERS_TABLE = "users"
)

type users struct {
	ID    string
	Token string
}

type UserStorageDB struct {
	db *pg.DB
}

func SetUp() (*UserStorageDB, error) {
	userStorage := UserStorageDB{}

	connStr := environment.GetPostgresDBURL()
	fmt.Println("connstr", connStr)
	opts, err := pg.ParseURL(connStr)
	fmt.Println("parsed", opts, err)

	if err != nil {
		fmt.Println("errore in parse url", opts, err)
		return nil, err
	}

	db := pg.Connect(opts)
	fmt.Println("post connect")

	if err == nil {
		userStorage.setDB(db)
		return &userStorage, nil
	}
	return nil, err
}

func (u *UserStorageDB) setDB(db *pg.DB) {
	u.db = db
}

func (u *UserStorageDB) Add(code string, value string) error {
	_, err := u.db.Exec("INSERT INTO users VALUES ($1, $2)", code, value)
	return err
}

func (u *UserStorageDB) AddUser(id string, token string) error {
	return u.db.Insert(&users{ID: id, Token: token})
}

func (u *UserStorageDB) lookupInTable(code string, table string) (bool, error) {
	// rows, err := u.db.Query("SELECT COUNT(*) as count FROM $1 WHERE id=$2", table, code)
	return false, nil
	// defer rows.Close()

	// if ok := rows.Next(); !ok {
	// 	panic("No result for lookup")
	// }

	// var count int
	// rows.Scan(&count)
	// if ok := rows.Next(); ok {
	// 	panic("Too much rows for select count")
	// }
	// return count == 1, err
}

func (u *UserStorageDB) LookupUser(id string) (bool, error) {
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

func (u *UserStorageDB) Lookup(code string) (bool, error) {
	return false, nil

	// rows, err := u.db.Query("SELECT COUNT(*) as count FROM users WHERE id=$1", code)
	// defer rows.Close()

	// if ok := rows.Next(); !ok {
	// 	panic("No result for lookup")
	// }

	// var count int
	// rows.Scan(&count)
	// if ok := rows.Next(); ok {
	// 	panic("Too much rows for select count")
	// }
	// return count == 1, err
}

func (u *UserStorageDB) remove(model interface{}) error {
	return u.db.Delete(model)
}

func (u *UserStorageDB) RemoveUser(id string) error {
	found, err := u.LookupUser(id)

	if !found {
		return fmt.Errorf("No user %s in table users", id)
	} else if err != nil {
		return err
	}

	return u.remove(&users{ID: id})
}

func (u *UserStorageDB) Remove(code string) error {

	ok, err := u.Lookup(code)
	fmt.Println("looking for", code, ok, err)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("No code %s", code)
	}

	// code = pq.QuoteIdentifier(code)
	_, err = u.db.Exec("DELETE FROM users WHERE id=$1", code)
	return err
}

func (u *UserStorageDB) GetUserToken(id string) (string, error) {
	user := users{ID: id}

	if err := u.db.Select(&user); err != nil {
		return "", nil
	}
	return user.Token, nil
}

func (u *UserStorageDB) Get(code string) (string, error) {
	return "", nil

	// code = pq.QuoteIdentifier(code)
	// rows, err := u.db.Query("SELECT token FROM users WHERE id=$1", code)
	// defer rows.Close()

	// if err != nil {
	// 	return "", err
	// }

	// if ok := rows.Next(); !ok {
	// 	return "", fmt.Errorf("Wrong number of lines to get token from table users")
	// }

	// var token string
	// rows.Scan(&token)

	// return token, nil
}

func (u *UserStorageDB) Pop(code string) (string, error) {

	value, err := u.Get(code)
	if err != nil {
		return "", fmt.Errorf("Unable to pop %s", code)
	}

	err = u.Remove(code)
	if err != nil {
		return "", err
	}
	return value, nil
}

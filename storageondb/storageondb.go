package storageondb

import (
	"database/sql"
	"fmt"

	"reactionbot/environment"

	"github.com/lib/pq"
)

type UserStorageDB struct {
	db *sql.DB
}

func SetUp() (UserStorageDB, error) {
	userStorage := UserStorageDB{}

	connStr := environment.GetPostgresDBURL()
	db, err := sql.Open("postgres", connStr)

	if err == nil {
		userStorage.setDB(db)
	}
	return userStorage, err
}

func (u *UserStorageDB) setDB(db *sql.DB) {
	u.db = db
}

func (u *UserStorageDB) Add(code string, value string) error {
	code = pq.QuoteIdentifier(code)
	value = pq.QuoteIdentifier(value)

	_, err := u.db.Exec("INSERT INTO users VALUES ($1, $2)", code, value)

	return err
}

func (u *UserStorageDB) Lookup(code string) (bool, error) {
	code = pq.QuoteIdentifier(code)

	rows, err := u.db.Query("SELECT COUNT(*) as count FROM users WHERE id=$1", code)
	defer rows.Close()

	if ok := rows.Next(); !ok {
		panic("No result for lookup")
	}

	var count int
	rows.Scan(&count)
	if ok := rows.Next(); ok {
		panic("Too much rows for select count")
	}
	return count == 1, err
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

	code = pq.QuoteIdentifier(code)
	_, err = u.db.Exec("DELETE FROM users WHERE id=$1", code)
	return err
}

func (u *UserStorageDB) Get(code string) (string, error) {

	code = pq.QuoteIdentifier(code)
	rows, err := u.db.Query("SELECT token FROM users WHERE id=$1", code)
	defer rows.Close()

	if err != nil {
		return "", err
	}

	if ok := rows.Next(); !ok {
		return "", fmt.Errorf("Wrong number of lines to get token from table users")
	}

	var token string
	rows.Scan(&token)

	return token, nil
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
	return value, err
}

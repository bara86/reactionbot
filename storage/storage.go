package storage

import (
	"fmt"
	"sync"
)

type UserStorage struct {
	keys sync.Map
}

func (u *UserStorage) LookupUser(user string) bool {
	_, ok := u.keys.Load(user)
	return ok
}

func (u *UserStorage) AddUser(user string, token string) error {
	u.keys.Store(user, token)
	return nil
}

func (u *UserStorage) RemoveUser(user string) error {
	u.keys.Delete(user)
	return nil
}

func (u *UserStorage) GetUser(user string) (string, error) {
	value, ok := u.keys.Load(user)
	if !ok {
		return "", fmt.Errorf("No user %s", user)
	}

	return value.(string), nil
}

func (u *UserStorage) PopUser(user string) (string, error) {
	value, err := u.GetUser(user)
	if err == nil {
		u.RemoveUser(user)
	}
	return value, err
}

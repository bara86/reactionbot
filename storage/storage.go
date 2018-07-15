package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reactionbot/environment"
	"sync"
)

type UserStorage struct {
	keys sync.Map
}

func SetUp() (*UserStorage, error) {
	userStorage := UserStorage{}

	err := userStorage.setUp()
	return &userStorage, err
}

func (u *UserStorage) setUp() error {
	u.keys = sync.Map{}

	fileName := environment.GetSaveFileName()
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	f.Close()

	if err != nil {
		return err
	}

	return u.loadFromFile(fileName)
}

func (u *UserStorage) loadFromFile(fileName string) error {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		return nil
	}
	return u.unmarshalJSON(content)
}

func (u *UserStorage) unmarshalJSON(data []byte) error {
	var tmpMap map[string]interface{}
	if err := json.Unmarshal(data, &tmpMap); err != nil {
		return err
	}

	for key, value := range tmpMap {
		u.keys.Store(key, value)
	}
	return nil
}

func (u *UserStorage) marshalJSON() ([]byte, error) {
	tmpMap := map[string]string{}

	u.keys.Range(func(k interface{}, v interface{}) bool {
		tmpMap[k.(string)] = v.(string)
		return true
	})

	return json.Marshal(tmpMap)
}

func (u *UserStorage) saveMap() error {

	marshalled, err := u.marshalJSON()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(environment.GetSaveFileName(), marshalled, 0644)

}

func (u *UserStorage) LookupUser(user string) bool {
	_, ok := u.keys.Load(user)
	return ok
}

func (u *UserStorage) AddUser(user string, token string) error {
	u.keys.Store(user, token)
	return u.saveMap()
}

func (u *UserStorage) RemoveUser(user string) error {
	u.keys.Delete(user)
	return u.saveMap()
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

package storageonfile

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
	lock sync.Mutex
}

func SetUp() (*UserStorage, error) {
	userStorage := UserStorage{}

	err := userStorage.setUp()
	return &userStorage, err
}

func (u *UserStorage) setUp() error {
	if v, _ := environment.GetSaveOnFile(); !v {
		return nil
	}

	fileName := environment.GetSaveFileName()
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	f.Close()

	if err != nil {
		return err
	}

	return u.loadFromFile(fileName)
}

func (u *UserStorage) loadFromFile(fileName string) error {
	u.lock.Lock()
	content, err := ioutil.ReadFile(fileName)
	u.lock.Unlock()
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

	u.lock.Lock()
	err = ioutil.WriteFile(environment.GetSaveFileName(), marshalled, 0644)
	u.lock.Unlock()
	return err
}

func (u *UserStorage) LookupUserToken(code string) (bool, error) {
	_, ok := u.keys.Load(code)
	return ok, nil
}

func (u *UserStorage) AddUserToken(code string, value string) error {
	u.keys.Store(code, value)
	return u.saveMap()
}

func (u *UserStorage) RemoveUserToken(code string) error {
	u.keys.Delete(code)
	return u.saveMap()
}

func (u *UserStorage) GetUserToken(code string) (string, error) {
	value, ok := u.keys.Load(code)
	if !ok {
		return "", fmt.Errorf("No code %s", code)
	}

	return value.(string), nil
}

func (u *UserStorage) PopUserToken(code string) (string, error) {
	value, err := u.GetUserToken(code)
	if err == nil {
		u.RemoveUserToken(code)
	}
	return value, err
}

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

func (u *UserStorage) Lookup(code string) bool {
	_, ok := u.keys.Load(code)
	return ok
}

func (u *UserStorage) Add(code string, value string) error {
	u.keys.Store(code, value)
	return u.saveMap()
}

func (u *UserStorage) Remove(code string) error {
	u.keys.Delete(code)
	return u.saveMap()
}

func (u *UserStorage) Get(code string) (string, error) {
	value, ok := u.keys.Load(code)
	if !ok {
		return "", fmt.Errorf("No code %s", code)
	}

	return value.(string), nil
}

func (u *UserStorage) Pop(code string) (string, error) {
	value, err := u.Get(code)
	if err == nil {
		u.Remove(code)
	}
	return value, err
}
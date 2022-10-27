package server

import (
	"encoding/json"
	"strconv"
)

type StringInt int
type StringFloat float64

type Objects struct {
	Users map[string]User `json:"users"`
	Apps  map[string]App  `json:"apps"`
}

type User struct {
	Name string    `json:"name"`
	Age  StringInt `json:"age"`
}

func (u User) GetName() string {
	return u.Name
}

func (u User) GetType() string {
	return "User"
}

type App struct {
	Name    string      `json:"name"`
	Created StringInt   `json:"created"`
	Price   StringFloat `json:"price"`
}

func (u App) GetName() string {
	return u.Name
}

func (u App) GetType() string {
	return "App"
}

func (st *StringInt) UnmarshalJSON(b []byte) error {
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringInt(v)
	case float64:
		*st = StringInt(int(v))
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		*st = StringInt(i)

	}
	return nil
}

func (st *StringFloat) UnmarshalJSON(b []byte) error {
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringFloat(v)
	case float64:
		*st = StringFloat(float64(v))
	case string:
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		*st = StringFloat(i)

	}
	return nil
}

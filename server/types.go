package server

import (
	"encoding/json"
	"strconv"
)

type StringInt int

type User struct {
	Name string    `json:"name"`
	Age  StringInt `json:"age"`
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

package redis

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strconv"
)

type Bulk []byte

type MultiBulk []Bulk

func (b Bulk) S() string {
	return string(b)
}

func (b Bulk) I() (i int) {
	i, _ = strconv.Atoi(string(b))
	return
}

func (b Bulk) B() (bl bool) {
	bl, _ = strconv.ParseBool(string(b))
	return
}

func (b Bulk) Json(obj interface{}) (err error) {
	return json.Unmarshal(b, obj)
}

func (b Bulk) Gob(obj interface{}) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(b))
	err = dec.Decode(obj)
	return
}

package redis

import (
	"testing"
)

func TestPing(t *testing.T) {
	s, err := RD.Ping()
	if err != nil {
		t.Error(err)
	}
	if s != "PONG" {
		t.Errorf("Not PONG!! %s", s)
	}
}

func TestGetUnExists(t *testing.T) {
	r, err := RD.Get("unexiststestkey")
	if err != nil {
		t.Error(err)
	}
	if r != nil {
		t.Error("Must be nil")
	}
}

func TestSetGet(t *testing.T) {
	err := RD.Set("testkey", []byte("testvalue"))
	if err != nil {
		t.Error(err)
	}
	r, err := RD.Get("testkey")
	if err != nil {
		t.Error(err)
	}
	if r.S() != "testvalue" {
		t.Errorf("Unexpected value %s", string(r))
	}
}


func TestAppend(t *testing.T) {
	RD.Set("testkey", []byte("testvalue"))
	l, err := RD.Append("testkey", []byte("_append"))
	if err != nil {
		t.Error(err)
	}
	if l != len("testvalue_append") {
		t.Error("Unexpected len")
	}
	r, err := RD.Get("testkey")
	if err != nil {
		t.Error(err)
	}
	if r.S() != "testvalue_append" {
		t.Errorf("Unexpected value %s", string(r))
	}
}

func TestIncrDescr(t *testing.T) {
	RD.Set("testincr", []byte("10"))
	incr, err := RD.Incr("testincr")
	if err != nil {
		t.Error(err)
	}
	if incr != 11 {
		t.Error("Unexpected increment ouput")
	}
	incr, err = RD.IncrBy("testincr", 4)
	if err != nil {
		t.Error(err)
	}
	if incr != 15 {
		t.Error("Unexpected increment ouput")
	}
	decr, err := RD.Decr("testincr")
	if err != nil {
		t.Error(err)
	}
	if decr != 14 {
		t.Error("Unexpected increment ouput")
	}
	decr, err = RD.DecrBy("testincr", 4)
	if err != nil {
		t.Error(err)
	}
	if decr != 10 {
		t.Error("Unexpected increment ouput")
	}
}

func TestGetSet(t *testing.T) {
	RD.Set("testgetset", []byte("old"))
	b, err := RD.GetSet("testgetset", []byte("new"))
	if err != nil {
		t.Error(err)
	}
	if b.S() != "old" {
		t.Error("Unexpected getset output")
	}
	b, err = RD.Get("testgetset")
	if err != nil {
		t.Error(err)
	}
	if b.S() != "new" {
		t.Error("Unexpected getset output")
	}
}

func TestGetRange(t *testing.T) {
	RD.Set("testgetrange", []byte("hello world"))
	b, err := RD.GetRange("testgetrange", 0, 4)
	if err != nil {
		t.Error(err)
	}
	if b.S() != "hello" {
		t.Error("Unexpected getrange output", b.S())
	}
}
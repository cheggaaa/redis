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
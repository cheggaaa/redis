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
	defer RD.Del("unexiststestkey")
	if err != nil {
		t.Error(err)
	}
	if r != nil {
		t.Error("Must be nil")
	}
}

func TestSetGet(t *testing.T) {
	err := RD.Set("testkey", []byte("testvalue"))
	defer RD.Del("testkey")
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
	defer RD.Del("testkey")
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
	defer RD.Del("testincr")
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
	defer RD.Del("testgetset")
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
	defer RD.Del("testgetrange")
	b, err := RD.GetRange("testgetrange", 0, 4)
	if err != nil {
		t.Error(err)
	}
	if b.S() != "hello" {
		t.Error("Unexpected getrange output", b.S())
	}
}

func TestMGetMSet(t *testing.T) {
	values := map[string][]byte{
		"testmset1": []byte("v1"),
		"testmset2": []byte("v2"),
		"testmset3": []byte("v3"),
		"testmset4": []byte("v4"),
	}

	keys := make([]string, 0)
	for k := range values {
		keys = append(keys, k)
	}

	defer RD.Del(keys...)

	err := RD.MSet(values)
	if err != nil {
		t.Error(err)
	}

	b, err := RD.Get("testmset2")
	if err != nil {
		t.Error(err)
	}
	if b.S() != "v2" {
		t.Error("Unsexpected value: ", b.S())
	}

	r, err := RD.MGet(keys...)
	if err != nil {
		t.Error(err)
	}

	if len(r) != len(values) {
		t.Error("Unexpected MGet length")
	}
	for k, v := range r {
		if string(values[k]) != v.S() {
			t.Error("Not equals!")
		}
	}
}

func TestBitOps(t *testing.T) {
	// SETBIT
	defer RD.Del("testsetbit")
	r, err := RD.SetBit("testsetbit", 7, 1)
	if err != nil {
		t.Error(err)
	}
	if r != 0 {
		t.Errorf("Unexpected SETBIT result: %d", r)
	}
	if s, _ := RD.Get("testsetbit"); s.S() != "\x01" {
		t.Errorf("Unexpected SETBIT result: %s", s.S())
	}

	// GETBIT
	r, err = RD.GetBit("testsetbit", 0)
	if err != nil {
		t.Error(err)
	}
	if r != 0 {
		t.Errorf("Unexpected GETBIT result: %d", r)
	}
	r, err = RD.GetBit("testsetbit", 7)
	if err != nil {
		t.Error(err)
	}
	if r != 1 {
		t.Errorf("Unexpected GETBIT result: %d", r)
	}

	// BITCOUNT
	RD.Set("testbitcount", []byte("ololo"))
	defer RD.Del("testbitcount")

	r, err = RD.BitCount("testbitcount", nil, nil)
	if err != nil {
		t.Error(err)
	}
	if r != 26 {
		t.Errorf("Unexpected result BITCOUNT: %d", r)
	}

	if r, _ = RD.BitCount("testbitcount", &NilInt{0}, &NilInt{0}); r != 6 {
		t.Errorf("Unexpected result BITCOUNT: %d", r)
	}
	if r, _ = RD.BitCount("testbitcount", &NilInt{0}, &NilInt{1}); r != 10 {
		t.Errorf("Unexpected result BITCOUNT: %d", r)
	}

	// BITOP
	RD.MSet(map[string][]byte{
		"testbitop1": []byte("olo1"),
		"testbitop2": []byte("olo2"),
	})
	defer RD.Del("testbitop1", "testbitop2", "testbitopand")
	r, err = RD.BitOp(BITOP_AND, "testbitopand", "testbitop1", "testbitop2")
	if err != nil {
		t.Error(err)
	}
	if r != 4 {
		t.Errorf("Unexpected result BITOP: %d", r)
	}
	if rs, _ := RD.Get("testbitopand"); rs.S() != "olo0" {
		t.Errorf("Unexpected result BITOP AND: %s", rs.S())
	}

}

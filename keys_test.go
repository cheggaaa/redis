package redis

import (
	"testing"
	"time"
)

var RD *Redis

func TestConnect(t *testing.T) {
	c, err := NewClient("")
	if err != nil {
		t.Error(err)
	}
	RD = c
}

func TestDel(t *testing.T) {
	RD.Set("testkey", []byte("testvalue"))
	RD.Set("testkey2", []byte("testvalue"))
	RD.Set("testkey3", []byte("testvalue"))
	c, err := RD.Del("testkey")
	if err != nil {
		t.Error(err)
	}
	if c != 1 {
		t.Errorf("Must be 1, expected %d", c)
	}

	c, err = RD.Del("testkey2", "testkey3")
	if err != nil {
		t.Error(err)
	}
	if c != 2 {
		t.Errorf("Must be 2, expected %d", c)
	}
}

func TestExists(t *testing.T) {
	RD.Set("testkey", []byte("testvalue"))
	is, err := RD.Exists("testkey")
	if err != nil {
		t.Error(err)
	}
	if !is {
		t.Error("Key must be exists")
	}
	is, err = RD.Exists("testkey2")
	if err != nil {
		t.Error(err)
	}
	if is {
		t.Error("Key must not be exists")
	}
}

func TestExpireAndTTL(t *testing.T) {
	RD.Set("testkey", []byte("testvalue"))
	ok, err := RD.Expire("testkey", 2)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Must be true")
	}

	ttl, err := RD.TTL("testkey")
	if err != nil {
		t.Error(err)
	}
	if ttl != 2 {
		t.Errorf("TTL must be 2, expected %d", ttl)
	}
	time.Sleep(time.Second)
	ttl, err = RD.TTL("testkey")
	if err != nil {
		t.Error(err)
	}
	if ttl != 1 {
		t.Errorf("TTL must be 1, expected %d", ttl)
	}
	time.Sleep(time.Second + time.Millisecond*100)
	ttl, err = RD.TTL("testkey")
	if err != nil {
		t.Error(err)
	}
	if ttl != -1 {
		t.Errorf("TTL must be -1, expected %d", ttl)
	}
}

func TestPExpireAndPTTL(t *testing.T) {
	RD.Set("testkey", []byte("testvalue"))
	ok, err := RD.PExpire("testkey", 100)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Must be true")
	}

	ttl, err := RD.PTTL("testkey")
	if err != nil {
		t.Error(err)
	}
	if ttl == 0 || ttl > 100 {
		t.Errorf("TTL must be 2, expected %d", ttl)
	}
	time.Sleep(time.Millisecond * 102)
	ttl, err = RD.PTTL("testkey")
	if err != nil {
		t.Error(err)
	}
	if ttl != -1 {
		t.Errorf("TTL must be -1, expected %d", ttl)
	}
}

func TestKeys(t *testing.T) {
	sets := map[string][]byte{
		"__testkeys1": []byte("tv"),
		"__testkeys2": []byte("tv"),
		"__testkeys3": []byte("tv"),
		"__testkeys4": []byte("tv"),
	}

	for k, v := range sets {
		RD.Set(k, v)
	}

	mb, err := RD.Keys("__testkeys*")
	if err != nil {
		t.Error(err)
	}
	if len(mb) != 4 {
		t.Errorf("Must be 4, expect %d", len(mb))
	}

	for _, key := range mb {
		_, ok := sets[key.S()]
		if !ok {
			t.Error("key must be exists")
		}
	}

	for k, _ := range sets {
		RD.Del(k)
	}

}

func TestDumpRestore(t *testing.T) {
	defer RD.Del("__testdump", "__testdumprestore")
	// dump
	RD.Set("__testdump", []byte("dump-restore-data"))
	dump, err := RD.Dump("__testdump")
	if err != nil {
		t.Error(err)
	}
	// restore
	err = RD.Restore("__testdumprestore", 2, dump)
	if err != nil {
		t.Error(err)
	}

	// check data
	val, _ := RD.Get("__testdumprestore")
	if val.S() != "dump-restore-data" {
		t.Error("Restore: failed data")
	}
}

func TestMove(t *testing.T) {
	RD.Set("__testmove", []byte("MOVEIT"))
	RD.Expire("__testmove", 1)
	moved, err := RD.Move("__testmove", 3)
	if err != nil {
		t.Error(err)
	}
	if !moved {
		t.Error("key was not moved")
	}
	RD.Select(3)
	data, _ := RD.Get("__testmove")
	if data.S() != "MOVEIT" {
		t.Error("key was not moved")
	}
	RD.Del("__testmove")
	RD.Select(0)
}

func TestPersist(t *testing.T) {
	defer RD.Del("__testpersist")
	RD.Set("__testpersist", []byte("LOLO"))
	RD.Expire("__testpersist", 2)
	timeoutRemoved, err := RD.Persist("__testpersist")
	if err != nil {
		t.Error(err)
	}
	if !timeoutRemoved {
		t.Error("timeout was not moved")
	}
}

func TestType(t *testing.T) {
	defer RD.Del("__testtype")
	RD.Set("__testtype", []byte("LOLO"))
	tp, err := RD.Type("__testtype")
	if err != nil {
		t.Error(err)
	}
	if tp != TYPE_STRING {
		t.Error("Type not STRING", tp)
	}
}

func TestRenameRenameNx(t *testing.T) {
	defer RD.Del("__testrenamed", "__testrenamednx")
	RD.Set("__testrename", []byte("LOLO"))
	RD.Set("__testrenamenx", []byte("LOLO"))
	err := RD.Rename("__testrename", "__testrenamed")
	if err != nil {
		t.Error(err)
	}
	if ex, _ := RD.Exists("__testrenamed"); !ex {
		t.Error("was not renamed")
	}
	ok, err := RD.RenameNx("__testrenamenx", "__testrenamednx")
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("was not renamed nx")
	}
	if ex, _ := RD.Exists("__testrenamednx"); !ex {
		t.Error("was not renamed nx")
	}

	ok, err = RD.RenameNx("__testrenamednx", "__testrenamed")
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Error("was renamed nx")
	}
}

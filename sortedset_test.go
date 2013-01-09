package redis

import (
	"testing"
)

var testZSet = map[int]string{
		1 : "one",
		2 : "two",
		3 : "three",
		4 : "four",
		5 : "five",
	}

func TestZAdd(t *testing.T) {
	RD.Del("testzset")
	defer RD.Del("testzset")
	
	i, err := RD.ZAdd("testzset", 1, "one")
	if err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Errorf("Must be 1, but %d expected", i)
	}
}

func TestZCard(t *testing.T) {
	defer RD.Del("testzset")
	for score, member := range testZSet {
		RD.ZAdd("testzset", score, member)
	}
	
	i, err := RD.ZCard("testzset")
	if err != nil {
		t.Error(err)
	}
	if i != 5 {
		t.Errorf("Must be 5, but %d expected", i)
	}
}

func TestZCount(t *testing.T) {
	defer RD.Del("testzset")
	for score, member := range testZSet {
		RD.ZAdd("testzset", score, member)
	}
		
	i, err := RD.ZCount("testzset", 0, 10)
	if err != nil {
		t.Error(err)
	}
	if i != 5 {
		t.Errorf("Must be 5, but %d expected", i)
	}
	
	i, err = RD.ZCount("testzset", 2, 3)
	if err != nil {
		t.Error(err)
	}
	if i != 2 {
		t.Errorf("Must be 2, but %d expected", i)
	}
	
	i, err = RD.ZCount("testzset", 0, 1)
	if err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Errorf("Must be 1, but %d expected", i)
	}
}


func TestZIncrBy(t *testing.T) {
	defer RD.Del("testzset")
	RD.ZAdd("testzset", 1, "one")
	
	i, err := RD.ZIncrBy("testzset", 1, "one")
	if err != nil {
		t.Error(err)
	}
	if i != 2 {
		t.Errorf("Must be 2, but %d expected", i)
	}
	
	i, err = RD.ZIncrBy("testzset", 10, "one")
	if err != nil {
		t.Error(err)
	}
	if i != 12 {
		t.Errorf("Must be 12, but %d expected", i)
	}
}


func TestZRange(t *testing.T) {
	
}

func TestZRangeByScoreLimit(t *testing.T) {

}

func TestZRangeByScore(t *testing.T) {

}

func TestZRank(t *testing.T) {

}


func TestZRem(t *testing.T) {

}


func TestZRemRangeByRank(t *testing.T) {

}

func TestZRemRangeByScore(t *testing.T) {

}

func TestZRevRange(t *testing.T) {

}


func TestZRevRangeByScoreLimit(t *testing.T) {

}

func TestZRevRangeByScore(t *testing.T) {

}


func TestZRevRank(t *testing.T) {

}

func TestZScore(t *testing.T) {

}

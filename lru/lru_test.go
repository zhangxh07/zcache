package lru_test

import (
	"fmt"
	"reflect"
	"testing"
	"zcache/lru"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestCache_Get(t *testing.T) {
	l := lru.New(int64(0), nil)
	l.Add("key1", String("1234"))
	if v, ok := l.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := l.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	} else {
		t.Log("key2 not exist")
	}
}

// 测试，当使用内存超过了设定值时，是否会触发“无用”节点的移除：
func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	l := lru.New(int64(cap), nil)
	l.Add(k1, String(v1))
	l.Add(k2, String(v2))
	l.Add(k3, String(v3))

	if _, ok := l.Get("key1"); ok || l.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	} else {
		t.Log("key1 is remove")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value lru.Value) {
		keys = append(keys, key) //存储的是因为超过最大内存而删除的key
	}
	l := lru.New(int64(10), callback)
	l.Add("key1", String("123456"))
	l.Add("k2", String("k2"))
	l.Add("k3", String("k3"))
	l.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	fmt.Println(keys)

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}

}

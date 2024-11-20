package chash

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int            //虚拟节点倍数
	keys     []int          //哈希环切片
	hashMap  map[int]string //虚拟节点与真实节点的映射表，键是虚拟节点的哈希值，值是真实节点的名称。
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	//如果函数为空，允许自定义hash函数，默认为 crc32.ChecksumIEEE 算法
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.    //把节点地址转换成虚拟节点hash只，并加入到hash(keys)环上,并在map中映射hash与节点地址(key)关系
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			fmt.Println(key, hash)
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
	fmt.Println(m.keys)
}

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))                   //计算key的hash值
	idx := sort.Search(len(m.keys), func(i int) bool { //search函数返回m.keys[i] >= hash中切片keys的最小索引
		return m.keys[i] >= hash
	})
	fmt.Println("hash:", hash, "idx:", idx)

	realnode := m.hashMap[m.keys[idx%len(m.keys)]]
	fmt.Println("realnode：", realnode)

	return realnode
}

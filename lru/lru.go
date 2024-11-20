package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes int64      //允许最大使用内存
	nbytes   int64      //当前已使用内存
	ls       *list.List //双向链表
	cahce    map[string]*list.Element

	//记录被移除key时的回调函数
	OnEvicted func(key string, value Value)
}

func New(maxBytes int64, one func(str string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ls:        list.New(),
		cahce:     make(map[string]*list.Element),
		OnEvicted: one,
	}
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// Get look ups a key's value
func (c Cache) Get(key string) (value Value, ok bool) {
	//判断健是否存在，如果存在则移到队尾
	if ele, ok := c.cahce[key]; ok {
		c.ls.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除
func (c *Cache) RemoveOldest() {
	//取队首节点
	ele := c.ls.Back()
	if ele != nil {
		//在链表里删除
		c.ls.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cahce, kv.key)                                //在字典里删除
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) //在当前使用内存中减去删除缓存的大小
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 新增
func (c *Cache) Add(key string, value Value) {
	//如果健存在就移到对尾，并更新value
	if ele, ok := c.cahce[key]; ok {
		c.ls.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//不存在则在队尾添加新节点
		ele := c.ls.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cahce[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len()) //更新当前已使用内存
		//fmt.Println("nbytes:", c.nbytes)
	}
	//如果超过当前可使用内存，则触发删除节点
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ls.Len()
}

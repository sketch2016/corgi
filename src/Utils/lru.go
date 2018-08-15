package Utils

import (
	"sync"
)

var lruTag = "LRU"

//LruKey lru key
type LruKey interface{}

//LruVal lru value
type LruVal interface{}

//LRUNode lru node
type lruNode struct {
	next *lruNode
	prev *lruNode
	key  interface{}
	val  interface{}
}

//LRUCache lru cache
type LRUCache struct {
	LRUMap  map[LruKey]*lruNode
	mapLock sync.RWMutex

	cap      int
	len      int
	head     *lruNode
	current  *lruNode
	listLock sync.RWMutex
}

//CreateLRUCache create LRU cache
func CreateLRUCache(cap int) *LRUCache {
	v := new(LRUCache)
	v.LRUMap = make(map[LruKey]*lruNode)
	v.len = 0
	v.cap = cap
	return v
}

//GetLength get length
func (p *LRUCache) GetLength() int {
	return p.len
}

//GetCap get cap
func (p *LRUCache) GetCap() int {
	return p.cap
}

//GetCache get cache
func (p *LRUCache) GetCache(key interface{}) (val interface{}) {
	p.mapLock.RLock()
	node, ok := p.LRUMap[key]
	p.mapLock.RUnlock()

	if ok {
		//LOGD(lruTag, "cache hit")
		///update list
		p.listLock.Lock()

		parent := node.prev
		child := node.next
		if parent != nil {
			parent.next = child
		}

		if child != nil {
			child.prev = parent
		}

		node.next = p.head
		p.head.prev = node
		p.head = node
		p.listLock.Unlock()

		return node.val
	}

	return nil
}

//AddCache add cache
func (p *LRUCache) AddCache(key interface{}, val interface{}) {
	//we should check wheter cache exist
	//fmt.Printf("AddCache val is %p \n", val)
	LOGD(lruTag, "add cache start")
	//time.Sleep(time.Duration(200) * time.Millisecond)
	p.mapLock.RLock()
	node, ok := p.LRUMap[key]
	p.mapLock.RUnlock()

	var newNode *lruNode
	var isNeedRemove = false

	if ok {
		newNode = node
	} else {
		newNode = new(lruNode)
		if p.len == p.cap {
			isNeedRemove = true
		} else {
			p.len++
		}
	}

	//cache already exits
	p.mapLock.Lock()
	newNode.key = key
	newNode.val = val
	p.LRUMap[key] = newNode

	if isNeedRemove {
		LOGD(lruTag, "delete node is ", p.current)
		delete(p.LRUMap, p.current.key)
	}

	p.mapLock.Unlock()

	p.listLock.Lock()
	defer func() {
		p.listLock.Unlock()
	}()

	if !isNeedRemove {
		if p.head == nil {
			p.head = newNode
			p.current = p.head
			return
		}
	} else {
		//remove current key
		LOGD(lruTag, "remove node is ", p.current)
		parent := p.current.prev
		parent.next = nil
		p.current = parent
	}

	//add new key to head
	newNode.next = p.head
	p.head.prev = newNode
	p.head = newNode

}

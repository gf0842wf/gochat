package zmap

import (
	"sync"
)

type SafeMap struct {
	lock *sync.RWMutex
	sm   map[interface{}]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		sm:   make(map[interface{}]interface{}),
	}
}

//Get from maps return the k's value
func (m *SafeMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.sm[k]; ok {
		return val
	}
	return nil
}

// Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *SafeMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.sm[k]; !ok {
		m.sm[k] = v
	} else if val != v {
		m.sm[k] = v
	} else {
		return false
	}
	return true
}

// Returns true if k is exist in the map.
func (m *SafeMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.sm[k]; !ok {
		return false
	}
	return true
}

func (m *SafeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.sm, k)
}

// 以下几个最好别加锁啊,太耗时了

// Returns the keys of map
// 协程安全版range
/*
for _, v := range sp.Keys() {

}
*/
func (m *SafeMap) Keys() []interface{} {
	m.lock.RLock()
	keys := make([]interface{}, len(m.sm))
	i := 0
	for _, v := range m.sm {
		keys[i] = v
		i++
	}
	return keys
}

// Returns the values of map
// 协程安全版range
/*
for _, v := range sp.Vaules() {

}
*/
func (m *SafeMap) Values() []interface{} {
	m.lock.RLock()
	values := make([]interface{}, len(m.sm))
	i := 0
	for _, v := range m.sm {
		values[i] = v
		i++
	}
	return values
}

// Return the keys and values of map
// 协程安全版range
/*
 */
func (m *SafeMap) Items() []interface{} { // [[k, v], ]
	m.lock.RLock()
	items := make([]interface{}, 0, len(m.sm))
	i := 0
	for k, v := range m.sm {
		items = append(items, []interface{}{k, v})
		i++
	}
	return items
}

// 用于range等操作, 非协程安全
func (m *SafeMap) SM() map[interface{}]interface{} {
	return m.sm
}

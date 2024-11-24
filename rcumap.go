package ut

import (
	"sync/atomic"
	"unsafe"
)

type RCUMap[K comparable, V any] struct {
	head unsafe.Pointer // 指向当前版本的数据
}

// Node 表示 RCUMap 中的一个节点
type Node[K comparable, V any] struct {
	data map[K]V
}

// NewRCUMap 创建一个新的 RCUMap
func NewRCUMap[K comparable, V any]() *RCUMap[K, V] {
	initialNode := &Node[K, V]{data: make(map[K]V)}
	return &RCUMap[K, V]{head: unsafe.Pointer(initialNode)}
}

// Get 获取值
func (m *RCUMap[K, V]) Get(key K) (V, bool) {
	current := (*Node[K, V])(atomic.LoadPointer(&m.head))
	value, ok := current.data[key]
	return value, ok
}

// Range
func (m *RCUMap[K, V]) Range(f func(key K, value V) bool) {
	current := (*Node[K, V])(atomic.LoadPointer(&m.head))
	for k, v := range current.data {
		if !f(k, v) {
			return
		}
	}
}

// Set 设置值
func (m *RCUMap[K, V]) Set(key K, value V) {
	for {
		current := (*Node[K, V])(atomic.LoadPointer(&m.head))
		newNode := &Node[K, V]{data: make(map[K]V)}

		// 复制当前数据
		for k, v := range current.data {
			newNode.data[k] = v
		}
		newNode.data[key] = value

		if atomic.CompareAndSwapPointer(&m.head, unsafe.Pointer(current), unsafe.Pointer(newNode)) {
			return
		}
	}
}

// Remove 删除值
func (m *RCUMap[K, V]) Remove(key K) {
	for {
		current := (*Node[K, V])(atomic.LoadPointer(&m.head))
		newNode := &Node[K, V]{data: make(map[K]V)}

		// 复制当前数据
		for k, v := range current.data {
			newNode.data[k] = v
		}
		delete(newNode.data, key)

		// 尝试更新 head
		if atomic.CompareAndSwapPointer(&m.head, unsafe.Pointer(current), unsafe.Pointer(newNode)) {
			return
		}
	}
}

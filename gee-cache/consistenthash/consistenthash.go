package consistenthash

import (
	"fmt"
	"hash/crc32"
	"sort"
)

// 哈希函数类型，将字节数组映射为uint32
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // 哈希函数
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环上的所有节点
	hashMap  map[int]string // 虚拟节点与真实节点的映射表
}

// New 创建一个一致性哈希算法的Map实例
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加真实节点到哈希环上
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(fmt.Sprint(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// 对哈希环上的节点进行排序
	sort.Ints(m.keys)
}

// Get 根据传入的key，返回对应的真实节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// 二分查找，找到大于等于hash的最小节点的位置
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 如果idx等于len(m.keys)，说明hash比所有节点的哈希值都大，应该返回第一个节点
	if idx == len(m.keys) {
		idx = 0
	}
	return m.hashMap[m.keys[idx]]
}

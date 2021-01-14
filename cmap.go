package cmap

import (
	"sync"
)

const (
	PIECE_NUMBER = 32
)
type CampSlice []*cmap
type cmap struct {
	// 内置sync.Mutex 锁
	mu   sync.Mutex
	data map[string]interface{}
}

// 返回计算出来的哈希值跟我们的分片的取余。
func (cs CampSlice)getid(key string) *cmap {
	return cs[FnvHash(key) % PIECE_NUMBER]
}

func (c CampSlice) Set(key string, value interface{}) {
	cp := c.getid(key)
	cp.mu.Lock()
	cp.data[key] = value
	cp.mu.Unlock()
}

func (c CampSlice) Get(key string) interface{} {
	cp := c.getid(key)
	cp.mu.Lock()
	result := cp.data[key]
	cp.mu.Unlock()
	return result
}

func NewCmap() CampSlice {
	s := make([]*cmap, PIECE_NUMBER)
	for i := 0; i < PIECE_NUMBER; i++ {
		s[i] = &cmap{
			data: make(map[string]interface{}),
		}
	}
	return s
}
func FnvHash(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
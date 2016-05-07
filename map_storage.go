package kademlia

import (
	"sync"
	"errors"
)

type MapStorage struct {
	sync.RWMutex
	m map[NodeID][]byte
}

func NewMapStorage() *MapStorage {
	storage := &MapStorage{
		m: make(map[NodeID][]byte, 0),
	}
	return storage
}

func (storage *MapStorage) Get(key NodeID) ([]byte, error) {
	storage.RLock()
	defer storage.RUnlock()
	if v, ok := storage.m[key]; ok {
		return v[:], nil
	}
	return nil, errors.New("Key not found")
}

func (storage *MapStorage) Put(key NodeID, value []byte) error {
	storage.Lock()
	defer storage.Unlock()
	storage.m[key] = value[:]
	return nil
}


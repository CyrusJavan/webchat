package chatservice

import (
	"fmt"
)

type MapStore struct {
	m map[string]string
}

var ErrNotFound error = fmt.Errorf("not found")

func (m MapStore) Get(k string) (string, error) {
	v, ok := m.m[k]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

func (m MapStore) Set(k string, v string) error {
	m.m[k] = v
	return nil
}

func NewMapStore() *MapStore {
	return &MapStore{
		m: make(map[string]string),
	}
}
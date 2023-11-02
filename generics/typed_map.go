package generics

import (
	"errors"
	"sync"
)

type MapItem interface{}

type TypedMap[T MapItem] struct {
	count int
	inner *sync.Map
}

var (
	ErrorNoSuchItem = errors.New("no such item")
	ErrorCastFail   = errors.New("failed to cast item")
)

func NewTypedMap[T MapItem]() *TypedMap[T] {
	return &TypedMap[T]{
		count: 0,
		inner: &sync.Map{},
	}
}

func (m *TypedMap[T]) Count() int {
	return m.count
}

func (m *TypedMap[T]) Add(key string, t T) {
	if _, exists := m.inner.LoadOrStore(key, t); !exists {
		m.count++
	}
}

func (m *TypedMap[T]) AddPtr(key string, t *T) {
	if _, exists := m.inner.LoadOrStore(key, t); !exists {
		m.count++
	}
}

func (m *TypedMap[T]) ItemOrDefault(key string, defaultT T) (T, error) {
	i, ok := m.inner.Load(key)
	if ok {
		t, ok := i.(T)
		if !ok {
			return defaultT, ErrorCastFail
		}
		return t, nil
	} else {
		return defaultT, ErrorNoSuchItem
	}
}

func (m *TypedMap[T]) ItemPtr(key string) (*T, error) {
	i, ok := m.inner.Load(key)
	if ok {
		t, ok := i.(*T)
		if !ok {
			return nil, ErrorCastFail
		}
		return t, nil
	} else {
		return nil, ErrorNoSuchItem
	}
}

func (m *TypedMap[T]) Range(f func(T) error) error {
	var err error = nil
	m.inner.Range(func(k interface{}, v interface{}) bool {
		t, ok := v.(T)
		if !ok {
			err = ErrorCastFail
			return false
		}

		err = f(t)
		return err == nil
	})
	return err
}

func (m *TypedMap[T]) RangePtr(f func(*T) error) error {
	var err error = nil
	m.inner.Range(func(k interface{}, v interface{}) bool {
		t, ok := v.(*T)
		if !ok {
			err = ErrorCastFail
			return false
		}

		err = f(t)
		return err == nil
	})
	return err
}

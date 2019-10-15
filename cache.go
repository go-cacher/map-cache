package map_cache

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-cacher/cacher"
)

/*MapCache MapCache */
type MapCache struct {
	cache sync.Map
}

type cacheData struct {
	value interface{}
	life  *time.Time
}

var NotFoundError = errors.New("data not found")
var WrongTypeError = errors.New("cast wrong type error")

func init() {
	cacher.Register(&MapCache{})
}

/*Get the exist value */
func (m MapCache) Get(key string) ([]byte, error) {
	if v, b := m.cache.Load(key); b {
		switch vv := v.(type) {
		case *cacheData:
			if vv.life != nil && vv.life.Before(time.Now()) {
				m.cache.Delete(key)
				return nil, NotFoundError
			}
			if v, b := vv.value.([]byte); b {
				return v, nil
			}
			return nil, WrongTypeError
		}
	}
	return nil, NotFoundError
}

//get a value with a default value
func (m MapCache) GetD(key string, v []byte) []byte {
	if ret, e := m.Get(key); e == nil {
		return ret
	}
	return v
}

// set a value to cache
func (m *MapCache) Set(key string, val []byte) error {
	m.cache.Store(key, &cacheData{
		value: val,
		life:  nil,
	})
	return nil
}

func (m *MapCache) SetWithTTL(key string, val []byte, ttl int64) error {
	t := time.Now().Add(time.Duration(ttl))
	m.cache.Store(key, &cacheData{
		value: val,
		life:  &t,
	})
	return nil
}

func (m MapCache) Has(key string) (bool, error) {
	if _, e := m.Get(key); e != nil {
		if errors.Is(e, NotFoundError) {
			return false, nil
		}
		return false, e
	}
	return true, nil
}

func (m *MapCache) Delete(key string) error {
	m.cache.Delete(key)
	return nil
}

func (m *MapCache) Clear() error {
	m.cache = sync.Map{}
	return nil
}

/*GetMultiple get multiple values */
func (m MapCache) GetMultiple(keys ...string) (map[string][]byte, error) {
	vals := make(map[string][]byte, len(keys))
	for _, key := range keys {
		if ret, e := m.Get(key); e == nil {
			vals[key] = ret
		}
		return nil, fmt.Errorf("%s:%+v", key, NotFoundError)
	}
	return vals, nil
}

/*SetMultiple set multiple values */
func (m *MapCache) SetMultiple(values map[string][]byte) error {
	for k, v := range values {
		e := m.Set(k, v)
		if e != nil {
			return fmt.Errorf("%s:%+v", k, e)
		}
	}
	return nil
}

/*DeleteMultiple delete multiple values */
func (m *MapCache) DeleteMultiple(keys ...string) error {
	for _, key := range keys {
		e := m.Delete(key)
		if e != nil {
			return fmt.Errorf("%s:%+v", key, e)
		}
	}
	return nil
}

// NewMapCache ...
func NewMapCache() *MapCache {
	return &MapCache{}
}

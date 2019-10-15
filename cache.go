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

func init() {
	cacher.Register(&MapCache{})
}

/*Get the exist value */
func (m MapCache) Get(key string) (interface{}, error) {
	if v, b := m.cache.Load(key); b {
		switch vv := v.(type) {
		case *cacheData:
			if vv.life != nil && vv.life.Before(time.Now()) {
				m.cache.Delete(key)
				return nil, NotFoundError
			}
			return vv.value, nil
		}
	}
	return nil, NotFoundError
}

//get a value with a default value
func (m MapCache) GetD(key string, v interface{}) interface{} {
	if ret, e := m.Get(key); e == nil {
		return ret
	}
	return v
}

// set a value to cache
func (m *MapCache) Set(key string, val interface{}) error {
	m.cache.Store(key, &cacheData{
		value: val,
		life:  nil,
	})
	return nil
}

func (m *MapCache) SetWithTTL(key string, val interface{}, ttl int64) error {
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
func (m MapCache) GetMultiple(keys ...string) (map[string]interface{}, error) {
	vals := make(map[string]interface{}, len(keys))
	for _, key := range keys {
		if ret, e := m.Get(key); e == nil {
			vals[key] = ret
		}
		return nil, fmt.Errorf("%s:%+v", key, NotFoundError)
	}
	return vals, nil
}

/*SetMultiple set multiple values */
func (m *MapCache) SetMultiple(values map[string]interface{}) error {
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

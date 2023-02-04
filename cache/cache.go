package cache

import (
	"context"
	"encoding/json"
	"errors"
	"kickcore/logging"
	"runtime"
	"time"
)

type CacheDriver interface {
	// Initialize cache
	Init() error

	// Pings cache connection
	PingContext(ctx context.Context) error

	// Inserts key-value to cache
	// Returns false if key is currently in cache
	Insert(key string, value []byte, date int64) (bool, error)

	// Selects the value specified by key in cache.
	Select(key string, value *[]byte) error

	// Selects the keys which are expired.
	SelectExpiredValues(expireAfter int64) ([]string, error)

	// Deletes the value specified by key in cache.
	Delete(key string) (bool, error)

	// Deletes the values which are specified by any key in cache.
	DeleteMany(keys []string) (int64, error)

	// Returns the length of cache
	Len() (int64, error)

	// Closes cache
	Close() error
}

type Cache struct {
	driver CacheDriver

	// It's true if you call c.Close()
	Closed bool
	Logger *logging.FileLogger
}

func NewCache(d CacheDriver, err error) (*Cache, error) {
	if err != nil {
		return nil, err
	}

	if err = d.Init(); err != nil {
		return nil, err
	}

	cache := Cache{driver: d}
	runtime.SetFinalizer(&cache, (*Cache).finalizer)
	return &cache, nil
}

func (c *Cache) finalizer() error {
	if c.Closed {
		return nil
	}

	c.Closed = true
	return c.driver.Close()
}

// Close the cache
func (c *Cache) Close() error { return c.finalizer() }

// Ping cache connection
func (c *Cache) Ping(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.driver.PingContext(ctx)
}

// Insert key-value into the cache.
func (c *Cache) Insert(apikey APICacheKey, key string, value []byte) (bool, error) {
	return c.driver.Insert(apikey.Key+key, value, time.Now().Unix()+apikey.ExtraTTL)
}

// Select value specified by the key.
func (c *Cache) Select(apikey APICacheKey, key string) ([]byte, error) {
	var value []byte
	err := c.driver.Select(apikey.Key+key, &value)
	return value, err
}

// Delete value specified by the key.
func (c *Cache) Delete(apikey APICacheKey, key string) (bool, error) {
	return c.driver.Delete(apikey.Key + key)
}

// Delete values which are specified by the keys.
func (c *Cache) DeleteMany(key []string) (int64, error) {
	return c.driver.DeleteMany(key)
}

// Cache length
func (c *Cache) Len() (int64, error) { return c.driver.Len() }

// CacheFunc first tries to returns value from cache, then if key not found in cache, call 'f'
// and (if it not returned error,) insert returned data into cache
//
// returns (data, data is in cache, error)
func (c *Cache) CacheFunc(apikey APICacheKey, key string, f func() ([]byte, error)) ([]byte, bool, error) {
	var value []byte = nil
	err := c.driver.Select(apikey.Key+key, &value)
	if value != nil {
		return value, true, err
	}

	value, err = f()
	if err != nil {
		return value, false, err
	}

	_, err = c.driver.Insert(apikey.Key+key, value, time.Now().Unix()+apikey.ExtraTTL)
	return value, err == nil, err
}

// Like c.CacheFunc, but recieve interface{} from 'f' and convert it to bytes by json.Marshal.
func (c *Cache) CacheFuncJSON(apikey APICacheKey, key string, f func() (interface{}, error)) ([]byte, bool, error) {
	var value []byte = nil
	err := c.driver.Select(apikey.Key+key, &value)
	if value != nil {
		return value, true, err
	}

	valueInterface, err := f()
	if err != nil {
		return nil, false, err
	}

	value, err = json.Marshal(valueInterface)
	if err != nil {
		return nil, false, err
	}

	_, err = c.driver.Insert(apikey.Key+key, value, time.Now().Unix()+apikey.ExtraTTL)
	return value, err == nil, err
}

// Expiration Machine - deletes values which are expired.
type ExpirationMachine struct {
	// Cache
	ActiveCache *Cache

	pool chan struct{}
}

func (e *ExpirationMachine) closeChannel() {
	if e.pool != nil {
		close(e.pool)
		e.pool = nil
	}
}

func (e *ExpirationMachine) deleteExpiredValues(expireAfterSeconds int64) (err error, exit bool) {
	if e.ActiveCache.Closed {
		return errors.New("cache is closed"), true
	}

	keys, err := e.ActiveCache.driver.SelectExpiredValues(expireAfterSeconds)
	if err != nil {
		return err, false
	}

	i, err := e.ActiveCache.driver.DeleteMany(keys)
	if err != nil {
		return err, false
	}

	if e.ActiveCache.Logger != nil {
		e.ActiveCache.Logger.Log(logging.LEVEL_DEBUG, "Number of expired values: %d", i)
	}
	return nil, false
}

func (e *ExpirationMachine) Start(interval time.Duration) error {
	if e.pool != nil {
		return errors.New("expiration machine already running")
	}

	err, _ := e.deleteExpiredValues(0)
	if err != nil {
		if e.ActiveCache.Logger != nil {
			e.ActiveCache.Logger.Log(logging.LEVEL_ERROR, "ExpirationMachine: while deleting: %s", err.Error())
		}
		return err
	}

	e.pool = make(chan struct{})

	go func() {
		if e.ActiveCache.Logger != nil {
			e.ActiveCache.Logger.Log(logging.LEVEL_DEBUG, "ExpirationMachine: running ...")
		}
		ticktack := time.NewTicker(interval)

		for {
			select {
			case <-ticktack.C:
				err, exit := e.deleteExpiredValues(0)
				if err != nil {
					if e.ActiveCache.Logger != nil {
						e.ActiveCache.Logger.Log(logging.LEVEL_ERROR, "ExpirationMachine: while deleting: %s", err.Error())
					}
				}
				if exit {
					e.closeChannel()
					return
				}
			case <-e.pool:
				return
			}
		}
	}()

	return nil
}

func (e *ExpirationMachine) Stop() error {
	if e.pool == nil {
		return errors.New("expiration machine isn't started")
	}

	e.ActiveCache.Logger.Log(logging.LEVEL_DEBUG, "ExpirationMachine: stopping ...")

	e.pool <- struct{}{}
	e.closeChannel()
	return nil
}

func (e *ExpirationMachine) IsStarted() bool { return e.pool != nil }

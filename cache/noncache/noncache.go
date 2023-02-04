package noncache

import (
	"context"
	"kickcore/cache"
)

type NonCache struct{}

func (c *NonCache) Init() error { return nil }

func (c NonCache) PingContext(_ context.Context) error { return nil }

func (c NonCache) Insert(_ string, _ []byte, _ int64) (bool, error) { return true, nil }

func (c NonCache) Select(_ string, _ *[]byte) error { return nil }

func (c NonCache) SelectExpiredValues(_ int64) ([]string, error) { return []string{}, nil }

func (c NonCache) Delete(_ string) (bool, error) { return true, nil }

func (c NonCache) DeleteMany(_ []string) (int64, error) { return 0, nil }

func (c NonCache) Len() (int64, error) { return 0, nil }

func (c NonCache) Close() error { return nil }

func Connect() (cache.CacheDriver, error) {
	c := new(NonCache)
	return c, nil
}

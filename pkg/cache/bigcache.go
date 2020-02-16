package cache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/tossp/tsgo/pkg/log"
)

var defConfig = bigcache.Config{
	// number of shards (must be a power of 2)
	Shards: 1024,

	// time after which entry can be evicted
	LifeWindow: 10 * time.Minute,

	// Interval between removing expired entries (clean up).
	// If set to <= 0 then no action is performed.
	// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
	CleanWindow: 5 * time.Minute,

	// rps * lifeWindow, used only in initial memory allocation
	MaxEntriesInWindow: 1000 * 10 * 60,

	// max entry size in bytes, used only in initial memory allocation
	MaxEntrySize: 500,

	// prints information about additional memory allocation
	Verbose: true,

	// cache will not allocate more memory than this limit, value in MB
	// if value is reached then the oldest entries can be overridden for the new ones
	// 0 value means no size limit
	HardMaxCacheSize: 8192,

	// callback fired when the oldest entry is removed because of its expiration time or no space left
	// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
	// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
	OnRemove: nil,

	// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
	// for the new entry, or because delete was called. A constant representing the reason will be passed through.
	// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
	// Ignored if OnRemove is specified.
	OnRemoveWithReason: nil,
}

func NewBigCahceConfig(maxCacheSize int, clean, eviction time.Duration) bigcache.Config {
	config := bigcache.DefaultConfig(eviction)
	config.CleanWindow = clean
	config.HardMaxCacheSize = maxCacheSize
	return config

}

func NewBigCache(config bigcache.Config) (cache *bigcache.BigCache, err error) {
	cache, err = bigcache.NewBigCache(config)
	if err != nil {
		log.Fatal("初始化缓存失败", err)
	}
	return
}

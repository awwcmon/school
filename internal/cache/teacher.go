package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"

	"school/internal/database"
	"school/internal/model"
)

const (
	// cache prefix key, must end with a colon
	teacherCachePrefixKey = "teacher:"
	// TeacherExpireTime expire time
	TeacherExpireTime = 5 * time.Minute
)

var _ TeacherCache = (*teacherCache)(nil)

// TeacherCache cache interface
type TeacherCache interface {
	Set(ctx context.Context, id string, data *model.Teacher, duration time.Duration) error
	Get(ctx context.Context, id string) (*model.Teacher, error)
	MultiGet(ctx context.Context, ids []string) (map[string]*model.Teacher, error)
	MultiSet(ctx context.Context, data []*model.Teacher, duration time.Duration) error
	Del(ctx context.Context, id string) error
	SetPlaceholder(ctx context.Context, id string) error
	IsPlaceholderErr(err error) bool
}

// teacherCache define a cache struct
type teacherCache struct {
	cache cache.Cache
}

// NewTeacherCache new a cache
func NewTeacherCache(cacheType *database.CacheType) TeacherCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Teacher{}
		})
		return &teacherCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Teacher{}
		})
		return &teacherCache{cache: c}
	}

	return nil // no cache
}

// GetTeacherCacheKey cache key
func (c *teacherCache) GetTeacherCacheKey(id string) string {
	return teacherCachePrefixKey + id
}

// Set write to cache
func (c *teacherCache) Set(ctx context.Context, id string, data *model.Teacher, duration time.Duration) error {
	if data == nil || id == "" {
		return nil
	}
	cacheKey := c.GetTeacherCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *teacherCache) Get(ctx context.Context, id string) (*model.Teacher, error) {
	var data *model.Teacher
	cacheKey := c.GetTeacherCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *teacherCache) MultiSet(ctx context.Context, data []*model.Teacher, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetTeacherCacheKey(v.ID.Hex())
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *teacherCache) MultiGet(ctx context.Context, ids []string) (map[string]*model.Teacher, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetTeacherCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Teacher)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*model.Teacher)
	for _, id := range ids {
		val, ok := itemMap[c.GetTeacherCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *teacherCache) Del(ctx context.Context, id string) error {
	cacheKey := c.GetTeacherCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *teacherCache) SetPlaceholder(ctx context.Context, id string) error {
	cacheKey := c.GetTeacherCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *teacherCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}

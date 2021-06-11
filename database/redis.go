package database

import (
	"context"
	"fmt"
	"time"

	"github.com/leopurba/go-article/config"

	cache "github.com/go-redis/cache/v8"
	redis "github.com/go-redis/redis/v8"
)

type RClient interface {
	Conn() *redis.Client
	Cache() *cache.Cache
	Close() error
}

func NewRClient() (RClient, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Cfg().RedisHost, config.Cfg().RedisPort),
		DB:       config.Cfg().RedisDatabase,
		Password: config.Cfg().RedisPassword,
		PoolSize: config.Cfg().RedisPoolSize,
	})

	err := db.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	dbcache := cache.New(&cache.Options{
		Redis:      db,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &client{db, dbcache}, nil
}

type client struct {
	db      *redis.Client
	dbcache *cache.Cache
}

func (c *client) Conn() *redis.Client {
	return c.db
}

func (c *client) Cache() *cache.Cache {
	return c.dbcache
}

func (c *client) Close() error {
	return c.db.Close()
}

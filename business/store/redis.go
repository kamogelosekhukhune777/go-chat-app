//
package store

import "github.com/go-redis/redis/v8"

// RedisCache is a wrapper struct around a Redis client.
// It lets us interact with Redis using higher-level methods (Set, Get, Delete).
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates and returns a new RedisCache instance.
//
// Parameters:
// - address: host:port of the Redis server (e.g., "localhost:6379")
// - password: password if Redis requires authentication (empty string if none)
// - db: database index (Redis supports multiple logical databases, usually 0 by default)
//
// It initializes the Redis client and returns a RedisCache object.
func NewRedisCache(address string, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     address,  // Redis server address
		Password: password, // password if set
		DB:       db,       // which logical DB to use
	})

	return &RedisCache{client: client}
}

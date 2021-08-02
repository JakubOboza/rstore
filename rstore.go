package rstore

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

//RedisStorage holds client info and implements fiber.Storage interface
type RedisStorage struct {
	addr        string
	db          int
	password    string
	redisClient *redis.Client
}

type RedisClientOption func(*RedisStorage)

//New returns fiber.Storage implementation only
func New(opts ...RedisClientOption) fiber.Storage {
	return NewClient(opts...)
}

//Newclient returns fiber.Storage plus extra methods
func NewClient(opts ...RedisClientOption) *RedisStorage {
	storage := &RedisStorage{}

	for _, option := range opts {
		option(storage)
	}

	if storage.addr == "" {
		storage.addr = "localhost:6379"
	}

	if storage.redisClient == nil {
		storage.redisClient = redis.NewClient(&redis.Options{
			Addr:     storage.addr,
			Password: storage.password,
			DB:       storage.db,
		})
	}

	return storage
}

//Client returns redis client pointer directly for extra setup?
func (rs *RedisStorage) Client() *redis.Client {
	return rs.redisClient
}

//ClientOptions

//Withclient sets entirely dev owned Client, all other options will be discarded if this is in use
func WithClient(client *redis.Client) RedisClientOption {
	return func(rs *RedisStorage) {
		rs.redisClient = client
	}
}

//WithAddr sets address of redis
func WithAddr(addr string) RedisClientOption {
	return func(rs *RedisStorage) {
		rs.addr = addr
	}
}

//WithDB sets DB number
func WithDB(DB int) RedisClientOption {
	return func(rs *RedisStorage) {
		rs.db = DB
	}
}

//WithPassword sets password for default redis client
func WithPassword(password string) RedisClientOption {
	return func(rs *RedisStorage) {
		rs.password = password
	}
}

// Get gets the value for the given key.
// It returns ErrNotFound if the storage does not contain the key.
func (rs *RedisStorage) Get(key string) ([]byte, error) {
	result := rs.redisClient.Get(context.Background(), key)
	val, err := result.Bytes()
	if redis.Nil == err {
		return val, fiber.ErrNotFound
	}
	return val, err
}

// Set stores the given value for the given key along with a
// time-to-live expiration value, 0 means live for ever
// Empty key or value will be ignored without an error.
func (rs *RedisStorage) Set(key string, val []byte, ttl time.Duration) error {
	result := rs.redisClient.Set(context.Background(), key, val, ttl)
	return result.Err()
}

// Delete deletes the value for the given key.
// It returns no error if the storage does not contain the key,
func (rs *RedisStorage) Delete(key string) error {
	result := rs.redisClient.Del(context.Background(), key)
	return result.Err()
}

// Reset resets the storage and delete all keys.
func (rs *RedisStorage) Reset() error {
	result := rs.redisClient.FlushAll(context.Background())
	return result.Err()
}

// Close closes the storage and will stop any running garbage
// collectors and open connections.
func (rs *RedisStorage) Close() error {
	return rs.redisClient.Close()
}

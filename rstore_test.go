package rstore

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func Test_RedisStore(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Errorf("You need to run redis on localhost:6379 for integration tests to work error: %v", err)
	}

	store := New(WithClient(client))

	key := "jakub"
	val := []byte("oboza")
	exp := 1 * time.Second

	store.Set(key, val, 0)
	store.Set(key, val, 0)

	result, err := store.Get(key)

	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, val, result)

	result, err = store.Get("empty")
	utils.AssertEqual(t, 0, len(result))
	utils.AssertEqual(t, fiber.ErrNotFound, err)

	err = store.Set(key, val, exp)
	utils.AssertEqual(t, nil, err)

	result, err = store.Get(key)
	utils.AssertEqual(t, val, result)
	utils.AssertEqual(t, nil, err)

	time.Sleep(1100 * time.Millisecond)

	result, err = store.Get(key)
	utils.AssertEqual(t, 0, len(result))
	utils.AssertEqual(t, fiber.ErrNotFound, err)

	store.Set(key, val, 0)
	result, err = store.Get(key)
	utils.AssertEqual(t, val, result)
	utils.AssertEqual(t, nil, err)

	store.Delete(key)
	result, err = store.Get(key)
	utils.AssertEqual(t, 0, len(result))
	utils.AssertEqual(t, fiber.ErrNotFound, err)

	store.Set("jakub", val, 0)
	store.Set("oboza", val, 0)
	err = store.Reset()
	utils.AssertEqual(t, nil, err)

	result, err = store.Get("jakub")
	utils.AssertEqual(t, 0, len(result))
	utils.AssertEqual(t, fiber.ErrNotFound, err)

	result, err = store.Get("oboza")
	utils.AssertEqual(t, 0, len(result))
	utils.AssertEqual(t, fiber.ErrNotFound, err)

}

func Test_RedisStoreClient(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Errorf("You need to run redis on localhost:6379 for integration tests to work error: %v", err)
	}

	//WithClient
	store := NewClient(WithClient(client))
	utils.AssertEqual(t, client, store.Client())

	//Close
	store.Close()
	_, err = client.Ping(context.Background()).Result()
	if err.Error() != "redis: client is closed" {
		t.Errorf("Client should be closed now but got %v", err)
	}

	//WithAddr
	//WithPassword
	//WithDB
	rscConfigs := NewClient(WithAddr("lol:123"), WithDB(13), WithPassword("Str@nkPazz"))

	utils.AssertEqual(t, "lol:123", rscConfigs.addr)
	utils.AssertEqual(t, 13, rscConfigs.db)
	utils.AssertEqual(t, "Str@nkPazz", rscConfigs.password)

}

func Benchmark_RedisStore(b *testing.B) {
	keyLength := 1000
	keys := make([]string, keyLength)
	for i := 0; i < keyLength; i++ {
		keys[i] = utils.UUID()
	}
	value := []byte("super good value for purpose of testing /s")

	ttl := 2 * time.Second
	b.Run("fiber_redis_store", func(b *testing.B) {
		d := New()
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for _, key := range keys {
				d.Set(key, value, ttl)
			}
			for _, key := range keys {
				_, _ = d.Get(key)
			}
			for _, key := range keys {
				d.Delete(key)

			}
		}
	})
}

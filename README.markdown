# Fiber Cache Redis Integration

This package lets You use redis as Your cache with [fiber.v2](https://github.com/gofiber/fiber) (golang web framework)

# Usage

When You are setting up cache import
```
  ...
	"github.com/JakubOboza/rstore"
	"github.com/gofiber/fiber/v2/middleware/cache"
```

Set the rstore as storage for config setup

```
cache.Config{
  Storage:      rstore.New(rstore.WithAddr("localhost:6379")),
}))

```

# Example

```
import (
	"time"

	"github.com/JakubOboza/rstore"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func main() {

	app := fiber.New()

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
		Storage:      rstore.New(rstore.WithAddr("localhost:6379")),
	}))

	app.Get("/cache-test", func(c *fiber.Ctx) error {
		return c.SendString(time.Now().String())
	})

	app.Listen(":3000")

}
```

Start redis on port 6379, start binary with example and
```
curl localhost:3000/cache-test
```

# Testing

To run integration tests please start redis on `localhost:6379` either via docker or on your own box. This is an integration package so it doesnt make much sense to mock redis.

You run tests by doing

```
make test
```

To run benchmarks just do:
```
make bench
```
Ofc because we dont use mock it will be heavily influenced by your redis performance. But at least You can compare to fiber.Memory and more less figure out/guess the impact.

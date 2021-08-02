# Fiber Cache Redis Integration

This package lets You use redis as Your cache with [fiber.v2](https://github.com/gofiber/fiber) (golang web framework)

# Usage

TODO: finish this section

# Example

TODO: full example of usage

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

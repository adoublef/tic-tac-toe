package structures_test

import (
	"testing"
	"tictactoe/pkg/structures"
	"time"

	"github.com/hyphengolang/prelude/testing/is"
)

// https://dev.to/bmf_san/implement-an-in-memory-cache-in-golang-lig
func TestCache(t *testing.T) {
	is := is.New(t)

	type entry struct {
		s string
		n int
	}

	// #0 -- initialise cache
	cache := structures.NewCache[string, *entry]()

	// #1 -- put value into cache
	cache.Put("foo", &entry{"foo", 1}, 3*time.Second)

	time.Sleep(2 * time.Second)

	// #2 -- get value from cache
	v, ok := cache.Get("foo")
	is.True(ok)      // true
	is.Equal(v.n, 1) // 1

	// #3 -- get value from cache
	_, ok = cache.Get("bar")
	is.True(!ok) // false

	// #4 -- delete value from cache
	cache.Delete("foo")

	// #5 -- get value from cache
	_, ok = cache.Get("foo")
	is.True(!ok) // false

	// #6 -- put value into cache, no expiry
	cache.Put("baz", &entry{"baz", 2}, 0)

	time.Sleep(2 * time.Second)

	// #7 -- get value from cache
	v, ok = cache.Get("baz")
	is.True(ok) // true
	is.Equal(v.n, 2)
}

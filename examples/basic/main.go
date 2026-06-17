package main

import (
    "fmt"

    "github.com/amitstephen-dev/goheapguard/pkg/heapguard"
)

type User struct {
    ID   int
    Name string
}

func main() {
    opts := heapguard.DefaultOptions()
    pool := heapguard.NewPool[User](opts)

    user := pool.Get()
    user.ID = 1
    user.Name = "Alice"

    fmt.Printf("User: %+v\n", user)

    pool.Put(user)

    user2 := pool.Get()
    fmt.Printf("User2: %+v\n", user2)

    pool.Put(user2)

    hits, misses, allocs, size := pool.Stats()
    fmt.Printf("Stats - Hits: %d, Misses: %d, Allocs: %d, Size: %d\n",
        hits, misses, allocs, size)
}
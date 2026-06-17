package heapguard

import "sync"

type Builder[T any] struct {
    pool  *Pool[T]
    obj   *T
    built bool
    mu    sync.Mutex
}

func NewBuilder[T any](pool *Pool[T]) *Builder[T] {
    return &Builder[T]{
        pool: pool,
    }
}

func (b *Builder[T]) Get() *T {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.built && b.obj != nil {
        b.pool.Put(b.obj)
    }

    b.obj = b.pool.Get()
    b.built = false
    return b.obj
}

func (b *Builder[T]) Build() *T {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.obj == nil {
        return nil
    }

    b.built = true
    obj := b.obj
    b.obj = nil
    return obj
}

func (b *Builder[T]) Reset() {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.obj != nil {
        b.pool.Put(b.obj)
        b.obj = nil
        b.built = false
    }
}
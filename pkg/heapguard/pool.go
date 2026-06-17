package heapguard

import (
    "runtime"
    "sync"
    "sync/atomic"
)

type Pool[T any] struct {
    pool     sync.Pool
    opts     Options
    resetFn  func(*T)
    hits     uint64
    misses   uint64
    allocs   uint64
    size     int32
    maxSize  int32
    mu       sync.Mutex
}

func NewPool[T any](opts Options) *Pool[T] {
    if opts.MaxSize <= 0 {
        opts.MaxSize = 1000
    }

    p := &Pool[T]{
        opts:    opts,
        maxSize: int32(opts.MaxSize),
    }

    p.pool.New = func() interface{} {
        atomic.AddUint64(&p.misses, 1)
        atomic.AddUint64(&p.allocs, 1)
        return new(T)
    }

    if opts.InitialSize > 0 {
        for i := 0; i < opts.InitialSize; i++ {
            obj := new(T)
            p.pool.Put(obj)
            atomic.AddInt32(&p.size, 1)
        }
    }

    if opts.EnableMetrics {
        runtime.SetFinalizer(p, func(p *Pool[T]) {
            if atomic.LoadInt32(&p.size) > 0 {
                println("WARNING: goheapguard pool finalized with",
                    atomic.LoadInt32(&p.size), "objects still in use")
            }
        })
    }

    return p
}

func (p *Pool[T]) Get() *T {
    obj := p.pool.Get()
    if obj != nil {
        atomic.AddUint64(&p.hits, 1)
        atomic.AddInt32(&p.size, -1)
        return obj.(*T)
    }
    atomic.AddUint64(&p.misses, 1)
    atomic.AddUint64(&p.allocs, 1)
    return new(T)
}

func (p *Pool[T]) Put(obj *T) {
    if obj == nil {
        return
    }

    if atomic.LoadInt32(&p.size) >= p.maxSize {
        return
    }

    p.reset(obj)
    atomic.AddInt32(&p.size, 1)
    p.pool.Put(obj)
}

func (p *Pool[T]) reset(obj *T) {
    if p.resetFn != nil {
        p.resetFn(obj)
        return
    }
    var zero T
    *obj = zero
}

func (p *Pool[T]) SetResetFunc(fn func(*T)) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.resetFn = fn
}

func (p *Pool[T]) Stats() (hits, misses, allocs uint64, size int32) {
    return atomic.LoadUint64(&p.hits),
        atomic.LoadUint64(&p.misses),
        atomic.LoadUint64(&p.allocs),
        atomic.LoadInt32(&p.size)
}

func (p *Pool[T]) Reset() {
    p.mu.Lock()
    defer p.mu.Unlock()
    for atomic.LoadInt32(&p.size) > 0 {
        obj := p.pool.Get()
        if obj != nil {
            atomic.AddInt32(&p.size, -1)
        }
    }
}
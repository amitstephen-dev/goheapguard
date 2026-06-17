package heapguard

import (
    "sync"
    "sync/atomic"
    "time"
)

type AutoPool[T any] struct {
    Pools     []*Pool[T]  // Exported now
    current   int32
    opts      Options
    mu        sync.RWMutex
    scaleOps  uint64
    lastScale time.Time
}

func NewAutoPool[T any](opts Options) *AutoPool[T] {
    if opts.MaxSize <= 0 {
        opts.MaxSize = 1000
    }

    ap := &AutoPool[T]{
        Pools:     make([]*Pool[T], 3),
        opts:      opts,
        lastScale: time.Now(),
    }

    ap.Pools[0] = NewPool[T](opts)
    ap.Pools[1] = NewPool[T](opts)
    ap.Pools[2] = NewPool[T](opts)

    ap.Pools[1].maxSize = int32(opts.MaxSize * 2)
    ap.Pools[2].maxSize = int32(opts.MaxSize * 4)

    if opts.EnableAutoTune {
        go ap.scaler()
    }

    return ap
}

func (ap *AutoPool[T]) scaler() {
    ticker := time.NewTicker(ap.opts.TuneInterval)
    defer ticker.Stop()

    for range ticker.C {
        stress := GetGCStress()
        ap.mu.Lock()

        var newIndex int32
        switch {
        case stress > 70:
            newIndex = 2
        case stress > 40:
            newIndex = 1
        default:
            newIndex = 0
        }

        if atomic.LoadInt32(&ap.current) != newIndex {
            atomic.StoreInt32(&ap.current, newIndex)
            atomic.AddUint64(&ap.scaleOps, 1)
            ap.lastScale = time.Now()
        }

        ap.mu.Unlock()
    }
}

func (ap *AutoPool[T]) Get() *T {
    idx := atomic.LoadInt32(&ap.current)
    return ap.Pools[idx].Get()
}

func (ap *AutoPool[T]) Put(obj *T) {
    idx := atomic.LoadInt32(&ap.current)
    ap.Pools[idx].Put(obj)
}

func (ap *AutoPool[T]) Stats() (totalHits, totalMisses, totalAllocs uint64, scaleOps uint64) {
    ap.mu.RLock()
    defer ap.mu.RUnlock()

    for _, pool := range ap.Pools {
        hits, misses, allocs, _ := pool.Stats()
        totalHits += hits
        totalMisses += misses
        totalAllocs += allocs
    }

    return totalHits, totalMisses, totalAllocs, atomic.LoadUint64(&ap.scaleOps)
}

func (ap *AutoPool[T]) ForceScale(index int32) {
    if index < 0 || index > 2 {
        return
    }
    atomic.StoreInt32(&ap.current, index)
    atomic.AddUint64(&ap.scaleOps, 1)
}

// SetResetFunc sets reset function for all pools
func (ap *AutoPool[T]) SetResetFunc(fn func(*T)) {
    ap.mu.Lock()
    defer ap.mu.Unlock()
    
    for _, pool := range ap.Pools {
        pool.SetResetFunc(fn)
    }
}
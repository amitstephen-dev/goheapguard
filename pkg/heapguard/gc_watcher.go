package heapguard

import (
    "runtime"
    "sync"
    "sync/atomic"
    "time"
)

var (
    gcMutex        sync.RWMutex
    lastGC         uint32
    gcCount        uint32
    lastPause      time.Duration
    gcStress       int32
    watcherRunning int32
)

func StartGCWatcher(interval time.Duration) {
    if !atomic.CompareAndSwapInt32(&watcherRunning, 0, 1) {
        return
    }

    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()

        for range ticker.C {
            updateGCStats()
        }
    }()
}

func updateGCStats() {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    gcMutex.Lock()
    defer gcMutex.Unlock()

    gcCount = memStats.NumGC
    if gcCount > 0 {
        lastPause = time.Duration(memStats.PauseTotalNs) * time.Nanosecond
    }

    stress := 10
    if memStats.NumGC > 0 {
        avgPauseNs := memStats.PauseTotalNs / uint64(memStats.NumGC)
        if avgPauseNs > 1_000_000 {
            stress = 90
        } else if avgPauseNs > 500_000 {
            stress = 70
        } else if avgPauseNs > 100_000 {
            stress = 50
        } else if avgPauseNs > 50_000 {
            stress = 30
        }
    }

    if stress > 100 {
        stress = 100
    }

    atomic.StoreInt32(&gcStress, int32(stress))
}

func GetGCStress() int {
    return int(atomic.LoadInt32(&gcStress))
}

func GetGCStats() (count uint32, totalPause time.Duration) {
    gcMutex.RLock()
    defer gcMutex.RUnlock()
    return gcCount, lastPause
}
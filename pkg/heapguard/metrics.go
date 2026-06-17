package heapguard

import (
    "fmt"
    "runtime"
    "time"
)

type MetricsCollector struct {
    startTime time.Time
}

var defaultMetrics = &MetricsCollector{
    startTime: time.Now(),
}

func GetGlobalMetrics() string {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    stress := GetGCStress()
    gcCount, gcPause := GetGCStats()

    return fmt.Sprintf(`=== goheapguard Metrics ===
GC Stress: %d/100
GC Count: %d
Total GC Pause: %v
Heap Alloc: %.2f MB
Heap Sys: %.2f MB
Number of GC: %d
`,
        stress,
        gcCount,
        gcPause,
        float64(memStats.HeapAlloc)/1024/1024,
        float64(memStats.HeapSys)/1024/1024,
        memStats.NumGC,
    )
}
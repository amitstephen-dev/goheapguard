package heapguard

import "time"

type Options struct {
    MaxSize       int
    InitialSize   int
    EnableMetrics bool
    EnableAutoTune bool
    TuneInterval  time.Duration
}

func DefaultOptions() Options {
    return Options{
        MaxSize:        1000,
        InitialSize:    100,
        EnableMetrics:  true,
        EnableAutoTune: true,
        TuneInterval:   5 * time.Second,
    }
}
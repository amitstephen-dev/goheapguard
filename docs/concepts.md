# Concepts & Architecture

## How goheapguard Works

### Object Pooling
Objects are created once and reused multiple times, eliminating GC pressure.

### GC Monitoring
The library monitors GC frequency and duration to detect stress.

### Auto-Scaling
Based on GC stress, the pool automatically adjusts its size.

### Reset Functions
Objects are cleaned before reuse to prevent data leaks.
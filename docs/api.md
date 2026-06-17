# API Reference

## Pool[T]

The main pool type for managing object reuse.

### Functions

#### NewPool[T any](opts Options) *Pool[T]
Creates a new pool with the given options.

#### (p *Pool[T]) Get() *T
Retrieves an object from the pool. Creates a new one if none available.

#### (p *Pool[T]) Put(obj *T)
Returns an object to the pool after resetting it.

#### (p *Pool[T]) SetResetFunc(fn func(*T))
Sets a custom reset function for clearing objects.

#### (p *Pool[T]) Stats() (hits, misses, allocs uint64, size int32)
Returns pool performance metrics.

## AutoPool[T]

An auto-scaling pool that adjusts size based on GC pressure.

### Functions

#### NewAutoPool[T any](opts Options) *AutoPool[T]
Creates a new auto-scaling pool.

#### (ap *AutoPool[T]) Get() *T
Retrieves an object from the current pool.

#### (ap *AutoPool[T]) Put(obj *T)
Returns an object to the current pool.

#### (ap *AutoPool[T]) Stats() (totalHits, totalMisses, totalAllocs uint64, scaleOps uint64)
Returns metrics from all pools.

## Builder[T]

A builder pattern for constructing complex objects.

### Functions

#### NewBuilder[T any](pool *Pool[T]) *Builder[T]
Creates a new builder.

#### (b *Builder[T]) Get() *T
Starts building a new object.

#### (b *Builder[T]) Build() *T
Finalizes and returns the built object.

#### (b *Builder[T]) Reset()
Resets the builder without returning an object.

## GC Watcher

#### StartGCWatcher(interval time.Duration)
Starts background GC monitoring.

#### GetGCStress() int
Returns current GC pressure (0-100).

#### GetGCStats() (count uint32, totalPause time.Duration)
Returns detailed GC statistics.

## JSON Parser

#### NewParser[T any](opts Options) *Parser[T]
Creates a new JSON parser with object pooling.

#### (p *Parser[T]) Parse(r io.Reader) (*T, error)
Parses JSON from an io.Reader.

#### (p *Parser[T]) ParseBytes(data []byte) (*T, error)
Parses JSON from a byte slice.

#### (p *Parser[T]) Release(obj *T)
Returns the parsed object to the pool.

## HTTP Middleware

#### Middleware(next http.HandlerFunc) http.HandlerFunc
Wraps a handler with GC-aware request handling.

#### WithPooledHandler[T any](pool *Pool[T], handler func(*T, http.ResponseWriter, *http.Request)) http.HandlerFunc
Wraps a handler with request-scoped object pooling.

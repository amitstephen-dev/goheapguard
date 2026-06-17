package json

import (
    "encoding/json"
    "io"

    "github.com/amitstephen-dev/goheapguard/pkg/heapguard"
)

type Parser[T any] struct {
    pool *heapguard.Pool[T]
}

func NewParser[T any](opts heapguard.Options) *Parser[T] {
    return &Parser[T]{
        pool: heapguard.NewPool[T](opts),
    }
}

func (p *Parser[T]) Parse(r io.Reader) (*T, error) {
    obj := p.pool.Get()

    decoder := json.NewDecoder(r)
    if err := decoder.Decode(obj); err != nil {
        p.pool.Put(obj)
        return nil, err
    }

    return obj, nil
}

func (p *Parser[T]) ParseBytes(data []byte) (*T, error) {
    obj := p.pool.Get()

    if err := json.Unmarshal(data, obj); err != nil {
        p.pool.Put(obj)
        return nil, err
    }

    return obj, nil
}

func (p *Parser[T]) Release(obj *T) {
    if obj != nil {
        p.pool.Put(obj)
    }
}
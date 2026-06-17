# goheapguard - GC-Aware Object Pooling for Go

[![Go Version](https://img.shields.io/github/go-mod/go-version/amitstephen-dev/goheapguard)](https://golang.org)
[![GitHub release](https://img.shields.io/github/v/release/amitstephen-dev/goheapguard)](https://github.com/amitstephen-dev/goheapguard/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/amitstephen-dev/goheapguard)](https://goreportcard.com/report/github.com/amitstephen-dev/goheapguard)
[![Go Reference](https://pkg.go.dev/badge/github.com/amitstephen-dev/goheapguard.svg)](https://pkg.go.dev/github.com/amitstephen-dev/goheapguard)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/amitstephen-dev/goheapguard)](https://github.com/amitstephen-dev/goheapguard/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/amitstephen-dev/goheapguard)](https://github.com/amitstephen-dev/goheapguard/issues)

goheapguard is a GC-aware object pooling library for Go. Reuse objects, reduce allocation pressure, monitor garbage collection activity, and build high-throughput applications with minimal overhead.

## ✨ Features

* **Generic Pools** - Type-safe object pools using Go generics
* **Object Reuse** - Reduce allocation pressure and GC overhead
* **Auto-Tuning** - Adjust pool behavior based on GC activity
* **Custom Reset Logic** - Clean objects before returning them to the pool
* **Built-in Metrics** - Monitor pool usage and allocation statistics
* **HTTP Integration** - Request-scoped object pooling middleware
* **JSON Utilities** - Reusable JSON parsing helpers
* **High Concurrency** - Designed for multi-core workloads

## 📋 Requirements

* Go 1.22 or later

## 🚀 Quick Start

### Installation

```bash
go get github.com/amitstephen-dev/goheapguard
```

### Basic Usage

```go
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
	pool := heapguard.NewPool[User](
		heapguard.DefaultOptions(),
	)

	user := pool.Get()

	user.ID = 1
	user.Name = "Alice"

	fmt.Printf("User: %+v\n", user)

	pool.Put(user)

	hits, misses, allocs, size := pool.Stats()

	fmt.Printf(
		"Hits=%d Misses=%d Allocs=%d Size=%d\n",
		hits,
		misses,
		allocs,
		size,
	)
}
```

## 📚 Core Concepts

### 1. Create a Pool

```go
opts := heapguard.DefaultOptions()

pool := heapguard.NewPool[User](opts)
```

### 2. Get an Object

```go
user := pool.Get()

user.ID = 1
user.Name = "Alice"
```

### 3. Return an Object

```go
pool.Put(user)
```

### 4. Custom Reset Logic

```go
pool.SetResetFunc(func(u *User) {
	u.ID = 0
	u.Name = ""
})
```

### 5. Pool Statistics

```go
hits, misses, allocs, size := pool.Stats()

fmt.Printf(
	"Hits=%d Misses=%d Allocs=%d Size=%d\n",
	hits,
	misses,
	allocs,
	size,
)
```

## 🌐 HTTP Integration

```go
package main

import (
	"net/http"

	"github.com/amitstephen-dev/goheapguard/pkg/heapguard"
	httpmw "github.com/amitstephen-dev/goheapguard/pkg/http"
)

type User struct {
	ID   int
	Name string
}

func main() {
	pool := heapguard.NewPool[User](
		heapguard.DefaultOptions(),
	)

	http.HandleFunc(
		"/users",
		httpmw.WithPooledHandler(
			pool,
			handleUser,
		),
	)

	http.ListenAndServe(":8080", nil)
}

func handleUser(
	user *User,
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Write([]byte("OK"))
}
```

## 📦 JSON Parsing

```go
package main

import (
	"strings"

	"github.com/amitstephen-dev/goheapguard/pkg/heapguard"
	heapjson "github.com/amitstephen-dev/goheapguard/pkg/json"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	parser := heapjson.NewParser[Person](
		heapguard.DefaultOptions(),
	)

	data := `{"name":"John","age":30}`

	person, err := parser.Parse(
		strings.NewReader(data),
	)
	if err != nil {
		panic(err)
	}

	parser.Release(person)
}
```

## 🔧 Configuration

### Default Configuration

```go
opts := heapguard.DefaultOptions()
```

### Custom Configuration

```go
opts := heapguard.Options{
	MaxSize:        2000,
	InitialSize:    500,
	EnableMetrics:  true,
	EnableAutoTune: true,
}
```

| Option         | Description                      |
| -------------- | -------------------------------- |
| MaxSize        | Maximum pool size                |
| InitialSize    | Initial number of pooled objects |
| EnableMetrics  | Enable pool statistics           |
| EnableAutoTune | Enable GC-aware tuning           |

## 💾 Pool Types

### Fixed Pool

```go
pool := heapguard.NewPool[User](opts)
```

Best for predictable workloads.

### Auto-Tuning Pool

```go
pool := heapguard.NewAutoPool[User](opts)
```

Automatically adjusts behavior based on runtime GC activity.

## 📊 Metrics

### Global Metrics

```go
metrics := heapguard.GetGlobalMetrics()

fmt.Println(metrics)
```

### HTTP Metrics Endpoint

```go
http.HandleFunc(
	"/debug/heapguard",
	func(w http.ResponseWriter, r *http.Request) {
		w.Write(
			[]byte(heapguard.GetGlobalMetrics()),
		)
	},
)
```

Example output:

```text
=== goheapguard Metrics ===
GC Stress: 10/100
GC Count: 5
Heap Alloc: 2.45 MB
Heap Sys: 8.12 MB
```

## 📁 Examples

Check the `examples/` directory for complete working examples:

* `basic/` - Basic object pooling
* `http/` - HTTP request pooling
* `json/` - JSON parsing example
* `metrics/` - Monitoring and statistics
* `autopool/` - Auto-tuning pool example

## 🏗️ Architecture

goheapguard combines:

* Generic object pools
* Runtime GC monitoring
* Custom reset hooks
* Pool metrics
* Auto-tuning capabilities

The library is designed to help reduce allocation pressure by reusing frequently allocated objects while keeping APIs simple and type-safe.

## 🎯 Use Cases

### High-Throughput APIs

Reduce allocation overhead in request-heavy services.

### Background Workers

Reuse temporary objects during task processing.

### Event Processing Systems

Handle large event streams with reduced memory churn.

### JSON-Heavy Applications

Reuse parsing structures across requests.

### Game Servers

Manage frequently created session and state objects.

### Data Processing Pipelines

Improve efficiency when processing large datasets.

## 🤝 Contributing

Contributions are welcome.

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Open a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

* Inspired by Go's `sync.Pool`
* Built for performance-focused Go applications
* Thanks to the Go community for feedback

---

**Made with ❤️ for Go developers**

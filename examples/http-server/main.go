package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/amitstephen-dev/goheapguard/pkg/heapguard"
    httpmw "github.com/amitstephen-dev/goheapguard/pkg/http"
)

type User struct {
    ID    int      `json:"id"`
    Name  string   `json:"name"`
    Email string   `json:"email"`
    Tags  []string `json:"tags"`
}

func (u *User) Reset() {
    u.ID = 0
    u.Name = ""
    u.Email = ""
    if u.Tags != nil {
        u.Tags = u.Tags[:0]
    }
}

func main() {
    heapguard.StartGCWatcher(2 * time.Second)

    opts := heapguard.DefaultOptions()
    opts.MaxSize = 1000
    pool := heapguard.NewAutoPool[User](opts)

    pool.SetResetFunc(func(u *User) {
        u.ID = 0
        u.Name = ""
        u.Email = ""
        if u.Tags != nil {
            u.Tags = u.Tags[:0]
        }
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    http.HandleFunc("/user", httpmw.WithPooledHandler(pool.Pools[0], handleUser))

    http.HandleFunc("/debug/heapguard", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(heapguard.GetGlobalMetrics()))
    })

    fmt.Println("🚀 Server starting on :8080")
    fmt.Println("📝 Test: curl -X POST http://localhost:8080/user -d '{\"id\":1,\"name\":\"John\",\"email\":\"john@test.com\"}' -H 'Content-Type: application/json'")
    fmt.Println("📊 Metrics: curl http://localhost:8080/debug/heapguard")
    
    http.ListenAndServe(":8080", nil)
}

func handleUser(user *User, w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", 405)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(user); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    fmt.Printf("✅ Received user: %+v\n", user)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status":  "ok",
        "message": "User processed with zero allocations",
    })
}
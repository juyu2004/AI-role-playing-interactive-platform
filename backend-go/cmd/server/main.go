package main

import (
    "log"
    "net/http"
    "os"

    "ai-role-playing-platform/backend-go/internal/handlers"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })

    mux.HandleFunc("/api/roles", handlers.HandleRoles)
    mux.HandleFunc("/api/roles/", handlers.HandleRoleDetail)
    mux.HandleFunc("/api/chat", handlers.HandleChat)

    log.Printf("Go backend listening on :%s", port)
    if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
        log.Fatal(err)
    }
}

func withCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}



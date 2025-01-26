package main

import (
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
    "github.com/proSamik/aiSecretary/server/handlers"
    "github.com/proSamik/aiSecretary/server/db"
)

func main() {
    log.Printf("Starting server initialization...")

    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found: %v", err)
    }

    // Verify environment variables
    requiredEnvVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
    for _, envVar := range requiredEnvVars {
        if os.Getenv(envVar) == "" {
            log.Fatalf("Required environment variable %s is not set", envVar)
        }
    }

    // Initialize database
    if err := db.InitDB(); err != nil {
        log.Fatalf("Error initializing database: %v", err)
    }

    log.Printf("Database initialized successfully")

    // Create router
    r := mux.NewRouter()

    // CORS Middleware
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            next.ServeHTTP(w, r)
        })
    })

    // Routes
    r.HandleFunc("/api/todos", handlers.CreateTodo).Methods("POST", "OPTIONS")
    r.HandleFunc("/api/todos", handlers.ListTodos).Methods("GET", "OPTIONS")
    r.HandleFunc("/api/todos/{id}", handlers.GetTodo).Methods("GET", "OPTIONS")
    r.HandleFunc("/api/todos/{id}", handlers.UpdateTodo).Methods("PUT", "OPTIONS")
    r.HandleFunc("/api/todos/{id}", handlers.DeleteTodo).Methods("DELETE", "OPTIONS")

    // WebSocket route
    r.HandleFunc("/ws", handlers.HandleWebSocket)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
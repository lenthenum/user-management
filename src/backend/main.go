package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"service_name", "go-user-api",
		"env", "production",
	)
}

func main() {
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        logger.Error("database connection failed", "error", err)
        os.Exit(1)
    }
    defer db.Close()

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")
    if err != nil {
        logger.Error("migration failed", "error", err)
        os.Exit(1)
    }

    router := mux.NewRouter()
    router.Use(requestIDMiddleware)

    router.HandleFunc("/ready", healthCheck(db)).Methods("GET")
    router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "OK")
    })

    router.HandleFunc("/api/go/users", getUsers(db)).Methods("GET")
    router.HandleFunc("/api/go/users", createUser(db)).Methods("POST")
    router.HandleFunc("/api/go/users/{id}", getUser(db)).Methods("GET")
    router.HandleFunc("/api/go/users/{id}", updateUser(db)).Methods("PUT")
    router.HandleFunc("/api/go/users/{id}", deleteUser(db)).Methods("DELETE")

    enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

    fmt.Println("Server starting on :8000...")
    if err := http.ListenAndServe(":8000", enhancedRouter); err != nil {
        logger.Error("server failed", "error", err)
        os.Exit(1)
    }
}

func healthCheck(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, _ := r.Context().Value("trace_id").(string)
		err := db.Ping()
		if err != nil {
			logger.Error("healthcheck: db unreachable", "error", err, "trace_id", traceID)
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "error": "database unreachable"})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Trace-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), "trace_id", traceID)
		w.Header().Set("X-Trace-ID", traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, _ := r.Context().Value("trace_id").(string)
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			logger.Error("query failed", "error", err, "trace_id", traceID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.Id, &u.Name, &u.Email); err != nil {
				logger.Error("scan failed", "error", err, "trace_id", traceID)
				continue
			}
			users = append(users, u)
		}

		logger.Info("getUsers called", "count", len(users), "trace_id", traceID)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			logger.Error("encode failed", "error", err, "trace_id", traceID)
		}
	}
}

func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, _ := r.Context().Value("trace_id").(string)
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.Id, &u.Name, &u.Email)
		if err != nil {
			logger.Warn("user not found", "id", id, "trace_id", traceID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(u); err != nil {
			logger.Error("encode failed", "error", err, "trace_id", traceID)
		}
	}
}

func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, _ := r.Context().Value("trace_id").(string)
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			logger.Error("decode failed", "error", err, "trace_id", traceID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.Id)
		if err != nil {
			logger.Error("insert failed", "error", err, "trace_id", traceID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("user created", "id", u.Id, "trace_id", traceID)
		if err := json.NewEncoder(w).Encode(u); err != nil {
			logger.Error("encode failed", "error", err, "trace_id", traceID)
		}
	}
}

func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, _ := r.Context().Value("trace_id").(string)
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			logger.Error("decode failed", "error", err, "trace_id", traceID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
		if err != nil {
			logger.Error("update failed", "error", err, "id", id, "trace_id", traceID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var updatedUser User
		err = db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&updatedUser.Id, &updatedUser.Name, &updatedUser.Email)
		if err != nil {
			logger.Error("post-update fetch failed", "error", err, "id", id, "trace_id", traceID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
			logger.Error("encode failed", "error", err, "trace_id", traceID)
		}
	}
}

func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, _ := r.Context().Value("trace_id").(string)
		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
		if err != nil {
			logger.Warn("delete failed: user not found", "id", id, "trace_id", traceID)
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
			if err != nil {
				logger.Error("delete execution failed", "error", err, "id", id, "trace_id", traceID)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			logger.Info("user deleted", "id", id, "trace_id", traceID)
			if err := json.NewEncoder(w).Encode("User deleted"); err != nil {
				logger.Error("encode failed", "error", err, "trace_id", traceID)
			}
		}
	}
}

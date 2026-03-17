package api

import (
	"database/sql"
	"net/http"

	sqliteinfra "tech-memo/internal/infrastructure/persistence/sqlite"
)

func BuildApp(dbPath string) (http.Handler, error) {
	db, err := sqliteinfra.Open(dbPath)
	if err != nil {
		return nil, err
	}

	return newRouter(db), nil
}

func newRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		if err := db.PingContext(r.Context()); err != nil {
			http.Error(w, "database unavailable", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"tech-memo API base is running"}`))
	})

	return mux
}

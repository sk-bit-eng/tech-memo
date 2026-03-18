package api

import (
	"database/sql"
	"net/http"
	"os"

	sqliteinfra "tech-memo/internal/infrastructure/persistence/sqlite"

	"github.com/gin-gonic/gin"
)

func NewHandler() (http.Handler, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "tech-memo.db"
	}

	db, err := sqliteinfra.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return newRouter(db), nil
}

func newRouter(db *sql.DB) http.Handler {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r
}

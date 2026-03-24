package api

import (
	"net/http"
	"os"

	sqlserverinfra "tech-memo/internal/infrastructure/persistence/sqlserver"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewHandler() (http.Handler, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "sqlserver://sa:Test@1234@localhost:1433?database=tech_memo"
	}

	db, err := sqlserverinfra.Open(dsn)
	if err != nil {
		return nil, err
	}
	return newRouter(db), nil
}

func newRouter(db *gorm.DB) http.Handler {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "teck-memo server is running",
		})
	})

	return r
}

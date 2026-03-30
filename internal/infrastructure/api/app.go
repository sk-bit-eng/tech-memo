package api

import (
	"net/http"
	"os"

	memoctrl "tech-memo/internal/adapter/controller/memo"
	todoctrl "tech-memo/internal/adapter/controller/todo"
	memorepo "tech-memo/internal/adapter/gateway/memo"
	todorepo "tech-memo/internal/adapter/gateway/todo"
	memouc "tech-memo/internal/application/usecase/memo"
	todouc "tech-memo/internal/application/usecase/todo"
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

	// DI
	memoRepo := memorepo.NewRepository(db)
	todoRepo := todorepo.NewRepository(db)
	memoUC := memouc.NewInteractor(memoRepo)
	todoUC := todouc.NewInteractor(todoRepo)
	memoCtrl := memoctrl.NewController(memoUC)
	todoCtrl := todoctrl.NewController(todoUC)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "teck-memo server is running"})
	})

	memos := r.Group("/memos")
	{
		memos.POST("", memoCtrl.Create)
		memos.GET("/:id", memoCtrl.GetByID)
		memos.PUT("/:id", memoCtrl.Update)
		memos.DELETE("/:id", memoCtrl.Delete)
		memos.PATCH("/:id/pin", memoCtrl.TogglePin)
	}

	users := r.Group("/users/:userID")
	{
		users.GET("/memos", memoCtrl.ListByUser)
		users.GET("/memos/search", memoCtrl.Search)
		users.GET("/memos/category/:categoryID", memoCtrl.ListByCategory)

		users.GET("/todos", todoCtrl.ListByUser)
		users.GET("/todos/search", todoCtrl.Search)
		users.GET("/todos/pending", todoCtrl.ListPending)
		users.GET("/todos/completed", todoCtrl.ListCompleted)
		users.GET("/todos/category/:categoryID", todoCtrl.ListByCategory)
	}

	todos := r.Group("/todos")
	{
		todos.POST("", todoCtrl.Create)
		todos.GET("/:id", todoCtrl.GetByID)
		todos.PUT("/:id", todoCtrl.Update)
		todos.DELETE("/:id", todoCtrl.Delete)
		todos.PATCH("/:id/pin", todoCtrl.TogglePin)
		todos.PATCH("/:id/complete", todoCtrl.Complete)
		todos.PATCH("/:id/incomplete", todoCtrl.Incomplete)
	}

	return r
}

// internal/infrastructure/api/di.go
package api

import (
	"net/http"

	dbgateway "tech-memo/internal/adapter/gateway"
	"tech-memo/internal/adapter/controller"
	"tech-memo/internal/application/usecase"
)

func BuildApp(dbPath string) (http.Handler, error) {
	gw, err := dbgateway.NewSQLiteMemoGateway(dbPath)
	if err != nil {
		return nil, err
	}

	uc := usecase.NewMemoInteracter(gw)
	ctrl := controller.NewMemoController(uc)
	h := NewMemoHandler(ctrl)
	return newRouter(h), nil
}

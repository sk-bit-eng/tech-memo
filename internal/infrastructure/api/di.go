// internal/infrastructure/api/di.go
package api

import (
	"net/http"

	dbgateway "tech-memo/internal/adapter/gateway"
	"tech-memo/internal/adapter/controller"
	"tech-memo/internal/application/interacter"
)

func BuildApp(dbPath string) (http.Handler, error) {
	gw, err := dbgateway.NewSQLiteMemoGateway(dbPath)
	if err != nil {
		return nil, err
	}

	uc := interacter.NewMemoInteracter(gw)
	ctrl := controller.NewMemoController(uc)
	h := NewMemoHandler(ctrl)
	return newRouter(h), nil
}

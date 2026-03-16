// internal/adapter/controller/memo_controller.go
package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"tech-memo/internal/application/usecase"
	"tech-memo/internal/domain"
)

type MemoController struct {
	uc usecase.MemoUseCase
}

func NewMemoController(uc usecase.MemoUseCase) *MemoController {
	return &MemoController{uc: uc}
}

func (c *MemoController) List(r *http.Request) ([]*domain.Memo, error) {
	q := r.URL.Query().Get("q")
	tag := r.URL.Query().Get("tag")
	switch {
	case q != "":
		return c.uc.Search(q)
	case tag != "":
		return c.uc.FindByTag(tag)
	default:
		return c.uc.GetAll()
	}
}

func (c *MemoController) GetByID(r *http.Request) (*domain.Memo, error) {
	id := extractID(r.URL.Path)
	return c.uc.GetByID(id)
}

type MemoRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	Language string   `json:"language"`
}

func (c *MemoController) Create(r *http.Request) (*domain.Memo, error) {
	var req MemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Tags == nil {
		req.Tags = []string{}
	}
	return c.uc.Create(req.Title, req.Content, req.Tags, req.Language)
}

func (c *MemoController) Update(r *http.Request) (*domain.Memo, error) {
	id := extractID(r.URL.Path)
	var req MemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Tags == nil {
		req.Tags = []string{}
	}
	return c.uc.Update(id, req.Title, req.Content, req.Tags, req.Language)
}

func (c *MemoController) Delete(r *http.Request) error {
	id := extractID(r.URL.Path)
	return c.uc.Delete(id)
}

func extractID(path string) string {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	return parts[len(parts)-1]
}

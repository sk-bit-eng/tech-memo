// internal/infrastructure/api/handler.go
package api

import (
	"net/http"
	"strings"

	"tech-memo/internal/adapter/controller"
	"tech-memo/internal/adapter/presenter"
)

type MemoHandler struct {
	ctrl *controller.MemoController
}

func NewMemoHandler(ctrl *controller.MemoController) *MemoHandler {
	return &MemoHandler{ctrl: ctrl}
}

func (h *MemoHandler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.ctrl.List(r)
	if err != nil {
		presenter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	presenter.WriteJSON(w, http.StatusOK, result)
}

func (h *MemoHandler) Get(w http.ResponseWriter, r *http.Request) {
	memo, err := h.ctrl.GetByID(r)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			presenter.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		presenter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	presenter.WriteJSON(w, http.StatusOK, memo)
}

func (h *MemoHandler) Create(w http.ResponseWriter, r *http.Request) {
	memo, err := h.ctrl.Create(r)
	if err != nil {
		presenter.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	presenter.WriteJSON(w, http.StatusCreated, memo)
}

func (h *MemoHandler) Update(w http.ResponseWriter, r *http.Request) {
	memo, err := h.ctrl.Update(r)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			presenter.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		presenter.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	presenter.WriteJSON(w, http.StatusOK, memo)
}

func (h *MemoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.ctrl.Delete(r); err != nil {
		if strings.Contains(err.Error(), "not found") {
			presenter.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		presenter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

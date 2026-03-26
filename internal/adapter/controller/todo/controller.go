package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	tododto "tech-memo/internal/application/dto/todo"
	todouc "tech-memo/internal/application/usecase/todo"
)

type Controller struct {
	uc todouc.UseCase
}

func NewController(uc todouc.UseCase) *Controller {
	return &Controller{uc: uc}
}

func (ctrl *Controller) GetByID(c *gin.Context) {
	todo, err := ctrl.uc.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (ctrl *Controller) ListByUser(c *gin.Context) {
	todos, err := ctrl.uc.ListByUser(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (ctrl *Controller) ListByCategory(c *gin.Context) {
	todos, err := ctrl.uc.ListByCategory(c.Param("userID"), c.Param("categoryID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (ctrl *Controller) ListPending(c *gin.Context) {
	todos, err := ctrl.uc.ListPending(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (ctrl *Controller) ListCompleted(c *gin.Context) {
	todos, err := ctrl.uc.ListCompleted(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (ctrl *Controller) Search(c *gin.Context) {
	todos, err := ctrl.uc.Search(c.Param("userID"), c.Query("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (ctrl *Controller) Create(c *gin.Context) {
	var input tododto.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := ctrl.uc.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func (ctrl *Controller) Update(c *gin.Context) {
	var input tododto.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = c.Param("id")
	todo, err := ctrl.uc.Update(input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	if err := ctrl.uc.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ctrl *Controller) TogglePin(c *gin.Context) {
	todo, err := ctrl.uc.TogglePin(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (ctrl *Controller) Complete(c *gin.Context) {
	if err := ctrl.uc.Complete(c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ctrl *Controller) Incomplete(c *gin.Context) {
	if err := ctrl.uc.Incomplete(c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

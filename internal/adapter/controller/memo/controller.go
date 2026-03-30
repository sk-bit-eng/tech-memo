package memo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	memodto "tech-memo/internal/application/dto/memo"
	memouc "tech-memo/internal/application/usecase/memo"
)

type Controller struct {
	uc memouc.UseCase
}

func NewController(uc memouc.UseCase) *Controller {
	return &Controller{uc: uc}
}

func (ctrl *Controller) GetByID(c *gin.Context) {
	memo, err := ctrl.uc.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memo)
}

func (ctrl *Controller) ListByUser(c *gin.Context) {
	memos, err := ctrl.uc.ListByUser(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memos)
}

func (ctrl *Controller) ListByCategory(c *gin.Context) {
	memos, err := ctrl.uc.ListByCategory(c.Param("userID"), c.Param("categoryID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memos)
}

func (ctrl *Controller) Search(c *gin.Context) {
	memos, err := ctrl.uc.Search(c.Param("userID"), c.Query("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memos)
}

func (ctrl *Controller) Create(c *gin.Context) {
	var input memodto.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	memo, err := ctrl.uc.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, memo)
}

func (ctrl *Controller) Update(c *gin.Context) {
	var input memodto.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = c.Param("id")
	memo, err := ctrl.uc.Update(input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memo)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	if err := ctrl.uc.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ctrl *Controller) TogglePin(c *gin.Context) {
	memo, err := ctrl.uc.TogglePin(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memo)
}

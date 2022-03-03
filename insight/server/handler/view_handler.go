package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ViewHandler struct{}

func (v *ViewHandler) GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "frontend", nil)
	return
}
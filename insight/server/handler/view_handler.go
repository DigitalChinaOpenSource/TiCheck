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

func (v *ViewHandler) GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "frontend", gin.H{})
	return
}

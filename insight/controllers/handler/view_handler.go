package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ViewHandler struct{}

func (v ViewHandler) GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"index" : "首页",
	})
	return
}
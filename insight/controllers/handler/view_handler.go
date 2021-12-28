package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ViewHandler struct{}

func (v ViewHandler) GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"index": "首页",
	})
	return
}

func (v ViewHandler) GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
	return
}

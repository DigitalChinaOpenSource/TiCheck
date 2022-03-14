package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/// gin api return object, for example:
/// c.JSON(http.StatusOK, &ApiResp{})
type ApiResp struct {
	Success bool        `json:"success"` // business success or not
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type ViewHandler struct{}

func (v *ViewHandler) GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "frontend", nil)
	return
}

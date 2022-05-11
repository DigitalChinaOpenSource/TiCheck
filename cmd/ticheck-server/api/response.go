package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespBody
// gin api return object, for example:
// c.JSON(http.StatusOK, &ApiResp{})
type RespBody struct {
	Success bool        `json:"success"` // business success or not
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// S api success return with empty body
// request and bussiness both success, return
// Front-end suggestion mode `message.success(content, [duration], onClose)`
func S(c *gin.Context) {
	c.JSON(http.StatusOK, &RespBody{
		Success: true,
		Msg:     "",
		Data:    nil,
	})
}

// SuccessWithMsg api success return with message
func SuccessWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, &RespBody{
		Success: true,
		Msg:     msg,
		Data:    nil,
	})
}

// Success api success return with message and data
func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, &RespBody{
		Success: true,
		Msg:     msg,
		Data:    data,
	})
}

// F api fail return with empty body, the request is ok, but the business is fail
// request successed but bussiness failed, return
// example: balance not enough...
// Front-end suggestion mode `message.info(content, [duration], onClose)`
func F(c *gin.Context) {
	c.JSON(http.StatusAccepted, &RespBody{
		Success: false,
		Msg:     "",
		Data:    nil,
	})
}

// FailWithMsg api fail return with message, the request is ok, but the business is fail
func FailWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusAccepted, &RespBody{
		Success: false,
		Msg:     msg,
		Data:    nil,
	})
}

// Fail api fail return with message and data, the request is ok, but the business is fail
func Fail(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusAccepted, &RespBody{
		Success: false,
		Msg:     msg,
		Data:    data,
	})
}

// B api BadRequest return with empty body
// request and bussiness both failed, return
// example: user not exist, page not found...
// Front-end suggestion mode `message.warning(content, [duration], onClose)`
func B(c *gin.Context) {
	c.JSON(http.StatusBadRequest, &RespBody{
		Success: false,
		Msg:     "",
		Data:    nil,
	})
}

// BadWithMsg api BadRequest return with message
func BadWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, &RespBody{
		Success: false,
		Msg:     msg,
		Data:    nil,
	})
}

// Bad api BadRequest return with message and data
func Bad(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusBadRequest, &RespBody{
		Success: false,
		Msg:     msg,
		Data:    data,
	})
}

// E api error return with empty body
// server exception, return
// example: nil pointer reference, slice out of index...
// Front-end suggestion mode `message.error(content, [duration], onClose)`
func E(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, &RespBody{
		Success: false,
		Msg:     "",
		Data:    nil,
	})
}

// ErrorWithMsg api error return  with message
func ErrorWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, &RespBody{
		Success: false,
		Msg:     msg,
		Data:    nil,
	})
}

// Error api error return with message and data
func Error(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusInternalServerError, &RespBody{
		Success: false,
		Msg:     msg,
		Data:    data,
	})
}

// AuthenticationFailed failed to verify token or password
func AuthenticationFailed(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, &RespBody{
		Success: false,
		Msg:     "authentication failed",
		Data:    nil,
	})
}

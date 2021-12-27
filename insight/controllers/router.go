package controllers

import (
	"TiCheck/insight/controllers/handler"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	authGroup := engine.Group("")
	{
		auth := &handler.AuthorizerHandler{}
		// 用户认证
		authGroup.POST("/auth", auth.AuthenticatedUser)
	}
}
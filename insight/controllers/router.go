package controllers

import (
	"TiCheck/insight/controllers/handler"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	engine.Static("/assets","insight/views/assets")
	engine.LoadHTMLGlob("insight/views/index.html")

	viewGroup := engine.Group("/")
	{
		view := &handler.ViewHandler{}
		viewGroup.GET("/", view.GetIndex)
	}

	authGroup := engine.Group("/auth")
	{
		auth := &handler.AuthorizerHandler{}
		// 用户认证
		authGroup.POST("/", auth.AuthenticatedUser)
	}
}
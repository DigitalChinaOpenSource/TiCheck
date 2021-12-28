package controllers

import (
	"TiCheck/insight/controllers/handler"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	engine.Static("/assets", "./views/assets")
	engine.LoadHTMLGlob("./views/*.html")

	viewGroup := engine.Group("/")
	{
		view := &handler.ViewHandler{}
		viewGroup.GET("/", view.GetIndex)

		viewGroup.GET("/login", view.GetLogin)
	}

	authGroup := engine.Group("/auth")
	{
		auth := &handler.AuthorizerHandler{}
		// 用户认证
		authGroup.POST("/", auth.AuthenticatedUser)
	}
}

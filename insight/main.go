package main

import (
	"TiCheck/insight/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	engine := gin.Default()
	engine.Use()

	// 初始化服务
	initService(engine)

	controllers.Register(engine)

	engine.Run(":8081")
}

// initService 初始化服务
func initService(r *gin.Engine) {
	// 服务状态
	status := true
	errMsg := "OK"

	// 访问根，用于测试是否能正常访问
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"RespondStatus": status,
			"Err":           errMsg,
		})
	})
}
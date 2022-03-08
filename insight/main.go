package main

import (
	"TiCheck/insight/server"
	"TiCheck/insight/server/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.Use()

	// 初始化服务
	initService(engine)

	// 初始化数据库
	err := model.InitDB()
	if err != nil {
		panic("can't connect to db")
	}

	server.Register(engine)

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

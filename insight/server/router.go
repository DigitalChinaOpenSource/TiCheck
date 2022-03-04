package server

import (
	handler2 "TiCheck/insight/server/handler"
	"TiCheck/insight/server/model"
	"github.com/gin-contrib/multitemplate"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	// 多模板
	engine.HTMLRender = createMyRender()

	// 加载静态资源
	engine.Static("/assets", "web/dist/assets")
	engine.Static("/css", "web/dist/css")
	engine.Static("/img", "web/dist/img")
	engine.Static("/js", "web/dist/js")
	engine.StaticFile("/avatar2.jpg", "web/dist/avatar2.jpg")
	engine.StaticFile("/logo.png", "web/dist/logo.png")

	// 初始化数据库
	err := model.InitDB()
	if err != nil {
		panic("can't connect to db")
	}

	viewGroup := engine.Group("/")
	{
		view := &handler2.ViewHandler{}

		// 打开首页
		viewGroup.GET("/", view.GetIndex)

		// 未定义的路由直接重定向到 Index
		engine.NoRoute(view.GetIndex)
	}

	sessionGroup := engine.Group("/session")
	session := &handler2.SessionHandler{
		Sessions: make(map[string]*handler2.Session, 0),
	}

	{
		// 用户认证
		sessionGroup.POST("/", session.AuthenticatedUser)

		// 退出用户
		sessionGroup.POST("/logout", session.Logout)
	}

	reportGroup := engine.Group("/report")
	reportGroup.Use(session.VerifyToken)
	{
		report := &handler2.ReportHandler{}

		// 获取历史巡检列表
		reportGroup.GET("/catalog", report.GetCatalog)

		reportGroup.GET("/frontend/auth/login", report.GetCatalog)

		// 通过id获得某次巡检结果
		reportGroup.GET("/id/:id", report.GetReport)

		// 获取最后一次巡检结果
		reportGroup.GET("/last", report.GetLastReport)

		// 获取巡检结果元信息
		reportGroup.GET("/meta", report.GetMeta)

		// 执行一次巡检
		reportGroup.GET("/", report.ExecuteCheck)

		// 下载所有的巡检报告
		reportGroup.GET("/download/all", report.DownloadAllReport)

		// 下载指定的一次巡检报告
		reportGroup.GET("/download/:id", report.DownloadReport)

		// 编辑配置脚本
		reportGroup.POST("/editconf/:script", report.EditConfig)
	}

	clusterGroup := engine.Group("/cluster")
	{
		cluster := &handler2.ClusterHandler{}

		//
		clusterGroup.GET("/:id", cluster.GetClusterInfo)
	}

	scriptGroup := engine.Group("/script")
	// test, ignore token
	// scriptGroup.Use(session.VerifyToken)
	{
		script := &handler2.ScriptHandler{}

		// 查看所有本地脚本
		scriptGroup.GET("/local", script.GetAllLocalScript)

		// 查看所有的远程仓库脚本，获取列表
		scriptGroup.GET("/remote", script.GetAllRemoteScript)

		// 查看指定远程脚本的介绍
		scriptGroup.GET("/remote/readme/:name", script.GetReadMe)

		// 下载指定名的脚本到本地
		scriptGroup.POST("/remote/download/:name", script.DownloadScript)
	}
}

func createMyRender() multitemplate.Renderer {
	p := multitemplate.NewRenderer()
	p.AddFromFiles("frontend", "insight/web/dist/index.html")
	return p
}

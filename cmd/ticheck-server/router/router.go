package router

import (
	"TiCheck/cmd/ticheck-server/handler"
	"TiCheck/config"

	"github.com/gin-contrib/multitemplate"

	"github.com/gin-gonic/gin"
)

var (
	web_prefix string
)

func Register(engine *gin.Engine) {

	// 使用静态资源需要在 web 目录下 npm run build
	if gin.Mode() == gin.ReleaseMode {
		web_prefix = config.GlobalConfig.WorkDir + "web/dist"
		// 多模板
		engine.HTMLRender = createMyRender()
		// 加载静态资源
		engine.Static("/assets", web_prefix+"/assets")
		engine.Static("/css", web_prefix+"/css")
		engine.Static("/img", web_prefix+"/img")
		engine.Static("/js", web_prefix+"/js")
		engine.StaticFile("/avatar2.jpg", web_prefix+"/avatar2.jpg")
		engine.StaticFile("/logo.png", web_prefix+"/logo.png")
	}

	viewGroup := engine.Group("/")
	{
		view := &handler.ViewHandler{}

		// 打开首页
		viewGroup.GET("/", view.GetIndex)

		// 未定义的路由直接重定向到 Index
		engine.NoRoute(view.GetIndex)
	}

	sessionGroup := engine.Group("/session")
	//session := &handler.SessionHandler{
	//	Users:    map[string]string{},
	//	Sessions: make(map[string]*handler.Session, 0),
	//}

	{
		// 用户认证
		sessionGroup.POST("", handler.SessionHelper.AuthenticatedUser)

		// 退出用户
		sessionGroup.POST("/logout", handler.SessionHelper.Logout)

		// 获取当前用户信息
		sessionGroup.GET("/info", handler.SessionHelper.GetUserInfo)
	}

	clusterGroup := engine.Group("/cluster")
	//clusterGroup.Use(session.VerifyToken)
	{
		cluster := &handler.ClusterHandler{}

		// Get cluster list
		clusterGroup.GET("/list", cluster.GetClusterList)

		// Get cluster information by id
		clusterGroup.GET("/info/:id", cluster.GetClusterInfo)

		// Get before updated cluster information
		clusterGroup.GET("/initial/:id", cluster.GetInitialClusterInfo)

		// Add cluster
		clusterGroup.POST("/add", cluster.PostClusterInfo)

		// Update cluster by id
		clusterGroup.PUT("/update/:id", cluster.UpdateClusterInfo)

		// Get cluster scheduler list by id
		clusterGroup.GET("/scheduler/:id", cluster.GetClusterSchedulerList)

		// Add a scheduler for this cluster
		clusterGroup.POST("/scheduler/add", cluster.PostClusterScheduler)

		// Update a scheduler by its id
		clusterGroup.PUT("/scheduler/update", cluster.UpdateScheduler)

		// Delete a scheduler for its id
		clusterGroup.DELETE("/scheduler/delete/:id", cluster.DeleteScheduler)

		// Get the cluster installed probe checklist
		clusterGroup.GET("/probe/:id", cluster.GetProbeList)

		// Get probe checklist that can be added to the cluster
		clusterGroup.GET("/probe/add/:id", cluster.GetAddProbeList)

		// Add a probe to the cluster
		clusterGroup.POST("/probe", cluster.AddProbeForCluster)

		// Update probe operator and threshold in the cluster checklist
		clusterGroup.PUT("/probe/config", cluster.UpdateProbeConfig)

		// Change status for probe, open or close a check probe
		clusterGroup.PUT("/probe/status", cluster.ChangeProbeStatus)

		clusterGroup.DELETE("/probe/:id", cluster.DeleteProbeForCluster)

		// Get some information before execute check
		clusterGroup.GET("/check/info/:id", cluster.GetExecuteInfo)

		// execute check
		clusterGroup.GET("/check/run/:id", cluster.ExecuteCheck)
	}

	reportGroup := engine.Group("/cluster/report")
	//reportGroup.Use(session.VerifyToken)
	{
		report := &handler.ReportHandler{}

		// 获取历史巡检列表
		reportGroup.GET("/all/:clusterID", report.GetReportList)

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

	storeGroup := engine.Group("/store")
	// test, ignore token
	// storeGroup.Use(session.VerifyToken)
	{
		sh := &handler.StoreHandler{}

		// 查看所有本地脚本
		storeGroup.GET("/local", sh.GetLocalScript)
		storeGroup.GET("/local/readme", sh.GetLocalReadme)

		// 查看所有的远程仓库脚本，获取列表
		storeGroup.GET("/remote", sh.GetAllRemoteScript)

		// 查看所有自定义脚本
		storeGroup.GET("/custom", sh.GetCustomScript)
		storeGroup.GET("/custom/readme/:id", sh.GetCustomReadme)
		storeGroup.POST("/custom", sh.UploadCustomScript)
		storeGroup.DELETE("/custom/:id", sh.DeleteCustomReadme)

		// 查看指定远程脚本的介绍
		storeGroup.GET("/remote/readme/:name", sh.GetReadMe)

		// 下载指定名的脚本到本地
		storeGroup.POST("/remote/download/:name", sh.DownloadScript)
	}
}

func createMyRender() multitemplate.Renderer {
	p := multitemplate.NewRenderer()
	p.AddFromFiles("frontend", web_prefix+"/index.html")
	return p
}

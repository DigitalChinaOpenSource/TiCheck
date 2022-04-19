package main

import (
	"TiCheck/cmd/ticheck-server/router"
	"TiCheck/executor"
	"TiCheck/internal/model"
	"TiCheck/internal/service"
	"TiCheck/util/logutil"
	"context"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.Use()

	// LogConfig should be read from the configuration file
	config := logutil.NewLogConfig(zap.DebugLevel, lumberjack.Logger{
		Filename:   logutil.DefaultLogFilePath,
		MaxSize:    logutil.DefaultLogMaxSize,
		MaxBackups: logutil.DefaultLogBackups,
		MaxAge:     logutil.DefaultLogAge,
		Compress:   logutil.DefaultLogCompress,
	})

	logutil.InitLog(config)

	err := model.InitDB()
	if err != nil {
		logutil.Logger.Panic("Can't connect to db", zap.Error(err))
	}
	logutil.Logger.Info("Completed database initialization.")

	// route register
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome To TiCheck Server.")
	})
	router.Register(engine)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: engine,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logutil.Logger.Fatal("Failed to initialize service", zap.Error(err))
		}
		logutil.Logger.Info("Completed TiCheck server initialization.")
	}()
	//testExe()
	service.CronService.Initialize()
	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logutil.Logger.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logutil.Logger.Fatal("The server forced to shutdown: ", zap.Error(err))
	}

	logutil.Logger.Info("The server has exited.")
}

func testExe() {

	exe := executor.CreateClusterExecutor(1, 0)

	resultCh := make(chan executor.CheckResult, 10)
	// ctx := context.WithValue(context.Background(), "", "")
	go exe.Execute(resultCh)
	// for {
	// 	select {
	// 	case result := <-resultCh:
	// 		fmt.Printf("%+v\n", result)
	// 		if result.IsFinished {
	// 			return
	// 		}
	// 	}

	// }

	// res := (&model.ClusterChecklist{}).GetEnabledCheckListTagGroup(1)
	// fmt.Println(res)
}

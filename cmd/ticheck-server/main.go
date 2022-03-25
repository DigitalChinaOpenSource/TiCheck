package main

import (
	"TiCheck/cmd/ticheck-server/router"
	"TiCheck/internal/model"
	"context"
	"fmt"
	"log"
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

	err := model.InitDB()
	if err != nil {
		panic("can't connect to db")
	}

	// route register
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome To Ticheck Server.")
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
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//testExe()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func testExe() {

	// exe := executor.CreateClusterExecutor(1, 0)

	// resultCh := make(chan executor.CheckResult, 10)
	// // ctx := context.WithValue(context.Background(), "", "")
	// go exe.Execute(resultCh)
	// for {
	// 	select {
	// 	case result := <-resultCh:
	// 		fmt.Printf("%+v\n", result)
	// 		if result.IsFinished {
	// 			return
	// 		}
	// 	}

	// }

	res := (&model.ClusterChecklist{}).GetEnabledCheckListTagGroup(1)
	fmt.Println(res)
}

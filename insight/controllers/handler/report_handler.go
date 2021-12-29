package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ReportHandler struct{}

func (r *ReportHandler) GetCatalog(c *gin.Context) {
	return
}

func (r *ReportHandler) GetReport(c *gin.Context) {
	return
}

func (r *ReportHandler) GetLastReport(c *gin.Context) {
	return
}

func (r *ReportHandler) GetMeta(c *gin.Context) {
	return
}

func (r *ReportHandler) ExecuteCheck(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err !=nil {
		return
	}

	defer ws.Close()

	for {

		err = ws.WriteJSON(nil)

		if err != nil {
			break
		}

		time.Sleep(time.Second)

	}


	return
}

func (r *ReportHandler) DownloadAllReport(c *gin.Context) {
	return
}

func (r *ReportHandler) DownloadReport(context *gin.Context) {

}
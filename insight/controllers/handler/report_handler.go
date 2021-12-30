package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	if err != nil {
		return
	}

	defer ws.Close()
	i := 0
	for {

		err = ws.WriteJSON(map[string]interface{}{
			"finished":        true,
			"check_class":     "集群",
			"check_name":      "存活的TiDB数量",
			"check_item":      "TiDB节点数",
			"check_result":    "正常",
			"check_value":     5,
			"check_threshold": "等于5",
			"check_time":      20211221063030,
		})
		i++
		if err != nil || i >= 10 {
			break
		}

		time.Sleep(time.Second * 3)

	}

	return
}

func (r *ReportHandler) DownloadAllReport(c *gin.Context) {
	return
}

func (r *ReportHandler) DownloadReport(context *gin.Context) {

}

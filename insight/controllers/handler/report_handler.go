package handler

import (
	"net/http"
	"time"

	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ReportHandler struct{
	con *sql.DB
}

type CheckDara struct {
	CheckTime string	`json:"check_time"`
	CheckClass string	`json:"check_class"`
	CheckName string	`json:"check_name"`
	Operator string		`json:"operator"`
	Threshold float64   `json:"threshold"`
	Duration int		`json:"duration"`
	CheckItem string	`json:"check_item"`
	CheckValue float64	`json:"check_value"`
	CheckStatus string	`json:"check_status"`

}

func (r *ReportHandler) GetCatalog(c *gin.Context) {
	conn, err := ConnectDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":err.Error(),
		})
	}

	_ = conn.Ping()

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

func (r *ReportHandler) DownloadReport(c *gin.Context) {

}

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "/report/report.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

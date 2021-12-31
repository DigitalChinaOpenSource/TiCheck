package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var DBInstance *sql.DB

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ReportHandler struct {
	con *sql.DB
}

type CheckHistory struct {
	CheckTime    int `json:"check_time"`
	NormalItems  int `json:"normal_items"`
	WarningItems int `json:"warning_items"`
	TotalItems   int `json:"total_items"`
	Duration     int `json:"duration"`
}

type CheckData struct {
	ID          string  `json:"id"`
	CheckTime   string  `json:"check_time"`
	CheckClass  string  `json:"check_class"`
	CheckName   string  `json:"check_name"`
	Operator    string  `json:"operator"`
	Threshold   float64 `json:"threshold"`
	Duration    int     `json:"duration"`
	CheckItem   string  `json:"check_item"`
	CheckValue  float64 `json:"check_value"`
	CheckStatus string  `json:"check_status"`
}

func (r *ReportHandler) GetCatalog(c *gin.Context) {

	length, _ := strconv.Atoi(c.Query("length"))
	start, _ := strconv.Atoi(c.Query("start"))

	err := ConnectDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	var total int
	count := DBInstance.QueryRow("select count(*) from check_history")
	count.Scan(&total)

	var result = []CheckHistory{}
	rows, _ := DBInstance.Query(fmt.Sprintf("select * from check_history order by check_time desc limit %v offset %v", length, start))
	for rows.Next() {
		r := CheckHistory{}
		rows.Scan(&r.CheckTime, &r.NormalItems, &r.WarningItems, &r.TotalItems, &r.Duration)
		result = append(result, r)
	}
	rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"draw":            (start % length) + 1,
		"recordsFiltered": total,
		"recordsTotal":    total,
		"data":            result,
	})
}

func (r *ReportHandler) GetReport(c *gin.Context) {
	id := c.Param("id")
	err := ConnectDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	var result = []CheckData{}
	rows, _ := DBInstance.Query(fmt.Sprintf("select * from check_data where check_time=%v order by id", id))
	for rows.Next() {
		r := CheckData{}
		rows.Scan(&r.ID, &r.CheckTime, &r.CheckClass, &r.CheckName, &r.Operator, &r.Threshold, &r.Duration, &r.CheckItem, &r.CheckValue, &r.CheckStatus)
		result = append(result, r)
	}
	rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"total": len(result),
		"id":    id,
		"data":  result,
	})

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

// DownloadReport 下载指定报告
func (r *ReportHandler) DownloadReport(c *gin.Context) {
	reportId := c.Param("id")
	fileName := reportId + ".csv"

	_, err := os.Open("../report/" + fileName)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the report is not found",
		})
	}

	c.Header("Content-Type", "application/x-xls")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.File("../report/" + fileName)
}

func ConnectDB() error {
	if DBInstance == nil {
		db, err := sql.Open("sqlite3", "../report/report.db")
		if err != nil {
			return err
		}
		err = db.Ping()
		if err != nil {
			return err
		}
		DBInstance = db
	}

	return nil
}

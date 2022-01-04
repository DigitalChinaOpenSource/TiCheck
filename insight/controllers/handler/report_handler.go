package handler

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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

type ReportHandler struct {
	Conn *sql.DB
}

type CheckHistory struct {
	CheckTime    int `json:"check_time"`
	NormalItems  int `json:"normal_items"`
	WarningItems int `json:"warning_items"`
	TotalItems   int `json:"total_items"`
	Duration     int `json:"duration"`
}

type CheckData struct {
	ID          int     `json:"id"`
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

	err := r.ConnectDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var total int
	count := r.Conn.QueryRow("select count(*) from check_history")
	count.Scan(&total)

	var result = []CheckHistory{}
	rows, _ := r.Conn.Query(fmt.Sprintf("select * from check_history order by check_time desc limit %v offset %v", length, start))
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
	err := r.ConnectDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	var result = []CheckData{}
	rows, _ := r.Conn.Query(fmt.Sprintf("select * from check_data where check_time=%v order by id", id))
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer ws.Close()

	err = r.ConnectDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 获取一个执行时间戳，作为执行时间，同时将其传给执行脚本
	executeTime := time.Now().Unix()

	// 监听脚本是否完成
	done := make(chan bool)
	go r.executeScript(executeTime, done)

	// 从数据库中获取实时结果
	resultCh := make(chan *CheckData, 10)
	go r.getResult(executeTime, resultCh, done)

	// 设置一分钟的超时时间
	ticker := time.NewTicker(time.Minute)

	select {
	case result := <-resultCh:
		err = ws.WriteJSON(result)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
	case <-done:
		c.JSON(http.StatusOK, gin.H{
			"finish": true,
		})
		return
	case <-ticker.C:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "execute check time out",
		})
		return
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

		return
	}

	c.Header("Content-Type", "application/x-xls")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.File("../report/" + fileName)
}

func (r *ReportHandler) ConnectDB() error {
	if r.Conn == nil {
		db, err := sql.Open("sqlite3", "../report/report.db")
		if err != nil {
			return err
		}
		err = db.Ping()
		if err != nil {
			return err
		}
		r.Conn = db
	}

	return nil
}

func (r *ReportHandler) executeScript(executeTime int64, done chan bool) {
	cmd := exec.Command("../run/run.sh", string(executeTime))
	cmd.Run()
	done <- true
}

func (r *ReportHandler) getResult(executeTime int64, ch chan *CheckData, done chan bool) {
	// 记录每一次查询到的最新数据，下一轮查询从这里开始
	var index int

	result := &CheckData{}

	select {
	case <-done:
		return
	default:
		querySQL := fmt.Sprintf("select * from check_data where check_time == %d and index > %d", executeTime, index)
		rows, err := r.Conn.Query(querySQL)
		if err != nil {
			return
		}

		for rows.Next() {
			//row.Scan(result)

			rows.Scan(&result.ID, &result.CheckTime, &result.CheckClass, &result.CheckName,
				&result.Operator, &result.Threshold, &result.Duration, &result.CheckItem,
				&result.CheckValue, &result.CheckStatus)

			if index <= result.ID {
				index = result.ID
			}

			ch <- result
		}

		time.Sleep(time.Second * 1)
	}
	return
}

package handler

import (
	"TiCheck/insight/server/model"
	"fmt"
	"net/http"
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

type CheckData struct {
	ID          int    `json:"id"`
	CheckTime   string `json:"check_time"`
	CheckClass  string `json:"check_class"`
	CheckName   string `json:"check_name"`
	Operator    string `json:"operator"`
	Threshold   string `json:"threshold"`
	Duration    int    `json:"duration"`
	CheckItem   string `json:"check_item"`
	CheckValue  string `json:"check_value"`
	CheckStatus string `json:"check_status"`
}

func (r *ReportHandler) GetReportByClusterID(c *gin.Context) {

	id := c.Param("clusterID")
	//length, _ := strconv.Atoi(c.Query("length"))
	//start, _ := strconv.Atoi(c.Query("start"))

	ch := &model.CheckHistory{}
	clusterID, err := strconv.Atoi(id)
	if err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	chs, err := ch.GetHistoryByClusterID(clusterID)

	if err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, chs)
	return
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
		ws.WriteJSON(gin.H{
			"error": err.Error(),
		})
		return
	}

	// 获取一个执行时间戳，作为执行时间，同时将其传给执行脚本
	executeTime := time.Now().Unix()

	// 监听脚本是否完成
	executionFinished := make(chan bool)

	// 监听结果获取是否完成
	getResultFinished := make(chan bool)
	go r.executeScript(executeTime, executionFinished)

	// 从数据库中获取实时结果
	resultCh := make(chan *CheckData, 10)
	go r.getResult(executeTime, resultCh, executionFinished, getResultFinished)

	// 设置一分钟的超时时间
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case result := <-resultCh:
			err = ws.WriteJSON(result)
			if err != nil {
				ws.WriteJSON(gin.H{
					"error": err.Error(),
				})
			}
		case <-getResultFinished:
			ws.WriteJSON(gin.H{
				"finish": true,
			})
			return
		case <-ticker.C:
			ws.WriteJSON(gin.H{
				"error": "execute check time out",
			})
			return
		}
	}
}

func (r *ReportHandler) DownloadAllReport(c *gin.Context) {
	return
}

// DownloadReport 下载指定报告
func (r *ReportHandler) DownloadReport(c *gin.Context) {
	c.JSON(http.StatusOK, nil)

	//reportId := c.Param("id")
	//fileName := reportId + ".csv"
	//
	//_, err := os.Open("../report/" + fileName)
	//
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": "the report is not found",
	//	})
	//
	//	return
	//}
	//
	//c.Header("Content-Type", "application/x-xls")
	//c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	//c.File("../report/" + fileName)
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

func (r *ReportHandler) EditConfig(c *gin.Context) {
	script := c.Param("script")
	cmd := exec.Command("../run/" + script + ".sh")
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
}

func (r *ReportHandler) executeScript(executeTime int64, executionFinished chan bool) {
	cmd := exec.Command("../run.sh", strconv.FormatInt(executeTime, 10))
	err := cmd.Run()
	if err != nil {
		print(err)
	}

	// sleep for 10 seconds before sending done signal
	time.Sleep(time.Second * 10)
	executionFinished <- true
}

func (r *ReportHandler) getResult(executeTime int64, ch chan *CheckData, executionFinished chan bool, getResultFinished chan bool) {
	// 记录每一次查询到的最新数据，下一轮查询从这里开始
	var index int

	result := &CheckData{}

	for {
		select {
		case <-executionFinished:
			getResultFinished <- true
			return
		default:
			querySQL := fmt.Sprintf("select * from check_data where check_time == %d and id > %d", executeTime, index)
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
	}
}

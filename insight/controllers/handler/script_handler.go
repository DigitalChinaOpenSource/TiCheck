package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

type ScriptHandler struct {
}

type ScriptList struct {
	Total int		    `json:"total"`
	Scripts []*Script    `json:"script_list"`
}

type Script struct {
	Name string	        `json:"name"`
	Download bool       `json:"download"`
}

// GetAllScript 获取远程仓库脚本列表
func (s *ScriptHandler) GetAllScript(c *gin.Context) {
	url := "https://api.github.com/repos/DigitalChinaOpenSource/TiCheck_ScriptWarehouse/contents/scripts"

	remoteList := make([]string, 0)
	localList := make([]string, 0)

	scriptList := &ScriptList{}

	files, err := ioutil.ReadDir("../script/")
	for _, f := range files {
		name := strings.Split(f.Name(), ".")
		localList = append(localList, name[0])
	}

	jsonMap, err := s.SendRequest(url)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
	}

	for i, _ := range jsonMap{
		switch data := jsonMap[i]["name"].(type) {
		case string:
			script := &Script{}
			remoteList = append(remoteList, data)
			isDownload := false
			for _, v := range localList {
				if data == v {
					isDownload = true
					break
				}
			}

			script.Name = data
			script.Download = isDownload

			scriptList.Total += 1
			scriptList.Scripts = append(scriptList.Scripts, script)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "can't find script for remote warehouse, please check whether the remote warehouse is valid : " + url,
			})
			return
		}

	}

	c.JSON(http.StatusOK, scriptList)
	return
}

// GetReadMe 获取远程仓库某个脚本的 Readme 文件并返回
func (s *ScriptHandler) GetReadMe(c *gin.Context) {

}

// DownloadScript 下载远程仓库脚本到本地
func (s *ScriptHandler) DownloadScript(c *gin.Context) {

}

func (s *ScriptHandler) SendRequest(url string) ([]map[string]interface{} ,error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err !=nil {
		return nil, err
	}

	jsonMap := make([]map[string]interface{}, 0)

	json.Unmarshal(body, &jsonMap)

	resp.Body.Close()

	return jsonMap, nil
}

package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ScriptHandler struct {
}

type RemoteScriptList struct {
	Total   int             `json:"total"`
	Scripts []*RemoteScript `json:"script_list"`
}

type RemoteScript struct {
	Name     string `json:"name"`
	Download bool   `json:"download"`
}

type LocalScriptList struct {
	Total   int            `json:"total"`
	Scripts []*LocalScript `json:"script_list"`
}

type LocalScript struct {
	Name string `json:"name"`
}

// GetAllLocalScript 获取本地所有脚本列表
func (s *ScriptHandler) GetAllLocalScript(c *gin.Context) {
	start := c.Query("start")
	length := c.Query("length")

	st, err := strconv.Atoi(start)
	if err != nil {
		st = 0
	}

	le, err := strconv.Atoi(length)
	if err != nil && le == 0 {
		le = 10
	}

	localList := &LocalScriptList{}
	files, err := ioutil.ReadDir("../script/")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if st > len(files) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "exceeded the maximum number of files",
		})
		return
	}

	for _, f := range files {
		name := strings.Split(f.Name(), ".")
		localList.Scripts = append(localList.Scripts, &LocalScript{name[0]})
		localList.Total += 1
	}

	if st+le > len(files) {
		le = len(files) - st
	}

	localList.Scripts = localList.Scripts[st : le+st]

	c.JSON(http.StatusOK, localList)
	return
}

// GetAllRemoteScript 获取远程仓库脚本列表
func (s *ScriptHandler) GetAllRemoteScript(c *gin.Context) {
	start := c.Query("start")
	length := c.Query("length")

	st, err := strconv.Atoi(start)
	if err != nil {
		st = 0
	}

	le, err := strconv.Atoi(length)
	if err != nil && le == 0 {
		le = 10
	}

	url := "https://api.github.com/repos/DigitalChinaOpenSource/TiCheck_ScriptWarehouse/contents/scripts"

	remoteList := make([]string, 0)
	localList := make([]string, 0)

	scriptList := &RemoteScriptList{}

	files, err := ioutil.ReadDir("../script/")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, f := range files {
		name := strings.Split(f.Name(), ".")
		localList = append(localList, name[0])
	}

	jsonMap, err := s.SendRequest(url)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i, _ := range jsonMap {
		data, ok := jsonMap[i]["name"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "can't find script for remote warehouse, please check whether the remote warehouse is valid : " + url,
			})
			return
		}

		script := &RemoteScript{}
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
	}

	if st+le > len(scriptList.Scripts) {
		le = len(scriptList.Scripts) - st
	}

	if st > len(scriptList.Scripts) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "exceeded the maximum number of scripts",
		})
		return
	}

	scriptList.Scripts = scriptList.Scripts[st : le+st]

	c.JSON(http.StatusOK, scriptList)
	return
}

// GetReadMe 获取远程仓库某个脚本的 Readme 文件并返回
func (s *ScriptHandler) GetReadMe(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "script name not specified",
		})
		return
	}

	url := "https://raw.githubusercontent.com/DigitalChinaOpenSource/TiCheck_ScriptWarehouse/main/scripts/" + name + "/readme.md"
	//url1 := "https://raw.githubusercontent.com/DigitalChinaOpenSource/TiCheck_ScriptWarehouse/main/scripts/" + name + "/"+ name + ".config"

	resp, err := http.Get(url)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to access remote warehouse: " + err.Error(),
		})
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get readme: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"readme": string(body),
	})

	return
}

// DownloadScript 下载远程仓库脚本到本地
func (s *ScriptHandler) DownloadScript(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "script name not specified",
		})
		return
	}

	isExist, err := s.CheckScriptIsExist(name)
	if isExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the script already exists locally",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	scriptFileUrl := "https://api.github.com/repos/DigitalChinaOpenSource/TiCheck_ScriptWarehouse/contents/scripts/" + name

	// 由于可能存在不同的脚本语言不同，脚本名的后缀可能是 .py 或者 .sh, 需要获取脚本名的全文（包括后缀）
	// 但规定该文件下只有三个文件，分表是 script_name.py || script_name.sh, script_name.config 和 readme.md

	jsonMap, err := s.SendRequest(scriptFileUrl)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var fileName string

	for i, _ := range jsonMap {
		value, ok := jsonMap[i]["name"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "can't find script for remote warehouse, please check whether the remote warehouse is valid : " + scriptFileUrl,
			})
			return
		}

		if value != name+".config" && value != "readme.md" {
			fileName = value
			break
		} else {
			if i == len(jsonMap)-1 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "can't find script for remote warehouse, please check whether the remote warehouse is valid : " + scriptFileUrl,
				})
				return
			}
		}
	}

	configUrl := "https://raw.githubusercontent.com/DigitalChinaOpenSource/TiCheck_ScriptWarehouse/main/scripts/" + name + "/" + name + ".config"

	scriptUrl := "https://raw.githubusercontent.com/DigitalChinaOpenSource/TiCheck_ScriptWarehouse/main/scripts/" + name + "/" + fileName

	if err := s.saveScript(scriptUrl, fileName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := s.updateConfig(configUrl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return
}

func (s *ScriptHandler) SendRequest(url string) ([]map[string]interface{}, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	resp.Body.Close()

	jsonMap := make([]map[string]interface{}, 0)

	json.Unmarshal(body, &jsonMap)

	return jsonMap, nil
}

// CheckScriptIsExist 检查改脚本本地是否已经存在
func (s *ScriptHandler) CheckScriptIsExist(name string) (bool, error) {
	files, err := ioutil.ReadDir("../script/")
	if err != nil {
		return false, err
	}

	for _, f := range files {
		localName := strings.Split(f.Name(), ".")
		if localName[0] == name {
			return true, nil
		}
	}

	return false, nil
}

// saveSctipt 通过 url 下载脚本文件保存到本地
// 本地文件保存目录： ../script/
func (s *ScriptHandler) saveScript(scriptUrl string, scriptName string) error {
	resp, err := http.Get(scriptUrl)

	if err != nil {
		return err
	}

	out, err := os.Create("../script/" + scriptName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// updateConfig 更新配置文件
// 配置文件目录：
func (s *ScriptHandler) updateConfig(configUrl string) error {
	resp, err := http.Get(configUrl)

	if err != nil {
		return err
	}

	f, err := os.OpenFile("../config/execution_config.csv", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// 查找文件末尾的偏移量
	n, _ := f.Seek(0, os.SEEK_END)
	// 从末尾的偏移量开始写入内容
	_, err = f.WriteAt(body, n)

	defer f.Close()

	if err != nil {
		return err
	}

	return nil
}

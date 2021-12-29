package handler

import (
	"bytes"
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"time"
)

type SessionHandler struct{
	user        string
	password    string
	token       string
}



func (s *SessionHandler)AuthenticatedUser(c *gin.Context){
	//cookie, _ := c.Cookie("TiCheckerToken")
	//if cookie != "" && cookie == s.token {
	//	c.JSON(http.StatusOK, gin.H{
	//		"token": s.token,
	//	})
	//}
	//
	//data, _ := ioutil.ReadAll(c.Request.Body)
	//jsonStr := string(data)
	//var jsonMap map[string]interface{}
	//
	//if err := json.Unmarshal([]byte(jsonStr), &jsonMap); err != nil {
	//	c.JSON(http.StatusOK, map[string]interface{}{
	//		"error":           err.Error(),
	//	})
	//}
	//
	//s.user = jsonMap["user"].(string)
	//s.password = jsonMap["password"].(string)

	s.user = c.PostForm("username")
	s.password = c.PostForm("password")

	if s.verifyDBUser() {
		http.SetCookie(c.Writer, &http.Cookie{
			Name: "TiCheckerToken",
			Value: s.token,
		})

		c.JSON(http.StatusOK, gin.H{
			"token": s.token,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "the user name or password is wrong.",
	})

	return
}

func (s *SessionHandler) verifyDBUser() bool {
	// 用户密码 目前先写死，后期用TiDB用户做验证
	if s.user == "tidb" && s.password == "password123" {
		if s.token == "" {
			s.UpdateToken()
		}

		return true
	}

	return false
}

// UpdateToken 更新验证Token, 之后每半小时自动更新一次
func (s *SessionHandler) UpdateToken() {
	s.CreateToken(64)

	go func() {
		for range time.Tick(time.Second*1800){
			s.CreateToken(64)
		}
	}()
}

func (s *SessionHandler) CreateToken(len int) {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0;i < len ;i++  {
		randomInt,_ := rand.Int(rand.Reader,bigInt)
		container += string(str[randomInt.Int64()])
	}
	s.token = container
}

func (s *SessionHandler) VerifyToken(c *gin.Context) {

	token, err := c.Cookie("TiCheckerToken")

	if err == nil && token == s.token {
		c.Next()
	}

	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

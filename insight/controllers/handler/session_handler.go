package handler

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type SessionHandler struct {
	user     string
	password string
	token    string
}

type Session struct {
	user     string
	password string
	token    string
}

func (s *SessionHandler) AuthenticatedUser(c *gin.Context) {

	s.user = c.PostForm("username")
	s.password = c.PostForm("password")

	if s.verifyDBUser() {
		c.SetCookie("TiCheckerToken", s.token, 3600, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"token": s.token,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "the user name or password is wrong.",
	})

}

func (s *SessionHandler) verifyDBUser() bool {
	// 用ping做账号密码验证
	conn := fmt.Sprintf("%s:%s@tcp(10.3.65.140:4000)/mysql", s.user, s.password)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return false
	}
	err = db.Ping()
	if err != nil {
		return false
	}

	defer db.Close()

	if s.token == "" {
		s.UpdateToken()
	}
	return true
}

// UpdateToken 更新验证Token, 之后每半小时自动更新一次
func (s *SessionHandler) UpdateToken() {
	s.CreateToken(64)

	go func() {
		for range time.Tick(time.Second * 1800) {
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
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	s.token = container
}

func (s *SessionHandler) VerifyToken(c *gin.Context) {

	token, err := c.Cookie("TiCheckerToken")

	if err == nil && token == s.token {
		c.Next()
		return
	}

	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

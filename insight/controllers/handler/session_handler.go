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
	Sessions map[string]*Session
}

type Session struct {
	user     string
	password string
	token    string
	ticker   *time.Ticker
}

func (sh *SessionHandler) AuthenticatedUser(c *gin.Context) {

	user := c.PostForm("username")
	password := c.PostForm("password")
	se := &Session{
		user: user,
		password: password,
	}
	if sh.UserIsExit(user) {
		se = sh.Sessions[user]
	} else {
		se = sh.CreateUser(user, password)
	}

	if se.verifyDBUser() {
		c.SetCookie("TiCheckerToken", se.token, 3600, "/", "", false, true)
		c.SetCookie("TiCheckerUser", se.user, 3600, "/", "", false, false)

		c.JSON(http.StatusOK, gin.H{
			"token": se.token,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "the user name or password is wrong.",
	})

}

func (se *Session) verifyDBUser() bool {
	// 用ping做账号密码验证
	conn := fmt.Sprintf("%s:%s@tcp(10.3.65.140:4000)/mysql", se.user, se.password)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return false
	}
	err = db.Ping()
	if err != nil {
		return false
	}

	defer db.Close()

	se.UpdateToken()

	return true
}

// UpdateToken 更新验证Token, 同时重置token更新时间
func (se *Session) UpdateToken() {
	se.CreateToken(64)
	se.ticker.Reset(time.Minute * 30)
}

func (se *Session) CreateToken(len int) {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	se.token = container
}

func (sh *SessionHandler) VerifyToken(c *gin.Context) {

	user, err0 := c.Cookie("TiCheckerUser")
	token, err1 := c.Cookie("TiCheckerToken")

	if err0 == nil && err1 ==nil && token == sh.Sessions[user].token {
		c.Next()
		return
	}

	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

// UserIsExit 判断用户是否已经存在
func (sh *SessionHandler) UserIsExit(user string) bool {
	if _,ok := sh.Sessions[user]; ok {
		return true
	}

	return false
}

// CreateUser 创建一个 user 并生成相应 token
func (sh *SessionHandler) CreateUser(user string, password string) *Session {
	se := &Session{
		user: user,
		password: password,
	}
	se.CreateToken(64)
	se.ticker = time.NewTicker(time.Minute * 30)
	sh.Sessions[user] = se
	return se
}

func (sh *SessionHandler) Logout(c *gin.Context) {
	user, err := c.Cookie("TiCheckerUser")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}

	if _,ok := sh.Sessions[user]; ok {
		delete(sh.Sessions, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return
}
package handler

import (
	"TiCheck/internal/model"
	"bytes"
	"crypto/rand"
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
	User   *model.User
	Token  string
	Ticker *time.Ticker
}

type UserReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type GetUserInfoResp struct {
	UserName string `json:"user_name,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (sh *SessionHandler) AuthenticatedUser(c *gin.Context) {
	userReq := &UserReq{}
	err := c.BindJSON(userReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "the request body is wrong",
		})

		return
	}

	se := &Session{
		User: &model.User{
			UserName:     userReq.UserName,
			UserPassword: userReq.Password,
		},
	}

	if !se.User.VerifyUser() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "the User name or password is wrong.",
		})

		return
	}

	if sh.UserIsExit(userReq.UserName) {
		se = sh.Sessions[userReq.UserName]
	} else {
		sh.CreateUser(se)
	}

	c.SetCookie("TiCheckerToken", se.Token, 3600, "/", "", false, true)
	c.SetCookie("TiCheckerUser", userReq.UserName, 3600, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{
		"token": se.Token,
	})
	return

}

func (sh *SessionHandler) Logout(c *gin.Context) {
	user, err := c.Cookie("TiCheckerUser")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}

	if _, ok := sh.Sessions[user]; ok {
		sh.Sessions[user].Ticker.Reset(time.Millisecond)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return
}

func (sh *SessionHandler) GetUserInfo(c *gin.Context) {
	user, err0 := c.Cookie("TiCheckerUser")
	token, err1 := c.Cookie("TiCheckerToken")

	if err0 == nil && err1 == nil && sh.UserIsExit(user) {
		if token == sh.Sessions[user].Token {
			userInfo := model.User{
				UserName: user,
			}
			if err := userInfo.GetUserInfoByName(); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, GetUserInfoResp{
				UserName: userInfo.UserName,
				FullName: userInfo.FullName,
				Email:    userInfo.Email,
			})

			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

// UpdateToken 更新Token, 同时重置token过期时间
func (se *Session) UpdateToken() {
	se.CreateToken(64)
	se.Ticker.Reset(time.Minute * 30)
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
	se.Token = container
}

func (sh *SessionHandler) VerifyToken(c *gin.Context) {

	user, err0 := c.Cookie("TiCheckerUser")
	token, err1 := c.Cookie("TiCheckerToken")

	if err0 == nil && err1 == nil && sh.UserIsExit(user) {
		if token == sh.Sessions[user].Token {
			c.Next()
			return
		}
	}

	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
	return
}

// UserIsExit 判断用户是否已经存在
func (sh *SessionHandler) UserIsExit(user string) bool {
	if _, ok := sh.Sessions[user]; ok {
		return true
	}

	return false
}

// CreateUser 创建一个 User 并生成相应 token
func (sh *SessionHandler) CreateUser(se *Session) {
	se.CreateToken(64)
	se.Ticker = time.NewTicker(time.Minute * 30)
	userName := se.User.UserName
	sh.Sessions[userName] = se
	go func() {
		<-se.Ticker.C
		se.Ticker.Stop()
		delete(sh.Sessions, userName)
		return
	}()

	return
}

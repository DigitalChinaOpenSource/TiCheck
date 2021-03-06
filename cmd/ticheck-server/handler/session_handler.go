package handler

import (
	"TiCheck/cmd/ticheck-server/api"
	"TiCheck/internal/model"
	"TiCheck/util/logutil"
	"bytes"
	"crypto/rand"
	"go.uber.org/zap"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const DefaultTokenLife = time.Hour

type SessionHandler struct {
	// token : user_name
	Users map[string]string
	// user_name : session
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

var SessionHelper = SessionHandler{
	Users:    map[string]string{},
	Sessions: make(map[string]*Session, 0),
}

func (sh *SessionHandler) AuthenticatedUser(c *gin.Context) {
	userReq := &UserReq{}
	err := c.BindJSON(userReq)
	if err != nil {
		logutil.Logger.Error("the request body can't be parsed correctly", zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	se := &Session{
		User: &model.User{
			UserName:     userReq.UserName,
			UserPassword: userReq.Password,
		},
	}

	if !se.User.VerifyUser() {
		logutil.Logger.Error("authentication failed", zap.String("user", se.User.UserName))
		api.AuthenticationFailed(c)
		return
	}

	if sh.UserIsExit(userReq.UserName) {
		se = sh.Sessions[userReq.UserName]
		delete(sh.Users, se.Token)
		se.UpdateToken()
		sh.Users[se.Token] = se.User.UserName
	} else {
		sh.CreateSession(se)
	}

	//c.SetCookie("TiCheckerToken", se.Token, 3600, "/", "", false, true)
	//c.SetCookie("TiCheckerUser", userReq.UserName, 3600, "/", "", false, false)

	api.Success(c, "", map[string]string{
		"token": se.Token,
	})
	return
}

func (sh *SessionHandler) Logout(c *gin.Context) {
	defer func() {
		api.S(c)
	}()

	//user, _ := c.Cookie("TiCheckerUser")
	tokens, ok := c.Request.Header["Access-Token"]
	if !ok || len(tokens) < 1 {
		return
	}

	se := sh.getSessionByToken(tokens[0])
	if se == nil {
		return
	}

	// reset ticker and let the token expire directly
	se.Ticker.Reset(time.Millisecond)

	return
}

func (sh *SessionHandler) GetUserInfo(c *gin.Context) {
	//user, err0 := c.Cookie("TiCheckerUser")
	//token, err1 := c.Cookie("TiCheckerToken")

	token, ok := c.Request.Header["Access-Token"]
	if !ok || len(token) < 1 {
		logutil.Logger.Debug("can't get token from request header.")
		api.BadWithMsg(c, "the token is invalid.")
		return
	}

	se := sh.getSessionByToken(token[0])
	if se == nil {
		logutil.Logger.Debug("the token is invalid.")
		api.BadWithMsg(c, "the token is invalid.")
		return
	}

	err := se.User.GetUserInfoByName()
	if err != nil {
		logutil.Logger.Error("user does not exist.", zap.Error(err))
		api.BadWithMsg(c, "user does not exist.")
		return
	}

	api.Success(c, "", map[string]string{
		"user_name": se.User.UserName,
		"email":     se.User.Email,
		"full_name": se.User.FullName,
	})

	return
}

// UpdateToken ??????Token, ????????????token????????????
func (se *Session) UpdateToken() {
	se.CreateToken(64)
	se.Ticker.Reset(DefaultTokenLife)
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
	//user, err0 := c.Cookie("TiCheckerUser")
	//token, err1 := c.Cookie("TiCheckerToken")
	token, ok := c.Request.Header["Access-Token"]
	if ok && len(token) > 0 {
		if se := sh.getSessionByToken(token[0]); se != nil {
			se.Ticker.Reset(DefaultTokenLife)
			c.Next()
			return
		}
	}

	c.Abort()
	api.AuthenticationFailed(c)
	return
}

// UserIsExit ??????????????????????????????
func (sh *SessionHandler) UserIsExit(user string) bool {
	if _, ok := sh.Sessions[user]; ok {
		return true
	}

	return false
}

// CreateSession ???????????? Session ??????????????? token
// Session will remain for 60 minutes
func (sh *SessionHandler) CreateSession(se *Session) {
	se.CreateToken(64)
	se.Ticker = time.NewTicker(DefaultTokenLife)
	userName := se.User.UserName
	sh.Users[se.Token] = userName
	sh.Sessions[userName] = se

	go sh.clearSession(se)

	return
}

func (sh *SessionHandler) clearSession(se *Session) {
	<-se.Ticker.C
	se.Ticker.Stop()
	logutil.Logger.Info("delete a session and token", zap.String("user", se.User.UserName))
	delete(sh.Sessions, se.User.UserName)
	delete(sh.Users, se.Token)
	return
}

func (sh *SessionHandler) getSessionByToken(token string) *Session {
	if token == "" {
		return nil
	}
	userName, ok := sh.Users[token]
	if !ok {
		return nil
	}

	return sh.Sessions[userName]
}

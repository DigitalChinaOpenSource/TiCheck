package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct{}

func (a SessionHandler) AuthenticatedUser(c *gin.Context) {

	// data, _ := ioutil.ReadAll(c.Request.Body)
	// jsonStr := string(data)
	// var jsonMap map[string]interface{}

	// if err := json.Unmarshal([]byte(jsonStr), &jsonMap); err != nil {
	// 	c.JSON(http.StatusOK, map[string]interface{}{
	// 		"RespondStatus": false,
	// 		"Err":           err.Error(),
	// 	})
	// }

	// user := jsonMap["username"]
	// password := jsonMap["password"]
	user := c.PostForm("username")
	password := c.PostForm("password")
	if user == "tidb" && password == "password123" {
		c.JSON(http.StatusOK, gin.H{
			"token": "xxx",
		})

		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "the user name or password is wrong.",
	})

}

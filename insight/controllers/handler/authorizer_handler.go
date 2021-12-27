package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type AuthorizerHandler struct{}

func (a AuthorizerHandler)AuthenticatedUser(c *gin.Context){

	data, _ := ioutil.ReadAll(c.Request.Body)
	jsonStr := string(data)
	var jsonMap map[string]interface{}

	if err := json.Unmarshal([]byte(jsonStr), &jsonMap); err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"RespondStatus": false,
			"Err":           err.Error(),
		})
	}

	user := jsonMap["user"]
	password := jsonMap["password"]

	if user == "root" && password == "password123" {
		c.JSON(http.StatusOK, gin.H{
			"Data": "OK",
			"Error": nil,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Data": nil,
		"Error": "the user name or password is wrong.",
	})

	return
}
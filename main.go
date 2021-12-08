package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bastianpahlke/go-test.git/types"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)
import "net/http"

func main() {
	r := gin.Default()

	r.GET("/userData/:userId", getDataForUser())

	r.Run()
}

func getDataForUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		userChan := make(chan *http.Response)
		dataChan := make(chan *http.Response)

		get := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%s", userId)
		go SendGetAsync(get, userChan)
		body, _ := json.Marshal("")
		post := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts?userId=%s", userId)
		go SendPostAsync(post, body, dataChan)

		userResponse := <- userChan
		defer userResponse.Body.Close()
		userBytes, _ := ioutil.ReadAll(userResponse.Body)
		var user types.User
		json.Unmarshal(userBytes, &user)

		dataResponse := <- dataChan
		defer dataResponse.Body.Close()
		dataBytes, _ := ioutil.ReadAll(dataResponse.Body)
		var posts []types.HtmlData
		json.Unmarshal(dataBytes, &posts)

		size := len(posts)

		var htmlData types.HtmlData
		if size > 0 {
			htmlData = posts[size - 1]
		} else {
			htmlData = types.HtmlData {
				Id: 0,
				Title: "Dummy Title",
				Body: "Dummy",
			}
		}

		c.String(http.StatusOK, "Data for user %s: \n Username: %s \n Email: %s \n \n Last Post: %s #%d \n %s", userId, user.Username,
			user.Email, htmlData.Title, htmlData.Id, htmlData.Body)
	}
}

func SendPostAsync(url string, body []byte, rc chan *http.Response) {
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	rc <- response
}

func SendGetAsync(url string, rc chan *http.Response) {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	rc <- response
}

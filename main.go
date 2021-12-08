package main

import (
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
		post := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts?userId=%s", userId)
		go SendGetAsync(post, dataChan)

		userResponse := <- userChan
		defer userResponse.Body.Close()
		userBytes, _ := ioutil.ReadAll(userResponse.Body)
		var user types.User
		json.Unmarshal(userBytes, &user)

		dataResponse := <- dataChan
		defer dataResponse.Body.Close()
		dataBytes, err := ioutil.ReadAll(dataResponse.Body)

		if err != nil {

			c.String(http.StatusOK, "Data for user %s: \n Username: %s \n Email: %s \n \n Last Post: %s #%d \n %s",
				userId, user.Username, user.Email, "Dummy title", 0, "No body")
		} else {
			var posts []types.HtmlData
			json.Unmarshal(dataBytes, &posts)

			size := len(posts)

			var htmlData types.HtmlData
			if size > 0 {
				htmlData = posts[size - 1]
			} else {
				htmlData = types.HtmlData {
					Id: 0,
					Title: fmt.Sprintf("raw: (%s)\n", string(dataBytes)),
					Body: "If you see this, something went wrong",
				}
			}

			c.String(http.StatusOK, "Data for user %s: \nUsername: %s \nEmail: %s \n \nLast Post: %s (#%d) \n\n%s",
				userId, user.Username, user.Email, htmlData.Title, htmlData.Id, htmlData.Body)
		}
	}
}

func SendGetAsync(url string, rc chan *http.Response) {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	rc <- response
}

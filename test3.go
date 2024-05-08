package main

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
)

var tools = make(map[string]interface{})
var jsonData = make(map[string]interface{})

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func handleRequestAndRedirect(c *gin.Context) {
	// We will get to this...
	res := c.Writer
	req := c.Request
	id := c.Param("id")

	adv, _ := url.Parse("http://127.0.0.1:1005/")
	req.URL, _ = url.Parse(id)
	// fmt.Println(req)
	httputil.NewSingleHostReverseProxy(adv).ServeHTTP(res, req)

}

func main() {
	r := gin.Default()

	r.GET("/advance", handleRequestAndRedirect)
	r.GET("/advance/:id", handleRequestAndRedirect)
	r.Use(static.Serve("/", static.LocalFile("./frontend/", false)))

	r.Run(":80")
}

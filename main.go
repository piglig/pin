package main

import (
	"net/http"
	"pin"
)

func main() {
	r := pin.New()
	r.GET("/", func(c *pin.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Pin</h1>")
	})

	r.POST("/hello", func(c *pin.Context) {
		c.JSON(http.StatusOK, c.Req.Header)
	})

	r.Run(":9999")
}

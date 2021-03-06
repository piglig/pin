package main

import (
	"log"
	"net/http"
	"pin"
	"time"
)

func onlyForV1() pin.HandlerFunc {
	return func(c *pin.Context) {
		t := time.Now()
		// c.Fail(500, "Invalid Server request")

		log.Printf("[%d] %s in %v for group v1", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func onlyDoSomething() pin.HandlerFunc {
	return func(c *pin.Context) {
		log.Printf("[%d] %s content-length %v for group v1", c.StatusCode, c.Req.RequestURI, c.Req.ContentLength)
	}
}

func main() {
	r := pin.New()
	r.GET("/", func(c *pin.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Pin</h1>")
	})

	v1 := r.Group("/v1")
	v1.Use(onlyForV1(), onlyDoSomething())
	{
		v1.GET("/", func(c *pin.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Pin</h1>")
		})

		v1.GET("/hello", func(c *pin.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	r.POST("/hello", func(c *pin.Context) {
		c.JSON(http.StatusOK, c.Req.Header)
	})

	r.GET("/hello/:name", func(c *pin.Context) {
		c.String(http.StatusOK, "hello %s\n", c.Param("name"))
	})

	r.GET("/assets/*filepath", func(c *pin.Context) {
		c.JSON(http.StatusOK, pin.H{
			"filepath": c.Param("filepath"),
		})
	})

	r.Run(":9999")
}

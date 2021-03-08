package main

import (
	"fmt"
	"log"
	"net/http"
	"pin"
	"text/template"
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

type Student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := pin.Default()

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")

	r.Static("/assets", "./static")

	r.GET("/test", func(c *pin.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	stu1 := &Student{Name: "Geektutu", Age: 20}
	stu2 := &Student{Name: "Jack", Age: 22}
	r.GET("/students", func(c *pin.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", pin.H{
			"title":  "pin",
			"stuArr": [2]*Student{stu1, stu2},
		})
	})

	r.GET("/panic", func(c *pin.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.GET("/date", func(c *pin.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", pin.H{
			"title": "pin",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	v1 := r.Group("/v1")
	v1.Use(onlyForV1(), onlyDoSomething())
	{
		v1.GET("/", func(c *pin.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
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

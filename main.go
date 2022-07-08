package main

import (
	"gee-web0703/gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandleFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		//c.String(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
func Logger() gee.HandleFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.NewEngine()
	r.Use(Logger()) // global midlleware
	r.GET("/:name/:test", func(c *gee.Context) {
		c.Json(200, gee.H{
			"name": c.Params["name"],
			"test": c.Params["test"],
		})
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Params["name"], c.Path)
		})
	}
	http.ListenAndServe(":8080", r)
}

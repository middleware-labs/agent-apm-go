~~~~# agent-apm-go

go get github.com/middleware-labs/agent-apm-go


```golang

package main

import (
	"github.com/gin-gonic/gin"
	tracer "github.com/middleware-labs/agent-apm-go/init"
	g "github.com/middleware-labs/agent-apm-go/packages/gin"
	"net/http"
)

func main() {
	go tracer.Execute()
	r := gin.Default()
	r.Use(g.Middleware("service1"))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8070")
}

```
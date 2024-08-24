package goeng

import (
	"github.com/gin-gonic/gin"
)

func Main() {

	service := newService()

	router := gin.Default()

	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	router.GET("/api/dict", service.getDictList)
	router.GET("/api/dict/:id", service.getDict)
	router.POST("/api/dict", service.createDict)

	router.POST("/api/word/:id", service.addWord)

	router.Run(":80")
	//router.RunTLS(":443", "server.pem", "server.key")
}

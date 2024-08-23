package goeng

import (
	"github.com/gin-gonic/gin"
)

func Main() {

	service := newService()

	router := gin.Default()

	router.GET("/api/dict", service.getDictList)
	router.GET("/api/dict/:id", service.getDict)
	router.POST("/api/dict", service.createDict)

	router.POST("/api/word/:id", service.addWord)

	router.Run(":8080")
}

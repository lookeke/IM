package router

import (
	"github.com/gin-gonic/gin"
	api "im/api/http"
	"im/service"
	"im/utils"
	"log"
)

func Server() {
	server := gin.Default()
	server.Use(utils.Cors())
	server.PUT("/register", service.Register)
	server.POST("/login", service.Login)
	server.PUT("/upload", service.Upload)
	server.GET("/chat", service.IM)
	server.GET("/roomList", service.RoomList)
	server.GET("/historyMessage", api.HistoryMessage)
	err := server.Run(":4000")
	if err != nil {
		log.Println("运行gin服务失败,请检查端口是否被占用", err.Error())
	}
}

package user

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup) {
	router.POST("/login", Login)             // 请求路径/user/login
	router.POST("/register", Register)       // 请求路径/user/register
	router.GET("/query", Query)              // 请求路径/user/query
	router.POST("/queryTarget", QueryTarget) // 请求路径/user/queryTarget
	router.PUT("/update", Update)            // 请求路径/user/update
}

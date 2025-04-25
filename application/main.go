package main

import (
	"CipProject/application/middle_ware"
	all_router "CipProject/application/router"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	engine.Use(
		middle_ware.CorsMiddleWare, // 1 解决跨域请求问题
		middle_ware.AuthMiddleWare, // 2 中间件，拦截未登录用户（Jwt Token实现）
	)

	// 调用router.go中方法实现路由分组
	all_router.InitRouter(engine)

	err := engine.Run(":8088")
	if err != nil {
		fmt.Println("应用程序启动失败:", err.Error())
	}
}

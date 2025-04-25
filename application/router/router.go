package all_router

import (
	"CipProject/application/controller/adminstrator"
	"CipProject/application/controller/code"
	"CipProject/application/controller/product"
	"CipProject/application/controller/subscribe"
	"CipProject/application/controller/user"
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) { //路由分组
	authorizeGroup := router.Group("/code")      // 代码项目操作请求
	subscribeGroup := router.Group("/subscribe") // 租赁请求
	userGroup := router.Group("/user")           // 普通账号操作请求
	administratorGroup := router.Group("/admin") // 组织管理员操作请求
	serviceGroup := router.Group("/product")     // 软件产品管理请求

	code.Router(authorizeGroup) // 调用controller中router
	subscribe.Router(subscribeGroup)
	user.Router(userGroup)
	adminstrator.Router(administratorGroup)
	product.Router(serviceGroup)
}

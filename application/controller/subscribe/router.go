package subscribe

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup) {
	router.POST("/new", New)                           // 请求路径/subscribe/new
	router.PUT("/cancel", Cancel)                      // 请求路径/subscribe/cancel
	router.POST("/review", Review)                     // 请求路径/subscribe/review
	router.POST("/confirm", Confirm)                   // 请求路径/subscribe/confirm
	router.GET("/queryByPartyA", QueryByPartyA)        // 请求路径/subscribe/queryByPartyA
	router.GET("/queryByPartyB", QueryByPartyB)        // 请求路径/subscribe/queryByPartyB
	router.POST("/query", Query)                       // 请求路径/subscribe/query
	router.PUT("/terminate", Terminate)                // 请求路径/subscribe/terminate
	router.PUT("/pay", Pay)                            // 请求路径/subscribe/pay
	router.POST("/addCart", AddCart)                   // 请求路径/subscribe/addCart
	router.GET("/queryCart", QueryCart)                // 请求路径/subscribe/queryCart
	router.POST("/queryCartDetails", QueryCartDetails) // 请求路径/subscribe/queryCartDetails
	router.DELETE("/removeCart", RemoveCart)           // 请求路径/subscribe/removeCart
}

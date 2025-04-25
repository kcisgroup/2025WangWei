package product

import "github.com/gin-gonic/gin"

func Router(r *gin.RouterGroup) {
	r.POST("/create", Create)                    // 请求路径是/product/create
	r.PUT("/start", Start)                       // 请求路径是/product/start
	r.PUT("/stop", Stop)                         // 请求路径是/product/stop
	r.DELETE("/destroy", Destroy)                // 请求路径是/product/destroy
	r.POST("/query", QueryById)                  // 请求路径是/product/query
	r.GET("/queryByOwner", QueryByOwner)         // 请求路径是/product/queryByOwner
	r.GET("/queryRaw", QueryRaw)                 // 请求路径是/product/queryRaw
	r.GET("/queryDestroyList", QueryDestroyList) // 请求路径是/product/queryDestroyList
}

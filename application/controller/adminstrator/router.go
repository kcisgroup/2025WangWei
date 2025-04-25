package adminstrator

import "github.com/gin-gonic/gin"

func Router(router *gin.RouterGroup) {
	router.POST("/putOnRecord", PutOnRecord) // 管理员备案员工信息，请求路径/admin/putOnRecord
}

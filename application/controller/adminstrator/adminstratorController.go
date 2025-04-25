package adminstrator

import (
	"CipProject/application/domain"
	"CipProject/application/service/userService"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PutOnRecord 组织管理员备案组织内的普通用户信息（上链）
func PutOnRecord(ctx *gin.Context) {
	// 接收(Role,Username)
	var user domain.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.Response{Msg: err.Error()})
		return
	}

	// 备案信息上链
	err = userService.PutRecord(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.Response{Msg: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "备案成功"})
}

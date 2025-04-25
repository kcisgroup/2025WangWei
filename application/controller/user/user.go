package user

import (
	"CipProject/application/common"
	"CipProject/application/domain"
	"CipProject/application/service/userService"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login 组织内普通用户登录
func Login(ctx *gin.Context) {
	var loginUser domain.User
	err := ctx.ShouldBind(&loginUser)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "参数接收失败，请重试"})
		return
	}

	// 验证用户名和密码是否正确
	// 1 先查区块链账本
	var user domain.User
	user, errGetUser := userService.GetUser(loginUser.Username)

	// 2 作判断
	// 2.1 验证不通过
	if errGetUser != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "该用户名不存在"})
		return
	}
	if user.Pwd != loginUser.Pwd {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "密码不正确"})
		return
	}

	// 2.2 验证通过，获取Token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "令牌生成失败，请重试"})
		return
	}

	if user.Status == "0" { // 通过Status字段判断登录用户是员工还是管理员，并通过Response的Data字段是否为空来通知前端
		ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "管理员登录成功", Data: "这是管理员", Token: token})
	} else {
		ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "创作者登录成功", Token: token})
	}
}

// Register 组织内普通用户注册（注册前提：组织管理员已将其内部普通用户信息在链上备案）
func Register(ctx *gin.Context) {
	var user domain.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "参数接收失败，请重试"})
		return
	}

	// 只允许员工注册，管理员提前由区块链管理员提供
	// 普通用户注册：“0”公司管理员，“1”普通员工
	user.Status = "1"

	// 注册信息上链
	saveErr := userService.SaveUser(user)
	if saveErr != nil { // 注册失败
		ctx.JSON(http.StatusOK, domain.Response{Msg: saveErr.Error()})
		return
	}
	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "注册成功"})
}

// Query 用户信息查询
func Query(ctx *gin.Context) {
	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: loginUser})
}

// QueryTarget 查询指定用户数据
func QueryTarget(ctx *gin.Context) {
	var user domain.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "参数接收失败，请重试"})
		return
	}

	userInfo, err := userService.GetUser(user.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: userInfo})
}

// Update 账户信息修改
func Update(ctx *gin.Context) {
	var user domain.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	user.Username = loginUser.Username
	err = userService.UpdateUser(user)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "修改成功"})
}

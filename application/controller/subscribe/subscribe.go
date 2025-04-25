package subscribe

import (
	"CipProject/application/common"
	"CipProject/application/domain"
	"CipProject/application/service/cartService"
	"CipProject/application/service/codeService"
	"CipProject/application/service/productService"
	"CipProject/application/service/subscribeService"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"math"
	"net/http"
	"strconv"
)

// New 甲方申请订阅，新增订阅申请表
func New(ctx *gin.Context) {
	// 接收参数（购物车id,code_id,realPrice,startTime,note）
	var sr domain.SubRecords
	err := ctx.ShouldBind(&sr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "参数接收失败，请重试"})
		return
	}

	cartId := sr.ID

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 1 生成订阅申请表ID
	sr.ID = common.NewSnowflake(3).Generate()
	// 2 甲方
	sr.PartyA = loginUser.Username
	// 3 订阅申请表目前状态
	sr.Status = domain.SubToExamine
	// 4 乙方
	cip, err := codeService.GetById(sr.CodeId)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "申请失败：此组件不存在"})
		return
	}
	sr.PartyB = cip.Author
	// 5 代码知识产权的cpi
	sr.Cpi = cip.Cpi
	// 6 申请表创建和更新日期时间
	sr.CreatedTime = common.GetCurrentTime()
	sr.UpdatedTime = sr.CreatedTime

	err = subscribeService.New(sr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "订阅申请单创建失败，请勿重复订阅同一组件"})
		return
	}

	// 订阅申请单创建成功，删除购物车记录
	_ = cartService.RemoveCart(cartId)

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: sr})
}

// Cancel 甲方取消订阅申请
func Cancel(ctx *gin.Context) {
	// 接收参数（订阅表id）
	var sr domain.SubRecords
	err := ctx.ShouldBind(&sr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "参数接收失败，请重试"})
		return
	}
	subInfo, err := subscribeService.GetById(sr.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此订阅申请单不存在"})
		return
	}

	if subInfo.Status == domain.SubToExamine {
		subInfo.Status = domain.SubCancel
		subInfo.UpdatedTime = common.GetCurrentTime()
		_ = subscribeService.Update(subInfo)
		ctx.JSON(http.StatusOK, domain.Response{Code: 1})
	} else {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "取消失败，该申请已被处理"})
		return
	}
}

// Review 乙方审核订阅申请表（"乙方同意，待甲方确认订阅条款"或"乙方已拒绝"）
func Review(ctx *gin.Context) {
	// 接收参数（订阅申请id,realPrice,note,status,startTime）或（订阅申请id,status）
	var sr domain.SubRecords
	err := ctx.ShouldBind(&sr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "参数接收失败，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	subInfo, err := subscribeService.GetById(sr.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此订阅申请不存在，本次处理失败"})
		return
	} else if subInfo.PartyB != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是此交易的乙方，无权处理"})
		return
	} else if subInfo.Status != domain.SubToExamine {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "甲方未申请或订阅已被处理，本次处理失败"})
		return
	}

	if sr.Status == domain.SubToConfirm {
		// 1 双方协商之后软件订阅价格
		subInfo.RealPrice = sr.RealPrice
		// 2 乙方批注
		subInfo.Note = sr.Note
		// 3 订阅申请表审核结果
		subInfo.Status = sr.Status
		// 4 修改订阅生效时间
		subInfo.StartTime = sr.StartTime
		// 5 更新时间
		subInfo.UpdatedTime = common.GetCurrentTime()
	} else if sr.Status == domain.SubRefusedPartyB {
		subInfo.Status = sr.Status
	}

	err = subscribeService.Update(subInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "审核失败，请重试"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1})
}

// Confirm 甲方确认（或拒绝）订阅条款
func Confirm(ctx *gin.Context) {
	// 接收参数（订阅申请id,status）
	var sr domain.SubRecords
	err := ctx.ShouldBind(&sr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "参数接收失败，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	subInfo, err := subscribeService.GetById(sr.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此订阅申请不存在，本次确认失败"})
		return
	} else if subInfo.PartyA != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是此交易甲方，无权审核"})
		return
	} else if subInfo.Status != domain.SubToConfirm {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "乙方未同意或已签订合同"})
		return
	}

	// 更新订阅申请表信息
	subInfo.Status = sr.Status
	subInfo.UpdatedTime = common.GetCurrentTime() // 更新时间

	if subInfo.Status == domain.SubRefusedPartyA {
		err = subscribeService.Update(subInfo)
		if err != nil {
			ctx.JSON(http.StatusOK, domain.Response{Msg: "确认失败，请重试"})
			return
		}
	} else if subInfo.Status == domain.SubApproved {
		var contract domain.SubContract

		// 处理待上链订阅信息
		_ = copier.Copy(&contract, &subInfo)
		contract.ID = strconv.FormatInt(subInfo.ID, 10)
		contract.RealPrice = strconv.Itoa(subInfo.RealPrice)
		contract.End = "待定"
		contract.Cost = "待定"

		err = subscribeService.Agree(subInfo, contract, loginUser.Role)
		if err != nil {
			ctx.JSON(http.StatusOK, domain.Response{Msg: "确认失败，请重试"})
			return
		}
	} else {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "确认失败，请重试"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1})
}

// QueryByPartyA 甲方查看历史订阅记录（使用者）
func QueryByPartyA(ctx *gin.Context) {
	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	subRecords, err := subscribeService.GetByPartyA(loginUser.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: subRecords})
}

// QueryByPartyB 乙方查看处理的订单（创作者）
func QueryByPartyB(ctx *gin.Context) {
	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	subs, err := subscribeService.GetByPartyB(loginUser.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "订单查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: subs})
}

// Query 查看订单详情
func Query(ctx *gin.Context) {
	// 接收参数（订阅记录id）
	var sr domain.SubRecords
	err := ctx.ShouldBind(&sr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "参数接收失败，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 查询得到订阅记录信息
	subInfo, err := subscribeService.GetById(sr.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您查询的订阅信息不存在"})
		return
	} else if subInfo.PartyA != loginUser.Username && subInfo.PartyB != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您未参与此交易，无法查看"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: subInfo})
}

// Terminate 终止订阅，结算费用
func Terminate(ctx *gin.Context) {
	// 接收参数（dockerId）
	var sr domain.SubRecords
	_ = ctx.ShouldBind(&sr)

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 确保服务已激活且已停止
	var service domain.DockerInfo
	if sr.DockerId == 0 { // 尚未使用组件
		ctx.JSON(http.StatusOK, domain.Response{Msg: "未激活，不可结束订阅"})
		return
	} else {
		service, _ = productService.GetById(sr.DockerId) // 检查服务状态
		if service.Status == domain.DockerRunning || service.Status == domain.DockerDestroyed {
			ctx.JSON(http.StatusOK, domain.Response{Msg: "请确保服务已停止"})
			return
		}
	}

	subInfo, err := subscribeService.GetByDockerId(sr.DockerId)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "订阅记录不存在"})
		return
	}

	// 检查登录用户身份是否为合法使用者
	if subInfo.PartyA != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您未订阅此产权，无权操作"})
		return
	}

	// 填充终止订阅信息
	subInfo.Cost = math.Round(subInfo.Term*float64(subInfo.RealPrice)*100) / 100
	subInfo.Status = domain.SubWithoutPay
	subInfo.EndTime = common.GetCurrentTime() // 订阅终止时间
	subInfo.UpdatedTime = subInfo.EndTime

	// 处理待上链合同信息
	var contract domain.SubContract
	_ = copier.Copy(&contract, &subInfo)
	contract.ID = strconv.FormatInt(subInfo.ID, 10)
	contract.RealPrice = strconv.Itoa(subInfo.RealPrice)
	contract.Cost = strconv.FormatFloat(subInfo.Cost, 'f', 2, 64)

	// 终止订阅，变更状态
	err = subscribeService.Terminate(subInfo, contract, loginUser.Role)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "终止失败，请重试"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "成功终止订阅"})
}

// Pay 支付
func Pay(ctx *gin.Context) {
	// 接收参数（订阅申请id）
	var sr domain.SubRecords
	err := ctx.ShouldBind(&sr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "支付失败：数据丢失，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)
	subRecords, err := subscribeService.GetById(sr.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "支付失败：数据丢失，请重试"})
		return
	} else if loginUser.Username != subRecords.PartyA {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "支付失败：您未参与此交易，无法支付"})
		return
	} else if subRecords.Status != domain.SubWithoutPay {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "支付失败：请先结算费用"})
		return
	}

	// Todo 支付功能
	subRecords.Status = domain.SubPayed
	err = subscribeService.Pay(subRecords)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "支付失败，请重试"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "支付成功"})
}

// AddCart 加入购物车
func AddCart(ctx *gin.Context) {
	// 接收参数（代码产权的codeId）
	var cart domain.Cart
	err := ctx.ShouldBind(&cart)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "加入购物车失败:数据丢失，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 设置购物车信息
	cart.Owner = loginUser.Username
	cart.CreatedTime = common.GetCurrentTime()
	cart.UpdatedTime = cart.CreatedTime
	cart.ID = common.NewSnowflake(4).Generate() // 雪花算法生成ID

	err = cartService.Add(cart)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "请不要将此组件重复加入购物车"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: cart})
}

// QueryCart 查看购物车列表
func QueryCart(ctx *gin.Context) {
	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	carts, err := cartService.GetCartByUsername(loginUser.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "购物车查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: carts})
}

// QueryCartDetails 查看购物车列详情
func QueryCartDetails(ctx *gin.Context) {
	var cart domain.Cart
	err := ctx.ShouldBind(&cart)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "接收参数失败"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	cartInfo, err := cartService.GetById(cart.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "购物车查询失败"})
		return
	} else if cartInfo.Owner != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "这不是您的购物车数据，无权操作"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: cartInfo})
}

// RemoveCart 移除购物车中指定数据
func RemoveCart(ctx *gin.Context) {
	// 接收参数（id）
	var cart domain.Cart
	err := ctx.ShouldBind(&cart)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "数据丢失，请重试"})
		return
	}

	err = cartService.RemoveCart(cart.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "移除失败"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1})
}

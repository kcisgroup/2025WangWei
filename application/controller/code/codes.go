package code

import (
	"CipProject/application/common"
	"CipProject/application/domain"
	"CipProject/application/ipfs"
	"CipProject/application/service/codeService"
	"CipProject/application/service/userService"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// SaveDraft 新增组件初稿（未上链）
func SaveDraft(ctx *gin.Context) {
	// 接收前端参数（文件同时接收）
	var codeDraft domain.CodeRecords
	err := ctx.ShouldBind(&codeDraft)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "接收数据失败，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 处理临时存储的数据
	// 1 创作者
	codeDraft.Author = loginUser.Username
	// 2 组件创建时间
	codeDraft.CreatedTime = common.GetCurrentTime()
	// 3 组件更新时间
	codeDraft.UpdatedTime = codeDraft.CreatedTime
	// 4 组件状态（无产权状态）
	codeDraft.Status = domain.CodeNoCip
	// 5 组件文件存储地址（链下IPFS系统）
	hash, errIpfs := ipfs.UploadIPFS(codeDraft.DraftFile)
	if errIpfs != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "文件保存失败：" + errIpfs.Error()})
		return
	}
	codeDraft.Addr = hash
	// 6 使用雪花算法生成id（16位）
	snowflake := common.NewSnowflake(1)
	codeDraft.ID = snowflake.Generate()
	// 7 所属公司
	codeDraft.Company = loginUser.Role

	// 组件初稿的元数据存储（链下MySQL数据库）
	errService := codeService.SaveDraft(codeDraft) // 组件名和版本号，确保唯一性
	if errService != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "上传失败，文件名和版本号已被注册"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "上传成功"})
}

// UpdateDraftMetaData 更新组件元数据（未上链）
func UpdateDraftMetaData(ctx *gin.Context) {
	// 接收参数（id,filename,version,language,description）
	var newDraft domain.CodeRecords
	err := ctx.ShouldBind(&newDraft)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "未成功接收更新数据，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 获取组件原始元数据（链下MySQL数据库）
	formerRecords, errGet := codeService.GetById(newDraft.ID)
	if errGet != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此组件未上传，无法更新信息"})
		return
	} else if formerRecords.Author != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是此组件的创作者，无权更新信息"})
		return
	}

	// 处理更新数据
	formerRecords.UpdatedTime = common.GetCurrentTime()
	formerRecords.Filename = newDraft.Filename
	formerRecords.Version = newDraft.Version
	formerRecords.Language = newDraft.Language
	formerRecords.Description = newDraft.Description
	formerRecords.InnerPort = newDraft.InnerPort

	// 更新草稿数据的元数据
	err = codeService.UpdateDraft(formerRecords)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "修改失败，请重试"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "修改成功"})
}

// UpdateDraftFile 更新组件文件（未上链）
func UpdateDraftFile(ctx *gin.Context) {
	// 接收参数（id,draftFile）
	var cr domain.CodeRecords
	err := ctx.ShouldBind(&cr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "未成功接收数据，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 获取组件原始元数据（链下MySQL数据库）
	formerRecords, errGet := codeService.GetById(cr.ID)
	if errGet != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此组件未上传，无法更新文件"})
		return
	} else if formerRecords.Author != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是此组件的创作者，无权更新文件"})
		return
	}

	// 持久化新的文件至IPFS中
	hash, errUploadIPFS := ipfs.UploadIPFS(cr.DraftFile)
	if errUploadIPFS != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "文件上传服务器失败，请重新上传"})
		return
	}

	// 更新代码产权元数据ipfsAddress
	formerRecords.Addr = hash
	formerRecords.UpdatedTime = common.GetCurrentTime()
	err = codeService.UpdateDraft(formerRecords)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "文件更新失败，请重试"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "您的文件上传成功"})
}

// DeleteDraft 删除组件初稿（未上链）
func DeleteDraft(ctx *gin.Context) {
	var draftToRemove domain.CodeRecords
	err := ctx.ShouldBind(&draftToRemove)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 检查登录用户身份权限
	draft, err := codeService.GetById(draftToRemove.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "删除失败，待删除的组件不存在"})
		return
	} else if draft.Author != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是此组件的创作者，无权删除"})
		return
	}
	// 执行删除操作
	errDelete := codeService.DeleteDraftById(draftToRemove.ID)
	if errDelete != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "删除成功"})
}

// QueryDraftBox 草稿箱（未上链）
func QueryDraftBox(ctx *gin.Context) {
	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	cipRecords, err := codeService.QueryDraftBox(loginUser.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: cipRecords})
}

// Register 申请代码产权
func Register(ctx *gin.Context) {
	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 接收参数，组件的（id,price,status）
	var records domain.CodeRecords
	err := ctx.ShouldBind(&records)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "接收参数失败", Data: loginUser})
		return
	}

	// 检查组件和创作者身份权限
	formerRecords, err := codeService.GetById(records.ID) // 获取此组件元数据（链下MySQL数据库）
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此组件不存在，请先上传组件，再申请", Data: loginUser})
		return
	} else if formerRecords.Author != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是此组件创作者，无权申请代码知识产权", Data: loginUser})
		return
	}

	// 处理新数据（MySQL）
	// 1 更新时间
	formerRecords.UpdatedTime = common.GetCurrentTime()
	// 2 设置产权状态（开放或禁用）
	formerRecords.Status = records.Status
	// 3 填充单价
	formerRecords.Price = records.Price
	// 4 计算cpi
	draftFile, errDownload := ipfs.DownloadIPFS(formerRecords.Addr)
	if errDownload != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "草稿文件数据失效，请重新上传", Data: loginUser})
		return
	}
	draftFileURL := filepath.Join(draftFile, formerRecords.Addr)
	readerCipFile, _ := os.Open(draftFileURL)         // 打开组件文件
	cipFileContent, _ := io.ReadAll(readerCipFile)    // 读取组件文件数据
	formerRecords.Cpi = common.SHA256(cipFileContent) // SHA-256单向哈希运算
	defer func() {
		_ = readerCipFile.Close()
		_ = os.RemoveAll(draftFile)
	}()
	// 5 授权时间（即上链时间）
	formerRecords.RegTime = formerRecords.UpdatedTime

	// 处理待上链数据
	var cip domain.CodeIntellectualProperty
	_ = copier.Copy(&cip, &formerRecords) // 复制结构体变量的值
	cip.RegTime = formerRecords.UpdatedTime
	cip.Price = strconv.Itoa(formerRecords.Price)

	// 确权
	err = codeService.Authorize(formerRecords, cip, loginUser.Role)
	if err != nil { // 确权失败
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error(), Data: loginUser})
		return
	}

	// 处理返回数据
	type applyResultDto struct {
		domain.CodeRecords
		domain.User
		FileAddr   string `json:"fileAddr" form:"fileAddr"`     // 组件文件存储地址
		AuthorAddr string `json:"authorAddr" form:"authorAddr"` // 创作者地址
	}
	var result applyResultDto
	_ = copier.Copy(&result, &formerRecords)
	_ = copier.Copy(&result, &loginUser)
	result.FileAddr = formerRecords.Addr
	result.AuthorAddr = loginUser.Addr

	// 确权成功
	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: result})
}

// UpdateCip 更新代码产权元数据（仅允许更改产权元数据中的price、status和UpdatedTime字段）
func UpdateCip(ctx *gin.Context) {
	// 接收参数cpi、price、status三个参数
	var cipForReceive domain.CodeIntellectualProperty
	err := ctx.ShouldBind(&cipForReceive)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "接收参数失败，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 检查登录用户身份权限
	cip, err := codeService.GetFromChain(cipForReceive.Cpi, loginUser.Role) // 获取该代码产权的元数据（链上获取）
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此组件未授权，修改失败"})
		return
	} else if cip.Author != loginUser.Username { // 不是创作者本人
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是此组件的创作者，修改失败"})
		return
	}

	// 修改updateTime字段
	cipForReceive.UpdatedTime = common.GetCurrentTime()
	// 执行更新操作（链上+链下）
	err = codeService.UpdateCip(cipForReceive, loginUser.Role)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "修改成功"})
}

// DeleteCIP 删除代码产权
func DeleteCIP(ctx *gin.Context) {
	cpi := ctx.Param("cpi")

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 检查登录用户身份权限
	cip, err := codeService.GetFromChain(cpi, loginUser.Role)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此代码产权不存在，删除失败"})
		return
	} else if cip.Author != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是该项知识产权的创作者，无权删除"})
		return
	}

	// 执行删除操作
	err = codeService.DeleteCipByCpi(cpi, loginUser.Role)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "删除成功"})
}

// QueryCipByFilename 创作者模糊查询本人的代码产权（已上链）
func QueryCipByFilename(ctx *gin.Context) {
	// 接收参数(filename)
	var cip domain.CodeRecords
	err := ctx.ShouldBind(&cip)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "接收参数失败，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	cip.Author = loginUser.Username

	// 模糊查询
	cips, err := codeService.GetByFilename(cip)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: cips})
}

// QueryDetails 创作者查看组件元数据（未上链和已上链皆可）
func QueryDetails(ctx *gin.Context) {
	// 接收参数(id)
	var cr domain.CodeRecords
	err := ctx.ShouldBind(&cr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "接收参数失败，请重试"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 获取产权元数据
	crDetails, err := codeService.GetById(cr.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "查询失败，此项代码产权不存在"})
	} else if crDetails.Author != loginUser.Username { // 检查查询者身份
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是此产权的拥有人，无法查看"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: crDetails})
}

// QueryCipInfo 使用者查看他人组件元数据+创作者身份信息
func QueryCipInfo(ctx *gin.Context) {
	// 接收参数（id）
	var cr domain.CodeRecords
	err := ctx.ShouldBind(&cr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "接收参数失败，请重试"})
		return
	}

	// 获取产权元数据
	crDetails, err := codeService.GetById(cr.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "查询失败，此项代码产权不存在"})
		return
	}

	// 获取创作者身份信息
	user, _ := userService.GetUser(crDetails.Author)

	var cipInfoDto domain.SubCipDto
	copier.Copy(&cipInfoDto, &crDetails)
	copier.Copy(&cipInfoDto, &user)

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: cipInfoDto})
}

// QueryByAuthor 创作者查看拥有的所有代码产权
func QueryByAuthor(ctx *gin.Context) {
	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	cipRecords, err := codeService.GetCipByAuthor(loginUser.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "查询成功", Data: cipRecords})
}

// QueryByDesc 使用者按功能筛选开放代码产权
func QueryByDesc(ctx *gin.Context) {
	var cip domain.CodeRecords
	err := ctx.ShouldBind(&cip)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "参数接收失败"})
		return
	}

	// 查询
	cips, err := codeService.GetByDesc(cip.Description)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "没有找到符合您需求的组件"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: cips})
}

// QueryAllCips 查看所有开放组件
func QueryAllCips(ctx *gin.Context) {
	cips, err := codeService.GetOpenCips()
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "查询成功", Data: cips})
}

// Check 检查代码产权真实性和有效性和产权文件完整性（结合区块链和存储库）
func Check(ctx *gin.Context) {
	// 接收参数（id）
	var cr domain.CodeRecords
	err := ctx.ShouldBind(&cr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "参数接收失败"})
		return
	}

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 1 结合存储库检查组件文件完整性
	// 1.1 根据id查询文件存储地址ipfsAddress
	cip, errGet := codeService.GetById(cr.ID)
	if errGet != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "未通过验证：此项代码产权不存在"})
		return
	}

	// 1.2 获取文件，检查文件是否存在
	ipFilePath, err := ipfs.DownloadIPFS(cip.Addr)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "未通过验证：此代码产权文件不存在"})
		return
	}

	// 1.3 对文件内容进行哈希运算，检查文件是否被篡改
	draftFileURL := filepath.Join(ipFilePath, cip.Addr)
	fmt.Println(draftFileURL)
	readerIpCoreFile, _ := os.Open(draftFileURL)         // 打开链下组件文件
	ipCoreFileContent, _ := io.ReadAll(readerIpCoreFile) // 读取链下组件文件数据
	cpi := common.SHA256(ipCoreFileContent)              // SHA-256运算
	defer func() {                                       // 删除临时文件
		_ = readerIpCoreFile.Close()
		_ = os.RemoveAll(ipFilePath)
	}()

	// 1.4 比对Mysql数据库中cpi
	if cpi != cip.Cpi {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "未通过验证：存储库中代码产权信息已被篡改，请进行申诉"})
		return
	}

	// 2 结合区块链进行代码产权校验
	// 2.1 查询代码产权是否存在于区块链上
	cipChain, err := codeService.GetFromChain(cpi, loginUser.Role)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "未通过验证：区块链上的代码产权信息已被篡改，请进行申诉"})
		return
	}
	// 2.2 检查其他信息
	if cip.Filename != cipChain.Filename || cip.Version != cipChain.Version || cip.Tag != cipChain.Tag ||
		cip.Author != cipChain.Author || cip.RegTime != cipChain.RegTime || cip.Language != cipChain.Language ||
		strconv.Itoa(cip.Price) != cipChain.Price || cip.Status != cipChain.Status ||
		cip.Description != cipChain.Description || cip.Addr != cipChain.Addr {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "未通过验证：代码产权信息已被篡改，请进行申诉"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "验证通过，此项代码产权未发现异常", Data: cipChain})
}

// Count 统计代码产权数量
func Count(ctx *gin.Context) {
	data := codeService.GetCount()
	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: data})
}

package product

import (
	"CipProject/application/common"
	"CipProject/application/domain"
	"CipProject/application/ipfs"
	"CipProject/application/service/codeService"
	"CipProject/application/service/productService"
	"CipProject/application/service/subscribeService"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

// Create 创建组件服务
func Create(ctx *gin.Context) {
	// 接收参数（cpi）
	var sub domain.SubRecords
	_ = ctx.ShouldBind(&sub)

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 获取订阅数据
	sr, err := subscribeService.GetIDByCpi(sub.Cpi)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此组件不存在"})
		return
	}
	subRecords, _ := subscribeService.GetById(sr.ID)

	// 验证代码产权（链上方式）
	cipMeta, err := codeService.GetFromChain(sub.Cpi, loginUser.Role)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此组件未获得知识产权，不允许使用"})
		return
	}

	// 检查登录用户是否为该组件的订阅者
	subMetaInfo, err := subscribeService.GetFromChain(strconv.FormatInt(sr.ID, 10), loginUser.Role) // 读取链上的订阅合同
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您未订阅此组件，启动失败"})
		return
	} else if subMetaInfo.PartyA != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您未订阅此组件，启动失败"})
		return
	}

	// 检查订阅是否开始生效
	if flag := common.After(subMetaInfo.Start, common.GetCurrentTime()); flag == false {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "未到订阅合同生效时间，请耐心等待"})
		return
	}

	// 检查组件容器是否已创建
	if subRecords.DockerId != 0 {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "请勿重复创建"})
		return
	}

	// 创建服务
	// 1 从IPFS中提取项目文件
	tempFile, err := ipfs.DownloadIPFS(cipMeta.Addr) // tempFile是文件路径，例如：/tmp/xxx文件，下载的是.zip文件
	defer os.RemoveAll(tempFile)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "组件数据丢失，请联系项目乙方"})
		return
	}

	// 处理文件目录
	destDir, _ := os.MkdirTemp("", "ciprtip_") // destDir是解文件压之后存放的目录，即ciprtip_随机符
	defer os.RemoveAll(destDir)

	// 2 项目文件解压（指定.zip文件）
	cmd := exec.Command("unzip", "-o", tempFile+"/*", "-d", destDir) // 创建解压指令的命令对象
	cmd.Stdout = os.Stdout                                           // 设置命令的输出和错误输出
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "组件数据处理失败，请重试"})
		return
	}

	// 3 切换工作目录，即Dockerfile所在的目录下
	err = os.Chdir(destDir + "/" + cipMeta.Filename)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "组件文件不存在，启动失败"})
		return
	}

	// 容器命名规则（唯一性）：ciprtip_cpi_使用者username_产权版本号version
	dockerName := "ciprtip_" + subRecords.Cpi + loginUser.Username + cipMeta.Version + common.GetTimestampAsString()
	// 镜像命名规则（唯一性）：ciprtip_代码产权唯一标识码cpi_使用者username
	imageName := "ciprtip_" + cipMeta.Cpi + loginUser.Username + common.GetTimestampAsString()
	// 镜像标签
	tag := cipMeta.Version
	// 镜像名及标签
	imageNameAndTag := imageName + ":" + tag

	// 4 创建Docker镜像
	err = common.CreateDockerImage(imageNameAndTag)
	if err != nil {
		defer func(imageNameAndTag string) {
			_ = common.RemoveDockerImage(imageNameAndTag)
		}(imageNameAndTag)

		ctx.JSON(http.StatusOK, domain.Response{Msg: "容器镜像构建失败，请重试"})
		return
	}

	// 5 运行容器docker run -id -p 服务器节点端口:innerPort --name=容器名 镜像名:标签
	crMeta, _ := codeService.GetPortByCpi(sub.Cpi)
	outerPort, err := common.CreateDocker(crMeta.InnerPort, dockerName, imageNameAndTag)
	if err != nil {
		defer func(imageNameAndTag string) {
			_ = common.RemoveDockerImage(imageNameAndTag)
		}(imageNameAndTag)

		defer func(dockerName string) {
			_ = common.RemoveDocker(dockerName)
		}(dockerName)

		defer func(dockerName string) {
			_ = common.StopDocker(dockerName)
		}(dockerName)

		ctx.JSON(http.StatusOK, domain.Response{Msg: "服务实例创建失败，" + err.Error()})
		return
	}

	publicIP, _ := common.GetPublicIP()

	// 组件服务创建成功，记录创建日志
	var newDocker = domain.DockerInfo{
		ID:          common.NewSnowflake(2).Generate(), // 调用 Generate 方法生成唯一 ID
		Cpi:         sub.Cpi,
		Owner:       subMetaInfo.PartyA, // 使用者
		Ip:          publicIP,
		Port:        outerPort,
		ImageName:   imageName,
		Tag:         tag,
		DockerName:  dockerName,
		Status:      domain.DockerRunning,
		CreatedTime: common.GetCurrentTime(),
		StopTime:    "",
	}
	newDocker.StartTime = newDocker.CreatedTime
	_ = productService.New(newDocker, subRecords.ID)

	//term, _ := subscribeService.GetTermByDocker(subRecords.DockerId)
	//Term: math.Round(term*24*10) / 10,

	createdDocker := domain.CreatedDockerDTO{
		ID:   newDocker.ID,
		Ip:   newDocker.Ip,
		Port: newDocker.Port,
	}
	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: createdDocker})
}

// Start 启动组件服务
func Start(ctx *gin.Context) {
	// 接收参数-->Docker id
	var docker domain.DockerInfo
	_ = ctx.ShouldBind(&docker)

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 获取服务数据
	service, err := productService.GetById(docker.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此服务不存在"})
		return
	}

	// 1 检查登录用户身份是否为合法使用者
	if service.Owner != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是此服务的使用者，无权启动"})
		return
	}

	// 2 检查运行状态
	if service.Status != domain.DockerStopped {
		ctx.JSON(http.StatusOK, domain.Response{Msg: service.Status})
		return
	}

	err = common.StartDocker(service.DockerName)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "服务启动失败"})
		return
	}

	// 更新服务的状态
	service.Status = domain.DockerRunning
	service.StartTime = common.GetCurrentTime()
	_ = productService.UpdateById(service)

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "服务启动成功"})
}

// Stop 停止服务运行
func Stop(ctx *gin.Context) {
	// 接收参数-->Docker id
	var docker domain.DockerInfo
	_ = ctx.ShouldBind(&docker)

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 获取服务数据
	service, err := productService.GetById(docker.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "服务不存在，请先创建"})
		return
	}

	// 1 检查登录用户身份是否为合法使用者
	if service.Owner != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是此服务使用者，无权停止"})
		return
	}

	// 2 检查运行状态
	if service.Status != domain.DockerRunning {
		ctx.JSON(http.StatusOK, domain.Response{Msg: service.Status + "请勿重试"})
		return
	}

	// 停止运行docker容器：docker stop 容器名
	err = common.StopDocker(service.DockerName)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "停止失败，请重试"})
		return
	}

	// 更新服务数据
	service.StopTime = common.GetCurrentTime()
	service.Status = domain.DockerStopped
	// 计算本次使用时长
	subInfo, _ := subscribeService.GetByDockerId(service.ID)
	term, _ := common.DateInterval(service.StartTime, service.StopTime)
	subInfo.Term = term + subInfo.Term

	// 更新容器状态，服务使用时长
	_ = productService.Update(service, subInfo)
	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "成功停止服务"})
}

// Destroy 销毁服务
func Destroy(ctx *gin.Context) {
	// 接收参数-->Docker id
	var docker domain.DockerInfo
	_ = ctx.ShouldBind(&docker)

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	// 获取服务数据
	service, err := productService.GetById(docker.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "此服务不存在"})
		return
	}

	// 1 检查登录用户身份是否为合法使用者
	if service.Owner != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您不是服务的使用者，无权销毁"})
		return
	}

	// 2 检查运行状态
	if service.Status == domain.DockerRunning || service.Status == domain.DockerStopped {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "请先终止订阅"})
		return
	} else if service.Status == domain.DockerDestroyed {
		ctx.JSON(http.StatusOK, domain.Response{Msg: domain.DockerDestroyed + "，请勿重试"})
		return
	}

	// 3 检查订阅合同状态
	subInfo, _ := subscribeService.GetByDockerId(service.ID)
	if subInfo.Status != domain.SubWithoutPay {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "当前订阅状态：" + subInfo.Status})
		return
	}

	// 销毁服务
	// 1 删除Docker实例
	err = common.RemoveDocker(service.DockerName)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "销毁实例失败，请重试"})
		return
	}

	// 2 删除Docker镜像
	imageNameAndTag := service.ImageName + ":" + service.Tag
	err = common.RemoveDockerImage(imageNameAndTag)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "销毁镜像失败，请重试"})
		return
	}

	// 删除Docker实例信息
	service.Status = domain.DockerDestroyed
	_ = productService.UpdateById(service)

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Msg: "成功销毁服务"})
}

// QueryById 查看服务详情
func QueryById(ctx *gin.Context) {
	// 接收参数-->Docker id
	var docker domain.DockerInfo
	_ = ctx.ShouldBind(&docker)

	fmt.Println(docker)

	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	dockerDto, err := productService.GetDetails(docker.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	} else if dockerDto.Owner != loginUser.Username {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "您无权访问此服务"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: dockerDto})
}

// QueryByOwner 查看已创建的服务列表
func QueryByOwner(ctx *gin.Context) {
	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	dockerDto, err := productService.GetByOwner(loginUser.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: dockerDto})
}

// QueryRaw 查看未创建的服务
func QueryRaw(ctx *gin.Context) {
	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	raw, err := productService.GetRawByOwner(loginUser.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: raw})
}

// QueryDestroyList 查看待销毁服务列表
func QueryDestroyList(ctx *gin.Context) {
	login, _ := ctx.Get("login_user")
	loginUser := login.(domain.User)

	dockerDto, err := productService.GetDestroyList(loginUser.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{Msg: "无查询结果"})
		return
	}

	ctx.JSON(http.StatusOK, domain.Response{Code: 1, Data: dockerDto})
}

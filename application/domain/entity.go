package domain

import (
	"mime/multipart"
)

// Cart 购物车
type Cart struct {
	ID          int64  `json:"id" form:"id" gorm:"column:id"`
	CodeId      int64  `json:"codeId" form:"codeId" gorm:"column:code_id"`
	Owner       string `json:"owner" form:"owner" gorm:"column:owner"`
	CreatedTime string `json:"createdTime" form:"createdTime" gorm:"column:created_time"`
	UpdatedTime string `json:"updatedTime" form:"updatedTime" gorm:"column:updated_time"`
}

// CodeRecords 代码产权和草稿项目备份，提高管理效率，只保存最新状态，结合区块链做数据真实性检查
type CodeRecords struct {
	ID          int64                 `json:"id" form:"id" gorm:"column:id"`                             // 主键，映射约束：结构体的字段驼峰命名，对应表列名的蛇形命名，如CreateTime对应表列名create_time
	Cpi         string                `json:"cpi" form:"cpi" gorm:"column:cpi"`                          // 产权唯一标识符
	Filename    string                `json:"filename" form:"filename" gorm:"column:filename"`           // 项目名
	Version     string                `json:"version" form:"version" gorm:"column:version"`              // 版本号
	Tag         string                `json:"tag" form:"tag" gorm:"column:tag"`                          // 标签
	Author      string                `json:"author" form:"author" gorm:"column:author"`                 // 创作者（对应用户实体Username）
	Company     string                `json:"company" form:"company" gorm:"column:company"`              // 创作者公司
	CreatedTime string                `json:"createdTime" form:"createdTime" gorm:"column:created_time"` // 创建时间
	RegTime     string                `json:"regTime" form:"regTime" gorm:"column:reg_time"`             // 授权时间（即上链时间）
	UpdatedTime string                `json:"updatedTime" form:"updatedTime" gorm:"column:updated_time"` // 更新时间
	Language    string                `json:"language" form:"language" gorm:"column:language"`           // 程序设计语言（Golang;Java;表示Python）
	Price       int                   `json:"price" form:"price" gorm:"price"`                           // 订阅价格（元/天）报价
	Status      string                `json:"status" form:"status" gorm:"column:status"`
	Description string                `json:"description" form:"description" gorm:"column:description"` // 项目描述
	Addr        string                `json:"addr" form:"addr" gorm:"column:addr"`                      // 文件的链下存储地址
	InnerPort   string                `json:"innerPort" form:"innerPort" gorm:"column:inner_port"`      // 项目内部访问端口
	DraftFile   *multipart.FileHeader `json:"draftFile" form:"draftFile" gorm:"-"`                      // 产权文件
}

// 代码产权状态
const (
	CodeNoCip    = "无知识产权"
	CodeDisabled = "禁用"
	CodeOpen     = "开放"
)

// SubRecords 订阅记录
type SubRecords struct {
	ID          int64   `json:"id" form:"id" gorm:"column:id"`
	CodeId      int64   `json:"codeId" form:"codeId" gorm:"column:code_id"`
	Cpi         string  `json:"cpi" form:"cpi" gorm:"column:cpi"`
	PartyA      string  `json:"partyA" form:"partyA" gorm:"column:party_a"`                // 甲方（使用者username）
	PartyB      string  `json:"partyB" form:"partyB" gorm:"column:party_b"`                // 乙方（创作者username）
	StartTime   string  `json:"startTime" form:"startTime" gorm:"column:start_time"`       // 服务订阅生效日期时间
	EndTime     string  `json:"endTime" form:"endTime" gorm:"column:end_time"`             // 服务订阅终止日期时间
	CreatedTime string  `json:"createdTime" form:"createdTime" gorm:"column:created_time"` // 订阅申请表生成日期时间
	UpdatedTime string  `json:"updatedTime" form:"updatedTime" gorm:"column:updated_time"` // 订阅申请表更新日期时间
	Term        float64 `json:"term" form:"term" gorm:"column:term"`                       // 订阅期限（单位：天）
	RealPrice   int     `json:"realPrice" form:"realPrice" gorm:"column:real_price"`       // 双方协商后价格（元/天）
	Cost        float64 `json:"cost" form:"cost" gorm:"column:cost"`                       // 当前费用
	Status      string  `json:"status" form:"status" gorm:"column:status"`
	Note        string  `json:"note" form:"note" gorm:"column:note"` // 备注
	DockerId    int64   `json:"dockerId" form:"dockerId" gorm:"column:docker_id"`
}

// 订阅申请单状态
const (
	SubToExamine     = "甲方申请订阅，待乙方审核"
	SubCancel        = "甲方取消订阅申请"
	SubToConfirm     = "乙方同意，待甲方确认订阅条款"
	SubRefusedPartyB = "乙方已拒绝"
	SubRefusedPartyA = "甲方拒绝订阅条款"
	SubApproved      = "双方已签订合同"
	SubWithoutPay    = "费用已清算但未支付"
	SubPayed         = "费用已支付，本次交易完成"
)

// DockerInfo 容器化服务元数据
type DockerInfo struct {
	ID          int64  `json:"id" form:"id" gorm:"column:id"`
	Cpi         string `json:"cpi" form:"cpi" gorm:"column:cpi"`                       // 代码产权唯一标识码
	Owner       string `json:"owner" form:"owner" gorm:"column:owner"`                 // 甲方（使用者username）
	Ip          string `json:"ip" form:"ip" gorm:"column:ip"`                          // ip地址
	Port        int    `json:"port" form:"port" gorm:"column:port"`                    // 端口号
	ImageName   string `json:"imageName" form:"imageName" gorm:"column:image_name"`    // 依赖的镜像名，命名规则为ciprtip_代码产权唯一标识码cpi_使用者username
	Tag         string `json:"tag" form:"tag" gorm:"column:tag"`                       // 镜像标签，命名规则为项目版本号
	DockerName  string `json:"dockerName" form:"dockerName" gorm:"column:docker_name"` // 容器名，命名规则为ciprtip_代码产权唯一标识码cpi_使用者username
	Status      string `json:"status" form:"status" gorm:"column:status"`
	CreatedTime string `json:"createdTime" form:"createdTime" gorm:"column:created_time"` // 服务创建日期时间
	StartTime   string `json:"startTime" form:"startTime" gorm:"column:start_time"`       // 上次启动日期时间
	StopTime    string `json:"stopTime" form:"stopTime" gorm:"column:stop_time"`          // 上次暂停日期时间
}

// 容器状态
const (
	DockerRunning   = "服务正在运行中"
	DockerStopped   = "服务已暂停"
	DockerToDestroy = "服务终止待消毁"
	DockerDestroyed = "服务已销毁"
)

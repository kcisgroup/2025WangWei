package domain

// User 用户信息，未保证安全性，不做链下备份
type User struct {
	Username string `json:"username" form:"username"` // 用户名默认使用用户所在公司的工号，限制其他用户非法注册
	Pwd      string `json:"pwd" form:"pwd"`
	Name     string `json:"name" form:"name"`
	Role     string `json:"role" form:"role"` // 就职公司/组织
	Tel      string `json:"tel" form:"tel"`
	Idnum    string `json:"idnum" form:"idnum"` // 身份证号
	Addr     string `json:"addr" form:"addr"`
	Gender   string `json:"gender" form:"gender"`
	Status   string `json:"status" form:"status"` // 0各公司管理员，1普通员工用户
}

// CodeIntellectualProperty 代码产权元数据
type CodeIntellectualProperty struct {
	Cpi         string `json:"cpi" form:"cpi"`
	Filename    string `json:"filename" form:"filename"`         // 项目名
	Version     string `json:"version" form:"version"`           // 创作者自定义的版本号信息
	Tag         string `json:"tag" form:"tag" gorm:"column:tag"` // 自定义标签
	Author      string `json:"author" form:"author"`             // 创作者（对应用户实体Username）
	RegTime     string `json:"regTime" form:"regTime"`           // 授权时间（即上链时间）
	UpdatedTime string `json:"updatedTime" form:"updatedTime"`   // 产权元数据更新时间
	Language    string `json:"language" form:"language"`         // 代码语言区分（0表示Golang，1表示Java，2表示Python）
	Price       string `json:"price" form:"price"`               // 租赁单价（元/天）
	Status      string `json:"status" form:"status"`             // 代码知识产权的状态，无知识产权；禁用；开放
	Description string `json:"description" form:"description"`   // 功能描述
	Addr        string `json:"addr" form:"addr"`                 // 文件的链下存储地址
}

// SubContract 服务订阅数据
type SubContract struct {
	ID        string `json:"id" form:"id" copier:"-"`
	Cpi       string `json:"cpi" form:"cpi"`
	PartyA    string `json:"partyA" form:"partyA"`                  // 甲方username（使用者）
	PartyB    string `json:"partyB" form:"partyB"`                  // 乙方username（创作者）
	Start     string `json:"start" form:"start" copier:"StartTime"` // 合同生效时间，格式为"2006-01-02 15:04:05"
	End       string `json:"end" form:"end" copier:"EndTime"`       // 服务截止使用时间，初始时"待定"
	RealPrice string `json:"realPrice" form:"realPrice" copier:"-"` // 双方协商后价格（元/天）
	Cost      string `json:"cost" form:"cost" copier:"-"`           // 费用，初始时"待定"
	Status    string `json:"status" form:"status"`                  // 交易状态：(1)双方签订合同，费用未清算;(2)费用已清算但未支付;(3)费用已支付，本次交易完成
	Note      string `json:"note" form:"note"`                      // 备注
}

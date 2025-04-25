package domain

// CartInfoDTO 购物车详情数据传输对象
type CartInfoDTO struct {
	Cart
	Filename    string `json:"filename" form:"filename" gorm:"column:filename"`
	Version     string `json:"version" form:"version" gorm:"column:version"`
	Tag         string `json:"tag" form:"tag" gorm:"column:tag"`
	Price       int    `json:"price" form:"price" gorm:"column:price"`
	Description string `json:"description" form:"description" gorm:"column:description"`
}

// SubInfoDTO 订阅申请数据传输对象（申请者查看）
type SubInfoDTO struct {
	ID          int64  `json:"id" form:"id" gorm:"column:id"`
	CodeId      int64  `json:"codeId" form:"codeId" gorm:"column:code_id"`
	PartyA      string `json:"partyA" form:"partyA" gorm:"column:party_a"`
	PartyB      string `json:"partyB" form:"partyB" gorm:"column:party_b"`
	EndTime     string `json:"endTime" form:"endTime" gorm:"column:end_time"`             // 软件订阅终止日期时间
	CreatedTime string `json:"createdTime" form:"createdTime" gorm:"column:created_time"` // 订阅申请表生成日期时间
	UpdatedTime string `json:"updatedTime" form:"updatedTime" gorm:"column:updated_time"` // 订阅申请表更新日期时间
	RealPrice   int    `json:"realPrice" form:"realPrice" gorm:"column:real_price"`
	Status      string `json:"status" form:"status" gorm:"column:status"`
	Note        string `json:"note" form:"note" gorm:"column:note"`
	Author      string `json:"author" form:"author" gorm:"column:author"`
	DockerId    int64  `json:"dockerId" form:"dockerId" gorm:"column:docker_id"`
	RawProductDTO
}

type SubCipDto struct {
	ID          int64  `json:"id" form:"id" gorm:"column:id"`
	Filename    string `json:"filename" form:"filename" gorm:"column:filename"`
	Version     string `json:"version" form:"version" gorm:"column:version"`
	Tag         string `json:"tag" form:"tag" gorm:"column:tag"`
	Company     string `json:"company" form:"company" gorm:"column:company"`
	Author      string `json:"author" form:"author" gorm:"column:author"`
	RegTime     string `json:"regTime" form:"regTime" gorm:"column:reg_time"`
	UpdatedTime string `json:"updatedTime" form:"updatedTime" gorm:"column:updated_time"`
	Language    string `json:"language" form:"language" gorm:"column:language"`
	Price       int    `json:"price" form:"price" gorm:"price"`
	Description string `json:"description" form:"description" gorm:"column:description"`
	Name        string `json:"name" form:"name"` // 创作者姓名
	Role        string `json:"role" form:"role"`
	Tel         string `json:"tel" form:"tel"` // 创作者手机号
}

type DockerInfoDTO struct {
	ID          int64   `json:"id" form:"id" gorm:"column:id"`
	Status      string  `json:"status" form:"status" gorm:"column:status"`
	Ip          string  `json:"ip" form:"ip" gorm:"column:ip"`
	Port        int     `json:"port" form:"port" gorm:"column:port"`
	CreatedTime string  `json:"createdTime" form:"createdTime" gorm:"column:created_time"`
	Term        float64 `json:"term" form:"term" gorm:"column:term"`
	RawProductDTO
}

type RawProductDTO struct {
	Cpi       string `json:"cpi" form:"cpi" gorm:"column:cpi"`
	Filename  string `json:"filename" form:"filename" gorm:"column:filename"`
	Version   string `json:"version" form:"version" gorm:"column:version"`
	Tag       string `json:"tag" form:"tag" gorm:"column:tag"`
	StartTime string `json:"startTime" form:"startTime" gorm:"column:start_time"`
	Language  string `json:"language" form:"language" gorm:"column:language"`
}

type DockerDestroy struct {
	ID         int64   `json:"id" form:"id" gorm:"column:id"`
	Term       float64 `json:"term" form:"term" gorm:"column:term"`
	Cost       float64 `json:"cost" form:"cost" gorm:"column:cost"`
	Cpi        string  `json:"cpi" form:"cpi" gorm:"column:cpi"`
	DockerName string  `json:"dockerName" form:"dockerName" gorm:"column:docker_name"`
}

type CountDto struct {
	TotalUsers    int64 `json:"totalUsers" form:"totalUsers" gorm:"column:total_users"`
	TotalDrafts   int64 `json:"totalDrafts" form:"totalDrafts" gorm:"column:total_drafts"`
	TotalCips     int64 `json:"totalCips" form:"totalCips" gorm:"column:total_cips"`
	TotalOpenCips int64 `json:"totalOpenCips" form:"totalOpenCips" gorm:"column:total_open_cips"`
	TotalDockers  int64 `json:"totalDockers" form:"totalDockers" gorm:"column:total_dockers"`
}

type CreatedDockerDTO struct {
	ID   int64  `json:"id" form:"id" gorm:"column:id"`
	Ip   string `json:"ip" form:"ip" gorm:"column:ip"`
	Port int    `json:"port" form:"port" gorm:"column:port"`
}

type DockerDetailsDTO struct {
	ID          int64   `json:"id" form:"id" gorm:"column:id"`
	Cpi         string  `gorm:"column:cpi"`
	Owner       string  `gorm:"column:owner"`
	Ip          string  `json:"ip" form:"ip" gorm:"column:ip"`
	Port        int     `json:"port" form:"port" gorm:"column:port"`
	DockerName  string  `json:"dockerName" form:"dockerName" gorm:"column:docker_name"`
	Status      string  `json:"status" form:"status" gorm:"column:status"`
	CreatedTime string  `json:"createdTime" form:"createdTime" gorm:"column:created_time"`
	Term        float64 `json:"term" form:"term" gorm:"column:term"` // 使用时长（小时）
	RealPrice   int     `json:"realPrice" form:"realPrice" gorm:"column:real_price"`
	Cost        float64 `json:"cost" form:"cost" gorm:"column:cost"`
	CurrentTime string  `json:"currentTime" form:"currentTime" gorm:"column:current_time"`
	Filename    string  `json:"filename" form:"filename" gorm:"column:filename"`
	Version     string  `json:"version" form:"version" gorm:"column:version"`
	Author      string  `json:"author" form:"author" gorm:"column:author"`
	Company     string  `json:"company" form:"company" gorm:"column:company"`
	Description string  `json:"description" form:"description" gorm:"column:description"`
}

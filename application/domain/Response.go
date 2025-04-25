package domain

import "github.com/dgrijalva/jwt-go"

// Response 响应数据格式
type Response struct {
	Code  int         `json:"code"` // 0表示失败，1表示成功
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Token string      `json:"token"` // Jwt Token
}

// MyClaims 用于Jwt Token
// 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	Username string `json:"username" form:"username"` // username塞入Token
	jwt.StandardClaims
}

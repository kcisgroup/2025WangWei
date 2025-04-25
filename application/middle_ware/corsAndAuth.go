package middle_ware

import (
	"CipProject/application/common"
	"CipProject/application/domain"
	"CipProject/application/service/userService"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// CorsMiddleWare 解决跨域请求问题
func CorsMiddleWare(c *gin.Context) {
	method := c.Request.Method
	origin := c.GetHeader("Origin")
	c.Header("Access-Control-Allow-Origin", origin)      // 注意这一行，不能配置为通配符“*”号
	c.Header("Access-Control-Allow-Credentials", "true") // 注意这一行，必须设定为true
	// 我们自定义的header字段都需要在这里声明
	c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Cookie, Origin, X-Requested-With, Content-Type, Accept, Authorization, Token, Timestamp, UserId")
	c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PUT")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,cache-control")

	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		//c.AbortWithStatus(http.StatusNoContent)
		c.AbortWithStatus(http.StatusOK)
	}
	// 处理请求
	c.Next()

}

// AuthMiddleWare 用户登录检查
func AuthMiddleWare(c *gin.Context) {
	switch c.Request.URL.Path {
	case "/user/login":
		break
	case "/user/register":
		break
	case "/user/exitLogin":
		break
	default:
		// 获取请求头中Authorization字段内容
		tokenString := c.GetHeader("Authorization")
		// 请求头Authorization内容为空或不含Bearer前缀；http有很多前缀，Bearer是OAuth2.0规范
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			c.JSON(http.StatusOK, domain.Response{Msg: "用户权限不足，请先登录"})
			c.Abort()
			return
		}
		tokenString = tokenString[7:]
		// 解析Jwt Token
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusOK, domain.Response{Msg: "用户权限不足，请先登录"})
			c.Abort()
			return
		}
		// 验证Token通过后，获取claims中的username
		username := claims.Username
		if username == "" {
			c.JSON(http.StatusOK, domain.Response{Msg: "用户权限不足，请先登录"})
			c.Abort()
			return
		}
		// username存在，查询user信息
		loginUser, errGet := userService.GetUser(username)
		if errGet != nil {
			c.JSON(http.StatusOK, domain.Response{Msg: "登录用户的信息有误，请重新登录"})
			c.Abort()
			return
		}
		// loginUser存入上下文，方便后面处理器使用
		c.Set("login_user", loginUser)
	}
	c.Next()
}

package common

import (
	"CipProject/application/domain"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"math"
	"net"
	"net/http"
	"sync"
	"time"
)

var jwtKey = []byte("jwt_token_of_coleW")

// ReleaseToken 生成Jwt token
func ReleaseToken(user domain.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 设置token过期时间，七天

	// 创建一个我们自己的声明的数据
	claims := &domain.MyClaims{
		Username: user.Username, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			// 注意INT不能直接与time.hour相乘，所以用time.duration转换一下
			ExpiresAt: expirationTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),     // token发放时间
			Issuer:    "coleW",               // 签发人
			Subject:   "login_user token",    // token主题
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析Jwt token
func ParseToken(tokenString string) (*jwt.Token, *domain.MyClaims, error) {
	// 解析token
	claims := &domain.MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}

// GetPrivateIP 获取本地连接IP地址
func GetPrivateIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// GetPublicIP 获取公网IP地址
func GetPublicIP() (string, error) {
	resp, err := http.Get("https://ifconfig.me/ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

// TimestampToDate 将时间戳转换为日期格式
func TimestampToDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

// GetTimestampAsString 获取时间戳字符串
func GetTimestampAsString() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

// GetCurrentTime 获取时间，格式：2006-01-02 15:04:05
func GetCurrentTime() string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return currentTime
}

// DateInterval 计算两个日期的时间间隔（单位：天数）
func DateInterval(start, end string) (float64, error) {
	startDate, errOfStart := time.Parse(time.DateTime, start) // 将字符串类型的时间转为time类型
	endDate, err2OfEnd := time.Parse(time.DateTime, end)

	if errOfStart != nil {
		return -1, errOfStart
	}
	if err2OfEnd != nil {
		return -1, err2OfEnd
	}

	flag := endDate.After(startDate) // 判断起始时间和截止时间是否规范
	if flag != true {
		return -1, errors.New("截止时间应晚于起始时间")
	}

	interval := endDate.Sub(startDate).Hours() / 24
	interval = math.Round(interval*100000) / 100000 // 保留五位小数
	return interval, nil
}

// After 判断a是否早于b
func After(a, b string) bool {
	startDate, errOfStart := time.Parse(time.DateTime, a) // 将字符串类型的时间转为time类型
	endDate, err2OfEnd := time.Parse(time.DateTime, b)
	if errOfStart != nil {
		return false
	}
	if err2OfEnd != nil {
		return false
	}

	return endDate.After(startDate)
}

// SHA256 SHA-256单向加密算法
func SHA256(input []byte) string {
	// 计算SHA256哈希值
	hash := sha256.Sum256(input)

	// 将哈希值转换为16进制字符串
	hashString := fmt.Sprintf("%x", hash)

	return hashString
}

// Snowflake 雪花算法的实现
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64 // 上次生成ID的时间戳
	nodeID        int64 // 节点ID
	sequence      int64 // 序列号
}

// NewSnowflake 创建一个Snowflake实例
func NewSnowflake(nodeID int64) *Snowflake {
	return &Snowflake{
		lastTimestamp: time.Now().UnixNano() / 1e6,
		nodeID:        nodeID,
		sequence:      0,
	}
}

func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	timestamp := time.Now().UnixNano() / 1e6
	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & 4095
		if s.sequence == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}
	s.lastTimestamp = timestamp
	startId := ((timestamp - 1288834974657) << 22) | (s.nodeID << 12) | s.sequence
	id := startId % 1e15
	return id
}

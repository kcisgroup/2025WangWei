package portsManager

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// portManager 端口管理器
type portManager struct {
	startPort int
	endPort   int
	allocated map[int]bool
	mutex     sync.Mutex
}

// pm 端口管理器实例（单例模式、缓存）
var pm *portManager

func GetPortManager() *portManager {
	return pm
}

func init() {
	// 获取系统已占用的端口
	pm = newPortManager(8001, 30000)
	occupiedPorts, err := getOccupiedPorts()
	if err != nil {
		_ = fmt.Sprintf("Error getting occupied ports:%s", err.Error())
		return
	}

	// 将已占用的端口添加到端口管理器
	for _, port := range occupiedPorts {
		pm.allocated[port] = true
	}

	fmt.Println(pm.allocated)
}

// newPortManager 创建端口管理器
func newPortManager(startPort, endPort int) *portManager {
	return &portManager{
		startPort: startPort,
		endPort:   endPort,
		allocated: make(map[int]bool),
	}
}

// AllocatePort 分配唯一的端口号
func (pm *portManager) AllocatePort() (int, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// 寻找第一个空闲的端口号
	for port := pm.startPort; port <= pm.endPort; port++ {
		flag := pm.isPortOccupied(port)
		if !pm.allocated[port] && flag == false {
			pm.allocated[port] = true
			return port, nil
		}
	}

	return 0, errors.New("当前没有空闲的端口")
}

// getOccupiedPorts 获取系统已被占用的端口号
func getOccupiedPorts() ([]int, error) {
	output, err := exec.Command("netstat", "-tln").Output() // output是netstat -tln命令执行结果的[]byte类型
	if err != nil {
		return nil, err
	}

	var occupiedPorts []int

	lines := strings.Split(string(output), "\n") // 以换行符分割
	for _, line := range lines[2:] {             // line举例：tcp   0   0   0.0.0.0:9051   0.0.0.0:*   LISTEN
		fields := strings.Fields(line) // 以空格分割
		if len(fields) >= 4 {
			address := fields[3]
			parts := strings.Split(address, ":") // 以:分割，取末尾字符串即为端口号
			portStr := parts[len(parts)-1]
			port, err := strconv.Atoi(portStr)
			if err == nil {
				occupiedPorts = append(occupiedPorts, port)
			}
		}
	}

	return occupiedPorts, nil
}

// DeallocatePort 释放端口号
func (pm *portManager) DeallocatePort(port int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	delete(pm.allocated, port)
}

// isPortOccupied 判断端口号是否被占用
func (pm *portManager) isPortOccupied(port int) bool {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}
	conn.Close()
	return false
}

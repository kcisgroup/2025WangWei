package common

import (
	"CipProject/application/portsManager"
	"os"
	"os/exec"
	"strconv"
)

func CreateDockerImage(imageNameAndTag string) error {
	// docker build --no-cache -t 镜像名:标签 $GOPATH/src/xxx项目名 Dockerfile文件路径
	cmd := exec.Command("docker", "build", "--no-cache", "-t", imageNameAndTag, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RemoveDockerImage(imageNameAndTag string) error {
	cmd := exec.Command("docker", "rmi", imageNameAndTag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CreateDocker(innerPort, dockerName, imageNameAndTag string) (int, error) {
	outerPort, err := portsManager.GetPortManager().AllocatePort()

	if err != nil {
		return -1, err
	}

	_p := strconv.Itoa(outerPort) + ":" + innerPort // 端口映射 宿主机端口:容器内端口
	cmd := exec.Command("docker", "run", "-id", "-p", _p, "--name", dockerName, imageNameAndTag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return outerPort, cmd.Run()
}

func StartDocker(dockerName string) error {
	cmd := exec.Command("docker", "start", dockerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func StopDocker(dockerName string) error {
	cmd := exec.Command("docker", "stop", dockerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RemoveDocker(dockerName string) error {
	cmd := exec.Command("docker", "rm", dockerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

package ipfs

import (
	shell "github.com/ipfs/go-ipfs-api"
	"io"
	"mime/multipart"
	"os"
)

type Ipfs struct {
	Url   string
	Shell *shell.Shell
}

// NewIpfs 初始化ipfs
func NewIpfs(url string) *Ipfs {
	sh := shell.NewShell(url)
	return &Ipfs{
		Url:   url,
		Shell: sh,
	}
}

// UploadIPFS 上传数据到ipfs
func UploadIPFS(file *multipart.FileHeader) (string, error) {
	// 文件类型转换
	uploadFile, errConvert := convertFileHeaderToFile(file)
	if errConvert != nil {
		return "", errConvert
	}

	// TODO 增加负载均衡算法
	ipfsNew := NewIpfs("127.0.0.1:5001") // ipfs私网中任一个节点均可

	// 上传文件
	hash, uploadErr := ipfsNew.Shell.Add(uploadFile) // 使用IPFS API上传文件
	return hash, uploadErr
}

// DownloadIPFS 从IPFS文件系统中下载文件，返回值是临时文件的路径，使用该函数后必须删除临时文件
func DownloadIPFS(hash string) (string, error) {
	ipfsNew := NewIpfs("127.0.0.1:5001") // ipfs私网中任一个节点均可，后期增加负载均衡算法

	ipfsDir, err := os.MkdirTemp("", "ipfs_")
	if err != nil {
		return "", err
	}
	// 获取文件
	errGet := ipfsNew.Shell.Get(hash, ipfsDir)
	return ipfsDir, errGet
}

// convertFileHeaderToFile 文件类型转换，multipart.FileHeader --> os.File
func convertFileHeaderToFile(fileHeader *multipart.FileHeader) (*os.File, error) {
	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "tempfile")
	if err != nil {
		return nil, err
	}

	// 将文件内容拷贝到临时文件中
	_, err = io.Copy(tempFile, file)
	if err != nil {
		tempFile.Close()
		return nil, err
	}

	// 将临时文件指针移到文件开头
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		tempFile.Close()
		return nil, err
	}

	return tempFile, nil
}

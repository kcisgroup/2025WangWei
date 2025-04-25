package productService

import (
	"CipProject/application/common"
	"CipProject/application/dao/productDao"
	"CipProject/application/domain"
	"math"
)

// New 创建软件服务
func New(docker domain.DockerInfo, subRecordsId int64) error {
	err := productDao.Save(docker, subRecordsId)
	return err
}

// UpdateById 更新Docker信息
func UpdateById(docker domain.DockerInfo) error {
	err := productDao.UpdateById(docker)
	return err
}

// Update 更新Docker信息及订阅表使用时长term
func Update(docker domain.DockerInfo, sub domain.SubRecords) error {
	err := productDao.Update(docker, sub)
	return err
}

// GetDetails 服务详情
func GetDetails(id int64) (domain.DockerDetailsDTO, error) {
	dockerDto, err := productDao.GetDetails(id)
	dockerDto.CurrentTime = common.GetCurrentTime()
	dockerDto.Cost = math.Round(dockerDto.Term*float64(dockerDto.RealPrice)*100) / 100
	dockerDto.Term = math.Round(dockerDto.Term*24*100) / 100
	return dockerDto, err
}

// GetByOwner 使用者查看服务列表
func GetByOwner(owner string) ([]domain.DockerInfoDTO, error) {
	dockerDto, err := productDao.GetByOwner(owner)
	return dockerDto, err
}

// GetRawByOwner 使用者查看未创建服务列表
func GetRawByOwner(owner string) ([]domain.RawProductDTO, error) {
	raw, err := productDao.GetRawByOwner(owner)
	return raw, err
}

// GetDestroyList 使用者查看待销毁服务列表
func GetDestroyList(username string) ([]domain.DockerDestroy, error) {
	dockerInfo, err := productDao.GetDestroyList(username)
	return dockerInfo, err
}

// GetById 根据id查询软件服务
func GetById(id int64) (domain.DockerInfo, error) {
	dockerInfo, err := productDao.GetById(id)
	return dockerInfo, err
}

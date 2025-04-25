package productDao

import (
	"CipProject/application/dao"
	"CipProject/application/domain"
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

// Save 容器信息入库
func Save(docker domain.DockerInfo, subRecordsId int64) error {
	db := dao.GetDB()
	tx := db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err := tx.Create(&docker).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(domain.SubRecords{}).Where("id", subRecordsId).Update("docker_id", docker.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// GetByOwner 使用者查看服务列表
func GetByOwner(username string) ([]domain.DockerInfoDTO, error) {
	var dockers []domain.DockerInfoDTO
	db := dao.GetDB()
	result := db.Model(&domain.DockerInfo{}).Select("docker_infos.id, docker_infos.cpi, docker_infos.status, "+
		"docker_infos.ip, docker_infos.port, docker_infos.created_time, code_records.filename, code_records.version, "+
		"code_records.tag").Joins("left join code_records on docker_infos.cpi = code_records.cpi").Where("docker_infos.owner = ?", username).Scan(&dockers)
	if result.Error != nil {
		return []domain.DockerInfoDTO{}, result.Error
	} else if result.RowsAffected == 0 {
		return []domain.DockerInfoDTO{}, errors.New("无查询结果")
	} else {
		return dockers, nil
	}
}

// GetDestroyList 使用者查看待销毁服务列表
func GetDestroyList(username string) ([]domain.DockerDestroy, error) {
	var dockers []domain.DockerDestroy
	db := dao.GetDB()
	result := db.Model(&domain.DockerInfo{}).Select("docker_infos.id, docker_infos.docker_name, "+
		"docker_infos.cpi, ROUND(sub_records.term * 24, 2) as term, sub_records.cost").Joins("left join sub_records "+
		"on docker_infos.id = sub_records.docker_id").Where("docker_infos.status = ? AND "+
		"docker_infos.owner = ?", domain.DockerToDestroy, username).Scan(&dockers)
	if result.Error != nil {
		return []domain.DockerDestroy{}, result.Error
	} else if result.RowsAffected == 0 {
		return []domain.DockerDestroy{}, errors.New("无查询结果")
	} else {
		return dockers, nil
	}
}

// GetRawByOwner 使用者查看未创建的服务
func GetRawByOwner(username string) ([]domain.RawProductDTO, error) {
	var raw []domain.RawProductDTO
	db := dao.GetDB()
	result := db.Model(&domain.SubRecords{}).Select("sub_records.cpi, sub_records.start_time, code_records.filename, "+
		"code_records.version, code_records.tag, code_records.language").Joins("left join code_records "+
		"on sub_records.cpi = code_records.cpi").Where("sub_records.party_a=? AND sub_records.docker_id=? "+
		"AND sub_records.status=?", username, 0, domain.SubApproved).Scan(&raw)
	if result.Error != nil {
		return []domain.RawProductDTO{}, result.Error
	} else if result.RowsAffected == 0 {
		return []domain.RawProductDTO{}, errors.New("无查询结果")
	} else {
		return raw, nil
	}
}

// GetById 根据id查询软件服务
func GetById(id int64) (domain.DockerInfo, error) {
	var docker = domain.DockerInfo{ID: id}
	db := dao.GetDB()

	result := db.First(&docker)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) == true {
		return domain.DockerInfo{}, errors.New("此组件服务不存在")
	} else if result.Error != nil {
		return domain.DockerInfo{}, result.Error
	} else {
		return docker, nil
	}
}

// GetDetails 服务详情
func GetDetails(id int64) (domain.DockerDetailsDTO, error) {
	var dockerDto = domain.DockerDetailsDTO{ID: id}
	db := dao.GetDB()
	tx := db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	err := tx.Model(&domain.DockerInfo{}).Select("cpi", "owner", "ip", "port", "docker_name", "status", "created_time").Where("id = ?", id).First(&dockerDto).Error
	if err != nil {
		tx.Rollback()
		return dockerDto, err
	}

	err = tx.Model(&domain.SubRecords{}).Select("term", "real_price", "cost", "party_b as author").Where("docker_id = ?", id).First(&dockerDto).Error
	if err != nil {
		tx.Rollback()
		return domain.DockerDetailsDTO{}, err
	}

	err = tx.Model(&domain.CodeRecords{}).Select("filename", "version", "tag", "description", "company").Where("cpi = ?", dockerDto.Cpi).First(&dockerDto).Error
	if err != nil {
		tx.Rollback()
		return domain.DockerDetailsDTO{}, err
	}

	tx.Commit()
	return dockerDto, err
}

// UpdateById 更新docker容器的信息
func UpdateById(docker domain.DockerInfo) error {
	db := dao.GetDB()
	err := db.Save(&docker).Error
	return err
}

// Update 更新docker状态和订阅表使用时长term
func Update(docker domain.DockerInfo, sub domain.SubRecords) error {
	db := dao.GetDB()
	tx := db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err := tx.Save(&docker).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(domain.SubRecords{}).Where("id", sub.ID).Update("term", sub.Term).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

package subscribeDao

import (
	"CipProject/application/dao"
	"CipProject/application/domain"
	"database/sql"
	"errors"
)

// Save 新增订阅表
func Save(sr domain.SubRecords) error {
	db := dao.GetDB()
	if err := db.Create(&sr).Error; err != nil {
		return err
	}
	return nil
}

// UpdateById 更新订阅申请表信息
func UpdateById(sr domain.SubRecords) error {
	db := dao.GetDB()
	result := db.Save(&sr)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

// Update 更新订阅表和Docker实例状态
func Update(sr domain.SubRecords) error {
	db := dao.GetDB()
	tx := db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err := tx.Save(&sr).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&domain.DockerInfo{}).Where("id = ?", sr.DockerId).Update("status", domain.DockerToDestroy).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// GetByPartyA 获取甲方（使用者）参与的所有订阅记录
func GetByPartyA(username string) ([]domain.SubInfoDTO, error) {
	var srs []domain.SubInfoDTO
	db := dao.GetDB()
	result := db.Model(&domain.SubRecords{}).Select("sub_records.id, sub_records.code_id, sub_records.start_time, "+
		"sub_records.end_time, sub_records.created_time, sub_records.updated_time, sub_records.real_price, "+
		"sub_records.party_a,sub_records.party_b,sub_records.status, sub_records.note, sub_records.docker_id, code_records.filename, "+
		"code_records.version, code_records.tag, code_records.author").Joins("LEFT JOIN code_records "+
		"ON sub_records.code_id = code_records.id").Where("sub_records.party_a = ?", username).Scan(&srs)

	if result.Error != nil {
		return []domain.SubInfoDTO{}, result.Error
	} else if result.RowsAffected == 0 {
		return []domain.SubInfoDTO{}, errors.New("无订阅记录")
	} else {
		return srs, nil
	}
}

// GetByPartyB 获取乙方（创作者）参与的所有订阅记录
func GetByPartyB(username string) ([]domain.SubRecords, error) {
	var srs []domain.SubRecords
	db := dao.GetDB()
	err := db.Where("party_b = ?", username).Find(&srs).Error
	return srs, err
}

// GetById 根据id获取订阅信息
func GetById(id int64) (domain.SubRecords, error) {
	var sr = domain.SubRecords{ID: id}
	db := dao.GetDB()
	err := db.First(&sr).Error
	return sr, err
}

// GetIDByCpi 根据cpi获取订阅信息
func GetIDByCpi(cpi string) (domain.SubRecords, error) {
	var sr domain.SubRecords
	db := dao.GetDB()
	err := db.Select("id, cpi").Where("cpi = ?", cpi).First(&sr).Error
	return sr, err
}

// GetByDocker 根据DockerId获取订阅id
func GetByDocker(id int64) (domain.SubRecords, error) {
	var sr domain.SubRecords
	db := dao.GetDB()
	err := db.Where("docker_id = ?", id).First(&sr).Error
	return sr, err
}

// GetTermByDocker 根据docker_id获取使用时长
func GetTermByDocker(id int64) (float64, error) {
	var sr domain.SubRecords
	db := dao.GetDB()
	if err := db.Select("term").Where("docker_id = ?", id).First(&sr).Error; err != nil {
		return 0, err
	}
	return sr.Term, nil
}

package codeDao

import (
	"CipProject/application/dao"
	"CipProject/application/domain"
	"errors"
	"fmt"
)

// SaveDraft 数字对象入库
func SaveDraft(ipRecords domain.CodeRecords) error {
	db := dao.GetDB()
	err := db.Create(&ipRecords).Error
	return err
}

// UpdateDraftById 更新代码项目元数据
func UpdateDraftById(ipRecords domain.CodeRecords) error {
	db := dao.GetDB()
	result := db.Save(&ipRecords)
	err := result.Error
	return err
}

// UpdateByCpi 根据cpi更新代码产权表的非零字段
func UpdateByCpi(cipInfo domain.CodeRecords) error {
	db := dao.GetDB()
	result := db.Model(&cipInfo).Where("cpi = ?", cipInfo.Cpi).Updates(cipInfo)
	err := result.Error
	return err
}

// DeleteDraftById 根据id删除数据
func DeleteDraftById(id int64) error {
	db := dao.GetDB()
	err := db.Delete(&domain.CodeRecords{}, id).Error
	return err
}

// DeleteByCpi 根据cpi删除产权元数据
func DeleteByCpi(cpi string) error {
	db := dao.GetDB()
	err := db.Where("cpi = ?", cpi).Delete(&domain.CodeRecords{}).Error
	return err
}

// GetDraftByAuthor 根据创作者用户名查询其拥有的草稿设计
func GetDraftByAuthor(author string) ([]domain.CodeRecords, error) {
	var cipRecords []domain.CodeRecords
	db := dao.GetDB()
	result := db.Where("author = ? AND status = ?", author, domain.CodeNoCip).Find(&cipRecords)

	if result.Error != nil {
		return []domain.CodeRecords{}, result.Error
	} else if result.RowsAffected == 0 {
		return []domain.CodeRecords{}, errors.New("无查询结果")
	} else {
		return cipRecords, nil
	}
}

// GetById 根据id查询
func GetById(id int64) (domain.CodeRecords, error) {
	var cr = domain.CodeRecords{ID: id}
	db := dao.GetDB()
	err := db.First(&cr).Error
	return cr, err
}

// GetByFilename 创作者根据产权名模糊查询
func GetByFilename(cip domain.CodeRecords) ([]domain.CodeRecords, error) {
	var crs []domain.CodeRecords
	db := dao.GetDB()

	condition := "%" + cip.Filename + "%"
	result := db.Where("filename LIKE ? AND author = ?", condition, cip.Author).Not("status = ?", domain.CodeNoCip).Find(&crs)
	if result.Error != nil {
		return []domain.CodeRecords{}, result.Error
	} else if result.RowsAffected == 0 {
		return []domain.CodeRecords{}, errors.New("无匹配的代码产权")
	} else {
		return crs, nil
	}
}

// GetByCpi 根据代码产权数字摘要查询
func GetByCpi(cpi string) (domain.CodeRecords, error) {
	var cr domain.CodeRecords
	db := dao.GetDB()
	err := db.Where("cpi = ?", cpi).First(&cr).Error
	return cr, err
}

// GetPortByCpi 根据cpi查询cpi和innerPort
func GetPortByCpi(cpi string) (domain.CodeRecords, error) {
	var cr domain.CodeRecords
	db := dao.GetDB()
	err := db.Select("cpi, inner_port").Where("cpi = ?", cpi).First(&cr).Error
	return cr, err
}

// GetByAuthor 根据创作者用户名查询其拥有的已授权的代码知识产权
func GetByAuthor(username string) ([]domain.CodeRecords, error) {
	var cipRecords []domain.CodeRecords
	db := dao.GetDB()

	result := db.Where("author = ?", username).Not("status = ?", domain.CodeNoCip).Find(&cipRecords)
	if result.Error != nil {
		return []domain.CodeRecords{}, result.Error
	} else if result.RowsAffected == 0 {
		return []domain.CodeRecords{}, errors.New("无查询结果")
	} else {
		return cipRecords, nil
	}
}

// GetByDesc 使用者根据功能模糊查询开放产权
func GetByDesc(desc string) ([]domain.CodeRecords, error) {
	var cips []domain.CodeRecords
	db := dao.GetDB()
	condition := fmt.Sprintf("%s%s%s", "%", desc, "%")

	result := db.Where("status = ? AND description LIKE ?", domain.CodeOpen, condition).Find(&cips)
	if result.Error != nil {
		return []domain.CodeRecords{}, result.Error
	} else if result.RowsAffected == 0 {
		return []domain.CodeRecords{}, errors.New("无查询结果")
	} else {
		return cips, nil
	}
}

// GetOpenCips 查看所有的待交易的产品的元数据信息
func GetOpenCips() ([]domain.CodeRecords, error) {
	var cips []domain.CodeRecords
	db := dao.GetDB()

	result := db.Where(&domain.CodeRecords{Status: domain.CodeOpen}).Find(&cips)
	if result.Error != nil {
		return []domain.CodeRecords{}, result.Error
	} else if result.RowsAffected == 0 {
		return []domain.CodeRecords{}, errors.New("未查询到记录")
	} else {
		return cips, nil
	}
}

// Count 统计
func Count() domain.CountDto {
	var data domain.CountDto
	db := dao.GetDB()

	db.Model(&domain.CodeRecords{}).Select("COUNT(DISTINCT author, company)").Scan(&data.TotalUsers)
	db.Model(&domain.CodeRecords{}).Count(&data.TotalDrafts)
	db.Model(&domain.CodeRecords{}).Where("status != ?", domain.CodeNoCip).Count(&data.TotalCips)
	db.Model(&domain.CodeRecords{}).Where("status = ?", domain.CodeOpen).Count(&data.TotalOpenCips)
	db.Model(&domain.DockerInfo{}).Distinct("cpi").Count(&data.TotalDockers)
	return data
}

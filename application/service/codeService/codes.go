package codeService

import (
	"CipProject/application/dao/codeDao"
	"CipProject/application/domain"
	companya_sdk "CipProject/application/sdk/companya"
	companyb_sdk "CipProject/application/sdk/companyb"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"strconv"
)

// SaveDraft 将数字对象元数据存储至Mysql中
func SaveDraft(draft domain.CodeRecords) error {
	err := codeDao.SaveDraft(draft) // 调用DAO层
	return err
}

// Authorize 授予代码产权
func Authorize(cr domain.CodeRecords, cip domain.CodeIntellectualProperty, role string) error {
	// 产权元数据上链
	chaincodeName := "authorizeCc"
	fcn := "authorize"
	args := [][]byte{[]byte(cip.Cpi), []byte(cip.Filename), []byte(cip.Version), []byte(cip.Tag), []byte(cip.Author), []byte(cip.RegTime),
		[]byte(cip.UpdatedTime), []byte(cip.Language), []byte(cip.Price), []byte(cip.Status), []byte(cip.Description), []byte(cip.Addr)}
	var err error
	err = nil
	switch role {
	case "companya":
		_, err = companya_sdk.ChannelExecute(chaincodeName, fcn, args)
	case "companyb":
		_, err = companyb_sdk.ChannelExecute(chaincodeName, fcn, args)
	}

	if err == nil || error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		// 更新Mysql表数据
		_ = codeDao.UpdateDraftById(cr)
		return nil
	} else {
		return err
	}
}

// UpdateCip 更新产权元数据
func UpdateCip(cipInfo domain.CodeIntellectualProperty, role string) error {
	// 更新链上数据
	chaincodeName := "authorizeCc"
	fcn := "update"
	args := [][]byte{[]byte(cipInfo.Cpi), []byte(cipInfo.UpdatedTime), []byte(cipInfo.Price), []byte(cipInfo.Status)}
	var err error
	err = nil
	switch role {
	case "companya":
		_, err = companya_sdk.ChannelExecute(chaincodeName, fcn, args)
	case "companyb":
		_, err = companyb_sdk.ChannelExecute(chaincodeName, fcn, args)
	}

	if err == nil || error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		price, _ := strconv.Atoi(cipInfo.Price)
		var cipRecords = domain.CodeRecords{
			Cpi:         cipInfo.Cpi,
			UpdatedTime: cipInfo.UpdatedTime,
			Status:      cipInfo.Status,
			Price:       price,
		}
		err = codeDao.UpdateByCpi(cipRecords)
		if err != nil {
			return errors.New("更新链下数据失败")
		}
		return nil
	} else {
		return errors.New("更新区块链数据失败")
	}
}

// UpdateDraft 更新代码项目元数据
func UpdateDraft(cr domain.CodeRecords) error {
	err := codeDao.UpdateDraftById(cr)
	return err
}

// DeleteDraftById 删除草稿设计的元数据
func DeleteDraftById(id int64) error {
	err := codeDao.DeleteDraftById(id)
	return err
}

// DeleteCipByCpi 根据cpi删除产权元数据
func DeleteCipByCpi(cpi, role string) error {
	// 删除链上产权的元数据
	chaincodeName := "authorizeCc"
	fcn := "delete"
	args := [][]byte{[]byte(cpi)}

	var err error
	err = nil
	switch role {
	case "companya":
		_, err = companya_sdk.ChannelExecute(chaincodeName, fcn, args)
	case "companyb":
		_, err = companyb_sdk.ChannelExecute(chaincodeName, fcn, args)
	}

	if err == nil || error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		// 删除链下Mysql产权元数据
		err = codeDao.DeleteByCpi(cpi)
		return err
	} else {
		return errors.New("删除失败")
	}
}

// QueryDraftBox 查看草稿箱
func QueryDraftBox(author string) ([]domain.CodeRecords, error) {
	drafts, err := codeDao.GetDraftByAuthor(author)
	return drafts, err
}

// GetFromChain 根据cpi查询链上产权元数据
func GetFromChain(cpi, role string) (domain.CodeIntellectualProperty, error) {
	chaincodeName := "authorizeCc"
	fcn := "query"
	args := [][]byte{[]byte(cpi)}

	var err error
	var rsp channel.Response
	err = nil

	switch role {
	case "companya":
		rsp, err = companya_sdk.ChannelQuery(chaincodeName, fcn, args)
	case "companyb":
		rsp, err = companyb_sdk.ChannelQuery(chaincodeName, fcn, args)
	}

	var cip domain.CodeIntellectualProperty
	if err != nil {
		return domain.CodeIntellectualProperty{}, errors.New("产权不存在")
	} else {
		cip.Cpi = cpi
		_ = json.Unmarshal(rsp.Payload, &cip)
		return cip, nil
	}
}

// GetById 根据id查询项目的元数据信息
func GetById(id int64) (domain.CodeRecords, error) {
	cr, err := codeDao.GetById(id)
	return cr, err
}

// GetByFilename 创作者根据`产权名`模糊查询代码产权的元数据信息
func GetByFilename(cip domain.CodeRecords) ([]domain.CodeRecords, error) {
	crs, err := codeDao.GetByFilename(cip)
	return crs, err
}

// GetByCpi 根据cpi查询项目的元数据信息
func GetByCpi(cpi string) (domain.CodeRecords, error) {
	cr, err := codeDao.GetByCpi(cpi)
	return cr, err
}

// GetCipByAuthor 根据创作者用户名查询全部代码知识产权
func GetCipByAuthor(author string) ([]domain.CodeRecords, error) {
	cipRecords, err := codeDao.GetByAuthor(author)
	return cipRecords, err
}

// GetPortByCpi 根据cpi查询组件内部访问端口号
func GetPortByCpi(cpi string) (domain.CodeRecords, error) {
	cr, err := codeDao.GetPortByCpi(cpi)
	return cr, err
}

// GetByDesc 根据功能需求查询项目
func GetByDesc(desc string) ([]domain.CodeRecords, error) {
	cips, err := codeDao.GetByDesc(desc)
	return cips, err
}

// GetOpenCips 查看所有开放项目
func GetOpenCips() ([]domain.CodeRecords, error) {
	cips, err := codeDao.GetOpenCips()
	return cips, err
}

// GetCount 统计项目数
func GetCount() domain.CountDto {
	return codeDao.Count()
}

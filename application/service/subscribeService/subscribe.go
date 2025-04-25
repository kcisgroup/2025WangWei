package subscribeService

import (
	"CipProject/application/dao/subscribeDao"
	"CipProject/application/domain"
	companya_sdk "CipProject/application/sdk/companya"
	companyb_sdk "CipProject/application/sdk/companyb"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// GetFromChain 获取软件订阅合同（链上）
func GetFromChain(id, role string) (domain.SubContract, error) {
	chaincodeName := "subscribeCc"
	fcn := "query"
	args := [][]byte{[]byte(id)}

	var err error = nil
	var rsp channel.Response

	switch role {
	case "companya":
		rsp, err = companya_sdk.ChannelQuery(chaincodeName, fcn, args)
	case "companyb":
		rsp, err = companyb_sdk.ChannelQuery(chaincodeName, fcn, args)
	}

	var subContract domain.SubContract
	if err != nil {
		return subContract, errors.New("合同查询失败")
	} else {
		subContract.ID = id
		err = json.Unmarshal(rsp.Payload, &subContract)
		if err != nil {
			return subContract, errors.New("合同查询失败:JSON数据传输失败")
		} else {
			return subContract, nil
		}
	}
}

// New 新增订阅申请表（MySQL）
func New(sr domain.SubRecords) error {
	err := subscribeDao.Save(sr)
	return err
}

// Agree 甲方确认通过订阅条款，生成订阅合同（链上+MySQL）
func Agree(sr domain.SubRecords, contract domain.SubContract, role string) error {
	// 更新链上数据
	chaincodeName := "subscribeCc"
	fcn := "save"
	args := [][]byte{[]byte(contract.ID), []byte(contract.Cpi), []byte(contract.PartyA), []byte(contract.PartyB), []byte(contract.Start),
		[]byte(contract.End), []byte(contract.RealPrice), []byte(contract.Cost), []byte(contract.Status), []byte(contract.Note)}

	var err error = nil
	switch role {
	case "companya":
		_, err = companya_sdk.ChannelExecute(chaincodeName, fcn, args)
	case "companyb":
		_, err = companyb_sdk.ChannelExecute(chaincodeName, fcn, args)
	}

	if err == nil || error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		errUpdate := subscribeDao.UpdateById(sr)
		return errUpdate
	} else {
		return err
	}
}

// Terminate 终止订阅
func Terminate(sr domain.SubRecords, contract domain.SubContract, role string) error {
	// 更新链上数据
	chaincodeName := "subscribeCc"
	fcn := "update"
	args := [][]byte{[]byte(contract.ID), []byte(contract.End), []byte(contract.Cost), []byte(contract.Status)}
	var err error = nil
	switch role {
	case "companya":
		_, err = companya_sdk.ChannelExecute(chaincodeName, fcn, args)
	case "companyb":
		_, err = companyb_sdk.ChannelExecute(chaincodeName, fcn, args)
	}

	if err == nil || error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		err = subscribeDao.Update(sr)
		return err
	} else {
		return errors.New("订阅数据上链失败")
	}
}

// Pay 付款
func Pay(sr domain.SubRecords) error {
	return subscribeDao.UpdateById(sr)
}

// GetByPartyA 获取甲方（使用者）参与的所有订阅数据
func GetByPartyA(username string) ([]domain.SubInfoDTO, error) {
	srs, err := subscribeDao.GetByPartyA(username)
	return srs, err
}

// GetByPartyB 获取乙方（创作者）参与的所有订阅数据
func GetByPartyB(username string) ([]domain.SubRecords, error) {
	rents, err := subscribeDao.GetByPartyB(username)
	return rents, err
}

// Update 更新SubRecords全字段（MySQL）
func Update(sr domain.SubRecords) error {
	err := subscribeDao.UpdateById(sr)
	return err
}

// GetById 根据id获取订阅信息
func GetById(id int64) (domain.SubRecords, error) {
	rent, err := subscribeDao.GetById(id)
	return rent, err
}

// GetIDByCpi 根据cpi获取订阅信息
func GetIDByCpi(cpi string) (domain.SubRecords, error) {
	rent, err := subscribeDao.GetIDByCpi(cpi)
	return rent, err
}

// GetByDockerId 根据DockerId获取订阅信息
func GetByDockerId(id int64) (domain.SubRecords, error) {
	rent, err := subscribeDao.GetByDocker(id)
	return rent, err
}

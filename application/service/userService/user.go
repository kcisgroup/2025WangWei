package userService

import (
	"CipProject/application/domain"
	accountservice_sdk "CipProject/application/sdk/accountservice"
	"encoding/json"
	"errors"
	"fmt"
)

// GetUser 根据username查询用户数据（链上查询）
func GetUser(username string) (domain.User, error) {
	chaincodeName := "userCc"
	fcn := "query"
	args := [][]byte{[]byte(username)}

	rsp, sdkErr := accountservice_sdk.ChannelQuery(chaincodeName, fcn, args) // 账号管理服务节点负责执行该操作

	var user domain.User
	if sdkErr != nil {
		return user, errors.New("该用户不存在")
	} else {
		user.Username = username
		err := json.Unmarshal(rsp.Payload, &user)
		if err != nil {
			return user, errors.New("用户查询失败:JSON数据传输错误")
		} else {
			return user, nil
		}
	}
}

// PutRecord 备案（链上）
func PutRecord(user domain.User) error {
	chaincodeName := "userCc"
	fcn := "putOnRecord"
	args := [][]byte{[]byte(user.Role), []byte(user.Username)}

	_, err := accountservice_sdk.ChannelExecute(chaincodeName, fcn, args) //用户账号管理服务节点负责操作

	if err == nil || error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		return nil
	} else { //备案失败
		//return errors.New("该工号已备案，请勿重复操作")
		return err
	}
}

// SaveUser 新增用户（FOAC链）
func SaveUser(user domain.User) error {
	chaincodeName := "userCc"
	fcn := "register"
	args := [][]byte{[]byte(user.Username), []byte(user.Pwd), []byte(user.Name), []byte(user.Role), []byte(user.Tel), []byte(user.Idnum), []byte(user.Addr), []byte(user.Gender), []byte(user.Status)}

	putOnRecordKey := fmt.Sprintf("%s-%s", user.Role, user.Username) // 根据role-username判断该工号是否备案，拒绝非法注册
	fmt.Println("================", putOnRecordKey)

	_, err := accountservice_sdk.ChannelExecute(chaincodeName, fcn, args) // 用户账号管理服务节点负责操作

	if err == nil || error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		return nil
	} else { // 注册失败
		return errors.New("该工号已被注册或未备案")
	}
}

// UpdateUser 更新账户部分信息（链上）
func UpdateUser(user domain.User) error {
	chaincodeName := "userCc"
	fcn := "update"
	args := [][]byte{[]byte(user.Username), []byte(user.Pwd), []byte(user.Name), []byte(user.Tel), []byte(user.Addr), []byte(user.Gender)}

	_, err := accountservice_sdk.ChannelExecute(chaincodeName, fcn, args) // 用户账号管理服务节点负责操作

	if err == nil || error.Error(err) == "Client Status Code: (5) TIMEOUT. Description: request timed out or been cancelled" {
		return nil
	} else {
		return errors.New("更新失败")
	}
}

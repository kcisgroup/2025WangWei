package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type AuthorizeCc struct {
}

func (p *AuthorizeCc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 12 {
		return shim.Error("初始化智能合约时，必须包含十二个参数")
	}

	// 代码产权唯一标识码
	cpi := args[0]
	// 组件名（代码知识产权名）
	filename := args[1]
	// 版本号
	version := args[2]
	// 标签
	tag := args[3]
	// 创作者username
	author := args[4]
	// 授权时间
	regTime := args[5]
	// 更新时间
	updatedTime := args[6]
	// 程序设计语言
	language := args[7]
	// 租赁价格（元/天）
	price := args[8]
	// 状态，0 禁用；1 开放
	status := args[9]
	// 功能描述
	description := args[10]
	// 文件链下存储地址
	addr := args[11]

	// 设定key
	key := fmt.Sprintf("%s", cpi)
	// 设定value
	value := fmt.Sprintf(`{"filename":"%s","version":"%s","tag":"%s","author":"%s","regTime":"%s","updatedTime":"%s","language":"%s","price":"%s","status":"%s","description":"%s","addr":"%s"}`, filename, version, tag, author, regTime, updatedTime, language, price, status, description, addr)
	// 上链
	err := stub.PutState(key, []byte(value))
	if err != nil {
		return shim.Error("代码知识产权上链失败")
	}
	return shim.Success([]byte("初始化成功"))
}

func (p *AuthorizeCc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	f, args := stub.GetFunctionAndParameters()

	if f == "authorize" { // 授权
		return p.authorize(stub, args)
	}

	if f == "query" { // 查询数字产权的元数据
		return p.query(stub, args)
	}

	if f == "update" { // 更新数字产权的元数据
		return p.update(stub, args)
	}
	if f == "delete" { // 删除数字产权的元数据
		return p.delete(stub, args)
	}

	return shim.Error("函数名称错误，必须是authorize、query、update和delete中的一种")
}

func (p *AuthorizeCc) authorize(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 12 {
		return shim.Error("除调用的函数外，必须包含十二个参数")
	}

	// 代码产权唯一标识码
	cpi := args[0]
	// 组件名（代码知识产权名）
	filename := args[1]
	// 版本号
	version := args[2]
	// 标签
	tag := args[3]
	// 创作者username
	author := args[4]
	// 授权时间
	regTime := args[5]
	// 更新时间
	updatedTime := args[6]
	// 程序设计语言
	language := args[7]
	// 租赁价格（元/天）
	price := args[8]
	// 状态，0 禁用；1 开放
	status := args[9]
	// 功能描述
	description := args[10]
	// 文件链下存储地址
	addr := args[11]

	// 设定value
	value := fmt.Sprintf(`{"filename":"%s","version":"%s","tag":"%s","author":"%s","regTime":"%s","updatedTime":"%s","language":"%s","price":"%s","status":"%s","description":"%s","addr":"%s"}`, filename, version, tag, author, regTime, updatedTime, language, price, status, description, addr)

	// 查询是否已经存在相同的cpi
	rsp, _ := stub.GetState(cpi)
	if string(rsp) != "" {
		return shim.Error("此组件已被确权，请勿重复申请")
	}

	// 写入区块链
	err := stub.PutState(cpi, []byte(value))
	if err != nil {
		return shim.Error("在代码知识产权上链过程中出现异常，导致确权失败，请尝试重新申请或联系相关工作人员")
	}
	return shim.Success([]byte("授权成功"))
}

func (p *AuthorizeCc) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("查询代码知识产权时，必须上次代码产权唯一标识码")
	}

	cpi := args[0]

	re, err := stub.GetState(cpi)
	if err != nil {
		return shim.Error("代码知识产权查询失败")
	}
	if (string(re) == "") || (re == nil) {
		return shim.Error("您查询的代码知识产权不存在")
	}

	return shim.Success(re)
}

func (p *AuthorizeCc) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("更新版权时，必须上传四个参数")
	}

	cpi := args[0]
	updatedTime := args[1]
	price := args[2]
	status := args[3]

	// 获取代码产权元数据
	re, err := stub.GetState(cpi)
	if err != nil {
		return shim.Error("产权信息修改失败")
	}
	if (string(re) == "") || (re == nil) {
		return shim.Error("此项代码知识产权不存在")
	}

	data := make(map[string]interface{})
	_ = json.Unmarshal(re, &data) // []byte转map

	// 更新
	data["updatedTime"] = updatedTime
	data["price"] = price
	data["status"] = status
	updateValue, _ := json.Marshal(data) // map转[]byte

	// 写入区块链
	updateErr := stub.PutState(cpi, updateValue)
	if updateErr != nil {
		return shim.Error("产权信息修改失败")
	}
	return shim.Success([]byte("产权信息修改成功"))
}

func (p *AuthorizeCc) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("删除版权时，必须上传代码产权唯一标识码")
	}

	cpi := args[0]

	err := stub.DelState(cpi)
	if err != nil {
		return shim.Error("代码知识产权删除失败")
	}

	return shim.Success([]byte("代码知识产权删除成功"))
}

func main() {
	err := shim.Start(new(AuthorizeCc))

	if err != nil {
		fmt.Println("启动fabric发生错误")
	}
}

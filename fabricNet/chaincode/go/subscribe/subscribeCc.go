package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SubscribeCc struct {
}

func (p *SubscribeCc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 10 {
		return shim.Error("初始化智能合约时，必须上传十个参数")
	}

	// 合同编号
	id := args[0]
	// 代码产权唯一标识码
	cpi := args[1]
	// 甲方username(使用者)
	partyA := args[2]
	// 乙方username(创作者)
	partyB := args[3]
	// 合同生效时间
	start := args[4]
	// 截止时间
	end := args[5]
	// 协商后价格（元/天）
	realPrice := args[6]
	// 费用
	cost := args[7]
	// 交易状态
	status := args[8]
	// 备注
	note := args[9]

	// 设定value
	value := fmt.Sprintf(`{"cpi":"%s","partyA":"%s","partyB":"%s","start":"%s","end":"%s","realPrice":"%s","cost":"%s","status":"%s","note":"%s"}`, cpi, partyA, partyB, start, end, realPrice, cost, status, note)

	// 写入区块链
	err := stub.PutState(id, []byte(value))
	if err != nil {
		return shim.Error("初始化失败")
	}

	return shim.Success([]byte("初始化成功"))
}

func (p *SubscribeCc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()

	if f == "save" { // 生成订阅合同
		return p.save(stub, args)
	}

	if f == "query" { // 查询订阅合同
		return p.query(stub, args)
	}

	if f == "update" { // 更新订阅合同
		return p.update(stub, args)
	}

	return shim.Error("函数名称错误")
}

func (p *SubscribeCc) save(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 10 {
		return shim.Error("生成订阅合同时，必须上传十个参数")
	}

	// 订阅合同编号
	id := args[0]
	// 代码产权唯一标识码
	cpi := args[1]
	// 甲方username(使用者)
	partyA := args[2]
	// 乙方username(创作者)
	partyB := args[3]
	// 订阅生效时间
	start := args[4]
	// 截止时间
	end := args[5]
	// 协商后价格（元/天）
	realPrice := args[6]
	// 费用
	cost := args[7]
	// 交易状态
	status := args[8]
	// 备注
	note := args[9]

	rsp, _ := stub.GetState(id) // 查询订阅合同是否已经存在
	if string(rsp) != "" {
		return shim.Error("此订阅已存在，请勿重复上传")
	}

	// 设定value
	value := fmt.Sprintf(`{"cpi":"%s","partyA":"%s","partyB":"%s","start":"%s","end":"%s","realPrice":"%s","cost":"%s","status":"%s","note":"%s"}`, cpi, partyA, partyB, start, end, realPrice, cost, status, note)

	// 合同写入区块链
	err := stub.PutState(id, []byte(value))
	if err != nil {
		return shim.Error("订阅合同上链失败")
	}

	return shim.Success([]byte("订阅合同上链成功"))
}

func (p *SubscribeCc) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("查询合同时，必须上传订阅合同编号")
	}

	id := args[0]

	re, err := stub.GetState(id)
	if err != nil {
		return shim.Error("订阅合同查询失败")
	}
	if (string(re) == "") || (re == nil) {
		return shim.Error("您查询的订阅合同不存在")
	}

	return shim.Success(re)

}

func (p *SubscribeCc) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("更新订阅合同时，必须上传四个数据")
	}

	id := args[0]
	end := args[1]
	cost := args[2]
	status := args[3]

	// 获取订阅合同信息
	sub, err := stub.GetState(id)
	if err != nil {
		return shim.Error("订阅合同更新失败")
	}
	if (string(sub) == "") || (sub == nil) {
		return shim.Error("订阅合同不存在，更新失败")
	}

	data := make(map[string]interface{})
	_ = json.Unmarshal(sub, &data)

	data["end"] = end
	data["cost"] = cost
	data["status"] = status
	updateValue, _ := json.Marshal(data)

	// 写入区块链
	updateErr := stub.PutState(id, updateValue)
	if updateErr != nil {
		return shim.Error("订阅合同更新失败")
	}

	return shim.Success([]byte("链上订阅合同更新成功"))
}

func main() {
	err := shim.Start(new(SubscribeCc))

	if err != nil {
		fmt.Println("启动fabric失败" + err.Error())
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type UserCc struct {
}

func (p *UserCc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 9 {
		return shim.Error("调用的函数外，必须上传九个参数")
	}

	//用户名，默认使用所在公司的备案工号
	username := args[0]
	//密码
	pwd := args[1]
	//姓名
	name := args[2]
	//所属组织/公司
	role := args[3]
	//联系方式
	tel := args[4]
	//身份证号
	idnum := args[5]
	//地址
	addr := args[6]
	//性别
	gender := args[7]
	//0各公司管理员，1普通员工用户
	status := args[8]

	//设定key:username
	key := fmt.Sprintf("%s", username)

	//设定value(json格式)
	value := fmt.Sprintf(`{"pwd":"%s","name":"%s","role":"%s","tel":"%s","idnum":"%s","addr":"%s","gender":"%s","status":"%s"}`, pwd, name, role, tel, idnum, addr, gender, status)

	err := stub.PutState(key, []byte(value))
	if err != nil {
		return shim.Error("初始化失败")
	}

	return shim.Success([]byte("初始化成功"))
}

func (p *UserCc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	f, args := stub.GetFunctionAndParameters()

	if f == "register" { //用户注册
		return p.register(stub, args)
	}

	if f == "query" { //用户信息查询，验证
		return p.query(stub, args)
	}

	if f == "putOnRecord" { //工号备案
		return p.putOnRecord(stub, args)
	}

	if f == "update" { // 更新用户信息
		return p.update(stub, args)
	}

	if f == "delete" { // 删除用户和备案信息
		return p.delete(stub, args)
	}

	return shim.Error("函数名称错误，必须是register、query、putOnRecord和update中一种")
}

func (p *UserCc) query(stub shim.ChaincodeStubInterface, args []string) pb.Response { //用户查询
	if len(args) != 1 {
		return shim.Error("查询用户时，必须上传用户名")
	}

	username := args[0]

	value, err := stub.GetState(username) // 根据username查询
	if err != nil {
		return shim.Error("用户信息查询失败")
	}
	if (string(value) == "") || (value == nil) {
		return shim.Error("您查询的用户信息不存在")
	}

	return shim.Success(value)
}

func (p *UserCc) register(stub shim.ChaincodeStubInterface, args []string) pb.Response { //用户注册
	if len(args) != 9 {
		return shim.Error("调用的函数外，必须上传九个参数")
	}

	// 用户名，默认使用所在公司的备案工号
	username := args[0]
	// 密码
	pwd := args[1]
	// 姓名
	name := args[2]
	// 所属组织/公司
	role := args[3]
	// 联系方式
	tel := args[4]
	// 身份证号
	idnum := args[5]
	// 地址
	addr := args[6]
	// 性别
	gender := args[7]
	// 0各公司管理员，1普通员工用户
	status := args[8]

	rsp, _ := stub.GetState(username) // 查询该工号是否已被注册
	if string(rsp) != "" {            // key不存在也不会报错，但是返回的空的
		return shim.Error("此用户名已被注册")
	}

	putOnRecordKey := fmt.Sprintf("%s-%s", role, username) // 根据role-username判断该工号是否备案，拒绝非法注册
	putOnRecordRsp, _ := stub.GetState(putOnRecordKey)
	if string(putOnRecordRsp) == "" {
		return shim.Error("此工号未备案")
	}

	//设定value(json格式)
	value := fmt.Sprintf(`{"pwd":"%s","name":"%s","role":"%s","tel":"%s","idnum":"%s","addr":"%s","gender":"%s","status":"%s"}`, pwd, name, role, tel, idnum, addr, gender, status)

	err := stub.PutState(username, []byte(value)) // 用户信息上链
	if err != nil {
		return shim.Error("注册失败")
	}

	return shim.Success([]byte("注册成功"))
}

func (p *UserCc) putOnRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response { //员工信息备案，由各公司管理员操作
	if len(args) != 2 {
		return shim.Error("备案工号时，必须上传组织名和工号")
	}

	// 组织
	role := args[0]
	// 工号
	username := args[1]

	key := fmt.Sprintf("%s-%s", role, username)

	rsp, _ := stub.GetState(key) // 验证此工号是否已备案
	if string(rsp) != "" {
		return shim.Error("此工号已存在，请勿重复备案")
	}

	value := fmt.Sprintf(`{"role":"%s","username":"%s","putOnRecord":"%s"}`, role, username, "success")
	err := stub.PutState(key, []byte(value)) //备案信息写入区块链
	if err != nil {
		return shim.Error("备案失败")
	}

	return shim.Success([]byte("备案成功"))
}

func (p *UserCc) update(stub shim.ChaincodeStubInterface, args []string) pb.Response { // 更新用户信息
	if len(args) != 6 {
		return shim.Error("更新用户信息时，必须上传六个参数")
	}

	// 用户名
	username := args[0]
	// 密码
	pwd := args[1]
	// 姓名
	name := args[2]
	// 手机号
	tel := args[3]
	// 地址
	addr := args[4]
	// 性别
	gender := args[5]

	re, err := stub.GetState(username) // 获取产权信息
	if err != nil {
		return shim.Error("更新失败")
	}
	if (string(re) == "") || (re == nil) {
		return shim.Error("此用户不存在，更新失败")
	}

	data := make(map[string]interface{})
	_ = json.Unmarshal(re, &data)

	// 更新
	data["pwd"] = pwd
	data["name"] = name
	data["tel"] = tel
	data["addr"] = addr
	data["gender"] = gender
	updateValue, _ := json.Marshal(data) // map转[]byte

	// 新数据上链
	err = stub.PutState(username, updateValue)
	if err != nil {
		return shim.Error("更新失败")
	}

	return shim.Success([]byte("更新成功"))
}

func (p *UserCc) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response { // 删除普通员工用户
	if len(args) != 2 {
		return shim.Error("删除用户时，必须上传两个参数")
	}

	// 用户名
	username := args[0]
	// 组织
	role := args[1]

	// 删除员工信息
	err := stub.DelState(username)
	if err != nil {
		return shim.Error("用户信息删除失败")
	}
	// 删除备案记录
	key := fmt.Sprintf("%s-%s", role, username)
	err = stub.DelState(key)
	if err != nil {
		return shim.Error("备案记录清除失败")
	}

	return shim.Success([]byte("用户信息删除成功"))
}

func main() {
	err := shim.Start(new(UserCc))

	if err != nil {
		fmt.Println("启动fabric发生错误")
	}
}

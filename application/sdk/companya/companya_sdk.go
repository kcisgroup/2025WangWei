package companya_sdk

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"time"
)

var (
	companya_sdk *fabsdk.FabricSDK
	config_path  string = "../application/conf/sdk_conf/companya-conf.yaml"
	channel_id   string = "masterchannel"
	p0a          string = "peer0.companya.cole.com"
	org          string = "CompanyaOrg"
	user         string = "Admin"
	err          error
)

func init() {
	companya_sdk, err = fabsdk.New(config.FromFile(config_path))

	if err != nil {
		ret_err := fmt.Sprintf("读取配置文件出错，错误信息：", err)
		panic(ret_err)
	}
}

func ChannelExecute(chaincode_name, execute_fcn string, args [][]byte) (channel.Response, error) {
	ctx := companya_sdk.ChannelContext(channel_id, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := channel.New(ctx)

	if err != nil {
		panic("实例化通道出错")
	}
	rsp, err := cli.Execute(channel.Request{
		ChaincodeID: chaincode_name,
		Fcn:         execute_fcn,
		Args:        args,
	}, channel.WithTargetEndpoints(p0a), channel.WithTimeout(fab.TimeoutType(3), time.Second*2))

	return rsp, err

}

func ChannelQuery(chaincode_name, query_fcn string, args [][]byte) (channel.Response, error) {
	ctx := companya_sdk.ChannelContext(channel_id, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := channel.New(ctx)

	if err != nil {
		ret_err := fmt.Sprintf("实例化通道出错:%v", err)
		panic(ret_err)
	}

	rsp, err := cli.Query(channel.Request{
		ChaincodeID: chaincode_name,
		Fcn:         query_fcn,
		Args:        args,
	}, channel.WithTargetEndpoints(p0a))
	return rsp, err

}

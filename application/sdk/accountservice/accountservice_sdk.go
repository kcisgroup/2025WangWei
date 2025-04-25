package accountservice_sdk

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"time"
)

var (
	accountservice_sdk *fabsdk.FabricSDK
	config_path        = "../application/conf/sdk_conf/accountservice-conf.yaml"
	channel_id         = "accountchannel"
	accountservice     = "peer0.accountservice.cole.com"
	org                = "AccountserviceOrg"
	user               = "Admin"
	err                error
)

func init() {
	accountservice_sdk, err = fabsdk.New(config.FromFile(config_path))

	if err != nil {
		ret_err := fmt.Sprintf("读取配置文件出错，错误信息：", err)
		panic(ret_err)
	}
}

func ChannelExecute(chaincodeName, executeFcn string, args [][]byte) (channel.Response, error) {
	ctx := accountservice_sdk.ChannelContext(channel_id, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := channel.New(ctx)

	if err != nil {
		panic("实例化通道出错")
	}
	rsp, err := cli.Execute(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         executeFcn,
		Args:        args,
	}, channel.WithTargetEndpoints(accountservice), channel.WithTimeout(fab.TimeoutType(3), time.Second*2))

	return rsp, err
}

func ChannelQuery(chaincode_name, query_fcn string, args [][]byte) (channel.Response, error) {
	ctx := accountservice_sdk.ChannelContext(channel_id, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := channel.New(ctx)

	if err != nil {
		ret_err := fmt.Sprintf("实例化通道出错:%v", err)
		panic(ret_err)
	}

	rsp, err := cli.Query(channel.Request{
		ChaincodeID: chaincode_name,
		Fcn:         query_fcn,
		Args:        args,
	}, channel.WithTargetEndpoints(accountservice))

	return rsp, err
}

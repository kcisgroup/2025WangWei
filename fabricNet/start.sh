#!/bin/bash
echo "第一步：下载镜像开始..."

#./download-dockerimages.sh   # 这个下载镜像有问题，换成下面命令

#docker pull hyperledger/fabric-tools:x86_64-1.0.0
#docker pull hyperledger/fabric-orderer:x86_64-1.0.0
#docker pull hyperledger/fabric-peer:x86_64-1.0.0
#docker pull hyperledger/fabric-couchdb:x86_64-1.0.0
#docker pull hyperledger/fabric-kafka:x86_64-1.0.0
#docker pull hyperledger/fabric-ca:x86_64-1.0.0
#docker pull hyperledger/fabric-ccenv:x86_64-1.0.0
#docker pull hyperledger/fabric-baseimage:x86_64-0.4.7
#docker pull hyperledger/fabric-javaenv:x86_64-1.0.0
#docker pull hyperledger/fabric-zookeeper:x86_64-1.0.0

#docker tag hyperledger/fabric-tools:x86_64-1.0.0 hyperledger/fabric-tools
#docker tag hyperledger/fabric-orderer:x86_64-1.0.0 hyperledger/fabric-orderer
#docker tag hyperledger/fabric-peer:x86_64-1.0.0 hyperledger/fabric-peer
#docker tag hyperledger/fabric-couchdb:x86_64-1.0.0 hyperledger/fabric-couchdb
#docker tag hyperledger/fabric-kafka:x86_64-1.0.0 hyperledger/fabric-kafka
#docker tag hyperledger/fabric-ca:x86_64-1.0.0 hyperledger/fabric-ca
#docker tag hyperledger/fabric-ccenv:x86_64-1.0.0 hyperledger/fabric-ccenv
#docker tag hyperledger/fabric-baseimage:x86_64-0.4.7 hyperledger/fabric-baseimage
#docker tag hyperledger/fabric-javaenv:x86_64-1.0.0 hyperledger/fabric-javaenv
#docker tag hyperledger/fabric-zookeeper:x86_64-1.0.0 hyperledger/fabric-zookeeper
echo "下载镜像完成"

echo "第二步:生成证书开始..."
cryptogen generate --config ./crypto-config.yaml
echo "生成证书结束"

echo "第三步:生成创世块及通道开始..."
mkdir channel-artifacts
configtxgen -profile GenGenesis -outputBlock ./channel-artifacts/genesis.block

configtxgen -profile GenChannel -channelID masterchannel -outputCreateChannelTx ./channel-artifacts/masterchannel.tx
configtxgen -profile GenChannel -channelID accountchannel -outputCreateChannelTx ./channel-artifacts/accountchannel.tx
echo "生成创世块及通道完成"

echo "第四步:启动docker-compose开始..."
docker-compose -f ./docker-compose-cli.yaml up -d
echo "启动docker-compose结束"

# 启动docker-compose需要点时间
sleep 5

echo "第六步:创建通道开始..."
docker exec cli_companya peer channel create -o orderer.cole.com:7050 -c masterchannel -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/masterchannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem
docker exec cli_accountservice peer channel create -o orderer.cole.com:7050 -c accountchannel -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/accountchannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem
echo  "创建通道结束"

echo "拷贝生成的block文件到宿主机"
#注意路径：将容器内masterchannel.block拷贝至宿主机项目IP-blockchain-project/fabricNet/
docker cp cli_companya:/opt/gopath/src/github.com/hyperledger/fabric/peer/masterchannel.block /home/ww/goworkspace/src/CipProject/fabricNet/
docker cp cli_accountservice:/opt/gopath/src/github.com/hyperledger/fabric/peer/accountchannel.block /home/ww/goworkspace/src/CipProject/fabricNet/

echo "拷贝block文件到新的宿主机"
#注意路径：将宿主机项目IP-blockchain-project/fabricNet/拷贝至容器中
docker cp /home/ww/goworkspace/src/CipProject/fabricNet/masterchannel.block cli_companya1:/opt/gopath/src/github.com/hyperledger/fabric/peer/
docker cp /home/ww/goworkspace/src/CipProject/fabricNet/masterchannel.block cli_companyb:/opt/gopath/src/github.com/hyperledger/fabric/peer/
docker cp /home/ww/goworkspace/src/CipProject/fabricNet/masterchannel.block cli_companyb1:/opt/gopath/src/github.com/hyperledger/fabric/peer/
docker cp /home/ww/goworkspace/src/CipProject/fabricNet/accountchannel.block cli_accountservice:/opt/gopath/src/github.com/hyperledger/fabric/peer/
docker cp /home/ww/goworkspace/src/CipProject/fabricNet/accountchannel.block cli_accountservice1:/opt/gopath/src/github.com/hyperledger/fabric/peer/

echo "第七步:当前peer节点加入通道开始..."
docker exec cli_companya peer channel join -b masterchannel.block
docker exec cli_companya1 peer channel join -b masterchannel.block

docker exec cli_companyb peer channel join -b masterchannel.block
docker exec cli_companyb1 peer channel join -b masterchannel.block

docker exec cli_accountservice peer channel join -b accountchannel.block
docker exec cli_accountservice1 peer channel join -b accountchannel.block
echo "当前peer节点加入结束"


echo "第八步:链码安装开始..."
docker exec cli_companya peer chaincode install -n authorizeCc -v 1.0 -p github.com/hyperledger/fabric/chaincode/go/authorization
docker exec cli_companya1 peer chaincode install -n authorizeCc -v 1.0 -p github.com/hyperledger/fabric/chaincode/go/authorization

docker exec cli_companyb peer chaincode install -n authorizeCc -v 1.0 -p github.com/hyperledger/fabric/chaincode/go/authorization
docker exec cli_companyb1 peer chaincode install -n authorizeCc -v 1.0 -p github.com/hyperledger/fabric/chaincode/go/authorization

docker exec cli_companya peer chaincode install -n subscribeCc -v 1.0 -p github.com/hyperledger/fabric/chaincode/go/subscribe
docker exec cli_companya1 peer chaincode install -n subscribeCc -v 1.0 -p github.com/hyperledger/fabric/chaincode/go/subscribe

docker exec cli_companyb peer chaincode install -n subscribeCc -v 1.0 -p github.com/hyperledger/fabric/chaincode/go/subscribe
docker exec cli_companyb1 peer chaincode install -n subscribeCc -v 1.0 -p github.com/hyperledger/fabric/chaincode/go/subscribe

docker exec cli_accountservice peer chaincode install -n userCc -v 1.0 -p github.com/hyperledger/fabric/chaincode/go/user
docker exec cli_accountservice1 peer chaincode install -n userCc -v 1.0 -p github.com/hyperledger/fabric/chaincode/go/user
echo "链码安装结束"

echo "第九步:链码初始化开始..."
docker exec cli_companya peer chaincode instantiate -o orderer.cole.com:7050 -C masterchannel -c '{"Args":["init","00000000000000000000000000000000","myFirstIP","0.1","test","231000","2024-01-13 20:30:00","2024-01-13 20:30:00","0","100","0","The function of this CIP is to instantiate chaincode","0000000000000000000000000000000000000000000000"]}' -n authorizeCc -P "OR ('CompanyaMSP.member','CompanybMSP.member')" -v 1.0 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem
echo "订阅合约初始化"
docker exec cli_companya peer chaincode instantiate -o orderer.cole.com:7050 -C masterchannel -c '{"Args":["init","2347164828300000","00000000000000000000000000000000","231000","230000","2024-01-13 20:30:00","2024-01-13 20:30:00","20","10000","费用已支付，交易完成","this is for instantiating chaincode"]}' -n subscribeCc -P "OR ('CompanyaMSP.member','CompanybMSP.member')" -v 1.0 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem
echo "用户合约初始化"
docker exec cli_accountservice peer chaincode instantiate -o orderer.cole.com:7050 -C accountchannel -c '{"Args":["init","241000","12345","testname","testrole","testTel","testidnum","testaddr","testgender","7"]}' -n userCc -P "OR ('AccountserviceMSP.member')" -v 1.0 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem
echo "链码初始化结束"

echo "第十步①:测试authorizeCc开始..."
docker exec cli_companya peer chaincode invoke -C masterchannel -n authorizeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["authorize","0000000011111111222222223333333344444444555555556666666677777777","myFirstIP","0.0.1","test","241000","2024-01-13 20:30:00","2024-01-13 20:30:00","0","100","0","The function of this CIP is to test authorize","Qma5h3TUai3BSow692g6om6Ngj98ukPXKPPGFk00000000"]}'
sleep 2
docker exec cli_companya peer chaincode query -C masterchannel -n authorizeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["query","0000000011111111222222223333333344444444555555556666666677777777"]}'
docker exec cli_companya1 peer chaincode query -C masterchannel -n authorizeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["query","0000000011111111222222223333333344444444555555556666666677777777"]}'
docker exec cli_companyb peer chaincode query -C masterchannel -n authorizeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["query","0000000011111111222222223333333344444444555555556666666677777777"]}'
docker exec cli_companyb1 peer chaincode query -C masterchannel -n authorizeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["query","0000000011111111222222223333333344444444555555556666666677777777"]}'
docker exec cli_companya peer chaincode invoke -C masterchannel -n authorizeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["update","0000000011111111222222223333333344444444555555556666666677777777","1693012644","150","1"]}'
docker exec cli_companya peer chaincode invoke -C masterchannel -n authorizeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["delete","0000000011111111222222223333333344444444555555556666666677777777"]}'
echo "测试authorizeCc结束"

echo "第十步②:测试subscribeCc开始..."
docker exec cli_companya peer chaincode invoke -C masterchannel -n subscribeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["save","8471712900000000","0000000011111111222222223333333344444444555555556666666677777777","231001","231001","2024-02-24 00:00:00","待定","待定","待定","双方已签订合同","测试用例"]}'
sleep 2
docker exec cli_companya peer chaincode query -C masterchannel -n subscribeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["query","8471712900000000"]}'
docker exec cli_companya1 peer chaincode query -C masterchannel -n subscribeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["query","8471712900000000"]}'
docker exec cli_companyb peer chaincode query -C masterchannel -n subscribeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["query","8471712900000000"]}'
docker exec cli_companyb1 peer chaincode invoke -C masterchannel -n subscribeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["update","8471712900000000","2025-04-05 00:00:00","2000","费用已支付，交易完成"]}'
sleep 2
docker exec cli_companyb1 peer chaincode query -C masterchannel -n subscribeCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["query","8471712900000000"]}'
echo "测试subscribeCc结束"

sleep 2
echo "第十步③:测试userCc开始，预置员工用户..."
docker exec cli_accountservice peer chaincode invoke -C accountchannel -n userCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["putOnRecord","companya","231001"]}'
sleep 2
docker exec cli_accountservice peer chaincode invoke -C accountchannel -n userCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["register","231001","12345","员工甲","companya","13785968572","320412200110100008","江苏常州","男","1"]}'
sleep 2
docker exec cli_accountservice1 peer chaincode query -C accountchannel -n userCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["query","231001"]}'

docker exec cli_accountservice peer chaincode invoke -C accountchannel -n userCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["putOnRecord","companyb","231002"]}'
sleep 2
docker exec cli_accountservice peer chaincode invoke -C accountchannel -n userCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["register","231002","12345","员工乙","companyb","13785968573","320112200010100016","江苏南京","男","1"]}'
sleep 2
docker exec cli_accountservice1 peer chaincode query -C accountchannel -n userCc -o orderer.cole.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/cole.com/msp/tlscacerts/tlsca.cole.com-cert.pem -c '{"Args":["query","231002"]}'

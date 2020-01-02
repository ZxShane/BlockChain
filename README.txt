
测试运行环境为CentOS7 
以下是按照步骤来启动Fabric集群与测试
为了方便 我们编写了./byfn.sh脚本文件 




----------------------------------------------------------------------------
             docker 初始化
docker-compose -f docker-compose-cli.yaml down --volumes --remove-orphans

docker rm -f $(docker ps -a | grep "hyperledger/*" | awk "{print \$1}")

docker volume prune

docker ps -a | grep Exited | awk '{print $1}'

docker stop $(docker ps -q)

docker rm $(docker ps -aq)

--------------------------------------------------------------------------------
通过crypto-config.yaml配置文件的配置项去生成对应的组织的节点的用户证书：
cryptogen generate --config=crypto-config.yaml --output ./crypto-config

生成创始区块到genesis.block文件中
configtxgen -profile myGenesis -channelID mychannel -outputBlock ./channel-artifacts/genesis.block

生成通道文件
configtxgen -profile myChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel


更新youai组织的锚节点
configtxgen -profile myChannel -outputAnchorPeersUpdate ./channel-artifacts/youaiAnchor.tx -channelID mychannel -asOrg youaiMSP

更新renai组织的锚节点
configtxgen -profile myChannel -outputAnchorPeersUpdate ./channel-artifacts/renaiAnchor.tx -channelID mychannel -asOrg renaiMSP


-------------------------------------------------------------------------------------------------
启动容器：
docker-compose -f docker-compose-cli.yaml up -d

显示容器的状态：
docker-compose -f docker-compose-cli.yaml ps

进入容器内部
docker exec -it cli bash

创建通道：
tlsfile=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer channel create -o orderer.example.com:7050 --tls true --cafile $tlsfile -c mychannel -f ./channel-artifacts/channel.tx

加入通道
peer channel join -b mychannel.block 
展示现有的通道
peer channel list

--------------------------------------------------------------------------------------------------------------
  安装 部署 链代码
peer chaincode install -n MER -v 1.0 -p github.com/MER

peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n MER -v 1.0 -c '{"Args":[""]}' -P "OR ('org1MSP.member','org2MSP.member')"


----------------------------------------------------------------------------------------------------------------
录入患者的基本信息
peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n MER -c '{"function":"create","Args":["1509050119","201907131440","zhouxian","m","19961019","hanzu","shanxi","1","1"]}'
查询患者基本信息
peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n MER -c '{"function":"isAllowQueryUserContent","Args":["1509050119"]}'

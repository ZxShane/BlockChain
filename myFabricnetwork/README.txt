
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
cryptogen generate --config=crypto-config.yaml

生成创始区块到genesis.block文件中
configtxgen -profile myGenesis -outputBlock ./channel-artifacts/genesis.block

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

tlsfile=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/c4BlockChain.com/orderers/orderer.c4BlockChain.com/msp/tlscacerts/tlsca.c4BlockChain.com-cert.pem

peer channel create -o orderer.c4BlockChain.com:7050 --tls true --cafile $tlsfile -c mychannel -f ./channel-artifacts/channel.tx

加入通道
peer channel join -b mychannel.block 
展示现有的通道
peer channel list

--------------------------------------------------------------------------------------------------------------
  安装链代码 通过打包连代码的方式来安装
peer chaincode package -n MER -v 1.0 -p github.com/chaincode/MER MER.package
peer chaincode install MER.package


初始化链代码
tlsfile=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/c4BlockChain.com/orderers/orderer.c4BlockChain.com/msp/tlscacerts/tlsca.c4BlockChain.com-cert.pem

peer chaincode instantiate -o orderer.c4BlockChain.com:7050 --tls --cafile $tlsfile -C mychannel -n MER -v 1.0 -c '{"Args":[""]}'


----------------------------------------------------------------------------------------------------------------
录入患者的基本信息
peer chaincode invoke -o orderer.c4BlockChain.com:7050 --tls --cafile $tlsfile -C mychannel -n MER -c '{"function":"create","Args":["1509050119","201907131440","zhouxian","m","19961019","hanzu","shanxi","1","1"]}'
查询患者基本信息
peer chaincode query -o orderer.c4BlockChain.com:7050 --tls --cafile $tlsfile -C mychannel -n MER -c '{"function":"isAllowQueryUserContent","Args":["1509050119"]}'


###########################################################
   注意！！！！！！
   上述的创建通道以及加入通道、安装连代码都是在youai组织的peer0节点进行的
   所以还需要进行组织、节点切换 使renai组织的peer0节点加入到通道中 并安装相同的连代码
   （此处的我们设计的fabric集群只有两个组织，一个组织内对应1个节点 若复杂组织则需要进行较多的切换）
###########################################################

  切换到组织renai的peer0节点
CORE_PEER_ADDRESS=peer0.renai.c4BlockChain.com:7051
CORE_PEER_LOCALMSPID=renaiMSP
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/renai.c4BlockChain.com/users/Admin@renai.c4BlockChain.com/msp

CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/renai.c4BlockChain.com/peers/peer0.renai.c4BlockChain.com/tls/ca.crt

CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/renai.c4BlockChain.com/peers/peer0.renai.c4BlockChain.com/tls/server.crt

CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/renai.c4BlockChain.com/peers/peer0.renai.c4BlockChain.com/tls/server.key

peer channel join -b mychannel.block
peer chaincode install MER.package

############################################################
至此我们将链代码部署到了两个节点 并将这两个组织（节点）加入到了通道当中

最后在renai组织的peer0节点查询患者基本信息
peer chaincode query -o orderer.c4BlockChain.com:7050 --tls --cafile $tlsfile -C mychannel -n MER -c '{"function":"isAllowQueryUserContent","Args":["1509050119"]}'




##############################################################
 新链代码的调试
##############################################################
docker exec -it cli bash

peer chaincode package -n MER -v 1.0.4 -p github.com/chaincode/MER MER.package
peer chaincode install MER.package

tlsfile=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/c4BlockChain.com/orderers/orderer.c4BlockChain.com/msp/tlscacerts/tlsca.c4BlockChain.com-cert.pem

peer chaincode instantiate -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -v 1.0.4 -c '{"Args":[""]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"aestest","Args":["0","zxshane","zhouxian"]}'


peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"createPatient","Args":["1509050119","201907131440","zhouxian","m","19961019","hanzu","shanxi","weihun"]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"createPatient","Args":["zxshane", "123", "kq71nTlrh4APQ4dtU42cyg==", "DyUDs+RJYRdaikaaqjucIw==", "Qe7AsPycgfq7wfmDavUKow==", "/Vp8jujDdOdHZY7eXfu+7g==", "fHELiv3xQjB1vBWDcQs/Yg==", "fHELiv3xQjB1vBWDcQs/Yg==", "uua9I7uNH+nZRrWMpP2FWw=="]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"patientregister","Args":["1234aaaa","huanzhe","123456","2345","123","34656"]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"patientLogin","Args":["zxshane","123","123"]}'

peer chaincode query -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"queryPaientBaseinfo","Args":["1509050119"]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"addMedicalContent","Args":["1509050119","01","xiaohua","02","futong","aaaaaaaaaaaaaa","chihuaiduzi","wu","wu","wu"]}'

peer chaincode query -o orderer.c4BlockChain.com:7050 --tls --cafile $tlsfile -C mychannel -n MER -c '{"function":"queryMedicalByID","Args":["1509050119","201907301002"]}'


peer chaincode query -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"queryMedicalNum","Args":["1607010308"]}'

peer chaincode query -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"queryReserverInfoNum","Args":["1509050119"]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"addReserveInfo","Args":["1509050119","01","zhouxian","m","19961019","hanzu","shanxi","weihun"]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"deleteReserverInfo","Args":["1509050119"]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"RSAPublicEncryptoAES","Args":["zxshane"]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"changeReserverState","Args":["1509050119"，"01"]}'

peer chaincode query -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"queryDoctorInfoByID","Args":["01"]}'

peer chaincode query -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"queryReserverInfoByID","Args":["1509050119","01"]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"createDoctorAuto","Args":[""]}'

peer chaincode invoke -o orderer.c4BlockChain.com:7050 -C mychannel -n MER -c '{"function":"transferPermission","Args":["01","1607010308"]}'

-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAkPxJUAbmaD+km3GdnenxoDg88jKb+VstkA2MbfdrRTZ3ezdtxtYlJypfBgZS
1BKfZMFKeIwNOqs+OzuKILQzUwZB4cVyYVa1TPq/DpJFEN/rrHgsSsf2VdwEHTP/onhXcZd3
zD6EZH3WO80ueMItcPbVVh5ZdZdA/UHVi9n5APo5tBUxnx8uyGau+M710iRK85Hs4vQ6pT12
MOOVAoGbtRW00gDA0BD6kBYjGFOnpUNH+5Xq84prfPrbThGI2X9CqaqzhA4HcuKzQL/ohF1q
5TH71kD2VfNz/TkOIE1ArHIgWRdAb3uVlcO7ti685QO8REzRP28+Bw3/8Lj90It6+wIDAQAB
-----END RSA PUBLIC KEY-----

-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAkPxJUAbmaD+km3GdnenxoDg88jKb+VstkA2MbfdrRTZ3ezdtxtYlJypf
BgZS1BKfZMFKeIwNOqs+OzuKILQzUwZB4cVyYVa1TPq/DpJFEN/rrHgsSsf2VdwEHTP/onhX
cZd3zD6EZH3WO80ueMItcPbVVh5ZdZdA/UHVi9n5APo5tBUxnx8uyGau+M710iRK85Hs4vQ6
pT12MOOVAoGbtRW00gDA0BD6kBYjGFOnpUNH+5Xq84prfPrbThGI2X9CqaqzhA4HcuKzQL/o
hF1q5TH71kD2VfNz/TkOIE1ArHIgWRdAb3uVlcO7ti685QO8REzRP28+Bw3/8Lj90It6+wID
AQABAoIBAEGciijXFony0zEtN2DxL9GL4bjRQliT9IiOORDCuR63SVbPfLRQ0Ltqp1n4np8u
VkeoWWU4K/xy5lSz2wx1wAxAdqwPSHXYYW+WwcN8WhkK3IJOV+z3lPjB+nKkx3jk8N2M6D/b
wtofQEYL0o3/gcTvTxgL3/whGN4DXvpNCCxyYKrD8dg5sSzGuzcBFFelqEnhfEP3l0y+n4ts
mHHx9V/9ltsS2fSK+zZedi0hIm7ftPG4iiWfv8vJrY3bwBpvvlIeg7g/pyqwijP0y4pYmPdR
UXPfEfATlfDwKNNNESiFRuIel8I+2xeVfXC7XxUKF8aRWJbXW8Bs1todXZL8NXECgYEA29I+
CPItcu8KVg38Ldgu9dco5dKdPlmQqHE7wdev5yJ1RUSph5OTntkBUqi+AMe6l9FVSmvckTsA
6yeuir4HqZZEvZED0b8hyhgwJ/5j68VJ3uNHjpXibholfr/e+FnVedrv+yjHHIGHkUbKsozD
4MbUPEJroI4uHr5Wkk1fWw8CgYEAqNj6oUtBGWDlpM9Nh7K/G8QF84a74Sv63drlZ6Dha8Xs
59aXlxXzFHl0GpYE3esi53/ZZIVQnhk4yGVr73dcqRAca3SJ4NCRnlEoFtkdQVJ4DJuMKecR
nBv1a4rZGtqyM0gWt0tPG4qxIlqZOxNxwlB0mfVmwQ3V/nM9zozF0VUCgYAtALDUkgf99LQw
A/Lxy8VpbSAhVOn+PsXfxjbOq4KGlkZd5P20FOFu7sxXiNZFQJ6RwDhu4QAp92NrwRb5rofR
D0OJb6vRgAjB4AvT1D/On/hMmkknBsZxdgbhGRTj0ThkFw90YtfInTgM5OpQfYMIfIwsvghc
uV71yk/c6dwwvwKBgBqV54ij/8EON7pmha+bHmoxyDoa+dQvh5WNFNfnRfchN/cdG8tHQnnz
0asp+eQzVNCcmc8xCouKLx2mkoMnCSj5h3AH7nm+fV8vKh/G2ctiP9LEXyJt5qDs6gyf1SVc
T/ixHhqIOhF9GfztxPi/TAcrgeCH+kDle89Pt+ig07jtAoGBAKBoxbdFctJyiBVxJEIvkanr
6+7gfdKJE4W6ueEaV8s8mhDM9RpRUfa/xoqs6KfLH/yU2UXXsiskxSXt8cyvKxtICnRYHVju
xLTMJU2zanYvb9yd/AqFCja/+dNxRJlsRntwQT9bTzm5DNvIZ3LDwVPS7XAdnhGChLdo2I0m
MBWI
-----END RSA PRIVATE KEY-----

6JetqVgn2B2508IT
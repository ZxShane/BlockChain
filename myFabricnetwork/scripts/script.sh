#!/bin/bash

echo
echo " ____    _____      _      ____    _____ "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _|"
echo "\___ \    | |     / _ \   | |_) |   | |  "
echo " ___) |   | |    / ___ \  |  _ <    | |  "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|  "
echo
echo "Build myFabricNetwork "
echo
CHANNEL_NAME="$1"
DELAY="$2"
LANGUAGE="$3"
TIMEOUT="$4"
VERBOSE="$5"
: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${LANGUAGE:="golang"}
: ${TIMEOUT:="10"}
: ${VERBOSE:="false"}
LANGUAGE=`echo "$LANGUAGE" | tr [:upper:] [:lower:]`
COUNTER=1
MAX_RETRY=5
CC_SRC_PATH="github.com/chaincode/MER"

echo "Channel name : "$CHANNEL_NAME

# import utils
. scripts/utils.sh

createChannel() {
	setGlobals 0 1

	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
                set -x
		peer channel create -o orderer.c4BlockChain.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx >&log.txt
		res=$?
                set +x
	else
				set -x
		peer channel create -o orderer.c4BlockChain.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
		res=$?
				set +x
	fi
	cat log.txt
	verifyResult $res "Channel creation failed"
	echo "===================== Channel '$CHANNEL_NAME' created ===================== "
	echo
}

joinChannel () {
   joinChannelWithRetry 0 1
   echo "===================== peer0.youai joined channel '$CHANNEL_NAME' ===================== "
   sleep $DELAY
   echo

   joinChannelWithRetry 1 1
   echo "===================== peer1.youai joined channel '$CHANNEL_NAME' ===================== "
   sleep $DELAY
   echo

   joinChannelWithRetry 0 2
   echo "===================== peer0.renai joined channel '$CHANNEL_NAME' ===================== "
   sleep $DELAY
   echo

   joinChannelWithRetry 1 2
   echo "===================== peer1.renai joined channel '$CHANNEL_NAME' ===================== "
   sleep $DELAY
   echo
   
	
}

## Create channel
echo "Creating channel..."
createChannel

## Join all the peers to the channel
echo "Having all peers join the channel..."
joinChannel

## Set the anchor peers for each org in the channel
echo "Updating anchor peers for youai..."
updateAnchorPeers 0 1
echo "Updating anchor peers for renai..."
updateAnchorPeers 0 2

## Install chaincode on peer0.youai and peer0.renai
#echo "Installing chaincode on peer0.youai..."
#installChaincode 0 1

#echo "Installing chaincode on peer1.youai..."
#installChaincode 1 1

#echo "Installing chaincode on peer0.renai..."
#installChaincode 0 2

#echo "Installing chaincode on peer1.renai..."
#installChaincode 1 2


# Instantiate chaincode on peer0.renai
#echo "Instantiating chaincode on peer0.renai..."
#instantiateChaincode 0 2

echo
echo "========= All GOOD, myFabricNetwork execution completed =========== "
echo

echo
echo " _____   _   _   ____   "
echo "| ____| | \ | | |  _ \  "
echo "|  _|   |  \| | | | | | "
echo "| |___  | |\  | | |_| | "
echo "|_____| |_| \_| |____/  "
echo

exit 0


#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This is a collection of bash functions used by different scripts

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/msp/tlscacerts/tlsca.dfarmorderer.com-cert.pem
PEER0_ORG1_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/tls/ca.crt
PEER0_ORG2_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/tls/ca.crt
PEER0_ORG3_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/tls/ca.crt
PEER0_ORG4_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/tls/ca.crt

# Set OrdererOrg.Admin globals
setOrdererGlobals() {
  CORE_PEER_LOCALMSPID="OrdererMSP"
  CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/msp/tlscacerts/tlsca.dfarmorderer.com-cert.pem
  CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/dfarmorderer.com/users/Admin@dfarmorderer.com/msp
}

# Set environment variables for the peer org
setGlobals() {
  ORG=$1
  if [ $ORG -eq 1 ]; then
    CORE_PEER_LOCALMSPID="DfarmadminMSP"
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/dfarmadmin.com/users/Admin@dfarmadmin.com/msp
    CORE_PEER_ADDRESS=peer0.dfarmadmin.com:7051
  elif [ $ORG -eq 2 ]; then
    CORE_PEER_LOCALMSPID="YngadminMSP"
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/yngadmin.com/users/Admin@yngadmin.com/msp
    CORE_PEER_ADDRESS=peer0.yngadmin.com:9051
  elif [ $ORG -eq 3 ]; then
    CORE_PEER_LOCALMSPID="DfarmclientMSP"
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/dfarmclient.com/users/Admin@dfarmclient.com/msp
    CORE_PEER_ADDRESS=peer0.dfarmclient.com:11051
  elif [ $ORG -eq 4 ]; then
    CORE_PEER_LOCALMSPID="KmadminMSP"
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG4_CA
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/kmadmin.com/users/Admin@kmadmin.com/msp
    CORE_PEER_ADDRESS=peer0.kmadmin.com:9031
  else
    echo "================== ERROR !!! ORG Unknown =================="
  fi

  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo $'\e[1;31m'!!!!!!!!!!!!!!! $2 !!!!!!!!!!!!!!!!$'\e[0m'
    echo
    exit 1
  fi
}

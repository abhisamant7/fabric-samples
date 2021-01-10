#!/bin/bash

source scriptUtils.sh

function createOrg1() {

  infoln "Enroll the CA admin"
  mkdir -p organizations/peerOrganizations/dfarmadmin.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/dfarmadmin.com/
  #  rm -rf $FABRIC_CA_CLIENT_HOME/fabric-ca-client-config.yaml
  #  rm -rf $FABRIC_CA_CLIENT_HOME/msp

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca-org1 --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/peerOrganizations/dfarmadmin.com/msp/config.yaml

  infoln "Register peer0"
  set -x
  fabric-ca-client register --caname ca-org1 --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Register user"
  set -x
  fabric-ca-client register --caname ca-org1 --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Register the org admin"
  set -x
  fabric-ca-client register --caname ca-org1 --id.name org1admin --id.secret org1adminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  mkdir -p organizations/peerOrganizations/dfarmadmin.com/peers
  mkdir -p organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com

  infoln "Generate the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-org1 -M ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/msp --csr.hosts peer0.dfarmadmin.com --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/dfarmadmin.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/msp/config.yaml

  infoln "Generate the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-org1 -M ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/tls --enrollment.profile tls --csr.hosts peer0.dfarmadmin.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/tls/server.key

  mkdir -p ${PWD}/organizations/peerOrganizations/dfarmadmin.com/msp/tlscacerts
  cp ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/dfarmadmin.com/msp/tlscacerts/ca.crt

  mkdir -p ${PWD}/organizations/peerOrganizations/dfarmadmin.com/tlsca
  cp ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/dfarmadmin.com/tlsca/tlsca.dfarmadmin.com-cert.pem

  mkdir -p ${PWD}/organizations/peerOrganizations/dfarmadmin.com/ca
  cp ${PWD}/organizations/peerOrganizations/dfarmadmin.com/peers/peer0.dfarmadmin.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/dfarmadmin.com/ca/ca.dfarmadmin.com-cert.pem

  mkdir -p organizations/peerOrganizations/dfarmadmin.com/users
  mkdir -p organizations/peerOrganizations/dfarmadmin.com/users/User1@dfarmadmin.com

  infoln "Generate the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca-org1 -M ${PWD}/organizations/peerOrganizations/dfarmadmin.com/users/User1@dfarmadmin.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/dfarmadmin.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/dfarmadmin.com/users/User1@dfarmadmin.com/msp/config.yaml

  mkdir -p organizations/peerOrganizations/dfarmadmin.com/users/Admin@dfarmadmin.com

  infoln "Generate the org admin msp"
  set -x
  fabric-ca-client enroll -u https://org1admin:org1adminpw@localhost:7054 --caname ca-org1 -M ${PWD}/organizations/peerOrganizations/dfarmadmin.com/users/Admin@dfarmadmin.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/dfarmadmin.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/dfarmadmin.com/users/Admin@dfarmadmin.com/msp/config.yaml

}

function createOrg2() {

  infoln "Enroll the CA admin"
  mkdir -p organizations/peerOrganizations/yngadmin.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/yngadmin.com/
  #  rm -rf $FABRIC_CA_CLIENT_HOME/fabric-ca-client-config.yaml
  #  rm -rf $FABRIC_CA_CLIENT_HOME/msp

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8054 --caname ca-org2 --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-org2.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-org2.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-org2.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-org2.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/peerOrganizations/yngadmin.com/msp/config.yaml

  infoln "Register peer0"
  set -x
  fabric-ca-client register --caname ca-org2 --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Register user"
  set -x
  fabric-ca-client register --caname ca-org2 --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Register the org admin"
  set -x
  fabric-ca-client register --caname ca-org2 --id.name org2admin --id.secret org2adminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  mkdir -p organizations/peerOrganizations/yngadmin.com/peers
  mkdir -p organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com

  infoln "Generate the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-org2 -M ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/msp --csr.hosts peer0.yngadmin.com --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/yngadmin.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/msp/config.yaml

  infoln "Generate the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-org2 -M ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/tls --enrollment.profile tls --csr.hosts peer0.yngadmin.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/tls/server.key

  mkdir -p ${PWD}/organizations/peerOrganizations/yngadmin.com/msp/tlscacerts
  cp ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/yngadmin.com/msp/tlscacerts/ca.crt

  mkdir -p ${PWD}/organizations/peerOrganizations/yngadmin.com/tlsca
  cp ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/yngadmin.com/tlsca/tlsca.yngadmin.com-cert.pem

  mkdir -p ${PWD}/organizations/peerOrganizations/yngadmin.com/ca
  cp ${PWD}/organizations/peerOrganizations/yngadmin.com/peers/peer0.yngadmin.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/yngadmin.com/ca/ca.yngadmin.com-cert.pem

  mkdir -p organizations/peerOrganizations/yngadmin.com/users
  mkdir -p organizations/peerOrganizations/yngadmin.com/users/User1@yngadmin.com

  infoln "Generate the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca-org2 -M ${PWD}/organizations/peerOrganizations/yngadmin.com/users/User1@yngadmin.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/yngadmin.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/yngadmin.com/users/User1@yngadmin.com/msp/config.yaml

  mkdir -p organizations/peerOrganizations/yngadmin.com/users/Admin@yngadmin.com

  infoln "Generate the org admin msp"
  set -x
  fabric-ca-client enroll -u https://org2admin:org2adminpw@localhost:8054 --caname ca-org2 -M ${PWD}/organizations/peerOrganizations/yngadmin.com/users/Admin@yngadmin.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/yngadmin.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/yngadmin.com/users/Admin@yngadmin.com/msp/config.yaml

}

function createOrderer() {

  infoln "Enroll the CA admin"
  mkdir -p organizations/ordererOrganizations/dfarmorderer.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/dfarmorderer.com
  #  rm -rf $FABRIC_CA_CLIENT_HOME/fabric-ca-client-config.yaml
  #  rm -rf $FABRIC_CA_CLIENT_HOME/msp

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca-orderer --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/ordererOrganizations/dfarmorderer.com/msp/config.yaml

  infoln "Register orderer"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Register the orderer admin"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  mkdir -p organizations/ordererOrganizations/dfarmorderer.com/orderers
  mkdir -p organizations/ordererOrganizations/dfarmorderer.com/orderers/dfarmorderer.com

  mkdir -p organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com

  infoln "Generate the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/msp --csr.hosts orderer.dfarmorderer.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/msp/config.yaml

  infoln "Generate the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/tls --enrollment.profile tls --csr.hosts orderer.dfarmorderer.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/tls/ca.crt
  cp ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/tls/signcerts/* ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/tls/server.crt
  cp ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/tls/keystore/* ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/tls/server.key

  mkdir -p ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/msp/tlscacerts
  cp ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/msp/tlscacerts/tlsca.dfarmorderer.com-cert.pem

  mkdir -p ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/msp/tlscacerts
  cp ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/orderers/orderer.dfarmorderer.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/msp/tlscacerts/tlsca.dfarmorderer.com-cert.pem

  mkdir -p organizations/ordererOrganizations/dfarmorderer.com/users
  mkdir -p organizations/ordererOrganizations/dfarmorderer.com/users/Admin@dfarmorderer.com

  infoln "Generate the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/users/Admin@dfarmorderer.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/dfarmorderer.com/users/Admin@dfarmorderer.com/msp/config.yaml

}



function createOrg3 {

  echo
	echo "Enroll the CA admin"
  echo
	mkdir -p ../organizations/peerOrganizations/dfarmclient.com/

	export FABRIC_CA_CLIENT_HOME=${PWD}/../organizations/peerOrganizations/dfarmclient.com/
#  rm -rf $FABRIC_CA_CLIENT_HOME/fabric-ca-client-config.yaml
#  rm -rf $FABRIC_CA_CLIENT_HOME/msp

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:11054 --caname ca-org3 --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-org3.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-org3.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-org3.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-org3.pem
    OrganizationalUnitIdentifier: orderer' > ${PWD}/../organizations/peerOrganizations/dfarmclient.com/msp/config.yaml

  echo
	echo "Register peer0"
  echo
  set -x
	fabric-ca-client register --caname ca-org3 --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  echo
  echo "Register user"
  echo
  set -x
  fabric-ca-client register --caname ca-org3 --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  echo
  echo "Register the org admin"
  echo
  set -x
  fabric-ca-client register --caname ca-org3 --id.name org3admin --id.secret org3adminpw --id.type admin --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

	mkdir -p ../organizations/peerOrganizations/dfarmclient.com/peers
  mkdir -p ../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com

  echo
  echo "## Generate the peer0 msp"
  echo
  set -x
	fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-org3 -M ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/msp --csr.hosts peer0.dfarmclient.com --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/../organizations/peerOrganizations/dfarmclient.com/msp/config.yaml ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/msp/config.yaml

  echo
  echo "## Generate the peer0-tls certificates"
  echo
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-org3 -M ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/tls --enrollment.profile tls --csr.hosts peer0.dfarmclient.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null


  cp ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/tls/tlscacerts/* ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/tls/ca.crt
  cp ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/tls/signcerts/* ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/tls/server.crt
  cp ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/tls/keystore/* ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/tls/server.key

  mkdir ${PWD}/../organizations/peerOrganizations/dfarmclient.com/msp/tlscacerts
  cp ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/tls/tlscacerts/* ${PWD}/../organizations/peerOrganizations/dfarmclient.com/msp/tlscacerts/ca.crt

  mkdir ${PWD}/../organizations/peerOrganizations/dfarmclient.com/tlsca
  cp ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/tls/tlscacerts/* ${PWD}/../organizations/peerOrganizations/dfarmclient.com/tlsca/tlsca.dfarmclient.com-cert.pem

  mkdir ${PWD}/../organizations/peerOrganizations/dfarmclient.com/ca
  cp ${PWD}/../organizations/peerOrganizations/dfarmclient.com/peers/peer0.dfarmclient.com/msp/cacerts/* ${PWD}/../organizations/peerOrganizations/dfarmclient.com/ca/ca.dfarmclient.com-cert.pem

  mkdir -p ../organizations/peerOrganizations/dfarmclient.com/users
  mkdir -p ../organizations/peerOrganizations/dfarmclient.com/users/User1@dfarmclient.com

  echo
  echo "## Generate the user msp"
  echo
  set -x
	fabric-ca-client enroll -u https://user1:user1pw@localhost:11054 --caname ca-org3 -M ${PWD}/../organizations/peerOrganizations/dfarmclient.com/users/User1@dfarmclient.com/msp --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/../organizations/peerOrganizations/dfarmclient.com/msp/config.yaml ${PWD}/../organizations/peerOrganizations/dfarmclient.com/users/User1@dfarmclient.com/msp/config.yaml

  mkdir -p ../organizations/peerOrganizations/dfarmclient.com/users/Admin@dfarmclient.com

  echo
  echo "## Generate the org admin msp"
  echo
  set -x
	fabric-ca-client enroll -u https://org3admin:org3adminpw@localhost:11054 --caname ca-org3 -M ${PWD}/../organizations/peerOrganizations/dfarmclient.com/users/Admin@dfarmclient.com/msp --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/../organizations/peerOrganizations/dfarmclient.com/msp/config.yaml ${PWD}/../organizations/peerOrganizations/dfarmclient.com/users/Admin@dfarmclient.com/msp/config.yaml

}

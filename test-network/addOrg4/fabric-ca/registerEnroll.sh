

function createOrg3 {

  echo
	echo "Enroll the CA admin"
  echo
	mkdir -p ../organizations/peerOrganizations/kmadmin.com/

	export FABRIC_CA_CLIENT_HOME=${PWD}/../organizations/peerOrganizations/kmadmin.com/
#  rm -rf $FABRIC_CA_CLIENT_HOME/fabric-ca-client-config.yaml
#  rm -rf $FABRIC_CA_CLIENT_HOME/msp

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9034 --caname ca-kmadmin --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9034-ca-kmadmin.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9034-ca-kmadmin.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9034-ca-kmadmin.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9034-ca-kmadmin.pem
    OrganizationalUnitIdentifier: orderer' > ${PWD}/../organizations/peerOrganizations/kmadmin.com/msp/config.yaml

  echo
	echo "Register peer0"
  echo
  set -x
	fabric-ca-client register --caname ca-kmadmin --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  echo
  echo "Register user"
  echo
  set -x
  fabric-ca-client register --caname ca-kmadmin --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  echo
  echo "Register the org admin"
  echo
  set -x
  fabric-ca-client register --caname ca-kmadmin --id.name org3admin --id.secret org3adminpw --id.type admin --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

	mkdir -p ../organizations/peerOrganizations/kmadmin.com/peers
  mkdir -p ../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com

  echo
  echo "## Generate the peer0 msp"
  echo
  set -x
	fabric-ca-client enroll -u https://peer0:peer0pw@localhost:9034 --caname ca-kmadmin -M ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/msp --csr.hosts peer0.kmadmin.com --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/../organizations/peerOrganizations/kmadmin.com/msp/config.yaml ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/msp/config.yaml

  echo
  echo "## Generate the peer0-tls certificates"
  echo
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:9034 --caname ca-kmadmin -M ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/tls --enrollment.profile tls --csr.hosts peer0.kmadmin.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null


  cp ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/tls/tlscacerts/* ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/tls/ca.crt
  cp ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/tls/signcerts/* ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/tls/server.crt
  cp ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/tls/keystore/* ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/tls/server.key

  mkdir ${PWD}/../organizations/peerOrganizations/kmadmin.com/msp/tlscacerts
  cp ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/tls/tlscacerts/* ${PWD}/../organizations/peerOrganizations/kmadmin.com/msp/tlscacerts/ca.crt

  mkdir ${PWD}/../organizations/peerOrganizations/kmadmin.com/tlsca
  cp ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/tls/tlscacerts/* ${PWD}/../organizations/peerOrganizations/kmadmin.com/tlsca/tlsca.kmadmin.com-cert.pem

  mkdir ${PWD}/../organizations/peerOrganizations/kmadmin.com/ca
  cp ${PWD}/../organizations/peerOrganizations/kmadmin.com/peers/peer0.kmadmin.com/msp/cacerts/* ${PWD}/../organizations/peerOrganizations/kmadmin.com/ca/ca.kmadmin.com-cert.pem

  mkdir -p ../organizations/peerOrganizations/kmadmin.com/users
  mkdir -p ../organizations/peerOrganizations/kmadmin.com/users/User1@kmadmin.com

  echo
  echo "## Generate the user msp"
  echo
  set -x
	fabric-ca-client enroll -u https://user1:user1pw@localhost:9034 --caname ca-kmadmin -M ${PWD}/../organizations/peerOrganizations/kmadmin.com/users/User1@kmadmin.com/msp --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/../organizations/peerOrganizations/kmadmin.com/msp/config.yaml ${PWD}/../organizations/peerOrganizations/kmadmin.com/users/User1@kmadmin.com/msp/config.yaml

  mkdir -p ../organizations/peerOrganizations/kmadmin.com/users/Admin@kmadmin.com

  echo
  echo "## Generate the org admin msp"
  echo
  set -x
	fabric-ca-client enroll -u https://org3admin:org3adminpw@localhost:9034 --caname ca-kmadmin -M ${PWD}/../organizations/peerOrganizations/kmadmin.com/users/Admin@kmadmin.com/msp --tls.certfiles ${PWD}/fabric-ca/org3/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/../organizations/peerOrganizations/kmadmin.com/msp/config.yaml ${PWD}/../organizations/peerOrganizations/kmadmin.com/users/Admin@kmadmin.com/msp/config.yaml

}

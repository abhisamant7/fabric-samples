package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	AGENT_RATE_PER_RECORD_TYPE string = "AgentRatePerAsset"
)

type AgentRateReq struct {
	State         string  `json:"STATE,omitempty"`
	Country       string  `json:"COUNTRY,omitempty"`
	ProduceName   string  `json:"PRODUCE,omitempty"`
	Variety       string  `json:"VARIETY,omitempty"`
	Role          int     `json:"ROLE,omitempty"`
	SourceID      string  `json:"SOURCEID,omitempty"`
	DestinationID string  `json:"DESTINATIONID,omitempty"`
	PerValue      float64 `json:"PERVALUE,omitempty"`
}

//AgentRatePerAsset asset to provide Markup and discount on final Buyer payment for a produce Variety

type AgentRatePerAsset struct {
	DocType       string  `json:"docType,omitempty"`
	SourceID      string  `json:"SOURCEID,omitempty"`
	DestinationID string  `json:"DESTINATIONID,omitempty"`
	Variety       string  `json:"VARIETY,omitempty"`
	ProduceName   string  `json:"PRODUCE,omitempty"`
	QualityType   string  `json:"QUALITYTYPE,omitempty"`
	PerValue      float64 `json:"PERVALUE,omitempty"`
}

// CreateAgentRatePerAsset to define DC rate
func CreateAgentRatePerAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys, dfarmkey []string
	if len(args) < 1 {
		logger.Errorf("CreateAgentRatePerAsset : Incorrect number of arguments.")
		return shim.Error("CreateAgentRatePerAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("CreateAgentRatePerAsset :Attribute approle not found to create ParticipantRateAsset"))
	}
	//Abhijeet add role logic here
	if "MARKETING" != strings.ToUpper(value) {
		return shim.Error(fmt.Sprintf("CreateAgentRatePerAsset :This User is not allowed to create ParticipantRateAsset"))
	}
	// Convert the arg to a AgentRatePerAsset Object
	logger.Infof("CreateAgentRatePerAsset: Arguments for ledgerapi %s : ", args[0])
	asset := AgentRateReq{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	logger.Infof("CreateAgentRatePerAsset :AgentRateReq is : %s ", asset)
	keys = append(keys, asset.DestinationID)
	keys = append(keys, asset.ProduceName)
	keys = append(keys, asset.Variety)
	Avalbytes, err = QueryAsset(stub, AGENT_RATE_PER_RECORD_TYPE, keys)
	if err != nil {
		dfarmkey = append(dfarmkey, asset.Country)
		dfarmkey = append(dfarmkey, asset.State)
		dfarmkey = append(dfarmkey, asset.ProduceName)
		dfarmrate, ok := getdFarmRateAsset(stub, dfarmkey)
		if !ok {
			logger.Errorf("CreateAgentRatePerAsset : Error reading dfram max marketing agent rate  %t", ok)
			return shim.Error(fmt.Sprintf("CreateAgentRatePerAsset : Error reading dfram max marketing agent rate  %t", ok))

		}
		if dfarmrate.MaxAgentPer < asset.PerValue {
			logger.Errorf("CreateAgentRatePerAsset : dFarm Max agent rate:  [ %f ] is less then Agent rate  : [ %f ]", dfarmrate.MaxAgentPer, asset.PerValue)
			return shim.Error(fmt.Sprintf("CreateAgentRatePerAsset : dFarm Max agent rate:  [ %f ] is less then Agent rate  : [ %f ]", dfarmrate.MaxAgentPer, asset.PerValue))

		}
		perasset := AgentRatePerAsset{}
		perasset.DocType = AGENT_RATE_PER_RECORD_TYPE
		perasset.SourceID = asset.SourceID
		perasset.ProduceName = asset.ProduceName
		perasset.DestinationID = asset.DestinationID
		perasset.Variety = asset.Variety
		perasset.PerValue = asset.PerValue
		Avalbytes, _ = json.Marshal(perasset)
		logger.Infof("CreateAgentRatePerAsset :AgentRatePerAsset  is : %s ", Avalbytes)
		//logger.Infof("CreateAgentRatePerAsset : AgentRatePerAsset  percentage is : %f ", perasset.PerValue)

		err = CreateAsset(stub, AGENT_RATE_PER_RECORD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreateAgentRatePerAsset : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreateAgentRatePerAsset : ParticipantRateAsset Object first time create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}
	readobj := AgentRatePerAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &readobj)
	readobj.PerValue = asset.PerValue
	Avalbytes, _ = json.Marshal(readobj)
	logger.Infof("CreateAgentRatePerAsset :AgentRatePerAsset Asset is : %v ", readobj)
	err = UpdateAssetWithoutGet(stub, AGENT_RATE_PER_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateAgentRatePerAsset : Error inserting Object first time  into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateAgentRatePerAsset : AgentRatePerAsset Object first time create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

//getAgentRatePerAsset to get AgentRatePerAsset from ledger based on 3 param
func getAgentRatePerAsset(stub shim.ChaincodeStubInterface, args []string) (AgentRatePerAsset, bool) {

	var err error
	var Avalbytes []byte

	if len(args) == 0 {
		logger.Errorf("getAgentRatePerAsset : Incorrect number of arguments.")
		return AgentRatePerAsset{}, false
	}
	logger.Infof("getAgentRatePerAsset :args is : %v ", args)
	//keys = append(keys, args[0]) //Destination id i.e Buyer ID
	//keys = append(keys, args[1]) //Produce Name
	//keys = append(keys, args[2]) //Variety

	Avalbytes, err = QueryAsset(stub, AGENT_RATE_PER_RECORD_TYPE, args)
	if err != nil {
		logger.Errorf("getAgentRatePerAsset : Error quering Object from LedgerState %s", err)
		return AgentRatePerAsset{}, false
	}
	agentrate := AgentRatePerAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &agentrate)
	return agentrate, true

}

// ListAgentRatePerAsset  Function will  query  record in ledger based on ID
func ListAgentRatePerAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var itr shim.StateQueryIteratorInterface
	var perlist []AgentRatePerAsset
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("ListAgentRatePerAsset :Attribute approle not found to List AgentRatePerAsset"))
	}
	user := strings.ToUpper(value)
	if user != "MARKETING" {
		return shim.Error(fmt.Sprintf("ListAgentRatePerAsset :This User is not allowed to List AgentRatePerAsset "))
	}

	itr, err = ListAllAsset(stub, AGENT_RATE_PER_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListAgentRatePerAsset : Instence not found in ledger")
		return shim.Error("pricerateitr : Instence not found in ledger")

	}
	defer itr.Close()
	for itr.HasNext() {
		data, derr := itr.Next()
		if derr != nil {
			logger.Errorf("ListAgentRatePerAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListAgentRatePerAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		readobj := AgentRatePerAsset{}
		_ = json.Unmarshal([]byte(databyte), &readobj)
		perlist = append(perlist, readobj)
	}
	Avalbytes, err = json.Marshal(perlist)
	logger.Infof("ListAgentRatePerAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListAgentRatePerAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListAgentRatePerAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))

}

// ListAgentRatePerAssetbyID  Function will  query  all record from DB with speficic storage id
func ListAgentRatePerAssetbyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var itr shim.StateQueryIteratorInterface
	var list []AgentRatePerAsset
	if len(args) < 1 {
		return shim.Error("ListAgentRatePerAssetbyID :Incorrect number of arguments. Expecting StorageId")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"SOURCEID\":\"%s\"}}", AGENT_RATE_PER_RECORD_TYPE, args[0])
	logger.Infof("ListAgentRatePerAssetbyID Query string is %s ", queryString)

	itr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListAgentRatePerAssetbyID : Instence not found in ledger")
		return shim.Error("itr : Instence not found in ledger")

	}
	defer itr.Close()
	for itr.HasNext() {
		data, derr := itr.Next()
		if derr != nil {
			logger.Errorf("ListAgentRatePerAssetbyID : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListAgentRatePerAssetbyID: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		readobj := AgentRatePerAsset{}
		_ = json.Unmarshal([]byte(databyte), &readobj)
		list = append(list, readobj)
	}
	Avalbytes, err = json.Marshal(list)
	logger.Infof("ListAgentRatePerAssetbyID Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListAgentRatePerAssetbyID : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListAgentRatePerAssetbyID: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

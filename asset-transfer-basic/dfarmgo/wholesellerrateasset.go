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
	WHOLESELLER_RATE_PER_RECORD_TYPE string = "WholesellerRatePerAsset"
)

type WholeSellerReq struct {
	State         string  `json:"STATE,omitempty"`
	Country       string  `json:"COUNTRY,omitempty"`
	ProduceName   string  `json:"PRODUCE,omitempty"`
	Variety       string  `json:"VARIETY,omitempty"`
	Role          int     `json:"ROLE,omitempty"`
	SourceID      string  `json:"SOURCEID,omitempty"`
	DestinationID string  `json:"DESTINATIONID,omitempty"`
	PerValue      float64 `json:"PERVALUE,omitempty"`
}

//WholeSellerRatePerAsset asset to provide Markup and discount on final Buyer payment for a produce Variety

type WholeSellerRatePerAsset struct {
	DocType       string  `json:"docType,omitempty"`
	SourceID      string  `json:"SOURCEID,omitempty"`
	DestinationID string  `json:"DESTINATIONID,omitempty"`
	Variety       string  `json:"VARIETY,omitempty"`
	ProduceName   string  `json:"PRODUCE,omitempty"`
	QualityType   string  `json:"QUALITYTYPE,omitempty"`
	PerValue      float64 `json:"PERVALUE,omitempty"`
}

// WholeSellerRatePerAsset to define DC rate
func CreateWholeSellerRatePerAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	if len(args) < 1 {
		logger.Errorf("CreateWholeSellerRatePerAsset : Incorrect number of arguments.")
		return shim.Error("CreateWholeSellerRatePerAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("CreateWholeSellerRatePerAsset :Attribute approle not found to create ParticipantRateAsset"))
	}
	//Abhijeet add role logic here
	if "WHOLESALER" != strings.ToUpper(value) {
		return shim.Error(fmt.Sprintf("CreateWholeSellerRatePerAsset :This User is not allowed to create ParticipantRateAsset"))
	}
	// Convert the arg to a WholeSellerRatePerAsset Object
	logger.Infof("CreateWholeSellerRatePerAsset: Arguments for ledgerapi %s : ", args[0])
	asset := AgentRateReq{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	logger.Infof("CreateWholeSellerRatePerAsset :AgentRateReq is : %s ", asset)
	keys = append(keys, asset.DestinationID)
	keys = append(keys, asset.ProduceName)
	keys = append(keys, asset.Variety)
	Avalbytes, err = QueryAsset(stub, WHOLESELLER_RATE_PER_RECORD_TYPE, keys)
	if err != nil {
		/*	wskey = append(wskey, asset.Country)
			wskey = append(wskey, asset.State)
			wskey = append(wskey, asset.ProduceName)
			dfarmrate, ok := getdFarmRateAsset(stub, wskey)
			if !ok {
				logger.Errorf("CreateWholeSellerRatePerAsset : Error reading dfram max marketing agent rate  %t", ok)
				return shim.Error(fmt.Sprintf("CreateWholeSellerRatePerAsset : Error reading dfram max marketing agent rate  %t", ok))

			}
			if dfarmrate.MaxAgentPer < asset.PerValue {
				logger.Errorf("CreateWholeSellerRatePerAsset : dFarm Max agent rate:  [ %f ] is less then Agent rate  : [ %f ]", dfarmrate.MaxAgentPer, asset.PerValue)
				return shim.Error(fmt.Sprintf("CreateWholeSellerRatePerAsset : dFarm Max agent rate:  [ %f ] is less then Agent rate  : [ %f ]", dfarmrate.MaxAgentPer, asset.PerValue))

			}
		*/
		perasset := WholeSellerRatePerAsset{}
		perasset.DocType = WHOLESELLER_RATE_PER_RECORD_TYPE
		perasset.SourceID = asset.SourceID
		perasset.ProduceName = asset.ProduceName
		perasset.DestinationID = asset.DestinationID
		perasset.Variety = asset.Variety
		perasset.PerValue = asset.PerValue
		Avalbytes, _ = json.Marshal(perasset)
		logger.Infof("CreateWholeSellerRatePerAsset :WholeSellerRatePerAsset  is : %s ", Avalbytes)
		//logger.Infof("CreateWholeSellerRatePerAsset : WholeSellerRatePerAsset  percentage is : %f ", perasset.PerValue)

		err = CreateAsset(stub, WHOLESELLER_RATE_PER_RECORD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreateWholeSellerRatePerAsset : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreateWholeSellerRatePerAsset : ParticipantRateAsset Object first time create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}
	readobj := WholeSellerRatePerAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &readobj)
	readobj.PerValue = asset.PerValue
	Avalbytes, _ = json.Marshal(readobj)
	logger.Infof("CreateWholeSellerRatePerAsset :WholeSellerRatePerAsset Asset is : %v ", readobj)
	err = UpdateAssetWithoutGet(stub, WHOLESELLER_RATE_PER_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateWholeSellerRatePerAsset : Error inserting Object first time  into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateWholeSellerRatePerAsset : WholeSellerRatePerAsset Object first time create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

//getWholeSellerRatePerAsset to get WholeSellerRatePerAsset from ledger based on 3 param
func getWholeSellerRatePerAsset(stub shim.ChaincodeStubInterface, args []string) (WholeSellerRatePerAsset, bool) {

	var err error
	var Avalbytes []byte

	if len(args) == 0 {
		logger.Errorf("getWholeSellerRatePerAsset : Incorrect number of arguments.")
		return WholeSellerRatePerAsset{}, false
	}
	logger.Infof("getWholeSellerRatePerAsset :args is : %v ", args)
	//keys = append(keys, args[0]) //Destination id i.e Buyer ID
	//keys = append(keys, args[1]) //Produce Name
	//keys = append(keys, args[2]) //Variety

	Avalbytes, err = QueryAsset(stub, WHOLESELLER_RATE_PER_RECORD_TYPE, args)
	if err != nil {
		logger.Errorf("getWholeSellerRatePerAsset : Error quering Object from LedgerState %s", err)
		return WholeSellerRatePerAsset{}, false
	}
	wsrate := WholeSellerRatePerAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &wsrate)
	return wsrate, true

}

// ListWholeSellerRatePerAsset  Function will  query  record in ledger based on ID
func ListWholeSellerRatePerAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var itr shim.StateQueryIteratorInterface
	var perlist []WholeSellerRatePerAsset
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("ListWholeSellerRatePerAsset :Attribute approle not found to List WholeSellerRatePerAsset"))
	}
	user := strings.ToUpper(value)
	//Abhijeet hange this logic
	if user != "WHOLESALER" {
		return shim.Error(fmt.Sprintf("ListWholeSellerRatePerAsset :This User is not allowed to List WholeSellerRatePerAsset "))
	}

	itr, err = ListAllAsset(stub, WHOLESELLER_RATE_PER_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListWholeSellerRatePerAsset : Instence not found in ledger")
		return shim.Error("pricerateitr : Instence not found in ledger")

	}
	defer itr.Close()
	for itr.HasNext() {
		data, derr := itr.Next()
		if derr != nil {
			logger.Errorf("ListWholeSellerRatePerAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListWholeSellerRatePerAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		readobj := WholeSellerRatePerAsset{}
		_ = json.Unmarshal([]byte(databyte), &readobj)
		perlist = append(perlist, readobj)
	}
	Avalbytes, err = json.Marshal(perlist)
	logger.Infof("ListWholeSellerRatePerAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListWholeSellerRatePerAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListWholeSellerRatePerAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))

}

// ListWholeSellerRatePerAssetbyID  Function will  query  all record from DB with speficic storage id
func ListWholeSellerRatePerAssetbyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var itr shim.StateQueryIteratorInterface
	var list []WholeSellerRatePerAsset
	if len(args) < 1 {
		return shim.Error("ListWholeSellerRatePerAssetbyID :Incorrect number of arguments. Expecting StorageId")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"SOURCEID\":\"%s\"}}", WHOLESELLER_RATE_PER_RECORD_TYPE, args[0])
	logger.Infof("ListWholeSellerRatePerAssetbyID Query string is %s ", queryString)

	itr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListWholeSellerRatePerAssetbyID : Instence not found in ledger")
		return shim.Error("itr : Instence not found in ledger")

	}
	defer itr.Close()
	for itr.HasNext() {
		data, derr := itr.Next()
		if derr != nil {
			logger.Errorf("ListWholeSellerRatePerAssetbyID : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListWholeSellerRatePerAssetbyID: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		readobj := WholeSellerRatePerAsset{}
		_ = json.Unmarshal([]byte(databyte), &readobj)
		list = append(list, readobj)
	}
	Avalbytes, err = json.Marshal(list)
	logger.Infof("ListWholeSellerRatePerAssetbyID Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListWholeSellerRatePerAssetbyID : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListWholeSellerRatePerAssetbyID: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

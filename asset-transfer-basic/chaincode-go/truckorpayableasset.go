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
	TRUCKOR_ASSET_RECORD_TYPE string = "TruckorPayable"
	TRANSPORTINFO_RECARD_TYPE string = "TransportInfo"
)

type TransportInfoAsset struct {
	DocType       string       `json:"docType,omitempty"`
	SourceID      string       `json:"SORUCEID,omitempty"`      //Primary Key
	DestinationID string       `json:"DESTINATIONID,omitempty"` //Primary Key
	Distance      DistanceUnit `json:"DISTANCEUNIT,omitempty"`
}

type TruckorQuantity struct {
	Verity     string `json:"verity,omitempty"`
	Weight     uint64 `json:"weight,omitempty"`
	WeightUnit string `json:"weightUnit,omitempty"`
}
type TruckerPayableAsset struct {
	GID         string            `json:"gid,omitempty"`
	TruckerID   string            `json:"TruckerId,omitempty"`
	Origin      string            `json:"origin,omitempty"`
	Destination string            `json:"destination,omitempty"`
	Quantities  []TruckorQuantity `json:"quantities,omitempty"`
	Payable     float64           `json:"payable,omitempty"`
	Currency    string            `json:"currency,omitempty"`
	Status      string            `json:"status,omitempty"`
}

//*********************** TransportInfoAsset Asset JSON Method *************************
// //Convert JSON  object to TransportInfoAsset Asset
func JsontoTransportInfoAsset(data []byte) (TransportInfoAsset, error) {
	obj := TransportInfoAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to TransportInfoAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert TransportInfoAsset  object to Json Message

func TransportInfoAssettoJson(obj TransportInfoAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

//JsontoTruckerPayableAsset to convert JSON  to asset object
func JsontoTruckerPayableAsset(data []byte) (TruckerPayableAsset, error) {
	obj := TruckerPayableAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to TruckerPayableAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

// CreateTransportInfoAsset Function will  insert record in ledger after receiving request from Client Application
func CreateTransportInfoAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	if len(args) < 1 {
		logger.Errorf("CreateTransportInfoAsset : Incorrect number of arguments.")
		return shim.Error("CreateTransportInfoAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("CreateTransportInfoAsset :Attribute approle not found to create TransportInfoAsset"))
	}
	user := strings.ToUpper(value)
	if user != "ADMIN" {
		return shim.Error(fmt.Sprintf("CreateTransportInfoAsset :This User is not allowed to create TransportInfoAsset "))
	}
	// Convert the arg to a TransportInfoAsset Object
	logger.Infof("CreateTransportInfoAsset: Arguments for ledgerapi %s : ", args[0])
	asset := TransportInfoAsset{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	logger.Infof("CreateTransportInfoAsset :Source ID is : %s ", asset.SourceID)
	logger.Infof("CreateTransportInfoAsset :Destination ID is : %s ", asset.DestinationID)
	keys = append(keys, asset.SourceID)
	keys = append(keys, asset.DestinationID)
	Avalbytes, err = QueryAsset(stub, TRANSPORTINFO_RECARD_TYPE, keys)
	if err != nil {
		asset.DocType = TRANSPORTINFO_RECARD_TYPE
		Avalbytes, _ = TransportInfoAssettoJson(asset)
		logger.Infof("CreateTransportInfoAsset :Ledger Data is  : %s ", Avalbytes)
		err = CreateAsset(stub, TRANSPORTINFO_RECARD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreateTransportInfoAsset : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreateTransportInfoAsset :  Object first time create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}
	read, _ := JsontoTransportInfoAsset([]byte(Avalbytes))
	read.Distance = asset.Distance
	Avalbytes, _ = TransportInfoAssettoJson(read)
	logger.Infof("CreateTransportInfoAsset :Transport  Asset is : %s ", Avalbytes)
	err = UpdateAssetWithoutGet(stub, TRANSPORTINFO_RECARD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateTransportInfoAsset : Error inserting Object first time  into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateTransportInfoAsset : Transport Object update failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

//QueryTransportInfoAsset to get TransportInfoAsset from ledger based on CASID
func QueryTransportInfoAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryTransportInfoAsset : Incorrect number of arguments.")
		return shim.Error("QueryTransportInfoAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("QueryTransportInfoAsset :Attribute approle not found to query TransportInfoAsset"))
	}
	user := strings.ToUpper(value)
	if user != "ADMIN" {
		return shim.Error(fmt.Sprintf("QueryTransportInfoAsset :This User is not allowed to Query TransportInfoAsset "))
	}
	logger.Infof("QueryTransportInfoAsset :Source ID  is : %s ", args[0])
	logger.Infof("QueryTransportInfoAsset :Destination ID is : %s ", args[1])
	keys = append(keys, args[0])
	keys = append(keys, args[1])

	Avalbytes, err = QueryAsset(stub, TRANSPORTINFO_RECARD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryTransportInfoAsset : Error Querying Object from LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryTransportInfoAsset : TransportInfoAsset object get failed %s", err))
	}

	return shim.Success([]byte(Avalbytes))
}

//getTransportInfoAsset to get TransportInfoAsset from ledger based on CASID
func getTransportInfoAsset(stub shim.ChaincodeStubInterface, sourceID string, destinationID string) (TransportInfoAsset, bool) {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(sourceID) == 0 || len(sourceID) == 0 {
		logger.Errorf("getTransportInfoAsset : Incorrect number of arguments.")
		return TransportInfoAsset{}, false
	}

	logger.Infof("QueryTransportInfoAsset :Source ID  is : %s ", sourceID)
	logger.Infof("QueryTransportInfoAsset :Destination ID is : %s ", destinationID)
	keys = append(keys, sourceID)
	keys = append(keys, destinationID)

	Avalbytes, err = QueryAsset(stub, TRANSPORTINFO_RECARD_TYPE, keys)
	if err != nil {
		logger.Errorf("getTransportInfoAsset : Error Querying Object from LedgerState %s", err)
		return TransportInfoAsset{}, false
	}
	asset, _ := JsontoTransportInfoAsset(Avalbytes)

	return asset, true
}

// ListTransportInfoAsset  Function will  query  record in ledger based on ID
func ListTransportInfoAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var pricerateitr shim.StateQueryIteratorInterface
	var priceratelist []TransportInfoAsset
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("QueryTransportInfoAsset :Attribute approle not found to List TransportInfoAsset"))
	}
	user := strings.ToUpper(value)
	if user != "ADMIN" {
		return shim.Error(fmt.Sprintf("QueryTransportInfoAsset :This User is not allowed to List TransportInfoAsset "))
	}

	pricerateitr, err = ListAllAsset(stub, TRANSPORTINFO_RECARD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListTransportInfoAsset : Instence not found in ledger")
		return shim.Error("pricerateitr : Instence not found in ledger")

	}
	defer pricerateitr.Close()
	for pricerateitr.HasNext() {
		data, derr := pricerateitr.Next()
		if derr != nil {
			logger.Errorf("ListTransportInfoAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListTransportInfoAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		pricerate, _ := JsontoTransportInfoAsset([]byte(databyte))
		priceratelist = append(priceratelist, pricerate)
	}
	Avalbytes, err = json.Marshal(priceratelist)
	logger.Infof("ListTransportInfoAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListTransportInfoAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListTransportInfoAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))

}

//Convert TruckerPayableAsset object to Json Message

func TruckerPayableAssettoJson(obj TruckerPayableAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateTruckerPayableAsset Function will  insert record in ledger after receiving request from Client Application
func CreateTruckerPayableAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateTruckerPayableAsset : Incorrect number of arguments.")
		return shim.Error("CreateTruckerPayableAsset : Incorrect number of arguments.")
	}

	// Convert the arg to a TruckerPayableAsset Object
	logger.Infof("CreateTruckerPayableAsset: Arguments for ledgerapi %s : ", args[0])

	asset, err := JsontoTruckerPayableAsset([]byte(args[0]))

	logger.Infof("CreateTruckerPayableAsset :Produce ID is : %s ", asset.GID)

	keys = append(keys, asset.GID)

	logger.Infof("CreateTruckerPayableAsset : Inserting object with data as  %s", args[0])
	Avalbytes = []byte(args[0])

	err = CreateAsset(stub, TRUCKOR_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateTruckerPayableAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateTruckerPayableAsset : TruckerPayableAsset object create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// QueryTruckerPayableAsset  Function will  query  record in ledger based on ID
func QueryTruckerPayableAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryTruckerPayableAsset : Incorrect number of arguments.")
		return shim.Error("QueryTruckerPayableAsset : Incorrect number of arguments.")
	}
	logger.Infof("QueryTruckerPayableAsset :GID is : %s ", args[0])

	keys = append(keys, args[0])
	Avalbytes, err = QueryAsset(stub, TRUCKOR_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryTruckerPayableAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryTruckerPayableAsset : TruckerPayableAsset object get failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

// UpdateTruckerPayableAsset Function will  update record in ledger after receiving request from Client Application
func UpdateTruckerPayableAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("UpdateTruckerPayableAsset : Incorrect number of arguments.")
		return shim.Error("UpdateTruckerPayableAsset : Incorrect number of arguments.")
	}

	// Convert the arg to a TruckerPayableAsset Object
	logger.Infof("UpdateTruckerPayableAsset: Arguments for ledgerapi %s : ", args[0])

	asset, err := JsontoTruckerPayableAsset([]byte(args[0]))

	logger.Infof("UpdateTruckerPayableAsset :GID is : %s ", asset.GID)

	keys = append(keys, asset.GID)

	logger.Infof("UpdateTruckerPayableAsset : updating object with data as  %s", args[0])
	Avalbytes = []byte(args[0])

	err = UpdateAsset(stub, TRUCKOR_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("UpdateTruckerPayableAsset : Error updating Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("UpdateTruckerPayableAsset : TruckerPayableAsset object update failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// ListTruckerPayableAsset  Function will  query  all record from DB
func ListAllTruckerPayableAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var keys []string
	var Avalbytes []byte
	var truckoritr shim.StateQueryIteratorInterface
	var truckorPayableList []TruckerPayableAsset

	truckoritr, err = ListAllAsset(stub, TRUCKOR_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListAllTruckerPayableAsset : Instence not found in ledger")
		return shim.Error("truckoritr : Instence not found in ledger")

	}
	defer truckoritr.Close()
	for truckoritr.HasNext() {
		data, derr := truckoritr.Next()
		if derr != nil {
			logger.Errorf("ListAllTruckerPayableAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListAllTruckerPayableAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		payable, _ := JsontoTruckerPayableAsset([]byte(databyte))
		truckorPayableList = append(truckorPayableList, payable)
	}
	Avalbytes, err = json.Marshal(truckorPayableList)
	logger.Infof("ListAllTruckerPayableAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListAllTruckerPayableAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListAllTruckerPayableAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

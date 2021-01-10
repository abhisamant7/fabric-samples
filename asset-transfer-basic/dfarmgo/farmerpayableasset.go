package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	FARPAYABLE_ASSET_RECORD_TYPE string = "FarmerPayable"
)
//Payable struct 
type Payable struct {
	Verity   string  `json:"verity,omitempty"`
	Currency string  `json:"currency,omitempty"` //usd ,Rs
	Amount   float64 `json:"amount,omitempty"`
	Status   string  `json:"status,omitempty"` //Paid,Unpaid

}
//FarmerPayableAsset struct 
type FarmerPayableAsset struct {
	FarmerID           string    `json:"farmerID,omitempty"`
	ProduceID          string    `json:"produceID,omitempty"`
	Payables           []Payable `json:"payables,omitempty"`
	ExpectedPayable    float64   `json:"expectedPayables,omitempty"`
	TotalPayable       float64   `json:"TotalFarmerPayable,omitempty"`
	TotalPayableStatus string    `json:"TotalPayableAssetstatus,omitempty"`
}

//JsontoFarmerPayableAsset to convert JSON  to asset object
func JsontoFarmerPayableAsset(data []byte) (FarmerPayableAsset, error) {
	obj := FarmerPayableAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to FarmerPayableAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert FarmerPayableAsset object to Json Message

func FarmerPayableAssetoJson(obj FarmerPayableAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateFarmerPayables Function will  insert record in ledger after receiving request from Client Application
func CreateFarmerPayables(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	// Convert the arg to a ProduceAsset Object
	logger.Infof("CreateFarmerPayables: Arguments for ledgerapi %s : ", args[0])

	asset, _ := JsontoFarmerPayableAsset([]byte(args[0]))

	logger.Infof("CreateFarmerPayables :Produce ID is : %s ", asset.ProduceID)

	keys = append(keys, asset.ProduceID)

	logger.Infof("CreateFarmerPayables : Inserting object with data as  %s", args[0])
	Avalbytes = []byte(args[0])

	err = CreateAsset(stub, PRODUCE_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateFarmerPayables : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateFarmerPayables : ProduceAsset object create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

func QueryFarmerPayables(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryFarmerPayables : Incorrect number of arguments.")
		return shim.Error("QueryFarmerPayables : Incorrect number of arguments.")
	}
	logger.Infof("QueryFarmerPayables :Produce Id  is : %s ", args[0])

	keys = append(keys, args[0])
	Avalbytes, err = QueryAsset(stub, FARPAYABLE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryFarmerPayables : Error Querying Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryFarmerPayables : QueryFarmerPayables object get failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

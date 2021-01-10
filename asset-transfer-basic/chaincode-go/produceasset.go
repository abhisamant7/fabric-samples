/*This file contains all the structure  and method for Produce
*
 */
package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	PRODUCE_ASSET_RECORD_TYPE string = "Produce"
)

type Gap struct {
	GapID          string `json:"GAPID,omitempty"`
	Status         string `json:"STATUS,omitempty"`
	ExpirationDate string `json:"EXPIRATIONDATTE,omitempty"`
}

type ProduceQuantity struct {
	VarietyType            string `json:"name,omitempty"`
	EstimatedTotalQuantity uint64 `json:"qty,omitempty"`
	// FinalQuantity          uint64 `json:"FINALQUANTITY,omitempty"`
	//UsableQuantity         uint64         `json:"USABLEQUANTITY,omitempty"`
	//AvailableQuantity      uint64         `json:"AVAILABLEQUANTITY,omitempty"`
	//DiscardedQuantity      uint64         `json:"DISCARDEDQUANTITY,omitempty"`
	//ProcessingQuantity uint64 `json:"PROCESSINGQUANTITY,omitempty"`
	//TableVerities          []TableVariety `json:"TABLEVARIETY,omitempty"`
	//Payment Related information
	//EstimatedTotalPayment uint64        `json:"EstimatedTotalPayment,omitempty"`
	//FinalTotalPayment float64 `json:"FinalTotalPayment,omitempty"`
	//PaymentInfos          []PaymentInfo `json:"PaymentInfos,omitempty"`
}

// ProduceAsset with key as ProduceID
type ProduceAsset struct {
	DocType           string            `json:"docType,omitempty"`
	ProduceID         string            `json:"PRODUCEID,omitempty"`
	ProduceName       string            `json:"PRODUCE,omitempty"`
	ProduceQuantities []ProduceQuantity `json:"PRODUCEQUANTITES,omitempty"` //key as variety
	FarmLocation      string            `json:"FARMLOCATION,omitempty"`
	PlantingDate      string            `json:"PLANTINGDATE,omitempty"`
	GAPInfo           Gap               `json:"GAPINFO,omitempty"`
	FarmerID          string            `json:"FARMERID,omitempty"`
	Status            string            `json:"STATUS,omitempty"` //Registered,Approved,In-transit,Cleaning done,Inspection done,Delivered to Buyer,Financial settlement done from buyer Side,
	BaseUnit          string            `json:"BASE_UNIT,omitempty"`
	Unit              Unit              `json:"SELECTED_UNIT,omitempty"`
	//financial settlement done from dFarm Side,Tracking closed
	//StatusHistories []StatusHistory `json:"statusHistories,omitempty"`
}

//JsontoProduceAsset to convert JSON  to asset object
func JsontoproduceAsset(data []byte) (ProduceAsset, error) {
	obj := ProduceAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to ProduceAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert ProduceAsset object to Json Message

func ProduceAssettoJson(obj ProduceAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateProduceAsset Function will  insert record in ledger after receiving request from Client Application
func CreateProduceAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateProduceAsset : Incorrect number of arguments.")
		return shim.Error("CreateProduceAsset : Incorrect number of arguments.")
	}

	// Convert the arg to a ProduceAsset Object
	logger.Infof("CreateProduceAsset: Arguments for ledgerapi %s : ", args[0])

	asset, err := JsontoproduceAsset([]byte(args[0]))
	asset.DocType = PRODUCE_ASSET_RECORD_TYPE

	logger.Infof("CreateProduceAsset :Produce ID is : %s ", asset.ProduceID)

	keys = append(keys, asset.ProduceID)

	logger.Infof("CreateProduceAsset : Inserting object with data as  %s", args[0])
	Avalbytes, _ = ProduceAssettoJson(asset)

	err = CreateAsset(stub, PRODUCE_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateProduceAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateProduceAsset : ProduceAsset object create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// QueryProduceAsset  Function will  query  record in ledger based on ID
func QueryProduceAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryProduceAsset : Incorrect number of arguments.")
		return shim.Error("QueryProduceAsset : Incorrect number of arguments.")
	}
	logger.Infof("QueryProduceAsset :Farmer ID is : %s ", args[0])

	keys = append(keys, args[0])
	Avalbytes, err = QueryAsset(stub, PRODUCE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryProduceAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryProduceAsset : ProduceAsset object get failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

// QueryProduceAssetInfo  Function will  query  record in ledger based on ID
func QueryProduceAssetInfo(stub shim.ChaincodeStubInterface, args []string) (ProduceAsset, bool) {
	var err error
	var Avalbytes []byte
	var keys []string
	var produce ProduceAsset

	if len(args) < 1 {
		logger.Errorf("QueryProduceAssetInfo : Incorrect number of arguments.")
		return produce, false
	}
	logger.Infof("QueryProduceAsset :Farmer ID is : %s ", args[0])

	keys = append(keys, args[0])
	Avalbytes, err = QueryAsset(stub, PRODUCE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryProduceAsset : Error inserting Object into LedgerState %s", err)
		return produce, false
	}
	produce, err = JsontoproduceAsset(Avalbytes)
	logger.Infof("produce Data %v", produce)
	return produce, true

}

// UpdateProduceAsset Function will  update record in ledger after receiving request from Client Application
func UpdateProduceAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("UpdateProduceAsset : Incorrect number of arguments.")
		return shim.Error("UpdateProduceAsset : Incorrect number of arguments.")
	}

	// Convert the arg to a ProduceAsset Object
	logger.Infof("UpdateProduceAsset: Arguments for ledgerapi %s : ", args[0])

	asset, err := JsontoproduceAsset([]byte(args[0]))

	logger.Infof("UpdateProduceAsset :Produce ID is : %s ", asset.ProduceID)

	keys = append(keys, asset.ProduceID)

	logger.Infof("UpdateProduceAsset : updating object with data as  %s", args[0])
	Avalbytes, err = QueryAsset(stub, PRODUCE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("UpdateProduceAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("UpdateProduceAsset : ProduceAsset object get failed %s", err))
	}
	assetread, _ := JsontoproduceAsset([]byte(Avalbytes))
	if len(asset.ProduceQuantities) != 0 {
		assetread.ProduceQuantities = asset.ProduceQuantities
	}
	if len(asset.ProduceName) != 0 {
		assetread.ProduceName = asset.ProduceName
	}
	if len(asset.FarmLocation) != 0 {
		assetread.FarmLocation = asset.FarmLocation
	}
	if len(asset.GAPInfo.GapID) != 0 {
		assetread.GAPInfo = asset.GAPInfo
	}
	if len(asset.PlantingDate) != 0 {
		assetread.PlantingDate = asset.PlantingDate
	}
	if len(asset.Status) != 0 {
		assetread.Status = asset.Status
	}
	logger.Infof("UpdateProduceAsset : updating object with data as  %v", assetread)
	Avalbytes, _ = ProduceAssettoJson(assetread)

	err = UpdateAssetWithoutGet(stub, PRODUCE_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("UpdateProduceAsset : Error updating Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("UpdateProduceAsset : ProduceAsset object update failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

//ListProduceAsset  Function will  query  all record from DB

func ListAllProduceAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var keys []string
	var Avalbytes []byte
	var produceitr shim.StateQueryIteratorInterface
	var produceList []ProduceAsset

	produceitr, err = ListAllAsset(stub, PRODUCE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListAllProduceAsset : Instence not found in ledger")
		return shim.Error("produceitr : Instence not found in ledger")

	}
	defer produceitr.Close()
	for produceitr.HasNext() {
		data, derr := produceitr.Next()
		if derr != nil {
			logger.Errorf("ListAllProduceAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListAllProduceAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		produce, _ := JsontoproduceAsset([]byte(databyte))
		produceList = append(produceList, produce)
	}
	Avalbytes, err = json.Marshal(produceList)
	logger.Infof("ListAllProduceAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListAllProduceAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListAllProduceAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

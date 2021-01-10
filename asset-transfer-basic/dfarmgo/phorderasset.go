package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	PHORDER_ASSET_RECORD_TYPE string = "PHOrder"
)

//JsontoPHOrderAsset to convert JSON  to asset object
func JsontoPHOrderAsset(data []byte) (PHOrderAsset, error) {
	obj := PHOrderAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to PHOrderAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//PHOrderAssettoJson Convert PHOrderAsset object to Json Message
func PHOrderAssettoJson(obj PHOrderAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

//CreatePHOrderAsset This function will be used by Farmer to book order for pickup by PH Center.
func CreatePHOrderAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreatePHOrderAsset : Incorrect number of arguments.")
		return shim.Error("CreatePHOrderAsset : Incorrect number of arguments.")
	}
	// Convert the arg to a PHOrderAsset Object
	logger.Infof("CreatePHOrderAsset: Arguments for ledgerapi %s : ", args[0])
	asset, err := JsontoPHOrderAsset([]byte(args[0]))
	logger.Infof("CreatePHOrderAsset :OrderId ID is : %s ", asset.OrderID)
	asset.DocType = PHORDER_ASSET_RECORD_TYPE
	timeinfo, err := stub.GetTxTimestamp()
	if err != nil {
		logger.Errorf("CreateCcOrderAsset : Error getting  timestamp  %s", err)
		return shim.Error(fmt.Sprintf("CreateCSOrderAsset : PHOrderAsset object create failed due to timestamp read %s", err))
	}
	logger.Infof("CreatePHOrderAsset: Time stamp is %+v ", timeinfo)
	asset.OrderUnixTime = timeinfo.GetSeconds()
	keys = append(keys, asset.OrderID)
	Avalbytes, _ = PHOrderAssettoJson(asset)
	err = CreateAsset(stub, PHORDER_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreatePHOrderAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreatePHOrderAsset : PHOrderAsset object create failed %s", err))
	}
	logger.Infof("CreatePHOrderAsset Asset : %s ", string(Avalbytes))
	return shim.Success([]byte(Avalbytes))
}
func QueryPHOrderAsset(stub shim.ChaincodeStubInterface, args []string) *PHOrderAsset {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryPHOrderAsset : Incorrect number of arguments.")
		return nil
	}
	// Convert the arg to a PHOrderAsset Object
	//logger.Infof("CreatePHOrderAsset: Arguments for ledgerapi %s : ", args[0])
	logger.Infof("QueryPHOrderAsset :OrderId ID is : %s ", args[0])
	keys = append(keys, args[0])

	Avalbytes, err = QueryAsset(stub, PHORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryPHOrderAsset : Error inserting Object into LedgerState %s", err)
		return nil
	}
	PHOrder := PHOrderAsset{}
	PHOrder, _ = JsontoPHOrderAsset(Avalbytes)
	logger.Infof("QueryPHOrderAsset :PHOrder  is : %v ", PHOrder)
	return &PHOrder
}

// ListPHOrderbyStorageId  Function will  query  all record from ledger with specific storage id
func ListPHOrderbyStorageId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var orderitr shim.StateQueryIteratorInterface
	var orderList []PHOrderAsset
	if len(args) < 1 {
		return shim.Error("ListPHOrderbyStorageId :Incorrect number of arguments. Expecting StorageId")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"DESTINATIONID\":\"%s\"}}", PHORDER_ASSET_RECORD_TYPE, args[0])
	logger.Infof("ListPHOrderbyStorageId Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListPHOrderbyStorageId : Instance not found in ledger")
		return shim.Error("orderitr : Instance not found in ledger")

	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ListPHOrderbyStorageId : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListPHOrderbyStorageId: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		order, _ := JsontoPHOrderAsset([]byte(databyte))
		orderList = append(orderList, order)
	}
	Avalbytes, err = json.Marshal(orderList)
	logger.Infof("ListPHOrderbyStorageId Responce for App : %s", Avalbytes)
	if err != nil {
		logger.Errorf("ListPHOrderbyStorageId : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListPHOrderbyStorageId: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

func QueryPHOrderAssetbyOrderID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryPHOrderAssetbyOrderID : Incorrect number of arguments.")
		return shim.Error("QueryPHOrderAssetbyOrderID : Incorrect number of arguments.")
	}
	logger.Infof("QueryPHOrderAssetbyOrderID :Order ID is : %s ", args[0])

	keys = append(keys, args[0])
	Avalbytes, err = QueryAsset(stub, PHORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryPHOrderAssetbyOrderID : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryPHOrderAssetbyOrderID : PHOrderAsset object get failed %s", err))
	}
	logger.Infof("QueryPHOrderAssetbyOrderID :asset is : %s ", string(Avalbytes))
	return shim.Success([]byte(Avalbytes))

}

//GetPHOrderList Function to get list of PH order based on produceId and Variety
func GetPHOrderList(stub shim.ChaincodeStubInterface, args []string) ([]string, bool) {
	var err error

	var orderitr shim.StateQueryIteratorInterface

	var orderlist []string
	if len(args) < 2 {
		logger.Infof("No of args are not 2  ")
		return orderlist, false
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PRID\":\"%s\",\"variety\":\"%s\"},\"fields\":[\"ORDERID\"]}", PHORDER_ASSET_RECORD_TYPE, args[0], args[1])
	logger.Infof("GetPHOrderList Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("GetPHOrderList : Instance not found in ledger")
		return orderlist, false
	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("GetPHOrderList : Cannot parse result set. Error : %v", derr)
			return orderlist, false
		}
		databyte := data.GetValue()
		var resultdata AssetData
		json.Unmarshal(databyte, &resultdata)
		logger.Infof("result is %v   ", resultdata)
		orderlist = append(orderlist, resultdata.OrderID)

	}
	logger.Infof("GetPHOrderList Query result  is %v ", orderlist)

	return orderlist, true

}

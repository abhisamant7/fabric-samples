package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	CCORDER_ASSET_RECORD_TYPE string = "CCOrder"
)

//JsontoCCOrderAsset to convert JSON  to asset object
func JsontoCCOrderAsset(data []byte) (CCOrderAsset, error) {
	obj := CCOrderAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to CCOrderAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert CCOrderAsset object to Json Message

func CCOrderAssettoJson(obj CCOrderAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// This function will be used by Farmer to book order for pickup by CCollection Center.

func CreateCCOrderAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateCCOrderAsset : Incorrect number of arguments.")
		return shim.Error("CreateCCOrderAsset : Incorrect number of arguments.")
	}
	// Convert the arg to a CCOrderAsset Object
	logger.Infof("CreateCCOrderAsset: Arguments for ledgerapi %s : ", args[0])
	asset, err := JsontoCCOrderAsset([]byte(args[0]))
	logger.Infof("CreateCCOrderAsset :OrderId ID is : %s ", asset.OrderID)
	asset.DocType = CCORDER_ASSET_RECORD_TYPE
	timeinfo, err := stub.GetTxTimestamp()
	if err != nil {
		logger.Errorf("CreateCcOrderAsset : Error getting  timestamp  %s", err)
		return shim.Error(fmt.Sprintf("CreateCSOrderAsset : CCOrderAsset object create failed due to timestamp read %s", err))
	}
	logger.Infof("CreateCCOrderAsset: Time stamp is %+v ", timeinfo)
	asset.OrderUnixTime = timeinfo.GetSeconds()
	keys = append(keys, asset.OrderID)
	Avalbytes, _ = CCOrderAssettoJson(asset)
	err = CreateAsset(stub, CCORDER_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateCCOrderAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateCCOrderAsset : CCOrderAsset object create failed %s", err))
	}
	logger.Infof("CreateCCOrderAsset Asset : %s ", string(Avalbytes))
	return shim.Success([]byte(Avalbytes))
}
func QueryCCOrderAsset(stub shim.ChaincodeStubInterface, args []string) *CCOrderAsset {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryCCOrderAsset : Incorrect number of arguments.")
		return nil
	}
	// Convert the arg to a CCOrderAsset Object
	//logger.Infof("CreateCCOrderAsset: Arguments for ledgerapi %s : ", args[0])
	logger.Infof("QueryCCOrderAsset :OrderId ID is : %s ", args[0])
	keys = append(keys, args[0])

	Avalbytes, err = QueryAsset(stub, CCORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryCCOrderAsset : Error inserting Object into LedgerState %s", err)
		return nil
	}
	CCOrder := CCOrderAsset{}
	CCOrder, _ = JsontoCCOrderAsset(Avalbytes)
	logger.Infof("QueryCCOrderAsset :CCOrder  is : %v ", CCOrder)
	return &CCOrder
}

// ListCCOrderbyStorageId  Function will  query  all record from ledger with specific storage id
func ListCCOrderbyStorageId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var orderitr shim.StateQueryIteratorInterface
	var orderList []CCOrderAsset
	if len(args) < 1 {
		return shim.Error("ListCCOrderbyStorageId :Incorrect number of arguments. Expecting StorageId")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"DESTINATIONID\":\"%s\"}}", CCORDER_ASSET_RECORD_TYPE, args[0])
	logger.Infof("ListCCOrderbyStorageId Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListCCOrderbyStorageId : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ListCCOrderbyStorageId : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListCCOrderbyStorageId: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		order, _ := JsontoCCOrderAsset([]byte(databyte))
		orderList = append(orderList, order)
	}
	Avalbytes, err = json.Marshal(orderList)
	logger.Infof("ListCCOrderbyStorageId Responce for App : %s", Avalbytes)
	if err != nil {
		logger.Errorf("ListCCOrderbyStorageId : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListCCOrderbyStorageId: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

func QueryCCOrderAssetbyOrderID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryCCOrderAssetbyOrderID : Incorrect number of arguments.")
		return shim.Error("QueryCCOrderAssetbyOrderID : Incorrect number of arguments.")
	}
	logger.Infof("QueryCCOrderAssetbyOrderID :Order ID is : %s ", args[0])

	keys = append(keys, args[0])
	Avalbytes, err = QueryAsset(stub, CCORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryCCOrderAssetbyOrderID : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryCCOrderAssetbyOrderID : CCOrderAsset object get failed %s", err))
	}
	logger.Infof("QueryCCOrderAssetbyOrderID :asset is : %s ", string(Avalbytes))
	return shim.Success([]byte(Avalbytes))

}

//GetCCOrderList Function to get list of CC order based on produceId and Variety

func GetCCOrderList(stub shim.ChaincodeStubInterface, args []string) ([]string, bool) {
	var err error

	var orderitr shim.StateQueryIteratorInterface

	var orderlist []string
	if len(args) < 2 {
		logger.Infof("No of args are not 2  ")
		return orderlist, false
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PRID\":\"%s\",\"variety\":\"%s\"},\"fields\":[\"ORDERID\"]}", CCORDER_ASSET_RECORD_TYPE, args[0], args[1])
	logger.Infof("GetCCOrderList Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("GetCCOrderList : Instence not found in ledger")
		return orderlist, false
	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("GetCCOrderList : Cannot parse result set. Error : %v", derr)
			return orderlist, false
		}
		databyte := data.GetValue()
		var resultdata AssetData
		json.Unmarshal(databyte, &resultdata)
		logger.Infof("result is %v   ", resultdata)
		orderlist = append(orderlist, resultdata.OrderID)

	}
	logger.Infof("GetCCOrderList Query result  is %v ", orderlist)

	return orderlist, true

}

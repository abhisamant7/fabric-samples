package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	CSORDER_ASSET_RECORD_TYPE string = "CSOrder"
)

//JsontoCSOrderAsset to convert JSON  to asset object
func JsontoCSOrderAsset(data []byte) (CSOrderAsset, error) {
	obj := CSOrderAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to CSOrderAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert CSOrderAsset object to Json Message

func CSOrderAssettoJson(obj CSOrderAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// This function will be used by Farmer to book order for pickup by CAS.

func CreateCSOrderAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateCSOrderAsset : Incorrect number of arguments.")
		return shim.Error("CreateCSOrderAsset : Incorrect number of arguments.")
	}
	// Convert the arg to a CSOrderAsset Object
	logger.Infof("CreateCSOrderAsset: Arguments for ledgerapi %s : ", args[0])
	asset, err := JsontoCSOrderAsset([]byte(args[0]))
	logger.Infof("CreateCSOrderAsset :OrderId ID is : %s ", asset.OrderID)
	asset.DocType = CSORDER_ASSET_RECORD_TYPE
	timeinfo, err := stub.GetTxTimestamp()
	if err != nil {
		logger.Errorf("CreateCSOrderAsset : Error getting  timestamp  %s", err)
		return shim.Error(fmt.Sprintf("CreateCSOrderAsset : CSOrderAsset object create failed due to timestamp read %s", err))
	}
	logger.Infof("CreateCSOrderAsset: Time stamp is %+v ", timeinfo)
	asset.OrderUnixTime = timeinfo.GetSeconds()
	keys = append(keys, asset.OrderID)
	Avalbytes, _ = CSOrderAssettoJson(asset)
	err = CreateAsset(stub, CSORDER_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateCSOrderAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateCSOrderAsset : CSOrderAsset object create failed %s", err))
	}
	logger.Infof("UpdateCSOrderStatus Asset : %s ", string(Avalbytes))
	return shim.Success([]byte(Avalbytes))
}
func QueryCSOrderAsset(stub shim.ChaincodeStubInterface, args []string) *CSOrderAsset {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateCSOrderAsset : Incorrect number of arguments.")
		return nil
	}
	// Convert the arg to a CSOrderAsset Object
	//logger.Infof("CreateCSOrderAsset: Arguments for ledgerapi %s : ", args[0])
	logger.Infof("UpdateCSOrderStatus :OrderId ID is : %s ", args[0])
	keys = append(keys, args[0])

	Avalbytes, err = QueryAsset(stub, CSORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("CreateCSOrderAsset : Error inserting Object into LedgerState %s", err)
		return nil
	}
	csOrder := CSOrderAsset{}
	csOrder, _ = JsontoCSOrderAsset(Avalbytes)
	logger.Infof("UpdateCSOrderStatus :csOrder  is : %v ", csOrder)
	return &csOrder
}

// ListCSOrderbyStorageId  Function will  query  all record from DB with speficic storage id
func ListCSOrderbyStorageId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var orderitr shim.StateQueryIteratorInterface
	var orderList []CSOrderAsset
	if len(args) < 1 {
		return shim.Error("ListCSOrderbyStorageId :Incorrect number of arguments. Expecting StorageId")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"DESTINATIONID\":\"%s\"}}", CSORDER_ASSET_RECORD_TYPE, args[0])
	logger.Infof("ListCSOrderbyStorageId Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListCSOrderbyStorageId : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ListCSOrderbyStorageId : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListCSOrderbyStorageId: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		order, _ := JsontoCSOrderAsset([]byte(databyte))
		orderList = append(orderList, order)
	}
	Avalbytes, err = json.Marshal(orderList)
	logger.Infof("ListCSOrderbyStorageId Responce for App : %s", Avalbytes)
	if err != nil {
		logger.Errorf("ListCSOrderbyStorageId : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListCSOrderbyStorageId: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

func QueryCSOrderAssetbyOrderID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryCSOrderAssetbyOrderID : Incorrect number of arguments.")
		return shim.Error("QueryCSOrderAssetbyOrderID : Incorrect number of arguments.")
	}
	logger.Infof("QueryCSOrderAssetbyOrderID :Order ID is : %s ", args[0])

	keys = append(keys, args[0])
	Avalbytes, err = QueryAsset(stub, CSORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryCSOrderAssetbyOrderID : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryCSOrderAssetbyOrderID : CSOrderAsset object get failed %s", err))
	}
	logger.Infof("QueryCSOrderAssetbyOrderID :asset is : %s ", string(Avalbytes))
	return shim.Success([]byte(Avalbytes))

}

//GetCASOrderList Function to get list of CAS order based on produceId and Variety

func GetCASOrderList(stub shim.ChaincodeStubInterface, args []string) ([]string, bool) {
	var err error

	var orderitr shim.StateQueryIteratorInterface

	var orderlist []string
	if len(args) < 2 {
		logger.Infof("No of args are not 2  ")
		return orderlist, false
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PRID\":\"%s\",\"variety\":\"%s\"},\"fields\":[\"ORDERID\"]}", CSORDER_ASSET_RECORD_TYPE, args[0], args[1])
	logger.Infof("GetCASOrderList Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("GetCASOrderList : Instence not found in ledger")
		return orderlist, false
	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("GetCASOrderList : Cannot parse result set. Error : %v", derr)
			return orderlist, false
		}
		databyte := data.GetValue()
		var resultdata AssetData
		json.Unmarshal(databyte, &resultdata)
		logger.Infof("result is %v   ", resultdata)
		orderlist = append(orderlist, resultdata.OrderID)

	}
	logger.Infof("GetCASOrderList Query result  is %v ", orderlist)

	return orderlist, true

}

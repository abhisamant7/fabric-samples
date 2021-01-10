package main

import (
	"encoding/json"
	"fmt"

	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	PC_ORDER_ASSET_RECORD_TYPE string = "PCOrderAsset"
)

//JsontoPCOrderAsset to convert JSON  to asset object
func JsontoPCOrderAsset(data []byte) (PCOrderAsset, error) {
	obj := PCOrderAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to CAStoPCOrderAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert CAStoPCOrderAsset object to Json Message

func PCOrderAssettoJson(obj PCOrderAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// This function will be used by Farmer to book order for pickup by CAS.

func CreatePCOrderAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreatePCOrderAsset : Incorrect number of arguments.")
		return shim.Error("CreatePCOrderAsset : Incorrect number of arguments.")
	}
	// Convert the arg to a PCOrderAsset Object
	logger.Infof("CreatePCOrderAsset: Arguments for ledgerapi %s : ", args[0])
	asset, err := JsontoPCOrderAsset([]byte(args[0]))
	asset.DocType = PC_ORDER_ASSET_RECORD_TYPE
	timeinfo, err := stub.GetTxTimestamp()
	if err != nil {
		logger.Errorf("CreatePCOrderAsset : Error getting  timestamp  %s", err)
		return shim.Error(fmt.Sprintf("CreatePCOrderAsset :PCOrderAsset object create failed due to timestamp read %s", err))
	}
	logger.Infof("CreatePCOrderAsset: Time stamp is %+v ", timeinfo)

	logger.Infof("CreatePCOrderAsset :OrderId ID is : %s ", asset.OrderID)
	logger.Infof("CreatePCOrderAsset :Asset is : %v ", asset)
	keys = append(keys, asset.OrderID)
	Avalbytes, _ = PCOrderAssettoJson(asset)
	err = CreateAsset(stub, PC_ORDER_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreatePCOrderAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreatePCOrderAsset : PCOrderAsset object create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}
func QueryPCOrderAsset(stub shim.ChaincodeStubInterface, args []string) *PCOrderAsset {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreatePCOrderAsset : Incorrect number of arguments.")
		return nil
	}
	// Convert the arg to a PCOrderAsset Object
	//logger.Infof("CreatePCOrderAsset: Arguments for ledgerapi %s : ", args[0])
	logger.Infof("UpdateCASPCOrderStatus :OrderId ID is : %s ", args[0])
	keys = append(keys, args[0])

	Avalbytes, err = QueryAsset(stub, PC_ORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("CreatePCOrderAsset : Error inserting Object into LedgerState %s", err)
		return nil
	}
	csOrder := PCOrderAsset{}
	csOrder, _ = JsontoPCOrderAsset(Avalbytes)
	logger.Infof("UpdateCASPCOrderStatus :Asset is : %v ", csOrder)
	return &csOrder
}

// ListFarmerCSOrder  Function will  query  all record from DB with speficic storage id
func ListPCOrderAssetbyStorageId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var orderitr shim.StateQueryIteratorInterface
	var orderList []PCOrderAsset
	if len(args) < 1 {
		return shim.Error("ListFarmerCSOrderbyStorageId :Incorrect number of arguments. Expecting StorageId")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"SOURCEID\":\"%s\"}}", PC_ORDER_ASSET_RECORD_TYPE, args[0])
	logger.Infof("ListFarmerCSOrderbyStorageId Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListFarmerCSOrderbyStorageId : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ListFarmerCSOrderbyStorageId : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListFarmerCSOrderbyStorageId: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		order, _ := JsontoPCOrderAsset([]byte(databyte))
		orderList = append(orderList, order)
	}
	Avalbytes, err = json.Marshal(orderList)
	logger.Infof("ListFarmerCSOrderbyStorageId Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListFarmerCSOrderbyStorageId : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListFarmerCSOrderbyStorageId: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

func ListPCOrderAssetbyPcID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var orderitr shim.StateQueryIteratorInterface
	var orderList []PCOrderAsset
	if len(args) < 1 {
		return shim.Error("ListPCOrderAssetbyPcID :Incorrect number of arguments. Expecting StorageId")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"DESTINATIONID\":\"%s\"}}", PC_ORDER_ASSET_RECORD_TYPE, args[0])
	logger.Infof("ListPCOrderAssetbyPcID Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListPCOrderAssetbyPcID : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ListPCOrderAssetbyPcID : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListPCOrderAssetbyPcID: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		order, _ := JsontoPCOrderAsset([]byte(databyte))
		orderList = append(orderList, order)
	}
	Avalbytes, err = json.Marshal(orderList)
	logger.Infof("ListPCOrderAssetbyPcID Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListPCOrderAssetbyPcID : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListPCOrderAssetbyPcID: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// ListPCOrderAsset  Function will  query  all record from DB with speficic storage id
func ListPCOrderAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var orderitr shim.StateQueryIteratorInterface
	var orderList []PCOrderAsset
	var keys []string
	if len(args) < 1 {
		return shim.Error("ListPCOrderAsset :Incorrect number of arguments. Expecting StorageId")
	}
	orderitr, err = ListAllAsset(stub, PC_ORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListPCOrderAsset : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ListPCOrderAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListPCOrderAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		order, _ := JsontoPCOrderAsset([]byte(databyte))
		orderList = append(orderList, order)
	}
	Avalbytes, err = json.Marshal(orderList)
	logger.Infof("ListPCOrderAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListPCOrderAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListPCOrderAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

type CombinedRate struct {
	Rate float64 `json:"combinedRatemitempty"`
}

func QueryPCOrderAssetbyOrderID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryPCOrderAssetbyOrderID : Incorrect number of arguments.")
		return shim.Error("QueryPCOrderAssetbyOrderID : Incorrect number of arguments.")
	}
	logger.Infof("QueryPCOrderAssetbyOrderID :Order ID is : %s ", args[0])

	keys = append(keys, args[0])
	Avalbytes, err = QueryAsset(stub, PC_ORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryPCOrderAssetbyOrderID : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryPCOrderAssetbyOrderID : PCOrderAsset object get failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

//GetSourceBasedPCOrderList is used to get PC orders created from a CAS order

func GetSourceBasedPCOrderList(stub shim.ChaincodeStubInterface, args []string) ([]string, bool) {
	var err error
	var orderitr shim.StateQueryIteratorInterface
	var orderlist []string

	if len(args) < 1 {
		logger.Infof("No of args are not 1  ")
		return orderlist, false
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"parentorderId\":\"%s\"},\"fields\":[\"ORDERID\"]}", PC_ORDER_ASSET_RECORD_TYPE, args[0])
	logger.Infof("GetSourceBasedPCOrderList Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("GetSourceBasedPCOrderList : Instence not found in ledger")
		return orderlist, false
	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("GetSourceBasedPCOrderList : Cannot parse result set. Error : %v", derr)
			return orderlist, false
		}
		var resultdata AssetData
		databyte := data.GetValue()
		json.Unmarshal(databyte, &resultdata)
		logger.Infof("result is %v   ", resultdata)
		orderlist = append(orderlist, resultdata.OrderID)
	}
	logger.Infof("GetSourceBasedPCOrderList Query result  is %v ", orderlist)

	return orderlist, true

}

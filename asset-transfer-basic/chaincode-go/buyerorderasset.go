package main

import (
	"encoding/json"
	"fmt"

	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	BUYER_ORDER_ASSET_RECORD_TYPE string = "BuyerOrderAsset"
)

//JsontoBuyerOrderAsset to convert JSON  to asset object
func JsontoBuyerOrderAsset(data []byte) (BuyerOrderAsset, error) {
	obj := BuyerOrderAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to BuyerOrderAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert CAStoPCOrderAsset object to Json Message

func BuyerOrderAssettoJson(obj BuyerOrderAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// This function will be used by Farmer to book order for pickup by CAS.

func CreateBuyerOrderAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateBuyerOrderAsset : Incorrect number of arguments.")
		return shim.Error("CreateBuyerOrderAsset : Incorrect number of arguments.")
	}
	// Convert the arg to a BuyerOrderAsset Object
	logger.Infof("CreateBuyerOrderAsset: Arguments for ledgerapi %s : ", args[0])
	asset, err := JsontoBuyerOrderAsset([]byte(args[0]))
	logger.Infof("UpdateBuyerOrderStatus :OrderId ID is : %s ", asset.OrderID)
	keys = append(keys, asset.OrderID)
	asset.DocType = BUYER_ORDER_ASSET_RECORD_TYPE
	timeinfo, err := stub.GetTxTimestamp()
	if err != nil {
		logger.Errorf("CreateBuyerOrderAsset : Error getting  timestamp  %s", err)
		return shim.Error(fmt.Sprintf("CreateBuyerOrderAsset : BuyerOrderAsset object create failed due to timestamp read %s", err))
	}
	logger.Infof("CreateBuyerOrderAsset: Time stamp is %+v ", timeinfo)
	asset.OrderUnixTime = timeinfo.GetSeconds()
	logger.Infof("UpdateBuyerOrderStatus :asset is : %v ", asset)
	Avalbytes, _ = BuyerOrderAssettoJson(asset)
	err = CreateAsset(stub, BUYER_ORDER_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateBuyerOrderAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateBuyerOrderAsset : BuyerOrderAsset object create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}
func QueryBuyerOrderAsset(stub shim.ChaincodeStubInterface, args []string) *BuyerOrderAsset {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryBuyerOrderAsset : Incorrect number of arguments.")
		return nil
	}
	// Convert the arg to a BuyerOrderAsset Object
	//logger.Infof("CreateBuyerOrderAsset: Arguments for ledgerapi %s : ", args[0])
	logger.Infof("QueryBuyerOrderAsset :OrderId ID is : %s ", args[0])
	keys = append(keys, args[0])

	Avalbytes, err = QueryAsset(stub, BUYER_ORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryBuyerOrderAsset : Error inserting Object into LedgerState %s", err)
		return nil
	}
	csOrder := BuyerOrderAsset{}
	csOrder, _ = JsontoBuyerOrderAsset(Avalbytes)
	return &csOrder
}

// ListBuyerOrderAssetbyBuyerID  Function will  query  all record from DB with speficic storage id
func ListBuyerOrderAssetbyBuyerID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var orderitr shim.StateQueryIteratorInterface
	var orderList []BuyerOrderAsset
	if len(args) < 1 {
		return shim.Error("ListBuyerOrderAssetbyBuyerID :Incorrect number of arguments. Expecting StorageId")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"buyerID\":\"%s\"}}", BUYER_ORDER_ASSET_RECORD_TYPE, args[0])
	logger.Infof("ListBuyerOrderAssetbyBuyerID Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListBuyerOrderAssetbyBuyerID : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ListBuyerOrderAssetbyBuyerID : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListBuyerOrderAssetbyBuyerID: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		order, _ := JsontoBuyerOrderAsset([]byte(databyte))
		orderList = append(orderList, order)
	}
	Avalbytes, err = json.Marshal(orderList)
	logger.Infof("ListBuyerOrderAssetbyBuyerID Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListBuyerOrderAssetbyBuyerID : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListBuyerOrderAssetbyBuyerID: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// ListFarmerCSOrder  Function will  query  all record from DB with speficic storage id
func ListBuyerOrderAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var orderitr shim.StateQueryIteratorInterface
	var orderList []BuyerOrderAsset
	var keys []string
	if len(args) < 1 {
		return shim.Error("ListBuyerOrderAsset :Incorrect number of arguments. Expecting StorageId")
	}
	orderitr, err = ListAllAsset(stub, BUYER_ORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListBuyerOrderAsset : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ListBuyerOrderAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListBuyerOrderAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		order, _ := JsontoBuyerOrderAsset([]byte(databyte))
		orderList = append(orderList, order)
	}
	Avalbytes, err = json.Marshal(orderList)
	logger.Infof("ListBuyerOrderAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListBuyerOrderAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListBuyerOrderAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

func QueryBuyerOrderAssetbyOrderID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryBuyerOrderAssetbyOrderID : Incorrect number of arguments.")
		return shim.Error(fmt.Sprintf("QueryBuyerOrderAssetbyOrderID :Incorrect number of arguments"))
	}
	// Convert the arg to a BuyerOrderAsset Object
	//logger.Infof("CreateBuyerOrderAsset: Arguments for ledgerapi %s : ", args[0])
	logger.Infof("QueryBuyerOrderAssetbyOrderID :OrderId ID is : %s ", args[0])
	keys = append(keys, args[0])

	Avalbytes, err = QueryAsset(stub, BUYER_ORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryBuyerOrderAssetbyOrderID : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryBuyerOrderAssetbyOrderID : Error querying BuyerOrder Order Object into LedgerState %s", err))
	}

	return shim.Success(Avalbytes)
}

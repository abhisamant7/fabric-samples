package main

import (
	"encoding/json"
	"fmt"

	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	DC_ORDER_ASSET_RECORD_TYPE string = "DCOrderAsset"
)

//JsontoDCOrderAsset to convert JSON  to asset object
func JsontoDCOrderAsset(data []byte) (DCOrderAsset, error) {
	obj := DCOrderAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to DCOrderAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert CAStoPCOrderAsset object to Json Message

func DCOrderAssettoJson(obj DCOrderAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// This function will be used by Farmer to book order for pickup by CAS.

func CreateDCOrderAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var AvalbytesRes []byte
	order := DCOrder{}
	if len(args) < 1 {
		logger.Errorf("CreateDCOrderAsset : Incorrect number of arguments.")
		return shim.Error("CreateDCOrderAsset : Incorrect number of arguments.")
	}
	// Convert the arg to a CASPCOrderAsset Object
	logger.Infof("CreateDCOrderAsset: Arguments for ledgerapi %s : ", args[0])
	err = json.Unmarshal([]byte(args[0]), &order)
	logger.Infof("CreateDCOrderAsset :OrderId ID is : %s ", order.OrderID)
	timeinfo, err := stub.GetTxTimestamp()
	if err != nil {
		logger.Errorf("CreateDCOrderAsset : Error getting  timestamp  %s", err)
		return shim.Error(fmt.Sprintf("CreateDCOrderAsset :DCOrderAsset object create failed due to timestamp read %s", err))
	}
	logger.Infof("CreateDCOrderAsset: Time stamp is %+v ", timeinfo)
	for i := 0; i < len(order.ParentInfo); i++ {
		var keys []string
		/*type DCOrderAsset struct {
			DocType        string              `json:"docType,omitempty"`
			PcID           string              `json:"PCID,omitempty"`
			DcID           string              `json:"DCID,omitempty"`
			Produce        string              `json:"Produce,omitempty"`
			Variety        string              `json:"Variety,omitempty"`
			RequiredDate   string              `json:"requiredDate,omitempty"`
			Qty            uint64              `json:"QTY,omitempty"`
			ProduceID      string              `json:"PRID,omitempty"`
			TableVarieties TableVarietyDC      `json:"QTY,omitempty"`
			PCOrderID      string              `json:"ParentOrderId,omitempty"`
			OrderID        string              `json:"ORDERID,omitempty"`
			ChildOrderID   string              `json:"ChildOrderID,omitempty"`
			Status         string              `json:"STATUS,omitempty"`
			Transports     []TransportaionPCDC `json:"Transportaion,omitempty"`
			TotaltransportationPrice float64 `json:"totalTransportationCost,omitempty"`
		}*/

		var asset = DCOrderAsset{DC_ORDER_ASSET_RECORD_TYPE, order.OrderType, order.SourceID, order.DestinationID, timeinfo.GetSeconds(), order.Produce, order.Variety, order.RequiredDate,
			order.Qty, order.ParentInfo[i].ProduceID, order.ParentInfo[i].TableVarieties, order.ParentInfo[i].ParentOrderId,
			order.OrderID, order.ParentInfo[i].ChildOrderID, order.Status, order.Transports, order.TotaltransportationPrice}
		keys = append(keys, order.ParentInfo[i].ChildOrderID)
		Avalbytes, _ = DCOrderAssettoJson(asset)
		err = CreateAsset(stub, DC_ORDER_ASSET_RECORD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreateDCOrderAsset : Error inserting Object into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreateDCOrderAsset : DCOrderAsset object create failed %s", err))
		}
		AvalbytesRes = append(AvalbytesRes, Avalbytes...)
	}
	return shim.Success([]byte(AvalbytesRes))
}
func QueryDCOrderAsset(stub shim.ChaincodeStubInterface, args []string) *DCOrderAsset {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateDCOrderAsset : Incorrect number of arguments.")
		return nil
	}
	// Convert the arg to a DCOrderAsset Object
	//logger.Infof("CreateDCOrderAsset: Arguments for ledgerapi %s : ", args[0])
	logger.Infof("UpdatePCDCOrderStatus :OrderId ID is : %s ", args[0]) // take the invoice ID
	keys = append(keys, args[0])
	Avalbytes, err = QueryAsset(stub, DC_ORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("CreateDCOrderAsset : Error inserting Object into LedgerState %s", err)
		return nil
	}
	csOrder := DCOrderAsset{}
	csOrder, _ = JsontoDCOrderAsset(Avalbytes)
	return &csOrder
}

// ListDCOrderAssetbyDcID  Function will  query  all record from DB with speficic storage id
func ListDCOrderAssetbyDcID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var orderitr shim.StateQueryIteratorInterface
	var orderList []DCOrderAsset
	if len(args) < 1 {
		return shim.Error("ListDCOrderAssetbyDcID :Incorrect number of arguments. Expecting StorageId")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"DESTINATIONID\":\"%s\"}}", DC_ORDER_ASSET_RECORD_TYPE, args[0])
	logger.Infof("ListDCOrderAssetbyDcID Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListDCOrderAssetbyDcID : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ListDCOrderAssetbyDcID : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListDCOrderAssetbyDcID: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		order, _ := JsontoDCOrderAsset([]byte(databyte))
		orderList = append(orderList, order)
	}
	Avalbytes, err = json.Marshal(orderList)
	logger.Infof("ListDCOrderAssetbyDcID Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListDCOrderAssetbyDcID : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListDCOrderAssetbyDcID: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// ListDCOrderAsset  Function will  query  all record from DB with speficic storage id
func ListDCOrderAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var orderitr shim.StateQueryIteratorInterface
	var orderList []DCOrderAsset
	var keys []string
	if len(args) < 1 {
		return shim.Error("ListDCOrderAsset :Incorrect number of arguments. Expecting StorageId")
	}
	orderitr, err = ListAllAsset(stub, DC_ORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListDCOrderAsset : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ListDCOrderAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListDCOrderAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		order, _ := JsontoDCOrderAsset([]byte(databyte))
		orderList = append(orderList, order)
	}
	Avalbytes, err = json.Marshal(orderList)
	logger.Infof("ListDCOrderAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListDCOrderAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListDCOrderAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

func QueryDCOrderAssetbyOrderID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryDCOrderAssetbyOrderID : Incorrect number of arguments.")
		return shim.Error("QueryDCOrderAssetbyOrderID : Incorrect number of arguments.")
	}
	logger.Infof("QueryDCOrderAssetbyOrderID :Order ID is : %s ", args[0])

	keys = append(keys, args[0])
	Avalbytes, err = QueryAsset(stub, DC_ORDER_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryDCOrderAssetbyOrderID : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryDCOrderAssetbyOrderID : DCOrderAsset object get failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

//GetDCPCOrderList function is used to get PC orders created from a CAS order

func GetDCPCOrderList(stub shim.ChaincodeStubInterface, args []string) ([]string, bool) {
	var err error
	var orderitr shim.StateQueryIteratorInterface
	var orderlist []string
	if len(args) < 1 {
		logger.Infof("No of args are not 1  ")
		return orderlist, false
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"ParentOrderId\":\"%s\"},\"fields\":[\"ChildOrderID\"]}", DC_ORDER_ASSET_RECORD_TYPE, args[0])
	logger.Infof("GetDCPCOrderList Query string is %s ", queryString)

	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("GetDCPCOrderList : Instence not found in ledger")
		return orderlist, false
	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("GetDCPCOrderList : Cannot parse result set. Error : %v", derr)
			return orderlist, false
		}
		databyte := data.GetValue()
		var resultdata AssetDataChild
		json.Unmarshal(databyte, &resultdata)
		logger.Infof("result is %v   ", resultdata)
		orderlist = append(orderlist, resultdata.ChildOrderID)
	}
	logger.Infof("GetDCPCOrderList list   is %v ", orderlist)
	return orderlist, true

}

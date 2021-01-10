package main

import (
	"encoding/json"
	"fmt"
	"sort"

	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//OrderInfo to track order information
type OrderInfo struct {
	OrderUnixTime int64  `json:"orderUnixTime,omitempty"`
	SourceID      string `json:"SOURCEID,omitempty"`
	DestinationID string `json:"DESTINATIONID,omitempty"`
	OrderID       string `json:"ORDERID,omitempty"`
	OrderType     uint64 `json:"ORDERTYPE,omitempty"`
}

// VerityInfo structure
type VerityInfo struct {
	Name       string      `json:"NAME,omitempty"`
	OrderInfos []OrderInfo `json:"TRACKINGINFOS,omitempty"`
}

// ProduceInfo structure
type ProduceInfo struct {
	ProduceID string       `json:"PRID,omitempty"`
	Varieties []VerityInfo `json:"VARIETIES,omitempty"`
}

// GetProduceTracking  Function will  query  record in ledger based on ID
func GetProduceTracking(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Avalbytes []byte
	producetrackinginfo := ProduceInfo{}
	produce, presult := QueryProduceAssetInfo(stub, args)

	if presult == true {
		producetrackinginfo.ProduceID = produce.ProduceID
		varietycount := len(produce.ProduceQuantities)
		logger.Infof("GetProduceTracking: Total variety fount is  %d : ", varietycount)
		if varietycount == 0 {
			logger.Errorf("GetProduceTracking : No variety found in  Ledger  for produceId  id as %s", args[0])
			return shim.Error(fmt.Sprintf("GetProduceTracking : No variety found in  Ledger  for produceId  id  as %s", args[0]))
		} //Variety count end loop
		var verityInfoList []VerityInfo
		//Variety parsing loop start
		for vid := 0; vid < varietycount; vid++ {
			logger.Infof("GetProduceTracking: Total variety Name is  %s : ", produce.ProduceQuantities[vid].VarietyType)
			var veritykey []string
			verityinfo := VerityInfo{}
			verityinfo.Name = produce.ProduceQuantities[vid].VarietyType
			veritykey = append(veritykey, produce.ProduceID, produce.ProduceQuantities[vid].VarietyType)

			verityinfo.OrderInfos = ProduceOrderList(stub, veritykey)
			logger.Infof("GetProduceTracking: verityinfo with order details is  %+v : ", verityinfo)
			verityInfoList = append(verityInfoList, verityinfo)
			logger.Infof("verityInfoList  is %v", verityInfoList)
		}
		producetrackinginfo.Varieties = verityInfoList
		logger.Infof("producetrackinginfo  is %v", producetrackinginfo)

	} else {
		logger.Errorf("GetProduceTracking : No Produce found in  Ledger  with  id as %s", args[0])
		return shim.Error(fmt.Sprintf("GetProduceTracking : No Produce found in  Ledger  with  id as %s", args[0]))
	}
	Avalbytes, _ = json.Marshal(&producetrackinginfo)

	return shim.Success([]byte(Avalbytes))
}

// ProduceOrderList to get list of all order associated with this produce id and variety
func ProduceOrderList(stub shim.ChaincodeStubInterface, keys []string) []OrderInfo {
	var err error
	var orderitr shim.StateQueryIteratorInterface
	var orderlist []OrderInfo
	var queryString string

	logger.Infof("ProduceOrderList:args is : %v ", keys)
	queryString = fmt.Sprintf("{\"selector\":{\"PRID\":\"%s\",\"variety\":\"%s\"},\"fields\":[\"ORDERID\",\"ORDERTYPE\",\"SOURCEID\",\"DESTINATIONID\",\"orderUnixTime\"]}", keys[0], keys[1])
	logger.Infof("ProduceOrderList Query string is %s ", queryString)
	orderitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ProduceOrderList : Instance not found in ledger")
		return nil
	}
	defer orderitr.Close()
	for orderitr.HasNext() {
		data, derr := orderitr.Next()
		if derr != nil {
			logger.Errorf("ProduceOrderList : Cannot parse result set. Error : %v", derr)
			return nil
		}
		databyte := data.GetValue()
		orderinfo := OrderInfo{}
		err = json.Unmarshal(databyte, &orderinfo)
		orderlist = append(orderlist, orderinfo)
	}
	logger.Infof("ProduceOrderList Query result  is %v ", orderlist)
	sort.Slice(orderlist, func(i, j int) bool { return orderlist[i].OrderUnixTime < orderlist[j].OrderUnixTime })
	logger.Infof("ProduceOrderList Query result  after sorting is %v ", orderlist)
	return orderlist
}

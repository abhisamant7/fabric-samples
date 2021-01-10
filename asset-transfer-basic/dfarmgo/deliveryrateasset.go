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
	DELIVERY_ASSET_RECORD_TYPE string = "DELIVERYRate"
)

type DeliveryRateReq struct {
	State       string         `json:"STATE,omitempty"`
	Country     string         `json:"COUNTRY,omitempty"`
	ProduceName string         `json:"PRODUCE,omitempty"`
	Role        int            `json:"ROLE,omitempty"`
	PriceRate   []DeliveryRate `json:"PRICERATE,omitempty"`
}
type DeliveryRateAsset struct {
	DocType     string         `json:"docType,omitempty"`
	State       string         `json:"STATE,omitempty"`
	Country     string         `json:"COUNTRY,omitempty"`
	ProduceName string         `json:"PRODUCE,omitempty"`
	PriceRate   []DeliveryRate `json:"PRICERATE,omitempty"`
}

type DeliveryRate struct {
	CurrencyUnit string      `json:"CURRENCY,omitempty"`
	QuantityUnit interface{} `json:"QUANTITYUNIT,omitempty"`
	DistanceUnit string      `json:"DISTANCEUNIT,omitempty"`
	Value        float64     `json:"VALUE,omitempty"`
}

func CreateDeliveryRatAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	if len(args) < 1 {
		logger.Errorf("CreateDeliveryRatAsset : Incorrect number of arguments.")
		return shim.Error("CreateDeliveryRatAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("CreateDeliveryRatAsset :Attribute approle not found to create ParticipantRateAsset"))
	}
	if "ADMIN" != strings.ToUpper(value) {
		return shim.Error(fmt.Sprintf("CreateDeliveryRatAsset :This User is not allowed to create ParticipantRateAsset"))
	}
	// Convert the arg to a ParticipantRateAsset Object
	logger.Infof("CreateDeliveryRatAsset: Arguments for ledgerapi %s : ", args[0])
	asset := DeliveryRateReq{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	logger.Infof("CreateDeliveryRatAsset :Country is : %s ", asset.Country)
	logger.Infof("CreateDeliveryRatAsset :state is : %s ", asset.State)
	logger.Infof("CreateDeliveryRatAsset :ProduceName is : %s ", asset.ProduceName)
	logger.Infof("CreateDeliveryRatAsset :DeliverRateReq is : %v ", asset)
	keys = append(keys, asset.Country)
	keys = append(keys, asset.State)
	keys = append(keys, asset.ProduceName)

	Avalbytes, err = QueryAsset(stub, DELIVERY_ASSET_RECORD_TYPE, keys)
	if err != nil {
		participant := DeliveryRateAsset{}
		participant.DocType = DELIVERY_ASSET_RECORD_TYPE
		participant.State = asset.State
		participant.Country = asset.Country
		participant.ProduceName = asset.ProduceName
		participant.PriceRate = asset.PriceRate
		Avalbytes, _ = json.Marshal(participant)
		logger.Infof("CreateDeliveryRatAsset :DeliveryRate Asset is : %s ", participant)
		logger.Infof("CreateDeliveryRatAsset :PCDeliveryRate: %s ", participant.PriceRate)
		err = CreateAsset(stub, DELIVERY_ASSET_RECORD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreateDeliveryRatAsset : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreateDeliveryRatAsset : ParticipantRateAsset Object first time create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}
	objread := DeliveryRateAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &objread)
	objread.PriceRate = asset.PriceRate
	Avalbytes, _ = json.Marshal(objread)
	logger.Infof("CreateDeliveryRatAsset :DeliverRate Asset is : %s ", objread)
	logger.Infof("CreateDeliveryRatAsset :DeliverRate  is : %s ", objread.PriceRate)

	err = UpdateAssetWithoutGet(stub, DELIVERY_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateDeliveryRatAsset : Error inserting Object first time  into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateDeliveryRatAsset : DeliveryRate Object first time create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

//getDeliveryRatAsset to get DeliveryRatAsset from ledger based on 3 param
func getDeliveryRatAsset(stub shim.ChaincodeStubInterface, args []string) (DeliveryRate, bool) {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) == 0 {
		logger.Errorf("DeliveryRatAsset : Incorrect number of arguments.")
		return DeliveryRate{}, false
	}
	logger.Infof("DeliveryRatAsset :args is : %v ", args)
	keys = append(keys, args[0])
	keys = append(keys, args[1])
	keys = append(keys, args[2])

	Avalbytes, err = QueryAsset(stub, DELIVERY_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("DeliveryRatAsset : Error Querying Object from LedgerState %s", err)
		return DeliveryRate{}, false
	}
	readobj := DeliveryRateAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &readobj)

	return readobj.PriceRate[0], true
}

// ListAllDeliveryRatAsset to list all dfarm rate

func ListAllDeliveryRatAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var keys []string
	var Avalbytes []byte
	var pitr shim.StateQueryIteratorInterface
	var DeliveryRatAssetList []DeliveryRateAsset

	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("ListAllDeliveryRatAsset :This User is not allowed to List ParticipantRateAsset"))
	}
	user := strings.ToUpper(value)
	if "ADMIN" != user {
		return shim.Error(fmt.Sprintf("ListAllDeliveryRatAsset :This User is not allowed to List ParticipantRateAsset"))
	}

	pitr, err = ListAllAsset(stub, DELIVERY_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListAllDeliveryRatAsset : Instence not found in ledger")
		return shim.Error("pitr : Instence not found in ledger")

	}
	defer pitr.Close()
	for pitr.HasNext() {
		data, derr := pitr.Next()
		if derr != nil {
			logger.Errorf("ListAllDeliveryRatAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListAllDeliveryRatAsset: Cannot parse result set. Error : %v", derr))

		}
		databytes := data.GetValue()
		objread := DeliveryRateAsset{}
		_ = json.Unmarshal([]byte(databytes), &objread)

		DeliveryRatAssetList = append(DeliveryRatAssetList, objread)

	}
	Avalbytes, err = json.Marshal(DeliveryRatAssetList)
	logger.Infof("ListAllDeliveryRatAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListAllDeliveryRatAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListAllDeliveryRatAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

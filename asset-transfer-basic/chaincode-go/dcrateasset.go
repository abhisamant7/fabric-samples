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
	DC_ASSET_RECORD_TYPE        string = "DCRate"
	MARKUP_DISCOUNT_RECORD_TYPE string = "MarkupDiscountPer"
)

type DCRateReq struct {
	State       string  `json:"STATE,omitempty"`
	Country     string  `json:"COUNTRY,omitempty"`
	ProduceName string  `json:"PRODUCE,omitempty"`
	Role        int     `json:"ROLE,omitempty"`
	PriceRate   float32 `json:"PRICERATE,omitempty"`
}
type DCRateAsset struct {
	DocType     string  `json:"docType,omitempty"`
	State       string  `json:"STATE,omitempty"`
	Country     string  `json:"COUNTRY,omitempty"`
	ProduceName string  `json:"PRODUCE,omitempty"`
	Role        int     `json:"ROLE,omitempty"`
	PriceRate   float32 `json:"PRICERATE,omitempty"`
}

// CreateDCRateAsset to define DC rate
func CreateDCRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	if len(args) < 1 {
		logger.Errorf("CreateDCRateAsset : Incorrect number of arguments.")
		return shim.Error("CreateDCRateAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("CreateDCRateAsset :Attribute approle not found to create ParticipantRateAsset"))
	}
	if "ADMIN" != strings.ToUpper(value) {
		return shim.Error(fmt.Sprintf("CreateDCRateAsset :This User is not allowed to create ParticipantRateAsset"))
	}
	// Convert the arg to a ParticipantRateAsset Object
	logger.Infof("CreateDCRateAsset: Arguments for ledgerapi %s : ", args[0])
	asset := DCRateReq{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	logger.Infof("CreateDeliveryRatAsset :Country is : %s ", asset.Country)
	logger.Infof("CreateDeliveryRatAsset :state is : %s ", asset.State)
	logger.Infof("CreateDeliveryRatAsset :Role is : %s ", asset.Role)
	logger.Infof("CreateDeliveryRatAsset :ProduceName is : %s ", asset.ProduceName)
	logger.Infof("CreateDeliveryRatAsset :PCRateReq is : %v ", asset)
	keys = append(keys, asset.Country)
	keys = append(keys, asset.State)
	keys = append(keys, asset.ProduceName)
	Avalbytes, err = QueryAsset(stub, DC_ASSET_RECORD_TYPE, keys)
	if err != nil {
		participant := DCRateAsset{}
		participant.DocType = DC_ASSET_RECORD_TYPE
		participant.State = asset.State
		participant.Country = asset.Country
		participant.Role = asset.Role
		participant.ProduceName = asset.ProduceName
		participant.PriceRate = asset.PriceRate
		Avalbytes, _ = json.Marshal(participant)
		logger.Infof("CreateDCRateAsset :DeliveryRatAsset  is : %s ", participant)
		logger.Infof("CreateDCRateAsset : DeliveryRatAsset  rate is : %s ", participant.PriceRate)

		err = CreateAsset(stub, DC_ASSET_RECORD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreateDCRateAsset : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreateDCRateAsset : ParticipantRateAsset Object first time create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}
	readobj := DCRateAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &readobj)
	readobj.PriceRate = asset.PriceRate
	Avalbytes, _ = json.Marshal(readobj)
	logger.Infof("CreateDCRateAsset :DCRateAsset Asset is : %s ", readobj)
	logger.Infof("CreateDCRateAsset :DCRateAsset Rate Asset is : %s ", readobj.PriceRate)

	err = UpdateAssetWithoutGet(stub, DC_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateDCRateAsset : Error inserting Object first time  into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateDCRateAsset : DCRateAsset Object first time create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

//getDCRateAsset to get DCRateAsset from ledger based on 3 param
func getDCRateAsset(stub shim.ChaincodeStubInterface, args []string) (float32, bool) {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) == 0 {
		logger.Errorf("getDCRateAsset : Incorrect number of arguments.")
		return 0.0, false
	}
	logger.Infof("getDCRateAsset :args is : %v ", args)
	keys = append(keys, args[0])
	keys = append(keys, args[1])
	keys = append(keys, args[2])

	Avalbytes, err = QueryAsset(stub, DC_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("getDCRateAsset : Error Querying Object from LedgerState %s", err)
		return 0.0, false
	}
	readobj := DCRateAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &readobj)

	return readobj.PriceRate, true
}

// ListDCRateAsset  Function will  query  record in ledger based on ID
func ListDCRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var pricerateitr shim.StateQueryIteratorInterface
	var priceratelist []DCRateAsset
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("ListDCRateAsset :Attribute approle not found to List DCRateAsset"))
	}
	user := strings.ToUpper(value)
	if user != "ADMIN" {
		return shim.Error(fmt.Sprintf("ListDCRateAsset :This User is not allowed to List DCRateAsset "))
	}

	pricerateitr, err = ListAllAsset(stub, DC_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListDCRateAsset : Instence not found in ledger")
		return shim.Error("pricerateitr : Instence not found in ledger")

	}
	defer pricerateitr.Close()
	for pricerateitr.HasNext() {
		data, derr := pricerateitr.Next()
		if derr != nil {
			logger.Errorf("ListDCRateAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListDCRateAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		readobj := DCRateAsset{}
		_ = json.Unmarshal([]byte(databyte), &readobj)
		priceratelist = append(priceratelist, readobj)
	}
	Avalbytes, err = json.Marshal(priceratelist)
	logger.Infof("ListDCRateAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListDCRateAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListDCRateAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))

}

//MarkupDiscountPer asset to provide Markup and discount on final Buyer payment for a produce Variety

type MarkupDiscountPerAsset struct {
	DocType       string  `json:"docType,omitempty"`
	SourceID      string  `json:"SOURCEID,omitempty"`
	DestinationID string  `json:"DESTINATIONID,omitempty"`
	ProduceName   string  `json:"PRODUCE,omitempty"`
	Variety       string  `json:"VARIETY,omitempty"`
	PerMode       bool    `json:"PERMODE,omitempty"` //0 means markuo 1 means discount
	PerValue      float64 `json:"PERVALUE,omitempty"`
}

// CreateMarkupDiscountPerAsset to define DC rate
func CreateMarkupDiscountPerAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	if len(args) < 1 {
		logger.Errorf("CreateMarkupDiscountPerAsset : Incorrect number of arguments.")
		return shim.Error("CreateMarkupDiscountPerAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("CreateMarkupDiscountPerAsset :Attribute approle not found to create ParticipantRateAsset"))
	}
	if "ADMIN" != strings.ToUpper(value) {
		return shim.Error(fmt.Sprintf("CreateMarkupDiscountPerAsset :This User is not allowed to create ParticipantRateAsset"))
	}
	// Convert the arg to a MarkupDiscountPerAsset Object
	logger.Infof("CreateMarkupDiscountPerAsset: Arguments for ledgerapi %s : ", args[0])
	asset := MarkupDiscountPerAsset{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	logger.Infof("CreateMarkupDiscountPerAsset :SourceId is : %s ", asset.SourceID)
	logger.Infof("CreateMarkupDiscountPerAsset :DestinationID is : %s ", asset.DestinationID)
	logger.Infof("CreateMarkupDiscountPerAsset :ProduceName is : %s ", asset.ProduceName)
	logger.Infof("CreateMarkupDiscountPerAsset :Variety is : %v ", asset.Variety)
	keys = append(keys, asset.SourceID)
	keys = append(keys, asset.DestinationID)
	keys = append(keys, asset.ProduceName)
	keys = append(keys, asset.Variety)
	Avalbytes, err = QueryAsset(stub, MARKUP_DISCOUNT_RECORD_TYPE, keys)
	if err != nil {
		perasset := MarkupDiscountPerAsset{}
		perasset.DocType = MARKUP_DISCOUNT_RECORD_TYPE
		Avalbytes, _ = json.Marshal(perasset)
		logger.Infof("CreateMarkupDiscountPerAsset :MarkupDiscountPerAsset  is : %s ", Avalbytes)
		//logger.Infof("CreateMarkupDiscountPerAsset : MarkupDiscountPerAsset  percentage is : %f ", perasset.PerValue)

		err = CreateAsset(stub, MARKUP_DISCOUNT_RECORD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreateMarkupDiscountPerAsset : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreateMarkupDiscountPerAsset : ParticipantRateAsset Object first time create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}
	readobj := MarkupDiscountPerAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &readobj)
	readobj.PerMode = asset.PerMode
	readobj.PerValue = asset.PerValue
	Avalbytes, _ = json.Marshal(readobj)
	logger.Infof("CreateMarkupDiscountPerAsset :MarkupDiscountPerAsset Asset is : %v ", readobj)
	err = UpdateAssetWithoutGet(stub, MARKUP_DISCOUNT_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateMarkupDiscountPerAsset : Error inserting Object first time  into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateMarkupDiscountPerAsset : MarkupDiscountPerAsset Object first time create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

//getMarkupDiscountPerAsset to get MarkupDiscountPerAsset from ledger based on 3 param
func getMarkupDiscountPerAsset(stub shim.ChaincodeStubInterface, args []string) (MarkupDiscountPerAsset, bool) {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) == 0 {
		logger.Errorf("getMarkupDiscountPerAsset : Incorrect number of arguments.")
		return MarkupDiscountPerAsset{}, false
	}
	logger.Infof("getMarkupDiscountPerAsset :args is : %v ", args)
	keys = append(keys, args[0])
	keys = append(keys, args[1])
	keys = append(keys, args[2])
	keys = append(keys, args[3])

	Avalbytes, err = QueryAsset(stub, MARKUP_DISCOUNT_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("getMarkupDiscountPerAsset : Error Querying Object from LedgerState %s", err)
		return MarkupDiscountPerAsset{}, false
	}
	readobj := MarkupDiscountPerAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &readobj)

	return readobj, true
}

// ListMarkupDiscountPerAsset  Function will  query  record in ledger based on ID
func ListMarkupDiscountPerAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var itr shim.StateQueryIteratorInterface
	var perlist []MarkupDiscountPerAsset
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("ListMarkupDiscountPerAsset :Attribute approle not found to List MarkupDiscountPerAsset"))
	}
	user := strings.ToUpper(value)
	if user != "ADMIN" {
		return shim.Error(fmt.Sprintf("ListMarkupDiscountPerAsset :This User is not allowed to List MarkupDiscountPerAsset "))
	}

	itr, err = ListAllAsset(stub, MARKUP_DISCOUNT_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListMarkupDiscountPerAsset : Instence not found in ledger")
		return shim.Error("pricerateitr : Instence not found in ledger")

	}
	defer itr.Close()
	for itr.HasNext() {
		data, derr := itr.Next()
		if derr != nil {
			logger.Errorf("ListMarkupDiscountPerAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListMarkupDiscountPerAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		readobj := MarkupDiscountPerAsset{}
		_ = json.Unmarshal([]byte(databyte), &readobj)
		perlist = append(perlist, readobj)
	}
	Avalbytes, err = json.Marshal(perlist)
	logger.Infof("ListMarkupDiscountPerAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListMarkupDiscountPerAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListMarkupDiscountPerAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))

}

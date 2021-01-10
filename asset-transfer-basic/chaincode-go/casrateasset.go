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
	CAS_ASSET_RECORD_TYPE string = "CASRate"
)

type CASRateReq struct {
	State       string    `json:"STATE,omitempty"`
	Country     string    `json:"COUNTRY,omitempty"`
	Role        int       `json:"ROLE,omitempty"`
	CASID       string    `json:"CASID,omitempty"`
	ProduceName string    `json:"PRODUCE,omitempty"`
	PriceRate   []CASRate `json:"CASRATES,omitempty"`
}
type CASRateAsset struct {
	DocType     string    `json:"docType,omitempty"`
	Country     string    `json:"COUNTRY,omitempty"`
	State       string    `json:"STATE,omitempty"`
	CASID       string    `json:"CASID,omitempty"`
	ProduceName string    `json:"PRODUCE,omitempty"`
	PriceRate   []CASRate `json:"CASRATES,omitempty"`
}

type CASRate struct {
	CurrencyUnit string      `json:"CURRENCY,omitempty"`
	QuantityUnit interface{} `json:"QUANTITYUNIT,omitempty"`
	DurationUnit string      `json:"DURATIONUNIT,omitempty"`
	Value        float64     `json:"VALUE,omitempty"`
}

//*********************** CASRate Asset JSON Method *************************
// //Convert JSON  object to CASRate Asset
func JsontoCASRateAsset(data []byte) (CASRateAsset, error) {
	obj := CASRateAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to PCPackagingRateAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert CASRate Asset object to Json Message

func CASRateAssettoJson(obj CASRateAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateCASRateAsset Function will  insert record in ledger after receiving request from Client Application
func CreateCASRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	if len(args) < 1 {
		logger.Errorf("CreateCASRateAsset : Incorrect number of arguments.")
		return shim.Error("CreateCASRateAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("CreateCASRateAsset :Attribute approle not found to create CASRateAsset"))
	}
	user := strings.ToUpper(value)
	if user != "COLDSTORAGE" {
		return shim.Error(fmt.Sprintf("CreateCASRateAsset :This User is not allowed to create CASRateAsset "))
	}
	// Convert the arg to a ParticipantRateAsset Object
	logger.Infof("CreateCASRateAsset: Arguments for ledgerapi %s : ", args[0])
	asset := CASRateReq{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	logger.Infof("CreateCASRateAsset :state is : %s ", asset.CASID)
	logger.Infof("CreateCASRateAsset :CASRateReq is : %v ", asset)
	keys = append(keys, asset.CASID)
	Avalbytes, err = QueryAsset(stub, CAS_ASSET_RECORD_TYPE, keys)
	if err != nil {
		ledgerasset := CASRateAsset{CAS_ASSET_RECORD_TYPE, asset.Country, asset.State, asset.CASID, asset.ProduceName, asset.PriceRate}
		Avalbytes, _ = CASRateAssettoJson(ledgerasset)
		logger.Infof("CreateCASRateAsset :Ledger Data is  : %s ", Avalbytes)
		err = CreateAsset(stub, CAS_ASSET_RECORD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreateCASRateAsset : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreateCASRateAsset :  Object first time create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}
	casread, _ := JsontoCASRateAsset([]byte(Avalbytes))
	casread.PriceRate = asset.PriceRate
	Avalbytes, _ = CASRateAssettoJson(casread)
	logger.Infof("CreateCASRateAsset :CAS  Asset is : %s ", Avalbytes)
	err = UpdateAssetWithoutGet(stub, CAS_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateCASRateAsset : Error inserting Object first time  into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateCASRateAsset : CAS Object update failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

//QueryCASRateAsset to get CASRateAsset from ledger based on CASID
func QueryCASRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryCASRateAsset : Incorrect number of arguments.")
		return shim.Error("QueryCASRateAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("QueryCASRateAsset :Attribute approle not found to query CASRateAsset"))
	}
	user := strings.ToUpper(value)
	if user != "ADMIN" && user != "COLDSTORAGE" {
		return shim.Error(fmt.Sprintf("QueryCASRateAsset :This User is not allowed to Query CASRateAsset "))
	}
	logger.Infof("QueryCASRateAsset :CASID is : %s ", args[0])
	keys = append(keys, args[0])

	Avalbytes, err = QueryAsset(stub, CAS_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryCASRateAsset : Error Querying Object from LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryCASRateAsset : CASRateAsset object get failed %s", err))
	}

	return shim.Success([]byte(Avalbytes))
}

//getCASRateAsset to get CASRateAsset from ledger based on CASID
func getCASRateAsset(stub shim.ChaincodeStubInterface, casid string) (CASRate, bool) {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(casid) == 0 {
		logger.Errorf("getCASRateAsset : Incorrect number of arguments.")
		return CASRate{}, false
	}

	logger.Infof("getCASRateAsset :CASID is : %s ", casid)
	keys = append(keys, casid)

	Avalbytes, err = QueryAsset(stub, CAS_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("getCASRateAsset : Error Querying Object from LedgerState %s", err)
		return CASRate{}, false
	}
	casrate, _ := JsontoCASRateAsset(Avalbytes)

	return casrate.PriceRate[0], true
}

// ListCASRateAsset  Function will  query  record in ledger based on ID
func ListCASRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var pricerateitr shim.StateQueryIteratorInterface
	var priceratelist []CASRateAsset
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("QueryCASRateAsset :Attribute approle not found to List CASRateAsset"))
	}
	user := strings.ToUpper(value)
	if user != "ADMIN" && user != "COLDSTORAGE" {
		return shim.Error(fmt.Sprintf("QueryCASRateAsset :This User is not allowed to List CASRateAsset "))
	}

	pricerateitr, err = ListAllAsset(stub, CAS_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListCASRateAsset : Instence not found in ledger")
		return shim.Error("pricerateitr : Instence not found in ledger")

	}
	defer pricerateitr.Close()
	for pricerateitr.HasNext() {
		data, derr := pricerateitr.Next()
		if derr != nil {
			logger.Errorf("ListCASRateAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListCASRateAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		pricerate, _ := JsontoCASRateAsset([]byte(databyte))
		priceratelist = append(priceratelist, pricerate)
	}
	Avalbytes, err = json.Marshal(priceratelist)
	logger.Infof("ListCASRateAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListCASRateAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListCASRateAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))

}

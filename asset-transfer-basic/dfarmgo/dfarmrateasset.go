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
	DFARM_ASSET_RECORD_TYPE string = "DFARMRate"
)

type DfarmRateReq struct {
	State       string   `json:"STATE,omitempty"`
	Country     string   `json:"COUNTRY,omitempty"`
	ProduceName string   `json:"PRODUCE,omitempty"`
	Role        int      `json:"ROLE,omitempty"`
	PriceRate   DfarmPer `json:"PRICERATE,omitempty"`
}
type DfarmRateAsset struct {
	DocType     string   `json:"docType,omitempty"`
	State       string   `json:"STATE,omitempty"`
	Country     string   `json:"COUNTRY,omitempty"`
	ProduceName string   `json:"PRODUCE,omitempty"`
	PriceRate   DfarmPer `json:"PRICERATE,omitempty"`
}

type DfarmPer struct {
	Totalper     float64 `json:"DFARMPER,omitempty"`
	MaxAgentPer  float64 `json:"MAXAGENTPER,omitempty"`
	DemurragePer float64 `json:"DEMURRAGEPER,omitempty"`
}

func CreatedFarmRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	if len(args) < 1 {
		logger.Errorf("CreatedFarmRateAsset : Incorrect number of arguments.")
		return shim.Error("CreatedFarmRateAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("CreatedFarmRateAsset :Attribute approle not found to create ParticipantRateAsset"))
	}
	if "ADMIN" != strings.ToUpper(value) {
		return shim.Error(fmt.Sprintf("CreatedFarmRateAsset :This User is not allowed to create ParticipantRateAsset"))
	}
	// Convert the arg to a ParticipantRateAsset Object
	logger.Infof("CreatedFarmRateAsset: Arguments for ledgerapi %s : ", args[0])
	asset := DfarmRateReq{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	logger.Infof("CreateDeliveryRatAsset :Country is : %s ", asset.Country)
	logger.Infof("CreateDeliveryRatAsset :state is : %s ", asset.State)
	logger.Infof("CreateDeliveryRatAsset :ProduceName is : %s ", asset.ProduceName)
	logger.Infof("CreateDeliveryRatAsset :PCRateReq is : %v ", asset)
	keys = append(keys, asset.Country)
	keys = append(keys, asset.State)
	keys = append(keys, asset.ProduceName)
	Avalbytes, err = QueryAsset(stub, DFARM_ASSET_RECORD_TYPE, keys)
	if err != nil {
		participant := DfarmRateAsset{}
		participant.DocType = DFARM_ASSET_RECORD_TYPE
		participant.State = asset.State
		participant.Country = asset.Country
		participant.ProduceName = asset.ProduceName
		participant.PriceRate = asset.PriceRate
		Avalbytes, _ = json.Marshal(participant)
		logger.Infof("CreatedFarmRateAsset :DfarmRate Asset is : %s ", participant)
		logger.Infof("CreatedFarmRateAsset :DfarmRate   is : %s ", participant.PriceRate)

		err = CreateAsset(stub, DFARM_ASSET_RECORD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreatedFarmRateAsset : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreatedFarmRateAsset : ParticipantRateAsset Object first time create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}
	objread := DfarmRateAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &objread)
	objread.PriceRate = asset.PriceRate
	Avalbytes, _ = json.Marshal(objread)
	logger.Infof("CreatedFarmRateAsset :Participant Asset is : %s ", objread)
	logger.Infof("CreatedFarmRateAsset :Participant Asset is : %s ", objread.PriceRate)

	err = UpdateAssetWithoutGet(stub, DFARM_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreatedFarmRateAsset : Error inserting Object first time  into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreatedFarmRateAsset : FarmRateAsset Object first time create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

//getdFarmRateAssett to get dFarmRateAsset from ledger based on 3 param
func getdFarmRateAsset(stub shim.ChaincodeStubInterface, args []string) (DfarmPer, bool) {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) == 0 {
		logger.Errorf("dFarmRateAsset : Incorrect number of arguments.")
		return DfarmPer{}, false
	}
	logger.Infof("dFarmRateAsset :args is : %v ", args)
	keys = append(keys, args[0])
	keys = append(keys, args[1])
	keys = append(keys, args[2])

	Avalbytes, err = QueryAsset(stub, DFARM_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("dFarmRateAsset : Error Querying Object from LedgerState %s", err)
		return DfarmPer{}, false
	}
	readobj := DfarmRateAsset{}
	_ = json.Unmarshal([]byte(Avalbytes), &readobj)

	return readobj.PriceRate, true
}

// ListAlldfarmRateAsset to list all dfarm rate

func ListAlldfarmRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var keys []string
	var Avalbytes []byte
	var pitr shim.StateQueryIteratorInterface
	var DfarmRateAssetList []DfarmRateAsset

	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("ListAlldfarmRateAsset :This User is not allowed to List ParticipantRateAsset"))
	}
	user := strings.ToUpper(value)
	if "ADMIN" != user {
		return shim.Error(fmt.Sprintf("ListAlldfarmRateAsset :This User is not allowed to List ParticipantRateAsset"))
	}

	pitr, err = ListAllAsset(stub, DFARM_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListAlldfarmRateAsset : Instence not found in ledger")
		return shim.Error("pcpackagingitr : Instence not found in ledger")

	}
	defer pitr.Close()
	for pitr.HasNext() {
		data, derr := pitr.Next()
		if derr != nil {
			logger.Errorf("ListAlldfarmRateAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListAlldfarmRateAsset: Cannot parse result set. Error : %v", derr))

		}
		databytes := data.GetValue()
		objread := DfarmRateAsset{}
		_ = json.Unmarshal([]byte(databytes), &objread)

		DfarmRateAssetList = append(DfarmRateAssetList, objread)

	}
	Avalbytes, err = json.Marshal(DfarmRateAssetList)
	logger.Infof("ListAlldfarmRateAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListAlldfarmRateAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListAlldfarmRateAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

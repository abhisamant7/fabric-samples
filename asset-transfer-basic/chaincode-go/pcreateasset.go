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
	//PARTICIPANTRATE_ASSET_RECORD_TYPE string = "ParticipantPriceRate"
	PCPACKAGINGRATE_ASSET_RECORD_TYPE string = "PCPackagingRate"
)

type PCPackagingRateAsset struct {
	DocType          string             `json:"docType,omitempty"`
	ProduceName      string             `json:"PRODUCE,omitempty"`
	Country          string             `json:"COUNTRY,omitempty"`
	State            string             `json:"STATE,omitempty"`
	Variety          string             `json:"VARIETY,omitempty"`
	PCID             string             `json:"PCID,omitempty"`
	TableVerityPrice []QualityPriceRate `json:"TABLEVARIETYPRICE,omitempty"`
	Selectedunit     Unit               `json:"SELECTED_UNIT,omitempty"`
	Currecny         string             `json:"CURRENCY,omitempty"`
}

//*********************** PCPackagingRate Asset JSON Method *************************
// //Convert JSON  object to PCPackagingRate Asset
func JsontoPCPackagingRateAsset(data []byte) (PCPackagingRateAsset, error) {
	obj := PCPackagingRateAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to PCPackagingRateAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert PCPackagingRateAsset object to Json Message

func PCPackagingRateAssettoJson(obj PCPackagingRateAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}



// CreatePCPackagingRateAsset Function will  insert record in ledger after receiving request from Client Application
func CreatePCPackagingRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var AvailsReads []byte
	var match bool

	if len(args) < 1 {
		logger.Errorf("CreatePCPackagingRateAsset : Incorrect number of arguments.")
		return shim.Error("CreatePCPackagingRateAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	logger.Infof("CreatePCPackagingRateAsset: User Role: %s  ", value)
	if !found {
		return shim.Error(fmt.Sprintf("CreatePCPackagingRateAsset :Attribute approle not found to create PCPackagingRateAsset"))
	}

	if "PC" != strings.ToUpper(value) {
		return shim.Error(fmt.Sprintf("CreatePCPackagingRateAsset :This User is not allowed to create PCPackagingRateAsset"))
	}
	// Convert the arg to a PCPackagingRateAsset Object
	logger.Infof("CreatePCPackagingRateAsset: Arguments for ledgerapi %s : ", args[0])
	asset := PCPackagingRateAsset{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	asset.DocType = PCPACKAGINGRATE_ASSET_RECORD_TYPE
	logger.Infof("CreatePCPackagingRateAsset :state is : %s ", asset.State)
	logger.Infof("CreatePCPackagingRateAsset :PCPackagingRateAsset is : %v ", asset)
	//Do we need to add Country?
	keys = append(keys, asset.Country)
	keys = append(keys, asset.State)
	keys = append(keys, asset.ProduceName)
	keys = append(keys, asset.Variety)
	keys = append(keys, asset.PCID)
	Avalbytes, err = QueryAsset(stub, PCPACKAGINGRATE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Infof("CreatePCPackagingRateAsset :Participant Asset is : %s ", args[0])
		Avalbytes, _ = PCPackagingRateAssettoJson(asset)
		err = CreateAsset(stub, PCPACKAGINGRATE_ASSET_RECORD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreatePCPackagingRateAsset : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreatePCPackagingRateAsset : PCPackagingRateAsset Object first time create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}
	prateread := PCPackagingRateAsset{}
	logger.Infof("Update case Availablevyte", string(Avalbytes))
	_ = json.Unmarshal(Avalbytes, &prateread) // not working
	logger.Infof("Data received from ledger %v", prateread)

	if len(asset.TableVerityPrice) > 0 {
		if len(prateread.TableVerityPrice) > 0 {
			for i := 0; i < len(asset.TableVerityPrice); i++ {
				for j := 0; j < len(prateread.TableVerityPrice); j++ {
					if asset.TableVerityPrice[i].QualityName == prateread.TableVerityPrice[j].QualityName {
						prateread.TableVerityPrice[j].Rates = asset.TableVerityPrice[i].Rates
						match = true
					}
				}
				if match == false {
					prateread.TableVerityPrice = append(prateread.TableVerityPrice, asset.TableVerityPrice[i])
				}
			}
		} else {
			prateread.TableVerityPrice = asset.TableVerityPrice
		}

		AvailsReads, _ = PCPackagingRateAssettoJson(prateread)
		logger.Infof("CreatePCPackagingRateAsset :Participant Asset is : %s ", prateread)
		logger.Infof("Available read in update", string(AvailsReads))
		err = UpdateAssetWithoutGet(stub, PCPACKAGINGRATE_ASSET_RECORD_TYPE, keys, AvailsReads)
		if err != nil {
			logger.Errorf("CreatePCPackagingRateAsset : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreatePCPackagingRateAsset : PCPackagingRateAsset Object first time create failed %s", err))
		}
	}
	return shim.Success([]byte(Avalbytes))

}

func QueryPCPackagingRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 5 {
		logger.Errorf("QueryPCPackagingRateAsset : Incorrect number of arguments.")
		return shim.Error("QueryPCPackagingRateAsset : Incorrect number of arguments.")
	}
	logger.Infof("QueryPCPackagingRateAsset :State is : %s ", args[0])
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	logger.Infof("QueryPCPackagingRateAsset :User Value is : %s ", value)
	if !found {
		return shim.Error(fmt.Sprintf("QueryPCPackagingRateAsset :This User is not allowed to Query PCPackagingRateAsset"))
	}
	user := strings.ToUpper(value)

	if user == "ADMIN" || user == "PC" {
		keys = append(keys, args[0]) //Country
		keys = append(keys, args[1]) //State
		keys = append(keys, args[2]) //Produce Name
		keys = append(keys, args[3]) //Variety
		keys = append(keys, args[4]) //pcID

		Avalbytes, err = QueryAsset(stub, PCPACKAGINGRATE_ASSET_RECORD_TYPE, keys)
		if err != nil {
			logger.Errorf("QueryPCPackagingRateAsset : Error Querying Object from LedgerState %s", err)
			return shim.Error(fmt.Sprintf("QueryPCPackagingRateAsset : PCPackagingRateAsset object get failed %s", err))
		}

		return shim.Success([]byte(Avalbytes))
	}

	return shim.Error(fmt.Sprintf("CreateCASRateAsset :This User is not allowed to Query ParticipantRateAsset %s", user))
}

func getPCPackagingRateAsset(stub shim.ChaincodeStubInterface, args []string) (PCPackagingRateAsset, bool) {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 5 {
		logger.Errorf("getPCPackagingRateAsset : Incorrect number of arguments.")
		return PCPackagingRateAsset{}, false
	}
	logger.Infof("getPCPackagingRateAsset :Country is : %s ", args[0])
	logger.Infof("getPCPackagingRateAsset :State is : %s ", args[1])
	logger.Infof("getPCPackagingRateAsset :Produce Name is : %s ", args[2])
	logger.Infof("getPCPackagingRateAsset :Variety is : %s ", args[3])
	logger.Infof("getPCPackagingRateAsset :pcID is : %s ", args[4])

	keys = append(keys, args[0]) //Country
	keys = append(keys, args[1]) //State
	keys = append(keys, args[2]) //Produce Name
	keys = append(keys, args[3]) //Variety
	keys = append(keys, args[4]) //pcID

	Avalbytes, err = QueryAsset(stub, PCPACKAGINGRATE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryPCPackagingRateAsset : Error Querying Object from LedgerState %s", err)
		return PCPackagingRateAsset{}, false
	}
	pcpackagingrate, _ := JsontoPCPackagingRateAsset(Avalbytes)

	return pcpackagingrate, true
}

// chaincode by add by me
func ListAllPCPackagingRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var keys []string
	var Avalbytes []byte
	var pcpackagingitr shim.StateQueryIteratorInterface
	var PCPackagingRateAssetList []PCPackagingRateAsset

	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("ListAllPCPackagingRateAsset :This User is not allowed to List ParticipantRateAsset"))
	}
	user := strings.ToUpper(value)
	if "ADMIN" != user && user != "PC" {
		return shim.Error(fmt.Sprintf("ListAllPCPackagingRateAsset :This User is not allowed to List ParticipantRateAsset"))
	}

	pcpackagingitr, err = ListAllAsset(stub, PCPACKAGINGRATE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListAllPCPackagingRateAsset : Instence not found in ledger")
		return shim.Error("pcpackagingitr : Instence not found in ledger")

	}
	defer pcpackagingitr.Close()
	for pcpackagingitr.HasNext() {
		data, derr := pcpackagingitr.Next()
		if derr != nil {
			logger.Errorf("ListAllPCPackagingRateAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListAllPCPackagingRateAsset: Cannot parse result set. Error : %v", derr))

		}
		databytes := data.GetValue()

		pcpackaging, _ := JsontoPCPackagingRateAsset([]byte(databytes))
		PCPackagingRateAssetList = append(PCPackagingRateAssetList, pcpackaging)

	}
	Avalbytes, err = json.Marshal(PCPackagingRateAssetList)
	logger.Infof("ListAllPCPackagingRateAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListAllPCPackagingRateAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListAllPCPackagingRateAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

func ListPCPackagingRateAssetbyPCID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var pcpackagingitr shim.StateQueryIteratorInterface
	var PCPackagingRateAssetList []PCPackagingRateAsset

	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("ListAllPCPackagingRateAsset :This User is not allowed to List ParticipantRateAsset"))
	}
	user := strings.ToUpper(value)
	if "ADMIN" != user && user != "PC" {
		return shim.Error(fmt.Sprintf("ListAllPCPackagingRateAsset :This User is not allowed to List ParticipantRateAsset"))
	}

	if len(args) < 1 {
		return shim.Error("ListPCPackagingRateAssetbyPCID :Incorrect number of arguments. Expecting StorageId")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PCID\":\"%s\"}}", PCPACKAGINGRATE_ASSET_RECORD_TYPE, args[0])
	logger.Infof("ListPCPackagingRateAssetbyPCID Query string is %s ", queryString)

	pcpackagingitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListPCPackagingRateAssetbyPCID : Instence not found in ledger")
		return shim.Error("pcpackagingitr : Instence not found in ledger")

	}
	defer pcpackagingitr.Close()
	for pcpackagingitr.HasNext() {
		data, derr := pcpackagingitr.Next()
		if derr != nil {
			logger.Errorf("ListPCPackagingRateAssetbyPCID : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListPCPackagingRateAssetbyPCID: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		pcpackaging, _ := JsontoPCPackagingRateAsset([]byte(databyte))
		PCPackagingRateAssetList = append(PCPackagingRateAssetList, pcpackaging)
	}
	Avalbytes, err = json.Marshal(PCPackagingRateAssetList)
	logger.Infof("ListPCPackagingRateAssetbyPCID Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListPCPackagingRateAssetbyPCID : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListPCPackagingRateAssetbyPCID: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

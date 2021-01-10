/*This file have  all API for main Loyalty application that is instantiated from main.go */
package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//FuncTemplate : This Function is Tamplate for all function to Blockchain
type FuncTemplate func(stub shim.ChaincodeStubInterface, args []string) pb.Response

//Dfarmsc : This  structure will store function pointers  for  all function executed by this chaincode
type Dfarmsc struct {
	funcMap      map[string]FuncTemplate
	restartcheck bool
}

//Result : Status info
type Result struct {
	Status string
}

// SmConfigAsset to store organization related Info
type SmConfigAsset struct {
	DocType string   `json:"docType,omitempty"`
	ID      string   `json:"ID,omitempty"`
	Orglist []string `json:"orgList,omitempty"`
}

//Constent definition for  all function names
const (
	//This function is in this file
	SC string = "StatusCheck"
	//Functions spefic to different assets
	CGAS   string = "CreateGenAssets"
	LAS    string = "ListGenAssets"
	PR     string = "ProduceRegistration"
	GP     string = "GetProduce"
	GAP    string = "GetAllProduce"
	UP     string = "UpdateProduce"
	CFPR   string = "CreateFarmerPriceRateAsset"
	UFPR   string = "UpdateFarmerPriceRateAsset"
	QFPR   string = "QueryFarmerPriceRateAsset"
	LFPRBP string = "ListFarmerPriceRateAssetbyProduceID"
	CPMRA  string = "CreateProduceMVPRateAsset"
	QPMRA  string = "QueryProduceMVPRateAsset"
	UPMRA  string = "UpdateProduceMVPRateAsset"
	//Participant Rate start
	CCRA  string = "CreateCASRateAsset"
	QCRA  string = "QueryCASRateAsset"
	LCRA  string = "ListCASRateAsset"
	CDERA string = "CreateDeliveryRatAsset"
	LDERA string = "ListAllDeliveryRatAsset"
	CDRA  string = "CreateDCRateAsset"
	LDRA  string = "ListDCRateAsset"
	CFRA  string = "CreatedFarmRateAsset"
	LFRA  string = "ListAlldfarmRateAsset"
	//Participant Rate end

	UCAS string = "UpdateCASAssetAvailableforSaleQty"
	LACA string = "ListAllCASAsset"
	//Payment Asset Functions
	CPIA string = "CreateParticipantPaymentInvoiceAsset"
	QPIA string = "QueryParticipantPaymentInvoiceAsset"
	UPIA string = "UpdateParticipantPaymentInvoiceAsset"
	LPIA string = "ListParticipantPaymentInvoiceAsset"

	CPCRA string = "CreatePCPackagingRateAsset"
	QPCRA string = "QueryPCPackagingRateAsset"
	//CCOrder Functions
	CCCO   string = "CreateCCOrderAsset"
	LCCOBS string = "ListCCOrderbyStorageId"
	QCCO   string = "QueryCCOrderAssetbyOrderID"
	//CAS Order Functions
	CCO   string = "CreateCSOrderAsset"
	LCOBS string = "ListCSOrderbyStorageId"
	QCAO  string = "QueryCSOrderAssetbyOrderID"

	//PC Order Functions
	CPA   string = "CreatePCOrderAsset"
	LPA   string = "ListPCOrderAsset"
	LPABS string = "ListPCOrderAssetbyStorageId"
	//GPPSL  string = "GetProducePriceatStorageLevel"
	QPOAO string = "QueryPCOrderAssetbyOrderID"

	//PCDC Order Functions
	CDOA   string = "CreateDCOrderAsset"
	LDOA   string = "ListDCOrderAsset"
	LDOABD string = "ListDCOrderAssetbyDcID"
	QDOAO  string = "QueryDCOrderAssetbyOrderID"

	//Buyer Order Functions
	CBOA   string = "CreateBuyerOrderAsset"
	LBOA   string = "ListBuyerOrderAsset"
	LBOABB string = "ListBuyerOrderAssetbyBuyerID"
	//GPPBL  string = "GetProducePriceatBuyerLevel"

	QBOAO string = "QueryBuyerOrderAssetbyOrderID"
	GPIA  string = "GenerateParticipantPaymentInvoiceAsset"

	LAPPRA string = "ListAllPCPackagingRateAsset"
	LPPRAP string = "ListPCPackagingRateAssetbyPCID"

	LCPAP string = "ListCASPCOrderAssetbyPcID"

	GPT string = "GetProduceTracking"
	//Stats
	GPS string = "GetPaymentStat"
	//GPSBI string = "GetPaymentStatbyBuyerOrderID"
	//GPSUI string = "GetPaymentStatbyUserOrderID"
	//New Function for India Green Orbit Logic
	CTVRA string = "CreateTableVarietyPerRateAsset"
	UTVRA string = "UpdateTableVarietyPerRateAsset"
	QTVRA string = "QueryTableVarietyPerRateAsset"
	LTVRA string = "ListALLTableVarietyPerRateAsset"
	GPPR  string = "GetProducePriceRate"
	//TransportInfo functions
	CTIA string = "CreateTransportInfoAsset"
	LTIA string = "ListTransportInfoAsset"
	//Markup and Discount functions
	CARA  string = "CreateAgentRatePerAsset"
	LARA  string = "ListAgentRatePerAsset"
	LARID string = "ListAgentRatePerAssetbyID"

	//Whole Seller Per Function
	CWSRA string = "CreateWholeSellerRatePerAsset"
	LWSRA string = "ListWholeSellerRatePerAsset"
	LWSID string = "ListWholeSellerRatePerAssetbyID"
	UC    string = "UpdateConfig"
	//PH Center Methods
	CPHOA string = "CreatePHOrderAsset"
	QPHOA string = "QueryPHOrderAssetbyOrderID"
	LPHOA string = "ListPHOrderbyStorageId"
)

//initfunMap() : Chaincode initialization Function
//This function will create a  map that will get initialized at chaincode init
func (inv *Dfarmsc) initfunMap() {
	inv.funcMap = make(map[string]FuncTemplate)
	inv.funcMap[SC] = ChainCodeStatusCheck
	inv.funcMap[PR] = CreateProduceAsset
	inv.funcMap[GP] = QueryProduceAsset
	inv.funcMap[GAP] = ListAllProduceAsset
	inv.funcMap[UP] = UpdateProduceAsset
	inv.funcMap[CGAS] = CreateGenAssets
	inv.funcMap[LAS] = ListGenAssets
	//PriceRate Functions
	inv.funcMap[CFPR] = CreateFarmerPriceRateAsset
	inv.funcMap[UFPR] = UpdateFarmerPriceRateAsset
	inv.funcMap[QFPR] = QueryFarmerPriceRateAsset
	inv.funcMap[LFPRBP] = ListFarmerPriceRateAssetbyProduceID
	inv.funcMap[CPMRA] = CreateProduceMVPRateAsset
	inv.funcMap[QPMRA] = QueryProduceMVPRateAsset
	//Rate Function for Participant start
	inv.funcMap[CCRA] = CreateCASRateAsset
	inv.funcMap[QCRA] = QueryCASRateAsset
	inv.funcMap[LCRA] = ListCASRateAsset

	inv.funcMap[CDERA] = CreateDeliveryRatAsset
	inv.funcMap[LDERA] = ListAllDeliveryRatAsset

	inv.funcMap[CDRA] = CreateDCRateAsset
	inv.funcMap[LDRA] = ListDCRateAsset

	inv.funcMap[CFRA] = CreatedFarmRateAsset
	inv.funcMap[LFRA] = ListAlldfarmRateAsset

	////Rate Function for Participant end
	inv.funcMap[CPIA] = CreateParticipantPaymentInvoiceAsset
	inv.funcMap[QPIA] = QueryParticipantPaymentInvoiceAsset
	inv.funcMap[UPIA] = UpdateParticipantPaymentInvoiceAsset
	inv.funcMap[LPIA] = ListParticipantPaymentInvoiceAsset
	inv.funcMap[CPCRA] = CreatePCPackagingRateAsset
	inv.funcMap[QPCRA] = QueryPCPackagingRateAsset
	//CC Order Function
	inv.funcMap[CCCO] = CreateCCOrderAsset
	inv.funcMap[LCCOBS] = ListCCOrderbyStorageId
	inv.funcMap[QCCO] = QueryCCOrderAssetbyOrderID
	//FarmerCAS Order Functions
	inv.funcMap[CCO] = CreateCSOrderAsset
	inv.funcMap[LCOBS] = ListCSOrderbyStorageId
	inv.funcMap[QCAO] = QueryCSOrderAssetbyOrderID
	//CASPC Order Functions
	inv.funcMap[CPA] = CreatePCOrderAsset
	inv.funcMap[LPA] = ListPCOrderAsset
	inv.funcMap[LPABS] = ListPCOrderAssetbyStorageId
	//inv.funcMap[GPPSL] = GetProducePriceatStorageLevel
	inv.funcMap[QPOAO] = QueryPCOrderAssetbyOrderID

	//PCDC Order Functions
	inv.funcMap[CDOA] = CreateDCOrderAsset
	inv.funcMap[LDOA] = ListDCOrderAsset
	inv.funcMap[LDOABD] = ListDCOrderAssetbyDcID
	inv.funcMap[QDOAO] = QueryDCOrderAssetbyOrderID

	//Buyer Order Functions
	inv.funcMap[CBOA] = CreateBuyerOrderAsset
	inv.funcMap[LBOA] = ListBuyerOrderAsset
	inv.funcMap[LBOABB] = ListBuyerOrderAssetbyBuyerID
	//inv.funcMap[GPPBL] = GetProducePriceatBuyerLevel
	inv.funcMap[QBOAO] = QueryBuyerOrderAssetbyOrderID
	//inv.funcMap[GPPSL] = GetProducePriceatStorageLevel
	inv.funcMap[GPIA] = GenerateParticipantPaymentInvoiceAsset

	inv.funcMap[LAPPRA] = ListAllPCPackagingRateAsset
	inv.funcMap[LPPRAP] = ListPCPackagingRateAssetbyPCID
	inv.funcMap[LCPAP] = ListPCOrderAssetbyPcID
	inv.funcMap[GPT] = GetProduceTracking
	inv.funcMap[GPS] = GetPaymentStat
	//inv.funcMap[GPSBI] = GetPaymentStatbyBuyerOrderID
	//inv.funcMap[GPSUI] = GetPaymentStatbyUserOrderID
	//Green Orbit
	inv.funcMap[CTVRA] = CreateTableVarietyPerRateAsset
	inv.funcMap[UTVRA] = UpdateTableVarietyPerRateAsset
	inv.funcMap[QTVRA] = QueryTableVarietyPerRateAsset
	inv.funcMap[LTVRA] = ListALLTableVarietyPerRateAsset
	inv.funcMap[GPPR] = GetProducePriceRate
	//TransportInfo functions
	inv.funcMap[CTIA] = CreateTransportInfoAsset
	inv.funcMap[LTIA] = ListTransportInfoAsset
	//Markup and Discount functions
	inv.funcMap[CARA] = CreateAgentRatePerAsset
	inv.funcMap[LARA] = ListAgentRatePerAsset
	inv.funcMap[LARID] = ListAgentRatePerAssetbyID
	//Whole Seller Per Function
	inv.funcMap[CWSRA] = CreateWholeSellerRatePerAsset
	inv.funcMap[LWSRA] = ListWholeSellerRatePerAsset
	inv.funcMap[LWSID] = ListWholeSellerRatePerAssetbyID
	inv.funcMap[UC] = UpdateConfig
	//PH order function Addtion
	inv.funcMap[CPHOA] = CreatePHOrderAsset
	inv.funcMap[QPHOA] = QueryPHOrderAssetbyOrderID
	inv.funcMap[LPHOA] = ListPHOrderbyStorageId

}

// Init initialize data in Chaincode
func (inv *Dfarmsc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//var keys []string
	logger.Infof("Init ChaininCode dFarmsc")
	inv.initfunMap()
	//commneted after upgrade to fabric 1.3
	inv.restartcheck = true
	/*	logger.Infof("%+v", inv)
			fun, args := stub.GetFunctionAndParameters()
			logger.Infof("%v", args)
			logger.Infof("%v", fun)
		    if args == nil {
				logger.Errorf("Init : Error reading required parameter failed %s", err)
				return shim.Error(fmt.Sprintf("Init : Error reading required parameter failed %s", err))
			}
			orglist := strings.Split(args[0], ",")
			config := SmConfigAsset{DocType: "CONFIG", ID: "001", Orglist: orglist}
			keys = append(keys, "001")
			bytedata, _ := json.Marshal(config)
			err := CreateAsset(stub, "CONFIG", keys, bytedata)
			if err != nil {
				logger.Errorf("Init : Error inserting Object first time  into LedgerState %s", err)
				return shim.Error(fmt.Sprintf("Init : Config Object init failed %s", err))
			}*/
	return shim.Success(nil)
}

// Invoke : Chaincode Invoke Function
func (inv *Dfarmsc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Infof("Invoke ChaininCode Dfarmsc")
	/*if !isAllowed(stub) {
		shim.Error(fmt.Sprintf("isAllowed :This user is not allowed to access chaincode invoke func"))

	}*/
	funname, args := stub.GetFunctionAndParameters()

	if funname == "" {
		logger.Infof("Function Name is not passed correctly while invoking ChainCode")
	}

	if inv.restartcheck == false {
		inv.initfunMap()
		inv.restartcheck = true
		logger.Infof("%+v", inv)
	}
	exefun, ok := inv.funcMap[funname]
	logger.Infof("Invoke ChaininCode Dfarmsc for Function Name: %s", funname)
	if ok {
		return exefun(stub, args)
	}
	logger.Errorf("Function Name:= %s is not defined in ChaininCode", funname)
	return shim.Error(fmt.Sprintf("Invalid Function Name: %s", funname))
}

// ChainCodeStatusCheck function is called by Nodejs after intsalling and instantiating chaincode to check if it is up and running
func ChainCodeStatusCheck(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	logger.Infof("ChaininCode  Running Status Check")
	result := Result{}
	result.Status = "ChainCode Running Successfully"
	logger.Infof("ChaininCode  Running Status json data: %+v", result)
	availabeByte, _ := json.Marshal(result)
	logger.Infof("ChaininCode  Running Status json data: %v", availabeByte)
	return shim.Success(availabeByte)
}

func getConfig(stub shim.ChaincodeStubInterface) SmConfigAsset {

	var keys []string
	var bytedata []byte
	keys = append(keys, "001")
	bytedata, _ = QueryAsset(stub, "CONFIG", keys)
	config := SmConfigAsset{}
	_ = json.Unmarshal(bytedata, &config)
	return config
}

// UpdateConfig is called by Nodejs when adding new organizations
func UpdateConfig(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var keys []string
	keys = append(keys, "001")

	Avalbytes, err := QueryAsset(stub, "CONFIG", keys)
	if err != nil {
		logger.Infof("UpdateConfig :Ledger Data is  : %s ", Avalbytes)
		err = CreateAsset(stub, "CONFIG", keys, Avalbytes)
		if err != nil {
			logger.Errorf("UpdateConfig : Error inserting Object first time  into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("UpdateConfig :  Object first time create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}
	config := SmConfigAsset{}
	_ = json.Unmarshal(Avalbytes, &config)
	newOrg := strings.Split(args[0], ",")
	config.Orglist = append(config.Orglist, newOrg...)
	Avalbytes, _ = json.Marshal(config)
	err = UpdateAssetWithoutGet(stub, "CONFIG", keys, Avalbytes)
	if err != nil {
		logger.Errorf("UpdateConfig : Error inserting Object first time  into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("UpdateConfig : Config Object update failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

func isAllowed(stub shim.ChaincodeStubInterface) bool {
	value, found, _ := cid.GetAttributeValue(stub, "orgsName")
	logger.Infof("Request Organization  %s ", value)
	if !found {
		return false
	}

	config := getConfig(stub)
	logger.Infof("chaincode Organization list %v ", config)
	for i := 0; i < len(config.Orglist); i++ {
		if config.Orglist[i] == value {
			return true
		}
	}
	return false
}

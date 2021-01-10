/*This file contains all the structure  and method for Factoring Invoice Asset
*
 */
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	pb "github.com/hyperledger/fabric/protos/peer"
)

/*Farmer - 1
Cold Storage - 4
Pack House - 8
Distribution Center - 6
Buyer - 7
DFarm - 5*/

const (
	FARMERPRICERATE_ASSET_RECORD_TYPE        string = "FarmerPriceRate"
	PRODUCERATE_ASSET_RECORD_TYPE            string = "ProduceRateAsset"
	PRODUCEMVPRATE_ASSET_RECORD_TYPE         string = "ProduceMVPRateAsset"
	TABLE_VARIETY_PER_RATE_ASSET_RECORD_TYPE string = "TableVarietyPerRateAsset"
	FARMER                                   int    = 1
	CAS                                      int    = 4
	PC                                       int    = 8
	DC                                       int    = 6
	PCDCTRANS                                int    = 2
	DFARM                                    int    = 5
	AGENT                                    int    = 9
	DEMURRAGEACC                             int    = 3
	BUYER                                    int    = 7
	WHOLESELLER                              int    = 10

	/*TABLEQTY                                 float32 = 0.60 // % of table quantity in 100 item is 60
	PROCESSINGQTY                            float32 = 0.30
	WASTAGEQTY                               float32 = 0.10
	XLPER                                    float32 = 0.10 // XL % in Table variety
	LPER                                     float32 = 0.15
	MPER                                     float32 = 0.50
	SPER                                     float32 = 0.15
	XSPER                                    float32 = 0.10
	//Medium Variety will sell at price after all addition cost
	//Other Table variety Quality price will increase/decrease by % shown below
	XLPRICEINCREASEPER float32 = 0.30
	LPRICEINCREASEPER  float32 = 0.15
	SPRICEDRCREASEPER  float32 = 0.15
	XSPRICEDRCREASEPER float32 = 0.30
	//Unit Type and their value in Pound
	BUSHELLB        uint    = 40
	BINBUSHEL       uint    = 18
	BINLB           uint    = 720
	PRODUCEYELDSPER float32 = 0.75
	PRODUCECULLPER  float32 = 0.25*/
)

type Rate struct {
	QuantityUnit string  `json:"UNIT,omitempty"`     //pound,Ton,KG etc
	CurrencyUnit string  `json:"CURRENCY,omitempty"` //usd,euro,inr etc
	Value        float64 `json:"VALUE,omitempty"`
}
type FarmerPriceRateAsset struct {
	DocType          string `json:"docType,omitempty"`
	FarmerID         string `json:"FARMERID,omitempty"`
	ProduceID        string `json:"PRODUCEID,omitempty"` //primary key
	Variety          string `json:"VARIETY,omitempty"`   //primery key
	State            string `json:"STATE,omitempty"`
	Country          string `json:"COUNTRY,omitempty"`
	MVPPrice         []Rate `json:"MVPPRICE,omitempty"`
	BaseUnit         string `json:"BASE_UNIT,omitempty"`
	Selectunit       Unit   `json:"SELECTED_UNIT,omitempty"`
	IsaskingPriceset bool   `json:"ISASKINGPRICESET,omitempty"`
	AskingPrice      []Rate `json:"ASKINGPRICE,omitempty"` // farmer can set this price
}

type VarietyBasedRate struct {
	Variety string `json:"VARIETY,omitempty"`
	Rates   []Rate `json:"RATES,omitempty"`
}

type ProduceMVPRateAsset struct {
	DocType     string                      `json:"docType,omitempty"`
	ProduceName string                      `json:"PRODUCE,omitempty"`
	State       string                      `json:"STATE,omitempty"`
	Country     string                      `json:"COUNTRY,omitempty"`
	BaseUnit    string                      `json:"BASE_UNIT,omitempty"`
	SelectUnit  Unit                        `json:"SELECTED_UNIT,omitempty"`
	MVPRates    map[string]VarietyBasedRate `json:"MVPRATES,omitempty"` //Variety as a string key
}
type ProduceMVPRateResponse struct {
	ProduceName string             `json:"PRODUCE,omitempty"`
	State       string             `json:"STATE,omitempty"`
	Country     string             `json:"COUNTRY,omitempty"`
	BaseUnit    string             `json:"BASE_UNIT,omitempty"`
	SelectUnit  Unit               `json:"SELECTED_UNIT,omitempty"`
	MVPRates    []VarietyBasedRate `json:"MVPRATES,omitempty"` //Variety as a string key
}

//[] {QUALITYTYPE,MODE ,VALUE}
type QualityPerInfo struct {
	Type     string  `json:"TYPE,omitempty"`
	OpMode   bool    `json:"OPMODE,omitempty"` //0 means add/more 1 means less/subscract
	PerValue float32 `json:"PERVALUE,omitempty"`
}

//New Structure for India Price calculation for Produce

type TableVarietyPerRateAsset struct {
	DocType         string           `json:"docType,omitempty"`
	ProduceName     string           `json:"PRODUCE,omitempty"`
	State           string           `json:"STATE,omitempty"`
	Country         string           `json:"COUNTRY,omitempty"`
	Variety         string           `json:"VARIETY,omitempty"`
	QualityPerInfos []QualityPerInfo `json:"QUALITYPERINFOS,omitempty"`
}
type PriceInqueryRequest struct {
	ProduceName         string            `json:"PRODUCE,omitempty"`
	Variety             string            `json:"VARIETY,omitempty"`
	BaseUnit            string            `json:"BASE_UNIT,omitempty"`
	SelectUnit          Unit              `json:"SELECTED_UNIT,omitempty"`
	ParticipantTypeList []ParticipantInfo `json:"PARTICIPANTTYPELIST,omitempty"`
}
type PriceInqueryResponse struct {
	ProduceName       string             `json:"PRODUCE,omitempty"`
	Variety           string             `json:"VARIETY,omitempty"`
	BaseUnit          string             `json:"BASE_UNIT,omitempty"`
	SelectUnit        Unit               `json:"SELECTED_UNIT,omitempty"`
	IsPCIncluded      bool               `json:"ISPCINCLUDED,omitempty"`
	VarietyRates      []Rate             `json:"VARIETYRATES,omitempty"`
	QualityPriceRates []QualityPriceRate `json:"QUALITYPRICERATES,omitempty"` //Variety as a string key
}

//0,1-Farmer  incase array is empty or have value 0/1 we will return Farmer Marklet pirce
//2-CAS
//3-PC
//4-Transport
//5-DC
//6-Dfarm
//7-Buyer
type PriceRateKey struct {
	ProduceID   string `json:"PRODUCEID,omitempty"` //primary key
	Variety     string `json:"VARIETY,omitempty"`   //primery key
	ProduceName string `json:"PRODUCE,omitempty"`
	Country     string `json:"COUNTRY,omitempty"`
	State       string `json:"STATE,omitempty"`
}

//Function to convert MAP into Array
func RateMaptoArray(m map[string]VarietyBasedRate) []VarietyBasedRate {
	var ratelist []VarietyBasedRate
	for _, value := range m {
		ratelist = append(ratelist, value)
	}
	return ratelist
}

func ArraytoMap(m []ParticipantInfo) map[int]ParticipantInfo {
	var objmap map[int]ParticipantInfo
	objmap = make(map[int]ParticipantInfo)
	for i := 0; i < len(m); i++ {
		objmap[m[i].Type] = m[i]
	}
	return objmap
}
func ConvertPCpricearraytoMap(m []QualityPriceRate) map[string]QualityPriceRate {
	objmap := make(map[string]QualityPriceRate)
	for i := 0; i < len(m); i++ {
		objmap[m[i].QualityName] = m[i]
	}
	return objmap
}

//JsontoFarmerPriceRateAsset to convert JSON  to asset object
func JsontoFarmerPriceRateAsset(data []byte) (FarmerPriceRateAsset, error) {
	obj := FarmerPriceRateAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to FarmerPriceRateAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert FarmerPriceRateAsset object to Json Message

func FarmerPriceRateAssettoJson(obj FarmerPriceRateAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

//JsontoProduceMVPRateAsset to convert JSON  to asset object
func JsontoProduceMVPRateAsset(data []byte) (ProduceMVPRateAsset, error) {
	obj := ProduceMVPRateAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to ProduceMVPRateAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert ProduceMVPRateAsset object to Json Message

func ProduceMVPRateAssettoJson(obj ProduceMVPRateAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateFarmerPriceRateAsset Function will  insert record in ledger after receiving request from Client Application
func CreateFarmerPriceRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateFarmerPriceRateAsset : Incorrect number of arguments.")
		return shim.Error("CreateFarmerPriceRateAsset : Incorrect number of arguments.")
	}

	// Convert the arg to a FarmerPriceRateAsset Object
	logger.Infof("CreateFarmerPriceRateAsset: Arguments for ledgerapi %s : ", args[0])

	asset, err := JsontoFarmerPriceRateAsset([]byte(args[0]))
	asset.DocType = FARMERPRICERATE_ASSET_RECORD_TYPE
	logger.Infof("CreateFarmerPriceRateAsset :Produce ID is : %s ", asset.ProduceID)

	keys = append(keys, asset.ProduceID)
	keys = append(keys, asset.Variety)
	keys = append(keys, asset.State)

	logger.Infof("CreateFarmerPriceRateAsset : Inserting object with data as  %s", args[0])
	Avalbytes, _ = FarmerPriceRateAssettoJson(asset)

	err = CreateAsset(stub, FARMERPRICERATE_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateFarmerPriceRateAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateFarmerPriceRateAsset : FarmerPriceRateAsset object create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// QueryFarmerPriceRateAsset  Function will  query  record in ledger based on ID
func QueryFarmerPriceRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var keysglobal []string
	var uiRates []Rate

	if len(args) < 1 {
		logger.Errorf("QueryFarmerPriceRateAsset : Incorrect number of arguments.")
		return shim.Error("QueryFarmerPriceRateAsset : Incorrect number of arguments.")
	}
	logger.Infof("QueryFarmerPriceRateAsset :FarmerPriceRateKey is : %s ", args[0])
	requestkey := PriceRateKey{}
	err = json.Unmarshal([]byte(args[0]), &requestkey)
	keys = append(keys, requestkey.ProduceID)
	keys = append(keys, requestkey.Variety)
	keys = append(keys, requestkey.State)
	Avalbytes, err = QueryAsset(stub, FARMERPRICERATE_ASSET_RECORD_TYPE, keys)
	if err != nil {

		keysglobal = append(keysglobal, requestkey.State)       // State
		keysglobal = append(keysglobal, requestkey.ProduceName) // Produce Name
		Avalbytes, err = QueryAsset(stub, PRODUCERATE_ASSET_RECORD_TYPE, keysglobal)
		if err != nil {
			logger.Errorf("QueryFarmerPriceRateAsset : Error inserting Object into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("QueryFarmerPriceRateAsset : FarmerPriceRateAsset object get failed %s", err))
		}
		asset, _ := JsontoProduceMVPRateAsset([]byte(Avalbytes))
		mvprate := asset.MVPRates[args[1]]
		logger.Infof("CreateFarmerPriceRateAsset :mvprate is : %f ", mvprate)

		farmerrate := FarmerPriceRateAsset{}
		logger.Infof("CreateFarmerPriceRateAsset :farmerrate is : %f ", farmerrate)

		farmerrate.IsaskingPriceset = false
		uiRates := append(uiRates, mvprate.Rates[0])
		farmerrate.MVPPrice = uiRates
		logger.Infof("CreateFarmerPriceRateAsset :farmerrate.MVPPrice is : %f ", farmerrate.MVPPrice)

		farmerrate.ProduceID = args[0]
		farmerrate.Variety = args[1]
		Avalbytes, err = FarmerPriceRateAssettoJson(farmerrate)
		return shim.Success([]byte(Avalbytes))
	}

	return shim.Success([]byte(Avalbytes))

}

//UpdateFarmerPriceRateAsset Function will  update record in ledger after receiving request from Client Application
func UpdateFarmerPriceRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("UpdateFarmerPriceRateAsset : Incorrect number of arguments.")
		return shim.Error("UpdateFarmerPriceRateAsset : Incorrect number of arguments.")
	}

	// Convert the arg to a FarmerPriceRateAsset Object
	logger.Infof("UpdateFarmerPriceRateAsset: Arguments for ledgerapi %s : ", args[0])

	asset, err := JsontoFarmerPriceRateAsset([]byte(args[0]))

	logger.Infof("UpdateFarmerPriceRateAsset :Produce ID is : %s ", asset.ProduceID)
	asset.DocType = FARMERPRICERATE_ASSET_RECORD_TYPE
	keys = append(keys, asset.ProduceID)
	keys = append(keys, asset.Variety)
	keys = append(keys, asset.State)

	logger.Infof("UpdateFarmerPriceRateAsset : updating object with data as  %s", args[0])
	Avalbytes, _ = FarmerPriceRateAssettoJson(asset)

	err = UpdateAsset(stub, FARMERPRICERATE_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("UpdateFarmerPriceRateAsset : Error updating Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("UpdateFarmerPriceRateAsset : FarmerPriceRateAsset object update failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// ListFarmerPriceRateAsset  Function will  query  all record from DB
func ListFarmerPriceRateAssetbyProduceID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var pricerateitr shim.StateQueryIteratorInterface
	var priceratelist []FarmerPriceRateAsset
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PRODUCEID\":\"%s\"}}", FARMERPRICERATE_ASSET_RECORD_TYPE, args[0])
	logger.Infof("ListFarmerPriceRateAssetbyProduceID Query string is %s ", queryString)
	pricerateitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("ListFarmerPriceRateAssetbyProduceID : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer pricerateitr.Close()
	for pricerateitr.HasNext() {
		data, derr := pricerateitr.Next()
		if derr != nil {
			logger.Errorf("ListFarmerPriceRateAssetbyProduceID : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListFarmerPriceRateAssetbyProduceID: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		pricerate, _ := JsontoFarmerPriceRateAsset([]byte(databyte))
		priceratelist = append(priceratelist, pricerate)
	}
	Avalbytes, err = json.Marshal(priceratelist)
	logger.Infof("ListFarmerPriceRateAssetbyProduceID Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListFarmerPriceRateAssetbyProduceID : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListFarmerPriceRateAssetbyProduceID: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// CreateProduceMVPRateAsset Function will  insert record in ledger after receiving request from Client Application
func CreateProduceMVPRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateProduceMVPRateAsset : Incorrect number of arguments.")
		return shim.Error("CreateProduceMVPRateAsset : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("CreateProduceMVPRateAsset :Attribute approle not found  to create ProduceMVPRateAsset"))
	}
	if "ADMIN" != strings.ToUpper(value) {
		return shim.Error(fmt.Sprintf("CreateProduceMVPRateAsset :This User   is not ADMIN ,so  not allowed to create ProduceMVPRateAsset"))
	}
	// Convert the arg to a ProduceMVPRateAsset Object
	logger.Infof("CreateProduceMVPRateAsset: Arguments for ledgerapi %s : ", args[0])

	asset, err := JsontoProduceMVPRateAsset([]byte(args[0]))
	asset.DocType=PRODUCEMVPRATE_ASSET_RECORD_TYPE

	logger.Infof("CreateProduceMVPRateAsset :Produce Name is : %s ", asset.ProduceName)
	keys = append(keys, asset.ProduceName)
	keys = append(keys, asset.Country)
	keys = append(keys, asset.State)

	Avalbytes, err = QueryAsset(stub, PRODUCEMVPRATE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Infof("CreateProduceMVPRateAsset : Inserting object with data as  %s", args[0])
		Avalbytes,_= ProduceMVPRateAssettoJson(asset)

		err = CreateAsset(stub, PRODUCEMVPRATE_ASSET_RECORD_TYPE, keys, Avalbytes)
		if err != nil {
			logger.Errorf("CreateProduceMVPRateAsset : Error inserting Object into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("CreateProduceMVPRateAsset : ProduceMVPRateAsset object create failed %s", err))
		}
		return shim.Success([]byte(Avalbytes))
	}

	assetread, _ := JsontoProduceMVPRateAsset([]byte(Avalbytes))
	logger.Infof("CreateProduceMVPRateAsset : Got read object with data as  %d", assetread)
	for key, value := range asset.MVPRates {
		assetread.MVPRates[key] = value
	}
	logger.Infof("CreateProduceMVPRateAsset : Got updated object with data as  %d", assetread)
	Avalbytes, _ = ProduceMVPRateAssettoJson(assetread)

	err = UpdateAssetWithoutGet(stub, PRODUCEMVPRATE_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("UpdateProduceMVPRateAsset : Error updating Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("UpdateProduceMVPRateAsset : ProduceMVPRateAsset object update failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

// QueryProduceMVPRateAsset  Function will  query  record in ledger based on ID
func QueryProduceMVPRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var producerateitr shim.StateQueryIteratorInterface
	var producerateList []ProduceMVPRateResponse
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("QueryProduceMVPRateAsset :Attribute approle not found to  Query ProduceMVPRateAsset"))
	}
	user := strings.ToUpper(value)
	if user != "ADMIN" && user != "FARMER" {
		return shim.Error(fmt.Sprintf("QueryProduceMVPRateAsset : User is not allowed to Query ProduceMVPRateAsset"))
	}

	if len(args) == 3 {
		logger.Infof("QueryProduceMVPRateAsset :Produce Name  is : %s ", args[0])
		logger.Infof("QueryProduceMVPRateAsset :Country is : %s ", args[1])
		logger.Infof("QueryProduceMVPRateAsset :State  is : %s ", args[2])
		keys = append(keys, args[0]) //Producename
		keys = append(keys, args[1]) //Country
		keys = append(keys, args[2]) //State
		Avalbytes, err = QueryAsset(stub, PRODUCEMVPRATE_ASSET_RECORD_TYPE, keys)
		if err != nil {
			logger.Errorf("QueryProduceMVPRateAsset : Error inserting Object into LedgerState %s", err)
			return shim.Error(fmt.Sprintf("QueryProduceMVPRateAsset : ProduceMVPRateAsset object get failed %s", err))
		}
		producerate, _ := JsontoProduceMVPRateAsset([]byte(Avalbytes))
		producerateres := ProduceMVPRateResponse{}
		producerateres.ProduceName = producerate.ProduceName
		producerateres.Country = producerate.Country
		producerateres.State = producerate.State
		producerateres.SelectUnit = producerate.SelectUnit
		producerateres.BaseUnit = producerate.BaseUnit
		producerateres.MVPRates = RateMaptoArray(producerate.MVPRates)
		Avalbytes, _ = json.Marshal(producerateres)
	} else {

		producerateitr, err = ListAllAsset(stub, PRODUCEMVPRATE_ASSET_RECORD_TYPE, keys)
		if err != nil {
			logger.Errorf("ListProduceMVPRateAsset : Instence not found in ledger")
			return shim.Error("producerateitr : Instence not found in ledger")

		}
		defer producerateitr.Close()
		for producerateitr.HasNext() {
			data, derr := producerateitr.Next()
			if derr != nil {
				logger.Errorf("ListProduceMVPRateAsset : Cannot parse result set. Error : %v", derr)
				return shim.Error(fmt.Sprintf("ListProduceMVPRateAsset: Cannot parse result set. Error : %v", derr))

			}
			databyte := data.GetValue()

			producerate, _ := JsontoProduceMVPRateAsset([]byte(databyte))
			producerateres := ProduceMVPRateResponse{}
			producerateres.ProduceName = producerate.ProduceName
			producerateres.Country = producerate.Country
			producerateres.State = producerate.State
			producerateres.SelectUnit = producerate.SelectUnit
			producerateres.BaseUnit = producerate.BaseUnit
			producerateres.MVPRates = RateMaptoArray(producerate.MVPRates)
			producerateList = append(producerateList, producerateres)
		}

		Avalbytes, err = json.Marshal(producerateList)
		logger.Infof("ListProduceMVPRateAsset Responce for App : %v", Avalbytes)
		if err != nil {
			logger.Errorf("ListProduceMVPRateAsset : Cannot Marshal result set. Error : %v", err)
			return shim.Error(fmt.Sprintf("ListProduceMVPRateAsset: Cannot Marshal result set. Error : %v", err))
		}
	}
	return shim.Success([]byte(Avalbytes))
}

//getFarmerMarketRate with Product Name,Country,State,Varity
func getFarmerMarketRate(stub shim.ChaincodeStubInterface, args []string) (Rate, bool) {
	var err error
	var Avalbytes []byte
	var keys []string
	/*value, found, _ := cid.GetAttributeValue(stub, "approle")
	if !found {
		return shim.Error(fmt.Sprintf("QueryProduceMVPRateAsset :Attribute approle not found to  Query ProduceMVPRateAsset"))
	}
	user := strings.ToUpper(value)
	if user != "ADMIN" && user != "FARMER" {
		return shim.Error(fmt.Sprintf("QueryProduceMVPRateAsset : User is not allowed to Query ProduceMVPRateAsset"))
	}*/

	if len(args) == 4 {
		logger.Infof("getFarmerMarketRate :Produce Name  is : %s ", args[0])
		logger.Infof("getFarmerMarketRate :Country is : %s ", args[1])
		logger.Infof("getFarmerMarketRate :State  is : %s ", args[2])
		logger.Infof("getFarmerMarketRate :Variety  is : %s ", args[3])
		keys = append(keys, args[0]) //Producename
		keys = append(keys, args[1]) //Country
		keys = append(keys, args[2]) //State
		Avalbytes, err = QueryAsset(stub, PRODUCEMVPRATE_ASSET_RECORD_TYPE, keys)
		if err != nil {
			logger.Errorf("getFarmerMarketRate : Error getting Object into LedgerState %s", err)
			return Rate{}, false
		}
		producerate, _ := JsontoProduceMVPRateAsset([]byte(Avalbytes))
		vrates, found := producerate.MVPRates[args[3]]
		if !found {
			logger.Errorf("getFarmerMarketRate : Error getting Produce Market Price for required Variety %s", args[3])
			return Rate{}, false
		}
		rate := vrates.Rates[0]
		return rate, true
	}
	logger.Errorf("getFarmerMarketRate : No of argument is not correct %d", len(args))
	return Rate{}, false

}

//Handeling Percentage related information based on Market Price

//JsontoTableVarietyPerRateAsset to convert JSON  to asset object
func JsontoTableVarietyPerRateAsset(data []byte) (TableVarietyPerRateAsset, error) {
	obj := TableVarietyPerRateAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to TableVarietyPerRateAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert TableVarietyPerRateAsset object to Json Message

func TableVarietyPerRateAssettoJson(obj TableVarietyPerRateAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

//create TableVarietyPerRateAsset
// CreateTableVarietyPerRateAsset Function will  insert record in ledger after receiving request from Client Application
func CreateTableVarietyPerRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateTableVarietyPerRateAsset : Incorrect number of arguments.")
		return shim.Error("CreateTableVarietyPerRateAsset : Incorrect number of arguments.")
	}

	// Convert the arg to a CreateTableVarietyPerRateAsset Object
	logger.Infof("CreateTableVarietyPerRateAsset: Arguments for ledgerapi %s : ", args[0])

	asset, err := JsontoTableVarietyPerRateAsset([]byte(args[0]))
	asset.DocType = TABLE_VARIETY_PER_RATE_ASSET_RECORD_TYPE
	keys = append(keys, asset.Country)
	keys = append(keys, asset.State)
	keys = append(keys, asset.ProduceName)
	keys = append(keys, asset.Variety)
	logger.Infof("CreateTableVarietyPerRateAsset : Inserting object with data as  %s", args[0])
	Avalbytes, _ = TableVarietyPerRateAssettoJson(asset)

	err = CreateAsset(stub, TABLE_VARIETY_PER_RATE_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateTableVarietyPerRateAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateTableVarietyPerRateAsset : TableVarietyPerRateAsset object create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

//UpdateTableVarietyPerRateAsset Function will  update record in ledger after receiving request from Client Application
func UpdateTableVarietyPerRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("UpdateTableVarietyPerRateAsset : Incorrect number of arguments.")
		return shim.Error("UpdateTableVarietyPerRateAsset : Incorrect number of arguments.")
	}

	// Convert the arg to a TableVarietyPerRateAsset Object
	logger.Infof("UpdateTableVarietyPerRateAsset: Arguments for ledgerapi %s : ", args[0])

	asset, err := JsontoTableVarietyPerRateAsset([]byte(args[0]))

	asset.DocType = TABLE_VARIETY_PER_RATE_ASSET_RECORD_TYPE
	keys = append(keys, asset.Country)
	keys = append(keys, asset.State)
	keys = append(keys, asset.ProduceName)
	keys = append(keys, asset.Variety)

	logger.Infof("UpdateTableVarietyPerRateAsset : updating object with data as  %s", args[0])
	Avalbytes, _ = TableVarietyPerRateAssettoJson(asset)

	err = UpdateAsset(stub, TABLE_VARIETY_PER_RATE_ASSET_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("UpdateTableVarietyPerRateAsset : Error updating Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("UpdateTableVarietyPerRateAsset : TableVarietyPerRateAsset object update failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// QueryTableVarietyPerRateAsset  Function will  query  record in ledger based on ID
func QueryTableVarietyPerRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("QueryTableVarietyPerRateAsset : Incorrect number of arguments.")
		return shim.Error("QueryTableVarietyPerRateAsset : Incorrect number of arguments.")
	}
	logger.Infof("QueryTableVarietyPerRateAsset :FarmerPriceRateKey is : %s ", args[0])
	requestkey := PriceRateKey{}
	err = json.Unmarshal([]byte(args[0]), &requestkey)

	keys = append(keys, requestkey.Country)
	keys = append(keys, requestkey.State)
	keys = append(keys, requestkey.ProduceName)
	keys = append(keys, requestkey.Variety)
	logger.Infof("QueryTableVarietyPerRateAsset :FarmerPriceRateKey is : %s ", keys)
	Avalbytes, err = QueryAsset(stub, TABLE_VARIETY_PER_RATE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryTableVarietyPerRateAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryTableVarietyPerRateAsset : TableVarietyPerRateAsset object get failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))

}

// getTableVarietyPerRateAsset  Function will  query  record in ledger based on ID
func getTableVarietyPerRateAsset(stub shim.ChaincodeStubInterface, args []string) (TableVarietyPerRateAsset, bool) {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 4 {
		logger.Errorf("getTableVarietyPerRateAsset : Incorrect number of arguments.")
		return TableVarietyPerRateAsset{}, false
	}
	logger.Infof("getTableVarietyPerRateAsset :FarmerPriceRateKey is : %s ", args[0])
	requestkey := PriceRateKey{}
	err = json.Unmarshal([]byte(args[0]), &requestkey)

	keys = append(keys, args[0])
	keys = append(keys, args[1])
	keys = append(keys, args[2])
	keys = append(keys, args[3])
	logger.Infof("getTableVarietyPerRateAsset :FarmerPriceRateKey is : %s ", keys)
	Avalbytes, err = QueryAsset(stub, TABLE_VARIETY_PER_RATE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("getTableVarietyPerRateAsset : Error inserting Object into LedgerState %s", err)
		return TableVarietyPerRateAsset{}, false
	}
	perrate, _ := JsontoTableVarietyPerRateAsset([]byte(Avalbytes))
	logger.Infof("getTableVarietyPerRateAsset : successfull query Data is %v", perrate)
	return perrate, true
}

// ListALLTableVarietyPerRateAsset  Function will  query  all record from DB
func ListALLTableVarietyPerRateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var priceperrateitr shim.StateQueryIteratorInterface
	var priceperratelist []TableVarietyPerRateAsset
	priceperrateitr, err = ListAllAsset(stub, TABLE_VARIETY_PER_RATE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListALLTableVarietyPerRateAsset : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer priceperrateitr.Close()
	for priceperrateitr.HasNext() {
		data, derr := priceperrateitr.Next()
		if derr != nil {
			logger.Errorf("ListALLTableVarietyPerRateAsset : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListALLTableVarietyPerRateAsset: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		pricerate, _ := JsontoTableVarietyPerRateAsset([]byte(databyte))
		priceperratelist = append(priceperratelist, pricerate)
	}
	Avalbytes, err = json.Marshal(priceperratelist)
	logger.Infof("ListALLTableVarietyPerRateAsset Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListALLTableVarietyPerRateAsset : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListALLTableVarietyPerRateAsset: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

//GetPriceRate
/*---------------------------
This function will first read market price
then see list of required participant
and then return price based on calculated data
type PriceInqueryResponse struct {
	ProduceName         string             `json:"PRODUCE,omitempty"`
	Variety             string             `json:"VARIETY,omitempty"`
	ParticipantTypeList []int              `json:"PARTICIPANTTYPELIST,omitempty"`
	BaseUnit            string             `json:"BASE_UNIT,omitempty"`
	SelectUnit          Unit               `json:"SELECTED_UNIT,omitempty"`
	IsPCIncluded        bool               `json:"ISPCINCLUDED,omitempty"`
	VarietyRates        []Rate             `json:"VARIETYRATES,omitempty"`
	QualityPriceRates   []QualityPriceRate `json:"QUALITYPRICERATES,omitempty"` //Variety as a string key
}
----------------*/
func GetProducePriceRate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys, cskeys, akeys []string
	var pckeys []string
	var perkeys []string
	var participantRate ParticipantRate
	var casrate, farmerrate = 0.00, 0.00
	var dcrate, dfarmrate = 0.00, 0.00
	var ok bool
	var transportrate = 0.00

	if len(args) < 1 {
		logger.Errorf("GetProducePriceRate : Incorrect number of arguments.")
		return shim.Error("GetProducePriceRate : Incorrect number of arguments.")
	}
	value, found, _ := cid.GetAttributeValue(stub, "approle")
	logger.Infof("GetProducePriceRate: User Role: %s  ", value)
	if !found {
		return shim.Error(fmt.Sprintf("GetProducePriceRate :Attribute approle not found to GetProducePriceRate"))
	}

	if "ADMIN" != strings.ToUpper(value) && "BUYER" != strings.ToUpper(value) {
		return shim.Error(fmt.Sprintf("GetProducePriceRate :This User is not allowed to  GetProducePriceRate"))
	}

	// Convert the arg to a FarmerPriceRateAsset Object

	logger.Infof("GetProducePriceRate: Arguments for ledgerapi %s : ", args[0])
	asset := PriceInqueryRequest{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	logger.Infof("GetProducePriceRate :Request is  : %v ", asset)
	listmap := ArraytoMap(asset.ParticipantTypeList)
	var farmeruserinfo ParticipantInfo
	var farmerfound bool
	if farmeruserinfo, farmerfound = listmap[FARMER]; !farmerfound {
		return shim.Error(fmt.Sprintf("GetProducePriceRate :Farmer is not in user list"))

	}

	//get Mareket Price and create Response  msg
	keys = append(keys, asset.ProduceName)
	keys = append(keys, farmeruserinfo.Country)
	keys = append(keys, farmeruserinfo.State)
	Avalbytes, err = QueryAsset(stub, PRODUCERATE_ASSET_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("GetProducePriceRate : Error getting Produce Market Price %s", err)
		return shim.Error(fmt.Sprintf("GetProducePriceRate : ProduceRateAsset object query failed %s", err))
	}

	producerate, _ := JsontoProduceMVPRateAsset([]byte(Avalbytes))
	vrates, found := producerate.MVPRates[asset.Variety]
	if !found {
		logger.Errorf("GetProducePriceRate : Error getting Produce Market Price for required Variety %s", asset.Variety)
		return shim.Error(fmt.Sprintf("GetProducePriceRate : Error getting Produce Market Price for required Variety %s", asset.Variety))
	}
	farmerrate = vrates.Rates[0].Value
	logger.Infof("GetProducePriceRate:Market Price for BaseUnit %s is %f", farmerrate)
	response := PriceInqueryResponse{}
	response.BaseUnit = asset.BaseUnit
	response.SelectUnit = asset.SelectUnit
	response.ProduceName = asset.ProduceName
	response.IsPCIncluded = false
	response.Variety = asset.Variety

	if len(asset.ParticipantTypeList) < 2 {
		response.VarietyRates = append(response.VarietyRates, vrates.Rates[0])
		logger.Infof("GetProducePriceRate:Rate Going to return in response is %v", response.VarietyRates)
		responsedata, _ := json.Marshal(response)
		return shim.Success([]byte(responsedata))

	} //Handle Farmer and other entiry price logic

	participantRate, ok = getParticipantRateAsset(stub, asset.ProduceName, listmap)
	if !ok {
		logger.Errorf("GetProducePriceRate : Error Querying Object from LedgerState %s", err)
		return shim.Error(fmt.Sprintf("GetProducePriceRate : ParticipantRateAsset object get failed %s", err))
	}

	//check if CS in the list:

	if csinfo, csfound := listmap[CAS]; csfound {
		cskeys = append(cskeys, csinfo.ID)
		Avalbytes, err = QueryAsset(stub, CAS_ASSET_RECORD_TYPE, cskeys)
		if err != nil {
			logger.Errorf("GetProducePriceRate : Error Querying Object from LedgerState %s", err)
			return shim.Error(fmt.Sprintf("GetProducePriceRate : CASRateAsset object get failed %s", err))
		}
		csRate, _ := JsontoCASRateAsset(Avalbytes)
		casrate = csRate.PriceRate[0].Value

	}
	_, transfound := listmap[PCDCTRANS]
	_, dcfound := listmap[DC]
	_, dfarmfound := listmap[DFARM]
	_, agentfound := listmap[AGENT]
	_, demurfound := listmap[DEMURRAGEACC]
	_, buyerfound := listmap[BUYER]
	_, wsfound := listmap[WHOLESELLER]

	var agentperasset AgentRatePerAsset
	var wsperasset WholeSellerRatePerAsset
	akeys = append(akeys, listmap[BUYER].ID) //Destination id i.e Buyer ID
	akeys = append(akeys, asset.ProduceName) //Produce Name
	akeys = append(akeys, asset.Variety)     //Variety

	if agentfound && buyerfound {

		agentperasset, _ = getAgentRatePerAsset(stub, akeys)
	}
	if wsfound && buyerfound {
		wsperasset, _ = getWholeSellerRatePerAsset(stub, akeys)
	}

	// Processing Center exist
	if pcinfo, pcfound := listmap[PC]; pcfound {
		response.IsPCIncluded = true
		pckeys = append(pckeys, pcinfo.Country)    //Country
		pckeys = append(pckeys, pcinfo.State)      //State
		pckeys = append(pckeys, asset.ProduceName) //Produce Name
		pckeys = append(pckeys, asset.Variety)     //Variety
		pckeys = append(pckeys, pcinfo.ID)         //pcID
		packagingrate, ok := getPCPackagingRateAsset(stub, pckeys)
		if !ok {
			logger.Errorf("GetProducePriceRate : Error Querying Object from LedgerState %t", ok)
			return shim.Error(fmt.Sprintf("GetProducePriceRate : PCPackagingRateAsset object get failed %t", ok))
		}
		packagingrateMap := ConvertPCpricearraytoMap(packagingrate.TableVerityPrice)

		perkeys = append(perkeys, pcinfo.Country)
		perkeys = append(perkeys, pcinfo.State)
		perkeys = append(perkeys, asset.ProduceName)
		perkeys = append(perkeys, asset.Variety)
		perrate, ok := getTableVarietyPerRateAsset(stub, perkeys)
		if !ok {
			logger.Errorf("GetProducePriceRate : Error inserting Object into LedgerState %t", ok)
			return shim.Error(fmt.Sprintf("GetProducePriceRate : TableVarietyPerRateAsset object get failed %t", ok))
		}

		if transfound && dcfound {

			distance, _ := getTransportInfoAsset(stub, listmap[PC].ID, listmap[DC].ID)
			transportrate = float64(distance.Distance.Value) * participantRate.DeliveryRateVal.Value
			logger.Infof("GetProducePriceRate:Rate total Transport Rate : %f ", transportrate)
		} else {
			logger.Infof("GetProducePriceRate:Not able to get transport rate as DC is not in the list")
		}

		for i := 0; i < len(perrate.QualityPerInfos); i++ {
			qinfo := perrate.QualityPerInfos[i]
			var qualitytmpRate = 0.00
			var agentRate = 0.00
			var wsRate = 0.0
			var demurRate = 0.0

			if value, ok := packagingrateMap[qinfo.Type]; ok {
				if qinfo.OpMode {
					qualitytmpRate = farmerrate - farmerrate*float64(qinfo.PerValue)
					logger.Infof("GetProducePriceRate:Rate -ve % with final rate %f  and % are %f", qualitytmpRate, farmerrate*float64(qinfo.PerValue))

				} else { //if QualityPerInfos[i].OpMode else
					qualitytmpRate = farmerrate + farmerrate*float64(qinfo.PerValue)
					logger.Infof("GetProducePriceRate:Rate +ve % with final rate %f  and % are %f", qualitytmpRate, farmerrate*float64(qinfo.PerValue))

				}
				qualityrate := qualitytmpRate + value.Rates[0].Value
				logger.Infof("GetProducePriceRate:Rate total quality rate is : %f for type %s", qualityrate, qinfo.Type)

				if dcfound {
					dcrate = (qualityrate + casrate) * float64(participantRate.DCpercentage)
					logger.Infof("GetProducePriceRate:Rate dcrate %f", dcrate)

				}
				if wsfound && buyerfound {
					wsRate = (qualityrate + casrate + dcrate) * wsperasset.PerValue
					logger.Infof("GetProducePriceRate:wsRate  : %f ", wsRate)
				} else {
					logger.Infof("GetProducePriceRate:Not able to get ws Rate   as Buyer and agent  is not in the list")

				}

				if agentfound && buyerfound {
					agentRate = (qualityrate + casrate + dcrate) * agentperasset.PerValue
					logger.Infof("GetProducePriceRate:agentRate  : %f ", agentRate)
				} else {
					logger.Infof("GetProducePriceRate:Not able to get Agent  rate as Buyer and agent  is not in the list")

				}
				if dfarmfound {
					dfarmrate = (qualityrate + casrate + transportrate + dcrate + agentRate) * participantRate.Dfarmpercentage.Totalper
					logger.Infof("GetProducePriceRate:Rate dfarm Rate  : %f ", dfarmrate)

				}

				if demurfound {
					demurRate = (qualityrate + casrate + transportrate + dcrate + wsRate + agentRate + dfarmrate) * participantRate.Dfarmpercentage.DemurragePer
					logger.Infof("GetProducePriceRate:demurRate  : %f ", demurRate)
				}
				logger.Infof("rate Segment quality rate : %f,trnasportRate : %f ,casRate: %f ,DCrate %f ,DfarmRate: %f ", qualityrate, transportrate, casrate, dcrate, dfarmrate)

				finalrate := qualityrate + dcrate + dfarmrate + casrate + transportrate + agentRate + demurRate + wsRate
				finalrate = math.Round(finalrate*100) / 100
				rate := Rate{value.Rates[0].QuantityUnit, value.Rates[0].CurrencyUnit, finalrate}
				var rates []Rate
				rates = append(rates, rate)
				qrate := QualityPriceRate{value.QualityName, rates}
				logger.Infof("GetProducePriceRate:Rate quality Name %s with Rates %v", value.QualityName, rates)

				response.QualityPriceRates = append(response.QualityPriceRates, qrate)
			}

		}
		logger.Infof("GetProducePriceRate:Rate Going to return in response is %v", response.QualityPriceRates)
		responsedata, _ := json.Marshal(response)
		return shim.Success([]byte(responsedata))

	}
	var nopagentRate = 0.00
	var nopwsRate = 0.00
	var nopdemurRate = 0.0
	if dcfound {
		dcrate = (farmerrate + casrate) * float64(participantRate.DCpercentage)
	}
	if agentfound && buyerfound {
		nopagentRate = (farmerrate + casrate + dcrate) * agentperasset.PerValue
		logger.Infof("GetProducePriceRate:agentRate  : %f ", nopagentRate)
	} else {
		logger.Infof("GetProducePriceRate:Not able to get Agent  rate as Buyer and agent  is not in the list")

	}
	if wsfound && buyerfound {
		nopwsRate = (farmerrate + casrate + dcrate) * wsperasset.PerValue
		logger.Infof("GetProducePriceRate:wsRate  : %f ", nopwsRate)
	} else {
		logger.Infof("GetProducePriceRate:Not able to get WS  rate as Buyer and agent  is not in the list")

	}
	if dfarmfound {
		dfarmrate = (transportrate + farmerrate + casrate + dcrate + nopagentRate) * participantRate.Dfarmpercentage.Totalper
		logger.Infof("dFram Rate is : %f", dfarmrate)

	}
	if demurfound {
		nopdemurRate = (farmerrate + casrate + transportrate + dcrate + dfarmrate + nopagentRate + nopwsRate) * participantRate.Dfarmpercentage.DemurragePer
		logger.Infof("GetProducePriceRate:demurRate  : %f ", nopdemurRate)
	}

	logger.Infof("rate Segment Farmer rate  : %f,trnasportRate : %f ,casRate: %f ,DCrate %f ,DfarmRate: %f ", farmerrate, transportrate, casrate, dcrate, dfarmrate)
	totalRate := farmerrate + casrate + dcrate + transportrate + dfarmrate + nopagentRate + nopwsRate + nopdemurRate
	totalRate = math.Round(totalRate*100) / 100
	rate := Rate{vrates.Rates[0].QuantityUnit, vrates.Rates[0].CurrencyUnit, totalRate}
	var rates []Rate
	rates = append(rates, rate)
	response.VarietyRates = rates
	logger.Infof("GetProducePriceRate:Rate Going to return in response is %v", response.VarietyRates)
	responsedata, _ := json.Marshal(response)
	return shim.Success([]byte(responsedata))
}

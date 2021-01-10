package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	PAYMENT_INVOICE_RECORD_TYPE string = "ParticipantInvoicePaymentAsset"
)

type QualityPaymentInfo struct {
	QualityType string  `json:"QUALITYTYPE,omitempty"` //Grade A /Grade B
	Qty         uint64  `json:"Qty,omitempty"`
	Amount      float64 `json:"AMOUNT,omitempty"`
	Currency    string  `json:"CURRENCY,omitempty"`
}
type ParticipantInfo struct {
	Type    int    `json:"ROLE,omitempty"`
	ID      string `json:"MAINID,omitempty"`
	OrderID string `json:"ORDERID,omitempty"`
	Country string `json:"COUNTRY,omitempty"`
	State   string `json:"STATE,omitempty"`
}

type ParticipantPaymentInvoiceAsset struct {
	DocType                string               `json:"docType,omitempty"`
	PaymentInvoiceID       string               `json:"PAYMENTINVOICEID,omitempty"`
	OrderID                string               `json:"ORDERID,omitempty"`
	ProduceID              string               `json:"PRODUCEID,omitempty"`
	Variety                string               `json:"VARIETY,omitempty"`
	Qty                    uint64               `json:"QTY,omitempty"`
	BuyerID                string               `json:"buyerID,omitempty"`
	BuyerOrderID           string               `json:"BUYERORDERID,omitempty"`
	ID                     string               `json:"ID,omitempty"`
	Type                   int                  `json:"TYPE,omitempty"`
	Currency               string               `json:"CURRENCY,omitempty"`
	PayableAmount          float64              `json:"PAYABLEAMOUNT,omitempty"`
	QualityPaymentInfos    []QualityPaymentInfo `json:"QPINFOS,omitempty"`
	AccountPaymentResponse string               `json:"ACCOUNTPAYMENTRESPONSE,omitempty"`
	Status                 string               `json:"STATUS,omitempty"`
}

//Convert JSON  object to ParticipantPaymentInvoiceAsset
func JsontoParticipantPaymentInvoiceAsset(data []byte) (ParticipantPaymentInvoiceAsset, error) {
	obj := ParticipantPaymentInvoiceAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to ParticipantPaymentInvoiceAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//Convert ParticipantPaymentInvoiceAsset object to Json Message

func ParticipantPaymentInvoiceAssettoJson(obj ParticipantPaymentInvoiceAsset) ([]byte, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateParticipantPaymentInvoiceAsset Function will  insert record in ledger after receiving request from Client Application
func CreateParticipantPaymentInvoiceAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("CreateParticipantPaymentInvoiceAsset : Incorrect number of arguments.")
		return shim.Error("CreateParticipantPaymentInvoiceAsset : Incorrect number of arguments.")
	}

	// Convert the arg to a ParticipantPaymentInvoiceAsset Object
	logger.Infof("CreateParticipantPaymentInvoiceAsset: Arguments for ledgerapi %s : ", args[0])

	asset, err := JsontoParticipantPaymentInvoiceAsset([]byte(args[0]))
	asset.DocType = PAYMENT_INVOICE_RECORD_TYPE

	logger.Infof("CreateParticipantPaymentInvoiceAsset :Produce ID is : %s ", asset.ProduceID)

	keys = append(keys, asset.PaymentInvoiceID)
	keys = append(keys, asset.OrderID)
	keys = append(keys, strconv.FormatInt(int64(asset.Type), 10))
	keys = append(keys, asset.ID)

	logger.Infof("CreateParticipantPaymentInvoiceAsset : Inserting object with data as  %s", args[0])
	Avalbytes, _ = ParticipantPaymentInvoiceAssettoJson(asset)

	err = CreateAsset(stub, PAYMENT_INVOICE_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("CreateParticipantPaymentInvoiceAsset : Error inserting Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("CreateParticipantPaymentInvoiceAsset : ParticipantPaymentInvoiceAsset object create failed %s", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// QueryParticipantPaymentInvoiceAsset  Function will  query  record in ledger based on ID
func QueryParticipantPaymentInvoiceAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 4 {
		logger.Errorf("QueryParticipantPaymentInvoiceAsset : Incorrect number of arguments.")
		return shim.Error("QueryParticipantPaymentInvoiceAsset : Incorrect number of arguments.")
	}
	logger.Infof("QueryParticipantPaymentInvoiceAsset :ParticipantPaymentInvoiceAsset ID is : %s ", args[0])
	keys = append(keys, args[0])
	keys = append(keys, args[1])
	keys = append(keys, args[2])
	keys = append(keys, args[3])
	Avalbytes, err = QueryAsset(stub, PAYMENT_INVOICE_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("QueryParticipantPaymentInvoiceAsset : Error Querying Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("QueryParticipantPaymentInvoiceAsset : QueryParticipantPaymentInvoiceAsset object get failed %s", err))
	}

	return shim.Success([]byte(Avalbytes))

}

//UpdateParticipantPaymentInvoiceAsset Function will  update record in ledger after receiving request from Client Application
func UpdateParticipantPaymentInvoiceAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//args[0] payableAssetId
	//args[1] status
	//args[2] batchpaymentId

	var err error
	var Avalbytes []byte
	var keys []string

	if len(args) < 1 {
		logger.Errorf("UpdateParticipantPaymentInvoiceAsset : Incorrect number of arguments.")
		return shim.Error("UpdateParticipantPaymentInvoiceAsset : Incorrect number of arguments.")
	}
	asset, err := JsontoParticipantPaymentInvoiceAsset([]byte(args[0]))

	logger.Infof("UpdateParticipantPaymentInvoiceAsset :asset is : %v ", asset)

	keys = append(keys, asset.PaymentInvoiceID)
	keys = append(keys, asset.OrderID)
	keys = append(keys, strconv.FormatInt(int64(asset.Type), 10))
	keys = append(keys, asset.ID)
	Avalbytes, err = QueryAsset(stub, PAYMENT_INVOICE_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("UpdateParticipantPaymentInvoiceAsset : Asset is not available for this ParticipantPaymentInvoice")
		return shim.Error(fmt.Sprintf("UpdateParticipantPaymentInvoiceAsset : Asset is not available for this ParticipantPaymentInvoice"))
	}

	assetread, err := JsontoParticipantPaymentInvoiceAsset(Avalbytes)
	// Convert the arg to a ParticipantPaymentInvoiceAsset Object
	logger.Infof("UpdateParticipantPaymentInvoiceAsset: %v : ", assetread)

	assetread.Status = asset.Status
	assetread.AccountPaymentResponse = asset.AccountPaymentResponse
	logger.Infof("UpdateParticipantPaymentInvoiceAsset :assetread is : %v ", assetread)
	Avalbytes, _ = ParticipantPaymentInvoiceAssettoJson(assetread)
	logger.Infof("UpdateParticipantPaymentInvoiceAsset :updated data is : %s ", string(Avalbytes))

	err = UpdateAssetWithoutGet(stub, PAYMENT_INVOICE_RECORD_TYPE, keys, Avalbytes)
	if err != nil {
		logger.Errorf("UpdateParticipantPaymentInvoiceAsset : Error updating Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("UpdateParticipantPaymentInvoiceAsset : ParticipantPaymentInvoiceAsset object update failed %s", err))
	}
	err = stub.SetEvent("UpdatedPayableAsset", Avalbytes)
	if err != nil {
		logger.Errorf("UpdateParticipantPaymentInvoiceAsset : Error setting Object into Event %s", err)
	}
	return shim.Success([]byte(Avalbytes))
}

// ListPaymentInvoiceAsset  Function will  query  all record from DB
func ListParticipantPaymentInvoiceAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var paymentitr shim.StateQueryIteratorInterface
	var paymentlist []ParticipantPaymentInvoiceAsset
	var keys []string
	//queryString := fmt.Sprintf("{\"selector\":{\"PRODUCEID\":\"%s\"}}", args[0])
	//logger.Infof("ListParticipantPaymentInvoiceAssetbyProduceID Query string is %s ", queryString)
	paymentitr, err = ListAllAsset(stub, PAYMENT_INVOICE_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("ListParticipantPaymentInvoiceAssetbyProduceID : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer paymentitr.Close()
	for paymentitr.HasNext() {
		data, derr := paymentitr.Next()
		if derr != nil {
			logger.Errorf("ListParticipantPaymentInvoiceAssetbyProduceID : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("ListParticipantPaymentInvoiceAssetbyProduceID: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		payment, _ := JsontoParticipantPaymentInvoiceAsset([]byte(databyte))
		paymentlist = append(paymentlist, payment)
	}
	Avalbytes, err = json.Marshal(paymentlist)
	logger.Infof("ListParticipantPaymentInvoiceAssetbyProduceID Responce for App : %v", Avalbytes)
	if err != nil {
		logger.Errorf("ListParticipantPaymentInvoiceAssetbyProduceID : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListParticipantPaymentInvoiceAssetbyProduceID: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

func generateFarmerPayable(qtybreakdowns []TableVariety, farmerRate Rate, qtyper TableVarietyPerRateAsset) (float64, []QualityPaymentInfo) {
	var qualityPaymentInfos []QualityPaymentInfo
	var totalFarmerpayment float64
	for i := 0; i < len(qtybreakdowns); i++ {
		for j := 0; j < len(qtyper.QualityPerInfos); j++ {
			if qtybreakdowns[i].Name == qtyper.QualityPerInfos[j].Type {
				var qualityPaymentInfo QualityPaymentInfo
				qualityPaymentInfo.QualityType = qtyper.QualityPerInfos[j].Type
				qualityPaymentInfo.Currency = farmerRate.CurrencyUnit
				var finalqualityRate float64
				if qtyper.QualityPerInfos[j].OpMode {
					finalqualityRate = farmerRate.Value - farmerRate.Value*float64(qtyper.QualityPerInfos[j].PerValue)

				} else {
					finalqualityRate = farmerRate.Value + farmerRate.Value*float64(qtyper.QualityPerInfos[j].PerValue)

				}
				qualityPaymentInfo.Amount = qtybreakdowns[i].Value * finalqualityRate
				qualityPaymentInfo.Amount = math.Round(qualityPaymentInfo.Amount*100) / 100
				qualityPaymentInfo.Qty = uint64(qtybreakdowns[i].Value)
				logger.Debugf("generateFarmerPayable: Farmer quality pay for Quality [ %s ]  is  =  %f : ", qualityPaymentInfo.QualityType, qualityPaymentInfo.Amount)
				totalFarmerpayment = totalFarmerpayment + qualityPaymentInfo.Amount
				qualityPaymentInfos = append(qualityPaymentInfos, qualityPaymentInfo)
			}
		}
	}
	logger.Infof("generateFarmerPayable:Farmer Payment details as total %f and breakdown %v ", totalFarmerpayment, qualityPaymentInfos)
	return totalFarmerpayment, qualityPaymentInfos
}

// BreakdownRate is having compined rate for pc and quality persentage
type BreakdownRate struct {
	QualityName  string
	CurrencyUnit string
	Quantityunit string
	PCRatevalue  float64
	FarmerRate   float64
	DCRate       float64
	agentRate    float64
	wsRate       float64
	dfarmRate    float64
}

func getBreakdownRate(qualityper TableVarietyPerRateAsset, packageingRate PCPackagingRateAsset, marketPrice float64) map[string]BreakdownRate {
	breakdownmap := make(map[string]BreakdownRate)

	for pcqcount := 0; pcqcount < len(packageingRate.TableVerityPrice); pcqcount++ {
		for percount := 0; percount < len(qualityper.QualityPerInfos); percount++ {
			if packageingRate.TableVerityPrice[pcqcount].QualityName == qualityper.QualityPerInfos[percount].Type {
				pcasset := packageingRate.TableVerityPrice[pcqcount]
				var totalfarmerrate = 0.00
				if qualityper.QualityPerInfos[percount].OpMode {
					totalfarmerrate = marketPrice - marketPrice*float64(qualityper.QualityPerInfos[percount].PerValue)
				} else {
					totalfarmerrate = marketPrice + marketPrice*float64(qualityper.QualityPerInfos[percount].PerValue)
				}
				bqty := BreakdownRate{pcasset.QualityName, pcasset.Rates[0].CurrencyUnit, pcasset.Rates[0].QuantityUnit, pcasset.Rates[0].Value, totalfarmerrate, 0.00, 0.00, 0.00, 0.00}
				breakdownmap[pcasset.QualityName] = bqty
			}
		}
	}
	logger.Infof("Created Map Data: %v", breakdownmap)
	return breakdownmap
}

// GenerateParticipantPaymentInvoiceAsset Function will  insert record in ledger after receiving request from Client Application
func GenerateParticipantPaymentInvoiceAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte
	var keys, akeys []string
	var keybuyersset []string
	var userpaymentInfos []ParticipantPaymentInvoiceAsset
	var farmerRate Rate
	var qualityper TableVarietyPerRateAsset
	var casRate CASRate
	var agentRateper AgentRatePerAsset
	var wsRateper WholeSellerRatePerAsset
	var participantRate ParticipantRate
	var packageingRate PCPackagingRateAsset
	var found bool
	var breakdownrateMap map[string]BreakdownRate
	var totaldcrate float64
	var transportRate float64
	var pcid, dcid string
	var isquality int
	var agentRate = 0.00
	var wsRate = 0.00
	var dFrarmRate = 0.00
	if len(args) < 3 {
		logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Incorrect number of arguments.")
		return shim.Error("GenerateParticipantPaymentInvoiceAsset : Incorrect number of arguments.")
	}
	// Convert the arg to a ParticipantPaymentInvoiceAsset Object
	logger.Infof("GenerateParticipantPaymentInvoiceAsset: BuyerOrderId %s : ", args[0])
	logger.Infof("GenerateParticipantPaymentInvoiceAsset: PaymentInvoideId %s : ", args[1])
	logger.Infof("GenerateParticipantPaymentInvoiceAsset: Produce Name %s : ", args[2])

	keys = append(keys, args[1])
	Avalbytes, err = QueryAsset(stub, PAYMENT_INVOICE_RECORD_TYPE, keys)
	if err == nil {
		logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Asset alreday generated for this ParticipantPaymentInvoice")
		return shim.Error(fmt.Sprintf("GenerateParticipantPaymentInvoiceAsset : Asset alreday generated for this ParticipantPaymentInvoice"))
	}

	keybuyersset = append(keybuyersset, args[0])
	Avalbytes, err = QueryAsset(stub, BUYER_ORDER_ASSET_RECORD_TYPE, keybuyersset)
	if err != nil {
		logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error querying Buyer Order Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("GenerateParticipantPaymentInvoiceAsset : Error querying Buyer Order Object into LedgerState %s", err))
	}

	buyerorder, _ := JsontoBuyerOrderAsset(Avalbytes)
	buyerorder.PayableAssetID = args[1]
	Avalbytes, _ = BuyerOrderAssettoJson(buyerorder)
	err = UpdateAssetWithoutGet(stub, BUYER_ORDER_ASSET_RECORD_TYPE, keybuyersset, Avalbytes)
	if err != nil {
		logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error updating Buyer Order Object into LedgerState %s", err)
		return shim.Error(fmt.Sprintf("GenerateParticipantPaymentInvoiceAsset : Error queryupdatingBuyer Order Object into LedgerState %s", err))
	}

	participantMap := ArraytoMap(buyerorder.ParticipantInfos)
	if farmerparticipant, ok := participantMap[FARMER]; ok {
		farmerkey := []string{args[2], farmerparticipant.Country, farmerparticipant.State, buyerorder.Variety}
		if farmerRate, found = getFarmerMarketRate(stub, farmerkey); found {
			logger.Infof("Market Rate   for Variety %s  is : %v", buyerorder.Variety, farmerRate)

		} else {
			logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error Getting Farmer Market Rate ")
			return shim.Error(fmt.Sprintf("GenerateParticipantPaymentInvoiceAsset : Error Getting Farmer Market Rate %t", found))
		}
		isquality = len(buyerorder.QtyBreakdowns)
		if isquality != 0 {
			qualityperkey := []string{farmerparticipant.Country, farmerparticipant.State, args[2], buyerorder.Variety}
			if qualityper, found = getTableVarietyPerRateAsset(stub, qualityperkey); found {
				logger.Infof("Per Rate  is : %v", qualityper)
			} else {
				logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error Getting quality per Rate ")
				return shim.Error(fmt.Sprintf("GenerateParticipantPaymentInvoiceAsset : Error Getting quality per Rate %t", found))
			}

		}
	} else {
		logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error Getting Farmer in Participant List")
		return shim.Error(fmt.Sprintf("GenerateParticipantPaymentInvoiceAsset : Error Getting Farmer in Participant List %t", ok))

	}

	/*keys = append(keys, asset.Country)
	keys = append(keys, asset.State)
	keys = append(keys, asset.ProduceName)*/

	if participantRate, found = getParticipantRateAsset(stub, args[2], participantMap); found {
		logger.Infof("Participant Rate is : %v", participantRate)
	} else {
		logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error Getting Getting Participant Rate ")
		return shim.Error(fmt.Sprintf("GenerateParticipantPaymentInvoiceAsset : Error Getting Participant Rate %t", found))
	}
	//Getting Agent rate Percentage
	akeys = append(akeys, buyerorder.BuyerID) //Destination id i.e Buyer ID
	akeys = append(akeys, args[2])            //Produce Name
	akeys = append(akeys, buyerorder.Variety) //Variety
	agentRateper, _ = getAgentRatePerAsset(stub, akeys)
	//Getting Whole seller rate Percentage
	wsRateper, _ = getWholeSellerRatePerAsset(stub, akeys)

	for i := 0; i < len(buyerorder.ParticipantInfos); i++ {

		PayableAsset := ParticipantPaymentInvoiceAsset{}
		PayableAsset.DocType = PAYMENT_INVOICE_RECORD_TYPE
		PayableAsset.ProduceID = buyerorder.ProduceID
		PayableAsset.Variety = buyerorder.Variety
		PayableAsset.PaymentInvoiceID = args[1]
		PayableAsset.Qty = buyerorder.Qty
		PayableAsset.BuyerOrderID = buyerorder.OrderID
		PayableAsset.BuyerID = buyerorder.BuyerID
		PayableAsset.ID = buyerorder.ParticipantInfos[i].ID
		if len(buyerorder.ParticipantInfos[i].OrderID) > 0 {
			PayableAsset.OrderID = buyerorder.ParticipantInfos[i].OrderID
		} else {
			PayableAsset.OrderID = buyerorder.OrderID
		}
		PayableAsset.Type = buyerorder.ParticipantInfos[i].Type
		PayableAsset.Status = "pending"
		PayableAsset.Currency = farmerRate.CurrencyUnit

		switch role := buyerorder.ParticipantInfos[i].Type; role {
		case FARMER:
			//getFarmerMarketRate with Product Name,Country,State,Variety
			if isquality == 0 {
				PayableAsset.PayableAmount = farmerRate.Value * float64(buyerorder.Qty)
				PayableAsset.PayableAmount = math.Round(PayableAsset.PayableAmount*100) / 100
				logger.Infof("GenerateParticipantPaymentInvoiceAsset: Farmer Total Pay  %f : ", PayableAsset.PayableAmount)
				userpaymentInfos = append(userpaymentInfos, PayableAsset)
				break
			}
			PayableAsset.PayableAmount, PayableAsset.QualityPaymentInfos = generateFarmerPayable(buyerorder.QtyBreakdowns, farmerRate, qualityper)
			PayableAsset.PayableAmount = math.Round(PayableAsset.PayableAmount*100) / 100
			logger.Infof("GenerateParticipantPaymentInvoiceAsset: Farmer Total Pay  %f : ", PayableAsset)
			userpaymentInfos = append(userpaymentInfos, PayableAsset)
			logger.Infof("Updated ParticipantPaymentInfos is : %v", userpaymentInfos)
			break
			//pass array of quality payment and packaging rate  it will return pcpayment breakdown and total payment
		case CAS:
			if casRate, found = getCASRateAsset(stub, buyerorder.ParticipantInfos[i].ID); found {
				PayableAsset.Currency = casRate.CurrencyUnit
			} else {
				logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error Getting CAS Market Rate ")
				return shim.Error(fmt.Sprintf("GenerateParticipantPaymentInvoiceAsset : Error Getting CAS Market Rate %t", found))
			}
			PayableAsset.PayableAmount = float64(buyerorder.Qty) * casRate.Value
			PayableAsset.PayableAmount = math.Round(PayableAsset.PayableAmount*100) / 100
			logger.Infof("GenerateParticipantPaymentInvoiceAsset: CAS Total Pay  %f : ", PayableAsset.PayableAmount)
			userpaymentInfos = append(userpaymentInfos, PayableAsset)
			break

		case PC:
			/*	keys = append(keys, args[0]) //Country
				keys = append(keys, args[1]) //State
				keys = append(keys, args[2]) //Produce Name
				keys = append(keys, args[3]) //Variety
				keys = append(keys, args[4]) //pcID*/
			pcid = buyerorder.ParticipantInfos[i].ID
			pcpackingkey := []string{buyerorder.ParticipantInfos[i].Country, buyerorder.ParticipantInfos[i].State, args[2], buyerorder.Variety, buyerorder.ParticipantInfos[i].ID}

			if packageingRate, found = getPCPackagingRateAsset(stub, pcpackingkey); found {
				breakdownrateMap = getBreakdownRate(qualityper, packageingRate, farmerRate.Value)

				var qualityPaymentInfos []QualityPaymentInfo
				for qualitycount := 0; qualitycount < len(buyerorder.QtyBreakdowns); qualitycount++ {
					if breakdownrate, ok := breakdownrateMap[buyerorder.QtyBreakdowns[qualitycount].Name]; ok {
						var qualityPaymentInfo QualityPaymentInfo
						qualityPaymentInfo.QualityType = breakdownrate.QualityName
						qualityPaymentInfo.Currency = breakdownrate.CurrencyUnit
						qualityPaymentInfo.Qty = uint64(buyerorder.QtyBreakdowns[qualitycount].Value)
						qualityPaymentInfo.Amount = buyerorder.QtyBreakdowns[qualitycount].Value * breakdownrate.PCRatevalue
						qualityPaymentInfo.Amount = math.Round(qualityPaymentInfo.Amount*100) / 100
						PayableAsset.PayableAmount = PayableAsset.PayableAmount + qualityPaymentInfo.Amount
						logger.Infof("GenerateParticipantPaymentInvoiceAsset: Pack House quality pay for Quality [ %s ]  is  =  %f : ", qualityPaymentInfo.QualityType, qualityPaymentInfo.Amount)
						qualityPaymentInfos = append(qualityPaymentInfos, qualityPaymentInfo)
					} else {

						logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error getting PC rate for quality : %s ", breakdownrate.QualityName)

					}
				}
				PayableAsset.QualityPaymentInfos = qualityPaymentInfos
				PayableAsset.Currency = farmerRate.CurrencyUnit
				userpaymentInfos = append(userpaymentInfos, PayableAsset)
				logger.Infof("Updated ParticipantPaymentInfos at Pack House is : %v", userpaymentInfos)
				break

			} else {
				logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error Getting Pack House Rate ")
				return shim.Error(fmt.Sprintf("GenerateParticipantPaymentInvoiceAsset : Error Getting quality per Rate %t", found))
			}

		case DC:
			dcid = buyerorder.ParticipantInfos[i].ID
			if isquality == 0 {
				logger.Infof("FarmerRate  and CAS rate for dc payable calculation : %f ,% f", farmerRate.Value, casRate.Value)
				totaldcrate = (farmerRate.Value + casRate.Value) * float64(participantRate.DCpercentage)
				logger.Infof("DC payable rate with no Pack house %f", totaldcrate)
				PayableAsset.Currency = farmerRate.CurrencyUnit
				PayableAsset.PayableAmount = float64(buyerorder.Qty) * totaldcrate
				PayableAsset.PayableAmount = math.Round(PayableAsset.PayableAmount*100) / 100
				logger.Infof("GenerateParticipantPaymentInvoiceAsset: DC Total Pay  %f : ", PayableAsset.PayableAmount)
				userpaymentInfos = append(userpaymentInfos, PayableAsset)
				break
			}
			var qualityPaymentInfos []QualityPaymentInfo
			for qualitycount := 0; qualitycount < len(buyerorder.QtyBreakdowns); qualitycount++ {
				if breakdownrate, ok := breakdownrateMap[buyerorder.QtyBreakdowns[qualitycount].Name]; ok {
					var qualityPaymentInfo QualityPaymentInfo
					qualityPaymentInfo.QualityType = breakdownrate.QualityName
					qualityPaymentInfo.Qty = uint64(buyerorder.QtyBreakdowns[qualitycount].Value)
					qualityPaymentInfo.Currency = breakdownrate.CurrencyUnit
					qualityRatefordc := (breakdownrate.FarmerRate + breakdownrate.PCRatevalue + casRate.Value) * float64(participantRate.DCpercentage)
					logger.Infof("Input Rates Farmer quality rate:   %f , casRate: %f ,Packhouse rate: %f and dc per %f :for quality", breakdownrate.FarmerRate, casRate.Value, breakdownrate.PCRatevalue, participantRate.DCpercentage)
					logger.Infof("Final DC Rate for quality : %f", qualityRatefordc)
					qualityPaymentInfo.Amount = buyerorder.QtyBreakdowns[qualitycount].Value * qualityRatefordc
					qualityPaymentInfo.Amount = math.Round(qualityPaymentInfo.Amount*100) / 100
					PayableAsset.PayableAmount = PayableAsset.PayableAmount + qualityPaymentInfo.Amount
					logger.Infof("GenerateParticipantPaymentInvoiceAsset: DC  quality pay for Quality [ %s ]  is  =  %f : ", qualityPaymentInfo.QualityType, qualityPaymentInfo.Amount)
					qualityPaymentInfos = append(qualityPaymentInfos, qualityPaymentInfo)
					breakdownrate.DCRate = qualityRatefordc
					breakdownrateMap[buyerorder.QtyBreakdowns[qualitycount].Name] = breakdownrate
				} else {

					logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error getting DC rate for quality : %s ", breakdownrate.QualityName)

				}
			}
			PayableAsset.QualityPaymentInfos = qualityPaymentInfos
			PayableAsset.Currency = farmerRate.CurrencyUnit
			userpaymentInfos = append(userpaymentInfos, PayableAsset)
			logger.Infof("Updated ParticipantPaymentInfos at DC is : %v", userpaymentInfos)
			break

		case PCDCTRANS:
			distance, _ := getTransportInfoAsset(stub, pcid, dcid)
			transportRate = float64(distance.Distance.Value) * participantRate.DeliveryRateVal.Value
			logger.Infof("GenerateParticipantPaymentInvoiceAsset:  Dfarm Transport rate  %f : ", transportRate)
			PayableAsset.Currency = farmerRate.CurrencyUnit
			PayableAsset.PayableAmount = float64(buyerorder.Qty) * transportRate
			PayableAsset.PayableAmount = math.Round(PayableAsset.PayableAmount*100) / 100
			logger.Infof("GenerateParticipantPaymentInvoiceAsset:  Dfarm Transport Total Pay  %f : ", PayableAsset.PayableAmount)
			userpaymentInfos = append(userpaymentInfos, PayableAsset)
			break

		case DEMURRAGEACC:
			if isquality == 0 {

				logger.Infof("FarmerRate  and CAS rate  and DC rate for Demurrage payable calculation : %f ,%f and %f", farmerRate.Value, casRate.Value, totaldcrate)
				logger.Infof("AgentRate  and WS rate  and DFarm Rate  for Demurrage payable calculation : %f ,%f and %f", agentRate, wsRate, dFrarmRate)

				demurrrate := (farmerRate.Value + casRate.Value + transportRate + totaldcrate + agentRate + wsRate + dFrarmRate) * participantRate.Dfarmpercentage.DemurragePer
				logger.Infof("dFarm payable rate with no Pack house %f", demurrrate)
				PayableAsset.Currency = farmerRate.CurrencyUnit
				PayableAsset.PayableAmount = float64(buyerorder.Qty) * demurrrate
				PayableAsset.PayableAmount = math.Round(PayableAsset.PayableAmount*100) / 100
				logger.Infof("GenerateParticipantPaymentInvoiceAsset: dFarm Total Pay  %f : ", PayableAsset.PayableAmount)
				userpaymentInfos = append(userpaymentInfos, PayableAsset)
				break
			}
			var qualityPaymentInfos []QualityPaymentInfo
			for qualitycount := 0; qualitycount < len(buyerorder.QtyBreakdowns); qualitycount++ {
				if breakdownrate, ok := breakdownrateMap[buyerorder.QtyBreakdowns[qualitycount].Name]; ok {
					var qualityPaymentInfo QualityPaymentInfo
					qualityPaymentInfo.QualityType = breakdownrate.QualityName
					qualityPaymentInfo.Qty = uint64(buyerorder.QtyBreakdowns[qualitycount].Value)
					qualityPaymentInfo.Currency = breakdownrate.CurrencyUnit
					logger.Infof("GenerateParticipantPaymentInvoiceAsset:  Demurr Transport rate  %f : ", transportRate)
					logger.Infof("Input Rates Farmer quality rate:   %f , casRate: %f ,Packhouse rate: %f and dFarm  per %f :for quality", breakdownrate.FarmerRate, casRate.Value, breakdownrate.PCRatevalue, participantRate.Dfarmpercentage.DemurragePer)
					logger.Infof("Input Rates agent rate:   %f , wsrate: %f ,dfarmrate : %f  :for quality", breakdownrate.agentRate, breakdownrate.wsRate, breakdownrate.dfarmRate)

					qualityRatefordmurr := (breakdownrate.FarmerRate + breakdownrate.PCRatevalue + casRate.Value + breakdownrate.DCRate + transportRate + breakdownrate.agentRate + breakdownrate.dfarmRate + breakdownrate.wsRate) * participantRate.Dfarmpercentage.DemurragePer

					logger.Infof("Final Demurrage Rate for quality : %f", qualityRatefordmurr)
					qualityPaymentInfo.Amount = buyerorder.QtyBreakdowns[qualitycount].Value * qualityRatefordmurr
					qualityPaymentInfo.Amount = math.Round(qualityPaymentInfo.Amount*100) / 100
					PayableAsset.PayableAmount = PayableAsset.PayableAmount + qualityPaymentInfo.Amount
					PayableAsset.PayableAmount = math.Round(PayableAsset.PayableAmount*100) / 100
					logger.Infof("GenerateParticipantPaymentInvoiceAsset: Demurrage pay for Quality [ %s ]  is  =  %f : ", qualityPaymentInfo.QualityType, qualityPaymentInfo.Amount)
					qualityPaymentInfos = append(qualityPaymentInfos, qualityPaymentInfo)
				} else {

					logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error getting dFarm rate for quality : %s ", breakdownrate.QualityName)

				}
			}
			PayableAsset.Currency = farmerRate.CurrencyUnit
			PayableAsset.QualityPaymentInfos = qualityPaymentInfos
			userpaymentInfos = append(userpaymentInfos, PayableAsset)
			logger.Infof("Updated ParticipantPaymentInfos at dFarma is : %v", userpaymentInfos)
		case DFARM:

			agentPayableAsset := PayableAsset
			agentPayableAsset.Type = AGENT
			agentPayableAsset.ID = agentRateper.SourceID
			wsPayableAsset := PayableAsset
			wsPayableAsset.ID = wsRateper.SourceID
			wsPayableAsset.Type = WHOLESELLER
			if isquality == 0 {

				if agentRateper.PerValue > 0 {

					logger.Infof("FarmerRate  and CAS rate  and DC rate for Dfarm payable calculation : %f ,%f and %f", farmerRate.Value, casRate.Value, totaldcrate)
					agentRate = (farmerRate.Value + casRate.Value + totaldcrate) * agentRateper.PerValue
					logger.Infof("Agent payable rate with no Pack house %f", agentRate)
					agentapaybleamount := float64(buyerorder.Qty) * agentRate
					agentapaybleamount = math.Round(agentapaybleamount*100) / 100
					logger.Infof("GenerateParticipantPaymentInvoiceAsset: agent Total Pay  %f : ", agentapaybleamount)
					agentPayableAsset.PayableAmount = agentapaybleamount
					userpaymentInfos = append(userpaymentInfos, agentPayableAsset)
					logger.Infof("GeneratePaymentInvoiceAsset: Payable Asset at agent level  %v : ", userpaymentInfos)
				}
				if wsRateper.PerValue > 0 {
					logger.Infof("FarmerRate  and CAS rate  and DC rate for Dfarm payable calculation : %f ,%f and %f", farmerRate.Value, casRate.Value, totaldcrate)
					wsRate = (farmerRate.Value + casRate.Value + totaldcrate) * wsRateper.PerValue
					logger.Infof("Whole saller payable rate with no Pack house %f", wsRate)
					wsapaybleamount := float64(buyerorder.Qty) * wsRate
					wsapaybleamount = math.Round(wsapaybleamount*100) / 100
					logger.Infof("GenerateParticipantPaymentInvoiceAsset: WS Total Pay  %f : ", wsapaybleamount)
					wsPayableAsset.PayableAmount = wsapaybleamount
					userpaymentInfos = append(userpaymentInfos, wsPayableAsset)
					logger.Infof("GeneratePaymentInvoiceAsset: Payable Asset at Whole Seller level  %v : ", userpaymentInfos)
				}
				dFrarmRate = (farmerRate.Value + casRate.Value + transportRate + totaldcrate + agentRate) * float64(participantRate.Dfarmpercentage.Totalper)
				logger.Infof("dFarm payable rate with no Pack house %f", dFrarmRate)
				PayableAsset.PayableAmount = float64(buyerorder.Qty) * dFrarmRate
				PayableAsset.PayableAmount = math.Round(PayableAsset.PayableAmount*100) / 100
				logger.Infof("GenerateParticipantPaymentInvoiceAsset: dFarm Total Pay  %f : ", PayableAsset.PayableAmount)
				userpaymentInfos = append(userpaymentInfos, PayableAsset)
				logger.Infof("GeneratePaymentInvoiceAsset: Payable Asset at dFarm level  %v : ", userpaymentInfos)
				break
			}
			var qualityPaymentInfos, agentqualityPaymentInfos, wsqualityPaymentInfos []QualityPaymentInfo

			for qualitycount := 0; qualitycount < len(buyerorder.QtyBreakdowns); qualitycount++ {
				var agentRate = 0.0
				var wsRate = 0.0
				if breakdownrate, ok := breakdownrateMap[buyerorder.QtyBreakdowns[qualitycount].Name]; ok {
					var qualityPaymentInfo QualityPaymentInfo
					qualityPaymentInfo.QualityType = breakdownrate.QualityName
					qualityPaymentInfo.Qty = uint64(buyerorder.QtyBreakdowns[qualitycount].Value)
					qualityPaymentInfo.Currency = breakdownrate.CurrencyUnit
					logger.Infof("GenerateParticipantPaymentInvoiceAsset:  Dfarm Transport rate  %f : ", transportRate)
					if agentRateper.PerValue > 0 {
						var agentqualityPayInfo QualityPaymentInfo
						agentqualityPayInfo = qualityPaymentInfo
						logger.Infof("Input Rates  for Agent Farmer quality rate:   %f , CAS Rate : %f ,Packhouse rate: %f and Dc Rate %f :for quality %s ", breakdownrate.FarmerRate, casRate.Value, breakdownrate.PCRatevalue, breakdownrate.DCRate, breakdownrate.QualityName)
						agentRate = (breakdownrate.FarmerRate + breakdownrate.PCRatevalue + casRate.Value + breakdownrate.DCRate) * agentRateper.PerValue
						breakdownrate.agentRate = agentRate
						breakdownrateMap[buyerorder.QtyBreakdowns[qualitycount].Name] = breakdownrate
						logger.Infof("GenerateParticipantPaymentInvoiceAsset: Agent Rate  is %f : ", agentRate)
						agentqualityPayInfo.Amount = buyerorder.QtyBreakdowns[qualitycount].Value * agentRate
						agentqualityPayInfo.Amount = math.Round(agentqualityPayInfo.Amount*100) / 100
						agentPayableAsset.PayableAmount = agentPayableAsset.PayableAmount + agentqualityPayInfo.Amount
						logger.Infof("GenerateParticipantPaymentInvoiceAsset: Agent pay for Quality [ %s ]  is  =  %f : ", agentqualityPayInfo.QualityType, agentqualityPayInfo.Amount)
						agentqualityPaymentInfos = append(agentqualityPaymentInfos, agentqualityPayInfo)

					}

					if wsRateper.PerValue > 0 {
						var wsqualityPayInfo QualityPaymentInfo
						wsqualityPayInfo = qualityPaymentInfo
						logger.Infof("Input Rates  for Agent Farmer quality rate:   %f , CAS Rate : %f ,Packhouse rate: %f and Dc Rate %f :for quality %s ", breakdownrate.FarmerRate, casRate.Value, breakdownrate.PCRatevalue, breakdownrate.DCRate, breakdownrate.QualityName)
						wsRate = (breakdownrate.FarmerRate + breakdownrate.PCRatevalue + casRate.Value + breakdownrate.DCRate) * wsRateper.PerValue
						breakdownrate.wsRate = wsRate
						logger.Infof("GenerateParticipantPaymentInvoiceAsset: ws Rate  is %f : ", wsRate)
						wsqualityPayInfo.Amount = buyerorder.QtyBreakdowns[qualitycount].Value * wsRate
						wsqualityPayInfo.Amount = math.Round(wsqualityPayInfo.Amount*100) / 100
						wsPayableAsset.PayableAmount = wsPayableAsset.PayableAmount + wsqualityPayInfo.Amount
						logger.Infof("GenerateParticipantPaymentInvoiceAsset: WS pay for Quality [ %s ]  is  =  %f : ", wsqualityPayInfo.QualityType, wsqualityPayInfo.Amount)
						wsqualityPaymentInfos = append(wsqualityPaymentInfos, wsqualityPayInfo)

					}

					qualityRatefordFarm := (breakdownrate.FarmerRate + breakdownrate.PCRatevalue + casRate.Value + breakdownrate.DCRate + transportRate + agentRate) * participantRate.Dfarmpercentage.Totalper
					logger.Infof("Input Rates Farmer quality rate:   %f , casRate: %f ,Packhouse rate: %f and dFarm  per %f :for quality", breakdownrate.FarmerRate, casRate.Value, breakdownrate.PCRatevalue, participantRate.Dfarmpercentage.Totalper)
					logger.Infof("DC Rate is : %f ", breakdownrate.DCRate)
					logger.Infof("Final dFarm Rate for quality : %f", qualityRatefordFarm)
					breakdownrate.dfarmRate = qualityRatefordFarm
					breakdownrateMap[buyerorder.QtyBreakdowns[qualitycount].Name] = breakdownrate
					qualityPaymentInfo.Amount = buyerorder.QtyBreakdowns[qualitycount].Value * qualityRatefordFarm
					qualityPaymentInfo.Amount = math.Round(qualityPaymentInfo.Amount*100) / 100
					PayableAsset.PayableAmount = PayableAsset.PayableAmount + qualityPaymentInfo.Amount
					logger.Infof("GenerateParticipantPaymentInvoiceAsset: dFarm pay for Quality [ %s ]  is  =  %f : ", qualityPaymentInfo.QualityType, qualityPaymentInfo.Amount)
					qualityPaymentInfos = append(qualityPaymentInfos, qualityPaymentInfo)
				} else {

					logger.Errorf("GenerateParticipantPaymentInvoiceAsset : Error getting dFarm rate for quality : %s ", breakdownrate.QualityName)

				}
			}
			if agentRateper.PerValue > 0 {
				agentPayableAsset.QualityPaymentInfos = agentqualityPaymentInfos
				userpaymentInfos = append(userpaymentInfos, agentPayableAsset)
				logger.Infof("Updated ParticipantPaymentInfos at Agent is : %v", userpaymentInfos)
			}
			if wsRateper.PerValue > 0 {
				agentPayableAsset.QualityPaymentInfos = wsqualityPaymentInfos
				userpaymentInfos = append(userpaymentInfos, wsPayableAsset)
				logger.Infof("Updated ParticipantPaymentInfos at WS is : %v", userpaymentInfos)
			}
			PayableAsset.QualityPaymentInfos = qualityPaymentInfos
			PayableAsset.Currency = farmerRate.CurrencyUnit
			userpaymentInfos = append(userpaymentInfos, PayableAsset)
			logger.Infof("Updated ParticipantPaymentInfos at dFarma is : %v", userpaymentInfos)
			break
		default:
			logger.Infof("Correct Party is not defined for payable asset: %d", role)

		}
	}

	var Payment []byte
	Payment, _ = json.Marshal(userpaymentInfos)
	logger.Infof("GenerateParticipantPaymentInvoiceAsset :PayableAsset is %s ", string(Payment))
	err = createPayableAssetfromList(stub, userpaymentInfos)
	if err != nil {
		logger.Errorf("Issue with creating payable asset in ledger with error as %v", err)
	}
	return shim.Success([]byte(Payment))
}

func createPayableAssetfromList(stub shim.ChaincodeStubInterface, userpaymentInfos []ParticipantPaymentInvoiceAsset) (err error) {

	if len(userpaymentInfos) == 0 {
		logger.Errorf("createPayableAssetfromList : creation list is empty.")
		return fmt.Errorf("createPayableAssetfromList : creation list is empty")
	}

	for userpayablecount := 0; userpayablecount < len(userpaymentInfos); userpayablecount++ {

		var keys []string
		var Available []byte
		asset := userpaymentInfos[userpayablecount]
		keys = append(keys, asset.PaymentInvoiceID, asset.OrderID, strconv.FormatInt(int64(asset.Type), 10), asset.ID)
		logger.Infof("createPayableAssetfromList : Generated Key is : %s", keys)
		Available, _ = ParticipantPaymentInvoiceAssettoJson(asset)
		logger.Infof("createPayableAssetfromList : Byte going to be added in the ledger is : %s", Available)
		err := CreateAsset(stub, PAYMENT_INVOICE_RECORD_TYPE, keys, Available)
		if err != nil {
			logger.Errorf("createPayableAssetfromList : Create fialed with error %v .", err)
			return err
		}

	}
	return nil
}

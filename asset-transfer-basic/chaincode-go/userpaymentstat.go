package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type UserPaymentStatRequest struct {
	ProduceID        string `json:"PRODUCEID,omitempty"`
	Variety          string `json:"VARIETY,omitempty"`
	UserType         int    `json:"USERTYPE,omitempty"`
	UserID           string `json:"USERID,omitempty"`
	BuyerOrderID     string `json:"BUYERORDERID,omitempty"`
	UserOrderID      string `json:"USERORDERID,omitempty"`
	PayableInvoiceID string `json:"PAYABLEINVOICEID,omitempty"`
}

type PaymentStat struct {
	Variety                      string               `json:"VARIETY,omitempty"`
	OrderID                      string               `json:"ORDERID,omitempty"`
	Currency                     string               `json:"CURRENCY,omitempty"`
	AvailableQtyforSell          uint64               `json:"AVAILABLEQTYFORSELL,omitempty"`
	AvailableQtyBreakdownforSell []QualityPaymentInfo `json:"AVAILABLEQTYBREAKDOWNFORSELL,omitempty"`
	QtyforReceivedPayment        uint64               `json:"QTYFORRECEIVEDPAYMENT,omitempty"`
	ReceivedPayment              float64              `json:"RECEIVEDPAYMENT,omitempty"`
	BreakdownReceivedPayment     []QualityPaymentInfo `json:"BREAKDOWNRECEIVEDPAYMENT,omitempty"`
	QtyforPendingPayment         uint64               `json:"QTYFORPENDINGPAYMENT,omitempty"`
	BreakdownPendingPayment      []QualityPaymentInfo `json:"BREAKDOWNPENDINGPAYMENT,omitempty"`
	PendingPayment               float64              `json:"PENDINGPAYMENT,omitempty"`
}
type UserPaymentStatResponse struct {
	ProduceID    string        `json:"PRODUCEID,omitempty"`
	PaymentStats []PaymentStat `json:"PAYMENTSTATS,omitempty"`
}

//GetPaymentStat(stub shim.ChaincodeStubInterface, args []string) pb.Response)

func GetPaymentStat(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var statlist []PaymentStat
	var keys []string
	var queryString string
	var payableAssetitr shim.StateQueryIteratorInterface
	if len(args) < 1 {
		return shim.Error("GetPaymentStat :Incorrect number of arguments. ")
	}
	var response UserPaymentStatResponse
	request := UserPaymentStatRequest{}
	err = json.Unmarshal([]byte(args[0]), &request)
	if len(request.ProduceID) ==0 {
		return shim.Error("GetPaymentStat :Produce Id is not set ")
	}
	response.ProduceID = request.ProduceID
	if request.UserType==0 && len(request.UserID)==0 {
		return shim.Error("GetPaymentStat :Incorrect number of arguments.UserType or UserID is not set in the request  ")
	}
	if  len(request.PayableInvoiceID)>0  && len(request.UserOrderID)>0 {
	keys = append(keys, request.PayableInvoiceID)
	keys = append(keys, request.UserOrderID)
	keys = append(keys, strconv.FormatInt(int64(request.UserType), 10))
	keys = append(keys, request.UserID)
	Avalbytes, err = QueryAsset(stub, PAYMENT_INVOICE_RECORD_TYPE, keys)
	if err != nil {
		logger.Errorf("GetPaymentStat : Can not read payable asset with  Error : %v", err)
		return shim.Error(fmt.Sprintf("GetPaymentStat: Can not read payable asset with  Error : %v", err))
	}
	payment, _ := JsontoParticipantPaymentInvoiceAsset([]byte(Avalbytes))
	logger.Infof("GetPaymentStat Payment struct : %v", payment)
	stat := PaymentStat{}
	stat.Currency=payment.Currency
	stat.OrderID = payment.OrderID
	stat.Variety = payment.Variety
	if payment.Status == "pending" {
		stat.PendingPayment = stat.PendingPayment + payment.PayableAmount
		stat.BreakdownPendingPayment = append(stat.BreakdownPendingPayment, payment.QualityPaymentInfos...)
		stat.QtyforPendingPayment = stat.QtyforPendingPayment + payment.Qty
	} else {
		stat.ReceivedPayment = stat.ReceivedPayment + payment.PayableAmount
		stat.QtyforReceivedPayment = stat.QtyforReceivedPayment + payment.Qty
		stat.BreakdownReceivedPayment = append(stat.BreakdownReceivedPayment, payment.QualityPaymentInfos...)
	}
	statlist = append(statlist, stat)
	response.ProduceID = request.ProduceID
	response.PaymentStats = statlist
	Avalbytes, err = json.Marshal(response)
	logger.Infof("GetPaymentStat Responce for App : %s", string(Avalbytes))
	if err != nil {
		logger.Errorf("GetPaymentStat : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("GetPaymentStat: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
} else if len(request.UserOrderID)>0 {
	return GetPaymentStatbyUserOrderID(stub,request)
} else if len(request.BuyerOrderID)>0 {
	return GetPaymentStatbyBuyerOrderID(stub,request)

} else {
if len(request.Variety) == 0 {
	queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PRODUCEID\":\"%s\",\"TYPE\":%d,\"ID\":\"%s\"}}", PAYMENT_INVOICE_RECORD_TYPE, request.ProduceID, request.UserType, request.UserID)
} else {
	queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PRODUCEID\":\"%s\",\"VARIETY\":\"%s\",\"TYPE\":%d,\"ID\":\"%s\"}}", PAYMENT_INVOICE_RECORD_TYPE, request.ProduceID, request.Variety, request.UserType, request.UserID)
}
logger.Infof("GetPaymentStatbyBuyerOrderID Query string is %s ", queryString)

payableAssetitr, err = GenericQueryAsset(stub, queryString)
if err != nil {
	logger.Errorf("GetPaymentStatbyBuyerOrderID : Instence not found in ledger")
	return shim.Error("orderitr : Instence not found in ledger")

}
defer payableAssetitr.Close()
for payableAssetitr.HasNext() {
	data, derr := payableAssetitr.Next()
	if derr != nil {
		logger.Errorf("GetPaymentStatbyBuyerOrderID : Cannot parse result set. Error : %v", derr)
		return shim.Error(fmt.Sprintf("GetPaymentStatbyBuyerOrderID: Cannot parse result set. Error : %v", derr))

	}
	databyte := data.GetValue()

	payment, _ := JsontoParticipantPaymentInvoiceAsset([]byte(databyte))
	stat := PaymentStat{}
	stat.Currency=payment.Currency
	stat.OrderID = payment.OrderID
	stat.Variety = payment.Variety
	if payment.Status == "pending" {
		stat.PendingPayment = stat.PendingPayment + payment.PayableAmount
		stat.BreakdownPendingPayment = append(stat.BreakdownPendingPayment, payment.QualityPaymentInfos...)
		stat.QtyforPendingPayment = stat.QtyforPendingPayment + payment.Qty
	} else {
		stat.ReceivedPayment = stat.ReceivedPayment + payment.PayableAmount
		stat.QtyforReceivedPayment = stat.QtyforReceivedPayment + payment.Qty
		stat.BreakdownReceivedPayment = append(stat.BreakdownReceivedPayment, payment.QualityPaymentInfos...)

	}
	statlist = append(statlist, stat)

}
response.ProduceID = request.ProduceID
response.PaymentStats = statlist
Avalbytes, err = json.Marshal(response)
logger.Infof("GetPaymentStatbyBuyerOrderID Responce for App : %s", string(Avalbytes))
if err != nil {
	logger.Errorf("GetPaymentStatbyBuyerOrderID : Cannot Marshal result set. Error : %v", err)
	return shim.Error(fmt.Sprintf("GetPaymentStatbyBuyerOrderID: Cannot Marshal result set. Error : %v", err))
}
return shim.Success([]byte(Avalbytes))
}
}

// GetPaymentStatbyUserOrderID Function will  get record from ledger after receiving requist  from Client Application
func GetPaymentStatbyUserOrderID(stub shim.ChaincodeStubInterface, request UserPaymentStatRequest) pb.Response {
	var err error
	var Avalbytes []byte
	var payableAssetitr shim.StateQueryIteratorInterface
	var response UserPaymentStatResponse
	var queryString string
	var statlist []PaymentStat
	if len(request.Variety) == 0 {
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PRODUCEID\":\"%s\",\"ORDERID\":\"%s\",\"TYPE\":%d,\"ID\":\"%s\"}}", PAYMENT_INVOICE_RECORD_TYPE, request.ProduceID, request.UserOrderID, request.UserType, request.UserID)
	} else {
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PRODUCEID\":\"%s\",\"VARIETY\":\"%s\",\"ORDERID\":\"%s\",\"TYPE\":%d,\"ID\":\"%s\"}}", PAYMENT_INVOICE_RECORD_TYPE, request.ProduceID, request.Variety, request.UserOrderID, request.UserType, request.UserID)
	}
	logger.Infof("GetPaymentStatbyUserOrderID Query string is %s ", queryString)

	payableAssetitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("GetPaymentStatbyUserOrderID : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	defer payableAssetitr.Close()
	for payableAssetitr.HasNext() {
		data, derr := payableAssetitr.Next()
		if derr != nil {
			logger.Errorf("GetPaymentStatbyUserOrderID : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("GetPaymentStatbyUserOrderID: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()
		logger.Infof("GetPaymentStatbyUserOrderID Databyte is  : %s", string(databyte))

		payment, _ := JsontoParticipantPaymentInvoiceAsset([]byte(databyte))
		logger.Infof("GetPaymentStatbyUserOrderID payment is  : %v", payment)
		stat := PaymentStat{}
		stat.Currency=payment.Currency
		stat.OrderID = payment.OrderID
		stat.Variety = payment.Variety
		if payment.Status == "pending" {
			stat.PendingPayment = stat.PendingPayment + payment.PayableAmount
			stat.BreakdownPendingPayment = append(stat.BreakdownPendingPayment, payment.QualityPaymentInfos...)
			stat.QtyforPendingPayment = stat.QtyforPendingPayment + payment.Qty
		} else {
			stat.ReceivedPayment = stat.ReceivedPayment + payment.PayableAmount
			stat.QtyforReceivedPayment = stat.QtyforReceivedPayment + payment.Qty
			stat.BreakdownReceivedPayment = append(stat.BreakdownReceivedPayment, payment.QualityPaymentInfos...)
		}
		statlist = append(statlist, stat)

	}
	response.ProduceID = request.ProduceID
	response.PaymentStats = statlist
	Avalbytes, err = json.Marshal(response)
	logger.Infof("GetPaymentStatbyUserOrderID Responce for App : %s", string(Avalbytes))
	if err != nil {
		logger.Errorf("GetPaymentStatbyUserOrderID : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("GetPaymentStatbyUserOrderID: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}

// GetPaymentStatbyBuyerOrderID Function will  get record from ledger after receiving requist  from Client Application
func GetPaymentStatbyBuyerOrderID(stub shim.ChaincodeStubInterface, request UserPaymentStatRequest) pb.Response {
	var err error
	var Avalbytes []byte
	var payableAssetitr shim.StateQueryIteratorInterface
	var response UserPaymentStatResponse
	var queryString string

	response.ProduceID = request.ProduceID

	if err != nil {
		return shim.Error("GetPaymentStatbyBuyerOrderID :json PrasingError ")
	}
	if len(request.Variety) == 0 {
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PRODUCEID\":\"%s\",\"BUYERORDERID\":\"%s\",\"TYPE\":%d,\"ID\":\"%s\"}}", PAYMENT_INVOICE_RECORD_TYPE, request.ProduceID, request.BuyerOrderID, request.UserType, request.UserID)
	} else {
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"PRODUCEID\":\"%s\",\"VARIETY\":\"%s\",\"BUYERORDERID\":\"%s\",\"TYPE\":%d,\"ID\":\"%s\"}}", PAYMENT_INVOICE_RECORD_TYPE, request.ProduceID, request.Variety, request.BuyerOrderID, request.UserType, request.UserID)
	}
	logger.Infof("GetPaymentStatbyBuyerOrderID Query string is %s ", queryString)

	payableAssetitr, err = GenericQueryAsset(stub, queryString)
	if err != nil {
		logger.Errorf("GetPaymentStatbyBuyerOrderID : Instence not found in ledger")
		return shim.Error("orderitr : Instence not found in ledger")

	}
	var statlist []PaymentStat
	defer payableAssetitr.Close()
	for payableAssetitr.HasNext() {
		data, derr := payableAssetitr.Next()
		if derr != nil {
			logger.Errorf("GetPaymentStatbyBuyerOrderID : Cannot parse result set. Error : %v", derr)
			return shim.Error(fmt.Sprintf("GetPaymentStatbyBuyerOrderID: Cannot parse result set. Error : %v", derr))

		}
		databyte := data.GetValue()

		payment, _ := JsontoParticipantPaymentInvoiceAsset([]byte(databyte))
		stat := PaymentStat{}
		stat.Currency=payment.Currency
		stat.OrderID = payment.OrderID
		stat.Variety = payment.Variety
		if payment.Status == "pending" {
			stat.PendingPayment = stat.PendingPayment + payment.PayableAmount
			stat.BreakdownPendingPayment = append(stat.BreakdownPendingPayment, payment.QualityPaymentInfos...)
			stat.QtyforPendingPayment = stat.QtyforPendingPayment + payment.Qty
		} else {
			stat.ReceivedPayment = stat.ReceivedPayment + payment.PayableAmount
			stat.QtyforReceivedPayment = stat.QtyforReceivedPayment + payment.Qty
			stat.BreakdownReceivedPayment = append(stat.BreakdownReceivedPayment, payment.QualityPaymentInfos...)

		}
		statlist = append(statlist, stat)

	}
	response.ProduceID = request.ProduceID
	response.PaymentStats = statlist
	Avalbytes, err = json.Marshal(response)
	logger.Infof("GetPaymentStatbyBuyerOrderID Responce for App : %s", string(Avalbytes))
	if err != nil {
		logger.Errorf("GetPaymentStatbyBuyerOrderID : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("GetPaymentStatbyBuyerOrderID: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))
}


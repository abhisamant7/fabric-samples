package main
/*

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	stub      *shim.MockStub
	status200 int32
)

func TestGoproducesc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Goproducesc Suite")
}

var _ = Describe("Unit Test Produce BlockChain Smart Contract", func() {
	stub = shim.NewMockStub("Testing Produce Blockahin", new(Dfarmsc))
	BeforeSuite(func() {
		receivedStatus := stub.MockInit("000", nil).Status
		status200 = int32(200)
		Expect(receivedStatus).Should(Equal(status200))
		PayableAsset := PaymentInvoiceAsset{}
		PayableAsset.BuyerOrderID = "001"
		PayableAsset.DocType = PAYMENT_INVOICE_RECORD_TYPE
		PayableAsset.ProduceID = "Pro001"
		PayableAsset.Variety = "Gala"
		PayableAsset.PaymentInvoiceID = "Inv001"
		PayableAsset.Qty = 200
		PayableAsset.Status = "Pending"
		var payableinfo []ParticipantPaymentInfo
		farmerPayableInfo := ParticipantPaymentInfo{"Far001", 200.25, "Ordfar001", "", "pending"}
		casPayableInfo := ParticipantPaymentInfo{"Cas001", 20.25, "Ordcas001", "", "pending"}
		pcPayableInfo := ParticipantPaymentInfo{"Pc001", 21.25, "Ordpc001", "", "pending"}
		dcPayableInfo := ParticipantPaymentInfo{"Dc001", 29.25, "Orddc001", "", "pending"}
		dFarmPayableInfo := ParticipantPaymentInfo{"002", 25.25, "", "", "pending"}
		dFarmTarnsportPayableInfo := ParticipantPaymentInfo{"001", 20.29, "", "", "pending"}
		payableinfo = append(payableinfo, farmerPayableInfo)
		payableinfo = append(payableinfo, casPayableInfo)
		payableinfo = append(payableinfo, pcPayableInfo)
		payableinfo = append(payableinfo, dcPayableInfo)
		payableinfo = append(payableinfo, dFarmTarnsportPayableInfo)
		payableinfo = append(payableinfo, dFarmPayableInfo)
		PayableAsset.ParticipantPaymentInfos = payableinfo
		data01, _ := PaymentInvoiceAssettoJson(PayableAsset)
		var argTocreatePayableasset [][]byte
		argTocreatePayableasset = [][]byte{[]byte("CreatePaymentInvoiceAsset"), []byte(data01)}
		payload := stub.MockInvoke("000", argTocreatePayableasset)
		logger.Infof("Response : %s", string(payload.Payload))
		Expect(payload.Status).Should(Equal(status200))

		PayableAsset2 := PaymentInvoiceAsset{}
		PayableAsset2.BuyerOrderID = "001"
		PayableAsset2.DocType = PAYMENT_INVOICE_RECORD_TYPE
		PayableAsset2.ProduceID = "Pro001"
		PayableAsset2.Variety = "Gala"
		PayableAsset2.PaymentInvoiceID = "Inv002"
		PayableAsset2.Qty = 400
		PayableAsset2.Status = "Paid"
		var payableinfo2 []ParticipantPaymentInfo
		farmerPayableInfo2 := ParticipantPaymentInfo{"Far001", 200.25, "Ordfar001", "", "paid"}
		casPayableInfo2 := ParticipantPaymentInfo{"Cas001", 20.25, "Ordcas001", "", "paid"}
		pcPayableInfo2 := ParticipantPaymentInfo{"Pc001", 21.25, "Ordpc001", "", "paid"}
		dcPayableInfo2 := ParticipantPaymentInfo{"Dc001", 29.25, "Orddc001", "", "paid"}
		dFarmPayableInfo2 := ParticipantPaymentInfo{"002", 25.25, "", "", "paid"}
		dFarmTarnsportPayableInfo2 := ParticipantPaymentInfo{"001", 20.29, "", "", "paid"}
		payableinfo2 = append(payableinfo2, farmerPayableInfo2)
		payableinfo2 = append(payableinfo2, casPayableInfo2)
		payableinfo2 = append(payableinfo2, pcPayableInfo2)
		payableinfo2 = append(payableinfo2, dcPayableInfo2)
		payableinfo2 = append(payableinfo2, dFarmTarnsportPayableInfo2)
		payableinfo2 = append(payableinfo2, dFarmPayableInfo2)
		PayableAsset2.ParticipantPaymentInfos = payableinfo2
		data02, _ := PaymentInvoiceAssettoJson(PayableAsset2)
		argTocreatePayableasset = [][]byte{[]byte("CreatePaymentInvoiceAsset"), []byte(data02)}
		payload = stub.MockInvoke("000", argTocreatePayableasset)
		logger.Infof("Response : %s", string(payload.Payload))
		Expect(payload.Status).Should(Equal(status200))
	})
	Describe("Checking the Smart contract Running status and Function Call Error", func() {
		Context("Checking SmartContract running Status ", func() {
			It("Run CheckStatus Function and see if result is 200", func() {
				argToStatuscheck := [][]byte{[]byte("StatusCheck")}
				payload := stub.MockInvoke("000", argToStatuscheck)
				logger.Infof("Response : %s", string(payload.Payload))
				Expect(payload.Status).Should(Equal(status200))
			})
		})
		Context("Checking wrong function name logic ", func() {
			It("Check if called function exist", func() {
				argToStatuscheck := [][]byte{[]byte("TestFunction")}
				payload := stub.MockInvoke("000", argToStatuscheck)
				logger.Infof("Response : %s", string(payload.Payload))
				Expect(payload.Status).Should(Equal(int32(500)))
			})
		})
	})
})

*/
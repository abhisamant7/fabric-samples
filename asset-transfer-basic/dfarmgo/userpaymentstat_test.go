package main
/*

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unit Test User Payment Stat", func() {
	stub = shim.NewMockStub("Testing User Payment Stat", new(Dfarmsc))
	Describe("Checking stats for different Participant", func() {
		Context(" stat for Farmer with produce Id Pro001 , id far001  and type as 0", func() {
			It("Run GetPaymentStat Function and see if result is 200 for Farmer Stat", func() {
				request := UserPaymentStatRequest{"Pro001", "", 0, "Far001"}
				rdata, _ := json.Marshal(request)
				logger.Infof("Request : %s", string(rdata))
				argTostat := [][]byte{[]byte("GetPaymentStat"), []byte(rdata)}
				payload := stub.MockInvoke("000", argTostat)
				logger.Infof("Response : %s", string(payload.Payload))
				Expect(payload.Status).Should(Equal(status200))
			})
			It("Run GetPaymentStat Function and see if result is 200 for Farmer Statif farmerid is Far002", func() {
				request := UserPaymentStatRequest{"Pro001", "", 0, "Far002"}
				rdata, _ := json.Marshal(request)
				argTostat := [][]byte{[]byte("GetPaymentStat"), []byte(rdata)}
				payload := stub.MockInvoke("000", argTostat)
				logger.Infof("Response : %s", string(payload.Payload))
				Expect(payload.Status).Should(Equal(status200))
			})
		})
		Context("Stat for CAS ", func() {
			It("Check request with CAS id Cas001", func() {
				request := UserPaymentStatRequest{"Pro001", "", 1, "Cas001"}
				rdata, _ := json.Marshal(request)
				argTostat := [][]byte{[]byte("GetPaymentStat"), []byte(rdata)}
				payload := stub.MockInvoke("000", argTostat)
				logger.Infof("Response : %s", string(payload.Payload))
				Expect(payload.Status).Should(Equal(status200))
			})
		})
	})
})*/

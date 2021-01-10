package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*type BuyerOrderAsset struct {
	DocType          string            `json:"docTypemitempty"`
	OrderType        uint32            `json:"ORDERTYPEmitempty"`
	DCID             string            `json:"DCIDmitempty"`
	ProduceID        string            `json:"prIdmitempty"`
	Variety          string            `json:"varietymitempty"`
	Qty              uint64            `json:"totalQTYmitempty"`
	BuyerID          string            `json:"buyerIDmitempty"`
	TotalPrice       float64           `json:"pricemitempty"` // total price
	QtyBreakdowns    []TableVariety    `json:"QTYmitempty"`
	ParticipantInfos []ParticipantInfo `json:"PARTICIPANTINFOSmitempty"`
	OrderID          string            `json:"ORDERIDmitempty"`
	Status           string            `json:"STATUSmitempty"`
	//ParentID         string            `json:"parentorderIdmitempty"` // NEED To discuss with prasanna
	TransportRate  Rate   `json:"TRANSPORTRATEmitempty"`
	PayableAssetID string `json:"payableAssetIDmitempty"`
}*/

/*
* TestInvokePricerateAsset simulates an TableVariteyPerRateAsset transaction on the This cahincode
 */
func TestCreatebuyerOrderandGeneratePayment(t *testing.T) {
	logger.Infof("Entering TestCreatebuyerOrderandGeneratePayment")

	// Instantiate mockStub using this chaincode as the target chaincode to unit test
	stub := shim.NewMockStub("mockStub", new(Dfarmsc))
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	var partinfo []ParticipantInfo

	farmer := ParticipantInfo{FARMER, "Far001", "Order001", "US", "NC"}
	partinfo = append(partinfo, farmer)
	pc := ParticipantInfo{PC, "PC001", "Order002", "US", "NC"}
	partinfo = append(partinfo, pc)
	trans := ParticipantInfo{PCDCTRANS, "PCDC001", "Order004", "US", "NC"}
	partinfo = append(partinfo, trans)
	dc := ParticipantInfo{DC, "DC001", "Order003", "US", "NC"}
	partinfo = append(partinfo, dc)
	dfarm := ParticipantInfo{DFARM, "dfarm001", "", "US", "NC"}
	partinfo = append(partinfo, dfarm)
	var qtbreakdown []TableVariety
	greadA := TableVariety{"Gread A", 100.00, Unit{"KG", 0.00}}
	greadB := TableVariety{"Gread B", 100.00, Unit{"KG", 0.00}}
	greadC := TableVariety{"Gread C", 100.00, Unit{"KG", 0.00}}
	qtbreakdown = append(qtbreakdown, greadA, greadB, greadC)

	border := BuyerOrderAsset{}
	border.BuyerID = "order0034"
	border.DCID = "DC001"
	border.OrderID = "0rder012"
	border.OrderType = 9
	border.ParticipantInfos = partinfo
	border.ProduceID = "tomato003"
	border.Qty = 300
	border.QtyBreakdowns = qtbreakdown
	border.Status = "delivered"
	border.TotalPrice = 234.45
	//border.TransportRate = Rate{"KG", "INR", 2.43}
	border.Variety = "kola"

	data, _ := BuyerOrderAssettoJson(border)
	logger.Infof("BuyerOrder is %s : ", data)

	//get Paybale Asset:

	/*// Here we perform a "mock invoke" to invoke the function  method with associated parameters
	// The first parameter is the function we are invoking
	/*type QualityPerInfo struct {
		Type     string  `json:"TYPEmitempty"`
		OpMode   bool    `json:"OPMODEmitempty"` //0 means add/more 1 means less/subscract
		PerValue float32 `json:"PERVALUEmitempty"`
	}*/
	/*type TableVarietyPerRateAsset struct {
		DocType         string           `json:"docTypemitempty"`
		ProduceName     string           `json:"PRODUCEmitempty"`
		State           string           `json:"STATEmitempty"`
		Country         string           `json:"COUNTRYmitempty"`
		Variety         string           `json:"VARIETYmitempty"`
		QualityPerInfos []QualityPerInfo `json:"QUALITYPERINFOSmitempty"`
	}
	qualityinfo := QualityPerInfo{"Grade A", false, .25}
	qualityinfos = append(qualityinfos, qualityinfo)
	perrate := TableVarietyPerRateAsset{"", "Tomato", "NC", "USA", "Dankal", qualityinfos}
	data, _ := TableVarietyPerRateAssettoJson(perrate)
	logger.Infof("TestInvokePricerateAsset: Data is %s", string(data))
	result := stub.MockInvoke("001",
		[][]byte{[]byte("CreateTableVarietyPerRateAsset"),
			[]byte(data)})

	// We expect a shim.ok if all goes well
	if result.Status != shim.OK {
		t.Fatalf("Expected unauthorized user error to be returned")
	}
	logger.Infof("TestInvokePricerateAsset: Create result is %s", result)
	key := TestKey{"Dankal", "Tomato", "USA", "NC"}
	data, _ = json.Marshal(key)
	logger.Infof("TestInvokePricerateAsset: Data is %s", data)
	result = stub.MockInvoke("001",
		[][]byte{[]byte("QueryTableVarietyPerRateAsset"),
			[]byte(data)})

	// We expect a shim.ok if all goes well
	if result.Status != shim.OK {
		t.Fatalf("Expected unauthorized user error to be returned")
	}
	logger.Infof("TestInvokePricerateAsset: Query Result  is %s", result)*/

}

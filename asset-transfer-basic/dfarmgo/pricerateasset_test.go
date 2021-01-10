package main

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type TestKey struct {
	Variety     string `json:"VARIETYmitempty"` //primery key
	ProduceName string `json:"PRODUCEmitempty"`
	Country     string `json:"COUNTRYmitempty"`
	State       string `json:"STATEmitempty"`
}

/*
* TestInvokePricerateAsset simulates an TableVariteyPerRateAsset transaction on the This cahincode
 */
func TestInvokePricerateAsset(t *testing.T) {
	logger.Infof("Entering TestInvokePricerateAsset")

	// Instantiate mockStub using this chaincode as the target chaincode to unit test
	stub := shim.NewMockStub("mockStub", new(Dfarmsc))
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	var qualityinfos []QualityPerInfo

	// Here we perform a "mock invoke" to invoke the function  method with associated parameters
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
	}*/
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
	logger.Infof("TestInvokePricerateAsset: Query Result  is %s", result)

}

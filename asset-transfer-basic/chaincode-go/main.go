/*This is main program for Test  project  and it get called for chaincode Instantiation*/

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ============================================================================================================================
// Main function for
// ============================================================================================================================
var logger = shim.NewLogger("dFarmsc")


func main() {
	logger.SetLevel(shim.LogDebug)
	err := shim.Start(new(Dfarmsc))
	if err != nil {
		logger.Infof("Error starting dFarm Smart Contract - %s", err)
	}
}

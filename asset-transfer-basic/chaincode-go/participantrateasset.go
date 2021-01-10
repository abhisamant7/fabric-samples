package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type ParticipantRate struct {
	DeliveryRateVal DeliveryRate
	DCpercentage    float32
	Dfarmpercentage DfarmPer
}

func getParticipantRateAsset(stub shim.ChaincodeStubInterface, produceName string, participantinfo map[int]ParticipantInfo) (ParticipantRate, bool) {

	var dckeys, dfarmkey, transkey []string
	var partrate ParticipantRate

	if len(produceName) == 0 {
		logger.Errorf("getParticipantRateAsset : Produce name is empty.")
		return ParticipantRate{}, false
	}

	if dcpartinfo, ok := participantinfo[DC]; ok {
		logger.Infof("getParticipantRateAsset : DC Country is : %s ", dcpartinfo.Country)
		logger.Infof("getParticipantRateAsset : DC State is : %s ", dcpartinfo.State)
		logger.Infof("getParticipantRateAsset :DC Produce Name is : %s ", produceName)

		dckeys = append(dckeys, dcpartinfo.Country) //Country
		dckeys = append(dckeys, dcpartinfo.State)   //State
		dckeys = append(dckeys, produceName)        //Produce Name
		//get DC rate
		dcrate, dcok := getDCRateAsset(stub, dckeys)
		if dcok {
			partrate.DCpercentage = dcrate
		}
	}
	//getdfarm Rate
	if dfarmpartinfo, ok := participantinfo[DFARM]; ok {
		logger.Infof("getParticipantRateAsset :dFarm Country is : %s ", dfarmpartinfo.Country)
		logger.Infof("getParticipantRateAsset :dFarm State is : %s ", dfarmpartinfo.State)
		logger.Infof("getParticipantRateAsset :dFarm Produce Name is : %s ", produceName)

		dfarmkey = append(dfarmkey, dfarmpartinfo.Country) //Country
		dfarmkey = append(dfarmkey, dfarmpartinfo.State)   //State
		dfarmkey = append(dfarmkey, produceName)           //Produce Name
		ddfarmper, dfarmok := getdFarmRateAsset(stub, dfarmkey)
		if dfarmok {
			partrate.Dfarmpercentage = ddfarmper
		}
	}
	//get Transport Rate
	if partinfo, ok := participantinfo[PCDCTRANS]; ok {
		logger.Infof("getParticipantRateAsset :Transport Country is : %s ", partinfo.Country)
		logger.Infof("getParticipantRateAsset :Transport State is : %s ", partinfo.State)
		logger.Infof("getParticipantRateAsset :Transport Produce Name is : %s ", produceName)

		transkey = append(transkey, partinfo.Country) //Country
		transkey = append(transkey, partinfo.State)   //State
		transkey = append(transkey, produceName)      //Produce Name
		dtrans, transok := getDeliveryRatAsset(stub, transkey)
		if transok {
			partrate.DeliveryRateVal = dtrans
		}
	}
	logger.Infof("getParticipantRateAsset:Participant rates are %v ", partrate)
	return partrate, true
}

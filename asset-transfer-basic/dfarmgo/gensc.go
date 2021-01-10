package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//GenAsset : This go structure  will be used in all request and response all operation of generic data in Fabric Ledger
//No Attributs based check will be performed in Smart contract execution logic
type GenAsset struct {
	AssetName   string        `json:"assetName"`             //Name of Asset i.e IOT Device Info
	Keys        []string      `json:"keys"`                  //Primary Key to add this Asset to Ledger
	EntityCount int32         `json:"entityCount,omitempty"` //No. of Asset data
	AssetDatas  []interface{} `json:"assetDatas,omitempty"`  //List of asset data
	QueryString string        `json:"queryString,omitempty"` //Generic Couch query string ,this value must be not set if ledger world state is stateDB
	Bookmark    string        `json:"bookmark,omitempty"`
}

// GenAssetResult :This go structure  will contain information about result of Asset creation request
type GenAssetResult struct {
	Keys   []string `json:"keys,omitempty"`
	Result string   `json:"result,omitempty"`
}

//JsontoGenAsset Convert JSON   to GenAsset Asset
func JsontoGenAsset(data []byte) (GenAsset, error) {
	obj := GenAsset{}
	if data == nil {
		return obj, fmt.Errorf("Input data  for json to GenAsset is missing")
	}

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

//GenAssettoJSON Convert GenAsset  object to Json Message
func GenAssettoJSON(obj GenAsset) ([]byte, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, err
}

/*CreateGenAssets Function will  insert record in ledger  based on primary keys after receiving request from Client Application*/
func CreateGenAssets(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte

	if len(args) < 1 {
		logger.Errorf("CreateGenAssets : Incorrect number of arguments, need to add one argument to this function call.")
		return shim.Error("CreateGenAssets : Incorrect number of arguments, need to add one argument to this function call.")
	}

	asset := GenAsset{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	if err != nil {
		logger.Errorf("CreateGenAssets : Error Parsing Request data   as %s", err)
		return shim.Error(fmt.Sprintf("CreateGenAssets :  Error Parsing request data  as %s", err))

	}
	logger.Infof("CreateGenAssets :AssetName  : %s ", asset.AssetName)
	logger.Infof("CreateGenAssets :Asset Keys are : %v ", asset.Keys)
	logger.Infof("CreateGenAssets :No of asset in this request is : %d", asset.EntityCount)
	dataList := asset.AssetDatas //This variable will have collection of asset user wants to add to the list
	logger.Infof("CreateGenAssets :No of asset in this request is : %d", len(dataList))
	var resultList []GenAssetResult

	//Added this code to avoid any core dump due to wrong value of asset.EntityCount
	if len(dataList) < int(asset.EntityCount) {
		asset.EntityCount = int32(len(dataList))

	}
	for assetcount := int32(0); assetcount < asset.EntityCount; assetcount++ {
		var keys []string
		data := dataList[assetcount]
		var result = make(map[string]interface{})
		logger.Infof("CreateGenAssets :data is  : %v", data)
		bytedata, _ := json.Marshal(data)
		err = json.Unmarshal(bytedata, &result)
		for keycount := 0; keycount < len(asset.Keys); keycount++ {
			key := result[asset.Keys[keycount]].(string)
			keys = append(keys, key)
			err = CreateAsset(stub, asset.AssetName, keys, bytedata)
			var result GenAssetResult
			if err != nil {
				logger.Errorf("CreateGenAssets : Error inserting Object first time  into LedgerState %s", err)
				result.Keys = keys
				result.Result = fmt.Sprintf("Error inserting Object %s", err)
			} else {

				result.Keys = keys
				result.Result = "Data inserted successfully in the ledger"
			}
			resultList = append(resultList, result)

		}

	}
	jsonresult, _ := json.Marshal(resultList)
	asset.AssetDatas[0] = string(jsonresult)
	Avalbytes, _ = json.Marshal(asset)
	return shim.Success([]byte(Avalbytes))
}

//ListGenAssets  Function will  query  record in ledger based on generic query or based on only Asset Name
func ListGenAssets(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var keys []string
	var dataItr shim.StateQueryIteratorInterface
	var resmetadata *pb.QueryResponseMetadata = nil
	if len(args) < 1 {
		logger.Errorf("ListGenAssets : Incorrect number of arguments, need to add one argument to this function call.")
		return shim.Error("ListGenAssets : Incorrect number of arguments, need to add one argument to this function call.")
	}

	asset := GenAsset{}
	err = json.Unmarshal([]byte(args[0]), &asset)
	if err != nil {
		logger.Errorf("ListGenAssets : Error Parsing Request data   as %s", err)
		return shim.Error(fmt.Sprintf("ListGenAssets :  Error Parsing request data  as %s", err))

	}
	logger.Infof("ListGenAssets :AssetName  : %s ", asset.AssetName)
	logger.Infof("ListGenAssets :Asset Keys are : %v ", asset.Keys)
	logger.Infof("ListGenAssets :No of asset in this request is : %vd", asset.EntityCount)
	if len(asset.QueryString) != 0 {
		if asset.EntityCount == 0 {
			dataItr, err = GenericQueryAsset(stub, asset.QueryString)
			if err != nil {
				logger.Errorf("ListGenAssets : No data found in the I ledger for query string %s", asset.QueryString)
				return shim.Error(fmt.Sprintf("ListGenAssets :  Error  %s found with query string %s", err, asset.QueryString))
			}

			buffer, berr := constructQueryResponseFromIterator(dataItr)
			if berr != nil {
				logger.Errorf("ListGenAssets : Error:  %s  with generation of result string", berr)
				return shim.Error(fmt.Sprintf("ListGenAssets :  Error:  %s found when generating result string", berr))
			}
			asset.AssetDatas = append(asset.AssetDatas, buffer.String())
		} else {
			dataItr, resmetadata, err = GenericQueryAssetwithPeginations(stub, asset.QueryString, asset.EntityCount, asset.Bookmark)
			if err != nil {
				logger.Errorf("ListGenAssets : No data found in the I ledger for query string %s", asset.QueryString)
				return shim.Error(fmt.Sprintf("ListGenAssets :  Error  %s found with query string %s", err, asset.QueryString))
			}

			buffer, berr := constructQueryResponseFromIterator(dataItr)
			if berr != nil {
				logger.Errorf("ListGenAssets : Error:  %s  with generation of result string", berr)
				return shim.Error(fmt.Sprintf("ListGenAssets :  Error:  %s found when generating result string", berr))
			}
			bufferwithpagination := addPaginationMetadataToQueryResults(buffer, resmetadata)
			asset.AssetDatas = append(asset.AssetDatas, bufferwithpagination.String())
		}

	} else {
		dataItr, err = ListAllAsset(stub, asset.AssetName, keys)
		if err != nil {
			logger.Errorf("ListGenAssets : instance not found in ledger")
			return shim.Error("ListGenAssets : instance not found in ledger")

		}

		buffer, berr := constructQueryResponseFromIterator(dataItr)
		if berr != nil {
			logger.Errorf("ListGenAssets : Error:  %s  with generation of result string", berr)
			return shim.Error(fmt.Sprintf("ListGenAssets :  Error:  %s found when generating result string", berr))
		}
		asset.AssetDatas = append(asset.AssetDatas, buffer.String())
	}

	Avalbytes, _ = json.Marshal(asset)
	if err != nil {
		logger.Errorf("ListGenAssets : Cannot Marshal result set. Error : %v", err)
		return shim.Error(fmt.Sprintf("ListGenAssets: Cannot Marshal result set. Error : %v", err))
	}
	return shim.Success([]byte(Avalbytes))

}

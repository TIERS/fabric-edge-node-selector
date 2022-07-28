package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/dmonteroh/distributed-resources-smartcontract/inventory-sc/internal"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []internal.Asset{
		// {ID: "localhost", Name: "localhost-laptop", Owner: "Org1", Type: 0, State: 1, Properties: map[string]string{"GPU": "true"}},
		// {ID: "172.26.45.114", Name: "172.26.45.114-jetson", Owner: "Org1", Type: 0, State: 1, Properties: map[string]string{"GPU": "true"}},
	}

	for _, asset := range assets {
		tmpJson := []byte(asset.String())
		err := ctx.GetStub().PutState(asset.ID, tmpJson)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetKey string) (internal.Asset, error) {
	statJSON, err := ctx.GetStub().GetState(assetKey)
	if err != nil {
		return internal.Asset{}, fmt.Errorf("failed to read from world state: %v", err)
	}
	if statJSON == nil {
		return internal.Asset{}, fmt.Errorf("the Asset with key: %s does not exist", assetKey)
	}

	var asset internal.Asset
	err = json.Unmarshal(statJSON, &asset)
	if err != nil {
		return internal.Asset{}, err
	}

	return asset, nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetJson string) error {
	asset, err := internal.JsonToAsset(assetJson)
	if err != nil {
		return err
	}
	exists, err := s.AssetExists(ctx, asset.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Asset with key: %s already exists", asset.ID)
	}

	// RUN VALIDATIONS
	validJson := []byte(asset.String())

	return ctx.GetStub().PutState(asset.ID, validJson)
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, assetJson string) error {
	asset, err := internal.JsonToAsset(assetJson)
	if err != nil {
		return err
	}
	exists, err := s.AssetExists(ctx, asset.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Asset with key: %s does not exist", asset.ID)
	}

	// RUN VALIDATIONS
	validJson := []byte(asset.String())

	return ctx.GetStub().PutState(asset.ID, validJson)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, assetKey string) error {
	exists, err := s.AssetExists(ctx, assetKey)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Asset with key: %s does not exist", assetKey)
	}

	return ctx.GetStub().DelState(assetKey)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, assetKey string) (bool, error) {
	asset, err := ctx.GetStub().GetState(assetKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return asset != nil, nil
}

// CURRENTLY IN TO-DO (NO OWNERSHIP REQUIRED ATM)
// TransferAsset updates the owner field of asset with given id in world state.
// func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, statIP string, newStatIP string) (string, error) {
// 	statObject, err := s.ReadAsset(ctx, statIP)
// 	if err != nil {
// 		return "", err
// 	}

// 	statObject.ID = newStatIP

// 	return statObject.String(), ctx.GetStub().PutState(newStatIP, []byte(statObject.String()))
// }

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]internal.Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return iteratorSlicer(resultsIterator)
}

// https://stackoverflow.com/questions/66685696/couchdb-mango-query-match-any-key-with-array-item
// GENERATE VIEW TO BETTER SEARCH PROPERTIES INSIDE INVENTORY ASSETS

func (s *SmartContract) GetServerAssets(ctx contractapi.TransactionContextInterface) ([]internal.Asset, error) {
	assetQuery := `{"selector":{"type":0,"state":1}}`
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetServerGPUAssets(ctx contractapi.TransactionContextInterface) ([]internal.Asset, error) {
	assetQuery := `{"selector":{"type":0,"state":1,"properties.gpu":1}}`
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetServerAssetsExceptId(ctx contractapi.TransactionContextInterface, excludeId string) ([]internal.Asset, error) {
	assetQuery := fmt.Sprintf(`{"selector":{"type":0,"state":1,"$not":{"id":"%s"}}}`, excludeId)
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetRobotAssets(ctx contractapi.TransactionContextInterface) ([]internal.Asset, error) {
	assetQuery := `{"selector":{"type":1,"state":1}}`
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetRobotAssetsExceptId(ctx contractapi.TransactionContextInterface, excludeId string) ([]internal.Asset, error) {
	assetQuery := fmt.Sprintf(`{"selector":{"type":1,"state":1,"$not":{"id":"%s"}}}`, excludeId)
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetSensorAssets(ctx contractapi.TransactionContextInterface) ([]internal.Asset, error) {
	assetQuery := `{"selector":{"type":2,"state":1}}`
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetSensorAssetsExceptId(ctx contractapi.TransactionContextInterface, excludeId string) ([]internal.Asset, error) {
	assetQuery := fmt.Sprintf(`{"selector":{"type":2,"state":1,"$not":{"id":"%s"}}}`, excludeId)
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetSensorAndRobotAssets(ctx contractapi.TransactionContextInterface) ([]internal.Asset, error) {
	assetQuery := `{"selector":{"type":{"$in":[1,2]},"state":1}}`
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetSensorAndRobotAssetsExceptId(ctx contractapi.TransactionContextInterface, excludeId string) ([]internal.Asset, error) {
	assetQuery := fmt.Sprintf(`{"selector":{"type":{"$in":[1,2]},"state":1,"$not":{"id":"%s"}}}`, excludeId)
	return stringQuery(ctx, assetQuery)
}

func stringQuery(ctx contractapi.TransactionContextInterface, queryString string) ([]internal.Asset, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return iteratorSlicer(resultsIterator)
}

func iteratorSlicer(resultsIterator shim.StateQueryIteratorInterface) ([]internal.Asset, error) {
	var assets []internal.Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		fmt.Println(queryResponse.Value)
		fmt.Println(string(queryResponse.Value))
		asset, err := internal.JsonToAsset(string(queryResponse.Value))
		if err != nil {
			return nil, err
		}
		fmt.Println(asset)
		assets = append(assets, asset)
	}
	return assets, nil
}

// Function for testing CouchDB Queries
func (s *SmartContract) ExecuteQuery(ctx contractapi.TransactionContextInterface, assetQuery string) ([]string, error) {
	var result []string
	resultsIterator, err := ctx.GetStub().GetQueryResult(assetQuery)
	if err != nil {
		return result, err
	}
	defer resultsIterator.Close()
	if resultsIterator.HasNext() {
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return nil, err
			}
			result = append(result, string(queryResponse.Value))
		}
	} else {
		return nil, fmt.Errorf("failed to query chaincode. No results found for iterator")
	}

	return result, nil

}

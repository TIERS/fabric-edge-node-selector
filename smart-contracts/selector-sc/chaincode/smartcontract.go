package chaincode

import (
	"fmt"
	"sort"

	"github.com/dmonteroh/distributed-resources-smartcontract/selector-sc/internal"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	stats := []internal.StoredSelection{}

	for _, stat := range stats {
		err := ctx.GetStub().PutState(stat.ID, []byte(stat.String()))
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetKey string) (internal.StoredSelection, error) {
	statJSON, err := ctx.GetStub().GetState(assetKey)
	if err != nil {
		return internal.StoredSelection{}, fmt.Errorf("failed to read from world state: %v", err)
	}
	if statJSON == nil {
		return internal.StoredSelection{}, fmt.Errorf("the Asset with key: %s does not exist", assetKey)
	}

	asset, err := internal.JsonToStoredSelection(string(statJSON))
	if err != nil {
		return internal.StoredSelection{}, fmt.Errorf("failed to read from world state: %v", err)
	}

	return asset, nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetJson string) error {
	asset, err := internal.JsonToStoredSelection(assetJson)
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
	asset, err := internal.JsonToStoredSelection(assetJson)
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

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]internal.StoredSelection, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return iteratorSlicer(resultsIterator)
}

func (s *SmartContract) GetAllSelectionTarget(ctx contractapi.TransactionContextInterface, asset string) ([]internal.StoredSelection, error) {
	assetQuery := fmt.Sprintf(`{"selector": {"target": "%s"}}`, asset)
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetAllSelectionServer(ctx contractapi.TransactionContextInterface, asset string) ([]internal.StoredSelection, error) {
	assetQuery := fmt.Sprintf(`{"selector": {"assetID": "%s"}}`, asset)
	return stringQuery(ctx, assetQuery)
}

// Inernal Functions
func iteratorSlicer(resultsIterator shim.StateQueryIteratorInterface) ([]internal.StoredSelection, error) {
	var assets []internal.StoredSelection
	if resultsIterator.HasNext() {
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return nil, err
			}
			asset, err := internal.JsonToStoredSelection(string(queryResponse.Value))
			if err != nil {
				return nil, err
			}
			assets = append(assets, asset)
		}
	} else {
		return nil, fmt.Errorf("failed to query chaincode. No results found for iterator")
	}

	sort.SliceStable(assets, func(i, j int) bool {
		return assets[i].Timestamp.TimeSeconds > assets[j].Timestamp.TimeSeconds
	})

	return assets, nil
}

func stringQuery(ctx contractapi.TransactionContextInterface, queryString string) ([]internal.StoredSelection, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return iteratorSlicer(resultsIterator)
}

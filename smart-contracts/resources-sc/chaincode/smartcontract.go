package chaincode

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/dmonteroh/distributed-resources-smartcontract/resources-sc/internal"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	stats := []internal.StoredStat{}

	for _, stat := range stats {
		err := ctx.GetStub().PutState(stat.DrcHost.HostID, []byte(stat.String()))
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, statIP string, statJSON string) error {
	exists, err := s.AssetExists(ctx, statIP)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Stats for %s already exists", statIP)
	}

	toStore, err := internal.JsonToStoredStat(statJSON)
	if err != nil {
		return err
	}
	// toStore := internal.ConvertToStorage(tmpStat)
	// toStore.ID = statIP
	// RUN VALIDATION

	return ctx.GetStub().PutState(statIP, []byte(toStore.String()))
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, statIP string) (*internal.StoredStat, error) {
	statJSON, err := ctx.GetStub().GetState(statIP)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if statJSON == nil {
		return nil, fmt.Errorf("the Stats for %s do not exist", statIP)
	}

	var stat internal.StoredStat
	err = json.Unmarshal(statJSON, &stat)
	if err != nil {
		return nil, err
	}

	return &stat, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, statIP string, statJSON string) error {
	exists, err := s.AssetExists(ctx, statIP)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Stats for %s do not exist", statIP)
	}

	tmpStat, err := internal.DrcJsonToStruct(statJSON)
	if err != nil {
		return err
	}
	toStore := internal.ConvertToStorage(tmpStat)
	toStore.ID = statIP

	return ctx.GetStub().PutState(statIP, []byte(toStore.String()))
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, statIP string) error {
	exists, err := s.AssetExists(ctx, statIP)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Stats for %s do not exist", statIP)
	}

	return ctx.GetStub().DelState(statIP)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, statIP string) (bool, error) {
	statJSON, err := ctx.GetStub().GetState(statIP)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return statJSON != nil, nil
}

func iteratorSlicer(resultsIterator shim.StateQueryIteratorInterface) ([]internal.StoredStat, error) {
	var assets []internal.StoredStat
	if resultsIterator.HasNext() {
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return nil, err
			}
			asset, err := internal.JsonToStoredStat(string(queryResponse.Value))
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

func stringQuery(ctx contractapi.TransactionContextInterface, queryString string) ([]internal.StoredStat, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return iteratorSlicer(resultsIterator)
}

func (s *SmartContract) GetAssetResource(ctx contractapi.TransactionContextInterface, hostname string) ([]internal.StoredStat, error) {
	assetQuery := fmt.Sprintf(`{"selector": {"hostname": "%s"}}`, hostname)
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetAssetResourceListTime(ctx contractapi.TransactionContextInterface, hostname string, minutes int) ([]internal.StoredStat, error) {
	timeStart := time.Now()
	timeEnd := timeStart.Add(time.Duration(-time.Duration(minutes) * time.Minute))
	assetQuery := fmt.Sprintf(`{"selector": {"hostname": "%s","timestamp.timeSeconds": {"$lt": %d,"$gte": %d}}}`, hostname, timeStart.Unix(), timeEnd.Unix())
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetLastResourceSummary(ctx contractapi.TransactionContextInterface, hostname string) (internal.StatSummary, error) {
	var statSummary internal.StatSummary
	assetQuery := fmt.Sprintf(`{"selector": {"hostname": "%s"}, "sort": [{"timestamp.timeSeconds": "desc"}], "limit": 2,"use_index": "resource_index"}`, hostname)
	queryResult, err := stringQuery(ctx, assetQuery)
	if err != nil {
		return statSummary, err
	}

	statSummary = internal.SummarizeStoredStat(queryResult[0])

	return statSummary, nil
}

func (s *SmartContract) GetSummaryAnalysisTime(ctx contractapi.TransactionContextInterface, hostname string, minutes int) (internal.StatAnalysis, error) {
	var statAnalysis internal.StatAnalysis
	storedStatList, err := s.GetAssetResourceListTime(ctx, hostname, minutes)
	if err != nil {
		return statAnalysis, err
	}
	var statSummarySlice []internal.StatSummary
	for _, stat := range storedStatList {
		var statSummary = internal.SummarizeStoredStat(stat)
		statSummarySlice = append(statSummarySlice, statSummary)
	}

	statAnalysis.Hostname = hostname
	statAnalysis.Duration = minutes
	statAnalysis.StatSummary = statSummarySlice
	statAnalysis = internal.AnalizeStatSummary(statAnalysis)

	return statAnalysis, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*internal.StoredStat, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var statObjects []*internal.StoredStat
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var statObject internal.StoredStat
		err = json.Unmarshal(queryResponse.Value, &statObject)
		if err != nil {
			return nil, err
		}
		statObjects = append(statObjects, &statObject)
	}

	sort.SliceStable(statObjects, func(i, j int) bool {
		return statObjects[i].Timestamp.TimeSeconds > statObjects[j].Timestamp.TimeSeconds
	})

	return statObjects, nil
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

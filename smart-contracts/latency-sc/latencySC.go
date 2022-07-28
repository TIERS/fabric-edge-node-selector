/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/dmonteroh/distributed-resources-smartcontract/latency-sc/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating resources-sc chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting resources-sc chaincode: %v", err)
	}
}

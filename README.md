# Hyperledger Fabric Smart Contracts and Application for Edge Node selection

**Distributed Ledger Technologies for Managing Heterogenous Computing and Sensing Systems at the Edge**. Daniel Andres Montero Hernandez and Jorge Peña Queralta and Tomi Westerlund.

The increased popularity of Internet of Things (IoT) devices, ranging from simple sensors to powerful embedded computers, has created the need for solutions capable of processing and storing information near those assets. Edge Computing (EC) has become a staple architecture when designing solutions for IoT, as it optimizes the workload and capacity of systems dependent of the Cloud, by placing the required computing power near to where the information is being produced and consumed. An issue with these solutions, is that reaching consensus regarding the state of the network becomes more challenging as they scale in size. Distributed Ledger Technology (DLT) can be described as a network of distributed databases that incorporate cryptography and algorithms to reach consensus among the participants. DLT has gained traction over the past years, particularly due to the popularity of Blockchain, the most well-known type of DLT implementation. In addition to the capability of reaching consensus, another key concept that brings EC and DLT together, is the reliability and trust that the latter offers through transparent and traceable transactions. In this thesis, we present the design and development of a proof-of-concept system that uses DLT Smart Contracts (SC) as the core for efficiently selecting Edge Nodes for offloading services. We present the experiments conducted to demonstrate the efficacy of the system and our conclusions regarding the usage of Hyperledger Fabric for managing systems at the edge.

## Installation

**For Go**: Please follow installation instructions of [Download and install - The Go Programming Language](https://go.dev/doc/install)

**For Smart Contracts and Gateway Application**: Please follow installation instructions of [Hyperledger Fabric and the Test Network](https://hyperledger-fabric.readthedocs.io/en/latest/getting_started.html)

**For Daemon**: Please follow installation instructions of [Install Docker Compose](https://docs.docker.com/compose/install/)

## Running System

Remember to Git Clone this repo!
```
git clone https://github.com/TIERS/fabric-edge-node-selector.git
```

In order to install the Smart Contracts into the Test Network, make sure to start the network and start the deployment using the following commands:

```
cd ~/hyperledger/fabric-samples/test-network/
./network.sh down 
./network.sh up createChannel -c mychannel -ca -s couchdb 
 
./network.sh deployCC -ccn resources-sc -ccp ~/fabric-edge-node-selector/smart-contracts/resources-sc -ccl go 
./network.sh deployCC -ccn inventory-sc -ccp ~/fabric-edge-node-selector/smart-contracts/inventory-sc -ccl go 
./network.sh deployCC -ccn latency-sc -ccp ~/fabric-edge-node-selector/smart-contracts/latency-sc -ccl go 
./network.sh deployCC -ccn selector-sc -ccp ~/fabric-edge-node-selector/smart-contracts/selector-sc -ccl go 
```

Once the Smart Contracts are installed, we will move the Application into the correct Hyperledger Folder:
```
cd ~/hyperledger/fabric-samples/asset-transfer-basic/application-go
rm -r *
cp -r ~/fabric-edge-node-selector/gateway-application/. ~/hyperledger/fabric-samples/asset-transfer-basic/application-go
go build distributedResources.go
./distributedResources
```

Note: Remember to delete the contents of the /wallet directory every time you use the ./network.sh up command

Before running the Daemon, remember to modify the .env file in the distributed-resource-collector folder to match your preferences.
The Daemon is executed by running the following command:
```
cd ~/fabric-edge-node-selector/distributed-resource-collector
docker-compose up
```
Some versions of docker-compose will complain about the version of the docker-compose.yaml, change from "3.8" to "3.7" if neccesary.

Finally, to interact with the system, install a REST Client such as Postman and import the Endpoints and Methods inside the postman-configuration folder.

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"bytes"
)


// ==== CHAINCODE RUN (CHAINCODE CONTAINER) ==================

// CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=scc:0 ./servicemarbles

// ==== IMPORT PACKAGE (CLI) ==================
// go get github.com/hyperledger/fabric/protos/ledger/queryresult

// ==== CHAINCODE INSTALLATION (CLI) ==================

// peer chaincode install -p chaincodedev/chaincode/servicemarbles -n scc -v 0

// ==== CHAINCODE INSTANTIATION (CLI) ==================

// peer chaincode instantiate -n scc -v 0 -c '{"Args":[]}' -C ch2

// ==== CHAINCODE EXECUTION SAMPLES (CLI) ==================

// ==== Invoke servicemarbles ====
// peer chaincode invoke -C ch2 -n scc -c '{"function": "helloWorld", "Args":[]}'
// ==== INITIALIZATION FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "initLedger", "Args":[]}'

// ==== GENERAL FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "read", "Args":["idagent1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "getValue", "Args":["idagent2"]}' -v 0
// peer chaincode invoke -C ch2 -n scc -c '{"function": "readEverything", "Args":[]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "allLedger", "Args":[]}'

// ==== CREATE ASSET FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "initService", "Args":["idservice5","service1","description1"]}
// peer chaincode invoke -C ch2 -n scc -c '{"function": "initAgent", "Args":["idagent10","agent10","address10"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "initServiceAgentRelation", "Args":["idservice1","idagent2","2","6","8"]}'

// ==== GET ASSET ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "getService", "Args":["idservice1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "getAgent", "Args":["idagent10"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "getServiceRelationAgent", "Args":["idservice1idagent1"]}'

// ==== GET HISTORY ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "getHistory", "Args":["idservice5"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "getServiceHistory", "Args":["idservice5"]}'


// ==== RANGE QUERY (USING COMPOSITE INDEX) ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "byService", "Args":["idservice1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "byAgent", "Args":["idAgent10"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "getAgentsByService", "Args":["idservice1"]}'

// ==== DELETE ASSET ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "deleteService", "Args":["idservice1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "deleteAgent", "Args":["idagent1"]}'

// ==== CALLS IN THE REAL PROJECT ====
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "helloWorld", "Args":[]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "initLedger", "Args":[]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "allLedger", "Args":[]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "getHistory", "Args":["service5"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "initAgent", "Args":["idagent10","agent10","address10"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "initService", "Args":["idservice10","service10","description10"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "getService", "Args":["idservice1"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "getAgent", "Args":["idagent1"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "getServiceRelationAgent", "Args":["idservice1idagent1"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "initServiceAgentRelation", "Args":["idservice1","idagent2","3","5","7"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "getAgentsByService", "Args":["idservice1"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "getService", "Args":["idservice5"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "deleteService", "Args":["idservice5"]}'



// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}

// Init initialize the chaincode
// The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
// Best practice is to have any Ledger initialization in separate function -- see initLedger()
//======================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	// Route to the appropriate handler function to interact with the ledger appropriately
	switch function {
	case "initLedger":
		return initLedger(stub)
	case "initService":
		return initService(stub, args)
	case "initAgent":
		return initAgent(stub, args)
	case "initServiceAgentRelation":
		// Already with reference integrity controls (service already exist, agent already exist, relation don't already exist)
		return initServiceAgentRelation(stub, args)
	case "getHistory":
		return getHistory(stub,args)
	case "getServiceHistory":
		return getServiceHistory(stub,args)
	case "getService":
		return queryService(stub,args)
	case "getAgent":
		return queryAgent(stub,args)
	case "getServiceRelationAgent":
		return queryServiceRelationAgent(stub,args)
	case "byService":
		return queryByServiceAgentRelation(stub,args)
	case "byAgent":
		return queryByAgentServiceRelation(stub,args)
	case "getAgentsByService":
		// also with only one record result return always a JSONArray
		return getServiceRelationAgentByServiceWithCostAndTime(stub,args)
	case "deleteService":
		return deleteService(stub,args)
	case "deleteAgent":
		return deleteAgent(stub,args)
	case "write":
		return  write(stub,args)
	case "read":
		return read(stub,args)
	case "readEverything":
		return readEverything(stub)
	case "allLedger":
		return readAllLedger(stub)
	case "getValue":
		return getValue(stub, args)
	case "helloWorld":
		fmt.Println("Ciao")
		// in := []byte(`{"Hello":"HelloWorld"}`)
		// var raw map[string]interface{}
		// json.Unmarshal(in, &raw)
		// out, _ := json.Marshal(raw)
		var buffer bytes.Buffer

		buffer.WriteString("[{\"Hello\":\"HelloWorld\"}]")

		return shim.Success(buffer.Bytes())
	default:
		return shim.Error("Invalid Smart Contract function Name.")
	}

	// error out
	fmt.Println("Received unknown invoke function Name - " + function)
	return shim.Error("Received unknown invoke function Name - '" + function + "'")
}

// ============================================================================================================================
// Query - legacy function
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call - Query()")
}
/*
Package main is the entry point of the hyperledger fabric chaincode and implements the shim.ChaincodeStubInterface
*/
/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
*/
package main

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	a "github.com/pavva91/assets"
	gen "github.com/pavva91/generalcc"
	// a "github.com/pavva91/trustreputationledger/assets"
	// gen "github.com/pavva91/trustreputationledger/generalcc"
	// in "github.com/pavva91/trustreputationledger/invokeapi"
	in "github.com/pavva91/invokeapi"
)

const(
	InitLedger    = "InitLedger"
	CreateService = "CreateService"
	CreateAgent   = "CreateAgent"
	CreateServiceAgentRelation = "CreateServiceAgentRelation"
	CreateServiceAndServiceAgentRelationWithStandardValue = "CreateServiceAndServiceAgentRelationWithStandardValue"
	CreateServiceAndServiceAgentRelation = "CreateServiceAndServiceAgentRelation"
	GetServiceHistory = "GetServiceHistory"
	GetServiceNotFoundError = "GetServiceNotFoundError"
	GetAgentNotFoundError = "GetAgentNotFoundError"
	GetServiceRelationAgent = "GetServiceRelationAgent"
	ByService = "byService"
	ByAgent = "byAgent"
	GetAgentsByService = "GetAgentsByService"
	GetServicesByAgent = "GetServicesByAgent"
	DeleteService = "DeleteService"
	DeleteAgent = "DeleteAgent"
	CreateActivity = "CreateActivity"
	GetActivity = "GetActivity"
	ByExecutedServiceTxId = "byExecutedServiceTxId"
	ByDemanderExecuter = "byDemanderExecuter"
	GetActivitiesByServiceTxId = "GetActivitiesByServiceTxId"
	GetActivitiesByDemanderExecuterTimestamp = "GetActivitiesByDemanderExecuterTimestamp"
	CreateReputation = "CreateReputation"
	ModifyReputationValue = "ModifyReputationValue"
	ModifyOrCreateReputationValue = "ModifyOrCreateReputationValue"
	GetReputationNotFoundError = "GetReputationNotFoundError"
	ByAgentServiceRole = "byAgentServiceRole"
	GetReputationsByAgentServiceRole = "GetReputationsByAgentServiceRole"
	Write = "Write"
	Read = "Read"
	ReadEverything = "ReadEverything"
	GetHistory = "GetHistory"
	GetReputationHistory = "GetReputationHistory"
	AllStateDB = "AllStateDB"
	GetValue = "GetValue"
	HelloWorld = "HelloWorld"

)



// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
	testMode bool
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = false
	err := shim.Start(simpleChaincode)
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}

// ============================================================================================================================
// Init - initialize the chaincode
// ============================================================================================================================
// The Init method is called when the Smart Contract "trustreputationledger" is instantiated by the blockchain network
// Best practice is to have any Ledger initialization in separate function -- see InitLedger()
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// TEST BEHAVIOUR
	if t.testMode {
		a.InitLedger(stub)
	}
	// NORMAL BEHAVIOUR
	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	switch function {
	// AGENT, SERVICE, AGENT SERVICE RELATION INVOKES

	// CREATE:
	case InitLedger:
		response := a.InitLedger(stub)
		return response
	case CreateService:
		return in.CreateService(stub, args)
	case CreateAgent:
		return in.CreateAgent(stub, args)
	case CreateServiceAgentRelation:
		// Already with reference integrity controls (service already exist, agent already exist, relation don't already exist)
		// Standard reputation value from inside
		return in.CreateServiceAgentRelation(stub, args)
	case CreateServiceAndServiceAgentRelationWithStandardValue:
		// If service doesn't exist it will create with a standard value of reputation defined inside the function
		return in.CreateServiceAndServiceAgentRelationWithStandardValue(stub, args)
	case CreateServiceAndServiceAgentRelation:
		// If service doesn't exist it will create
		return in.CreateServiceAndServiceAgentRelation(stub, args)

		// GET:
	case GetServiceHistory:
		return a.GetServiceHistory(stub, args)
	case GetServiceNotFoundError:
		return in.QueryService(stub, args)
	case GetAgentNotFoundError:
		return in.QueryAgent(stub, args)
	case GetServiceRelationAgent:
		return in.QueryServiceRelationAgent(stub, args)

		// RANGE QUERY:
	case ByService:
		return in.QueryByServiceAgentRelation(stub, args)
	case ByAgent:
		return in.QueryByAgentServiceRelation(stub, args)
	case GetAgentsByService:
		// also with only one record result return always a JSONArray
		return in.GetServiceRelationAgentByServiceWithCostAndTime(stub, args)
	case GetServicesByAgent:
		// also with only one record result return always a JSONArray
		return in.GetServiceRelationAgentByAgentWithCostAndTime(stub, args)

		// DELETE:
	case DeleteService:
		return a.DeleteService(stub, args)
	case DeleteAgent:
		return a.DeleteAgent(stub, args)

	// ACTIVITY INVOKES
	// CREATE:
	case CreateActivity:
		return in.CreateActivity(stub, args)
		// GET:
	case GetActivity:
		return in.QueryActivity(stub, args)
		// RANGE QUERY:
	case ByExecutedServiceTxId:
		return in.QueryByExecutedServiceTx(stub, args)
	case ByDemanderExecuter:
		return in.QueryByDemanderExecuter(stub, args)
	case GetActivitiesByServiceTxId:
		// also with only one record result return always a JSONArray
		return in.GetActivitiesByExecutedServiceTxId(stub, args)
	case GetActivitiesByDemanderExecuterTimestamp:
		// also with only one record result return always a JSONArray
		return in.GetActivitiesByDemanderExecuterTimestamp(stub, args)

	// REPUTATION INVOKES
	// CREATE:
	case CreateReputation:
		return in.CreateReputation(stub, args)
		// MODIFTY:
	case ModifyReputationValue:
		return in.ModifyReputationValue(stub, args)
	case ModifyOrCreateReputationValue:
		return in.ModifyOrCreateReputationValue(stub, args)

		// GET:
	case GetReputationNotFoundError:
		return in.QueryReputation(stub, args)
		// RANGE QUERY:
	case ByAgentServiceRole:
		return in.QueryByAgentServiceRole(stub, args)
	case GetReputationsByAgentServiceRole:
		// also with only one record result return always a JSONArray
		return in.GetReputationsByAgentServiceRole(stub, args)

		// GENERAL INVOKES
	case Write:
		return gen.Write(stub, args)
	case Read:
		return gen.Read(stub, args)
	case ReadEverything:
		return a.ReadEverything(stub)
	case GetHistory:
		// Get Chain Transaction Log of that assetId
		return gen.GetHistory(stub, args)
	case GetReputationHistory:
		return in.GetReputationHistory(stub, args)
	case AllStateDB:
		return gen.ReadAllStateDB(stub)
	case GetValue:
		return gen.GetValue(stub, args)
	case HelloWorld:
		fmt.Println("Ciao")
		var buffer bytes.Buffer
		buffer.WriteString("[{\"Hello\":\"HelloWorld\"}]")
		return shim.Success(buffer.Bytes())
	default:
		// Error Output
		fmt.Println("Received unknown in function Name - " + function)
		return shim.Error("Invalid Smart Contract function Name.")
	}


}

// ============================================================================================================================
// Query - legacy function
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call - Query()")
}

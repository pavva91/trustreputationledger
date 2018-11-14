/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
*/
package invokeapi

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pavva91/arglib"
	// a "github.com/pavva91/trustreputationledger/assets"
	a "github.com/pavva91/assets"

)

var serviceRelationAgentInvokeCallLog = shim.NewLogger("serviceRelationAgentInvokeCall")


// ========================================================================================================================
// Init Service Agent Relation - wrapper of CreateServiceAgentRelationAndReputation called from chiancode's Invoke
// ========================================================================================================================
func CreateServiceAgentRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1       2       3
	// "ServiceId", "AgentId", "Cost", "Time"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 4)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	// ==== Set Variables from Passed Arguments ====
	serviceId := args[0]
	agentId := args[1]
	cost := args[2]
	time := args[3]

	// ==== Check if already existing service ====
	service, errS := a.GetServiceNotFoundError(stub, serviceId)
	if errS != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to find service by id " + serviceId)
		return shim.Error("Failed to find service by id " + errS.Error())
	}
	serviceRelationAgentInvokeCallLog.Info("Service Already existing ok")

	// ==== Check if already existing agent ====
	agent, errA := a.GetAgentNotFoundError(stub, agentId)
	if errA != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to find agent by id " + agentId)
		return shim.Error("Failed to find agent by id: " + errA.Error())
	}
	serviceRelationAgentInvokeCallLog.Info("Agent Already existing ok")


	// ==== Check, Create, Indexing ServiceRelationAgent ====

	serviceRelationAgent, serviceRelationError := a.CheckingCreatingIndexingServiceRelationAgent(serviceId,agentId, cost, time, stub)
	if serviceRelationError != nil {
		return shim.Error("Error saving ServiceRelationAgent: " + serviceRelationError.Error())

	}

	// ==== ServiceRelationAgent saved and indexed. Set Event ====

	eventPayload:="Created Service RelationAgent: " + serviceId + " with agent: " + agentId
	payloadAsBytes := []byte(eventPayload)
	eventError := stub.SetEvent("ServiceRelationAgentCreatedEvent",payloadAsBytes)
	if eventError != nil {
		serviceRelationAgentInvokeCallLog.Info("Error in event Creation: " + eventError.Error())
	}else {
		serviceRelationAgentInvokeCallLog.Info("Event Create ServiceRelationAgent OK")
	}
	// ==== ServiceRelationAgent saved & indexed. Return success ====
	serviceRelationAgentInvokeCallLog.Info("Service: " + service.Name + " mapped with agent: " + agent.Name + " with cost: " + serviceRelationAgent.Cost + " and time: " + serviceRelationAgent.Time)
	return shim.Success(nil)
}

// ========================================================================================================================
// Init Service Agent Relation - wrapper of CreateServiceAgentRelationAndReputation called from chiancode's Invoke
// ========================================================================================================================
func CreateServiceAgentRelationAndReputation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1       2       3
	// "ServiceId", "AgentId", "Cost", "Time"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 4)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	// ==== Set Variables from Passed Arguments ====
	serviceId := args[0]
	agentId := args[1]
	cost := args[2]
	time := args[3]

	// ==== Check if already existing service ====
	service, errS := a.GetServiceNotFoundError(stub, serviceId)
	if errS != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to find service by id " + serviceId)
		return shim.Error("Failed to find service by id " + errS.Error())
	}
	serviceRelationAgentInvokeCallLog.Info("Service Already existing ok")

	// ==== Check if already existing agent ====
	agent, errA := a.GetAgentNotFoundError(stub, agentId)
	if errA != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to find agent by id " + agentId)
		return shim.Error("Failed to find agent by id: " + errA.Error())
	}
	serviceRelationAgentInvokeCallLog.Info("Agent Already existing ok")


	// ==== Check, Create, Indexing ServiceRelationAgent ====

	serviceRelationAgent, serviceRelationError := a.CheckingCreatingIndexingServiceRelationAgent(serviceId,agentId, cost, time, stub)
	if serviceRelationError != nil {
		return shim.Error("Error saving ServiceRelationAgent: " + serviceRelationError.Error())

	}

	// ==== Check, Create, Indexing Reputation ====
	initReputationValue := "6"
	reputation,reputationError := a.CheckingCreatingIndexingReputation(agentId,serviceId,a.Executer,initReputationValue,stub)
	if reputationError != nil {
		return shim.Error("Error saving Agent reputation: " + reputationError.Error())
	}

	// ==== ServiceRelationAgent saved and indexed. Set Event ====

	eventPayload:="Created Service RelationAgent: " + serviceId + " with agent: " + agentId
	payloadAsBytes := []byte(eventPayload)
	eventError := stub.SetEvent("ServiceRelationAgentAndReputationCreatedEvent",payloadAsBytes)
	if eventError != nil {
		serviceRelationAgentInvokeCallLog.Info("Error in event Creation: " + eventError.Error())
	}else {
		serviceRelationAgentInvokeCallLog.Info("Event Create ServiceRelationAgent OK")
	}
	// ==== ServiceRelationAgent saved & indexed. Return success ====
	serviceRelationAgentInvokeCallLog.Info("Service: " + service.Name + " mapped with agent: " + agent.Name + " with cost: " + serviceRelationAgent.Cost + " and time: " + serviceRelationAgent.Time + " nella relazione con reputazione iniziale: "+ reputation.Value)
	return shim.Success(nil)
}

// ========================================================================================================================
// Modify Service Relation Agent Cost - wrapper of ModifyServiceRelationAgentCost called from chiancode's Invoke
// ========================================================================================================================
func ModifyServiceRelationAgentCost(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0            1
	// "relationId", "newRelationCost"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 2)
	if argumentSizeError != nil {
		serviceRelationAgentInvokeCallLog.Error(argumentSizeError.Error())
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		serviceRelationAgentInvokeCallLog.Error(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	relationId := args[0]
	newRelationCost := args[1]

	// ==== get the serviceRelationAgent ====
	serviceRelationAgent, getError := a.GetServiceRelationAgentNotFoundError(stub, relationId)
	if getError != nil {
		serviceRelationAgentInvokeCallLog.Info(getError.Error())
		return shim.Error(getError.Error())
	}

	// ==== modify the serviceRelationAgent ====
	modifyError := a.ModifyServiceRelationAgentCost(serviceRelationAgent, newRelationCost, stub)
	if modifyError != nil {
		serviceRelationAgentInvokeCallLog.Error(modifyError.Error())
		return shim.Error(modifyError.Error())
	}

	// ==== ServiceRelationAgent Cost modified. Set Event ====

	eventPayload:="Modified Service RelationAgent: " + serviceRelationAgent.ServiceId + " with agent: " + serviceRelationAgent.AgentId + "from old cost value: " + serviceRelationAgent.Cost + "to new cost value: " + newRelationCost
	payloadAsBytes := []byte(eventPayload)
	eventError := stub.SetEvent("ServiceRelationAgentCostModifiedEvent",payloadAsBytes)
	if eventError != nil {
		serviceRelationAgentInvokeCallLog.Info("Error in event Creation: " + eventError.Error())
	}else {
		serviceRelationAgentInvokeCallLog.Info("Event Modifiy ServiceRelationAgent Cost OK")
	}

	serviceRelationAgentInvokeCallLog.Info("Modify Service RelationAgent Time OK")
	return shim.Success(nil)
}

// ========================================================================================================================
// Modify Service Relation Agent Time - wrapper of ModifyServiceRelationAgentTime called from chiancode's Invoke
// ========================================================================================================================
func ModifyServiceRelationAgentTime(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0            1
	// "relationId", "newRelationTime"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 2)
	if argumentSizeError != nil {
		serviceRelationAgentInvokeCallLog.Error(argumentSizeError.Error())
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		serviceRelationAgentInvokeCallLog.Error(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	relationId := args[0]
	newRelationTime := args[1]

	// ==== get the serviceRelationAgent ====
	serviceRelationAgent, getError := a.GetServiceRelationAgentNotFoundError(stub, relationId)
	if getError != nil {
		serviceRelationAgentInvokeCallLog.Error(getError.Error())
		return shim.Error(getError.Error())
	}

	// ==== modify the serviceRelationAgent ====
	modifyError := a.ModifyServiceRelationAgentTime(serviceRelationAgent, newRelationTime, stub)
	if modifyError != nil {
		serviceRelationAgentInvokeCallLog.Info(modifyError.Error())
		return shim.Error(modifyError.Error())
	}

	// ==== ServiceRelationAgent Time modified. Set Event ====

	eventPayload:="Modified Service RelationAgent: " + serviceRelationAgent.ServiceId + " with agent: " + serviceRelationAgent.AgentId + "from old time value: " + serviceRelationAgent.Time + "to new time value: " + newRelationTime
	payloadAsBytes := []byte(eventPayload)
	eventError := stub.SetEvent("ServiceRelationAgentTimeModifiedEvent",payloadAsBytes)
	if eventError != nil {
		serviceRelationAgentInvokeCallLog.Info("Error in event Creation: " + eventError.Error())
	}else {
		serviceRelationAgentInvokeCallLog.Info("Event Modifiy ServiceRelationAgent Time OK")
	}

	serviceRelationAgentInvokeCallLog.Info("Modify Service RelationAgent Time OK")
	return shim.Success(nil)
}

// ============================================================================================================================
// Query ServiceRelationAgent - wrapper of GetServiceNotFoundError called from the chaincode invoke
// ============================================================================================================================
func QueryServiceRelationAgent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "relationId"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	relationId := args[0]

	// ==== get the serviceRelationAgent ====
	serviceRelationAgent, err := a.GetServiceRelationAgent(stub, relationId)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to find serviceRelationAgent by id " + relationId)
		return shim.Error(err.Error())
	} else {
		serviceRelationAgentInvokeCallLog.Info("Service ID: " + serviceRelationAgent.ServiceId + ", Agent: " + serviceRelationAgent.AgentId + ", with Cost: " + serviceRelationAgent.Cost + ", with Time: " + serviceRelationAgent.Time)
		// ==== Marshal the byService query result ====
		serviceAsJSON, err := json.Marshal(serviceRelationAgent)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(serviceAsJSON)
	}
}

// ========================================================================================================================
// Query by Service Agent Relation - wrapper of GetByService called from chiancode's Invoke
// ========================================================================================================================
func QueryByServiceAgentRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ServiceId"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	serviceRelationAgentInvokeCallLog.Info("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== Check if already existing service ====
	service, err := a.GetServiceNotFoundError(stub, serviceId)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to find service  by id " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Run the byService query ====
	byServiceQuery, err := a.GetByService(serviceId, stub)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to get service relation " + serviceId)
		return shim.Error(err.Error())
	}

	fmt.Printf("Agents that expose the service: %s, with Description: %s\n", service.Name, service.Description)

	// ==== Print the byService query result ====
	err = a.PrintByServiceResultsIterator(byServiceQuery, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success with payload====
	return shim.Success(nil)
}

// ========================================================================================================================
// GetAgentsByService - wrapper of GetByService called from chiancode's Invoke, for looking for agents that provide certain service
// ========================================================================================================================
func GetAgentsByService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ServiceId"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	serviceRelationAgentInvokeCallLog.Info("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== Check if already existing service ====
	service, err := a.GetServiceNotFoundError(stub, serviceId)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("The service doesn't exist " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Run the byService query ====
	byServiceQuery, err := a.GetByService(serviceId, stub)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("The service " + service.Name + " is not mapped with any agent " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Get the Agents for the byService query result ====
	agentSlice, err := a.GetAgentSliceFromByServiceQuery(byServiceQuery, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Marshal the byService query result ====
	agentsAsJSON, err := json.Marshal(agentSlice)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success with payload====
	return shim.Success(agentsAsJSON)
}

// ========================================================================================================================
// GetServiceRelationAgentByServiceWithCostAndTime - wrapper of GetByService called from chiancode's Invoke, for looking for agents that provide certain service
// ========================================================================================================================
func GetServiceRelationAgentByServiceWithCostAndTime(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ServiceId"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	serviceRelationAgentInvokeCallLog.Info("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== Check if already existing service ====
	service, err := a.GetServiceNotFoundError(stub, serviceId)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("The service doesn't exist " + serviceId)
		return shim.Error("The service doesn't exist: " + err.Error())
	}

	// ==== Run the byService query ====
	byServiceQueryIterator, err := a.GetByService(serviceId, stub)
	// byServiceQueryIterator, err := stub.GetStateByPartialCompositeKey("service~agent~relation", []string{serviceId})
	defer byServiceQueryIterator.Close()

	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("The service " + service.Name + " is not mapped with any agent " + serviceId)
		return shim.Error(err.Error())
	}
	if byServiceQueryIterator != nil {
		serviceRelationAgentInvokeCallLog.Info(&byServiceQueryIterator)
	}

	// ==== Get the Agents for the byService query result ====
	serviceRelationSlice, err := a.GetServiceRelationSliceFromRangeQuery(byServiceQueryIterator, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	byServiceQueryIterator.Close()
	// ==== Marshal the byService query result ====
	fmt.Print(serviceRelationSlice)
	agentsByServiceAsBytes, err := json.Marshal(serviceRelationSlice)
	if err != nil {
		return shim.Error(err.Error())
	}
	serviceRelationAgentInvokeCallLog.Info(agentsByServiceAsBytes)

	stringOut := string(agentsByServiceAsBytes)
	serviceRelationAgentInvokeCallLog.Info(stringOut)
	if stringOut == "null" {
		serviceRelationAgentInvokeCallLog.Info("Service exists but has no existing relationships with agents")
		return shim.Error("Service exists but has no existing relationships with agents")
	}

	// ==== Return success with agentsByServiceSliceAsBytes as payload ====
	return shim.Success(agentsByServiceAsBytes)
}

// ========================================================================================================================
// GetServiceRelationAgentByServiceWithCostAndTimeNotFoundError - wrapper of GetByService called from chiancode's Invoke, for looking for agents that provide certain service, return Error if not found
// ========================================================================================================================
func GetServiceRelationAgentByAgentWithCostAndTimeNotFoundError(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "AgentId"

	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	serviceRelationAgentInvokeCallLog.Info("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	agentId := args[0]

	// ==== Check if already existing agent ====
	agent, err := a.GetAgentNotFoundError(stub, agentId)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("The agent doesn't exist " + agentId)
		return shim.Error("The agent doesn't exist: " + err.Error())
	}

	// ==== Run the byAgent query ====
	byAgentQueryIterator, err := a.GetByAgent(agentId, stub)
	defer byAgentQueryIterator.Close()

	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("The agent " + agent.Name + " is not mapped with any service " + agentId)
		return shim.Error(err.Error())
	}
	if byAgentQueryIterator != nil {
		serviceRelationAgentInvokeCallLog.Info(&byAgentQueryIterator)
	}

	// ==== Get the Agents for the byService query result ====
	serviceRelationSlice, err := a.GetServiceRelationSliceFromRangeQuery(byAgentQueryIterator, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	byAgentQueryIterator.Close()
	// ==== Marshal the byService query result ====
	fmt.Print(serviceRelationSlice)
	servicesByAgentAsBytes, err := json.Marshal(serviceRelationSlice)
	if err != nil {
		return shim.Error(err.Error())
	}
	serviceRelationAgentInvokeCallLog.Info(servicesByAgentAsBytes)

	stringOut := string(servicesByAgentAsBytes)
	serviceRelationAgentInvokeCallLog.Info(stringOut)
	if stringOut == "null" {
		serviceRelationAgentInvokeCallLog.Info("Agent exists but has no existing relationships with services")
		return shim.Error("Service exists but has no existing relationships with agents")
	}

	// ==== AgentServiceRelation saved & indexed. Return success with payload ====
	return shim.Success(servicesByAgentAsBytes)
}

// ========================================================================================================================
// GetServiceRelationAgentByServiceWithCostAndTime - wrapper of GetByAgent called from chiancode's Invoke, for looking for services provided by the agent, return null if not found
// ========================================================================================================================
func GetServiceRelationAgentByAgentWithCostAndTime(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "AgentId"

	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	serviceRelationAgentInvokeCallLog.Info("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	agentId := args[0]

	// ==== Check if already existing agent ====
	agent, err := a.GetAgentNotFoundError(stub, agentId)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("The agent doesn't exist " + agentId)
		return shim.Error("The agent doesn't exist: " + err.Error())
	}

	// ==== Run the byAgent query ====
	byAgentQueryIterator, err := a.GetByAgent(agentId, stub)
	defer byAgentQueryIterator.Close()

	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("The agent " + agent.Name + " is not mapped with any service " + agentId)
		return shim.Error(err.Error())
	}
	if byAgentQueryIterator != nil {
		serviceRelationAgentInvokeCallLog.Info(&byAgentQueryIterator)
	}

	// ==== Get the Agents for the byService query result ====
	serviceRelationSlice, err := a.GetServiceRelationSliceFromRangeQuery(byAgentQueryIterator, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	byAgentQueryIterator.Close()
	// ==== Marshal the byService query result ====
	fmt.Print(serviceRelationSlice)
	servicesByAgentAsBytes, err := json.Marshal(serviceRelationSlice)
	if err != nil {
		return shim.Error(err.Error())
	}
	serviceRelationAgentInvokeCallLog.Info(servicesByAgentAsBytes)

	stringOut := string(servicesByAgentAsBytes)
	serviceRelationAgentInvokeCallLog.Info("stringOut Value: " + stringOut)

	// ==== AgentServiceRelation saved & indexed. Return success with payload ====
	return shim.Success(servicesByAgentAsBytes)
}

// ========================================================================================================================
// Query by Agent Service Relation - wrapper of GetByAgent called from chiancode's Invoke
// ========================================================================================================================
func QueryByAgentServiceRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "AgentId"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	agentId := args[0]

	// ==== Check if already existing agent ====
	agent, err := a.GetAgentNotFoundError(stub, agentId)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to find agent  by id " + agentId)
		return shim.Error(err.Error())
	}

	// ==== Run the byAgent query ====
	byAgentQuery, err := a.GetByAgent(agentId, stub)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to get agent relation " + agentId)
		return shim.Error(err.Error())
	}

	fmt.Printf("The agent %s expose the services:\n", agent.Name)

	// ==== Print the byService query result ====
	printError := a.PrintByAgentResultsIterator(byAgentQuery, stub)
	if printError != nil {
		return shim.Error(printError.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success ====
	return shim.Success(nil)
}

// ============================================================================================================================
// Remove Service Agent Relation - wrapper of DeleteServiceRelationAgent a marble from state and from marble index Shows Off DelState() - "removing"" a key/value from the ledger
// UNSAFE function
// ============================================================================================================================
func DeleteServiceAgentRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	serviceRelationAgentInvokeCallLog.Info("starting delete serviceRelationAgent agent relation")

	//   0
	// "RelationId"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// input sanitation
	err := arglib.SanitizeArguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	relationId := args[0]

	// get the serviceRelationAgent
	serviceRelationAgent, err := a.GetServiceRelationAgent(stub, relationId)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to find serviceRelationAgent by relationId " + relationId)
		return shim.Error(err.Error())
	}

	// remove the serviceRelationAgent
	err = a.DeleteServiceRelationAgent(stub, relationId) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Printf("Deleted serviceRelationAgent that maps ServiceId: %s, with AgentId: %s of Cost: %s, Time: %s\n", serviceRelationAgent.ServiceId, serviceRelationAgent.AgentId, serviceRelationAgent.Cost, serviceRelationAgent.Time)
	return shim.Success(nil)
}

// =====================================================================================================================
// DeleteServiceRelationAgentApi() - remove a service from state and from service index
//
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs:
//      0
//     RelationId
// =====================================================================================================================
func DeleteServiceRelationAgentAndIndexes(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "RelationId"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		serviceRelationAgentInvokeCallLog.Error(argumentSizeError)
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}
	// input sanitation
	err := arglib.SanitizeArguments(args)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Error(err.Error())
		return shim.Error(err.Error())
	}

	// get args into variables
	relationId := args[0]

	// get the serviceRelationAgent
	serviceRelationAgent, err := a.GetServiceRelationAgentNotFoundError(stub, relationId)
	if err != nil {
		serviceRelationAgentInvokeCallLog.Info("Failed to find serviceRelationAgent by relationId " + relationId)
		return shim.Error(err.Error())
	}

	// ==== remove the serviceRelationAgent ====
	err = stub.DelState(relationId) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete serviceRelationAgent: " + err.Error())
	}

	// ==== remove the indexes ====
	indexNameService := "service~agent~relation"
	err = a.DeleteServiceIndex(stub, indexNameService,serviceRelationAgent.ServiceId,serviceRelationAgent.AgentId,serviceRelationAgent.RelationId)
	if err != nil {
		return shim.Error("Failed to delete serviceRelationAgent Agent Index: " + err.Error())
	}

	indexNameAgent := "agent~service~relation"
	err = a.DeleteAgentIndex(stub, indexNameAgent,serviceRelationAgent.AgentId,serviceRelationAgent.ServiceId,serviceRelationAgent.RelationId)
	if err != nil {
		return shim.Error("Failed to delete serviceRelationAgent Agent Index: " + err.Error())
	}

	// ==== ServiceRelationAgent and indexed deleted. Set Event ====
	eventPayload:="Deleted Service RelationAgent: " + serviceRelationAgent.RelationId + ", of service: " + serviceRelationAgent.ServiceId + ", with agent: " + serviceRelationAgent.AgentId
	payloadAsBytes := []byte(eventPayload)
	eventError := stub.SetEvent("ServiceRelationAgentDeletedEvent",payloadAsBytes)
	if eventError != nil {
		serviceRelationAgentInvokeCallLog.Error("Error in event Creation: " + eventError.Error())
	}else {
		serviceRelationAgentInvokeCallLog.Info("Event Delete ServiceRelationAgent OK")
	}

	// ==== ServiceRelationAgent saved & indexed. Return success ====
	serviceRelationAgentInvokeCallLog.Info("Deleted serviceRelationAgent: " + serviceRelationAgent.RelationId)
	return shim.Success(nil)
}


package invokeapi

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pavva91/arglib"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
	// a "github.com/pavva91/trustreputationledger/assets"
	a "github.com/pavva91/assets"

)

// =====================================================================================================================
// Init Service And Service Agent Relation - Same as InitServiceAgentRelation, but if the service doesn't exist
// it will create the service (and relative indexes) first
// =====================================================================================================================
func CreateServiceAndServiceAgentRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0            1             2                     3         4       5         6
	// "ServiceId", "ServiceName", "ServiceDescription", "AgentId", "Cost", "Time","initReputationValue"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 7)
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
	serviceName := args[1]
	serviceDescription := args[2]
	agentId := args[3]
	cost := args[4]
	time := args[5]
	initReputationValue := args[6]

	// ==== Check if already existing agent ====
	agent, errA := a.GetAgentNotFoundError(stub, agentId)
	if errA != nil {
		fmt.Println("Failed to find agent by id " + agentId)
		return shim.Error("Failed to find agent by id: " + errA.Error())
	}

	// ==== Check if already existing service ====
	service, errS := a.GetServiceNotFoundError(stub, serviceId)
	if errS != nil {
		// se il servizio non esiste lo creo
		fmt.Println("Failed to find service by id " + serviceId)
		errorCreateAndIndex := a.CreateAndIndexService(serviceId, serviceName, serviceDescription, stub)
		if errorCreateAndIndex != nil {
			return shim.Error("Error in creating and indexing service: " + errorCreateAndIndex.Error())
		}
	}

	// ==== Check, Create, Indexing ServiceRelationAgent ====
	serviceRelationAgent, serviceRelationError := a.CheckingCreatingIndexingServiceRelationAgent(serviceId, agentId, cost, time, stub)
	if serviceRelationError != nil {
		return shim.Error("Error saving ServiceRelationAgent: " + serviceRelationError.Error())

	}

	// ==== Check, Create, Indexing Reputation ====
	reputation,reputationError := a.CheckingCreatingIndexingReputation(serviceId,agentId,"EXECUTER",initReputationValue,stub)
	if reputationError != nil {
		return shim.Error("Error saving Agent reputation: " + reputationError.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success ====
	fmt.Println("Servizio: " + service.Name + " mappato con l'agente: " + agent.Name + " al costo: " + serviceRelationAgent.Cost + " e tempo: " + serviceRelationAgent.Time + " nella relazione con reputazione iniziale: "+ reputation.Value)
	return shim.Success(nil)
}

// ========================================================================================================================
// Init Service And Service Agent Relation With the Standard Value of Reputation (initReputationValue := 6)- Same as InitServiceAgentRelation, but if the service doesn't exist
// it will create the service (and relative indexes) first
// ========================================================================================================================
func CreateServiceAndServiceAgentRelationWithStandardValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0            1             2                     3         4       5
	// "ServiceId", "ServiceName", "ServiceDescription", "AgentId", "Cost", "Time"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 6)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	initReputationValue := "6.0"

	serviceId := args[0]
	serviceName := args[1]
	serviceDescription := args[2]
	agentId := args[3]
	cost := args[4]
	time := args[5]

	// ==== Check if already existing agent ====
	agent, errA := a.GetAgentNotFoundError(stub, agentId)
	if errA != nil {
		fmt.Println("Failed to find agent by id " + agentId)
		return shim.Error("Failed to find agent by id: " + errA.Error())
	}

	// ==== Check if already existing service ====
	service, errS := a.GetServiceNotFoundError(stub, serviceId)
	if errS != nil {
		// se il servizio non esiste lo creo
		fmt.Println("Failed to find service by id " + serviceId)
		errorCreateAndIndex := a.CreateAndIndexService(serviceId, serviceName, serviceDescription, stub)
		if errorCreateAndIndex != nil {
			return shim.Error("Error in creating and indexing service: " + errorCreateAndIndex.Error())
		}
	}

	// ==== Check, Create, Indexing ServiceRelationAgent ====

	serviceRelationAgent, serviceRelationError := a.CheckingCreatingIndexingServiceRelationAgent(serviceId, agentId, cost, time, stub)
	if serviceRelationError != nil {
		return shim.Error("Error saving ServiceRelationAgent: " + serviceRelationError.Error())

	}

	// ==== Check, Create, Indexing Reputation ====

	reputation,reputationError := a.CheckingCreatingIndexingReputation(agentId,serviceId,"EXECUTER",initReputationValue,stub)
	if reputationError != nil {
		return shim.Error("Error saving Agent reputation: " + reputationError.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success ====
	fmt.Println("Servizio: " + service.Name + " mappato con l'agente: " + agent.Name + " al costo: " + serviceRelationAgent.Cost + " e tempo: " + serviceRelationAgent.Time + " nella relazione con reputazione iniziale: "+ reputation.Value)
	return shim.Success(nil)
}

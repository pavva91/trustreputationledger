package invokeapi

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"encoding/json"
	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/pavva91/arglib"
	"github.com/pavva91/servicemarbles/model"
)

// ========================================================================================================================
// Init Service Agent Relation - wrapper of CreateServiceAgentRelation called from chiancode's Invoke
// ========================================================================================================================
func InitServiceAgentRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1       2       3         4
	// "ServiceId", "AgentId", "Cost", "Time", "AgentReputation"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 5)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	fmt.Println("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]
	agentId := args[1]
	cost := args[2]
	time := args[3]
	agentReputation, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return shim.Error("Wrong emptyValue inserted in AgentReputation, need float64: " + err.Error())
	}

	fmt.Println(serviceId)
	fmt.Println(agentId)
	fmt.Println(cost)
	fmt.Println(time)

	// ==== Check if already existing service ====
	service, errS := model.GetService(stub, serviceId)
	if errS != nil {
		fmt.Println("Failed to find service by id " + serviceId)
		return shim.Error("Failed to find service by id " + errS.Error())
	}
	fmt.Println("Service ok")

	// ==== Check if already existing agent ====
	agent, errA := model.GetAgent(stub, agentId)
	if errA != nil {
		fmt.Println("Failed to find agent by id " + agentId)
		return shim.Error("Failed to find agent by id: " + errA.Error())
	}

	// ==== Check if serviceRelationAgent already exists ====
	// TODO: Definire come creare relationId, per ora è composto dai due ID (serviceId + agentId)
	serviceRelationAgentId := serviceId + agentId
	agent2AsBytes, err := stub.GetState(serviceRelationAgentId)
	if err != nil {
		return shim.Error("Failed to get service agent relation: " + err.Error())
	} else if agent2AsBytes != nil {
		fmt.Println("This service agent relation already exists with relationId: " + serviceRelationAgentId)
		return shim.Error("This service agent relation already exists with relationId: " + serviceRelationAgentId)
	}

	// ==== Actual creation of serviceRelationAgent  ====
	relationId := serviceId + agentId
	serviceRelationAgent, err := model.CreateServiceAgentRelation(relationId, serviceId, agentId, cost, time, agentReputation, stub)
	if err != nil {
		return shim.Error("Failed to create service agent relation of service " + service.Name + " with agent " + agent.Name)
	}

	// ==== Indexing of serviceRelationAgent by Service ====

	// index create
	serviceAgentIndexKey, serviceIndexError := model.CreateServiceIndex(serviceRelationAgent, stub)
	if serviceIndexError != nil {
		return shim.Error(serviceIndexError.Error())
	}
	//  Note - passing a 'nil' emptyValue will effectively delete the key from state, therefore we pass null character as emptyValue
	//  Save index entry to state. Only the key Name is needed, no need to store a duplicate copy of the ServiceAgentRelation.
	emptyValue := []byte{0x00}
	// index save
	putStateError := stub.PutState(serviceAgentIndexKey, emptyValue)
	if putStateError != nil {
		return shim.Error("Error  saving Service index: " + putStateError.Error())
	}

	// ==== Indexing of serviceRelationAgent by Agent ====

	// index create
	agentServiceIndexKey, agentIndexError := model.CreateAgentIndex(serviceRelationAgent, stub)
	if agentIndexError != nil {
		return shim.Error(agentIndexError.Error())
	}
	// index save
	putStateAgentIndexError := stub.PutState(agentServiceIndexKey, emptyValue)
	if putStateAgentIndexError != nil {
		return shim.Error("Error  saving Agent index: " + putStateAgentIndexError.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success ====
	fmt.Println("Servizio: " + service.Name + " mappato con l'agente: " + agent.Name + " nella relazione con reputazione: " + strconv.FormatFloat(serviceRelationAgent.AgentReputation, 'f', 6, 64) + " - end init serviceRelationAgent")
	return shim.Success(nil)
}

// ========================================================================================================================
// Init Service And Service Agent Relation - Same as InitServiceAgentRelation, but if the service doesn't exist
// it will create the service (and relative indexes) first
// ========================================================================================================================
func InitServiceAndServiceAgentRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0            1             2                     3         4       5       6
	// "ServiceId", "ServiceName", "ServiceDescription", "AgentId", "Cost", "Time", "AgentReputation"
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

	serviceId := args[0]
	serviceName := args[1]
	serviceDescription := args[2]
	agentId := args[3]
	cost := args[4]
	time := args[5]
	agentReputation, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return shim.Error("Wrong emptyValue inserted in AgentReputation, need float64: " + err.Error())
	}

	fmt.Println(serviceId)
	fmt.Println(serviceName)
	fmt.Println(serviceDescription)
	fmt.Println(agentId)
	fmt.Println(cost)
	fmt.Println(time)

	// ==== Check if already existing agent ====
	agent, errA := model.GetAgent(stub, agentId)
	if errA != nil {
		fmt.Println("Failed to find agent by id " + agentId)
		return shim.Error("Failed to find agent by id: " + errA.Error())
	}

	// ==== Check if already existing service ====
	service, errS := model.GetService(stub, serviceId)
	if errS != nil {
		// se il servizio non esiste lo creo
		fmt.Println("Failed to find service by id " + serviceId)
		errorCreateAndIndex := model.CreateAndIndexService(serviceId, serviceName, serviceDescription, stub)
		if errorCreateAndIndex != nil {
			return shim.Error("Error in creating and indexing service: " + errorCreateAndIndex.Error())
		}
	}

	// ==== Check if serviceRelationAgent already exists ====
	// TODO: Definire come creare relationId, per ora è composto dai due ID (serviceId + agentId)
	serviceRelationAgentId := serviceId + agentId
	agent2AsBytes, err := stub.GetState(serviceRelationAgentId)
	if err != nil {
		return shim.Error("Failed to get service agent relation: " + err.Error())
	} else if agent2AsBytes != nil {
		fmt.Println("This service agent relation already exists with relationId: " + serviceRelationAgentId)
		return shim.Error("This service agent relation already exists with relationId: " + serviceRelationAgentId)
	}

	// ==== Actual creation of serviceRelationAgent  ====
	relationId := serviceId + agentId
	serviceRelationAgent, err := model.CreateServiceAgentRelation(relationId, serviceId, agentId, cost, time, agentReputation, stub)
	if err != nil {
		return shim.Error("Failed to create service agent relation of service " + service.Name + " with agent " + agent.Name)
	}

	// ==== Indexing of serviceRelationAgent by Service ====

	// index create
	serviceAgentIndexKey, serviceIndexError := model.CreateServiceIndex(serviceRelationAgent, stub)
	if serviceIndexError != nil {
		return shim.Error(serviceIndexError.Error())
	}
	//  Note - passing a 'nil' emptyValue will effectively delete the key from state, therefore we pass null character as emptyValue
	//  Save index entry to state. Only the key Name is needed, no need to store a duplicate copy of the ServiceAgentRelation.
	emptyValue := []byte{0x00}
	// index save
	putStateError := stub.PutState(serviceAgentIndexKey, emptyValue)
	if putStateError != nil {
		return shim.Error("Error  saving Service index: " + putStateError.Error())
	}

	// ==== Indexing of serviceRelationAgent by Agent ====

	// index create
	agentServiceIndexKey, agentIndexError := model.CreateAgentIndex(serviceRelationAgent, stub)
	if agentIndexError != nil {
		return shim.Error(agentIndexError.Error())
	}
	// index save
	putStateAgentIndexError := stub.PutState(agentServiceIndexKey, emptyValue)
	if putStateAgentIndexError != nil {
		return shim.Error("Error  saving Agent index: " + putStateAgentIndexError.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success ====
	fmt.Println("Servizio: " + service.Name + " mappato con l'agente: " + agent.Name + " nella relazione con reputazione: " + strconv.FormatFloat(serviceRelationAgent.AgentReputation, 'f', 6, 64) + " - end init serviceRelationAgent")
	return shim.Success(nil)
}

// ============================================================================================================================
// Query ServiceRelationAgent - wrapper of GetService called from the chaincode invoke
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
	serviceRelationAgent, err := model.GetServiceRelationAgent(stub, relationId)
	if err != nil {
		fmt.Println("Failed to find serviceRelationAgent by id " + relationId)
		return shim.Error(err.Error())
	} else {
		fmt.Println("Service ID: " + serviceRelationAgent.ServiceId + ", Agent: " + serviceRelationAgent.AgentId + ", with Cost: " + serviceRelationAgent.Cost + ", with Time: " + serviceRelationAgent.Time + ", with Reputation: " + strconv.FormatFloat(serviceRelationAgent.AgentReputation, 'f', 6, 64))
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

	fmt.Println("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== Check if already existing service ====
	service, err := model.GetService(stub, serviceId)
	if err != nil {
		fmt.Println("Failed to find service  by id " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Run the byService query ====
	byServiceQuery, err := model.GetByService(serviceId, stub)
	if err != nil {
		fmt.Println("Failed to get service relation " + serviceId)
		return shim.Error(err.Error())
	}

	fmt.Printf("Agents that expose the service: %s, with Description: %s\n", service.Name, service.Description)

	// ==== Print the byService query result ====
	err = model.PrintByServiceResultsIterator(byServiceQuery, stub)
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

	fmt.Println("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== Check if already existing service ====
	service, err := model.GetService(stub, serviceId)
	if err != nil {
		fmt.Println("The service doesn't exist " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Run the byService query ====
	byServiceQuery, err := model.GetByService(serviceId, stub)
	if err != nil {
		fmt.Println("The service " + service.Name + " is not mapped with any agent " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Get the Agents for the byService query result ====
	agentSlice, err := model.GetAgentSliceFromByServiceQuery(byServiceQuery, stub)
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

	fmt.Println("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== Check if already existing service ====
	service, err := model.GetService(stub, serviceId)
	if err != nil {
		fmt.Println("The service doesn't exist " + serviceId)
		return shim.Error("The service doesn't exist: " + err.Error())
	}

	// ==== Run the byService query ====
	byServiceQueryIterator, err := model.GetByService(serviceId, stub)
	// byServiceQueryIterator, err := stub.GetStateByPartialCompositeKey("service~agent~relation", []string{serviceId})
	defer byServiceQueryIterator.Close()

	if err != nil {
		fmt.Println("The service " + service.Name + " is not mapped with any agent " + serviceId)
		return shim.Error(err.Error())
	}
	if byServiceQueryIterator != nil {
		fmt.Println(&byServiceQueryIterator)
	}

	// ==== Get the Agents for the byService query result ====
	serviceRelationSlice, err := model.GetServiceRelationSliceFromRangeQuery(byServiceQueryIterator, stub)
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
	fmt.Println(agentsByServiceAsBytes)

	stringOut := string(agentsByServiceAsBytes)
	fmt.Println(stringOut)
	if stringOut == "null" {
		fmt.Println("Service exists but has no existing relationships with agents")
		return shim.Error("Service exists but has no existing relationships with agents")
	}

	// ==== AgentServiceRelation saved & indexed. Return success with payload ====
	return shim.Success(agentsByServiceAsBytes)
}

// ========================================================================================================================
// GetServiceRelationAgentByServiceWithCostAndTime - wrapper of GetByService called from chiancode's Invoke, for looking for agents that provide certain service
// ========================================================================================================================
func GetServiceRelationAgentByAgentWithCostAndTime(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "AgentId"

	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	fmt.Println("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	agentId := args[0]

	// ==== Check if already existing agent ====
	agent, err := model.GetAgent(stub, agentId)
	if err != nil {
		fmt.Println("The agent doesn't exist " + agentId)
		return shim.Error("The agent doesn't exist: " + err.Error())
	}

	// ==== Run the byService query ====
	byAgentQueryIterator, err := model.GetByAgent(agentId, stub)
	defer byAgentQueryIterator.Close()

	if err != nil {
		fmt.Println("The agent " + agent.Name + " is not mapped with any service " + agentId)
		return shim.Error(err.Error())
	}
	if byAgentQueryIterator != nil {
		fmt.Println(&byAgentQueryIterator)
	}

	// ==== Get the Agents for the byService query result ====
	serviceRelationSlice, err := model.GetServiceRelationSliceFromRangeQuery(byAgentQueryIterator, stub)
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
	fmt.Println(servicesByAgentAsBytes)

	stringOut := string(servicesByAgentAsBytes)
	fmt.Println(stringOut)
	if stringOut == "null" {
		fmt.Println("Service exists but has no existing relationships with agents")
		return shim.Error("Service exists but has no existing relationships with agents")
	}

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
	agent, err := model.GetAgent(stub, agentId)
	if err != nil {
		fmt.Println("Failed to find agent  by id " + agentId)
		return shim.Error(err.Error())
	}

	// ==== Run the byAgent query ====
	byAgentQuery, err := model.GetByAgent(agentId, stub)
	if err != nil {
		fmt.Println("Failed to get agent relation " + agentId)
		return shim.Error(err.Error())
	}

	fmt.Printf("The agent %s expose the services:\n", agent.Name)

	// ==== Print the byService query result ====
	printError := model.PrintByAgentResultsIterator(byAgentQuery, stub)
	if printError != nil {
		return shim.Error(printError.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success ====
	return shim.Success(nil)
}

// ============================================================================================================================
// Remove Service Agent Relation - wrapper of DeleteServiceAgentRelation a marble from state and from marble index Shows Off DelState() - "removing"" a key/value from the ledger
// UNSAFE function, TODO: you have to remove also the indexes
// ============================================================================================================================
func RemoveServiceAgentRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting delete serviceRelationAgent agent relation")

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
	serviceRelationAgent, err := model.GetServiceRelationAgent(stub, relationId)
	if err != nil {
		fmt.Println("Failed to find serviceRelationAgent by relationId " + relationId)
		return shim.Error(err.Error())
	}

	// remove the serviceRelationAgent
	err = model.DeleteServiceAgentRelation(stub, relationId) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Printf("Deleted serviceRelationAgent that maps ServiceId: %s, with AgentId: %s of Cost: %s, Time: %s, Agent reputation: %s\n", serviceRelationAgent.ServiceId, serviceRelationAgent.AgentId, serviceRelationAgent.Cost, serviceRelationAgent.Time,
		strconv.FormatFloat(serviceRelationAgent.AgentReputation, 'f', 6, 64))
	return shim.Success(nil)
}


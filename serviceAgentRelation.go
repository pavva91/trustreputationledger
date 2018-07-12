package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"strconv"
	"errors"
)

type ServiceRelationAgent struct {
	ServiceId       string  `json:"ServiceId"`
	AgentId         string  `json:"AgentId"`
	Cost            string  `json:"Cost"`            //TODO: Usare float64
	Time            string  `json:"Time"`            //TODO: Usare float64
	AgentReputation float64 `json:"AgentReputation"` // Reputation as agent in Executor Role
}

type OutputStructure struct {
	agentList []Agent
	serviceRelationAgent []ServiceRelationAgent
}

// ========================================================================================================================
// Init Service Agent Relation - wrapper of createServiceAgentRelation called from chiancode's Invoke
// ========================================================================================================================
func initServiceAgentRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1       2       3         4
	// "ServiceId", "AgentId", "Cost", "Time", "AgentReputation"
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	fmt.Println("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := sanitize_arguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]
	agentId := args[1] // ID INCREMENTALE DEVE ESSERE PASSATO DA JAVA APPLICATION (PER ORA UGUALE AL NOME)
	cost := args[2]
	time := args[3]
	agentReputation, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return shim.Error("Wrong value inserted in AgentReputation, need float64: " + err.Error())
	}

	// ==== Check if already existing agent ====
	agent, err := getAgent(stub,agentId)
	if err != nil{
		fmt.Println("Failed to find agent by id " + agentId)
		return shim.Error(err.Error())
	}

	// ==== Check if already existing service ====
	service, err := getService(stub,serviceId)
	if err != nil{
		fmt.Println("Failed to find service by id " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Check if serviceRelationAgent already exists ====
	serviceRelationAgentId := serviceId + agentId
	agent2AsBytes, err := stub.GetState(serviceRelationAgentId)
	if err != nil {
		return shim.Error("Failed to get service agent relation: " + err.Error())
	} else if agent2AsBytes != nil {
		fmt.Println("This service agent relation already exists: " + cost)
		return shim.Error("This service agent relation already exists: " + cost)
	}

	// ==== Actual creation of serviceRelationAgent  ====
	serviceRelationAgent,err := createServiceAgentRelation(serviceId, agentId, cost, time, agentReputation, stub)
	if err != nil {
		return shim.Error("Failed to create service agent relation of service " + service.Name + " with agent " + agent.Name)
	}

	relationId := serviceId+agentId

	// ==== Indexing of serviceRelationAgent by Service ====

	// index create
	serviceAgentIndexKey, serviceIndexError := createServiceIndex(relationId,serviceRelationAgent,stub)
	if serviceIndexError != nil {
		return shim.Error(serviceIndexError.Error())
	}

	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	//  Save index entry to state. Only the key Name is needed, no need to store a duplicate copy of the ServiceAgentRelation.
	value := []byte{0x00}

	// index save
	putStateError:=stub.PutState(serviceAgentIndexKey, value)
	if putStateError != nil {
		return shim.Error(putStateError.Error())
	}

	// ==== Indexing of serviceRelationAgent by Agent ====

	// index create
	agentServiceIndexKey, agentIndexError := createAgentIndex(relationId,serviceRelationAgent,stub)
	if agentIndexError != nil {
		return shim.Error(agentIndexError.Error())
	}

	// index save
	putStateAgentIndexError:=stub.PutState(agentServiceIndexKey, value)
	if putStateAgentIndexError != nil {
		return shim.Error(putStateAgentIndexError.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success ====
	fmt.Println("Servizio: " + service.Name + " mappato con l'agente: " + agent.Name + " nella relazione con reputazione: " + strconv.FormatFloat(serviceRelationAgent.AgentReputation,'f',6,64) + " - end init serviceRelationAgent")
	return shim.Success(nil)
}

// ============================================================================================================================
// Query ServiceRelationAgent - wrapper of getService called from the chaincode invoke
// ============================================================================================================================
func queryServiceRelationAgent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "relationId"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// ==== Input sanitation ====
	sanitizeError := sanitize_arguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	relationId := args[0]

	// ==== get the serviceRelationAgent ====
	serviceRelationAgent, err := getServiceRelationAgent(stub, relationId)
	if err != nil{
		fmt.Println("Failed to find serviceRelationAgent by id " + relationId)
		return shim.Error(err.Error())
	}else {
		fmt.Println("Service ID: " + serviceRelationAgent.ServiceId +", Agent: " + serviceRelationAgent.AgentId + ", with Cost: " + serviceRelationAgent.Cost + ", with Time: " + serviceRelationAgent.Time + ", with Reputation: " + strconv.FormatFloat(serviceRelationAgent.AgentReputation,'f',6,64))
		// ==== Marshal the byService query result ====
		serviceAsJSON, err := json.Marshal(serviceRelationAgent)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(serviceAsJSON)
	}
}

// ========================================================================================================================
// Query by Service Agent Relation - wrapper of getByService called from chiancode's Invoke
// ========================================================================================================================
func queryByServiceAgentRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ServiceId"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fmt.Println("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := sanitize_arguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== Check if already existing service ====
	service, err := getService(stub,serviceId)
	if err != nil{
		fmt.Println("Failed to find service  by id " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Run the byService query ====
	byServiceQuery, err := getByService(serviceId,stub)
	if err != nil{
		fmt.Println("Failed to get service relation " + serviceId)
		return shim.Error(err.Error())
	}

	fmt.Printf("Agents that expose the service: %s, with Description: %s\n",service.Name,service.Description)

	// ==== Print the byService query result ====
	err = printByServiceResultsIterator(byServiceQuery,stub)
	if err != nil{
		return shim.Error(err.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success with payload====
	return shim.Success(nil)
}

// ========================================================================================================================
// getAgentsByService - wrapper of getByService called from chiancode's Invoke, for looking for agents that provide certain service
// ========================================================================================================================
func getAgentsByService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ServiceId"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fmt.Println("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := sanitize_arguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== Check if already existing service ====
	service, err := getService(stub,serviceId)
	if err != nil{
		fmt.Println("The service doesn't exist " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Run the byService query ====
	byServiceQuery, err := getByService(serviceId,stub)
	if err != nil{
		fmt.Println("The service " + service.Name + " is not mapped with any agent " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Get the Agents for the byService query result ====
	agentSlice, err := getAgentSliceFromByServiceQuery(byServiceQuery,stub)
	if err != nil{
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
// getServiceRelationAgentByServiceWithCostAndTime - wrapper of getByService called from chiancode's Invoke, for looking for agents that provide certain service
// ========================================================================================================================
func getServiceRelationAgentByServiceWithCostAndTime(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ServiceId"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fmt.Println("- start init serviceRelationAgent")

	// ==== Input sanitation ====
	sanitizeError := sanitize_arguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== Check if already existing service ====
	service, err := getService(stub,serviceId)
	if err != nil{
		fmt.Println("The service doesn't exist " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Run the byService query ====
	byServiceQuery, err := getByService(serviceId,stub)
	if err != nil{
		fmt.Println("The service " + service.Name + " is not mapped with any agent " + serviceId)
		return shim.Error(err.Error())
	}

	// ==== Get the Agents for the byService query result ====
	serviceRelationSlice, err := getServiceRelationSliceFromByServiceQuery(byServiceQuery,stub)
	if err != nil{
		return shim.Error(err.Error())
	}

	// ==== Marshal the byService query result ====
	agentsByServiceAsJSON, err := json.Marshal(serviceRelationSlice)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success with payload ====
	return shim.Success(agentsByServiceAsJSON)
}






// ========================================================================================================================
// Query by Agent Service Relation - wrapper of getByAgent called from chiancode's Invoke
// ========================================================================================================================
func queryByAgentServiceRelation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "AgentId"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// ==== Input sanitation ====
	sanitizeError := sanitize_arguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	agentId := args[0]

	// ==== Check if already existing agent ====
	agent, err := getAgent(stub, agentId)
	if err != nil{
		fmt.Println("Failed to find agent  by id " + agentId)
		return shim.Error(err.Error())
	}

	// ==== Run the byAgent query ====
	byAgentQuery, err := getByAgent(agentId,stub)
	if err != nil{
		fmt.Println("Failed to get agent relation " + agentId)
		return shim.Error(err.Error())
	}

	fmt.Printf("The agent %s expose the services:\n", agent.Name)

	// ==== Print the byService query result ====
	printError := printByAgentResultsIterator(byAgentQuery,stub)
	if printError != nil{
		return shim.Error(printError.Error())
	}


	// ==== AgentServiceRelation saved & indexed. Return success ====
	return shim.Success(nil)
}

// ============================================================================================================================
// Remove Service Agent Relation - wrapper of deleteServiceAgentRelation a marble from state and from marble index Shows Off DelState() - "removing"" a key/value from the ledger
// ============================================================================================================================
func removeServiceAgentRelation(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting delete serviceRelationAgent agent relation")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	relationId := args[0]

	// get the serviceRelationAgent
	serviceRelationAgent, err := getServiceRelationAgent(stub, relationId)
	if err != nil{
		fmt.Println("Failed to find serviceRelationAgent by relationId " + relationId)
		return shim.Error(err.Error())
	}

	// remove the serviceRelationAgent
	err = deleteServiceAgentRelation(stub,relationId) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Printf("Deleted serviceRelationAgent that maps ServiceId: %s, with AgentId: %s of Cost: %s, Time: %s, Agent reputation: %s\n",serviceRelationAgent.ServiceId,serviceRelationAgent.AgentId,serviceRelationAgent.Cost,serviceRelationAgent.Time,
		strconv.FormatFloat(serviceRelationAgent.AgentReputation,'f',6,64) )
	return shim.Success(nil)
}

// ============================================================
// createServiceAgentMapping - create a new mapping service agent
// ============================================================
func  createServiceAgentRelation(serviceId string, agentId string, cost string, time string, agentReputation float64, stub shim.ChaincodeStubInterface) (*ServiceRelationAgent,error) {
	// ==== Create marble object and marshal to JSON ====
	serviceRelationAgent := &ServiceRelationAgent{serviceId, agentId, cost, time, agentReputation}
	serviceRelationAgentJSONAsBytes, _ := json.Marshal(serviceRelationAgent)

	// === Save marble to state ===
	relationId := serviceId+agentId
	stub.PutState(relationId, serviceRelationAgentJSONAsBytes)

	return serviceRelationAgent,nil
}

// ============================================================================================================================
// Create Service Based Index - to do query based on Service
// ============================================================================================================================
func createServiceIndex(relationId string, serviceRelationAgent *ServiceRelationAgent, stub shim.ChaincodeStubInterface)  (serviceAgentIndexKey string, err error){
	//  ==== Index the serviceAgentRelation to enable service-based range queries, e.g. return all x services ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on service~agent~relation.
	//  This will enable very efficient state range queries based on composite keys matching service~agent~relation
	indexName := "service~agent~relation"
	serviceAgentIndexKey, err = stub.CreateCompositeKey(indexName, []string{serviceRelationAgent.ServiceId, serviceRelationAgent.AgentId, relationId})
	if err != nil {
		return serviceAgentIndexKey, err
	}
	return serviceAgentIndexKey,nil
}

// ============================================================================================================================
// Create Agent Based Index - to do query based on Agent
// ============================================================================================================================
func createAgentIndex(relationId string, serviceRelationAgent *ServiceRelationAgent, stub shim.ChaincodeStubInterface)  (agentServiceIndex string, err error){
	//  ==== Index the serviceAgentRelation to enable service-based range queries, e.g. return all x agents ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on agent~service~relation.
	//  This will enable very efficient state range queries based on composite keys matching agent~service~relation
	indexName := "agent~service~relation"
	agentServiceIndex, err = stub.CreateCompositeKey(indexName, []string{serviceRelationAgent.AgentId, serviceRelationAgent.ServiceId, relationId})
	if err != nil {
		return agentServiceIndex, err
	}
	return agentServiceIndex,nil
}

// ============================================================================================================================
// Get Service Agent Relation - get the service agent relation asset from ledger
// ============================================================================================================================
func getServiceRelationAgent(stub shim.ChaincodeStubInterface, relationId string) (ServiceRelationAgent, error) {
	var serviceRelationAgent ServiceRelationAgent
	serviceRelationAgentAsBytes, err := stub.GetState(relationId) //getState retreives a key/value from the ledger
	if err != nil {                                            //this seems to always succeed, even if key didn't exist
		return serviceRelationAgent, errors.New("Failed to get serviceRelationAgent - " + relationId)
	}
	json.Unmarshal(serviceRelationAgentAsBytes, &serviceRelationAgent) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo ServiceRelationAgent)

	return serviceRelationAgent, nil
}

// ============================================================================================================================
// Get the service query on ServiceRelationAgent - Execute the query based on service composite index
// ============================================================================================================================
func getByService(serviceId string, stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, error){
	// Query the service~agent~relation index by service
	// This will execute a key range query on all keys starting with 'service'
	serviceAgentResultsIterator, err := stub.GetStateByPartialCompositeKey("service~agent~relation", []string{serviceId})
	if err != nil {
		return serviceAgentResultsIterator, err
	}
	// defer serviceAgentResultsIterator.Close()
	return serviceAgentResultsIterator, nil
}

// ============================================================================================================================
// Get the agent query on ServiceRelationAgent - Execute the query based on agent composite index
// ============================================================================================================================
func getByAgent(serviceId string, stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, error){
	// Query the service~agent~relation index by service
	// This will execute a key range query on all keys starting with 'service'
	agentServiceResultsIterator, err := stub.GetStateByPartialCompositeKey("agent~service~relation", []string{serviceId})
	if err != nil {
		return agentServiceResultsIterator, err
	}
	// defer agentServiceResultsIterator.Close()
	return agentServiceResultsIterator, nil
}

// ============================================================================================================================
// Delete Service Agent Relation - delete from state and from marble index Shows Off DelState() - "removing"" a key/value from the ledger
// ============================================================================================================================
func deleteServiceAgentRelation(stub shim.ChaincodeStubInterface, relationId string) error {
	// remove the serviceRelationAgent
	err := stub.DelState(relationId) //remove the key from chaincode state
	if err != nil {
		return err
	}
	return nil
}

// ============================================================================================================================
// Delete Service Agent Relation - delete from state and from marble index Shows Off DelState() - "removing"" a key/value from the ledger
// ============================================================================================================================
func deleteServiceIndex(stub shim.ChaincodeStubInterface, indexName string, serviceId string, agentId string, relationId string) error {
	// remove the serviceRelationAgent
	// TODO: Capire come funziona, perché prima crea la composite key?
	agentServiceIndex, err := stub.CreateCompositeKey(indexName, []string{serviceId, agentId, relationId})
	if err != nil {
		return err
	}
	err = stub.DelState(agentServiceIndex) //remove the key from chaincode state
	if err != nil {
		return err
	}
	return nil
}

// ============================================================================================================================
// Delete Agent Service Relation - delete from state and from marble index Shows Off DelState() - "removing"" a key/value from the ledger
// ============================================================================================================================
func deleteAgentIndex(stub shim.ChaincodeStubInterface, indexName string, agentId string, serviceId string, relationId string) error {
	// remove the serviceRelationAgent
	agentServiceIndex, err := stub.CreateCompositeKey(indexName, []string{agentId, serviceId, relationId})
	if err != nil {
		return err
	}
	err = stub.DelState(agentServiceIndex) //remove the key from chaincode state
	if err != nil {
		return err
	}
	return nil
}

// ============================================================================================================================
// getAgentSliceFromByServiceQuery - Get the Agent and ServiceRelationAgent Slices from the result of query "byService"
// ============================================================================================================================
func getServiceRelationSliceFromByServiceQuery(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) ([]ServiceRelationAgent, error){
	var serviceRelationAgentSlice []ServiceRelationAgent

	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return nil,err
		}
		// get the service agent relation from service~agent~relation composite key
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		relationId:=compositeKeyParts[2]

		iserviceRelationAgent, err := getServiceRelationAgent(stub,relationId)
		serviceRelationAgentSlice = append(serviceRelationAgentSlice, iserviceRelationAgent)
		if err != nil {
			return nil,err
		}
		fmt.Printf("- found a relation RELATION ID: %s \n",  relationId)
	}
	fmt.Println(serviceRelationAgentSlice[0].AgentId)
	return serviceRelationAgentSlice,nil
}

// ============================================================================================================================
// getAgentSliceFromByServiceQuery - Get the Agent Slice from the result of query "byService"
// ============================================================================================================================
func getAgentSliceFromByServiceQuery(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) ([]Agent, error){
	var agentSlice []Agent
	for i := 0; queryIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble Name from the composite key
		responseRange, err := queryIterator.Next()
		if err != nil {
			return nil,err
		}
		// get the service agent relation from service~agent~relation composite key
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		agentId:=compositeKeyParts[1]

		iAgent,err := getAgent(stub,agentId)
		agentSlice =append(agentSlice, iAgent)

		if err != nil {
			return nil,err
		}
	}
	return agentSlice,nil
}

// ============================================================================================================================
// Print Results Iterator - Print on screen the general iterator of the composite index query result
// ============================================================================================================================
func printByServiceResultsIterator(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface)  error{
	for i := 0; queryIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble Name from the composite key
		responseRange, err := queryIterator.Next()
		if err != nil {
			return err
		}
		// get the service agent relation from service~agent~relation composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		serviceId:=compositeKeyParts[0]
		agentId:=compositeKeyParts[1]
		relationId:=compositeKeyParts[2]

		if err != nil {
			return err
		}
		fmt.Printf("- found a relation from OBJECT_TYPE:%s SERVICE ID:%s AGENT ID:%s RELATION ID: %s\n", objectType, serviceId, agentId, relationId)
	}
	return nil
}
// ============================================================================================================================
// Print Results Iterator - Print on screen the general iterator of the composite index query result
// ============================================================================================================================
func printByAgentResultsIterator(iteratorInterface shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) error{
	for i := 0; iteratorInterface.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble Name from the composite key
		responseRange, err := iteratorInterface.Next()
		if err != nil {
			return err
		}
		// get the service agent relation from service~agent~relation composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		agentId:=compositeKeyParts[0]
		serviceId:=compositeKeyParts[1]
		relationId:=compositeKeyParts[2]

		if err != nil {
			return err
		}
		fmt.Printf("- found a relation from OBJECT_TYPE:%s AGENT ID:%s SERVICE ID:%s  RELATION ID: %s\n", objectType, agentId, serviceId, relationId)
	}
	return nil
}

// ============================================================================================================================
// Print Results Iterator - Print on screen the general iterator of the composite index query result
// ============================================================================================================================
func printResultsIterator(iteratorInterface shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) error{
	for i := 0; iteratorInterface.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble Name from the composite key
		responseRange, err := iteratorInterface.Next()
		if err != nil {
			return err
		}
		// get the service agent relation from service~agent~relation composite key
		// get the agent service relation from agent~service~relation composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return err
		}
		fmt.Printf("- found a relation from OBJECT_TYPE:%s SERVICE ID:%s AGENT ID:%s RELATION ID: %s\n", objectType, compositeKeyParts[0], compositeKeyParts[1],compositeKeyParts[2])
	}
	return nil
}
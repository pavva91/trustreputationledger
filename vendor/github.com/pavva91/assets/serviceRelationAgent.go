/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
*/
package assets

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var serviceRelationAgentLog = shim.NewLogger("serviceRelationAgent")

type ServiceRelationAgent struct {
	RelationId      string  `json:"RelationId"`// relationId := serviceId + agentId
	ServiceId       string  `json:"ServiceId"`
	AgentId         string  `json:"AgentId"`
	Cost            string  `json:"Cost"`            //TODO: Usare float64
	Time            string  `json:"Time"`            //TODO: Usare float64
	// AgentReputation float64 `json:"AgentReputation"` //TODO: Se uso Reputation lo devo levare
}

// =====================================================================================================================
// createServiceAgentMapping - create a new mapping service agent
// =====================================================================================================================
func CreateServiceAgentRelation(relationId string, serviceId string, agentId string, cost string, time string,  stub shim.ChaincodeStubInterface) (*ServiceRelationAgent, error) {
	// ==== Create marble object and marshal to JSON ====
	serviceRelationAgent := &ServiceRelationAgent{relationId, serviceId, agentId, cost, time}
	serviceRelationAgentJSONAsBytes, _ := json.Marshal(serviceRelationAgent)

	// === Save marble to state ===
	stub.PutState(relationId, serviceRelationAgentJSONAsBytes)

	return serviceRelationAgent, nil
}

// =====================================================================================================================
// Create Service Based Index - to do query based on Service
// =====================================================================================================================
func CreateServiceIndex(serviceRelationAgent *ServiceRelationAgent, stub shim.ChaincodeStubInterface) (serviceAgentIndexKey string, err error) {
	//  ==== Index the serviceAgentRelation to enable service-based range queries, e.g. return all x services ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on service~agent~relation.
	//  This will enable very efficient state range queries based on composite keys matching service~agent~relation
	indexName := "service~agent~relation"
	serviceAgentIndexKey, err = stub.CreateCompositeKey(indexName, []string{serviceRelationAgent.ServiceId, serviceRelationAgent.AgentId, serviceRelationAgent.RelationId})
	if err != nil {
		return serviceAgentIndexKey, err
	}
	return serviceAgentIndexKey, nil
}

// =====================================================================================================================
// Create Agent Based Index - to do query based on Agent
// =====================================================================================================================
func CreateAgentIndex(serviceRelationAgent *ServiceRelationAgent, stub shim.ChaincodeStubInterface) (agentServiceIndex string, err error) {
	//  ==== Index the serviceAgentRelation to enable service-based range queries, e.g. return all x agents ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on agent~service~relation.
	//  This will enable very efficient state range queries based on composite keys matching agent~service~relation
	indexName := "agent~service~relation"
	agentServiceIndex, err = stub.CreateCompositeKey(indexName, []string{serviceRelationAgent.AgentId, serviceRelationAgent.ServiceId, serviceRelationAgent.RelationId})
	if err != nil {
		return agentServiceIndex, err
	}
	return agentServiceIndex, nil
}

// =====================================================================================================================
// CheckingCreatingIndexingServiceRelationAgent - Incapsulate the three tasks:
// 1. CHECKING
// 2. CREATING
// 3. INDEXING
// =====================================================================================================================
func CheckingCreatingIndexingServiceRelationAgent(serviceId string, agentId string, cost string, time string, stub shim.ChaincodeStubInterface) (*ServiceRelationAgent, error){

	// ==== Check if serviceRelationAgent already exists ====
	// TODO: Definire come creare relationId, per ora è composto dai due ID (serviceId + agentId)
	relationId := serviceId + agentId
	agent2AsBytes, err := stub.GetState(relationId)
	if err != nil {
		return nil,errors.New("Failed to get service agent relation: " + err.Error())
	} else if agent2AsBytes != nil {
		serviceLog.Info("This service agent relation already exists with relationId: " + relationId)
		return nil,errors.New("This service agent relation already exists with relationId: " + relationId)
	}

	// ==== Actual creation of serviceRelationAgent  ====
	serviceRelationAgent, err := CreateServiceAgentRelation(relationId, serviceId, agentId, cost, time, stub)
	if err != nil {
		return nil,errors.New("Failed to create service agent relation of service " + serviceId + " with agent " + agentId)
	}

	// ==== Indexing of serviceRelationAgent by Service ====

	// index create
	serviceAgentIndexKey, serviceIndexError := CreateServiceIndex(serviceRelationAgent, stub)
	if serviceIndexError != nil {
		return nil,errors.New(serviceIndexError.Error())
	}
	//  Note - passing a 'nil' emptyValue will effectively delete the key from state, therefore we pass null character as emptyValue
	//  Save index entry to state. Only the key Name is needed, no need to store a duplicate copy of the ServiceAgentRelation.
	emptyValue := []byte{0x00}
	// index save
	putStateError := stub.PutState(serviceAgentIndexKey, emptyValue)
	if putStateError != nil {
		return nil,errors.New("Error  saving Service index: " + putStateError.Error())
	}

	// ==== Indexing of serviceRelationAgent by Agent ====

	// index create
	agentServiceIndexKey, agentIndexError := CreateAgentIndex(serviceRelationAgent, stub)
	if agentIndexError != nil {
		return nil,errors.New(agentIndexError.Error())
	}
	// index save
	putStateAgentIndexError := stub.PutState(agentServiceIndexKey, emptyValue)
	if putStateAgentIndexError != nil {
		return nil,errors.New("Error  saving Agent index: " + putStateAgentIndexError.Error())
	}

	return serviceRelationAgent,nil
}

// =====================================================================================================================
// ModifyServiceRelationAgentCost - Modify the serviceRelationAgent cost of the asset passed as parameter
// =====================================================================================================================
func ModifyServiceRelationAgentCost(serviceRelationAgent ServiceRelationAgent, newRelationCost string, stub shim.ChaincodeStubInterface) (error) {

	serviceRelationAgent.Cost = newRelationCost

	serviceRelationAgentAsBytes, _ := json.Marshal(serviceRelationAgent)
	putStateError := stub.PutState(serviceRelationAgent.RelationId, serviceRelationAgentAsBytes)
	if putStateError != nil {
		serviceRelationAgentLog.Error(putStateError)
		return errors.New(putStateError.Error())
	}
	return nil
}

// =====================================================================================================================
// ModifyServiceRelationAgentTime - Modify the serviceRelationAgent time of the asset passed as parameter
// =====================================================================================================================
func ModifyServiceRelationAgentTime(serviceRelationAgent ServiceRelationAgent, newRelationTime string, stub shim.ChaincodeStubInterface) (error) {

	serviceRelationAgent.Time = newRelationTime

	serviceRelationAgentAsBytes, _ := json.Marshal(serviceRelationAgent)
	putStateError := stub.PutState(serviceRelationAgent.RelationId, serviceRelationAgentAsBytes)
	if putStateError != nil {
		serviceRelationAgentLog.Error(putStateError)
		return errors.New(putStateError.Error())
	}
	return nil
}

// ============================================================================================================================
// Get Service Agent Relation - get the service agent relation asset from ledger - return (nil,nil) if not found
// ============================================================================================================================
func GetServiceRelationAgent(stub shim.ChaincodeStubInterface, relationId string) (ServiceRelationAgent, error) {
	var serviceRelationAgent ServiceRelationAgent
	serviceRelationAgentAsBytes, err := stub.GetState(relationId) //getState retreives a key/value from the ledger
	if err != nil {                                               //this seems to always succeed, even if key didn't exist
		return serviceRelationAgent, errors.New("Error in finding service relation with agent: " + err.Error())
	}

	json.Unmarshal(serviceRelationAgentAsBytes, &serviceRelationAgent) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo ServiceRelationAgent)

	return serviceRelationAgent, nil
}

// =====================================================================================================================
// Get Service Agent Relation Not Found Error - get the service agent relation asset from ledger - throws error if not found (error!=nil ---> key not found)
// =====================================================================================================================
func GetServiceRelationAgentNotFoundError(stub shim.ChaincodeStubInterface, relationId string) (ServiceRelationAgent, error) {
	var serviceRelationAgent ServiceRelationAgent
	serviceRelationAgentAsBytes, err := stub.GetState(relationId) //getState retreives a key/value from the ledger
	if err != nil {                                               //this seems to always succeed, even if key didn't exist
		return serviceRelationAgent, errors.New("Error in finding service relation with agent: " + err.Error())
	}

	if serviceRelationAgentAsBytes == nil {
		return ServiceRelationAgent{}, errors.New("ServiceRelationAgent non found, RelationId: " + relationId)
	}
	json.Unmarshal(serviceRelationAgentAsBytes, &serviceRelationAgent) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo ServiceRelationAgent)

	return serviceRelationAgent, nil
}

// =====================================================================================================================
// Get the service query on ServiceRelationAgent - Execute the query based on service composite index
// =====================================================================================================================
func GetByService(serviceId string, stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, error) {
	// Query the service~agent~relation index by service
	// This will execute a key range query on all keys starting with 'service'
	serviceAgentResultsIterator, err := stub.GetStateByPartialCompositeKey("service~agent~relation", []string{serviceId})
	if err != nil {
		return serviceAgentResultsIterator, err
	}
	defer serviceAgentResultsIterator.Close()
	return serviceAgentResultsIterator, nil
}

// =====================================================================================================================
// Get the agent query on ServiceRelationAgent - Execute the query based on agent composite index
// =====================================================================================================================
// TODO: BUG, ritorna anche il record dell'agente
func GetByAgent(agentId string, stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, error) {
	// Query the service~agent~relation index by service
	// This will execute a key range query on all keys starting with 'agent'
	agentServiceResultsIterator, err := stub.GetStateByPartialCompositeKey("agent~service~relation", []string{agentId})
	if err != nil {
		return agentServiceResultsIterator, err
	}
	defer agentServiceResultsIterator.Close()
	return agentServiceResultsIterator, nil
}

// =====================================================================================================================
// Delete Service Relation Agent - delete from state DelState() - "removing"" a key/value from the ledger
// =====================================================================================================================
func DeleteServiceRelationAgent(stub shim.ChaincodeStubInterface, relationId string) error {
	// remove the serviceRelationAgent
	err := stub.DelState(relationId) //remove the key from chaincode state
	if err != nil {
		return err
	}
	return nil
}

// =====================================================================================================================
// Delete Service Index - Delete the index
// =====================================================================================================================
func DeleteServiceIndex(stub shim.ChaincodeStubInterface, indexName string, serviceId string, agentId string, relationId string) error {

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

// =====================================================================================================================
// Delete Agent Service Relation - delete from state and from marble index Shows Off DelState() - "removing"" a key/value from the ledger
// =====================================================================================================================
func DeleteAgentIndex(stub shim.ChaincodeStubInterface, indexName string, agentId string, serviceId string, relationId string) error {
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

// =====================================================================================================================
// GetAgentSliceFromByServiceQuery - Get the Agent and ServiceRelationAgent Slices from the result of query "byService"
// =====================================================================================================================
func GetServiceRelationSliceFromRangeQuery(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) ([]ServiceRelationAgent, error) {
	var serviceRelationAgentSlice []ServiceRelationAgent
	// get the service agent relation from service~agent~relation composite key
	// defer queryIterator.Close()

	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return nil, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		relationId := compositeKeyParts[2]

		iserviceRelationAgent, err := GetServiceRelationAgentNotFoundError(stub, relationId)
		serviceRelationAgentSlice = append(serviceRelationAgentSlice, iserviceRelationAgent)
		if err != nil {
			return nil, err
		}
		fmt.Printf("- found a relation RELATION ID: %s \n", relationId)
	}
	queryIterator.Close()
	return serviceRelationAgentSlice, nil
}

// =====================================================================================================================
// GetAgentSliceFromByServiceQuery - Get the Agent Slice from the result of query "byService"
// =====================================================================================================================
func GetAgentSliceFromByServiceQuery(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) ([]Agent, error) {
	var agentSlice []Agent
	for i := 0; queryIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble Name from the composite key
		responseRange, err := queryIterator.Next()
		if err != nil {
			return nil, err
		}
		// get the service agent relation from service~agent~relation composite key
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		agentId := compositeKeyParts[1]

		iAgent, err := GetAgentNotFoundError(stub, agentId)
		agentSlice = append(agentSlice, iAgent)

		if err != nil {
			return nil, err
		}
	}
	queryIterator.Close()
	return agentSlice, nil
}

// =====================================================================================================================
// Print Results Iterator - Print on screen the general iterator of the composite index query result
// =====================================================================================================================
func PrintByServiceResultsIterator(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) error {
	for i := 0; queryIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble Name from the composite key
		responseRange, err := queryIterator.Next()
		if err != nil {
			return err
		}
		// get the service agent relation from service~agent~relation composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		serviceId := compositeKeyParts[0]
		agentId := compositeKeyParts[1]
		relationId := compositeKeyParts[2]

		if err != nil {
			return err
		}
		fmt.Printf("- found a relation from OBJECT_TYPE:%s SERVICE ID:%s AGENT ID:%s RELATION ID: %s\n", objectType, serviceId, agentId, relationId)
	}
	return nil
}

// =====================================================================================================================
// Print Results Iterator - Print on screen the general iterator of the composite index query result
// =====================================================================================================================
func PrintByAgentResultsIterator(iteratorInterface shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) error {
	for i := 0; iteratorInterface.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble Name from the composite key
		responseRange, err := iteratorInterface.Next()
		if err != nil {
			return err
		}
		// get the service agent relation from service~agent~relation composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		agentId := compositeKeyParts[0]
		serviceId := compositeKeyParts[1]
		relationId := compositeKeyParts[2]

		if err != nil {
			return err
		}
		fmt.Printf("- found a relation from OBJECT_TYPE:%s AGENT ID:%s SERVICE ID:%s  RELATION ID: %s\n", objectType, agentId, serviceId, relationId)
	}
	return nil
}

// =====================================================================================================================
// Print Results Iterator - Print on screen the general iterator of the composite index query result
// =====================================================================================================================
func PrintResultsIterator(iteratorInterface shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) error {
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
		fmt.Printf("- found a relation from OBJECT_TYPE:%s SERVICE ID:%s AGENT ID:%s RELATION ID: %s\n", objectType, compositeKeyParts[0], compositeKeyParts[1], compositeKeyParts[2])
	}
	return nil
}

/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
 */
package assets

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"errors"
	"fmt"
)

// =====================================================================================================================
// Define the Agent's Reputation structure
// =====================================================================================================================
// - ReputationId
// - AgentId
// - ServiceId
// - AgentRole
// - Value
type Reputation struct {
	ReputationId        string `json:"ReputationId"`
	AgentId             string `json:"AgentId"`
	ServiceId           string `json:"ServiceId"`
	AgentRole           string `json:"AgentRole"` // TODO:Available roles: Executer, Demander
	Value               float64 `json:"Value"`  // Value of Reputation of the agent
}


//TODO: Don't delete reputation of a deleted agent

// =====================================================================================================================
// createReputation - create a new reputation identified as: service-agent-agentrole (Demander || Executer)
// =====================================================================================================================
func createReputation(reputationId string, serviceId string, agentId string, agentRole string, value float64, stub shim.ChaincodeStubInterface) (*Reputation, error) {
	// agentRoleNow := "Demander"
	// ==== Create marble object and marshal to JSON ====
	reputation := &Reputation{reputationId, agentId, serviceId,  agentRole, value}
	ReputationJSONAsBytes, _ := json.Marshal(reputation)

	// === Save marble to state ===
	stub.PutState(reputationId, ReputationJSONAsBytes)

	return reputation, nil
}

// =====================================================================================================================
// Create Agent Based Index - to do query based on Agent, Service and AgentRole
// =====================================================================================================================
func CreateAgentServiceRoleIndex(reputation *Reputation, stub shim.ChaincodeStubInterface) (agentServiceRoleIndex string, err error) {
	//  ==== Index the serviceAgentRelation to enable service-based range queries, e.g. return all x agents ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on agent~service~relation.
	//  This will enable very efficient state range queries based on composite keys matching agent~service~relation
	indexName := "agent~service~agentRole~reputation"
	agentServiceRoleIndex, err = stub.CreateCompositeKey(indexName, []string{reputation.AgentId, reputation.ServiceId, reputation.AgentRole, reputation.ReputationId})
	if err != nil {
		return agentServiceRoleIndex, err
	}
	return agentServiceRoleIndex, nil
}

// =====================================================================================================================
// Get Reputation - get the reputation asset from ledger - return (nil,nil) if not found
// =====================================================================================================================
func GetReputation(stub shim.ChaincodeStubInterface, reputationId string) (Reputation, error) {
	var serviceRelationAgent Reputation
	serviceRelationAgentAsBytes, err := stub.GetState(reputationId) //getState retreives a key/value from the ledger
	if err != nil {                                               //this seems to always succeed, even if key didn't exist
		return serviceRelationAgent, errors.New("Error in finding the reputation of the agent: " + error.Error(err))
	}

	// TODO: Levare trigger error ma gestire il payload null
	if serviceRelationAgentAsBytes == nil {
		return Reputation{}, errors.New("Reputation not found, ReputationId: " + reputationId)
	}
	json.Unmarshal(serviceRelationAgentAsBytes, &serviceRelationAgent) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo ServiceRelationAgent)

	return serviceRelationAgent, nil
}

// =====================================================================================================================
// Get Reputation Not Found Error - get the reputation asset from ledger - throws error if not found (error!=nil ---> key not found)
// =====================================================================================================================
func GetReputationNotFoundError(stub shim.ChaincodeStubInterface, reputationId string) (Reputation, error) {
	var serviceRelationAgent Reputation
	serviceRelationAgentAsBytes, err := stub.GetState(reputationId) //getState retreives a key/value from the ledger
	if err != nil {                                               //this seems to always succeed, even if key didn't exist
		return serviceRelationAgent, errors.New("Error in finding service relation with agent: " + error.Error(err))
	}

	// TODO: Levare trigger error ma gestire il payload null
	if serviceRelationAgentAsBytes == nil {
		return Reputation{}, errors.New("Service non found, ServiceId: " + reputationId)
	}
	json.Unmarshal(serviceRelationAgentAsBytes, &serviceRelationAgent) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo ServiceRelationAgent)

	return serviceRelationAgent, nil
}

// =====================================================================================================================
// Get the service query on ServiceRelationAgent - Execute the query based on service composite index
// =====================================================================================================================
func GetByAgentServiceRole(agentId string, serviceId string, agentRole string, stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, error) {
	// Query the service~agent~relation index by service
	// This will execute a key range query on all keys starting with 'service'
	serviceAgentResultsIterator, err := stub.GetStateByPartialCompositeKey("agent~service~agentRole~reputation", []string{agentId,serviceId,agentRole})
	if err != nil {
		return serviceAgentResultsIterator, err
	}
	defer serviceAgentResultsIterator.Close()
	return serviceAgentResultsIterator, nil
}

// =====================================================================================================================
// Delete Reputation - "removing"" a key/value from the ledger
// =====================================================================================================================
func DeleteReputation(stub shim.ChaincodeStubInterface, reputationId string) error {
	// remove the serviceRelationAgent
	err := stub.DelState(reputationId) //remove the key from chaincode state
	if err != nil {
		return err
	}
	return nil
}

// =====================================================================================================================
// Delete Service Agent Role Reputation - "removing"" the key/value from the ledger relative to the index
// =====================================================================================================================
func DeleteAgentServiceRoleIndex(stub shim.ChaincodeStubInterface, indexName string, agentId string, serviceId string, agentRole string, reputationId string) error {
	// remove the serviceRelationAgent
	agentServiceRoleIndex, err := stub.CreateCompositeKey(indexName, []string{agentId, serviceId, agentRole, reputationId})
	if err != nil {
		return err
	}
	err = stub.DelState(agentServiceRoleIndex) //remove the key from chaincode state
	if err != nil {
		return err
	}
	return nil
}

// =====================================================================================================================
// GetAgentSliceFromByServiceQuery - Get the Agent and ServiceRelationAgent Slices from the result of query "byService"
// =====================================================================================================================
func GetReputationSliceFromRangeQuery(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) ([]Reputation, error) {
	var serviceRelationAgentSlice []Reputation
	defer queryIterator.Close()

	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return nil, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		reputationId := compositeKeyParts[3]

		iserviceRelationAgent, err := GetReputation(stub, reputationId)
		serviceRelationAgentSlice = append(serviceRelationAgentSlice, iserviceRelationAgent)
		if err != nil {
			return nil, err
		}
		fmt.Printf("- found a reputation REPUTATION ID: %s \n", reputationId)
	}
	return serviceRelationAgentSlice, nil
}


// =====================================================================================================================
// Print Results Iterator - Print on screen the general iterator of the composite index query result
// =====================================================================================================================
func PrintByAgentServiceRoleReputationResultsIterator(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) error {
	for i := 0; queryIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble Name from the composite key
		responseRange, err := queryIterator.Next()
		if err != nil {
			return err
		}
		// get the service agent relation from service~agent~relation composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		agentId := compositeKeyParts[0]
		serviceId := compositeKeyParts[1]
		agentRole := compositeKeyParts[2]
		reputationId := compositeKeyParts[3]


		if err != nil {
			return err
		}
		fmt.Printf("- found a relation from OBJECT_TYPE:%s AGENT ID:%s SERVICE ID:%s AGENT ROLE: %s RELATION ID: %s\n", objectType, agentId, serviceId, agentRole, reputationId)
	}
	return nil
}
/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
 */
package assets

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pavva91/arglib"
)

// =====================================================================================================================
// Define the Agent structure, with 3 properties.  Structure tags are used by encoding/json library
// =====================================================================================================================
// - AgentId
// - Name
// - Address
type Agent struct {
	AgentId string `json:"AgentId"`
	Name    string `json:"Name"`
	Address string `json:"Address"`
}

// =====================================================================================================================
// CreateAgent - create a new agent and return the created agent
// =====================================================================================================================
func CreateAgent(agentId string, agentName string, agentAddress string, stub shim.ChaincodeStubInterface) *Agent {
	// ==== Create marble object and marshal to JSON ====
	agent := &Agent{AgentId: agentId, Name: agentName, Address: agentAddress}
	agentJSONAsBytes, _ := json.Marshal(agent)

	// === Save marble to state ===
	stub.PutState(agent.AgentId, agentJSONAsBytes)
	return agent
}
// =====================================================================================================================
// Get Agent Not Found Error - get an agent asset from ledger- throws error if not found (error!=nil ---> key not found)
// =====================================================================================================================
func GetAgentNotFoundError(stub shim.ChaincodeStubInterface, agentId string) (Agent, error) {
	var agent Agent
	agentAsBytes, err := stub.GetState(agentId) //getState retreives agent key/value from the ledger
	if err != nil {                             //this seems to always succeed, even if key didn't exist
		return agent, errors.New("Error in finding agent - " + error.Error(err))
	}
	fmt.Println(agentAsBytes)
	fmt.Println(agent)

	if agentAsBytes == nil {
		return agent, errors.New("Agent non found, AgentId: " + agentId)
	}

	json.Unmarshal(agentAsBytes, &agent) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo Agent)

	fmt.Println(agent)

	return agent, nil
}
// =====================================================================================================================
// Get Agent - get an agent asset from ledger - return (nil,nil) if not found
// =====================================================================================================================
func GetAgent(stub shim.ChaincodeStubInterface, agentId string) (Agent, error) {
	var agent Agent
	agentAsBytes, err := stub.GetState(agentId) //getState retreives agent key/value from the ledger
	if err != nil {                             //this seems to always succeed, even if key didn't exist
		return agent, errors.New("Error in finding agent - " + error.Error(err))
	}
	fmt.Println(agentAsBytes)
	fmt.Println(agent)


	json.Unmarshal(agentAsBytes, &agent) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo Agent)

	fmt.Println(agent)

	return agent, nil
}

func GetAllAgents(stub shim.ChaincodeStubInterface) ([]Agent, error) {
	var agents []Agent
	// ---- Get All Agents ---- //
	agentsIterator, err := stub.GetStateByRange("idagent0", "idagent99999999999999999999999999999999999")
	if err != nil {
		return nil, err
	}
	defer agentsIterator.Close()

	for agentsIterator.HasNext() {
		aKeyValue, err := agentsIterator.Next()
		if err != nil {
			return nil, err
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on agent id - ", queryKeyAsStr)
		var agent Agent
		json.Unmarshal(queryValAsBytes, &agent) //un stringify it aka JSON.parse()
		agents = append(agents, agent)
	}
	fmt.Println("agent array - ", agents)
	return agents, nil
}

// =====================================================================================================================
// DeleteAgent() - remove a agent from state and from agent index
//
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs:
//      0
//     // =====================================================================================================================
// ============================================================================================================================
func DeleteAgent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting delete_marble")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// input sanitation
	err := arglib.SanitizeArguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	agentId := args[0]

	// get the service
	service, err := GetServiceNotFoundError(stub, agentId)
	if err != nil {
		fmt.Println("Failed to find service by AgentId " + agentId)
		return shim.Error(err.Error())
	}

	// TODO: Delete anche (prima) le relazioni del servizio con gli agenti
	err = DeleteAllAgentServiceRelations(agentId, stub)
	if err != nil {
		return shim.Error("Failed to delete agent service relation: " + err.Error())
	}

	// remove the agent
	err = stub.DelState(agentId) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete agent: " + err.Error())
	}

	fmt.Println("Deleted agent: " + service.Name)
	return shim.Success(nil)
}

// =====================================================================================================================
// DeleteAllAgentServiceRelations - delete all the Agent relations with service (aka: Reference Integrity)
// =====================================================================================================================
func DeleteAllAgentServiceRelations(agentId string, stub shim.ChaincodeStubInterface) error {
	agentServiceResultsIterator, err := GetByAgent(agentId, stub)
	if err != nil {
		return err
	}
	for i := 0; agentServiceResultsIterator.HasNext(); i++ {
		responseRange, err := agentServiceResultsIterator.Next()
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

		fmt.Printf("Delete the relation: from composite key OBJECT_TYPE:%s AGENT ID:%s SERVICE ID:%s RELATION ID: %s\n", objectType, agentId, serviceId, relationId)

		// remove the serviceRelationAgent
		err = DeleteServiceAgentRelation(stub, relationId) //remove the key from chaincode state
		if err != nil {
			return err
		}

		// remove the agent index
		err = DeleteAgentIndex(stub, objectType, agentId, serviceId, relationId) //remove the key from chaincode state
		if err != nil {
			return err
		}

		// TODO: Devo rimuovere anche dall'index?

	}
	return nil
}

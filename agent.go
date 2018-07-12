package main

import ("github.com/hyperledger/fabric/core/chaincode/shim"
pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"errors"
)

// ===================================================================================
// Define the Agent structure, with 3 properties.  Structure tags are used by encoding/json library
// ===================================================================================
// - AgentId
// - Name
// - Address
type Agent struct {
	AgentId string `json:"AgentId"`
	Name    string `json:"Name"`
	Address string `json:"Address"`
}

// ============================================================================================================================
// Init Agent - wrapper of createAgent called from the chaincode invoke
// ============================================================================================================================
func  initAgent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1                 2
	// "AgentId", "agentName", "agentAddress"
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// ==== Input sanitation ====
	sanitizeError := sanitize_arguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	agentId := args[0] // ID INCREMENTALE DEVE ESSERE PASSATO DA JAVA APPLICATION (PER ORA UGUALE AL NOME)
	agentName := args[1]
	agentAddress := args[2]

	// ==== Check if Agent already exists ====
	agentAsBytes, err := stub.GetState(agentId)
	if err != nil {
		return shim.Error("Failed to get agent: " + err.Error())
	} else if agentAsBytes != nil {
		fmt.Println("This agent already exists: " + agentName)
		return shim.Error("This agent already exists: " + agentName)
	}

	agent := createAgent(agentId, agentName, agentAddress, stub)

	// indexAgent(agent, stub)
	// TODO: index agent, sarà da fare lo stesso se riesco a fare queste due tabelle?
	// ==== Service2 saved and indexed. Return success ====
	fmt.Println("Servizio: " + agent.Name + " creato - end init agent")
	return shim.Success(nil)
}

// ============================================================================================================================
// Query Agent - wrapper of getAgent called from the chaincode invoke
// ============================================================================================================================
func queryAgent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	// ==== get the agent ====
	agent, err := getAgent(stub, agentId)
	if err != nil{
		fmt.Println("Failed to find agent by id " + agentId)
		return shim.Error(err.Error())
	}else {
		fmt.Println("Agent: " + agent.Name + ", with Address: " + agent.Address + " found")
		// ==== Marshal the byService query result ====
		agentAsJSON, err := json.Marshal(agent)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(agentAsJSON)
	}
}

// ============================================================
// createAgent - create a new agent and return the created agent
// ============================================================
func createAgent(agentId string, agentName string, agentAddress string, stub shim.ChaincodeStubInterface) *Agent {
	// ==== Create marble object and marshal to JSON ====
	agent := &Agent{AgentId:agentId, Name:agentName, Address:agentAddress}
	agentJSONAsBytes, _ := json.Marshal(agent)

	// === Save marble to state ===
	stub.PutState(agent.AgentId, agentJSONAsBytes)
	return agent
}

// ============================================================================================================================
// Get Agent - get an agent asset from ledger
// ============================================================================================================================
func getAgent(stub shim.ChaincodeStubInterface, idAgent string) (Agent, error) {
	var agent Agent
	agentAsBytes, err := stub.GetState(idAgent) //getState retreives agent key/value from the ledger
	if err != nil {                                          //this seems to always succeed, even if key didn't exist
		return agent, errors.New("Failed to find agent - " + idAgent)
	}
	json.Unmarshal(agentAsBytes, &agent) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo Service)

	return agent, nil
}

func getAllAgents(stub shim.ChaincodeStubInterface) ([]Agent,error) {
	var agents []Agent
	// ---- Get All Agents ---- //
	agentsIterator, err := stub.GetStateByRange("idagent0", "idagent99999999999999999999999999999999999")
	if err != nil {
		return nil,err
	}
	defer agentsIterator.Close()

	for agentsIterator.HasNext() {
		aKeyValue, err := agentsIterator.Next()
		if err != nil {
			return nil,err
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on agent id - ", queryKeyAsStr)
		var agent Agent
		json.Unmarshal(queryValAsBytes, &agent) //un stringify it aka JSON.parse()
		agents = append(agents, agent)
	}
	fmt.Println("agent array - ", agents)
	return agents,nil
}

// ============================================================================================================================
// deleteAgent() - remove a agent from state and from agent index
//
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs:
//      0
//     ServiceId
// ============================================================================================================================
func deleteAgent(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting delete_marble")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	agentId := args[0]

	// get the service
	service, err := getService(stub, agentId)
	if err != nil{
		fmt.Println("Failed to find service by AgentId " + agentId)
		return shim.Error(err.Error())
	}

	// TODO: Delete anche (prima) le relazioni del servizio con gli agenti
	err=deleteAllAgentServiceRelations(agentId,stub)
	if err != nil {
		return shim.Error("Failed to delete agent service relation: "+ err.Error())
	}

	// remove the agent
	err = stub.DelState(agentId) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete agent: "+ err.Error())
	}

	fmt.Println("Deleted agent: " + service.Name)
	return shim.Success(nil)
}

// ============================================================
// deleteAllAgentServiceRelations - delete all the Agent relations with service (aka: Reference Integrity)
// ============================================================
func deleteAllAgentServiceRelations(agentId string, stub shim.ChaincodeStubInterface) error{
	agentServiceResultsIterator, err := getByAgent(agentId,stub)
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

		agentId:=compositeKeyParts[0]
		serviceId:=compositeKeyParts[1]
		relationId:=compositeKeyParts[2]

		if err != nil {
			return err
		}

		fmt.Printf("Delete the relation: from composite key OBJECT_TYPE:%s AGENT ID:%s SERVICE ID:%s RELATION ID: %s\n", objectType, agentId, serviceId, relationId)

		// remove the serviceRelationAgent
		err = deleteServiceAgentRelation(stub, relationId) //remove the key from chaincode state
		if err != nil {
			return err
		}

		// remove the agent index
		err = deleteAgentIndex(stub,objectType,agentId,serviceId,relationId) //remove the key from chaincode state
		if err != nil {
			return err
		}

		// TODO: Devo rimuovere anche dall'index?

	}
	return nil
}



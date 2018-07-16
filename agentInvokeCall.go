package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"encoding/json"
	pb "github.com/hyperledger/fabric/protos/peer"

)

// ============================================================================================================================
// Init Agent - wrapper of createAgent called from the chaincode invoke
// ============================================================================================================================
func initAgent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1                 2
	// "AgentId", "agentName", "agentAddress"
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// ==== Input sanitation ====
	sanitizeError := sanitizeArguments(args)
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
	sanitizeError := sanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	agentId := args[0]

	// ==== get the agent ====
	agent, err := getAgent(stub, agentId)
	if err != nil {
		fmt.Println("Failed to find agent by id " + agentId)
		return shim.Error("Failed to find agent by id: " + err.Error())
	} else {
		fmt.Println("Agent: " + agent.Name + ", with Address: " + agent.Address + " found")
		// ==== Marshal the byService query result ====
		agentAsJSON, err := json.Marshal(agent)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(agentAsJSON)
	}
}


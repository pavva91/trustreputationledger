/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
 */
package invokeapi
// WHEN CHANGE THAT NAME REFACTOR DOESN'T WORK NOW THAT IS A PACKAGE, YOU MODIFY HERE, SAVE AND FROM CLI
// DO: govendor update +vendor
// POI VAI A MODIFICARE LE VARIE CHIAMATE DEL PACKAGE INTERESSATE DAL CAMBIAMENTO

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"encoding/json"
	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/pavva91/arglib"
	// a "github.com/pavva91/trustreputationledger/assets"
	a "github.com/pavva91/assets"
)

// =====================================================================================================================
// Init Agent - wrapper of CreateAgent called from the chaincode invoke
// =====================================================================================================================
func CreateAgent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1                 2
	// "AgentId", "agentName", "agentAddress"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 3)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
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

	agent := a.CreateAgent(agentId, agentName, agentAddress, stub)

	// TODO: index agent, sarà da fare lo stesso se riesco a fare queste due tabelle?

	// ==== Agent saved and indexed. Set Event ====
	eventPayload:="Created Agent: " + agentId
	payloadAsBytes := []byte(eventPayload)
	eventError := stub.SetEvent("AgentCreatedEvent",payloadAsBytes)
	if eventError != nil {
		fmt.Println("Error in event Creation: " + eventError.Error())
	}else {
		fmt.Println("Event Create Agent OK")
	}

	// ==== Agent saved and indexed and event setted. Return success ====
	fmt.Println("Agent: " + agent.Name + " created - end init agent")
	return shim.Success(nil)
}

// ========================================================================================================================
// Modify Agent Name - wrapper of ModifyAgentName called from chiancode's Invoke
// ========================================================================================================================
func ModifyAgentName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0            1
	// "agentId", "newAgentName"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 2)
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
	newAgentName := args[1]

	// ==== get the agent ====
	agent, getError := a.GetAgentNotFoundError(stub, agentId)
	if getError != nil {
		fmt.Println("Failed to find agent by id " + agentId)
		return shim.Error(getError.Error())
	}

	// ==== modify the agent ====
	modifyError := a.ModifyAgentName(agent, newAgentName, stub)
	if modifyError != nil {
		fmt.Println("Failed to modify the agent name: " + newAgentName)
		return shim.Error(modifyError.Error())
	}

	return shim.Success(nil)
}

// ========================================================================================================================
// Modify Agent Address - wrapper of ModifyAgentAddress called from chiancode's Invoke
// ========================================================================================================================
func ModifyAgentAddress(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0            1
	// "agentId", "newAgentAddress"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 2)
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
	newAgentAddress := args[1]

	// ==== get the agent ====
	agent, getError := a.GetAgentNotFoundError(stub, agentId)
	if getError != nil {
		fmt.Println("Failed to find agent by id " + agentId)
		return shim.Error(getError.Error())
	}

	// ==== modify the agent ====
	modifyError := a.ModifyAgentAddress(agent, newAgentAddress, stub)
	if modifyError != nil {
		fmt.Println("Failed to modify the agent address: " + newAgentAddress)
		return shim.Error(modifyError.Error())
	}

	return shim.Success(nil)
}

// =====================================================================================================================
// Query Agent Not Found Error - wrapper of GetAgentNotFoundError called from the chaincode invoke
// =====================================================================================================================
func QueryAgentNotFoundError(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	// ==== get the agent ====
	agent, err := a.GetAgentNotFoundError(stub, agentId)
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
// =====================================================================================================================
// Query Agent - wrapper of GetAgent called from the chaincode invoke
// =====================================================================================================================
func QueryAgent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	// ==== get the agent ====
	agent, err := a.GetAgent(stub, agentId)
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


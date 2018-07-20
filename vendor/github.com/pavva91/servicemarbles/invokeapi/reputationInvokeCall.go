/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
*/
package invokeapi

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pavva91/arglib"
	"fmt"
	"encoding/json"
	pb "github.com/hyperledger/fabric/protos/peer"
	a "github.com/pavva91/servicemarbles/assets"


)

/*
For now we want that the Activity assets can only be added on the ledger (NO MODIFY, NO DELETE)
 */
// ========================================================================================================================
// Create Executed Service Evaluation - wrapper of CreateServiceAgentRelation called from chiancode's Invoke
// ========================================================================================================================
func CreateReputation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0         1            2             3
	// "AgentId", "ServiceId", "AgentRole", "Value"
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

	agentId := args[0]
	serviceId := args[1]
	agentRole := args[2]
	value := args[3]

	// ==== Check if already existing agent ====
	agent, errA := a.GetAgentNotFoundError(stub, agentId)
	if errA != nil {
		fmt.Println("Failed to find Agent by id " + agentId)
		return shim.Error("Failed to find Agent by id: " + errA.Error())
	}

	// ==== Check if already existing service ====
	service, errS := a.GetServiceNotFoundError(stub, serviceId)
	if errS != nil {
		fmt.Println("Failed to find service by id " + serviceId)
		return shim.Error("Failed to find service by id " + errS.Error())
	}

	// ==== Check if AgentRole == Demander || Executer ====
	if ("DEMANDER"!=agentRole && "EXECUTER"!=agentRole){
		return shim.Error("Wrong Agent Role: " + agentRole + ", use \"DEMANDER\"or \"EXECUTER\"")
	}

	// ==== Check if reputation already exists ====
	// TODO: Definire come creare reputationId, per ora è composto dai tre ID (agentId + serviceId + agentRole)
	reputationId := agentId + serviceId + agentRole
	reputationAsBytes, err := stub.GetState(reputationId)
	if err != nil {
		return shim.Error("Failed to get service agent reputation: " + err.Error())
	} else if reputationAsBytes != nil {
		fmt.Println("This service agent reputation already exists with reputationId: " + reputationId)
		return shim.Error("This service agent reputation already exists with reputationId: " + reputationId)
	}

	// ==== Actual creation of Reputation  ====
	reputation, err := a.CreateReputation(reputationId, agentId, serviceId, agentRole, value, stub)
	if err != nil {
		return shim.Error("Failed to create service agent relation of service " + service.Name + " with agent " + agent.Name)
	}

	// ==== Indexing of reputation by Service Tx Id ====

	// index create
	agentReputationIndex, serviceIndexError := a.CreateAgentServiceRoleIndex(reputation, stub)
	if serviceIndexError != nil {
		return shim.Error(serviceIndexError.Error())
	}
	//  Note - passing a 'nil' emptyValue will effectively delete the key from state, therefore we pass null character as emptyValue
	//  Save index entry to state. Only the key Name is needed, no need to store a duplicate copy of the ServiceAgentRelation.
	emptyValue := []byte{0x00}
	// index save
	putStateError := stub.PutState(agentReputationIndex, emptyValue)
	if putStateError != nil {
		return shim.Error("Error saving Agent Reputation index: " + putStateError.Error())
	}

	// ==== Reputation saved & indexed. Return success ====
	fmt.Println("ReputationId: " + reputation.ReputationId + " of agent: " + reputation.AgentId + " in role of: " + reputation.AgentRole + " relative to the service: " + reputation.ServiceId)
	return shim.Success(nil)
}

// ========================================================================================================================
// Create Executed Service Evaluation - wrapper of CreateServiceAgentRelation called from chiancode's Invoke
// ========================================================================================================================
func ModifyReputationValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0            1
	// "reputationId", "newReputationValue"
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

	reputationId := args[0]
	newReputationValue := args[1]

	// ==== get the reputation ====
	reputation, getError := a.GetReputationNotFoundError(stub, reputationId)
	if getError != nil {
		fmt.Println("Failed to find reputation by id " + reputationId)
		return shim.Error(getError.Error())
	}

	// ==== modify the reputation ====
	modifyError := a.ModifyReputationValue(reputation,newReputationValue,stub)
	if modifyError != nil {
		fmt.Println("Failed to modify the reputation value: " + newReputationValue)
		return shim.Error(modifyError.Error())
	}

	return shim.Success(nil)
}


// ============================================================================================================================
// Query Reputation - wrapper of GetReputation called from the chaincode invoke
// ============================================================================================================================
func QueryReputation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "reputationId"
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

	reputationId := args[0]

	// ==== get the reputation ====
	reputation, err := a.GetReputationNotFoundError(stub, reputationId)
	if err != nil {
		fmt.Println("Failed to find reputation by id " + reputationId)
		return shim.Error(err.Error())
	} else {
		fmt.Println("Reputation ID: " + reputation.ReputationId + ", of Agent: " + reputation.AgentId + ", Agent Role: " + reputation.AgentRole + ", of the Service: " + reputation.ServiceId)
		// ==== Marshal the Get Service Evaluation query result ====
		evaluationAsJSON, err := json.Marshal(reputation)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(evaluationAsJSON)
	}
}

// ========================================================================================================================
// QueryByAgentServiceRole - wrapper of GetByAgentServiceRole called from chiancode's Invoke
// TODO: Per come è impostato l'id ora è "inutile", però in vista di refactor ID sarà utile
// ========================================================================================================================
func QueryByAgentServiceRole(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0           1            2
	// "agentId", "serviceId", "agentRole"
	argumentSizeError := arglib.ArgumentSizeLimitVerification(args, 3)
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

	var byAgentServiceRoleQuery shim.StateQueryIteratorInterface
	var err error

	// ==== Run the byAgentServiceRole query ====
	switch len(args) {
	case 3:
		serviceId := args[1]
		agentRole := args[2]
		byAgentServiceRoleQuery, err = a.GetByAgentServiceRole(agentId, serviceId, agentRole, stub)
		if err != nil {
			fmt.Println("Failed to get reputation for this agent: " + agentId + ", in this service: " + serviceId + ", in this role: " + agentRole)
			return shim.Error(err.Error())
		}
	case 2:
		serviceId := args[1]
		byAgentServiceRoleQuery, err = a.GetByAgentService(agentId, serviceId, stub)
		if err != nil {
			fmt.Println("Failed to get reputation for this agent: " + agentId + ", in this service: " + serviceId)
			return shim.Error(err.Error())
		}
	case 1:
		byAgentServiceRoleQuery, err = a.GetByAgentOnly(agentId, stub)
		if err != nil {
			fmt.Println("Failed to get reputation for this agent: " + agentId)
			return shim.Error(err.Error())
		}
	}

	// ==== Print the byService query result ====
	err = a.PrintByAgentServiceRoleReputationResultsIterator(byAgentServiceRoleQuery, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


// =====================================================================================================================
// GetReputationsByAgentServiceRole - wrapper of GetByAgentServiceRole called from chiancode's Invoke,
// for looking for serviceEvaluations of a certain Agent-Service-AgentRole triple
// return: Reputations As JSON (IS EVERYTHING is WORKING is only ONE Reputation the result)
// =====================================================================================================================
func GetReputationsByAgentServiceRole(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0        1            2
	// "agentId", "serviceId", "agentRole"
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

	agentId := args[0]
	serviceId := args[1]
	agentRole := args[2]

	// TODO: With empty string works (FIX)
	indexName := "agent~service~agentRole~reputation"
	byAgentServiceRoleQuery, err := stub.GetStateByPartialCompositeKey(indexName, []string{"idservice12"})


	// ==== Run the byAgentServiceRole query ====
	// byAgentServiceRoleQuery, err := a.GetByAgentServiceRole(agentId, serviceId, agentRole, stub)
	if err != nil {
		fmt.Println("Failed to get reputation for this agent: " + agentId + ", in this service: " + serviceId + ", in this role: " + agentRole)
		return shim.Error(err.Error())
	}

	// ==== Get the ServiceEvaluations for the byDemanderExecuter query result ====
	reputations, err := a.GetReputationSliceFromRangeQuery(byAgentServiceRoleQuery, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Marshal the byServiceTxId query result ====
	serviceEvaluationsAsJSON, err := json.Marshal(reputations)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(serviceEvaluationsAsJSON)
}

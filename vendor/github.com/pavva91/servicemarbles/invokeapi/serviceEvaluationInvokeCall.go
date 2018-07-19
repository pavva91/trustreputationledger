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
For now we want that the ServiceEvaluation assets can only be added on the ledger (NO MODIFY, NO DELETE)
 */
// ========================================================================================================================
// Create Executed Service Evaluation - wrapper of CreateServiceAgentRelation called from chiancode's Invoke
// ========================================================================================================================
func CreateServiceEvaluation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1                   2                     3                   4                        5         6
	// "WriterAgentId", "DemanderAgentId", "ExecuterAgentId", "ExecutedServiceId", "ExecutedServiceTxId", "Timestamp", "Value"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 7)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	fmt.Println("- start init Service Evaluation")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	writerAgentId := args[0]
	demanderAgentId := args[1]
	executerAgentId := args[2]
	executedServiceId := args[3]
	executedServiceTxId := args[4]
	timestamp := args[5]
	value := args[6]

	var writerAgent a.Agent


	// TODO: Check if already existing DemanderAgent
	// ==== Check if already existing demanderAgent ====
	demanderAgent, errA := a.GetAgentNotFoundError(stub, demanderAgentId)
	if errA != nil {
		fmt.Println("Failed to find demanderAgent by id " + demanderAgentId)
		return shim.Error("Failed to find demanderAgent by id: " + errA.Error())
	}
	// TODO: Check if already existing ExecuterAgent
	// ==== Check if already existing executerAgent ====
	executerAgent, errA := a.GetAgentNotFoundError(stub, executerAgentId)
	if errA != nil {
		fmt.Println("Failed to find executerAgent by id " + executerAgentId)
		return shim.Error("Failed to find executerAgent by id: " + errA.Error())
	}
	// TODO: Check if WriterAgent == DemanderAgent || ExecuterAgent
	// ==== Check if WriterAgent == DemanderAgent || ExecuterAgent ====
	switch true {
	case demanderAgentId==writerAgentId:
		writerAgent=demanderAgent
	case executerAgentId==writerAgentId:
		writerAgent=executerAgent
	default:
		return shim.Error("Wrong Writer Agent Id: " + writerAgentId)
	}
	// TODO: Check if already existing ExecutedServiceId
	// ==== Check if already existing executedService ====
	executedService, errS := a.GetServiceNotFoundError(stub, executedServiceId)
	if errS != nil {
		fmt.Println("Failed to find executedService by id " + executedServiceId)
		return shim.Error("Failed to find executedService by id " + errS.Error())
	}

	fmt.Println("Service ok")


	// ==== Check if serviceEvaluation already exists ====
	// TODO: Definire come creare evaluationId, per ora è composto dai due ID (writerAgentId + demanderAgentId + executerAgentId + ExecutedServiceTxId)
	evaluationId := writerAgentId + demanderAgentId + executerAgentId + executedServiceTxId
	serviceEvaluationAsBytes, err := stub.GetState(evaluationId)
	if err != nil {
		return shim.Error("Failed to get executedService demanderAgent relation: " + err.Error())
	} else if serviceEvaluationAsBytes != nil {
		fmt.Println("This executedService demanderAgent relation already exists with relationId: " + evaluationId)
		return shim.Error("This executedService demanderAgent relation already exists with relationId: " + evaluationId)
	}

	// ==== Actual creation of Service Evaluation  ====
	serviceEvaluation, err := a.CreateServiceEvaluation(evaluationId, writerAgentId, demanderAgentId, executerAgentId, executedServiceId, executedServiceTxId, timestamp, value, stub)
	if err != nil {
		return shim.Error("Failed to create executedService demanderAgent relation of executedService " + executedService.Name + " with demanderAgent " + demanderAgent.Name)
	}

	// ==== Indexing of serviceEvaluation by Service Tx Id ====

	// index create
	serviceTxIndexKey, serviceIndexError := a.CreateServiceTxIndex(serviceEvaluation, stub)
	if serviceIndexError != nil {
		return shim.Error(serviceIndexError.Error())
	}
	//  Note - passing a 'nil' emptyValue will effectively delete the key from state, therefore we pass null character as emptyValue
	//  Save index entry to state. Only the key Name is needed, no need to store a duplicate copy of the ServiceAgentRelation.
	emptyValue := []byte{0x00}
	// index save
	putStateError := stub.PutState(serviceTxIndexKey, emptyValue)
	if putStateError != nil {
		return shim.Error("Error  saving Service index: " + putStateError.Error())
	}

	// ==== Indexing of serviceEvaluation by Agent ====

	// index create
	demanderExecuterIndexKey, agentIndexError := a.CreateDemanderExecuterIndex(serviceEvaluation, stub)
	if agentIndexError != nil {
		return shim.Error(agentIndexError.Error())
	}
	// index save
	putStateDemanderExecuterIndexError := stub.PutState(demanderExecuterIndexKey, emptyValue)
	if putStateDemanderExecuterIndexError != nil {
		return shim.Error("Error  saving Agent index: " + putStateDemanderExecuterIndexError.Error())
	}

	// ==== AgentServiceRelation saved & indexed. Return success ====
	fmt.Println("Servizio: " + executedService.Name + " evaluated by: " + writerAgent.Name + " relative to the transaction: " + executedServiceTxId)
	return shim.Success(nil)
}


// ============================================================================================================================
// Query ServiceRelationAgent - wrapper of GetService called from the chaincode invoke
// ============================================================================================================================
func QueryServiceEvaluation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "evaluationId"
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

	evaluationId := args[0]

	// ==== get the serviceEvaluation ====
	serviceEvaluation, err := a.GetServiceEvaluationNotFoundError(stub, evaluationId)
	if err != nil {
		fmt.Println("Failed to find serviceEvaluation by id " + evaluationId)
		return shim.Error(err.Error())
	} else {
		fmt.Println("Evaluation ID: " + serviceEvaluation.EvaluationId + ", Writer Agent: " + serviceEvaluation.WriterAgentId + ", Demander Agent: " + serviceEvaluation.DemanderAgentId + ", Executer Agent: " + serviceEvaluation.ExecuterAgentId + ", of the Service: " + serviceEvaluation.ExecutedServiceId + ", with Timestamp: " + serviceEvaluation.Timestamp + ", with Evaluation: " + serviceEvaluation.Value)
		// ==== Marshal the Get Service Evaluation query result ====
		evaluationAsJSON, err := json.Marshal(serviceEvaluation)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(evaluationAsJSON)
	}
}

// ========================================================================================================================
// Query by Executed Service Tx Id - wrapper of GetByExecutedServiceTxId called from chiancode's Invoke
// ========================================================================================================================
func QueryByExecutedServiceTx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "executedServiceTxId"
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

	executedServiceTxId := args[0]

	// ==== Run the byExecutedServiceTx query ====
	byExecutedServiceTxIdQuery, err := a.GetByExecutedServiceTx(executedServiceTxId, stub)
	if err != nil {
		fmt.Println("Failed to get service evaluation for this serviceTxId: " + executedServiceTxId)
		return shim.Error(err.Error())
	}

	// ==== Print the byService query result ====
	err = a.PrintByExecutedServiceTxIdResultsIterator(byExecutedServiceTxIdQuery, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ========================================================================================================================
// Query by Demander Executer - wrapper of GetByDemanderExecuter called from chiancode's Invoke
// ========================================================================================================================
func QueryByDemanderExecuter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0                1
	// "demanderAgentId", "executerAgentId"
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

	demanderAgentId := args[0]
	executerAgentId := args[1]


	// ==== Run the byExecutedServiceTx query ====
	byExecutedServiceTxIdQuery, err := a.GetByDemanderExecuter(demanderAgentId, executerAgentId, stub)
	if err != nil {
		fmt.Println("Failed to get service evaluation for this demander: " + demanderAgentId + " and executer: " + executerAgentId)
		return shim.Error(err.Error())
	}

	// ==== Print the byService query result ====
	err = a.PrintByDemanderExecuterResultsIterator(byExecutedServiceTxIdQuery, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// =====================================================================================================================
// GetServiceEvaluationsByExecutedServiceTxId - wrapper of GetByExecutedServiceTxId called from chiancode's Invoke,
// for looking for serviceEvaluations of a certain executedServiceTxId
// return: ServiceEvaluations As JSON
// =====================================================================================================================
func GetServiceEvaluationsByExecutedServiceTxId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ExecutedServiceTxId"
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

	executedServiceTxId := args[0]

	// ==== Run the byService query ====
	byServiceQuery, err := a.GetByExecutedServiceTx(executedServiceTxId, stub)
	if err != nil {
		fmt.Println("The service Tx Id " + executedServiceTxId + " is not mapped with any service evaluation.")
		return shim.Error(err.Error())
	}

	// ==== Get the Agents for the byServiceTxId query result ====
	serviceEvaluations, err := a.GetServiceEvaluationSliceFromServiceTxIdRangeQuery(byServiceQuery, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Marshal the byServiceTxId query result ====
	serviceEvaluationsAsJSON, err := json.Marshal(serviceEvaluations)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(serviceEvaluationsAsJSON)
}

// =====================================================================================================================
// GetServiceEvaluationsByDemanderExecuter - wrapper of GetByDemanderExecuter called from chiancode's Invoke,
// for looking for serviceEvaluations of a certain Demander-Executer couple
// return: ServiceEvaluations As JSON
// =====================================================================================================================
func GetServiceEvaluationsByDemanderExecuter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0          1
	// "Demander", "Executer"
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

	demanderAgentId := args[0]
	executerAgentId := args[1]

	// ==== Run the ByDemanderExecuter query ====
	byExecutedServiceTxIdQuery, err := a.GetByDemanderExecuter(demanderAgentId, executerAgentId, stub)
	if err != nil {
		fmt.Println("Failed to get service evaluation for this demander: " + demanderAgentId + " and executer: " + executerAgentId)
		return shim.Error(err.Error())
	}

	// ==== Get the ServiceEvaluations for the byDemanderExecuter query result ====
	serviceEvaluations, err := a.GetServiceEvaluationSliceFromDemanderExecuterRangeQuery(byExecutedServiceTxIdQuery, stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Marshal the byServiceTxId query result ====
	serviceEvaluationsAsJSON, err := json.Marshal(serviceEvaluations)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(serviceEvaluationsAsJSON)
}

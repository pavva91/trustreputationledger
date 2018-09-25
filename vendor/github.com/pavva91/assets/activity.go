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

// =====================================================================================================================
// Define the Service Evaluation structure
// =====================================================================================================================
// - ReputationId
// - AgentId
// - ServiceId
// - AgentRole
// - ExecutedServiceId
// - ExecutedServiceTxId
// - ExecutedServiceTimestamp
// - Outcome
// - Value
// - IsFinalEvaluation
// UNIVOCAL: WriterAgentId, DemanderAgentId, ExecuterAgentId, ExecutedServiceTxId
type Activity struct {
	EvaluationId             string `json:"EvaluationId"`
	WriterAgentId            string `json:"WriterAgentId"` // WriterAgentId = DemanderAgentId || ExecuterAgentId
	DemanderAgentId          string `json:"DemanderAgentId"`
	ExecuterAgentId          string `json:"ExecuterAgentId"`
	ExecutedServiceId        string `json:"ExecutedServiceId"`
	ExecutedServiceTxid      string `json:"ExecutedServiceTxid"` // Relativo all'esecuzione del servizio (TODO: a cosa serve?)
	ExecutedServiceTimestamp string `json:"ExecutedServiceTimestamp"`
	Value                    string `json:"Value"`
}

// ============================================================
// Create Service Evaluation - create a new service evaluation
// ============================================================
func CreateActivity(evaluationId string, writerAgentId string, demanderAgentId string, executerAgentId string, executedServiceId string, executedServiceTxId string, timestamp string, value string, stub shim.ChaincodeStubInterface) (*Activity, error) {
	// ==== Create marble object and marshal to JSON ====
	serviceEvaluation := &Activity{evaluationId, writerAgentId, demanderAgentId, executerAgentId, executedServiceId, executedServiceTxId, timestamp, value}
	serviceEvaluationJSONAsBytes, _ := json.Marshal(serviceEvaluation)

	// === Save Service Evaluation to state ===
	stub.PutState(evaluationId, serviceEvaluationJSONAsBytes)

	return serviceEvaluation, nil
}

// =====================================================================================================================
// Create Executed Service Transaction(Tx) Index - to do query based on Executed Service Tx Id
// =====================================================================================================================
func CreateServiceTxIndex(activity *Activity, stub shim.ChaincodeStubInterface) (serviceTxIndexKey string, err error) {
	indexName := "serviceTx~evaluation"
	serviceTxIndexKey, err = stub.CreateCompositeKey(indexName, []string{activity.ExecutedServiceTxid, activity.EvaluationId})
	if err != nil {
		return serviceTxIndexKey, err
	}
	return serviceTxIndexKey, nil
}

// ============================================================================================================================
// Create Demander Agent - Executer Agent - Timestamp - Evaluation Id Index - to do query based on Demander-Executer-Timestamp Evaluations
// ============================================================================================================================
func CreateDemanderExecuterTimestampIndex(activity *Activity, stub shim.ChaincodeStubInterface) (agentServiceIndex string, err error) {
	indexName := "demander~executer~timestamp~evaluation"
	agentServiceIndex, err = stub.CreateCompositeKey(indexName, []string{activity.DemanderAgentId, activity.ExecuterAgentId, activity.ExecutedServiceTimestamp, activity.EvaluationId})
	if err != nil {
		return agentServiceIndex, err
	}
	return agentServiceIndex, nil
}

func CheckingCreatingIndexingActivity(writerAgentId string, demanderAgentId string, executerAgentId string, executedServiceId string, executedServiceTxId string, timestamp string, value string, stub shim.ChaincodeStubInterface) (*Activity, error) {
	// ==== Check if serviceEvaluation already exists ====
	// TODO: Definire come creare evaluationId, per ora è composto dai due ID (writerAgentId + demanderAgentId + executerAgentId + ExecutedServiceTxId)
	evaluationId := writerAgentId + demanderAgentId + executerAgentId + executedServiceTxId
	serviceEvaluationAsBytes, err := stub.GetState(evaluationId)
	if err != nil {
		return nil, errors.New("Failed to get executedService demanderAgent relation: " + err.Error())
	} else if serviceEvaluationAsBytes != nil {
		fmt.Println("This executedService demanderAgent relation already exists with relationId: " + evaluationId)
		return nil, errors.New("This executedService demanderAgent relation already exists with relationId: " + evaluationId)
	}

	// ==== Actual creation of Service Evaluation  ====
	serviceEvaluation, err := CreateActivity(evaluationId, writerAgentId, demanderAgentId, executerAgentId, executedServiceId, executedServiceTxId, timestamp, value, stub)
	if err != nil {
		return nil, errors.New("Failed to create executedService demanderAgent relation of executedService " + executedServiceId + " with demanderAgent " + executedServiceId)
	}

	// ==== Indexing of serviceEvaluation by Service Tx Id ====

	// index create
	serviceTxIndexKey, serviceIndexError := CreateServiceTxIndex(serviceEvaluation, stub)
	if serviceIndexError != nil {
		return nil, errors.New(serviceIndexError.Error())
	}
	//  Note - passing a 'nil' emptyValue will effectively delete the key from state, therefore we pass null character as emptyValue
	//  Save index entry to state. Only the key Name is needed, no need to store a duplicate copy of the ServiceAgentRelation.
	emptyValue := []byte{0x00}
	// index save
	putStateError := stub.PutState(serviceTxIndexKey, emptyValue)
	if putStateError != nil {
		return nil, errors.New("Error  saving Service index: " + putStateError.Error())
	}

	// ==== Indexing of serviceEvaluation by Agent ====

	// index create
	demanderExecuterIndexKey, agentIndexError := CreateDemanderExecuterTimestampIndex(serviceEvaluation, stub)
	if agentIndexError != nil {
		return nil, errors.New(agentIndexError.Error())
	}
	// index save
	putStateDemanderExecuterIndexError := stub.PutState(demanderExecuterIndexKey, emptyValue)
	if putStateDemanderExecuterIndexError != nil {
		return nil, errors.New("Error  saving Agent index: " + putStateDemanderExecuterIndexError.Error())
	}

	return serviceEvaluation, nil
}

// =====================================================================================================================
// Get Service Agent Relation - get the service agent relation asset from ledger - return (nil,nil) if not found
// =====================================================================================================================
func GetActivity(stub shim.ChaincodeStubInterface, evaluationId string) (Activity, error) {
	var serviceRelationAgent Activity
	serviceRelationAgentAsBytes, err := stub.GetState(evaluationId) // getState retreives a key/value from the ledger
	if err != nil { // this seems to always succeed, even if key didn't exist
		return serviceRelationAgent, errors.New("Error in finding service relation with agent: " + error.Error(err))
	}

	json.Unmarshal(serviceRelationAgentAsBytes, &serviceRelationAgent) // un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo Activity?)

	return serviceRelationAgent, nil
}

// =====================================================================================================================
// Get Service Agent Relation Not Found Error - get the service agent relation asset from ledger - throws error if not found (error!=nil ---> key not found)
// =====================================================================================================================
func GetActivityNotFoundError(stub shim.ChaincodeStubInterface, evaluationId string) (Activity, error) {
	var serviceRelationAgent Activity
	serviceRelationAgentAsBytes, err := stub.GetState(evaluationId) // getState retreives a key/value from the ledger
	if err != nil { // this seems to always succeed, even if key didn't exist
		return serviceRelationAgent, errors.New("Error in finding service evaluation: " + error.Error(err))
	}

	if serviceRelationAgentAsBytes == nil {
		return Activity{}, errors.New("Service Evaluation non found, EvaluationId: " + evaluationId)
	}
	json.Unmarshal(serviceRelationAgentAsBytes, &serviceRelationAgent) // un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo Activity)

	return serviceRelationAgent, nil
}

// =====================================================================================================================
// Get the service query on ServiceRelationAgent - Execute the query based on service composite index
// =====================================================================================================================
func GetByExecutedServiceTx(executedServiceTxId string, stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, error) {
	// Query the service~agent~relation index by service
	// This will execute a key range query on all keys starting with 'service'
	indexName := "serviceTx~evaluation"
	executedServiceTxResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{executedServiceTxId})
	if err != nil {
		return executedServiceTxResultsIterator, err
	}
	return executedServiceTxResultsIterator, nil
}

// =====================================================================================================================
// Get the agent query on ServiceRelationAgent - Execute the query based on agent composite index
// =====================================================================================================================
func GetByDemanderExecuterTimestamp(demanderAgentId string, executerAgentId string, timestamp string, stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, error) {
	// Query the service~agent~relation index by service
	// This will execute a key range query on all keys starting with 'service'
	indexName := "demander~executer~timestamp~evaluation"
	demanderExecuterResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{demanderAgentId, executerAgentId, timestamp})
	if err != nil {
		return demanderExecuterResultsIterator, err
	}
	return demanderExecuterResultsIterator, nil
}

// =====================================================================================================================
// Delete Service Evaluation - "removing"" a key/value from the ledger
// =====================================================================================================================
func DeleteServiceEvaluation(stub shim.ChaincodeStubInterface, evaluationId string) error {
	// remove the serviceRelationAgent
	err := stub.DelState(evaluationId) // remove the key from chaincode state
	if err != nil {
		return err
	}
	return nil
}

// =====================================================================================================================
// Delete Executed Service Tx Index - "removing"" a key/value from the ledger
// =====================================================================================================================
func DeleteExecutedServiceTxIndex(stub shim.ChaincodeStubInterface, executedServiceTxId string, evaluationId string) error {
	// remove the serviceRelationAgent
	indexName := "serviceTx~evaluation"

	agentServiceIndex, err := stub.CreateCompositeKey(indexName, []string{executedServiceTxId, evaluationId})
	if err != nil {
		return err
	}
	err = stub.DelState(agentServiceIndex) // remove the key from chaincode state
	if err != nil {
		return err
	}
	return nil
}

// =====================================================================================================================
// Delete Agent Service Relation - delete from state and from marble index Shows Off DelState() - "removing"" a key/value from the ledger
// =====================================================================================================================
func DeleteDemanderExecuterIndex(stub shim.ChaincodeStubInterface, demanderAgentId string, executerAgentId string, evaluationId string) error {
	// remove the serviceRelationAgent
	indexName := "demander~executer~evaluation"

	agentServiceIndex, err := stub.CreateCompositeKey(indexName, []string{demanderAgentId, executerAgentId, evaluationId})
	if err != nil {
		return err
	}
	err = stub.DelState(agentServiceIndex) // remove the key from chaincode state
	if err != nil {
		return err
	}
	return nil
}

// =====================================================================================================================
// GetServiceRelationSliceFromServiceTxRangeQuery - Get the Activity Slices from the result of query "GetByExecutedServiceTx"
// =====================================================================================================================
func GetActivitySliceFromServiceTxIdRangeQuery(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) ([]Activity, error) {
	var serviceEvaluations []Activity
	defer queryIterator.Close()

	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return nil, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		evaluationId := compositeKeyParts[1]

		iserviceRelationAgent, err := GetActivity(stub, evaluationId)
		serviceEvaluations = append(serviceEvaluations, iserviceRelationAgent)
		if err != nil {
			return nil, err
		}
		fmt.Printf("- found a relation EVALUATION ID: %s \n", evaluationId)
	}
	return serviceEvaluations, nil
}

// =====================================================================================================================
// GetActivitySliceFromDemanderExecuterTimestampRangeQuery - Get the Agent and Activity Slices from the result of query "GetByDemanderExecuterTimestamp"
// =====================================================================================================================
func GetActivitySliceFromDemanderExecuterTimestampRangeQuery(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) ([]Activity, error) {
	var serviceEvaluations []Activity
	// USE DEFER BECAUSE it will close also in case of error throwing (premature return)
	defer queryIterator.Close()

	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return nil, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		evaluationId := compositeKeyParts[3]

		iserviceRelationAgent, err := GetActivity(stub, evaluationId)
		serviceEvaluations = append(serviceEvaluations, iserviceRelationAgent)
		if err != nil {
			return nil, err
		}
		fmt.Printf("- found a relation EVALUATION ID: %s , VALUE: %s\n", iserviceRelationAgent.EvaluationId, iserviceRelationAgent.Value)
	}
	return serviceEvaluations, nil
}

// =====================================================================================================================
// Print Service Tx Results Iterator - Print on screen the iterator of the executed service tx id query result
// =====================================================================================================================
func PrintByExecutedServiceTxIdResultsIterator(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) error {
	// USE DEFER BECAUSE it will close also in case of error throwing (premature return)
	defer queryIterator.Close()
	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return err
		}
		// get the service agent relation from service~agent~relation composite key
		indexName, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		executedServiceTxId := compositeKeyParts[0]
		evaluationId := compositeKeyParts[1]

		if err != nil {
			return err
		}
		fmt.Printf("- found a relation from OBJECT_TYPE:%s EXECUTED SERVICE TX ID:%s EVALUATION ID: %s\n", indexName, executedServiceTxId, evaluationId)
	}
	return nil
}

// =====================================================================================================================
// Print Demander Executer Results Iterator - Print on screen the general iterator of the demander executer index query result
// =====================================================================================================================
func PrintByDemanderExecuterTimestampResultsIterator(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) error {
	defer queryIterator.Close()
	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return err
		}
		indexName, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		demanderAgentId := compositeKeyParts[0]
		executerAgentId := compositeKeyParts[1]
		evaluationId := compositeKeyParts[3]

		if err != nil {
			return err
		}
		fmt.Printf("- found a relation from OBJECT_TYPE:%s Demander AGENT ID:%s Executer AGENT ID:%s  EVALUATION ID: %s\n", indexName, demanderAgentId, executerAgentId, evaluationId)
	}
	return nil
}

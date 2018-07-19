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

// ===================================================================================
// Define the Service Evaluation structure
// ===================================================================================
// - ReputationId
// - AgentId
// - ServiceId
// - AgentRole
// - ExecutedServiceId
// - ExecutedServiceTxId
// - Timestamp
// - Outcome
// - Value
// - IsFinalEvaluation
// UNIVOCAL: WriterAgentId, DemanderAgentId, ExecuterAgentId, ExecutedServiceTxId
type ServiceEvaluation struct {
	EvaluationId        string `json:"EvaluationId"`
	WriterAgentId       string `json:"WriterAgentId"`// WriterAgentId = DemanderAgentId || ExecuterAgentId
	DemanderAgentId     string `json:"DemanderAgentId"`
	ExecuterAgentId     string `json:"ExecuterAgentId"`
	ExecutedServiceId   string `json:"ExecutedServiceId"`
	ExecutedServiceTxid string `json:"ExecutedServiceTxid"` //Relativo all'esecuzione del servizio
	Timestamp           string `json:"Timestamp"`
	Value               string `json:"Value"`
	// Outcome             string `json:"Outcome"` // TODO: Da levare
	// IsFinalEvaluation   string `json:"IsFinalEvaluation"` // TODO: Da levare
}

// ============================================================
// Create Service Evaluation - create a new service evaluation
// ============================================================
func CreateServiceEvaluation(evaluationId string, writerAgentId string, demanderAgentId string, executerAgentId string, executedServiceId string, executedServiceTxId string, timestamp string, value string, stub shim.ChaincodeStubInterface) (*ServiceEvaluation, error) {
	// ==== Create marble object and marshal to JSON ====
	serviceEvaluation := &ServiceEvaluation{evaluationId, writerAgentId, demanderAgentId, executerAgentId, executedServiceId, executedServiceTxId,timestamp,value}
	serviceEvaluationJSONAsBytes, _ := json.Marshal(serviceEvaluation)

	// === Save Service Evaluation to state ===
	stub.PutState(evaluationId, serviceEvaluationJSONAsBytes)

	return serviceEvaluation, nil
}

// ============================================================================================================================
// Create Executed Service Transaction(Tx) Index - to do query based on Executed Service Tx Id
// ============================================================================================================================
func CreateServiceTxIndex(serviceEvaluation *ServiceEvaluation, stub shim.ChaincodeStubInterface) (serviceTxIndexKey string, err error) {
	indexName := "serviceTx~evaluation"
	serviceTxIndexKey, err = stub.CreateCompositeKey(indexName, []string{serviceEvaluation.ExecutedServiceTxid, serviceEvaluation.EvaluationId})
	if err != nil {
		return serviceTxIndexKey, err
	}
	return serviceTxIndexKey, nil
}

// ============================================================================================================================
// Create Demander Agent - Executer Agent - Evaluation Id Index - to do query based on Demander-Executer Evaluations
// ============================================================================================================================
func CreateDemanderExecuterIndex(serviceEvaluation *ServiceEvaluation, stub shim.ChaincodeStubInterface) (agentServiceIndex string, err error) {
	indexName := "demander~executer~evaluation"
	agentServiceIndex, err = stub.CreateCompositeKey(indexName, []string{serviceEvaluation.DemanderAgentId, serviceEvaluation.ExecuterAgentId, serviceEvaluation.EvaluationId})
	if err != nil {
		return agentServiceIndex, err
	}
	return agentServiceIndex, nil
}

// ============================================================================================================================
// Get Service Agent Relation - get the service agent relation asset from ledger - return (nil,nil) if not found
// ============================================================================================================================
func GetServiceEvaluation(stub shim.ChaincodeStubInterface, evaluationId string) (ServiceEvaluation, error) {
	var serviceRelationAgent ServiceEvaluation
	serviceRelationAgentAsBytes, err := stub.GetState(evaluationId) //getState retreives a key/value from the ledger
	if err != nil {                                               //this seems to always succeed, even if key didn't exist
		return serviceRelationAgent, errors.New("Error in finding service relation with agent: " + error.Error(err))
	}

	json.Unmarshal(serviceRelationAgentAsBytes, &serviceRelationAgent) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo ServiceEvaluation?)

	return serviceRelationAgent, nil
}

// ============================================================================================================================
// Get Service Agent Relation Not Found Error - get the service agent relation asset from ledger - throws error if not found (error!=nil ---> key not found)
// ============================================================================================================================
func GetServiceEvaluationNotFoundError(stub shim.ChaincodeStubInterface, evaluationId string) (ServiceEvaluation, error) {
	var serviceRelationAgent ServiceEvaluation
	serviceRelationAgentAsBytes, err := stub.GetState(evaluationId) //getState retreives a key/value from the ledger
	if err != nil {                                               //this seems to always succeed, even if key didn't exist
		return serviceRelationAgent, errors.New("Error in finding service evaluation: " + error.Error(err))
	}

	if serviceRelationAgentAsBytes == nil {
		return ServiceEvaluation{}, errors.New("Service Evaluation non found, EvaluationId: " + evaluationId)
	}
	json.Unmarshal(serviceRelationAgentAsBytes, &serviceRelationAgent) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo ServiceEvaluation)

	return serviceRelationAgent, nil
}

// ============================================================================================================================
// Get the service query on ServiceRelationAgent - Execute the query based on service composite index
// ============================================================================================================================
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

// ============================================================================================================================
// Get the agent query on ServiceRelationAgent - Execute the query based on agent composite index
// ============================================================================================================================
func GetByDemanderExecuter(demanderAgentId string, executerAgentId string, stub shim.ChaincodeStubInterface) (shim.StateQueryIteratorInterface, error) {
	// Query the service~agent~relation index by service
	// This will execute a key range query on all keys starting with 'service'
	indexName := "demander~executer~evaluation"
	demanderExecuterResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{demanderAgentId, executerAgentId})
	if err != nil {
		return demanderExecuterResultsIterator, err
	}
	return demanderExecuterResultsIterator, nil
}

// ============================================================================================================================
// Delete Service Evaluation - "removing"" a key/value from the ledger
// ============================================================================================================================
func DeleteServiceEvaluation(stub shim.ChaincodeStubInterface, evaluationId string) error {
	// remove the serviceRelationAgent
	err := stub.DelState(evaluationId) //remove the key from chaincode state
	if err != nil {
		return err
	}
	return nil
}

// ============================================================================================================================
// Delete Executed Service Tx Index - "removing"" a key/value from the ledger
// ============================================================================================================================
func DeleteExecutedServiceTxIndex(stub shim.ChaincodeStubInterface, executedServiceTxId string, evaluationId string) error {
	// remove the serviceRelationAgent
	indexName := "serviceTx~evaluation"

	agentServiceIndex, err := stub.CreateCompositeKey(indexName, []string{executedServiceTxId, evaluationId})
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
func DeleteDemanderExecuterIndex(stub shim.ChaincodeStubInterface, demanderAgentId string, executerAgentId string, evaluationId string) error {
	// remove the serviceRelationAgent
	indexName := "demander~executer~evaluation"

	agentServiceIndex, err := stub.CreateCompositeKey(indexName, []string{demanderAgentId, executerAgentId, evaluationId})
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
// GetServiceRelationSliceFromServiceTxRangeQuery - Get the ServiceEvaluation Slices from the result of query "GetByExecutedServiceTx"
// ============================================================================================================================
func GetServiceEvaluationSliceFromServiceTxIdRangeQuery(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) ([]ServiceEvaluation, error) {
	var serviceEvaluations []ServiceEvaluation
	defer queryIterator.Close()

	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return nil, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		evaluationId := compositeKeyParts[1]

		iserviceRelationAgent, err := GetServiceEvaluation(stub, evaluationId)
		serviceEvaluations = append(serviceEvaluations, iserviceRelationAgent)
		if err != nil {
			return nil, err
		}
		fmt.Printf("- found a relation EVALUATION ID: %s \n", evaluationId)
	}
	return serviceEvaluations, nil
}

// ============================================================================================================================
// GetServiceEvaluationSliceFromDemanderExecuterRangeQuery - Get the Agent and ServiceEvaluation Slices from the result of query "GetByDemanderExecuter"
// ============================================================================================================================
func GetServiceEvaluationSliceFromDemanderExecuterRangeQuery(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) ([]ServiceEvaluation, error) {
	var serviceEvaluations []ServiceEvaluation
	// USE DEFER BECAUSE it will close also in case of error throwing (premature return)
	defer queryIterator.Close()

	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return nil, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		evaluationId := compositeKeyParts[2]

		iserviceRelationAgent, err := GetServiceEvaluation(stub, evaluationId)
		serviceEvaluations = append(serviceEvaluations, iserviceRelationAgent)
		if err != nil {
			return nil, err
		}
		fmt.Printf("- found a relation EVALUATION ID: %s \n", evaluationId)
	}
	return serviceEvaluations, nil
}

// ============================================================================================================================
// Print Service Tx Results Iterator - Print on screen the iterator of the executed service tx id query result
// ============================================================================================================================
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

// ============================================================================================================================
// Print Demander Executer Results Iterator - Print on screen the general iterator of the demander executer index query result
// ============================================================================================================================
func PrintByDemanderExecuterResultsIterator(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) error {
	defer queryIterator.Close()
	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return err
		}
		indexName, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)

		demanderAgentId := compositeKeyParts[0]
		executerAgentId := compositeKeyParts[1]
		evaluationId := compositeKeyParts[2]

		if err != nil {
			return err
		}
		fmt.Printf("- found a relation from OBJECT_TYPE:%s DEMANDER AGENT ID:%s EXECUTER AGENT ID:%s  EVALUATION ID: %s\n", indexName, demanderAgentId, executerAgentId, evaluationId)
	}
	return nil
}


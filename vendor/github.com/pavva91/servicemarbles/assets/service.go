/*
Package assets represent the assets with relatives base functions (create, indexing, queries) that can be stored in the ledger of hyperledger fabric blockchain
 */
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

// ===================================================================================
// Define the Service structure, with 3 properties. trying(https://medium.com/@wishmithasmendis/from-rdbms-to-key-value-store-data-modeling-techniques-a2874906bc46)
// ===================================================================================
// - ServiceId
// - Name
// - Description
type Service struct {
	ServiceId   string `json:"ServiceId"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

// ============================================================
// CreateService - create a new service and return the created agent
// ============================================================
func CreateService(serviceId string, serviceName string, serviceDescription string, stub shim.ChaincodeStubInterface) (*Service, error) {
	// ==== Create marble object and marshal to JSON ====
	service := &Service{serviceId, serviceName, serviceDescription}
	service2JSONAsBytes, err := json.Marshal(service)
	if err != nil {
		return service, errors.New("Failed Marshal service: " + service.Name)
	}

	// === Save marble to state ===
	stub.PutState(serviceId, service2JSONAsBytes)
	return service, nil
}

// ============================================================================================================================
// Create Service's Name based Index - to do query based on Name of the Service
// ============================================================================================================================
func CreateNameIndex(serviceToIndex *Service, stub shim.ChaincodeStubInterface) (nameServiceIndexKey string, err error) {
	//  ==== Index the serviceAgentRelation to enable service-based range queries, e.g. return all x services ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on service~agent~relation.
	//  This will enable very efficient state range queries based on composite keys matching service~agent~relation
	indexName := "Name~ServiceId"
	nameServiceIndexKey, err = stub.CreateCompositeKey(indexName, []string{serviceToIndex.Name, serviceToIndex.ServiceId})
	if err != nil {
		return nameServiceIndexKey, err
	}
	return nameServiceIndexKey, nil
}

// ============================================================================================================================
// Create Service and create and save the index - Atomic function of 3 the subfunctions: save, index, saveindex
// ============================================================================================================================
func CreateAndIndexService(serviceId string, serviceName string, serviceDescription string, stub shim.ChaincodeStubInterface) error {
	service, err := CreateService(serviceId, serviceName, serviceDescription, stub)
	if err != nil {
		return errors.New("Failed to create the service: " + err.Error())
	}

	// ==== Indexing of service by Name (to do query by Name, if you want) ====
	// index create
	nameIndexKey, nameIndexError := CreateNameIndex(service, stub)
	if nameIndexError != nil {
		return errors.New(nameIndexError.Error())
	}
	fmt.Println(nameIndexKey)
	// TODO: Mettere a Posto (fare un create Service index

	saveIndexError :=  SaveIndex(nameIndexKey, stub)
	if saveIndexError != nil {
		return errors.New(saveIndexError.Error())
	}
	return nil

}

// ============================================================================================================================
// Get Service - get the service asset from ledger - return (nil,nil) if not found
// ============================================================================================================================
func GetService(stub shim.ChaincodeStubInterface, serviceId string) (Service, error) {
	var service Service
	serviceAsBytes, err := stub.GetState(serviceId) //getState retreives a key/value from the ledger
	if err != nil {                                 //this seems to always succeed, even if key didn't exist
		return service, errors.New("Error in finding service: " + error.Error(err))
	}
	fmt.Println(serviceAsBytes)
	fmt.Println(service)

	json.Unmarshal(serviceAsBytes, &service) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo Service)

	fmt.Println(service)
	return service, nil
}

// ============================================================================================================================
// Get Service Not Found Error - get the service asset from ledger - throws error if not found (error!=nil ---> key not found)
// ============================================================================================================================
func GetServiceNotFoundError(stub shim.ChaincodeStubInterface, serviceId string) (Service, error) {
	var service Service
	serviceAsBytes, err := stub.GetState(serviceId) //getState retreives a key/value from the ledger
	if err != nil {                                 //this seems to always succeed, even if key didn't exist
		return service, errors.New("Error in finding service: " + error.Error(err))
	}
	fmt.Println(serviceAsBytes)
	fmt.Println(service)
	// TODO: Levare trigger error ma gestire il payload null

	if serviceAsBytes == nil {
		return service, errors.New("Service non found, ServiceId: " + serviceId)
	}
	json.Unmarshal(serviceAsBytes, &service) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo Service)

	fmt.Println(service)
	return service, nil
}

// ============================================================================================================================
// Get Service as Bytes - get the service as bytes from ledger
// ============================================================================================================================
func GetServiceAsBytes(stub shim.ChaincodeStubInterface, idService string) ([]byte, error) {
	serviceAsBytes, err := stub.GetState(idService) //getState retreives a key/value from the ledger
	if err != nil {                                 //this seems to always succeed, even if key didn't exist
		return serviceAsBytes, errors.New("Failed to get service - " + idService)
	}
	return serviceAsBytes, nil
}

// ============================================================================================================================
// DeleteService() - remove a service from state and from service index
//
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs:
//      0
//     ServiceId
// ============================================================================================================================
func DeleteService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting delete_marble")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// input sanitation
	err := arglib.SanitizeArguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	serviceId := args[0]

	// get the service
	service, err := GetServiceNotFoundError(stub, serviceId)
	if err != nil {
		fmt.Println("Failed to find service by ServiceId " + serviceId)
		return shim.Error(err.Error())
	}

	// TODO: Delete anche (prima) le relazioni del servizio con gli agenti
	err = DeleteAllServiceAgentRelations(serviceId, stub)
	if err != nil {
		return shim.Error("Failed to delete service agent relation: " + err.Error())
	}

	// remove the service
	err = stub.DelState(serviceId) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete service: " + err.Error())
	}

	fmt.Println("Deleted service: " + service.Name)
	return shim.Success(nil)
}

// ============================================================
// DeleteAllServiceAgentRelations - delete all the Service relations with agent (aka: Reference Integrity)
// ============================================================
func DeleteAllServiceAgentRelations(serviceId string, stub shim.ChaincodeStubInterface) error {
	serviceAgentResultsIterator, err := GetByService(serviceId, stub)
	if err != nil {
		return err
	}
	for i := 0; serviceAgentResultsIterator.HasNext(); i++ {
		responseRange, err := serviceAgentResultsIterator.Next()
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

		fmt.Printf("Delete the relation: from composite key OBJECT_TYPE:%s SERVICE ID:%s AGENT ID:%s RELATION ID: %s\n", objectType, serviceId, agentId, relationId)

		// remove the serviceRelationAgent
		err = DeleteServiceAgentRelation(stub, relationId) //remove the key from chaincode state
		if err != nil {
			return err
		}

		// remove the service index
		err = deleteServiceIndex(stub, objectType, serviceId, agentId, relationId) //remove the key from chaincode state
		if err != nil {
			return err
		}

	}
	return nil
}

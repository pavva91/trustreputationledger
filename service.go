package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"errors"
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

// ============================================================================================================================
// Init Service - wrapper of createService called from the chaincode invoke
// ============================================================================================================================
func initService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1                 2
	// "ServiceId", "serviceName", "serviceDescription"
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	fmt.Println("- start init service")

	// ==== Input sanitation ====
	sanitizeError := sanitize_arguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]
	serviceName := args[1]
	serviceDescription := args[2]

	// ==== Check if service already exists ====
	serviceAsBytes, err := getServiceAsBytes(stub,serviceId)
	if err != nil {
		return shim.Error("Failed to get service: " + err.Error())
	} else if serviceAsBytes != nil {
		fmt.Println("This service already exists: " + serviceName)
		return shim.Error("This service already exists: " + serviceName)
	}

	service,err := createService(serviceId, serviceName, serviceDescription, stub)
	if err != nil {
		return shim.Error("Failed to create the service: " + err.Error())
	}

	// ==== Indexing of service by Name (to do query by Name, if you want) ====
	// index create
	nameIndexKey, nameIndexError := createNameIndex(service,stub)
	if nameIndexError != nil {
		return shim.Error(nameIndexError.Error())
	}
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	//  Save index entry to state. Only the key Name is needed, no need to store a duplicate copy of the marble.
	value := []byte{0x00}
	// index save
	putStateError:=stub.PutState(nameIndexKey, value)
	if putStateError != nil {
		return shim.Error(putStateError.Error())
	}

	// ==== Service saved and indexed. Return success ====
	fmt.Println("Servizio: " + service.Name + " creato - end init service")
	return shim.Success(nil)
}
// ============================================================================================================================
// Query Service - wrapper of getService called from the chaincode invoke
// ============================================================================================================================
func queryService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ServiceId"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// ==== Input sanitation ====
	sanitizeError := sanitize_arguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== get the service ====
	service, err := getService(stub, serviceId)
	if err != nil{
		fmt.Println("Failed to find service by id " + serviceId)
		return shim.Error(err.Error())
	}else {
		fmt.Println("Service ID: " + service.ServiceId +", Service: " + service.Name + ", with Description: " + service.Description + " found")
		// ==== Marshal the byService query result ====
		serviceAsJSON, err := json.Marshal(service)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(serviceAsJSON)
	}
}

// ============================================================
// createService - create a new service and return the created agent
// ============================================================
func createService(serviceId string, serviceName string, serviceDescription string, stub shim.ChaincodeStubInterface) (*Service,error) {
	// ==== Create marble object and marshal to JSON ====
	service := &Service{serviceId, serviceName, serviceDescription}
	service2JSONAsBytes, err := json.Marshal(service)
	if err != nil {
		return service, errors.New("Failed Marshal service: " + service.Name)
	}

	// === Save marble to state ===
	stub.PutState(serviceId, service2JSONAsBytes)
	return service,nil
}

// ============================================================================================================================
// Create Service's Name based Index - to do query based on Name of the Service
// ============================================================================================================================
func createNameIndex(serviceToIndex *Service, stub shim.ChaincodeStubInterface)  (nameServiceIndexKey string, err error){
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
	return nameServiceIndexKey,nil
}

// ============================================================================================================================
// Get Service - get the service asset from ledger
// ============================================================================================================================
func getService(stub shim.ChaincodeStubInterface, idService string) (Service, error) {
	var service Service
	serviceAsBytes, err := stub.GetState(idService) //getState retreives a key/value from the ledger
	if err != nil {                                            //this seems to always succeed, even if key didn't exist
		return service, errors.New("Failed to get service - " + idService)
	}
	json.Unmarshal(serviceAsBytes, &service) //un stringify it aka JSON.parse()

	// TODO: Inserire controllo di tipo (Verificare sia di tipo Service)

	return service, nil
}

// ============================================================================================================================
// Get Service as Bytes - get the service as bytes from ledger
// ============================================================================================================================
func getServiceAsBytes(stub shim.ChaincodeStubInterface, idService string) ([]byte, error) {
	serviceAsBytes, err := stub.GetState(idService) //getState retreives a key/value from the ledger
	if err != nil {                                            //this seems to always succeed, even if key didn't exist
		return serviceAsBytes, errors.New("Failed to get service - " + idService)
	}
	return serviceAsBytes, nil
}

// ============================================================================================================================
// deleteService() - remove a service from state and from service index
//
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs:
//      0
//     ServiceId
// ============================================================================================================================
func deleteService(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting delete_marble")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	serviceId := args[0]

	// get the service
	service, err := getService(stub, serviceId)
	if err != nil{
		fmt.Println("Failed to find service by ServiceId " + serviceId)
		return shim.Error(err.Error())
	}

	// TODO: Delete anche (prima) le relazioni del servizio con gli agenti
	err=deleteAllServiceAgentRelations(serviceId,stub)
	if err != nil {
		return shim.Error("Failed to delete service agent relation: "+ err.Error())
	}

	// remove the service
	err = stub.DelState(serviceId) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete service: "+ err.Error())
	}

	fmt.Println("Deleted service: " + service.Name)
	return shim.Success(nil)
}

// ============================================================
// deleteAllServiceAgentRelations - delete all the Service relations with agent (aka: Reference Integrity)
// ============================================================
func deleteAllServiceAgentRelations(serviceId string, stub shim.ChaincodeStubInterface) error{
	serviceAgentResultsIterator, err := getByService(serviceId,stub)
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

		serviceId:=compositeKeyParts[0]
		agentId:=compositeKeyParts[1]
		relationId:=compositeKeyParts[2]

		if err != nil {
			return err
		}

		fmt.Printf("Delete the relation: from composite key OBJECT_TYPE:%s SERVICE ID:%s AGENT ID:%s RELATION ID: %s\n", objectType, serviceId, agentId, relationId)

		// remove the serviceRelationAgent
		err = deleteServiceAgentRelation(stub, relationId) //remove the key from chaincode state
		if err != nil {
			return err
		}

		// remove the service index
		err = deleteServiceIndex(stub,objectType,serviceId,agentId,relationId) //remove the key from chaincode state
		if err != nil {
			return err
		}



		}
	return nil
}

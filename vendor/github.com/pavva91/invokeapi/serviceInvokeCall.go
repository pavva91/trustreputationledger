/*
Package invokeapi is the middle layer between the Chaincode entry point (main package) and the Assets (assets package)
that is called directly from the chaincode's Invoke funtions and aggregate the calls to the assets to follow the
"business logic"
*/
/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
*/
package invokeapi

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pavva91/arglib"
	"strconv"

	// a "github.com/pavva91/trustreputationledger/assets"
	a "github.com/pavva91/assets"
)

var serviceInvokeCallLog = shim.NewLogger("serviceInvokeCall")
// =====================================================================================================================
// Create Leaf Service - wrapper of CreateLeafService called from the chaincode invoke
// =====================================================================================================================
func CreateLeafService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1                 2
	// "ServiceId", "serviceName", "serviceDescription"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 3)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}
	serviceInvokeCallLog.Info("- start create Leaf service")

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]
	serviceName := args[1]
	serviceDescription := args[2]

	// ==== Check if service already exists ====
	serviceAsBytes, err := a.GetServiceAsBytes(stub, serviceId)
	if err != nil {
		return shim.Error("Failed to get service: " + err.Error())
	} else if serviceAsBytes != nil {
		serviceInvokeCallLog.Info("This service already exists: " + serviceId)
		return shim.Error("This service already exists: " + serviceId)
	}

	service, err := a.CreateLeafService(serviceId, serviceName, serviceDescription, stub)
	if err != nil {
		return shim.Error("Failed to create the service: " + err.Error())
	}

	// ==== Indexing of service by Name (to do query by Name, if you want) ====
	// index create
	nameIndexKey, nameIndexError := a.CreateNameIndex(service, stub)
	if nameIndexError != nil {
		return shim.Error(nameIndexError.Error())
	}
	serviceInvokeCallLog.Info(nameIndexKey)

	// index save
	saveIndexError := a.SaveIndex(nameIndexKey, stub)
	if saveIndexError != nil {
		return shim.Error(saveIndexError.Error())
	}

	// ==== Service saved and indexed. Set Event ====

	eventPayload:="Created Service: " + serviceId
	payloadAsBytes := []byte(eventPayload)
	eventError := stub.SetEvent("ServiceCreatedEvent",payloadAsBytes)
	if eventError != nil {
		serviceInvokeCallLog.Error("Error in event Creation: " + eventError.Error())
	}else {
		serviceInvokeCallLog.Info("Event Create Service OK")
	}
	// ==== Service saved and indexed. Return success ====
	serviceInvokeCallLog.Info("ServiceId: " + service.ServiceId + ", Name: " + service.Name + ", Description: " + service.Description + " Succesfully Created - End Create Service")
	return shim.Success(nil)
}
// =====================================================================================================================
// Create Composite Service - wrapper of CreateLeafService called from the chaincode invoke
// =====================================================================================================================
func CreateCompositeService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1                 2               3
	// "ServiceId", "serviceName", "serviceDescription", "serviceComposition"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 4)
	if argumentSizeError != nil {
		serviceInvokeCallLog.Error(argumentSizeError.Error())
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		serviceInvokeCallLog.Error(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]
	serviceName := args[1]
	serviceDescription := args[2]
	serviceCompositionAsString := args[3]

	serviceComposition := arglib.ParseStringToStringSlice(serviceCompositionAsString)

	serviceInvokeCallLog.Info("- Start Create Composite service: " + serviceId)

	// ==== Check if service already exists ====
	serviceAsBytes, err := a.GetServiceAsBytes(stub, serviceId)
	if err != nil {
		serviceInvokeCallLog.Error(err.Error())
		return shim.Error("Failed to get service: " + err.Error())
	} else if serviceAsBytes != nil {
		serviceInvokeCallLog.Error("This service already exists: " + serviceId)
		return shim.Error("This service already exists: " + serviceId)
	}

	service, err := a.CreateCompositeService(serviceId, serviceName, serviceDescription, serviceComposition, stub)
	if err != nil {
		serviceInvokeCallLog.Error(err.Error())
		return shim.Error("Failed to create the service: " + err.Error())
	}

	// ==== Indexing of service by Name (to do query by Name, if you want) ====
	// index create
	nameIndexKey, nameIndexError := a.CreateNameIndex(service, stub)
	if nameIndexError != nil {
		serviceInvokeCallLog.Error(nameIndexError.Error())
		return shim.Error(nameIndexError.Error())
	}
	serviceInvokeCallLog.Info("nameIndexKey: " + nameIndexKey)

	// index save
	saveIndexError := a.SaveIndex(nameIndexKey, stub)
	if saveIndexError != nil {
		serviceInvokeCallLog.Error(saveIndexError.Error())
		return shim.Error(saveIndexError.Error())
	}

	// ==== Service saved and indexed. Set Event ====

	eventPayload:="Created Service: " + serviceId
	payloadAsBytes := []byte(eventPayload)
	eventError := stub.SetEvent("ServiceCreatedEvent",payloadAsBytes)
	if eventError != nil {
		serviceInvokeCallLog.Error("Error in event Creation: " + eventError.Error())
	}else {
		serviceInvokeCallLog.Info("Event Create Service OK")
	}
	// ==== Service saved and indexed. Return success ====
	serviceInvokeCallLog.Info("ServiceId: " + service.ServiceId + ", Name: " + service.Name + ", Description: " + service.Description + " Succesfully Created - End Create Service")
	return shim.Success(nil)
}
// =====================================================================================================================
// Create Service - wrapper of CreateLeafService called from the chaincode invoke
// =====================================================================================================================
func CreateService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1                 2               3
	// "ServiceId", "serviceName", "serviceDescription", "serviceComposition"
	fourArgumentSizeError := arglib.ArgumentSizeVerification(args, 4)
	if fourArgumentSizeError != nil {
		return CreateLeafService(stub,args)
	}else{

		serviceId := args[0]
		serviceName := args[1]
		serviceDescription := args[2]
		serviceCompositionAsString := args[3]

		serviceComposition := arglib.ParseStringToStringSlice(serviceCompositionAsString)

		// ==== Check if service already exists ====
		serviceAsBytes, err := a.GetServiceAsBytes(stub, serviceId)
		if err != nil {
			return shim.Error("Failed to get service: " + err.Error())
		} else if serviceAsBytes != nil {
			serviceInvokeCallLog.Info("This service already exists: " + serviceId)
			return shim.Error("This service already exists: " + serviceId)
		}

		service, err := a.CreateService(serviceId, serviceName, serviceDescription, serviceComposition, stub)
		if err != nil {
			return shim.Error("Failed to create the service: " + err.Error())
		}

		// ==== Indexing of service by Name (to do query by Name, if you want) ====
		// index create
		nameIndexKey, nameIndexError := a.CreateNameIndex(service, stub)
		if nameIndexError != nil {
			return shim.Error(nameIndexError.Error())
		}
		serviceInvokeCallLog.Info("nameIndexKey: " + nameIndexKey)

		// index save
		saveIndexError := a.SaveIndex(nameIndexKey, stub)
		if saveIndexError != nil {
			return shim.Error(saveIndexError.Error())
		}

		// ==== Service saved and indexed. Set Event ====
		transientMap, err := stub.GetTransient()
		transientData, ok := transientMap["event"]
		serviceInvokeCallLog.Info("OK: " + strconv.FormatBool(ok))
		serviceInvokeCallLog.Info(transientMap)
		serviceInvokeCallLog.Info(transientData)
		// eventError := stub.SetEvent("ServiceCreatedEvent", transientData)

		// TODO: Meaningful Event Payload
		eventPayload:="Created Service: " + serviceId
		payloadAsBytes := []byte(eventPayload)
		eventError := stub.SetEvent("ServiceCreatedEvent",payloadAsBytes)
		if eventError != nil {
			serviceInvokeCallLog.Error("Error in event Creation: " + eventError.Error())
			return shim.Error(eventError.Error())
		}else {
			serviceInvokeCallLog.Info("Event Create Service OK")
		}
		// ==== Service saved and indexed. Return success ====
		serviceInvokeCallLog.Info("ServiceId: " + service.ServiceId + ", Name: " + service.Name + ", Description: " + service.Description + " Succesfully Created - End Create Service")
		return shim.Success(nil)
	}
}

// ========================================================================================================================
// Modify Service Name - wrapper of ModifyAgentAddress called from chiancode's Invoke
// ========================================================================================================================
func ModifyServiceName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0            1
	// "serviceId", "newServiceName"
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

	serviceId := args[0]
	newServiceName := args[1]

	// ==== get the service ====
	service, getError := a.GetServiceNotFoundError(stub, serviceId)
	if getError != nil {
		serviceInvokeCallLog.Info("Failed to find service by id " + serviceId)
		serviceInvokeCallLog.Error(getError.Error())
		return shim.Error(getError.Error())
	}

	// ==== modify the service ====
	modifyError := a.ModifyServiceName(service, newServiceName, stub)
	if modifyError != nil {
		serviceInvokeCallLog.Info("Failed to modify the service name: " + newServiceName)
		serviceInvokeCallLog.Error(modifyError.Error())
		return shim.Error(modifyError.Error())
	}
	serviceInvokeCallLog.Infof("Service: " + service.Name + " modified - end modify service")

	return shim.Success(nil)
}

// ========================================================================================================================
// Modify Service Description - wrapper of ModifyAgentAddress called from chiancode's Invoke
// ========================================================================================================================
func ModifyServiceDescription(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0            1
	// "serviceId", "newServiceDescription"
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

	serviceId := args[0]
	newServiceDescription := args[1]

	// ==== get the service ====
	service, getError := a.GetServiceNotFoundError(stub, serviceId)
	if getError != nil {
		serviceInvokeCallLog.Info("Failed to find service by id " + serviceId)
		serviceInvokeCallLog.Error(getError.Error())
		return shim.Error(getError.Error())
	}

	// ==== modify the service ====
	modifyError := a.ModifyServiceDescription(service, newServiceDescription, stub)
	if modifyError != nil {
		serviceInvokeCallLog.Info("Failed to modify the service description: " + newServiceDescription)
		serviceInvokeCallLog.Error(modifyError.Error())
		return shim.Error(modifyError.Error())
	}

	return shim.Success(nil)
}

// =====================================================================================================================
// Query Service Not Found Error - wrapper of GetServiceNotFoundError called from the chaincode invoke
// =====================================================================================================================
func QueryServiceNotFoundError(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ServiceId"
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

	serviceId := args[0]

	// ==== get the service ====
	service, err := a.GetServiceNotFoundError(stub, serviceId)
	if err != nil {
		serviceInvokeCallLog.Info("Failed to find service by id " + serviceId)
		serviceInvokeCallLog.Error(err.Error())
		return shim.Error(err.Error())
	} else {
		serviceInvokeCallLog.Info("Service ID: " + service.ServiceId + ", Service: " + service.Name + ", with Description: " + service.Description + " found")
		// ==== Marshal the byService query result ====
		serviceAsJSON, err := json.Marshal(service)
		if err != nil {
			serviceInvokeCallLog.Error(err.Error())
			return shim.Error(err.Error())
		}
		return shim.Success(serviceAsJSON)
	}
}

// =====================================================================================================================
// Query Service - wrapper of GetService called from the chaincode invoke
// =====================================================================================================================
func QueryService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ServiceId"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		serviceInvokeCallLog.Error(argumentSizeError.Error())
		return shim.Error(argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		serviceInvokeCallLog.Error(sanitizeError.Error())
		return shim.Error(sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== get the service ====
	service, err := a.GetService(stub, serviceId)
	if err != nil {
		serviceInvokeCallLog.Info("Failed to find service by id " + serviceId)
		serviceInvokeCallLog.Error(err.Error())
		return shim.Error(err.Error())
	} else {
		serviceInvokeCallLog.Info("Service ID: " + service.ServiceId + ", Service: " + service.Name + ", with Description: " + service.Description + " found")
		// ==== Marshal the byService query result ====
		serviceAsJSON, err := json.Marshal(service)
		if err != nil {
			serviceInvokeCallLog.Error(err.Error())
			return shim.Error(err.Error())
		}
		return shim.Success(serviceAsJSON)
	}
}

// ========================================================================================================================
// Query by  - wrapper of GetByServiceName called from chiancode's Invoke
// ========================================================================================================================
func QueryByServiceName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "ServiceName"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		serviceInvokeCallLog.Error(argumentSizeError.Error())
		return shim.Error(argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := arglib.SanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		serviceInvokeCallLog.Error(sanitizeError.Error())
		return shim.Error(sanitizeError.Error())
	}

	serviceName := args[0]

	// ==== Run the byServiceName query ====
	byServiceNameQueryIterator, err := a.GetByServiceName(serviceName, stub)
	if err != nil {
		serviceInvokeCallLog.Info("Failed to get service by name: " + serviceName)
		serviceInvokeCallLog.Error(err.Error())
		return shim.Error(err.Error())
	}
	if byServiceNameQueryIterator != nil {
		serviceInvokeCallLog.Info(&byServiceNameQueryIterator)
	}

	// ==== Get the Services for the byServiceName query result ====
	servicesSlice, err := a.GetServiceSliceFromRangeQuery(byServiceNameQueryIterator, stub)
	if err != nil {
		serviceInvokeCallLog.Error(err.Error())
		return shim.Error(err.Error())
	}

	// ==== Marshal the byServiceName query result ====
	servicesByNameAsBytes, err := json.Marshal(servicesSlice)
	if err != nil {
		serviceInvokeCallLog.Error(err.Error())
		return shim.Error(err.Error())
	}
	serviceInvokeCallLog.Info(servicesByNameAsBytes)

	stringOut := string(servicesByNameAsBytes)
	if stringOut == "null" {
		serviceInvokeCallLog.Error("Doesn't exist a service with the name: " + serviceName)
		return shim.Error("Doesn't exist a service with the name: " + serviceName)
	}

	// ==== Return success with servicesByNameAsBytes as payload ====
	return shim.Success(servicesByNameAsBytes)
}

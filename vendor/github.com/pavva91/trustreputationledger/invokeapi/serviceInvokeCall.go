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
	a "github.com/pavva91/trustreputationledger/assets"
)

// =====================================================================================================================
// Init Service - wrapper of CreateService called from the chaincode invoke
// =====================================================================================================================
func CreateService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1                 2
	// "ServiceId", "serviceName", "serviceDescription"
	argumentSizeError := arglib.ArgumentSizeVerification(args, 3)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}
	fmt.Println("- start init service")

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
		fmt.Println("This service already exists: " + serviceName)
		return shim.Error("This service already exists: " + serviceName)
	}

	service, err := a.CreateService(serviceId, serviceName, serviceDescription, stub)
	if err != nil {
		return shim.Error("Failed to create the service: " + err.Error())
	}

	// ==== Indexing of service by Name (to do query by Name, if you want) ====
	// index create
	nameIndexKey, nameIndexError := a.CreateNameIndex(service, stub)
	if nameIndexError != nil {
		return shim.Error(nameIndexError.Error())
	}
	fmt.Println(nameIndexKey)

	// index save
	saveIndexError := a.SaveIndex(nameIndexKey, stub)
	if saveIndexError != nil {
		return shim.Error(saveIndexError.Error())
	}

	// ==== Service saved and indexed. Return success ====
	fmt.Println("Servizio: " + service.Name + " creato - end init service")
	return shim.Success(nil)
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
		fmt.Println("Failed to find service by id " + serviceId)
		return shim.Error(getError.Error())
	}

	// ==== modify the service ====
	modifyError := a.ModifyServiceName(service, newServiceName, stub)
	if modifyError != nil {
		fmt.Println("Failed to modify the service name: " + newServiceName)
		return shim.Error(modifyError.Error())
	}

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
		fmt.Println("Failed to find service by id " + serviceId)
		return shim.Error(getError.Error())
	}

	// ==== modify the service ====
	modifyError := a.ModifyServiceDescription(service, newServiceDescription, stub)
	if modifyError != nil {
		fmt.Println("Failed to modify the service description: " + newServiceDescription)
		return shim.Error(modifyError.Error())
	}

	return shim.Success(nil)
}

// =====================================================================================================================
// Query Service - wrapper of GetServiceNotFoundError called from the chaincode invoke
// =====================================================================================================================
func QueryService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
		fmt.Println("Failed to find service by id " + serviceId)
		return shim.Error(err.Error())
	} else {
		fmt.Println("Service ID: " + service.ServiceId + ", Service: " + service.Name + ", with Description: " + service.Description + " found")
		// ==== Marshal the byService query result ====
		serviceAsJSON, err := json.Marshal(service)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(serviceAsJSON)
	}
}

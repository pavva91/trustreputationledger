package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"encoding/json"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// Init Service - wrapper of createService called from the chaincode invoke
// ============================================================================================================================
func initService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0               1                 2
	// "ServiceId", "serviceName", "serviceDescription"
	argumentSizeError := argumentSizeVerification(args, 3)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}
	fmt.Println("- start init service")

	// ==== Input sanitation ====
	sanitizeError := sanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]
	serviceName := args[1]
	serviceDescription := args[2]

	// ==== Check if service already exists ====
	serviceAsBytes, err := getServiceAsBytes(stub, serviceId)
	if err != nil {
		return shim.Error("Failed to get service: " + err.Error())
	} else if serviceAsBytes != nil {
		fmt.Println("This service already exists: " + serviceName)
		return shim.Error("This service already exists: " + serviceName)
	}

	service, err := createService(serviceId, serviceName, serviceDescription, stub)
	if err != nil {
		return shim.Error("Failed to create the service: " + err.Error())
	}

	// ==== Indexing of service by Name (to do query by Name, if you want) ====
	// index create
	nameIndexKey, nameIndexError := createNameIndex(service, stub)
	if nameIndexError != nil {
		return shim.Error(nameIndexError.Error())
	}

	// index save
	saveIndexError := saveIndex(nameIndexKey, stub)
	if saveIndexError != nil {
		return shim.Error(saveIndexError.Error())
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
	argumentSizeError := argumentSizeVerification(args, 1)
	if argumentSizeError != nil {
		return shim.Error("Argument Size Error: " + argumentSizeError.Error())
	}

	// ==== Input sanitation ====
	sanitizeError := sanitizeArguments(args)
	if sanitizeError != nil {
		fmt.Print(sanitizeError)
		return shim.Error("Sanitize error: " + sanitizeError.Error())
	}

	serviceId := args[0]

	// ==== get the service ====
	service, err := getService(stub, serviceId)
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


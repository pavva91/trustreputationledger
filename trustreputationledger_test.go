/*
 * Copyright 2018 IBM All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an 'AS IS' BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	lib "github.com/pavva91/arglib"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	a "github.com/pavva91/assets"
)

var testLog = shim.NewLogger("trustreputationledger_test")

const (
	ExistingServiceId          = "idservice1"
	ExistingServiceName        = "service1"
	ExistingServiceDescription = "service Description 1"
	ExistingAgentId            = "idagent1"
	ExistingAgentName      = "agent1"
	ExistingAgentAddress   = "address1"
	NewServiceId           = "idservice6"
	NewServiceName         = "service6"
	NewServiceDescription  = "service Description 6"
	NewServiceSameNameId = "idservice1000"
	ServiceComposition = "asdf,fdas"
	NullServiceComposition = ""
	NewAgentId             = "idagent6"
	NewAgentName           = "agent6"
	NewAgentAddress        = "address6"
	ServiceAgentServiceId  = ExistingServiceId
	ServiceAgentAgentId    = ExistingAgentId
	ServiceAgentCost       = "8"
	ServiceAgentTime       = "6"
	ReputationValue        = "6"
	ExecuterAgentId = "idagent99"
	DemanderAgentId = "idagent98"
	WritingExecuterAgentId = ExecuterAgentId
	WritingDemanderAgentId = DemanderAgentId
	ExecutedServiceId = "idservice99"
	ExecutedServiceTxId = "execServiceTxId"
	ExecutedServiceTimestamp = "execServiceTimestamp"
	ActivityValue = "10"


	EXPORTER = "LumberInc"
	EXPBANK = "LumberBank"
	EXPBALANCE = 100000
	IMPORTER = "WoodenToys"
	IMPBANK = "ToyBank"
	IMPBALANCE = 200000
	CARRIER = "UniversalFrieght"
	REGAUTH = "ForestryDepartment"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		testLog.Info("Init failed", string(res.Message))
		t.FailNow()
	}else{
		testLog.Info("Init OK", string(res.Message))

	}
}

func checkNoState(t *testing.T, stub *shim.MockStub, name string) {
	bytes := stub.State[name]
	if bytes != nil {
		testLog.Info("State", name, "should be absent; found value")
		t.FailNow()
	}else {
		testLog.Info("State", name, "is absent as it should be")
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		testLog.Info("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		testLog.Info("State value", name, "was", string(bytes), "and not", value, "as expected")
		t.FailNow()
	}else{
		testLog.Info("State value", name, "is", string(bytes), "as expected")
	}
}

func checkBadQuery(t *testing.T, stub *shim.MockStub, function string, name string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status == shim.OK {
		testLog.Info("Query", name, "unexpectedly succeeded")
		t.FailNow()
	}else {
		testLog.Info("Query", name, "failed as espected, with message: ",res.Message)

	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, function string, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status != shim.OK {
		testLog.Info("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		testLog.Info("Query", name, "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		testLog.Info("Query value", name, "was", payload, "and not", value, "as expected")
		t.FailNow()
	}else{
		testLog.Info("Query value", name, "is", payload, "as expected")
	}
}

func checkQueryArgs(t *testing.T, stub *shim.MockStub, args [][]byte, value string) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		testLog.Info("Query", string(args[1]), "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		testLog.Info("Query", string(args[1]), "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		testLog.Info("Query value", string(args[1]), "was", payload, "and not", value, "as expected")
		t.FailNow()
	}else {
		testLog.Info("Query value", string(args[1]), "is", payload, "as expected")

	}
}

func checkBadInvoke(t *testing.T, stub *shim.MockStub, functionAndArgs []string) {
	functionAndArgsAsBytes := lib.ParseStringSliceToByteSlice(functionAndArgs)
	res := stub.MockInvoke("1", functionAndArgsAsBytes)
	if res.Status == shim.OK {
		testLog.Info("Invoke", functionAndArgs, "unexpectedly succeeded")
		t.FailNow()
	}else {
		testLog.Info("Invoke", functionAndArgs, "failed as espected, with message: "+ res.Message)
	}
}

// func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
// 	res := stub.MockInvoke("1", args)
// 	if res.Status != shim.OK {
// 		testLog.Info("Invoke", args, "failed", string(res.Message))
// 		t.FailNow()
// 	}else {
// 		testLog.Info("Invoke", args, "successful", string(res.Message))
// 	}
// }
func checkInvoke(t *testing.T, stub *shim.MockStub, functionAndArgs []string) {
	functionAndArgsAsBytes := lib.ParseStringSliceToByteSlice(functionAndArgs)
	res := stub.MockInvoke("1", functionAndArgsAsBytes)
	if res.Status != shim.OK {
		testLog.Info("Invoke", functionAndArgs, "failed", string(res.Message))
		t.FailNow()
	}else {
		testLog.Info("Invoke", functionAndArgs, "successful", string(res.Message))
	}
}

func getInitArguments() [][]byte {
	return [][]byte{}
}

// =====================================================================================================================
// TestTrustReputationInit - Test the 'Init' function
// =====================================================================================================================
func TestTrustReputationInit(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	stub := shim.NewMockStub("Test Init", simpleChaincode)

	// Init
	checkInit(t, stub, getInitArguments())



}

// TEST CREATE:

// =====================================================================================================================
// TestServiceCreation - Test the 'CreateService' function
// =====================================================================================================================
func TestServiceCreation(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Service Creation", simpleChaincode)

	var functionAndArgs []string
	functionName:= CreateService

	// Invoke 'CreateLeafService'
	serviceId := NewServiceId
	serviceName := NewServiceName
	serviceDescription := NewServiceDescription

	serviceCompositionAsString := "asd,fda"

	serviceComposition := lib.ParseStringToStringSlice(serviceCompositionAsString)
	args := []string{serviceId,serviceName,serviceDescription, serviceCompositionAsString}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	service := &a.Service{ServiceId: serviceId, Name: serviceName, Description: serviceDescription, ServiceComposition:serviceComposition}
	serviceAsBytes, _ := json.Marshal(service)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceId})
	checkState(t, mockStub, serviceId, string(serviceAsBytes))
	testLog.Info(serviceComposition)
	serviceCompositionJsonRappresentation := "["
	for i := 0; i<len(service.ServiceComposition) ;i++  {
		if i == 0 {
			serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation + "\""+service.ServiceComposition[i]+ "\""
		}else {
			serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation + ",\""+service.ServiceComposition[i]+ "\""
		}
	}
	serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation +"]"


	expectedResp := "{\"ServiceId\":\""+ serviceId + "\",\"Name\":\""+ serviceName + "\",\"Description\":\""+ serviceDescription + "\",\"ServiceComposition\":"+ serviceCompositionJsonRappresentation + "}"
	checkQuery(t, mockStub, "GetServiceNotFoundError", serviceId, expectedResp)
}

// =====================================================================================================================
// TestServiceCreationWithEmptyServiceComposition - Test the 'CreateService' function
// =====================================================================================================================
func TestServiceCreationWithEmptyServiceComposition(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Service Creation", simpleChaincode)

	var functionAndArgs []string
	functionName:= CreateService

	// Invoke 'CreateLeafService'
	serviceId := NewServiceId
	serviceName := NewServiceName
	serviceDescription := NewServiceDescription

	serviceCompositionAsString := NullServiceComposition

	serviceComposition := lib.ParseStringToStringSlice(serviceCompositionAsString)
	args := []string{serviceId,serviceName,serviceDescription, serviceCompositionAsString}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	service := &a.Service{ServiceId: serviceId, Name: serviceName, Description: serviceDescription, ServiceComposition:serviceComposition}
	serviceAsBytes, _ := json.Marshal(service)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceId})
	checkState(t, mockStub, serviceId, string(serviceAsBytes))

	expectedResp := "{\"ServiceId\":\""+ serviceId + "\",\"Name\":\""+ serviceName + "\",\"Description\":\""+ serviceDescription + "\",\"ServiceComposition\":null}"
	checkQuery(t, mockStub, "GetServiceNotFoundError", serviceId, expectedResp)
}

// =====================================================================================================================
// TestServiceCreationWithEmptyServiceComposition - Test the 'CreateService' function
// =====================================================================================================================
func TestServiceCreationWithMissingServiceComposition(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Service Creation", simpleChaincode)

	var functionAndArgs []string
	functionName:= CreateService

	// Invoke 'CreateLeafService'
	serviceId := NewServiceId
	serviceName := NewServiceName
	serviceDescription := NewServiceDescription


	args := []string{serviceId,serviceName,serviceDescription}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	service := &a.Service{ServiceId: serviceId, Name: serviceName, Description: serviceDescription}
	serviceAsBytes, _ := json.Marshal(service)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceId})
	checkState(t, mockStub, serviceId, string(serviceAsBytes))


	testLog.Info(len(service.ServiceComposition))

	expectedResp := "{\"ServiceId\":\""+ serviceId + "\",\"Name\":\""+ serviceName + "\",\"Description\":\""+ serviceDescription + "\",\"ServiceComposition\":null}"
	checkQuery(t, mockStub, "GetServiceNotFoundError", serviceId, expectedResp)
}
// =====================================================================================================================
// TestLeafServiceCreation - Test the 'CreateLeafService' function
// =====================================================================================================================
func TestLeafServiceCreation(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Service Creation", simpleChaincode)

	var functionAndArgs []string
	functionName:= CreateLeafService

	// Invoke 'CreateLeafService'
	serviceId := NewServiceId
	serviceName := NewServiceName
	serviceDescription := NewServiceDescription



	args := []string{serviceId,serviceName,serviceDescription}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	service := &a.Service{ServiceId: serviceId, Name: serviceName, Description: serviceDescription}
	serviceAsBytes, _ := json.Marshal(service)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceId})
	checkState(t, mockStub, serviceId, string(serviceAsBytes))

	expectedResp := "{\"ServiceId\":\""+ serviceId + "\",\"Name\":\""+ serviceName + "\",\"Description\":\""+ serviceDescription + "\",\"ServiceComposition\":null}"
	checkQuery(t, mockStub, "GetServiceNotFoundError", serviceId, expectedResp)
}
// =====================================================================================================================
// TestCompositeServiceCreation - Test the 'CreateLeafService' function
// =====================================================================================================================
func TestCompositeServiceCreation(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Service Creation", simpleChaincode)

	var functionAndArgs []string
	functionName:= CreateCompositeService

	// Invoke 'CreateLeafService'
	serviceId := NewServiceId
	serviceName := NewServiceName
	serviceDescription := NewServiceDescription

	serviceCompositionAsString := ServiceComposition
	serviceComposition := lib.ParseStringToStringSlice(serviceCompositionAsString)

	args := []string{serviceId,serviceName,serviceDescription,serviceCompositionAsString}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	service := &a.Service{ServiceId: serviceId, Name: serviceName, Description: serviceDescription, ServiceComposition:serviceComposition}
	serviceAsBytes, _ := json.Marshal(service)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceId})
	checkState(t, mockStub, serviceId, string(serviceAsBytes))


	serviceCompositionJsonRappresentation := "["
	for i := 0; i<len(service.ServiceComposition) ;i++  {
		if i == 0 {
			serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation + "\""+service.ServiceComposition[i]+ "\""
		}else {
			serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation + ",\""+service.ServiceComposition[i]+ "\""
		}
	}
	serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation +"]"


	expectedResp := "{\"ServiceId\":\""+ serviceId + "\",\"Name\":\""+ serviceName + "\",\"Description\":\""+ serviceDescription + "\",\"ServiceComposition\":"+serviceCompositionJsonRappresentation+"}"
	checkQuery(t, mockStub, "GetServiceNotFoundError", serviceId, expectedResp)
}
// =====================================================================================================================
// TestCompositeServiceCreationWithNullValue - Test the 'CreateLeafService' function
// =====================================================================================================================
func TestCompositeServiceCreationWithNullValue(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Service Creation", simpleChaincode)

	var functionAndArgs []string
	functionName:= CreateCompositeService

	// Invoke 'CreateLeafService'
	serviceId := NewServiceId
	serviceName := NewServiceName
	serviceDescription := NewServiceDescription

	serviceCompositionAsString := NullServiceComposition
	// serviceComposition := lib.ParseStringToStringSlice(serviceCompositionAsString)

	args := []string{serviceId,serviceName,serviceDescription,serviceCompositionAsString}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkBadInvoke(t, mockStub, functionAndArgs)
	checkBadQuery(t, mockStub, "GetServiceNotFoundError", serviceId)
}
// =====================================================================================================================
// TestExistingLeafServiceCreation - Test the 'CreateLeafService' function when trying to insert an already existing record
// =====================================================================================================================
func TestExistingLeafServiceCreation(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Already Existing Service Creation", simpleChaincode)

	// Init
	checkInit(t, mockStub, getInitArguments())

	var functionAndArgs []string
	functionName:= CreateLeafService

	// Invoke 'CreateService'
	existingServiceId := ExistingServiceId
	serviceName := ExistingServiceName
	serviceDescription := ExistingServiceDescription

	args := []string{existingServiceId,serviceName,serviceDescription}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkBadInvoke(t, mockStub, functionAndArgs)


	service := &a.Service{ServiceId: existingServiceId, Name: serviceName, Description: serviceDescription}
	serviceBytes, _ := json.Marshal(service)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{existingServiceId})
	checkState(t, mockStub, existingServiceId, string(serviceBytes))

	expectedResp := "{\"ServiceId\":\""+ existingServiceId + "\",\"Name\":\""+ serviceName + "\",\"Description\":\""+ serviceDescription + "\",\"ServiceComposition\":null}"
	checkQuery(t, mockStub, "GetServiceNotFoundError", existingServiceId, expectedResp)
}

// =====================================================================================================================
// TestAgentCreation - Test the 'CreateAgent' function
// =====================================================================================================================
func TestAgentCreation(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Agent Creation", simpleChaincode)

	// Init
	// checkInit(t, mockStub, getInitArguments())

	var functionAndArgs []string
	functionName:= CreateAgent

	// Invoke 'CreateAgent'
	agentId := NewAgentId
	agentName := NewAgentName
	agentAddress := NewAgentAddress

	args := []string{agentId,agentName,agentAddress}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	agent := &a.Agent{agentId, agentName, agentAddress}
	agentAsBytes, _ := json.Marshal(agent)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{agentId})
	checkState(t, mockStub, agentId, string(agentAsBytes))

	expectedResp := "{\"AgentId\":\""+ agentId + "\",\"Name\":\""+ agentName + "\",\"Address\":\""+ agentAddress + "\"}"
	checkQuery(t, mockStub, "GetAgentNotFoundError", agentId, expectedResp)


}
// =====================================================================================================================
// TestExistingAgentCreation - Test the 'CreateAgent' function when trying to insert an already existing record
// =====================================================================================================================
func TestExistingAgentCreation(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Already Existing Agent Creation", simpleChaincode)

	// Init
	checkInit(t, mockStub, getInitArguments())

	var functionAndArgs []string
	functionName:= CreateAgent

	// Invoke 'CreateAgent'
	agentId := ExistingAgentId
	agentName := ExistingAgentName
	agentAddress := ExistingAgentAddress

	args := []string{agentId,agentName,agentAddress}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkBadInvoke(t, mockStub, functionAndArgs)

	agent := &a.Agent{agentId, agentName, agentAddress}
	agentAsBytes, _ := json.Marshal(agent)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{agentId})
	checkState(t, mockStub, agentId, string(agentAsBytes))

	expectedResp := "{\"AgentId\":\""+ agentId + "\",\"Name\":\""+ agentName + "\",\"Address\":\""+ agentAddress + "\"}"
	checkQuery(t, mockStub, "GetAgentNotFoundError", agentId, expectedResp)
}
// =====================================================================================================================
// TestServiceAgentRelationCreation - Test the 'CreateServiceAgentRelation' function
// =====================================================================================================================
func TestServiceAgentRelationCreation(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test ServiceAgentRelation Creation", simpleChaincode)

	// Init
	checkInit(t, mockStub, getInitArguments())

	var functionAndArgs []string
	functionName:= CreateServiceAgentRelation

	// Invoke 'CreateServiceAgentRelation'
	serviceId := ServiceAgentServiceId
	agentId := ServiceAgentAgentId
	cost := ServiceAgentCost
	time := ServiceAgentTime

	args := []string{serviceId,agentId,cost,time}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	relationId := serviceId + agentId

	serviceRelationAgent := &a.ServiceRelationAgent{relationId, serviceId, agentId, cost, time}
	serviceRealationAgentAsBytes, _ := json.Marshal(serviceRelationAgent)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{agentId})
	checkState(t, mockStub, relationId, string(serviceRealationAgentAsBytes))


	expectedResp := "{\"RelationId\":\""+ relationId +"\",\"ServiceId\":\""+ serviceId +"\",\"AgentId\":\""+ agentId + "\",\"Cost\":\""+ cost + "\",\"Time\":\""+ time + "\"}"
	checkQuery(t, mockStub, GetServiceRelationAgent, relationId, expectedResp)
}
// =====================================================================================================================
// TestServiceAndServiceAgentRelationWithStandardValueCreationNewService - Test the 'CreateServiceAndServiceAgentRelationWithStandardValue' function adding a new service
// =====================================================================================================================
func TestServiceAndServiceAgentRelationWithStandardValueCreationNewService(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test ServiceAndServiceAgentRelationWithStandardValue Creation of a New Service", simpleChaincode)
	// Init
	checkInit(t, mockStub, getInitArguments())

	var functionAndArgs []string
	functionName := CreateServiceAndServiceAgentRelationWithStandardValue

	// "ServiceId", "ServiceName", "ServiceDescription", "AgentId", "Cost", "Time"
	// Invoke 'CreateServiceAndServiceAgentRelationWithStandardValue'
	serviceId := NewServiceId
	serviceName := NewServiceName
	serviceDescription := NewServiceDescription
	agentId := ExistingAgentId
	cost := ServiceAgentCost
	time := ServiceAgentTime

	args := []string{serviceId,serviceName,serviceDescription,agentId,cost,time}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	relationId := serviceId + agentId

	serviceRelationAgent := &a.ServiceRelationAgent{relationId, serviceId, agentId, cost, time}
	serviceRealationAgentAsBytes, _ := json.Marshal(serviceRelationAgent)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceName})
	checkState(t, mockStub, relationId, string(serviceRealationAgentAsBytes))

	expectedResp := "{\"RelationId\":\""+ relationId +"\",\"ServiceId\":\""+ serviceId +"\",\"AgentId\":\""+ agentId + "\",\"Cost\":\""+ cost + "\",\"Time\":\""+ time + "\"}"
	checkQuery(t, mockStub, GetServiceRelationAgent, relationId, expectedResp)
}
// =====================================================================================================================
// TestServiceAndServiceAgentRelationWithStandardValueExistingService - Test the 'CreateServiceAndServiceAgentRelationWithStandardValue' function using an existing service
// =====================================================================================================================
func TestServiceAndServiceAgentRelationWithStandardValueExistingService(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test ServiceAndServiceAgentRelationWithStandardValue of an Existing Service", simpleChaincode)
	// Init
	checkInit(t, mockStub, getInitArguments())

	var functionAndArgs []string
	functionName := CreateServiceAndServiceAgentRelationWithStandardValue

	// "ServiceId", "ServiceName", "ServiceDescription", "AgentId", "Cost", "Time"
	// Invoke 'CreateServiceAndServiceAgentRelationWithStandardValue'
	serviceId := ExistingServiceId
	serviceName := ExistingServiceName
	serviceDescription := ExistingServiceDescription
	agentId := ExistingAgentId
	cost := ServiceAgentCost
	time := ServiceAgentTime

	args := []string{serviceId,serviceName,serviceDescription,agentId,cost,time}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	relationId := serviceId + agentId

	serviceRelationAgent := &a.ServiceRelationAgent{relationId, serviceId, agentId, cost, time}
	serviceRealationAgentAsBytes, _ := json.Marshal(serviceRelationAgent)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceName})
	checkState(t, mockStub, relationId, string(serviceRealationAgentAsBytes))

	expectedResp := "{\"RelationId\":\""+ relationId +"\",\"ServiceId\":\""+ serviceId +"\",\"AgentId\":\""+ agentId + "\",\"Cost\":\""+ cost + "\",\"Time\":\""+ time + "\"}"
	checkQuery(t, mockStub, GetServiceRelationAgent, relationId, expectedResp)
}
// =====================================================================================================================
// TestServiceAndServiceAgentRelationWithStandardValueCreationNewService - Test the 'CreateServiceAndServiceAgentRelation' function adding a new service with passed reputation
// =====================================================================================================================
func TestServiceAndServiceAgentRelationCreationNewService(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test ServiceAndServiceAgentRelation Creation of a New Service", simpleChaincode)

	// Init
	checkInit(t, mockStub, getInitArguments())

	var functionAndArgs []string
	functionName:= CreateServiceAndServiceAgentRelation

	// "ServiceId", "ServiceName", "ServiceDescription", "AgentId", "Cost", "Time", "InitReputationValue"
	// Invoke 'CreateServiceAndServiceAgentRelation'
	serviceId := NewServiceId
	serviceName := NewServiceName
	serviceDescription := NewServiceDescription
	agentId := ExistingAgentId
	cost := ServiceAgentCost
	time := ServiceAgentTime
	initReputationValue := ReputationValue

	args := []string{serviceId,serviceName,serviceDescription,agentId,cost,time,initReputationValue}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	relationId := serviceId + agentId

	serviceRelationAgent := &a.ServiceRelationAgent{relationId, serviceId, agentId, cost, time}
	serviceRealationAgentAsBytes, _ := json.Marshal(serviceRelationAgent)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceName})
	checkState(t, mockStub, relationId, string(serviceRealationAgentAsBytes))

	expectedResp := "{\"RelationId\":\""+ relationId +"\",\"ServiceId\":\""+ serviceId +"\",\"AgentId\":\""+ agentId + "\",\"Cost\":\""+ cost + "\",\"Time\":\""+ time + "\"}"
	checkQuery(t, mockStub, "GetServiceRelationAgent", relationId, expectedResp)

	agentRole := a.Executer

	reputationId := agentId + serviceId + agentRole

	reputation := &a.Reputation{reputationId, agentId, serviceId, agentRole, initReputationValue}
	reputationAsBytes, _ := json.Marshal(reputation)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceName})
	checkState(t, mockStub, reputationId, string(reputationAsBytes))
	expectedResp2 := "{\"ReputationId\":\""+ reputationId +"\",\"AgentId\":\""+ agentId +"\",\"ServiceId\":\""+ serviceId +"\",\"AgentRole\":\""+ agentRole +"\",\"Value\":\""+ initReputationValue +"\"}"

	checkQuery(t, mockStub, GetReputationNotFoundError, reputationId, expectedResp2)
}
// =====================================================================================================================
// TestServiceAndServiceAgentRelationExistingService - Test the 'CreateServiceAndServiceAgentRelation' function adding an existing service
// =====================================================================================================================
func TestServiceAndServiceAgentRelationExistingService(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test ServiceAndServiceAgentRelation Creation of a New Service", simpleChaincode)

	// Init
	checkInit(t, mockStub, getInitArguments())

	var functionAndArgs []string
	functionName:= CreateServiceAndServiceAgentRelation

	// "ServiceId", "ServiceName", "ServiceDescription", "AgentId", "Cost", "Time", "InitReputationValue"
	// Invoke 'CreateServiceAndServiceAgentRelation'
	serviceId := ExistingServiceId
	serviceName := ExistingServiceName
	serviceDescription := ExistingServiceDescription
	agentId := ExistingAgentId
	cost := ServiceAgentCost
	time := ServiceAgentTime
	initReputationValue := ReputationValue

	args := []string{serviceId,serviceName,serviceDescription,agentId,cost,time,initReputationValue}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	relationId := serviceId + agentId

	serviceRelationAgent := &a.ServiceRelationAgent{relationId, serviceId, agentId, cost, time}
	serviceRealationAgentAsBytes, _ := json.Marshal(serviceRelationAgent)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceName})
	checkState(t, mockStub, relationId, string(serviceRealationAgentAsBytes))

	expectedResp := "{\"RelationId\":\""+ relationId +"\",\"ServiceId\":\""+ serviceId +"\",\"AgentId\":\""+ agentId + "\",\"Cost\":\""+ cost + "\",\"Time\":\""+ time + "\"}"
	checkQuery(t, mockStub, "GetServiceRelationAgent", relationId, expectedResp)

	agentRole := a.Executer

	reputationId := agentId + serviceId + agentRole

	reputation := &a.Reputation{reputationId, agentId, serviceId, agentRole, initReputationValue}
	reputationAsBytes, _ := json.Marshal(reputation)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceName})
	checkState(t, mockStub, reputationId, string(reputationAsBytes))
	expectedResp2 := "{\"ReputationId\":\""+ reputationId +"\",\"AgentId\":\""+ agentId +"\",\"ServiceId\":\""+ serviceId +"\",\"AgentRole\":\""+ agentRole +"\",\"Value\":\""+ initReputationValue +"\"}"

	checkQuery(t, mockStub, GetReputationNotFoundError, reputationId, expectedResp2)
}
// =====================================================================================================================
// TestExecuterActivityCreation - Test the 'CreateActivity' function called from an Reputation.EXECUTER
// =====================================================================================================================
func TestExecuterActivityCreation(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Executer Activity Creation", simpleChaincode)

	// Init
	checkInit(t, mockStub, getInitArguments())

	//   0               1                   2                     3                   4                        5         6
	// "WriterAgentId", "DemanderAgentId", "ExecuterAgentId", "ExecutedServiceId", "ExecutedServiceTxId", "ExecutedServiceTimestamp", "Value"
	var functionAndArgs []string
	functionName:= CreateActivity

	// Invoke 'CreateServiceAgentRelation'
	writerAgentId := WritingExecuterAgentId
	demanderAgentId := DemanderAgentId
	executerAgentId := ExecuterAgentId
	executedServiceId := ExecutedServiceId
	executedServiceTxId := ExecutedServiceTxId
	executedServiceTimestamp := ExecutedServiceTimestamp
	activityValue := ActivityValue

	args := []string{writerAgentId, demanderAgentId, executerAgentId, executedServiceId,executedServiceTxId,executedServiceTimestamp, activityValue}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	evaluationId := writerAgentId + demanderAgentId + executerAgentId + executedServiceTxId

	activity := &a.Activity{evaluationId, writerAgentId, demanderAgentId, executerAgentId, executedServiceId,executedServiceTxId,executedServiceTimestamp, activityValue}
	activityAsBytes, _ := json.Marshal(activity)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{demanderAgentId})
	checkState(t, mockStub, evaluationId, string(activityAsBytes))

	expectedResp := "{\"EvaluationId\":\""+ evaluationId +"\",\"WriterAgentId\":\""+ writerAgentId +"\",\"DemanderAgentId\":\""+ demanderAgentId + "\",\"ExecuterAgentId\":\""+ executerAgentId + "\",\"ExecutedServiceId\":\""+ executedServiceId + "\",\"ExecutedServiceTxid\":\""+ executedServiceTxId + "\",\"ExecutedServiceTimestamp\":\""+ executedServiceTimestamp + "\",\"Value\":\""+ activityValue + "\"}"
	checkQuery(t, mockStub, GetActivity, evaluationId, expectedResp)
}
// =====================================================================================================================
// TestDemanderActivityCreation - Test the 'CreateActivity' function called from an Reputation.DEMANDER
// =====================================================================================================================
func TestDemanderActivityCreation(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Demander Activity Creation", simpleChaincode)

	// Init
	checkInit(t, mockStub, getInitArguments())

	//   0               1                   2                     3                   4                        5         6
	// "WriterAgentId", "DemanderAgentId", "ExecuterAgentId", "ExecutedServiceId", "ExecutedServiceTxId", "ExecutedServiceTimestamp", "Value"
	var functionAndArgs []string
	functionName:= CreateActivity

	// Invoke 'CreateServiceAgentRelation'
	writerAgentId := WritingDemanderAgentId
	demanderAgentId := DemanderAgentId
	executerAgentId := ExecuterAgentId
	executedServiceId := ExecutedServiceId
	executedServiceTxId := ExecutedServiceTxId
	executedServiceTimestamp := ExecutedServiceTimestamp
	activityValue := ActivityValue

	args := []string{writerAgentId, demanderAgentId, executerAgentId, executedServiceId,executedServiceTxId,executedServiceTimestamp, activityValue}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	evaluationId := writerAgentId + demanderAgentId + executerAgentId + executedServiceTxId

	activity := &a.Activity{evaluationId, writerAgentId, demanderAgentId, executerAgentId, executedServiceId,executedServiceTxId,executedServiceTimestamp, activityValue}
	activityAsBytes, _ := json.Marshal(activity)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{demanderAgentId})
	checkState(t, mockStub, evaluationId, string(activityAsBytes))

	expectedResp := "{\"EvaluationId\":\""+ evaluationId +"\",\"WriterAgentId\":\""+ writerAgentId +"\",\"DemanderAgentId\":\""+ demanderAgentId + "\",\"ExecuterAgentId\":\""+ executerAgentId + "\",\"ExecutedServiceId\":\""+ executedServiceId + "\",\"ExecutedServiceTxid\":\""+ executedServiceTxId + "\",\"ExecutedServiceTimestamp\":\""+ executedServiceTimestamp + "\",\"Value\":\""+ activityValue + "\"}"
	checkQuery(t, mockStub, GetActivity, evaluationId, expectedResp)
}

// =====================================================================================================================
// TestQueryByServiceName - Test the 'GetServicesByName' function
// =====================================================================================================================
func TestQueryByServiceName(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Get Services by Service Name", simpleChaincode)

	// CREATION OF SERVICE 1:
	var functionAndArgsCreateService1 []string
	createServiceFunctionName := CreateService
	newServiceId1 := NewServiceId
	sameServiceName := NewServiceName
	newServiceDescription1 := NewServiceDescription
	serviceCompositionAsString1 := "asd,fda"
	serviceComposition := lib.ParseStringToStringSlice(serviceCompositionAsString1)
	args1 := []string{newServiceId1, sameServiceName, newServiceDescription1, serviceCompositionAsString1}
	functionAndArgsCreateService1 = append(functionAndArgsCreateService1, createServiceFunctionName)
	functionAndArgsCreateService1 = append(functionAndArgsCreateService1,args1...)

	checkInvoke(t, mockStub, functionAndArgsCreateService1)

	serviceCompositionJsonRappresentation := "["
	for i := 0; i<len(serviceComposition) ;i++  {
		if i == 0 {
			serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation + "\""+serviceComposition[i]+ "\""
		}else {
			serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation + ",\""+serviceComposition[i]+ "\""
		}
	}
	serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation +"]"

	// CREATION OF SERVICE 2 (WITH THE SAME NAME)
	var functionAndArgsCreateService2 []string
	newServiceId2 := NewServiceSameNameId
	newServiceDescription2 := NewServiceDescription
	serviceCompositionAsString2 := "blu,les"
	serviceComposition2 := lib.ParseStringToStringSlice(serviceCompositionAsString2)
	args2 := []string{newServiceId2, sameServiceName, newServiceDescription2, serviceCompositionAsString2}
	functionAndArgsCreateService2 = append(functionAndArgsCreateService2, createServiceFunctionName)
	functionAndArgsCreateService2 = append(functionAndArgsCreateService2, args2...)

	checkInvoke(t, mockStub, functionAndArgsCreateService2)

	serviceCompositionJsonRappresentation2 := "["
	for i := 0; i<len(serviceComposition) ;i++  {
		if i == 0 {
			serviceCompositionJsonRappresentation2 = serviceCompositionJsonRappresentation2 + "\""+ serviceComposition2[i]+ "\""
		}else {
			serviceCompositionJsonRappresentation2 = serviceCompositionJsonRappresentation2 + ",\""+ serviceComposition2[i]+ "\""
		}
	}
	serviceCompositionJsonRappresentation2 = serviceCompositionJsonRappresentation2 +"]"

	// VERIFY THE QUERY GetServicesByName with the newly created services
	//   0
	// "serviceName"
	var functionAndArgs []string
	functionName:= GetServicesByName

	serviceName := sameServiceName

	args := []string{serviceName}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args...)

	expectedResp := "[{\"ServiceId\":\""+ newServiceId2 + "\",\"Name\":\""+ sameServiceName + "\",\"Description\":\""+ newServiceDescription2 + "\",\"ServiceComposition\":"+ serviceCompositionJsonRappresentation2 + "},{\"ServiceId\":\""+ newServiceId1 + "\",\"Name\":\""+ sameServiceName + "\",\"Description\":\""+ newServiceDescription1 + "\",\"ServiceComposition\":"+ serviceCompositionJsonRappresentation + "}]"
	checkQuery(t, mockStub, functionName, serviceName, expectedResp)

}

func TestDeleteService(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Delete Service", simpleChaincode)


	// CREATION OF SERVICE 1:
	var functionAndArgsCreateService1 []string
	createServiceFunctionName := CreateService
	newServiceId1 := NewServiceId
	newServiceName1 := NewServiceName
	newServiceDescription1 := NewServiceDescription
	serviceCompositionAsString1 := "asd,fda"
	serviceComposition := lib.ParseStringToStringSlice(serviceCompositionAsString1)
	args1 := []string{newServiceId1, newServiceName1, newServiceDescription1, serviceCompositionAsString1}
	functionAndArgsCreateService1 = append(functionAndArgsCreateService1, createServiceFunctionName)
	functionAndArgsCreateService1 = append(functionAndArgsCreateService1,args1...)

	checkInvoke(t, mockStub, functionAndArgsCreateService1)

	serviceCompositionJsonRappresentation := "["
	for i := 0; i<len(serviceComposition) ;i++  {
		if i == 0 {
			serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation + "\""+serviceComposition[i]+ "\""
		}else {
			serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation + ",\""+serviceComposition[i]+ "\""
		}
	}
	serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation +"]"


	// VERIFY THE QUERY BEFORE THE DELETE GetServicesByName with the newly created services
	//   0
	// "serviceName"
	var functionAndArgs []string
	functionName:= GetService

	args3 := []string{newServiceId1}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args3...)
	// {"ServiceId":"idservice6","Name":"service6","Description":"service Description 6","ServiceComposition":["asd","fda"]}
	expectedRespBeforeDelete := "{\"ServiceId\":\""+ newServiceId1 + "\",\"Name\":\""+ newServiceName1 + "\",\"Description\":\""+ newServiceDescription1 + "\",\"ServiceComposition\":"+ serviceCompositionJsonRappresentation + "}"
	checkQuery(t, mockStub, functionName, newServiceId1, expectedRespBeforeDelete)


	// DELETE THE SERVICE RELATION AGENT
	var functionAndArgsDelete []string

	functionNameDelete := DeleteService
	argsDelete := []string{newServiceId1}

	functionAndArgsDelete = append(functionAndArgsDelete, functionNameDelete)
	functionAndArgsDelete = append(functionAndArgsDelete, argsDelete...)

	checkInvoke(t, mockStub, functionAndArgsDelete)


	// VERIFY THE QUERY AFTER THE DELETE GetServicesByName with the newly created services
	//   0
	// "serviceName"
	var functionAndArgs2 []string

	functionAndArgs2 = append(functionAndArgs2, functionName)
	functionAndArgs2 = append(functionAndArgs2, args3...)


	expectedRespAfterDelete := "{\"ServiceId\":\"\",\"Name\":\"\",\"Description\":\"\",\"ServiceComposition\":null}"
	checkQuery(t, mockStub, functionName, newServiceId1, expectedRespAfterDelete)

}

func TestDeleteServiceRelationAgent(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Delete ServiceRelationAgent", simpleChaincode)

	// CREATION OF AGENT 1:
	var functionAndArgsAgentCreation []string
	createAgentFunctionName := CreateAgent
	newAgentId1 := NewAgentId
	newAgentName := NewAgentName
	newAgentAddress := NewAgentAddress
	args := []string{newAgentId1, newAgentName, newAgentAddress}
	functionAndArgsAgentCreation = append(functionAndArgsAgentCreation, createAgentFunctionName)
	functionAndArgsAgentCreation = append(functionAndArgsAgentCreation, args...)

	checkInvoke(t, mockStub, functionAndArgsAgentCreation)


	// CREATION OF SERVICE 1:
	var functionAndArgsCreateService1 []string
	createServiceFunctionName := CreateService
	newServiceId1 := NewServiceId
	sameServiceName := NewServiceName
	newServiceDescription1 := NewServiceDescription
	serviceCompositionAsString1 := "asd,fda"
	serviceComposition := lib.ParseStringToStringSlice(serviceCompositionAsString1)
	args1 := []string{newServiceId1, sameServiceName, newServiceDescription1, serviceCompositionAsString1}
	functionAndArgsCreateService1 = append(functionAndArgsCreateService1, createServiceFunctionName)
	functionAndArgsCreateService1 = append(functionAndArgsCreateService1,args1...)

	checkInvoke(t, mockStub, functionAndArgsCreateService1)

	serviceCompositionJsonRappresentation := "["
	for i := 0; i<len(serviceComposition) ;i++  {
		if i == 0 {
			serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation + "\""+serviceComposition[i]+ "\""
		}else {
			serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation + ",\""+serviceComposition[i]+ "\""
		}
	}
	serviceCompositionJsonRappresentation = serviceCompositionJsonRappresentation +"]"

	// CREATION OF SERVICE RELATION AGENT 1:
	var functionAndArgsServiceRelationAgentCreation []string
	createServiceRelationAgentFunctionName := CreateServiceAgentRelation
	newServiceRelationAgentId := NewServiceId + NewAgentId
	newServiceId := NewServiceId
	newAgentId := NewAgentId
	newCost := "7"
	newTime := "3"
	args2 := []string{ newServiceId, newAgentId, newCost, newTime}
	functionAndArgsServiceRelationAgentCreation = append(functionAndArgsServiceRelationAgentCreation, createServiceRelationAgentFunctionName)
	functionAndArgsServiceRelationAgentCreation = append(functionAndArgsServiceRelationAgentCreation, args2...)

	checkInvoke(t, mockStub, functionAndArgsServiceRelationAgentCreation)

	// VERIFY THE QUERY BEFORE THE DELETE GetServicesByName with the newly created services
	//   0
	// "serviceName"
	var functionAndArgs []string
	functionName:= GetServiceRelationAgent

	args3 := []string{newServiceRelationAgentId}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args3...)
	// {"RelationId":"idservice6idagent6","ServiceId":"idservice6","AgentId":"idagent6","Cost":"7","Time":"3"}
	expectedRespBeforeDelete := "{\"RelationId\":\""+ newServiceRelationAgentId + "\",\"ServiceId\":\""+ newServiceId + "\",\"AgentId\":\""+ newAgentId + "\",\"Cost\":\""+ newCost + "\",\"Time\":\""+ newTime + "\"}"
	checkQuery(t, mockStub, functionName, newServiceRelationAgentId, expectedRespBeforeDelete)


	// DELETE THE SERVICE RELATION AGENT
	var functionAndArgsDelete []string

	functionNameDelete := DeleteServiceRelationAgent
	argsDelete := []string{newServiceRelationAgentId}

	functionAndArgsDelete = append(functionAndArgsDelete, functionNameDelete)
	functionAndArgsDelete = append(functionAndArgsDelete, argsDelete...)

	checkInvoke(t, mockStub, functionAndArgsDelete)


	// VERIFY THE QUERY AFTER THE DELETE GetServicesByName with the newly created services
	//   0
	// "serviceName"
	var functionAndArgs2 []string

	functionAndArgs2 = append(functionAndArgs2, functionName)
	functionAndArgs2 = append(functionAndArgs2, args3...)

	expectedRespAfterDelete := "{\"RelationId\":\"\",\"ServiceId\":\"\",\"AgentId\":\"\",\"Cost\":\"\",\"Time\":\"\"}"
	checkQuery(t, mockStub, functionName, newServiceRelationAgentId, expectedRespAfterDelete)

	// VERIFY THE QUERY ON THE INDEX AFTER THE DELETE GetServicesByName with the newly created services
	//   0
	// "serviceName"
	testLog.Info("Test the GetAgentsByService Query")
	var functionAndArgs3 []string

	functionNameIndexQuery := GetAgentsByService

	args4 := []string{newServiceId}

	functionAndArgs3 = append(functionAndArgs3, functionNameIndexQuery)
	functionAndArgs3 = append(functionAndArgs3, args4...)

	expectedRespAfterDeleteOnIndex := "{\"RelationId\":\"\",\"ServiceId\":\"\",\"AgentId\":\"\",\"Cost\":\"\",\"Time\":\"\"}"
	checkQuery(t, mockStub, functionName, newServiceRelationAgentId, expectedRespAfterDeleteOnIndex)

	// VERIFY THE QUERY ON THE INDEX AFTER THE DELETE GetServicesByName with the newly created services
	//   0
	// "serviceName"
	testLog.Info("Test the GetAgentsByService Query")
	var functionAndArgs4 []string

	functionNameIndexGetServicesByAgentQuery := GetServicesByAgent

	args5 := []string{newServiceId}

	functionAndArgs4 = append(functionAndArgs4, functionNameIndexGetServicesByAgentQuery)
	functionAndArgs4 = append(functionAndArgs4, args5...)

	expectedRespAfterDeleteOnIndex2 := "{\"RelationId\":\"\",\"ServiceId\":\"\",\"AgentId\":\"\",\"Cost\":\"\",\"Time\":\"\"}"
	checkQuery(t, mockStub, functionName, newServiceRelationAgentId, expectedRespAfterDeleteOnIndex2)

}

/*
func TestTradeWorkflow_LetterOfCredit(t *testing.T) {
	scc := new(TradeWorkflowChaincode)
	scc.testMode = true
	stub := shim.NewMockStub("Trade Workflow", scc)

	// Init
	checkInit(t, stub, getInitArguments())

	// Invoke 'requestTrade' and 'acceptTrade'
	tradeID := "2ks89j9"
	amount := 50000
	descGoods := "Wood for Toys"
	checkInvoke(t, stub, [][]byte{[]byte("requestTrade"), []byte(tradeID), []byte(strconv.Itoa(amount)), []byte(descGoods)})
	checkInvoke(t, stub, [][]byte{[]byte("acceptTrade"), []byte(tradeID)})

	// Invoke 'requestLC'
	checkInvoke(t, stub, [][]byte{[]byte("requestLC"), []byte(tradeID)})
	letterOfCredit := &LetterOfCredit{"", "", EXPORTER, amount, []string{}, REQUESTED}
	letterOfCreditBytes, _ := json.Marshal(letterOfCredit)
	lcKey, _ := stub.CreateCompositeKey("LetterOfCredit", []string{tradeID})
	checkState(t, stub, lcKey, string(letterOfCreditBytes))

	expectedResp := "{\"Status\":\"REQUESTED\"}"
	checkQuery(t, stub, "getLCStatus", tradeID, expectedResp)

	// Invoke bad 'issueLC' and verify unchanged state
	checkBadInvoke(t, stub, [][]byte{[]byte("issueLC")})
	badTradeID := "abcd"
	checkBadInvoke(t, stub, [][]byte{[]byte("issueLC"), []byte(badTradeID)})
	checkState(t, stub, lcKey, string(letterOfCreditBytes))

	// Invoke 'acceptLC' prematurely and verify failure and unchanged state
	checkBadInvoke(t, stub, [][]byte{[]byte("acceptLC"), []byte(badTradeID)})
	checkState(t, stub, lcKey, string(letterOfCreditBytes))
	checkQuery(t, stub, "getLCStatus", tradeID, expectedResp)

	// Invoke 'issueLC'
	lcID := "lc8349"
	expirationDate := "12/31/2018"
	doc1 := "E/L"
	doc2 := "B/L"
	checkInvoke(t, stub, [][]byte{[]byte("issueLC"), []byte(tradeID), []byte(lcID), []byte(expirationDate), []byte(doc1), []byte(doc2)})
	letterOfCredit = &LetterOfCredit{lcID, expirationDate, EXPORTER, amount, []string{doc1, doc2}, ISSUED}
	letterOfCreditBytes, _ = json.Marshal(letterOfCredit)
	checkState(t, stub, lcKey, string(letterOfCreditBytes))

	expectedResp = "{\"Status\":\"ISSUED\"}"
	checkQuery(t, stub, "getLCStatus", tradeID, expectedResp)

	// Invoke 'acceptLC'
	checkInvoke(t, stub, [][]byte{[]byte("acceptLC"), []byte(tradeID)})
	letterOfCredit = &LetterOfCredit{lcID, expirationDate, EXPORTER, amount, []string{doc1, doc2}, ACCEPTED}
	letterOfCreditBytes, _ = json.Marshal(letterOfCredit)
	checkState(t, stub, lcKey, string(letterOfCreditBytes))

	expectedResp = "{\"Status\":\"ACCEPTED\"}"
	checkQuery(t, stub, "getLCStatus", tradeID, expectedResp)
}

func TestTradeWorkflow_ExportLicense(t *testing.T) {
	scc := new(TradeWorkflowChaincode)
	scc.testMode = true
	stub := shim.NewMockStub("Trade Workflow", scc)

	// Init
	checkInit(t, stub, getInitArguments())

	// Invoke 'requestTrade', 'acceptTrade', 'requestLC', 'issueLC', 'acceptLC'
	tradeID := "2ks89j9"
	amount := 50000
	descGoods := "Wood for Toys"
	checkInvoke(t, stub, [][]byte{[]byte("requestTrade"), []byte(tradeID), []byte(strconv.Itoa(amount)), []byte(descGoods)})
	checkInvoke(t, stub, [][]byte{[]byte("acceptTrade"), []byte(tradeID)})
	checkInvoke(t, stub, [][]byte{[]byte("requestLC"), []byte(tradeID)})
	lcID := "lc8349"
	lcExpirationDate := "12/31/2018"
	doc1 := "E/L"
	doc2 := "B/L"
	checkInvoke(t, stub, [][]byte{[]byte("issueLC"), []byte(tradeID), []byte(lcID), []byte(lcExpirationDate), []byte(doc1), []byte(doc2)})
	checkInvoke(t, stub, [][]byte{[]byte("acceptLC"), []byte(tradeID)})

	// Issue 'requestEL'
	checkInvoke(t, stub, [][]byte{[]byte("requestEL"), []byte(tradeID)})
	exportLicense := &ExportLicense{"", "", EXPORTER, CARRIER, descGoods, REGAUTH, REQUESTED}
	exportLicenseBytes, _ := json.Marshal(exportLicense)
	elKey, _ := stub.CreateCompositeKey("ExportLicense", []string{tradeID})
	checkState(t, stub, elKey, string(exportLicenseBytes))

	expectedResp := "{\"Status\":\"REQUESTED\"}"
	checkQuery(t, stub, "getELStatus", tradeID, expectedResp)

	elID := "el979"
	elExpirationDate := "4/30/2019"

	// Invoke bad 'issueEL' and verify unchanged state
	checkBadInvoke(t, stub, [][]byte{[]byte("issueEL")})
	badTradeID := "abcd"
	checkBadInvoke(t, stub, [][]byte{[]byte("issueEL"), []byte(badTradeID), []byte(elID), []byte(elExpirationDate)})
	checkState(t, stub, elKey, string(exportLicenseBytes))
	checkQuery(t, stub, "getELStatus", tradeID, expectedResp)

	// Invoke 'issueEL' and verify state change
	checkInvoke(t, stub, [][]byte{[]byte("issueEL"), []byte(tradeID), []byte(elID), []byte(elExpirationDate)})
	exportLicense = &ExportLicense{elID, elExpirationDate, EXPORTER, CARRIER, descGoods, REGAUTH, ISSUED}
	exportLicenseBytes, _ = json.Marshal(exportLicense)
	checkState(t, stub, elKey, string(exportLicenseBytes))

	expectedResp = "{\"Status\":\"ISSUED\"}"
	checkQuery(t, stub, "getELStatus", tradeID, expectedResp)
}

func TestTradeWorkflow_ShipmentInitiation(t *testing.T) {
	scc := new(TradeWorkflowChaincode)
	scc.testMode = true
	stub := shim.NewMockStub("Trade Workflow", scc)

	// Init
	checkInit(t, stub, getInitArguments())

	// Invoke 'requestTrade', 'acceptTrade', 'requestLC', 'issueLC', 'acceptLC', 'requestEL', 'issueEL'
	tradeID := "2ks89j9"
	amount := 50000
	descGoods := "Wood for Toys"
	checkInvoke(t, stub, [][]byte{[]byte("requestTrade"), []byte(tradeID), []byte(strconv.Itoa(amount)), []byte(descGoods)})
	checkInvoke(t, stub, [][]byte{[]byte("acceptTrade"), []byte(tradeID)})
	checkInvoke(t, stub, [][]byte{[]byte("requestLC"), []byte(tradeID)})
	lcID := "lc8349"
	lcExpirationDate := "12/31/2018"
	doc1 := "E/L"
	doc2 := "B/L"
	checkInvoke(t, stub, [][]byte{[]byte("issueLC"), []byte(tradeID), []byte(lcID), []byte(lcExpirationDate), []byte(doc1), []byte(doc2)})
	checkInvoke(t, stub, [][]byte{[]byte("acceptLC"), []byte(tradeID)})
	checkInvoke(t, stub, [][]byte{[]byte("requestEL"), []byte(tradeID)})
	elID := "el979"
	elExpirationDate := "4/30/2019"
	checkInvoke(t, stub, [][]byte{[]byte("issueEL"), []byte(tradeID), []byte(elID), []byte(elExpirationDate)})

	// Invoke 'prepareShipment'
	checkInvoke(t, stub, [][]byte{[]byte("prepareShipment"), []byte(tradeID)})
	slKey, _ := stub.CreateCompositeKey("Shipment", []string{"Location", tradeID})
	checkState(t, stub, slKey, SOURCE)

	expectedResp := "{\"Location\":\"SOURCE\"}"
	checkQuery(t, stub, "getShipmentLocation", tradeID, expectedResp)

	// Invoke bad 'acceptShipmentAndIssueBL' and verify unchanged state
	checkBadInvoke(t, stub, [][]byte{[]byte("acceptShipmentAndIssueBL")})
	badTradeID := "abcd"
	blID := "bl06678"
	blExpirationDate := "8/31/2018"
	sourcePort := "Woodlands Port"
	destinationPort := "Market Port"
	checkBadInvoke(t, stub, [][]byte{[]byte("acceptShipmentAndIssueBL"), []byte(badTradeID), []byte(blID), []byte(blExpirationDate), []byte(sourcePort), []byte(destinationPort)})
	blKey, _ := stub.CreateCompositeKey("BillOfLading", []string{tradeID})
	checkNoState(t, stub, blKey)
	checkBadQuery(t, stub, "getBillOfLading", tradeID)

	// Invoke 'acceptShipmentAndIssueBL' and verify state change
	checkInvoke(t, stub, [][]byte{[]byte("acceptShipmentAndIssueBL"), []byte(tradeID), []byte(blID), []byte(blExpirationDate), []byte(sourcePort), []byte(destinationPort)})
	billOfLading := &BillOfLading{blID, blExpirationDate, EXPORTER, CARRIER, descGoods, amount, IMPBANK, sourcePort, destinationPort}
	billOfLadingBytes, _ := json.Marshal(billOfLading)
	checkState(t, stub, blKey, string(billOfLadingBytes))
	checkQuery(t, stub, "getBillOfLading", tradeID, string(billOfLadingBytes))
}

func TestTradeWorkflow_PaymentFulfilment(t *testing.T) {
	scc := new(TradeWorkflowChaincode)
	scc.testMode = true
	stub := shim.NewMockStub("Trade Workflow", scc)

	// Init
	checkInit(t, stub, getInitArguments())

	// Invoke 'requestTrade', 'acceptTrade', 'requestLC', 'issueLC', 'acceptLC', 'requestEL', 'issueEL', 'prepareShipment', 'acceptShipmentAndIssueBL'
	tradeID := "2ks89j9"
	amount := 50000
	descGoods := "Wood for Toys"
	checkInvoke(t, stub, [][]byte{[]byte("requestTrade"), []byte(tradeID), []byte(strconv.Itoa(amount)), []byte(descGoods)})
	checkInvoke(t, stub, [][]byte{[]byte("acceptTrade"), []byte(tradeID)})
	checkInvoke(t, stub, [][]byte{[]byte("requestLC"), []byte(tradeID)})
	lcID := "lc8349"
	lcExpirationDate := "12/31/2018"
	doc1 := "E/L"
	doc2 := "B/L"
	checkInvoke(t, stub, [][]byte{[]byte("issueLC"), []byte(tradeID), []byte(lcID), []byte(lcExpirationDate), []byte(doc1), []byte(doc2)})
	checkInvoke(t, stub, [][]byte{[]byte("acceptLC"), []byte(tradeID)})
	checkInvoke(t, stub, [][]byte{[]byte("requestEL"), []byte(tradeID)})
	elID := "el979"
	elExpirationDate := "4/30/2019"
	checkInvoke(t, stub, [][]byte{[]byte("issueEL"), []byte(tradeID), []byte(elID), []byte(elExpirationDate)})
	checkInvoke(t, stub, [][]byte{[]byte("prepareShipment"), []byte(tradeID)})
	blID := "bl06678"
	blExpirationDate := "8/31/2018"
	sourcePort := "Woodlands Port"
	destinationPort := "Market Port"
	checkInvoke(t, stub, [][]byte{[]byte("acceptShipmentAndIssueBL"), []byte(tradeID), []byte(blID), []byte(blExpirationDate), []byte(sourcePort), []byte(destinationPort)})

	// Invoke 'requestPayment'
	checkInvoke(t, stub, [][]byte{[]byte("requestPayment"), []byte(tradeID)})
	paymentKey, _ := stub.CreateCompositeKey("Payment", []string{tradeID})
	checkState(t, stub, paymentKey, REQUESTED)

	// Invoke 'makePayment'
	checkInvoke(t, stub, [][]byte{[]byte("makePayment"), []byte(tradeID)})
	checkNoState(t, stub, paymentKey)
	// Verify account and payment balances
	payment := amount/2
	expBalanceStr := strconv.Itoa(EXPBALANCE + payment)
	impBalanceStr := strconv.Itoa(IMPBALANCE - payment)
	checkState(t, stub, expBalKey, expBalanceStr)
	checkState(t, stub, impBalKey, impBalanceStr)
	tradeAgreement := &TradeAgreement{amount, descGoods, ACCEPTED, payment}
	tradeAgreementBytes, _ := json.Marshal(tradeAgreement)
	tradeKey, _ := stub.CreateCompositeKey("Trade", []string{tradeID})
	checkState(t, stub, tradeKey, string(tradeAgreementBytes))

	// Check queries
	checkBadQuery(t, stub, "getAccountBalance", tradeID)
	expectedResp := "{\"Balance\":\"" + expBalanceStr + "\"}"
	checkQueryArgs(t, stub, [][]byte{[]byte("getAccountBalance"), []byte(tradeID), []byte("exporter")}, expectedResp)

	expectedResp = "{\"Balance\":\"" + impBalanceStr + "\"}"
	checkQueryArgs(t, stub, [][]byte{[]byte("getAccountBalance"), []byte(tradeID), []byte("importer")}, expectedResp)

	// Deliver shipment to final location
	checkInvoke(t, stub, [][]byte{[]byte("updateShipmentLocation"), []byte(tradeID), []byte(DESTINATION)})
	slKey, _ := stub.CreateCompositeKey("Shipment", []string{"Location", tradeID})
	checkState(t, stub, slKey, DESTINATION)

	// Invoke 'requestPayment' and 'makePayment'
	checkInvoke(t, stub, [][]byte{[]byte("requestPayment"), []byte(tradeID)})
	checkState(t, stub, paymentKey, REQUESTED)
	checkInvoke(t, stub, [][]byte{[]byte("makePayment"), []byte(tradeID)})
	checkNoState(t, stub, paymentKey)

	// Verify account and payment balances, and check queries
	expBalanceStr = strconv.Itoa(EXPBALANCE + amount)
	impBalanceStr = strconv.Itoa(IMPBALANCE - amount)
	checkState(t, stub, expBalKey, expBalanceStr)
	checkState(t, stub, impBalKey, impBalanceStr)
	tradeAgreement = &TradeAgreement{amount, descGoods, ACCEPTED, amount}
	tradeAgreementBytes, _ = json.Marshal(tradeAgreement)
	checkState(t, stub, tradeKey, string(tradeAgreementBytes))

	expectedResp = "{\"Balance\":\"" + expBalanceStr + "\"}"
	checkQueryArgs(t, stub, [][]byte{[]byte("getAccountBalance"), []byte(tradeID), []byte("exporter")}, expectedResp)

	expectedResp = "{\"Balance\":\"" + impBalanceStr + "\"}"
	checkQueryArgs(t, stub, [][]byte{[]byte("getAccountBalance"), []byte(tradeID), []byte("importer")}, expectedResp)
}*/

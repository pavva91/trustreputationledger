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
	"fmt"
	lib "github.com/pavva91/arglib"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	a "github.com/pavva91/assets"
)

const (
	ExistingServiceId          = "idservice1"
	ExistingServiceName        = "service1"
	ExistingServiceDescription = "service Description 1"
	ExistingAgentId            = "idagent1"
	ExistingAgentName          = "agent1"
	ExistingAgentAddress       = "address1"
	NewServiceId               = "idservice6"
	NewServiceName             = "service6"
	NewServiceDescription      = "service Description 6"
	NewAgentId                 = "idagent6"
	NewAgentName               = "agent6"
	NewAgentAddress            = "address6"
	ServiceAgentServiceId      = ExistingServiceId
	ServiceAgentAgentId        = ExistingAgentId
	ServiceAgentCost           = "8"
	ServiceAgentTime           = "6"
	ReputationValue            = "6"

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
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}else{
		fmt.Println("Init OK", string(res.Message))

	}
}

func checkNoState(t *testing.T, stub *shim.MockStub, name string) {
	bytes := stub.State[name]
	if bytes != nil {
		fmt.Println("State", name, "should be absent; found value")
		t.FailNow()
	}else {
		fmt.Println("State", name, "is absent as it should be")
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was", string(bytes), "and not", value, "as expected")
		t.FailNow()
	}else{
		fmt.Println("State value", name, "is", string(bytes), "as expected")
	}
}

func checkBadQuery(t *testing.T, stub *shim.MockStub, function string, name string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status == shim.OK {
		fmt.Println("Query", name, "unexpectedly succeeded")
		t.FailNow()
	}else {
		fmt.Println("Query", name, "failed as espected")

	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, function string, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		fmt.Println("Query value", name, "was", payload, "and not", value, "as expected")
		t.FailNow()
	}else{
		fmt.Println("Query value", name, "is", payload, "as expected")
	}
}

func checkQueryArgs(t *testing.T, stub *shim.MockStub, args [][]byte, value string) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Query", string(args[1]), "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", string(args[1]), "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		fmt.Println("Query value", string(args[1]), "was", payload, "and not", value, "as expected")
		t.FailNow()
	}else {
		fmt.Println("Query value", string(args[1]), "is", payload, "as expected")

	}
}

func checkBadInvoke(t *testing.T, stub *shim.MockStub, functionAndArgs []string) {
	functionAndArgsAsBytes := lib.ParseStringSliceToByteSlice(functionAndArgs)
	res := stub.MockInvoke("1", functionAndArgsAsBytes)
	if res.Status == shim.OK {
		fmt.Println("Invoke", functionAndArgs, "unexpectedly succeeded")
		t.FailNow()
	}else {
		fmt.Println("Invoke", functionAndArgs, "failed as espected")
	}
}

// func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
// 	res := stub.MockInvoke("1", args)
// 	if res.Status != shim.OK {
// 		fmt.Println("Invoke", args, "failed", string(res.Message))
// 		t.FailNow()
// 	}else {
// 		fmt.Println("Invoke", args, "successful", string(res.Message))
// 	}
// }
func checkInvoke(t *testing.T, stub *shim.MockStub, functionAndArgs []string) {
	functionAndArgsAsBytes := lib.ParseStringSliceToByteSlice(functionAndArgs)
	res := stub.MockInvoke("1", functionAndArgsAsBytes)
	if res.Status != shim.OK {
		fmt.Println("Invoke", functionAndArgs, "failed", string(res.Message))
		t.FailNow()
	}else {
		fmt.Println("Invoke", functionAndArgs, "successful", string(res.Message))
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

	// Invoke 'CreateService'
	serviceId := NewServiceId
	serviceName := NewServiceName
	serviceDescription := NewServiceDescription

	args := []string{serviceId,serviceName,serviceDescription}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkInvoke(t, mockStub, functionAndArgs)

	service := &a.Service{serviceId, serviceName, serviceDescription}
	serviceAsBytes, _ := json.Marshal(service)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{serviceId})
	checkState(t, mockStub, serviceId, string(serviceAsBytes))

	expectedResp := "{\"ServiceId\":\""+ serviceId + "\",\"Name\":\""+ serviceName + "\",\"Description\":\""+ serviceDescription + "\"}"
	checkQuery(t, mockStub, "GetServiceNotFoundError", serviceId, expectedResp)
}
// =====================================================================================================================
// TestServiceCreation - Test the 'CreateService' function when trying to insert an already existing record
// =====================================================================================================================
func TestExistingServiceCreation(t *testing.T) {
	simpleChaincode := new(SimpleChaincode)
	simpleChaincode.testMode = true
	mockStub := shim.NewMockStub("Test Already Existing Service Creation", simpleChaincode)

	// Init
	checkInit(t, mockStub, getInitArguments())

	var functionAndArgs []string
	functionName:= CreateService

	// Invoke 'CreateService'
	existingServiceId := ExistingServiceId
	serviceName := ExistingServiceName
	serviceDescription := ExistingServiceDescription

	args := []string{existingServiceId,serviceName,serviceDescription}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs,args...)

	checkBadInvoke(t, mockStub, functionAndArgs)


	service := &a.Service{existingServiceId, serviceName, serviceDescription}
	serviceBytes, _ := json.Marshal(service)
	// tradeKey, _ := mockStub.CreateCompositeKey("Trade", []string{existingServiceId})
	checkState(t, mockStub, existingServiceId, string(serviceBytes))

	expectedResp := "{\"ServiceId\":\""+ existingServiceId + "\",\"Name\":\""+ serviceName + "\",\"Description\":\""+ serviceDescription + "\"}"
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

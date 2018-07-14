/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import ("github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"encoding/json"
	"bytes"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"strconv"
)

// ===============================================
// getValue - get a generic variable from ledger
// ===============================================
func getValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var agentId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting the key of the value to query")
	}

	agentId = args[0]
	valAsbytes, err := stub.GetState(agentId) //get the agent from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + agentId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Agent does not exist: " + agentId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}





//
// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
//
// 	"github.com/hyperledger/fabric/core/chaincode/shim"
// 	pb "github.com/hyperledger/fabric/protos/peer"
// )
//

// ============================================================================================================================
// Read - read a generic variable from ledger
//
// Shows Off GetState() - reading a key/value from the ledger
//
// Inputs - Array of strings
//  0
//  key
//  "abc"
//
// Returns Payload:
// SUCCESS (found key value): shim.Success(json.RawMessage)
// FAIL (not found key-value): shim.Error
// ============================================================================================================================
func read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	var outJson json.RawMessage

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key of the var to query")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)   //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	// Trasformo risposta da bytes a JSON (così ritorna null in caso di risultato vuoto)
	json.Unmarshal(valAsbytes, &outJson)
	out, _ := json.Marshal(outJson)


	// We are crazy, we work directly with []byte :P
	// if bytes.Equal(out,[]byte{110,117,108,108}){
	// 	return shim.Error("Key not found in the Ledger")
	// }

	// Normal people work with string
	stringOut := string(out)
	fmt.Print("Raw bytes: ")
	fmt.Println(out)
	fmt.Println("String: " + stringOut)
	if stringOut == "null" {
		return shim.Error("Key not found in the Ledger")

	}

	return shim.Success(out)                  //send it onward
}

// ============================================================================================================================
// Get everything we need (agents + services)
//
// Inputs - none
//
// Returns:
// }
// ============================================================================================================================
func readEverything(stub shim.ChaincodeStubInterface) pb.Response {
	type Everything struct {
		Agents   []Agent   `json:"Agents"`
		Services []Service `json:"Services"`
	}
	var everything Everything

	// ---- Get All Services ---- //
	resultsIterator, err := stub.GetStateByRange("idservice0", "idservice9999999999999999999")
	// resultsIterator, err := stub.GetStateByRange("", "")

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		aKeyValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on service id - ", queryKeyAsStr)
		var service Service
		json.Unmarshal(queryValAsBytes, &service)                  //un stringify it aka JSON.parse()
		everything.Services = append(everything.Services, service) //add this service to the list
	}
	fmt.Println("service array - ", everything.Services)

	// ---- Get All Agents ---- //
	ownersIterator, err := stub.GetStateByRange("idagent0", "idagent9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer ownersIterator.Close()

	for ownersIterator.HasNext() {
		aKeyValue, err := ownersIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on agent id - ", queryKeyAsStr)
		var agent Agent
		json.Unmarshal(queryValAsBytes, &agent) //un stringify it aka JSON.parse()
		everything.Agents = append(everything.Agents, agent) //add this service to the list
	}
	fmt.Println("agent array - ", everything.Agents)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything)              //convert to array of bytes
	return shim.Success(everythingAsBytes)
}

// ============================================================================================================================
// Get all the ledger
//
// Inputs - none
//
// Returns:
// }
// ============================================================================================================================
func readAllLedger(stub shim.ChaincodeStubInterface) pb.Response {

	var buffer bytes.Buffer
	buffer.WriteString("[")

	// ---- Get All the ledger ---- //
	resultsIterator, err := stub.GetStateByRange("", "")

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		aKeyValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryResultKey := aKeyValue.Key
		queryResultValue := aKeyValue.Value
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
			buffer.WriteString("\n")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResultKey)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResultValue))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getMarblesByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
// ============================================================================================================================
// Get history of a general asset
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "m01490985296352SjAyM"
// ============================================================================================================================
func getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var history []queryresult.KeyModification
	var service Service

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	serviceId := args[0]
	fmt.Printf("- start getHistory: %s\n", serviceId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(serviceId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx queryresult.KeyModification
		tx.TxId = historyData.TxId                  //copy transaction id over
		json.Unmarshal(historyData.Value, &service) //un stringify it aka JSON.parse()
		if historyData.Value == nil {                  //service has been deleted
			var emptyBytes []byte
			tx.Value = emptyBytes //copy nil service
		} else {
			json.Unmarshal(historyData.Value, &service) //un stringify it aka JSON.parse()
			tx.Value = historyData.Value //copy service over
			tx.Timestamp = historyData.Timestamp
			tx.IsDelete = historyData.IsDelete
		}
		history = append(history, tx)              //add this tx to the list
	}
	// fmt.Printf("- getHistoryForService returning:\n%s", history)
	prettyPrintHistory(history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}

func prettyPrintHistory(history []queryresult.KeyModification){
	for i := 0; i< len(history); i++ {
		fmt.Printf("Value version: %s:\n",strconv.Itoa(i))
		fmt.Println("Timestamp: " + history[i].Timestamp.String())
		fmt.Println("Value: " + string(history[i].Value))
		fmt.Println("TxId: " + history[i].TxId)
		fmt.Println("IsDelete: " + strconv.FormatBool(history[i].IsDelete))
		fmt.Println("=====================================================================")
	}
}

// ============================================================================================================================
// Get history of service
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "m01490985296352SjAyM"
// ============================================================================================================================
func getServiceHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type ServiceHistory struct {
		TxId    string   `json:"txId"`
		Value   Service   `json:"value"`
		// Timestamp *google_protobuf.Timestamp
		IsDelete bool `json:"isDelete"`
	}
	var history []ServiceHistory
	var service Service

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	serviceId := args[0]
	fmt.Printf("- start getHistoryForService: %s\n", serviceId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(serviceId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx ServiceHistory
		tx.TxId = historyData.TxId                  //copy transaction id over
		json.Unmarshal(historyData.Value, &service) //un stringify it aka JSON.parse()
		if historyData.Value == nil {                  //service has been deleted
			var emptyService Service
			tx.Value = emptyService //copy nil service
		} else {
			json.Unmarshal(historyData.Value, &service) //un stringify it aka JSON.parse()
			tx.Value = service //copy service over
			tx.IsDelete = historyData.IsDelete
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForService returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}

// ============================================================================================================================
// Get history of service
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "m01490985296352SjAyM"
// ============================================================================================================================
func getAgentHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AgentHistory struct {
		TxId    string   `json:"txId"`
		Value   Agent   `json:"value"`
		IsDelete bool `json:"isDelete"`
	}
	var history []AgentHistory
	var agent Agent

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	serviceId := args[0]
	fmt.Printf("- start getHistoryForAgent: %s\n", serviceId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(serviceId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AgentHistory
		tx.TxId = historyData.TxId                //copy transaction id over
		json.Unmarshal(historyData.Value, &agent) //un stringify it aka JSON.parse()
		if historyData.Value == nil {                  //agent has been deleted
			var emptyAgent Agent
			tx.Value = emptyAgent //copy nil agent
		} else {
			json.Unmarshal(historyData.Value, &agent) //un stringify it aka JSON.parse()
			tx.Value = agent                          //copy agent over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForAgent returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}

// ============================================================================================================================
// Get history of ServiceRelationAgent
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "m01490985296352SjAyM"
// ============================================================================================================================
func getServiceRelationAgentHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId    string   `json:"txId"`
		Value   ServiceRelationAgent   `json:"value"`
	}
	var history []AuditHistory
	var serviceRelationAgent ServiceRelationAgent

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	relationId := args[0]
	fmt.Printf("- start getHistoryForServiceRelationAgent: %s\n", relationId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(relationId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxId = historyData.TxId                               //copy transaction id over
		json.Unmarshal(historyData.Value, &serviceRelationAgent) //un stringify it aka JSON.parse()
		if historyData.Value == nil {                  //serviceRelationAgent has been deleted
			var emptyServiceRelationAgent ServiceRelationAgent
			tx.Value = emptyServiceRelationAgent //copy nil serviceRelationAgent
		} else {
			json.Unmarshal(historyData.Value, &serviceRelationAgent) //un stringify it aka JSON.parse()
			tx.Value = serviceRelationAgent                          //copy serviceRelationAgent over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForServiceRelationAgent returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}


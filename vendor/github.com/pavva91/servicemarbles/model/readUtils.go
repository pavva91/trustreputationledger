package model

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"encoding/json"
	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/pavva91/servicemarbles/generalcc"
)

// ============================================================================================================================
// Get everything we need (agents + services)
//
// Inputs - none
//
// Returns:
// }
// ============================================================================================================================
func ReadEverything(stub shim.ChaincodeStubInterface) pb.Response {
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
		json.Unmarshal(queryValAsBytes, &agent)              //un stringify it aka JSON.parse()
		everything.Agents = append(everything.Agents, agent) //add this service to the list
	}
	fmt.Println("agent array - ", everything.Agents)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything) //convert to array of bytes
	return shim.Success(everythingAsBytes)
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
func GetServiceHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type ServiceHistory struct {
		TxId  string  `json:"txId"`
		Value Service `json:"value"`
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
		if historyData.Value == nil {               //service has been deleted
			var emptyService Service
			tx.Value = emptyService //copy nil service
		} else {
			json.Unmarshal(historyData.Value, &service) //un stringify it aka JSON.parse()
			tx.Value = service                          //copy service over
			tx.IsDelete = historyData.IsDelete
		}
		history = append(history, tx) //add this tx to the list
	}
	fmt.Printf("- getHistoryForService returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history) //convert to array of bytes
	return shim.Success(historyAsBytes)
}

// ============================================================================================================================
// Get history of agent
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "m01490985296352SjAyM"
// ============================================================================================================================
func GetAgentHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AgentHistory struct {
		TxId     string `json:"txId"`
		Value    Agent  `json:"value"`
		IsDelete bool   `json:"isDelete"`
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
		if historyData.Value == nil {             //agent has been deleted
			var emptyAgent Agent
			tx.Value = emptyAgent //copy nil agent
		} else {
			json.Unmarshal(historyData.Value, &agent) //un stringify it aka JSON.parse()
			tx.Value = agent                          //copy agent over
		}
		history = append(history, tx) //add this tx to the list
	}
	fmt.Printf("- getHistoryForAgent returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history) //convert to array of bytes
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
func GetServiceRelationAgentHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId  string               `json:"txId"`
		Value ServiceRelationAgent `json:"value"`
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
		if historyData.Value == nil {                            //serviceRelationAgent has been deleted
			var emptyServiceRelationAgent ServiceRelationAgent
			tx.Value = emptyServiceRelationAgent //copy nil serviceRelationAgent
		} else {
			json.Unmarshal(historyData.Value, &serviceRelationAgent) //un stringify it aka JSON.parse()
			tx.Value = serviceRelationAgent                          //copy serviceRelationAgent over
		}
		history = append(history, tx) //add this tx to the list
	}
	fmt.Printf("- getHistoryForServiceRelationAgent returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history) //convert to array of bytes
	return shim.Success(historyAsBytes)
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
func GetHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var history []queryresult.KeyModification
	var service Service

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]
	fmt.Printf("- start GetHistory: %s\n", key)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(key)
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
		if historyData.Value == nil {               //service has been deleted
			var emptyBytes []byte
			tx.Value = emptyBytes //copy nil service
		} else {
			json.Unmarshal(historyData.Value, &service) //un stringify it aka JSON.parse()
			tx.Value = historyData.Value                //copy service over
			tx.Timestamp = historyData.Timestamp
			tx.IsDelete = historyData.IsDelete
		}
		history = append(history, tx) //add this tx to the list
	}
	// fmt.Printf("- getHistoryForService returning:\n%s", history)
	generalcc.PrettyPrintHistory(history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history) //convert to array of bytes
	return shim.Success(historyAsBytes)
}

/*
Package generalcc implements a simple library for common fabric hyperledger's chaincode functions.
*/
/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
*/

package generalcc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pavva91/arglib"
	"strconv"
	)

// =====================================================================================================================
// GetValue - get a generic variable from ledger
// =====================================================================================================================
func GetValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

// =====================================================================================================================
// Read - Read a generic variable from ledger
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
// =====================================================================================================================
func Read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	var outJson json.RawMessage

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key of the var to query")
	}

	// input sanitation
	err = arglib.SanitizeArguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key) //get the var from ledger
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

	return shim.Success(out) //send it onward
}

// =====================================================================================================================
// Get all the Ledger's Current State Data (State Database) - The ledger’s current state data represents the latest
// values for all keys ever included in the chain transaction log.
// (https://hyperledger-fabric.readthedocs.io/en/release-1.1/ledger.html)
//
// Inputs - none
//
// Returns:
// }
// =====================================================================================================================
func ReadAllStateDB(stub shim.ChaincodeStubInterface) pb.Response {

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
		// Record is a JSON object, so we Write as-is
		buffer.WriteString(string(queryResultValue))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getMarblesByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}



// TODO: Trovare il modo di generalizzare senza usare assets.Service
// =====================================================================================================================
// Get history of a general asset in the Chain - The chain is a transaction log, structured as hash-linked blocks
// (https://hyperledger-fabric.readthedocs.io/en/release-1.1/ledger.html)
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "m01490985296352SjAyM"
// =====================================================================================================================
func GetHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type KeyModificationWrapper struct {
		RealValue interface{} `json:"InterfaceValue"`
		Tx        queryresult.KeyModification
	}
	var sliceReal []KeyModificationWrapper

	var history []queryresult.KeyModification
	var value interface{}

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

		var singleReal KeyModificationWrapper
		var tx queryresult.KeyModification
		singleReal.Tx.TxId = historyData.TxId     //copy transaction id over
		json.Unmarshal(historyData.Value, &value) //un stringify it aka JSON.parse()
		if historyData.Value == nil {             //value has been deleted
			var emptyBytes []byte
			singleReal.Tx.Value = emptyBytes //copy nil value
		} else {
			json.Unmarshal(historyData.Value, &value) //un stringify it aka JSON.parse()
			singleReal.Tx.Value = historyData.Value   //copy value over
			singleReal.Tx.Timestamp = historyData.Timestamp
			singleReal.Tx.IsDelete = historyData.IsDelete
			singleReal.RealValue = value
		}
		history = append(history, tx) //add this Tx to the list
		sliceReal = append(sliceReal, singleReal)
	}
	// fmt.Printf("- getHistoryForService returning:\n%s", history)
	PrettyPrintHistory(history)

	//change to array of bytes
	// historyAsBytes, _ := json.Marshal(history) //convert to array of bytes

	realAsBytes, _ := json.Marshal(sliceReal)
	return shim.Success(realAsBytes)
}

func PrettyPrintHistory(history []queryresult.KeyModification) {
	for i := 0; i < len(history); i++ {
		fmt.Printf("Value version: %s:\n", strconv.Itoa(i))
		fmt.Println("ExecutedServiceTimestamp: " + history[i].Timestamp.String())
		fmt.Println("Value: " + string(history[i].Value))
		fmt.Println("TxId: " + history[i].TxId)
		fmt.Println("IsDelete: " + strconv.FormatBool(history[i].IsDelete))
		fmt.Println("=====================================================================")
	}
}

// =====================================================================================================================
// Print Results Iterator - Print on screen the general iterator of the composite index query result
// =====================================================================================================================
func PrintResultsIterator(queryIterator shim.StateQueryIteratorInterface, stub shim.ChaincodeStubInterface) error {
	// USE DEFER BECAUSE it will close also in case of error throwing (premature return)
	defer queryIterator.Close()
	for i := 0; queryIterator.HasNext(); i++ {
		responseRange, err := queryIterator.Next()
		if err != nil {
			return err
		}
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return err
		}
		i := 0
		for _, keyPart := range compositeKeyParts {
			fmt.Printf("Found a Relation OBJECT_TYPE:%s KEYPART %s: %s", objectType, i, keyPart)
			i++
		}
	}
	return nil
}

func GetNextIncrementalKey(keyPrefix string, stub shim.ChaincodeStubInterface)(string, error ){
	// TODO: Levare nextIncrementalKey prefix di 3 lettere in testa

	startKey := keyPrefix+""
	endKey := keyPrefix+""
	i:=0

	//i need to get the last IncrementalKey on the ledger
	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return "",err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		resultsIterator.Next()
		i=i+1
	}

	nextIncrementalKey := keyPrefix + strconv.Itoa(i)
	return nextIncrementalKey,nil
}

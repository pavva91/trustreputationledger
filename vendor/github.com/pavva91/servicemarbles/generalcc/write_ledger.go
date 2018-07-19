/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
*/

package generalcc

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pavva91/arglib"
)

// ============================================================================================================================
// simpleWrite() - generic simpleWrite variable into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - Array of strings
//    0   ,    1
//   key  ,  value
//  "abc" , "test"
// ============================================================================================================================
func Write(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var key, value string
	var err error
	fmt.Println("starting simpleWrite")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	// input sanitation
	err = arglib.SanitizeArguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //simpleWrite the variable into the ledger
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end simpleWrite")
	return shim.Success(nil)
}

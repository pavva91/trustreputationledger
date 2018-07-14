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

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


// ============================================================
// initLedger - create a batch of new agents and services
// ============================================================
func initLedger(stub shim.ChaincodeStubInterface) pb.Response {
	services := []Service{
		Service{ServiceId:"idservice1", Name:"service1", Description: "service Description 1"},
		Service{ServiceId:"idservice2", Name:"service2", Description:"service Description 2"},
		Service{ServiceId:"idservice3", Name:"service3", Description:"service Description 3"},
		Service{ServiceId:"idservice4", Name:"service4", Description:"service Description 4"},
		Service{ServiceId:"idservice5", Name:"service5", Description:"service Description 5"},
	}
	agents := []Agent{
		Agent{AgentId:"idagent1", Name:"agent1", Address:"address1"},
		Agent{AgentId:"idagent2", Name:"agent2", Address:"address2"},
		Agent{AgentId:"idagent3", Name:"agent3", Address:"address3"},
		Agent{AgentId:"idagent4", Name:"agent4", Address:"address4"},
		Agent{AgentId:"idagent5", Name:"agent5", Address:"address5"},
	}

	// non funziona ( come chiamare, si può fare?)
	initServiceAgentRelation(stub, []string{"idservice1idagent1","idservice1","idagent1","5","3","9"})
	initServiceAgentRelation(stub, []string{"idservice1idagent2","idservice1","idagent2","6","2","8"})

	for i := 0; i < len(services); i++ {
		fmt.Println("i is ", i)
		serviceAsBytes, _ := json.Marshal(services[i])
		fmt.Println(serviceAsBytes)
		err:=stub.PutState(services[i].ServiceId, serviceAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("Addeds", services[i])
	}
	for i := 0; i < len(agents); i++ {
		fmt.Println("i is ", i)
		agentAsBytes, _ := json.Marshal(agents[i])
		err:=stub.PutState(agents[i].AgentId,agentAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("Added", agents[i])
	}

	return shim.Success(nil)
}


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
func write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	var err error
	fmt.Println("starting simpleWrite")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]                                   //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value))         //simpleWrite the variable into the ledger
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end simpleWrite")
	return shim.Success(nil)
}

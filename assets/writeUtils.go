/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
*/
package assets

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	"errors"
)

// ============================================================
// InitLedger - create a batch of new agents and services
// ============================================================
func InitLedger(stub shim.ChaincodeStubInterface) pb.Response {
	services := []Service{
		Service{ServiceId: "idservice1", Name: "service1", Description: "service Description 1"},
		Service{ServiceId: "idservice2", Name: "service2", Description: "service Description 2"},
		Service{ServiceId: "idservice3", Name: "service3", Description: "service Description 3"},
		Service{ServiceId: "idservice4", Name: "service4", Description: "service Description 4"},
		Service{ServiceId: "idservice5", Name: "service5", Description: "service Description 5"},
	}
	agents := []Agent{
		Agent{AgentId: "idagent1", Name: "agent1", Address: "address1"},
		Agent{AgentId: "idagent2", Name: "agent2", Address: "address2"},
		Agent{AgentId: "idagent3", Name: "agent3", Address: "address3"},
		Agent{AgentId: "idagent4", Name: "agent4", Address: "address4"},
		Agent{AgentId: "idagent5", Name: "agent5", Address: "address5"},
	}

	// non funziona ( come chiamare, si può fare?)
	// InitServiceAgentRelation(stub, []string{"idservice1idagent1", "idservice1", "idagent1", "5", "3", "9"})
	// InitServiceAgentRelation(stub, []string{"idservice1idagent2", "idservice1", "idagent2", "6", "2", "8"})

	for i := 0; i < len(services); i++ {
		fmt.Println("i is ", i)
		serviceAsBytes, _ := json.Marshal(services[i])
		fmt.Println(serviceAsBytes)
		err := stub.PutState(services[i].ServiceId, serviceAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("Addeds", services[i])
	}
	for i := 0; i < len(agents); i++ {
		fmt.Println("i is ", i)
		agentAsBytes, _ := json.Marshal(agents[i])
		err := stub.PutState(agents[i].AgentId, agentAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("Added", agents[i])
	}

	return shim.Success(nil)
}

func SaveIndex(indexKey string, stub shim.ChaincodeStubInterface) error {
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	//  Save index entry to state. Only the key Name is needed, no need to store a duplicate copy of the marble.
	value := []byte{0x00}
	// index save
	putStateError := stub.PutState(indexKey, value)
	if putStateError != nil {
		return errors.New(putStateError.Error())
	}
	return nil
}

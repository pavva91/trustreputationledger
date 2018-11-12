/*
Created by Valerio Mattioli @ HES-SO (valeriomattioli580@gmail.com
*/
package assets

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"errors"
)

var writeUtilsLog = shim.NewLogger("writeUtils")
// =====================================================================================================================
// InitLedger - create a batch of new agents and services (TEST)
// =====================================================================================================================
func InitLedger(stub shim.ChaincodeStubInterface) pb.Response {
	services := []Service{
		Service{ServiceId: "idservice1", Name: "service1", Description: "service Description 1"},
		Service{ServiceId: "idservice2", Name: "service2", Description: "service Description 2"},
		Service{ServiceId: "idservice3", Name: "service3", Description: "service Description 3"},
		Service{ServiceId: "idservice4", Name: "service4", Description: "service Description 4"},
		Service{ServiceId: "idservice5", Name: "service5", Description: "service Description 5"},
		Service{ServiceId: "idservice99", Name: "service99", Description: "service Description 99"},
	}
	agents := []Agent{
		Agent{AgentId: "idagent1", Name: "agent1", Address: "address1"},
		Agent{AgentId: "idagent2", Name: "agent2", Address: "address2"},
		Agent{AgentId: "idagent3", Name: "agent3", Address: "address3"},
		Agent{AgentId: "idagent4", Name: "agent4", Address: "address4"},
		Agent{AgentId: "idagent5", Name: "agent5", Address: "address5"},
		Agent{AgentId: "idagent98", Name: "agent98", Address: "address98"},
		Agent{AgentId: "idagent99", Name: "agent99", Address: "address99"},
	}
	serviceRelationAgents := []ServiceRelationAgent{
		ServiceRelationAgent{"idservice99idagent99","idservice99","idagent99" ,"5","7"},
	}
	reputations := []Reputation{
		Reputation{"idagent99idservice99EXECUTER","idagent99","idservice99" ,"EXECUTER","9"},
		Reputation{"idagent98idservice99DEMANDER","idagent98","idservice99" ,"DEMANDER","8"},
	}


	// non funziona ( come chiamare, si può fare?)
	// InitServiceAgentRelation(stub, []string{"idservice1idagent1", "idservice1", "idagent1", "5", "3", "9"})
	// InitServiceAgentRelation(stub, []string{"idservice1idagent2", "idservice1", "idagent2", "6", "2", "8"})

	for i := 0; i < len(services); i++ {
		serviceLog.Info("i is ", i)
		serviceAsBytes, _ := json.Marshal(services[i])
		serviceLog.Info(serviceAsBytes)
		err := stub.PutState(services[i].ServiceId, serviceAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		serviceLog.Info("Addeds", services[i])
	}
	for i := 0; i < len(agents); i++ {
		serviceLog.Info("i is ", i)
		agentAsBytes, _ := json.Marshal(agents[i])
		err := stub.PutState(agents[i].AgentId, agentAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		serviceLog.Info("Added", agents[i])
	}
	for i := 0; i < len(serviceRelationAgents); i++ {
		serviceLog.Info("i is ", i)
		serviceRelationAgentsAsBytes, _ := json.Marshal(serviceRelationAgents[i])
		err := stub.PutState(serviceRelationAgents[i].RelationId, serviceRelationAgentsAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		serviceLog.Info("Added", serviceRelationAgents[i])
	}
	for i := 0; i < len(reputations); i++ {
		serviceLog.Info("i is ", i)
		reputationsAsBytes, _ := json.Marshal(reputations[i])
		err := stub.PutState(reputations[i].ReputationId, reputationsAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		serviceLog.Info("Added", reputations[i])
	}

	return shim.Success(nil)
}

// =====================================================================================================================
// SaveIndex - save the index
// =====================================================================================================================
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

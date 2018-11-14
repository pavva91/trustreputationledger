#!/usr/bin/env bash
// ==== CHAINCODE RUN (CHAINCODE CONTAINER) ==================

// CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=tcc:0 ./servicemarbles
// CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=trustreputationledger:0 ./trustreputationledger

// ==== IMPORT PACKAGE (CLI) ==================
// go get github.com/hyperledger/fabric/protos/ledger/queryresult

// ==== RUN CLI CONTAINER to DIRECTLY INVOKE the CHAINCODE on the RUNNING NETWORK ==================
// docker exec -it cli bash

// ==== SEE LOGS ON the RUNNING NETWORK ==================
// docker logs <CHAINCODE_CONTAINER_NAME> (E.G.: dev-peer1.org1.example.com-trustreputationledger-1.0 ...)

// ==== CHAINCODE INSTALLATION (CLI) ==================

// peer chaincode install -p chaincodedev/chaincode/trustreputationledger -n trustreputationledger -v 0

// ==== CHAINCODE INSTANTIATION (CLI) ==================

// peer chaincode instantiate -n trustreputationledger -v 0 -c '{"Args":[]}' -C ch2

// ==== CHAINCODE EXECUTION SAMPLES (CLI) ==================

// ==== Invoke servicemarbles ====
// peer chaincode invoke -C ch2 -n trustreputationledger -c '{"function": "HelloWorld", "Args":[]}'
// ==== INITIALIZATION FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n trustreputationledger -c '{"function": "InitLedger", "Args":[]}'

// ==== GENERAL FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "Read", "Args":["idagent1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetValue", "Args":["idagent2"]}' -v 0
// peer chaincode invoke -C ch2 -n scc -c '{"function": "ReadEverything", "Args":[]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "AllStateDB", "Args":[]}'

// ==== CREATE ASSET FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateService", "Args":["idservice5","service1","description1asdfasdf"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateAgent", "Args":["idagent10","agent10","address10"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateServiceAgentRelation", "Args":["idservice1","idagent1","2","6"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateServiceAndServiceAgentRelationWithStandardValue", "Args":["idservice1","service1","description1","idagent1","2","6"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateActivity", "Args":["idagent1","idagent4", "idagent1","idservice1","asdfCIAOsfasdfa","2018-07-23 16:51:01.2","2"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateReputation", "Args":["idagent1","idservice1", "DEMANDER","6"]}'

// ==== MODIFY ASSET FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "ModifyReputationValue", "Args":["idagent1idservice1EXECUTER","8"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "ModifyOrCreateReputationValue", "Args":["idagent1","idservice1","EXECUTER","1.0"]}'


// ==== GET ASSET ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetServiceNotFoundError", "Args":["idservice1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetAgentNotFoundError", "Args":["idagent10"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetServiceRelationAgent", "Args":["idservice1idagent1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetActivity", "Args":["idagent3idagent3idagent3asdfasfasdfa"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetReputationNotFoundError", "Args":["idagent1idservice1EXECUTER"]}'


// ==== GET HISTORY ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetServiceHistory2", "Args":["idagent2"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetHistory", "Args":["idagent1idservice1EXECUTER"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetHistory", "Args":["idagent1idservice1DEMANDER"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetReputationHistory", "Args":["idagent1idservice1EXECUTER"]}'

// ==== RANGE QUERY (USING COMPOSITE INDEX) ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "byService", "Args":["idservice1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "byAgent", "Args":["idAgent10"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetAgentsByService", "Args":["idservice1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "getServicesByAgent", "Args":["idagent1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "byExecutedServiceTxId", "Args":["asdfasfasdfa"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "byDemanderExecuter", "Args":["idagent3","idagent3"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetEvaluationsByServiceTxId", "Args":["asdfasfasdfa"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetActivitiesByDemanderExecuterTimestamp", "Args":["idagent4","idagent1","2018-07-23 16:51:01.2"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "byAgentServiceRole", "Args":["idagent5","idservice4","EXECUTER"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetReputationsByAgentServiceRole", "Args":["idagent5","idservice4","DEMANDER"]}'

// ==== DELETE ASSET ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "DeleteService", "Args":["idservice1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "DeleteAgent", "Args":["idagent1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "DeleteServiceRelationAgent", "Args":["dinnerambassador"]}'




// ==== CALLS IN THE REAL PROJECT ====
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "helloWorld", "Args":[]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "InitLedger", "Args":[]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "AllStateDB", "Args":[]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetHistory", "Args":["half_board"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetHistory", "Args":["S1"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetReputationHistory", "Args":["parc_hotellunchEXECUTER"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetReputationHistory", "Args":["a2S1DEMANDER"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "InitAgent", "Args":["idagent10","agent10","address10"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "InitService", "Args":["idservice10","service10","description10"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetServiceNotFoundError", "Args":["\"asdf\""]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetAgentNotFoundError", "Args":["idagent1"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "ModifyServiceRelationAgentCost", "Args":["breakfastambassador","10"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetServiceRelationAgent", "Args":["breakfastambassador"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "InitServiceAgentRelation", "Args":["idservice1","idagent2","3","5","7"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetAgentsByService", "Args":["CIAO"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetServicesByAgent", "Args":["d"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "byAgent", "Args":["a1"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetServiceNotFoundError", "Args":["idservice5"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "DeleteService", "Args":["half_boardambassador"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "DeleteServiceRelationAgent", "Args":["dambassador"]}'
// peer chaincode invoke -C servicech -n trustreputationledger -c '{"function": "GetActivitiesByDemanderExecuterTimestamp", "Args":["a3","a1","2018-07-23 16:51:01.2"]}'


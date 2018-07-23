// ==== CHAINCODE RUN (CHAINCODE CONTAINER) ==================

// CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=scc:0 ./servicemarbles

// ==== IMPORT PACKAGE (CLI) ==================
// go get github.com/hyperledger/fabric/protos/ledger/queryresult

// ==== CHAINCODE INSTALLATION (CLI) ==================

// peer chaincode install -p chaincodedev/chaincode/servicemarbles -n scc -v 0

// ==== CHAINCODE INSTANTIATION (CLI) ==================

// peer chaincode instantiate -n scc -v 0 -c '{"Args":[]}' -C ch2

// ==== CHAINCODE EXECUTION SAMPLES (CLI) ==================

// ==== Invoke servicemarbles ====
// peer chaincode invoke -C ch2 -n scc -c '{"function": "HelloWorld", "Args":[]}'
// ==== INITIALIZATION FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "InitLedger", "Args":[]}'

// ==== GENERAL FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "Read", "Args":["idagent1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetValue", "Args":["idagent2"]}' -v 0
// peer chaincode invoke -C ch2 -n scc -c '{"function": "ReadEverything", "Args":[]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "AllStateDB", "Args":[]}'

// ==== CREATE ASSET FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateService", "Args":["idservice5","service1","description1asdfasdf"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateAgent", "Args":["idagent10","agent10","address10"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateServiceAgentRelation", "Args":["idservice1","idagent1","2","6"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateServiceAndServiceAgentRelationWithStandardValue", "Args":["idservice10","service10","description10","idagent2","2","6"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateActivity", "Args":["idagent4","idagent4", "idagent1","idservice1","asdfCIAOsfasdfa","2018-07-23 16:51:01.2","2"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "CreateReputation", "Args":["idagent5","idservice4", "DEMANDER","6"]}'

// ==== MODIFY ASSET FUNCTIONS ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "ModifyReputationValue", "Args":["idservice1idagent1EXECUTER","8"]}'


// ==== GET ASSET ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetServiceNotFoundError", "Args":["idservice1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetAgentNotFoundError", "Args":["idagent10"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetServiceRelationAgent", "Args":["idservice1idagent1"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetActivity", "Args":["idagent3idagent3idagent3asdfasfasdfa"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetReputation", "Args":["idagent5idservice4EXECUTER"]}'


// ==== GET HISTORY ==================
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetServiceHistory2", "Args":["idagent2"]}'
// peer chaincode invoke -C ch2 -n scc -c '{"function": "GetHistory", "Args":["idservice1idagent1EXECUTER"]}'

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




// ==== CALLS IN THE REAL PROJECT ====
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "helloWorld", "Args":[]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "InitLedger", "Args":[]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "AllStateDB", "Args":[]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "GetServiceHistory2", "Args":["service5"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "InitAgent", "Args":["idagent10","agent10","address10"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "InitService", "Args":["idservice10","service10","description10"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "GetServiceNotFoundError", "Args":["idservice1"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "GetAgentNotFoundError", "Args":["idagent1"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "GetServiceRelationAgent", "Args":["idservice1idagent1"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "InitServiceAgentRelation", "Args":["idservice1","idagent2","3","5","7"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "GetAgentsByService", "Args":["CIAO"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "byAgent", "Args":["a1"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "GetServiceNotFoundError", "Args":["idservice5"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "DeleteService", "Args":["idservice5"]}'
// peer chaincode invoke -C servicech -n servicemarbles -c '{"function": "GetActivitiesByDemanderExecuterTimestamp", "Args":["a3","a1","2018-07-23 16:51:01.2"]}'


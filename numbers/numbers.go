package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value as parameter")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval))) //making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "add" { //Adds a number to current value
		return t.add(stub, args)
	} else if function == "subtract" { //Subtracts a number from current value
		return t.subtract(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	valAsbytes, err := stub.GetState("abc") //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}

// ============================================================================================================================
// Add - add variable to value in chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var addVal, value int
	var valueAsBytes []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	addVal, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value as parameter")
	}

	//Get the value in state
	valueAsBytes, err = stub.GetState("abc")
	if err != nil {
		return nil, errors.New("Failed to get value from state")
	}

	value, err = strconv.Atoi(string(valueAsBytes))
	if err != nil {
		return nil, errors.New("Expecting integer value as parameter")
	}

	// [CLAUSE]
	if addVal > 0 {

		// [PERFORMANCE]
		value += addVal
		
		err = stub.PutState("abc", []byte(strconv.Itoa(value)))

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// ============================================================================================================================
// Subtract - subtract variable from value in chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) subtract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var subVal, value int
	var valueAsBytes []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	subVal, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value as parameter")
	}

	//Get the value in state
	valueAsBytes, err = stub.GetState("abc")
	if err != nil {
		return nil, errors.New("Failed to get value from state")
	}

	value, err = strconv.Atoi(string(valueAsBytes))
	if err != nil {
		return nil, errors.New("Expecting integer value as parameter")
	}

	// [CLAUSE]
	if subVal > 0 {

		// [PERFORMANCE]
		value -= subVal

		err = stub.PutState("abc", []byte(strconv.Itoa(value)))

		if err != nil {
			return nil, err
		}

	}

	return nil, nil
}

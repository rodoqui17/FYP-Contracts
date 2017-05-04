package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Car struct {
	V5cID 					string `json:"v5cID"`
	Model           string `json:"model"`
	Reg             string `json:"reg"`
	Owner           string `json:"owner"`
	Colour          string `json:"colour"`
	Scrapped        string `json:"scrapped"`
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

	//Args		0
	//				manufacturer

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
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
	} else if function == "create_car" {
		return t.create_car(stub, args)
	}  else if function == "register_car" {
		return t.register_car(stub, args)
	} else if function == "transfer_car" {
		return t.transfer_car(stub, args)
	} else if function == "scrap_car" {
		return t.scrap_car(stub, args)
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
	if function == "get_car" {
		return t.get_car(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}


func (t * SimpleChaincode) get_car(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//var c Car

	bytes, err := stub.GetState(args[0]);
	if err != nil {	return nil, err }

	// err = json.Unmarshal(bytes, &c);
	// if err != nil {	return nil, err }

	return bytes, nil
}

func (t *SimpleChaincode) create_car(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//Args	0						1						2
	//			v5c_ID			model				colour

	var err error
	var c Car
	manufacturer := "BMW"
	caller := "BMW"

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3.")
	}

	v5c_ID         := "\"v5cID\":\""+args[0]+"\", "
	model          := "\"Model\":\""+args[1]+"\", "
	reg            := "\"Reg\":\"UNDEFINED\", "
	owner          := "\"Owner\":\""+manufacturer+"\", "
	colour         := "\"Colour\":\""+args[2]+"\", "
	scrapped       := "\"Scrapped\":\"False\""

	car_json := "{"+v5c_ID+model+reg+owner+colour+scrapped+"}" // Concatenates the variables to create the total JSON object

	fmt.Println("car object is  " + car_json)

	//Convert to a car object
	err = json.Unmarshal([]byte(car_json), &c)
	if err != nil {
		fmt.Println("Unmarshal error")
		return nil, errors.New("Invalid JSON object")
	}

	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Marshal error")
		return nil, errors.New("Error creating car record")
	}

	// [CLAUSE]
	if caller == manufacturer {
		err = stub.PutState(c.V5cID, bytes)
		if err != nil {
			fmt.Println("Put state error")
			return nil, err
		}
	}

	return nil, nil
}

func (t *SimpleChaincode) register_car(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//Args	0						1
	//			v5c_ID			reg

	var err error
	var c Car
	var new_bytes []byte

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}

	bytes, err := stub.GetState(args[0]);
	if err != nil {
		fmt.Println("Get state error")
		return nil, err
	}

	err = json.Unmarshal(bytes, &c);
	if err != nil {
		fmt.Println("Unmarshal error")
		return nil, err
	}

	caller := c.Owner

	fmt.Println("car object is now " + string(new_bytes) + caller)

	// [CLAUSE]
	if caller == c.Owner {

		//[PERFORMANCE]
		c.Reg = args[1]

		new_bytes, err = json.Marshal(c)
		if err != nil {
			fmt.Println("Marshal error")
			return nil, errors.New("Error creating car record")
		}

		fmt.Println("car object is now " + string(new_bytes))

		err = stub.PutState(c.V5cID, new_bytes)
		if err != nil {
			fmt.Println("Put state error")
			return nil, err
		}
	}

	return nil, nil
}

func (t *SimpleChaincode) transfer_car(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//Args	0						1
	//			v5c_ID			new_owner

	var err error
	var c Car
	var new_bytes []byte

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}

	bytes, err := stub.GetState(args[0]);
	if err != nil {
		fmt.Println("Get state error")
		return nil, err
	}

	err = json.Unmarshal(bytes, &c);
	if err != nil {
		fmt.Println("Unmarshal error")
		return nil, err
	}

	caller := c.Owner

	fmt.Println("car object is now " + string(new_bytes))

	// [CLAUSE]
	if caller == c.Owner {

		//[PERFORMANCE]
		c.Owner = args[1]

		new_bytes, err = json.Marshal(c)
		if err != nil {
			fmt.Println("Marshal error")
			return nil, errors.New("Error creating car record")
		}

		fmt.Println("car object is now " + string(new_bytes))

		err = stub.PutState(c.V5cID, new_bytes)
		if err != nil {
			fmt.Println("Put state error")
			return nil, err
		}
	}

	return nil, nil
}

func (t *SimpleChaincode) scrap_car(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//Args	0
	//			v5c_ID

	var err error
	var c Car
	var new_bytes []byte

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}

	bytes, err := stub.GetState(args[0]);
	if err != nil {
		fmt.Println("Get state error")
		return nil, err
	}

	err = json.Unmarshal(bytes, &c);
	if err != nil {
		fmt.Println("Unmarshal error")
		return nil, err
	}

	caller := c.Owner

	// [CLAUSE]
	if caller == c.Owner {

		//[PERFORMANCE]
		c.Scrapped = "True"

		new_bytes, err = json.Marshal(c)
		if err != nil {
			fmt.Println("Marshal error")
			return nil, errors.New("Error creating car record")
		}

		fmt.Println("car object is now " + string(new_bytes))

		err = stub.PutState(c.V5cID, new_bytes)
		if err != nil {
			fmt.Println("Put state error")
			return nil, err
		}
	}

	return nil, nil
}

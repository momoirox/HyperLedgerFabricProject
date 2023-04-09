package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"girhub.com/fist/chaincode/data"
	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// Hello is a simple handler
type Cars struct {
	l        *log.Logger
	contract *gateway.Contract
}

// NewHello creates a new hello handler with the given logger
func NewCars(l *log.Logger, contract *gateway.Contract) *Cars {
	return &Cars{l, contract}
}

func (c *Cars) AddCarMalfunction(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	carId := vars["car"]
	description := vars["description"]
	repairPrice := vars["repairPrice"]

	c.l.Println("Handle AddCarMalfunction")

	result, err := c.contract.SubmitTransaction("AddMalfunction", carId, description, repairPrice)
	if err != nil {
		errors := strings.Split(err.Error(), ":")
		message := fmt.Sprintf("Failed to evaluate transaction: %s\n", errors[len(errors)-1])
		fmt.Printf(message)
		http.Error(rw, message, http.StatusConflict)
		return
	}
	fmt.Println(string(result))

}

func (c *Cars) RepairCar(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	carId := vars["car"]

	c.l.Println("Handle repairCar")

	result, err := c.contract.SubmitTransaction("RepairCar", carId)
	if err != nil {
		errors := strings.Split(err.Error(), ":")
		message := fmt.Sprintf("Failed to evaluate transaction: %s\n", errors[len(errors)-1])
		fmt.Printf(message)
		http.Error(rw, message, http.StatusConflict)
		return
	}
	fmt.Println(string(result))

}

func (c *Cars) ChangeCarColor(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	carId := vars["car"]
	newColour := vars["color"]

	c.l.Println("Handle changeCarColor")

	result, err := c.contract.SubmitTransaction("ChangeCarColour", carId, newColour)
	if err != nil {
		errors := strings.Split(err.Error(), ":")
		message := fmt.Sprintf("Failed to evaluate transaction: %s\n", errors[len(errors)-1])
		fmt.Printf(message)
		http.Error(rw, message, http.StatusConflict)
		return
	}
	fmt.Println(string(result))

}

func (c *Cars) TransferCarOwnership(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	carId := vars["car"]
	newOwnerId := vars["owner"]
	flag := vars["flag"]
	acceptMalfunctionedBool := false

	switch flag {
	case "yes":
		acceptMalfunctionedBool = true
	case "no":
		acceptMalfunctionedBool = false
	default:
		acceptMalfunctionedBool = false
	}

	c.l.Println("Handle transferCarOwnership")

	result, err := c.contract.SubmitTransaction("ChangeOwner", carId, newOwnerId, fmt.Sprintf("%t", acceptMalfunctionedBool))
	if err != nil {
		errors := strings.Split(err.Error(), ":")
		message := fmt.Sprintf("Failed to evaluate transaction: %s\n", errors[len(errors)-1])
		fmt.Printf(message)
		http.Error(rw, message, http.StatusConflict)
		return
	}
	fmt.Println(string(result))

}

func (c *Cars) GetCarsByColorAndOwner(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	color := vars["color"]
	ownerId := vars["owner"]

	c.l.Println("Handle GET car by color & owner")

	result, err := c.contract.EvaluateTransaction("QueryCarsByColorAndOwner", color, ownerId)
	if err != nil {
		errors := strings.Split(err.Error(), ":")
		message := fmt.Sprintf("Failed to evaluate transaction: %s\n", errors[len(errors)-1])
		fmt.Printf(message)
		http.Error(rw, message, http.StatusConflict)
		return
	}
	fmt.Println(string(result))

	cars := []data.Car{}

	err = json.Unmarshal(result, &cars)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

	for _, car := range cars {
		car.ToJSON(rw)
	}

}

func (c *Cars) GetCarsByColor(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	color := vars["color"]

	c.l.Println("Handle GET car by color")

	result, err := c.contract.EvaluateTransaction("QueryCarsByColor", color)
	if err != nil {
		errors := strings.Split(err.Error(), ":")
		message := fmt.Sprintf("Failed to evaluate transaction: %s\n", errors[len(errors)-1])
		fmt.Printf(message)
		http.Error(rw, message, http.StatusConflict)

	}
	fmt.Println(string(result))

	cars := []data.Car{}

	err = json.Unmarshal(result, &cars)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

	for _, car := range cars {
		car.ToJSON(rw)
	}
}

func (c *Cars) GetPerson(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	personId := vars["id"]

	c.l.Println("Handle GET Person")

	result, err := c.contract.EvaluateTransaction("QueryPerson", personId)
	if err != nil {
		errors := strings.Split(err.Error(), ":")
		message := fmt.Sprintf("Failed to evaluate transaction: %s\n", errors[len(errors)-1])
		fmt.Printf(message)
		http.Error(rw, message, http.StatusConflict)
		return
	}
	fmt.Println(string(result))

	person := data.Person{}

	err = json.Unmarshal(result, &person)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

	person.ToJSON(rw)
}

func (c *Cars) GetCar(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	carId := vars["id"]

	c.l.Println("Handle GET Cars")

	result, err := c.contract.EvaluateTransaction("QueryCar", carId)
	if err != nil {
		errors := strings.Split(err.Error(), ":")
		message := fmt.Sprintf("Failed to evaluate transaction: %s\n", errors[len(errors)-1])
		fmt.Printf(message)
		http.Error(rw, message, http.StatusConflict)
		return
	}
	fmt.Println(string(result))

	car := data.Car{}
	err = json.Unmarshal(result, &car)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

	car.ToJSON(rw)
}

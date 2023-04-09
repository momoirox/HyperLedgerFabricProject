package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type CarMalfunction struct {
	Description string
	RepairPrice float32
}

type Car struct {
	Id              string
	Brand           string
	Model           string
	Year            int
	Colour          string
	OwnerId         string
	Price           float32
	MalfunctionList []CarMalfunction
}

type Person struct {
	Id      string
	Name    string
	Surname string
	Email   string
	Money   float32
}

type QueryResult struct {
	Key    string `json:"Key"`
	Record *Car
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Car{
		{Id: "car1", Brand: "Toyota", Year: 2001, Model: "Prius", Colour: "blue", OwnerId: "person1", Price: 100.00, MalfunctionList: []CarMalfunction{
			{Description: "Broken Tail/Head Lights", RepairPrice: 40},
			{Description: "Warning Lights", RepairPrice: 50},
		}},
		{Id: "car2", Brand: "Ford", Year: 2001, Model: "Mustang", Colour: "red", OwnerId: "person1", Price: 200.00, MalfunctionList: []CarMalfunction{
			{Description: "Bad Fuel Economy", RepairPrice: 40},
		}},
		{Id: "car3", Brand: "Fiat", Year: 2001, Model: "XXL", Colour: "pink", OwnerId: "person1", Price: 300.00, MalfunctionList: []CarMalfunction{
			{Description: "Flat Tires", RepairPrice: 50},
		}},
		{Id: "car4", Brand: "Hyundai", Year: 2001, Model: "Tucson", Colour: "green", OwnerId: "person2", Price: 400.00, MalfunctionList: []CarMalfunction{
			{Description: "Rusting", RepairPrice: 100},
		}},
		{Id: "car5", Brand: "Volkswagen", Year: 2001, Model: "Passat", Colour: "yellow", OwnerId: "person3", Price: 500.00, MalfunctionList: []CarMalfunction{
			{Description: "Bad Brakes", RepairPrice: 10},
			{Description: "Overheating", RepairPrice: 15},
		}},
		{Id: "car6", Brand: "Tesla", Year: 2001, Model: "S", Colour: "black", OwnerId: "person3", Price: 600.00, MalfunctionList: []CarMalfunction{
			{Description: "Airbags That Injure", RepairPrice: 20},
		}},
	}
	persons := []Person{
		{Id: "person1", Name: "Jean-Jacques", Surname: "Rousseau", Email: "rousseau@gmail.com", Money: 8900.99},
		{Id: "person2", Name: "Marco", Surname: "Polo", Email: "polo@gmail.com", Money: 3230.33},
		{Id: "person3", Name: "Amadeo", Surname: "Avogadro", Email: "avogadro@gmail.com", Money: 3333.33},
	}

	for _, car := range cars {
		carAsBytes, _ := json.Marshal(car)
		err := ctx.GetStub().PutState(car.Id, carAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

		//  ==== Index the marble to enable color-based range queries, e.g. return all blue marbles ====
		//  An 'index' is a normal key/value entry in state.
		//  The key is a composite key, with the elements that you want to range query on listed first.
		//  In our case, the composite key is based on indexName=color~name.
		//  This will enable very efficient state range queries based on composite keys matching indexName=color~*
		indexName := "Colour~OwnerId~Id"
		colorOwnerIndexKey, err := ctx.GetStub().CreateCompositeKey(indexName, []string{car.Colour, car.OwnerId, car.Id})
		if err != nil {
			return err
		}

		value := []byte{0x00}
		err = ctx.GetStub().PutState(colorOwnerIndexKey, value)
		if err != nil {
			return err
		}
	}

	for _, person := range persons {
		personAsBytes, err := json.Marshal(person)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(person.Id, personAsBytes)
		if err != nil {
			return fmt.Errorf("Failed to put persons to world state. %v", err)
		}
	}

	return nil
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, carNumber string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carNumber)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}

//QueryPerson returns the person stored in the world state with given id
func (s *SmartContract) QueryPerson(ctx contractapi.TransactionContextInterface, personId string) (*Person, error) {
	personAsBytes, err := ctx.GetStub().GetState(personId)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if personAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", personId)
	}

	person := new(Person)
	_ = json.Unmarshal(personAsBytes, person)

	return person, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]*Car, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	retList := []*Car{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var car *Car
		err = json.Unmarshal(response.Value, &car)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		retList = append(retList, car)
	}

	return retList, nil
}

func (s *SmartContract) QueryCarsByColor(ctx contractapi.TransactionContextInterface, color string) ([]*Car, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Colour~OwnerId~Id", []string{color})
	if err != nil {
		return nil, err
	}

	defer iterator.Close()

	retList := make([]*Car, 0)

	for i := 0; iterator.HasNext(); i++ {
		responseRange, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, err
		}

		retCarId := compositeKeyParts[2]

		car, err := s.QueryCar(ctx, retCarId)
		if err != nil {
			return nil, err
		}

		retList = append(retList, car)
	}

	return retList, nil

}

func (s *SmartContract) QueryCarsByOwner(ctx contractapi.TransactionContextInterface, OwnerId string) ([]*Car, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Colour~OwnerId~Id", []string{OwnerId})
	if err != nil {
		return nil, err
	}

	defer iterator.Close()

	retList := make([]*Car, 0)

	for i := 0; iterator.HasNext(); i++ {
		responseRange, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, err
		}

		carId := compositeKeyParts[2]

		car, err := s.QueryCar(ctx, carId)
		if err != nil {
			return nil, err
		}

		retList = append(retList, car)
	}

	return retList, nil
}

func (s *SmartContract) QueryCarsByColorAndOwner(ctx contractapi.TransactionContextInterface, Colour string, OwnerId string) ([]*Car, error) {

	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Colour~OwnerId~Id", []string{Colour, OwnerId})
	if err != nil {
		return nil, err
	}

	defer iterator.Close()

	retList := make([]*Car, 0)

	for i := 0; iterator.HasNext(); i++ {
		responseRange, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, err
		}

		retCarId := compositeKeyParts[2]

		car, err := s.QueryCar(ctx, retCarId)
		if err != nil {
			return nil, err
		}

		retList = append(retList, car)
	}

	return retList, nil
}

func (s *SmartContract) ChangeOwner(ctx contractapi.TransactionContextInterface, carId string, newOwnerId string, acceptCarWithMalfunction bool) error {
	car, err := s.QueryCar(ctx, carId)
	if err != nil {
		return err
	}

	if car.OwnerId == newOwnerId {
		return fmt.Errorf("This person already owns this car!")
	}

	price := car.Price

	newOwner, err := s.QueryPerson(ctx, newOwnerId)
	if err != nil {
		return err
	}

	oldOwner, err := s.QueryPerson(ctx, car.OwnerId)
	if err != nil {
		return err
	}

	if !acceptCarWithMalfunction && len(car.MalfunctionList) > 0 {
		return fmt.Errorf("This car has malfunctions, purchase cannot be made! ")
	}
	if acceptCarWithMalfunction && len(car.MalfunctionList) > 0 {
		for _, malfunction := range car.MalfunctionList {
			price -= malfunction.RepairPrice
		}
	}
	if newOwner.Money < price {
		return fmt.Errorf("The buyer doesn't have enough money to buy the car! ")
	}

	newOwner.Money -= price
	oldOwner.Money += price

	car.OwnerId = newOwnerId

	carAsBytes, err := json.Marshal(car)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(carId, carAsBytes)
	if err != nil {
		return err
	}

	indexName := "Colour~OwnerId~Id"
	colorOwnerIdIndexKey, err := ctx.GetStub().CreateCompositeKey(indexName, []string{car.Colour, car.OwnerId, car.Id})
	if err != nil {
		return err
	}
	value := []byte{0x00}
	err = ctx.GetStub().PutState(colorOwnerIdIndexKey, value)
	if err != nil {
		return err
	}
	oldColorOwnerIDIndexKey, err := ctx.GetStub().CreateCompositeKey(indexName, []string{car.Colour, oldOwner.Id, car.Id})
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(oldColorOwnerIDIndexKey)
	if err != nil {
		return err
	}

	oldOwnerAsBytes, err := json.Marshal(oldOwner)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(oldOwner.Id, oldOwnerAsBytes)
	if err != nil {
		return err
	}

	newOwnerAsBytes, err := json.Marshal(newOwner)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(newOwner.Id, newOwnerAsBytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) ChangeCarColour(ctx contractapi.TransactionContextInterface, carNumber string, newColour string) error {
	car, err := s.QueryCar(ctx, carNumber)

	if err != nil {
		return err
	}

	oldColour := car.Colour
	car.Colour = newColour

	carAsBytes, _ := json.Marshal(car)

	err = ctx.GetStub().PutState(carNumber, carAsBytes)

	if err != nil {
		return err
	}

	//Must change the entry
	indexName := "Colour~OwnerId~Id"
	newColorOwnerIndexKey, err := ctx.GetStub().CreateCompositeKey(indexName, []string{newColour, car.OwnerId, car.Id})
	if err != nil {
		return err
	}

	value := []byte{0x00}
	err = ctx.GetStub().PutState(newColorOwnerIndexKey, value)
	if err != nil {
		return err
	}

	oldColorOwnerIndexKey, err := ctx.GetStub().CreateCompositeKey(indexName, []string{oldColour, car.OwnerId, car.Id})
	if err != nil {
		return err
	}

	err = ctx.GetStub().DelState(oldColorOwnerIndexKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) AddMalfunction(ctx contractapi.TransactionContextInterface, carId string, description string, price float32) error {
	car, err := s.QueryCar(ctx, carId)
	if err != nil {
		return err
	}

	newMalfunction := CarMalfunction{
		Description: description,
		RepairPrice: price,
	}

	car.MalfunctionList = append(car.MalfunctionList, newMalfunction)

	totalMalfunctionsPrice := float32(0)
	for _, malfunction := range car.MalfunctionList {
		totalMalfunctionsPrice += malfunction.RepairPrice
	}

	if totalMalfunctionsPrice > car.Price {
		err = ctx.GetStub().DelState(carId)
		if err != nil {
			return err
		}
	} else {
		carAsBytes, _ := json.Marshal(car)
		err = ctx.GetStub().PutState(car.Id, carAsBytes)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}
	return nil
}

func (s *SmartContract) RepairCar(ctx contractapi.TransactionContextInterface, carId string) error {
	car, err := s.QueryCar(ctx, carId)
	if err != nil {
		return err
	}
	owner, err := s.QueryPerson(ctx, car.OwnerId)
	if err != nil {
		return err
	}

	price := float32(0)
	for _, malfuction := range car.MalfunctionList {
		price += malfuction.RepairPrice
	}
	if owner.Money < price {
		return fmt.Errorf("The owner has no enough money to repair the car.")
	}

	owner.Money -= price
	car.MalfunctionList = []CarMalfunction{}

	carAsBytes, err := json.Marshal(car)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(carId, carAsBytes)
	if err != nil {
		return err
	}

	ownerAsBytes, err := json.Marshal(owner)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(owner.Id, ownerAsBytes)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Blockchain is a public database distributed between peers
	// Db doesnt rely on trusting the nodes
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}

}

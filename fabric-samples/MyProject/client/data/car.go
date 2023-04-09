package data

import (
	"encoding/json"
	"io"
)

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

func (p *Car) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Person) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

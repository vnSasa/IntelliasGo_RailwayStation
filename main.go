package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	etalonTime = "15:04:05"
	fileJSON = "data.json"
)

type Trains []Train

type Train struct {
	TrainID            int		`json:"trainId"` 
	DepartureStationID int		`json:"departureStationID"`
	ArrivalStationID   int		`json:"arrivalStationId"`
	Price              float32	`json:"price"`
	ArrivalTime        time.Time	`json:"arrivalTime"`
	DepartureTime      time.Time	`json:"departureTime"`
}

func readDataJSON() (data string) {
	jsonFile, err := os.ReadFile(fileJSON)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func (t *Train) UnmarshalJSON(data []byte)error {
	var byteValue map[string]interface{}
	err := json.Unmarshal(data, &byteValue)
	if err != nil {
		return err
	}

	id, _ := byteValue["trainId"].(float64)
	t.TrainID = int(id)

	depStationId, _ := byteValue["departureStationID"].(float64)
	t.DepartureStationID = int(depStationId)

	arrStationId, _ := byteValue["arrivalStationId"].(float64)
	t.ArrivalStationID = int(arrStationId)

	price, _ := byteValue["price"].(float64)
	t.Price = float32(price)

	arrTime, _ := byteValue["arrivalTime"].(string)
	t.ArrivalTime, _ = time.Parse(etalonTime, arrTime)

	depTime, _ := byteValue["departureTime"].(string)
	t.DepartureTime, _ = time.Parse(etalonTime, depTime)

	return nil
}

func userRequest() (input string) {
	fmt.Scanln(& input)
	return input
}

func main() {
	//	... запит даних від користувача
	fmt.Println("enter departure station id: ")
	depStationId := userRequest()
	fmt.Println("enter arrival station id: ")
	arrStationId := userRequest()
	fmt.Println("enter criteria: ")
	criteria := userRequest()

	fmt.Println(depStationId, arrStationId, criteria)
	//result, err := FindTrains(departureStation, arrivalStation, criteria))
	//	... обробка помилки
	//	... друк result
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// ... код
	return nil, nil // маєте повернути правильні значення
}

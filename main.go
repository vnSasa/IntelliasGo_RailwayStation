package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"errors"
	"strconv"
	"strings"
)

const (
	etalonTime = "15:04:05"
	fileJSON = "data.json"
)

var (
	UnsupportedCriteria			= errors.New("unsupported criteria")
	EmptyDepartureStation		= errors.New("empty departure station")
	EmptyArrivalStation			= errors.New("empty arrival station")
	BadArrivalStationInput		= errors.New("bad arrival station input")
	BadDepartureStationInput	= errors.New("bad departure station input")
)

type Trains []Train

type Train struct {
	TrainID					int		`json:"trainId"` 
	DepartureStationID		int		`json:"departureStationId"`
	ArrivalStationID		int		`json:"arrivalStationId"`
	Price					float32		`json:"price"`
	ArrivalTime				time.Time		`json:"arrivalTime"`
	DepartureTime			time.Time		`json:"departureTime"`
}

func readDataJSON() (data Trains) {
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

	depStationId, _ := byteValue["departureStationId"].(float64)
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
	departureStation := userRequest()
	fmt.Println("enter arrival station id: ")
	arrivalStation := userRequest()
	fmt.Println("enter criteria: ")
	criteria := userRequest()

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	//	... обробка помилки
	if err != nil {
		fmt.Printf("invalid data entered - %v", err)
	}
	//	... друк result
	for _, v := range result {
		fmt.Println(v)
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// ... код
	var trains Trains
	trains = readDataJSON()
	trains = filteredTrains(trains, departureStation, arrivalStation)
	return trains, nil // маєте повернути правильні значення
}

func filteredTrains(trains Trains, departureStation string, arrivalStation string) (validTrains Trains) {
	for _, v := range trains {
		depStation := strconv.Itoa(v.DepartureStationID)
		arrStation := strconv.Itoa(v.ArrivalStationID)
		if strings.Compare(depStation, departureStation) == 0 && strings.Compare(arrStation, arrivalStation) == 0 {
			validTrains = append(validTrains, v)
		}
	}
	return validTrains
}
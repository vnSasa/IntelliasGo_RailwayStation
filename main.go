package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"errors"
	"strconv"
	"strings"
	"sort"
)

const (
	etalonTime = "15:04:05"
	fileJSON = "data.json"
)

var (
	UnsupportedCriteria		= errors.New("unsupported criteria")
	EmptyDepartureStation		= errors.New("empty departure station")
	EmptyArrivalStation		= errors.New("empty arrival station")
	BadArrivalStationInput		= errors.New("bad arrival station input")
	BadDepartureStationInput	= errors.New("bad departure station input")
)

type Trains []Train

type Train struct {
	TrainID				int		`json:"trainId"` 
	DepartureStationID		int		`json:"departureStationId"`
	ArrivalStationID		int		`json:"arrivalStationId"`
	Price				float32		`json:"price"`
	ArrivalTime			time.Time	`json:"arrivalTime"`
	DepartureTime			time.Time	`json:"departureTime"`
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

func outputDataJSON(train Train) {
	fmt.Printf("TrainID: %v, DepartureStationID: %v, ArrivalStationID: %v, "+
				"Price: %v, ArrivalTime: %v, DepartureTime: %v\n",	
				train.TrainID, train.DepartureStationID, train.ArrivalStationID, train.Price,
				train.ArrivalTime.Format("15:04:05"), train.DepartureTime.Format("15:04:05"))
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
	criteria := strings.ToLower(userRequest())

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	//	... обробка помилки
	if err != nil {
		fmt.Printf("invalid data entered - %v", err)
		return
	}
	if result == nil {
		fmt.Println("no trains were found according to the specified data")
	}
	//	... друк result
	for _, v := range result {
		outputDataJSON(v)
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// ... код
	var trains Trains
	trains = nil
	err := inputValidation(trains, departureStation, arrivalStation, criteria)
	if err == nil {
		trains = readDataJSON()
		trains = filteredTrains(trains, departureStation, arrivalStation)
		trains = filteredByCriteria(trains, criteria)
	}
	
	return trains, err // маєте повернути правильні значення
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

func filteredByCriteria(trains Trains, criteria string) Trains {
	switch criteria {
	case "price" :
		sort.SliceStable(trains, func(i, j int)bool {
			return trains[i].Price < trains[j].Price
		})
	case "arrival-time" :
		sort.SliceStable(trains, func(i, j int)bool {
			return trains[i].ArrivalTime.Before(trains[j].ArrivalTime)
		})
	case "departure-time" :
		sort.SliceStable(trains, func(i, j int)bool {
			return trains[i].DepartureTime.Before(trains[j].DepartureTime)
		})
	}
	return limitedOutput(trains)
}

func limitedOutput(trains Trains) (outputedTrains Trains) {
	for _, v := range trains {
		outputedTrains = append(outputedTrains, v)
		if len(outputedTrains) == 3 {
			return outputedTrains
		}
	}
	return
}

func inputValidation(trains Trains, departureStation, arrivalStation, criteria string) (error) {
	// empty input
	if strings.Compare(departureStation, "") == 0 {
		return EmptyDepartureStation
	}
	if strings.Compare(arrivalStation, "") == 0 {
		return EmptyArrivalStation
	}
	// bad input
	depStation, err := strconv.Atoi(departureStation)
	if err != nil || depStation < 1 {
		return BadDepartureStationInput
	}
	arrStation, err := strconv.Atoi(arrivalStation)
	if err != nil || arrStation < 1 {
		return BadArrivalStationInput
	}
	// unsupported criteria input
	err = UnsupportedCriteria
	if strings.Compare(criteria, "price") == 0 || strings.Compare(criteria, "arrival-time") == 0 || strings.Compare(criteria, "departure-time") == 0 {
		err = nil
	}

	return err
}
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
	unsupportedCriteria	= errors.New("unsupported criteria")
	emptyDepartureStation	= errors.New("empty departure station")
	emptyArrivalStation	= errors.New("empty arrival station")
	badArrivalStationInput	= errors.New("bad arrival station input")
	badDepartureStationInput	= errors.New("bad departure station input")
)

type Trains []Train

type Train struct {
	TrainID	int	`json:"trainId"` 
	DepartureStationID	int	`json:"departureStationId"`
	ArrivalStationID	int	`json:"arrivalStationId"`
	Price	float32	`json:"price"`
	ArrivalTime	time.Time	`json:"arrivalTime"`
	DepartureTime	time.Time	`json:"departureTime"`
}

func readDataJSON() (data Trains, err error) {
	jsonFile, err := os.ReadFile(fileJSON)
	if err != nil { return nil, err }
	
	err = json.Unmarshal(jsonFile, &data)
	if err != nil { return nil, err }
	
	return data, nil
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
	if err != nil { return err }

	if id, ok := byteValue["trainId"].(float64); ok {
		t.TrainID = int(id)
	}

	if depStationId, ok := byteValue["departureStationId"].(float64); ok {
		t.DepartureStationID = int(depStationId)
	}
	
	if arrStationId, ok := byteValue["arrivalStationId"].(float64); ok {
		t.ArrivalStationID = int(arrStationId)
	}

	if price, ok := byteValue["price"].(float64); ok {
		t.Price = float32(price)
	}

	if arrTime, ok := byteValue["arrivalTime"].(string); ok {
		t.ArrivalTime, _ = time.Parse(etalonTime, arrTime)
	}
	
	if depTime, ok := byteValue["departureTime"].(string); ok {
		t.DepartureTime, _ = time.Parse(etalonTime, depTime)
	}

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

	//	... обробка помилки
	
	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		fmt.Printf("invalid data entered - %v", err)
		return
	}
	
	if len(result) == 0 {
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
	
	trains, err := readDataJSON()
	if err != nil {
		return nil, err
	}

	err = inputValidation(trains, departureStation, arrivalStation, criteria)
	if err != nil {
		return nil, err
	}
	
	trains = filteredByCriteria(filteredTrains(trains, departureStation, arrivalStation), criteria)

	return trains, err // маєте повернути правильні значення
}

func filteredTrains(trains Trains, departureStation, arrivalStation string) (validTrains Trains) {
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
	if len(departureStation) == 0 {
		return emptyDepartureStation
	}
	
	if len(arrivalStation) == 0 {
		return emptyArrivalStation
	}
	
	// bad input
	depStation, err := strconv.Atoi(departureStation)
	if err != nil || depStation < 1 {
		return badDepartureStationInput
	}
	
	arrStation, err := strconv.Atoi(arrivalStation)
	if err != nil || arrStation < 1 {
		return badArrivalStationInput
	}
	
	// unsupported criteria input
	switch criteria {
		case "price" : return nil
		case "arrival-time" : return nil
		case "departure-time" : return nil
		default : err = unsupportedCriteria
	}

	return err
}
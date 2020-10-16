package logic

import (
	"encoding/json"
	"github.com/jinzhu/now"
	"strconv"
	"time"
)

// inputLoad represents a inputLoad json input
type inputLoad struct {
	LoadId     string     `json:"id"`
	CustomerId string     `json:"customer_id"`
	Amount     loadAmount `json:"load_amount"`
	Time       time.Time  `json:"time"`
}

// loadResponse response given to a load
type loadResponse struct {
	LoadId     string `json:"id"`
	CustomerId string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
}

// customerLoadId couple load / customer
type customerLoadId struct {
	LoadId     string
	CustomerId string
}

// loadAmount represents the amount value of a load
type loadAmount struct {
	Value float64
}

// UnmarshalJSON implementation of parsing of $123.45 to a loadAmount
func (l *loadAmount) UnmarshalJSON(b []byte) error {
	amountStr := string(b)
	numberAmountStr := amountStr[2 : len(amountStr)-1]
	amount, err := strconv.ParseFloat(numberAmountStr, 64)
	if err != nil {
		return err
	}
	l.Value = amount
	return nil
}

type FinanceLogic struct {
	CustomersLoads map[string][]inputLoad
	TreatedLoadIds map[customerLoadId]interface{}
}

type LoadParser interface {
	ParseLoads(parsingChannel chan string) ([]string, []error)
}

func NewFinanceLogic() LoadParser {
	return &FinanceLogic{
		CustomersLoads: make(map[string][]inputLoad),
		TreatedLoadIds: make(map[customerLoadId]interface{}),
	}
}

func (logic *FinanceLogic) validateLoadAndFill(load inputLoad) bool {
	customerLoads, customerExist := logic.CustomersLoads[load.CustomerId]
	if !customerExist {
		customerLoads = make([]inputLoad, 0)
	}
	validated := logic.validateLoad(load, customerLoads)
	if validated {
		logic.CustomersLoads[load.CustomerId] = append(customerLoads, load)
	}
	return validated

}

func (logic *FinanceLogic) validateLoad(load inputLoad, customerLoads []inputLoad) bool {
	dayStart := now.With(load.Time).BeginningOfDay().Add(-time.Second)
	dayEnd := now.With(load.Time).EndOfDay()
	weekEnd := now.With(load.Time).EndOfWeek()
	weekStart := now.With(load.Time).BeginningOfWeek()
	dayAmountSum := float64(0)
	dayAmountCount := 0
	weekAmountSum := float64(0)
	for _, storedLoad := range customerLoads {
		if storedLoad.Time.After(dayStart) && storedLoad.Time.Before(dayEnd) {
			dayAmountCount++
			dayAmountSum += storedLoad.Amount.Value
		}
		if storedLoad.Time.After(weekStart) && storedLoad.Time.Before(weekEnd) {
			weekAmountSum += storedLoad.Amount.Value
		}
	}
	if dayAmountSum+load.Amount.Value > 5000 {
		return false
	}
	if dayAmountCount > 2 {
		return false
	}
	if weekAmountSum+load.Amount.Value > 20000 {
		return false
	}
	return true
}

func (logic *FinanceLogic) ParseLoads(parsingChannel chan string) ([]string, []error) {
	loadResponses := make([]string, 0)
	loadErrors := make([]error, 0)
	for line := range parsingChannel {
		var loadTry inputLoad
		err := json.Unmarshal([]byte(line), &loadTry)
		if err != nil {
			loadErrors = append(loadErrors, err)
		}
		if !logic.addLoadToTreated(loadTry) {
			loadStatus := logic.validateLoadAndFill(loadTry)
			loadResponse := loadResponse{
				LoadId:     loadTry.LoadId,
				CustomerId: loadTry.CustomerId,
				Accepted:   loadStatus,
			}
			loadResponseString, err := json.Marshal(loadResponse)
			if err != nil {
				loadErrors = append(loadErrors, err)
			} else {
				loadResponses = append(loadResponses, string(loadResponseString))
			}
		}

	}
	return loadResponses, loadErrors
}

// addLoadToTreated adds load to the list of treated ones and returns true if already exists
func (logic *FinanceLogic) addLoadToTreated(load inputLoad) bool {
	customerLoadId := customerLoadId{
		LoadId:     load.LoadId,
		CustomerId: load.CustomerId,
	}
	_, loadExist := logic.TreatedLoadIds[customerLoadId]
	if loadExist {
		return true
	}
	logic.TreatedLoadIds[customerLoadId] = nil
	return false
}

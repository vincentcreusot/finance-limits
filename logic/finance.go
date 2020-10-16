package logic

import (
	"encoding/json"
	"github.com/jinzhu/now"
	"strconv"
	"time"
)

// Load represents a Load json input
type Load struct {
	LoadId     string     `json:"id"`
	CustomerId string     `json:"customer_id"`
	Amount     LoadAmount `json:"load_amount"`
	Time       time.Time  `json:"time"`
}

// LoadResponse response given to a load
type LoadResponse struct {
	LoadId     string `json:"id"`
	CustomerId string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
}

type CustomerLoadId struct {
	LoadId     string
	CustomerId string
}
type LoadAmount struct {
	Value float64
}

func (l *LoadAmount) UnmarshalJSON(b []byte) error {
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
	CustomersLoads map[string][]Load
	TreatedLoadIds map[CustomerLoadId]interface{}
}

type LoadParser interface {
	ParseLoads(parsingChannel chan string) ([]string, []error)
}

func NewFinanceLogic() LoadParser {
	return &FinanceLogic{
		CustomersLoads: make(map[string][]Load),
		TreatedLoadIds: make(map[CustomerLoadId]interface{}),
	}
}

func (logic *FinanceLogic) validateLoadAndFill(load Load) bool {
	customerLoads, customerExist := logic.CustomersLoads[load.CustomerId]
	if !customerExist {
		customerLoads = make([]Load, 0)
	}
	validated := logic.validateLoad(load, customerLoads)
	if validated {
		logic.CustomersLoads[load.CustomerId] = append(customerLoads, load)
	}
	return validated

}

func (logic *FinanceLogic) validateLoad(load Load, customerLoads []Load) bool {
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
		var loadTry Load
		err := json.Unmarshal([]byte(line), &loadTry)
		if err != nil {
			loadErrors = append(loadErrors, err)
		}
		if !logic.addLoadToTreated(loadTry) {
			loadStatus := logic.validateLoadAndFill(loadTry)
			loadResponse := LoadResponse{
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
func (logic *FinanceLogic) addLoadToTreated(load Load) bool {
	customerLoadId := CustomerLoadId{
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

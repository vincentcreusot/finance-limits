package logic

import (
	"encoding/json"
	"github.com/jinzhu/now"
	"strconv"
	"time"
)

const (
	dayMaxAmount  = 5000
	dayMaxCount   = 3
	weekMaxAmount = 20000
)

// inputLoad represents a inputLoad json input
type inputLoad struct {
	LoadID     string     `json:"id"`
	CustomerID string     `json:"customer_id"`
	Amount     loadAmount `json:"load_amount"`
	Time       time.Time  `json:"time"`
}

// loadResponse response given to a load
type loadResponse struct {
	LoadID     string `json:"id"`
	CustomerID string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
}

// customerLoadID couple load / customer
type customerLoadID struct {
	LoadID     string
	CustomerID string
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

// FinanceLogic LoadParser implementation for holding history maps
type FinanceLogic struct {
	CustomersLoads map[string][]inputLoad
	TreatedLoadIds map[customerLoadID]interface{}
}

// LoadParser interface for defining how to parse loads
type LoadParser interface {
	ParseLoads(parsingChannel chan string) ([]string, []error)
}

// NewFinanceLogic creates a LoadParser implementation
func NewFinanceLogic() *FinanceLogic {
	return &FinanceLogic{
		CustomersLoads: make(map[string][]inputLoad),
		TreatedLoadIds: make(map[customerLoadID]interface{}),
	}
}

// validateLoadAndFillHistory deals with load history for each customer and validate
func (logic *FinanceLogic) validateLoadAndFillHistory(load inputLoad) bool {
	customerLoads, customerExist := logic.CustomersLoads[load.CustomerID]
	if !customerExist {
		customerLoads = make([]inputLoad, 0)
	}
	validated := validateLoad(load, customerLoads)
	if validated {
		logic.CustomersLoads[load.CustomerID] = append(customerLoads, load)
	}
	return validated
}

// validateLoad validates a load using load history given as parameter
func validateLoad(load inputLoad, customerLoads []inputLoad) bool {
	dayStart := now.With(load.Time).BeginningOfDay().Add(-time.Second) // removing one second for comparison
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

	if dayAmountSum+load.Amount.Value > dayMaxAmount {
		return false
	}
	if dayAmountCount >= dayMaxCount {
		return false
	}
	if weekAmountSum+load.Amount.Value > weekMaxAmount {
		return false
	}
	return true
}

// ParseLoads parse the loads given in a channel
func (logic *FinanceLogic) ParseLoads(parsingChannel chan string) ([]string, []error) {
	loadResponses := make([]string, 0)
	loadErrors := make([]error, 0)
	for line := range parsingChannel {
		var loadTry inputLoad
		err := json.Unmarshal([]byte(line), &loadTry)
		if err != nil {
			loadErrors = append(loadErrors, err)
		} else {
			if logic.addCustomerLoadToTreated(loadTry) { // do not treat if (loadid, customerid)  couple already exists
				loadStatus := logic.validateLoadAndFillHistory(loadTry)
				loadResponse := loadResponse{
					LoadID:     loadTry.LoadID,
					CustomerID: loadTry.CustomerID,
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

	}
	return loadResponses, loadErrors
}

// addCustomerLoadToTreated adds load to the list of treated ones and returns false if not added (already exists)
func (logic *FinanceLogic) addCustomerLoadToTreated(load inputLoad) bool {
	customerLoadID := customerLoadID{
		LoadID:     load.LoadID,
		CustomerID: load.CustomerID,
	}
	_, loadExist := logic.TreatedLoadIds[customerLoadID]
	if loadExist {
		return false
	}
	logic.TreatedLoadIds[customerLoadID] = nil
	return true
}

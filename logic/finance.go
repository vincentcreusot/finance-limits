package logic

import "time"

type Load struct {
	LoadId string `json:"id"`
	CustomerId string `json:"customer_id"`
	Amount float32 `json:"load_amount"`
	Time time.Time `json:"time"`
}

func validateLoad(loadLine string) {

}

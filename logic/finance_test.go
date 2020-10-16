package logic

import (
	"encoding/json"
	"testing"
)

func Test_UnmarshalJSON(t *testing.T) {
	str := `{"id": "1234","customer_id": "1234","load_amount": "$123.45","time": "2018-01-01T00:00:00Z"}`
	
	l := Load{}
	err := json.Unmarshal([]byte(str), &l)
	if err != nil {
		t.Errorf("Error unmarshalling amount %s", err)
	}
	if l.Amount.Value != 123.45 {
		t.Error("Bad amount parsed")
	}
}
package spiral

import (
	"encoding/json"
	"strconv"
)

type OrderbookReturn struct {
	Symbol       string          `json:"symbol"`
	LastUpdateId int64           `json:"last_update_id"`
	Data         []OrderbookData `json:"data"`
	errorResponse
}

type OrderbookData struct {
	Price float64
	Size  float64
	Side  side
}

func (r *OrderbookData) UnmarshalJSON(bs []byte) error {
	arr := []interface{}{}
	if err := json.Unmarshal(bs, &arr); err != nil {
		return err
	}

	var err error
	if r.Price, err = strconv.ParseFloat(arr[0].(string), 64); err != nil {
		return err
	}
	if r.Size, err = strconv.ParseFloat(arr[1].(string), 64); err != nil {
		return err
	}
	r.Side = side(arr[2].(string))

	return nil
}

// Orderbook represents an orderbook from spiral api.
type Orderbook struct {
	Ask []OrderBookItem `json:"ask,struct"`
	Bid []OrderBookItem `json:"bid,struct"`
}

// OrderBookItem for Ask and Bid field.
type OrderBookItem struct {
	Price float64 `json:"price,string"`
	Size  float64 `json:"size,string"`
}

// UnmarshalJSON for OrderBook function
func (t *Orderbook) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Orderbook
	aux := &struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

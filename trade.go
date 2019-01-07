package spiral

type Trade struct {
	ID        int64   `json:"id"`
	Side      string  `json:"side"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price,string"`
	Quantity  float64 `json:"quantity,string"`
	Fee       float64 `json:"fee,string"`
	Timestamp int64   `json:"timestamp"`
}

type TradesReturn struct {
	Trades []Trade `json:"trades"`
	errorResponse
}

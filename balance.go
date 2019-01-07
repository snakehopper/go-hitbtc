package spiral

type Balance struct {
	Currency  string  `json:"currency"`
	Available float64 `json:"available,string"`
	Locked    float64 `json:"locked,string"`
	Timestamp int64   `json:"timestamp"`
}

type BalanceReturn struct {
	Data []Balance `json:"data"`
	errorResponse
}

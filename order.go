package spiral

type orderGetRequest struct {
	Symbol    string `json:"symbol"`
	Side      side   `json:"side"`
	Filter    string `json:"filter"`
	Count     int    `json:"count"`
	Reverse   bool   `json:"reverse"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

type Orders struct {
	Id             int64       `json:"id"`
	ClientOrderId  string      `json:"clt_ord_id"`
	Symbol         string      `json:"symbol"`
	Side           side        `json:"side"`
	Price          float64     `json:"price,string"`
	FilledPrice    float64     `json:"filled_price,string"`
	Quantity       float64     `json:"quantity,string"`
	FilledQuantity float64     `json:"filled_quantity,string"`
	Type           orderType   `json:"type"`
	Status         orderStatus `json:"status"`
	CreateTime     int64       `json:"create_time"`
	UpdateTime     int64       `json:"update_time"`
}

type OrdersReturn struct {
	Orders []Orders `json:"orders"`
	errorResponse
}

type PlaceReturn struct {
	Order PlaceData `json:"order"`
	errorResponse
}

type PlaceData struct {
	Id            int64       `json:"id"`
	ClientOrderId string      `json:"clt_ord_id"`
	UserId        int64       `json:"user_id"`
	Symbol        string      `json:"symbol"`
	Side          side        `json:"side"`
	Price         float64     `json:"price,string"`
	Quantity      float64     `json:"quantity,string"`
	Type          orderType   `json:"type"`
	Status        orderStatus `json:"status"`
	CreateTime    int64       `json:"create_time"`
	UpdateTime    int64       `json:"update_time"`
}

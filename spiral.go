// Package spiral is an implementation of the HitBTC API in Golang.
package spiral

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	API_BASE = "https://api.spiral.exchange/api/v1" // Spiral API endpoint
)

// New returns an instantiated HitBTC struct
func New(apiKey, apiSecret string) *Spiral {
	client := NewClient(apiKey, apiSecret)
	return &Spiral{client}
}

// NewWithCustomHttpClient returns an instantiated HitBTC struct with custom http client
func NewWithCustomHttpClient(apiKey, apiSecret string, httpClient *http.Client) *Spiral {
	client := NewClientWithCustomHttpConfig(apiKey, apiSecret, httpClient)
	return &Spiral{client}
}

// NewWithCustomTimeout returns an instantiated HitBTC struct with custom timeout
func NewWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) *Spiral {
	client := NewClientWithCustomTimeout(apiKey, apiSecret, timeout)
	return &Spiral{client}
}

// handleErr gets JSON response from spiral API en deal with error
func handleErr(r errorResponse) error {
	switch r.ErrorCode {
	case 0:
		return nil
	default:
		return errors.New(r.Message)
	}
}

// Spiral represent a Spiral client
type Spiral struct {
	client *client
}

// SetDebug sets enable/disable http request/response dump
func (b *Spiral) SetDebug(enable bool) {
	b.client.debug = enable
}

// GetCurrencies is used to get all supported currencies at Spiral along with other meta data.
func (b *Spiral) GetCurrencies() (currencies []Currency, err error) {
	r, err := b.client.do("GET", "currencies", nil, false)
	if err != nil {
		return
	}
	var resp CurrencyResponse
	if err = json.Unmarshal(r, &resp); err != nil {
		return
	}

	if err = handleErr(resp.errorResponse); err != nil {
		return
	}

	currencies = resp.Data
	return
}

// GetSymbols is used to get the open and available trading markets at Spiral along with other meta data.
func (b *Spiral) GetSymbols() (symbols []Symbol, err error) {
	r, err := b.client.do("GET", "products", nil, false)
	if err != nil {
		return
	}
	var response SymbolsResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response.errorResponse); err != nil {
		return
	}

	symbols = response.Data
	return
}

// GetKLines is used to fetch trading symbol kline data.
func (b *Spiral) GetKLines(market string, p period, limit int) (kline []KLine, err error) {
	params := map[string]string{
		"symbol": market,
		"period": string(p),
		"limit":  strconv.Itoa(limit),
	}
	r, err := b.client.do("GET", "klines", params, false)
	if err != nil {
		return
	}
	var response KLineReturn
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response.errorResponse); err != nil {
		return
	}

	kline = response.Data

	return
}

// GetOrderbook is used to get the current order book for a market.
func (b *Spiral) GetOrderbook(market string, limit int) (orderbook Orderbook, err error) {
	params := map[string]string{
		"symbol": market,
		"limit":  strconv.Itoa(limit),
	}
	r, err := b.client.do("GET", "orderbook", params, false)
	if err != nil {
		return
	}
	var response OrderbookReturn
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response.errorResponse); err != nil {
		return
	}

	count := len(response.Data)
	for i := 0; i < count; i++ {
		row := response.Data[i]
		switch row.Side {
		case BidSide:
			orderbook.Bid = append(orderbook.Bid, OrderBookItem{
				Price: row.Price,
				Size:  row.Size,
			})
		case AskSide:
			orderbook.Ask = append(orderbook.Ask, OrderBookItem{
				Price: row.Price,
				Size:  row.Size,
			})
		}
	}
	//rearrange bid array to make arr[0] is best deal price
	reverseSlice(orderbook.Bid)

	if len(orderbook.Bid) < 1 || len(orderbook.Ask) < 1 {
		err = fmt.Errorf("GetOrderBook() error, can not get enough Bid or Ask")
		return
	}

	return
}

func reverseSlice(a []OrderBookItem) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

// GetBalances is used to retrieve all balances from your account
func (b *Spiral) GetBalances() (balances []Balance, err error) {
	r, err := b.client.do("GET", "wallet/balances", nil, true)
	if err != nil {
		return
	}
	var response BalanceReturn
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response.errorResponse); err != nil {
		return
	}

	balances = response.Data
	return
}

// GetBalance is used to retrieve the balance from your account for a specific currency.
// currency: a string literal for the currency (ex: LTC)
func (b *Spiral) GetBalance(currency string) (balance Balance, err error) {
	params := map[string]string{
		"currency": currency,
	}

	r, err := b.client.do("GET", "wallet/balances", params, true)
	if err != nil {
		return
	}
	var response BalanceReturn
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response.errorResponse); err != nil {
		return
	}

	balance = response.Data[0]
	return
}

// GetTrades used to retrieve your trade history.
// market string literal for the market (ie. BTC/LTC). If set to "all", will return for all market
func (b *Spiral) GetTrades(symbol string, count int) (trades []Trade, err error) {
	payload := map[string]string{
		"symbol": symbol,
		"count":  "1000",
		//"start":  strSymbol
		//"reverse":strSymbol
		//"start_time":strSymbol
		//"end_time": strSymbol
	}
	if count > 0 {
		payload["count"] = strconv.Itoa(count)
	}

	r, err := b.client.do("GET", "trades", payload, true)
	if err != nil {
		return
	}
	var response TradesReturn
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response.errorResponse); err != nil {
		return
	}

	trades = response.Trades
	return
}

// CancelOrder cancels a pending order
func (b *Spiral) CancelOrder(orderId string) (err error) {
	params := map[string]string{
		"order_id": orderId,
	}
	r, err := b.client.do("DELETE", "order", params, true)
	if err != nil {
		return
	}
	var response errorResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}

	return nil
}

func (b *Spiral) CancelAllOrder(symbol, filter string) error {
	params := map[string]string{
		"symbol": symbol,
		"filter": filter,
	}
	r, err := b.client.do("DELETE", "order/all", params, true)
	if err != nil {
		return err
	}
	var response errorResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return err
	}
	if err = handleErr(response); err != nil {
		return err
	}

	return nil
}

// GetOrder gets a pending order data.
func (b *Spiral) GetOrder(orderId string) (orders []Orders, err error) {
	payload := make(map[string]string)
	payload["clientOrderId"] = orderId
	r, err := b.client.do("GET", "order", payload, true)
	if err != nil {
		return
	}
	var response OrdersReturn
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response.errorResponse); err != nil {
		return
	}

	orders = response.Orders
	return
}

// GetOrderHistory gets the history of orders for an user.
func (b *Spiral) GetOrderHistory(req orderGetRequest) (orders []Orders, err error) {
	params := map[string]string{
		"symbol":  req.Symbol,
		"side":    string(req.Side),
		"filter":  req.Filter,
		"count":   strconv.Itoa(req.Count),
		"reverse": fmt.Sprint(req.Reverse),
	}
	r, err := b.client.do("GET", "order", params, true)
	if err != nil {
		return
	}
	var response OrdersReturn
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response.errorResponse); err != nil {
		return
	}

	orders = response.Orders
	return
}

// GetOpenOrders gets the open orders of an user.
func (b *Spiral) GetOpenOrders(count int) (orders []Orders, err error) {
	filter := map[string]interface{}{"open": true}
	filterBytes, err := json.Marshal(filter)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"count":  strconv.Itoa(count),
		"filter": string(filterBytes),
	}

	r, err := b.client.do("GET", "order", params, true)
	if err != nil {
		return
	}
	var response OrdersReturn
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response.errorResponse); err != nil {
		return
	}

	orders = response.Orders
	return
}

// PlaceOrder creates a new order.
func (b *Spiral) PlaceOrder(requestOrder Orders) (resp PlaceReturn, err error) {
	payload := make(map[string]string)

	payload["clt_ord_id"] = requestOrder.ClientOrderId
	payload["symbol"] = requestOrder.Symbol
	payload["side"] = string(requestOrder.Side)
	payload["type"] = string(requestOrder.Type)
	payload["quantity"] = fmt.Sprintf("%.8f", requestOrder.Quantity)
	payload["price"] = fmt.Sprintf("%.8f", requestOrder.Price)

	r, err := b.client.do("POST", "order", payload, true)
	if err != nil {
		return
	}
	var response PlaceReturn
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response.errorResponse); err != nil {
		return
	}

	resp = response
	return
}

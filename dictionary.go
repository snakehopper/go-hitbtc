package spiral

type side string
type orderType string
type orderStatus string
type currency string
type period string

const (
	BidSide side = "bid"
	AskSide side = "ask"

	LimitOrderType  orderType = "limit"
	MarketOrderType orderType = "market"

	Submitted       orderStatus = "submitted"
	Accepted        orderStatus = "accepted"
	Waiting         orderStatus = "waiting"
	Rejected        orderStatus = "rejected"
	PartialFilled   orderStatus = "partial_filled"
	Filled          orderStatus = "filled"
	CancelRequested orderStatus = "cancel_requested"
	CancelRejected  orderStatus = "cancel_rejected"
	Cancelled       orderStatus = "cancelled"
	ModifyRequested orderStatus = "modify_requested"
	ModifyRejected  orderStatus = "modify_rejected"
	Modified        orderStatus = "modified"
	Unknown         orderStatus = "unknown"

	USDT currency = "USDT"
	BTC  currency = "BTC"
	ETH  currency = "ETH"
	LTC  currency = "LTC"
	BCH  currency = "BCH"

	Period1Minute   period = "1"
	Period5Minutes  period = "5"
	Period15Minutes period = "15"
	Period1Hour     period = "60"
)

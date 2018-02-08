package kiteconnect

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

// Order represents a single order entry.
type Order struct {
	AccountID string `json:"account_id"`
	PlacedBy  string `json:"placed_by"`

	OrderID                 string `json:"order_id"`
	ExchangeOrderID         string `json:"exchange_order_id"`
	ParentOrderID           string `json:"parent_order_id"`
	Status                  string `json:"status"`
	StatusMessage           string `json:"status_message"`
	OrderTimestamp          string `json:"order_timestamp"`
	ExchangeUpdateTimestamp string `json:"exchange_update_timestamp"`
	ExchangeTimestamp       string `json:"exchange_timestamp"`
	Meta                    string `json:"meta"`
	RejectedBy              string `json:"rejected_by"`
	Variety                 string `json:"variety"`

	Exchange        string `json:"exchange"`
	TradingSymbol   string `json:"tradingsymbol"`
	InstrumentToken int    `json:"instrument_token"`

	OrderType         string  `json:"order_type"`
	TransactionType   string  `json:"transaction_type"`
	Validity          string  `json:"validity"`
	Product           string  `json:"product"`
	Quantity          float64 `json:"quantity"`
	DisclosedQuantity float64 `json:"disclosed_quantity"`
	Price             float64 `json:"price"`
	TriggerPrice      float64 `json:"trigger_price"`

	AveragePrice      float64 `json:"average_price"`
	FilledQuantity    float64 `json:"filled_quantity"`
	PendingQuantity   float64 `json:"pending_quantity"`
	CancelledQuantity float64 `json:"cancelled_quantity"`
}

// Orders is a list of Order entries.
type Orders []Order

// OrderParams represents parameters for placing an order.
type OrderParams struct {
	Exchange        string `url:"exchange,omitempty"`
	Tradingsymbol   string `url:"tradingsymbol,omitempty"`
	Validity        string `url:"validity,omitempty"`
	Product         string `url:"product,omitempty"`
	OrderType       string `url:"order_type,omitempty"`
	TransactionType string `url:"transaction_type,omitempty"`

	Quantity          int     `url:"quantity,omitempty"`
	DisclosedQuantity int     `url:"disclosed_quantity,omitempty"`
	Price             float64 `url:"price,omitempty"`
	TriggerPrice      float64 `url:"trigger_price,omitempty"`

	Squareoff        float64 `url:"squareoff,omitempty"`
	Stoploss         float64 `url:"stoploss,omitempty"`
	TrailingStoploss float64 `url:"trailing_stoploss,omitempty"`

	Tag string `json:"tag" url:"tag,omitempty"`
}

// OrderResponse represents the result of a successful order placement.
type OrderResponse struct {
	OrderID string `json:"order_id"`
}

// Trade type as returned by the c++ api.
type Trade struct {
	AveragePrice      float64 `json:"average_price"`
	Quantity          float64 `json:"quantity"`
	TradeID           string  `json:"trade_id"`
	Product           string  `json:"product"`
	OrderTimestamp    string  `json:"order_timestamp"`
	ExchangeTimestamp string  `json:"exchange_timestamp"`
	ExchangeOrderID   string  `json:"exchange_order_id"`
	OrderID           string  `json:"order_id"`
	TransactionType   string  `json:"transaction_type"`
	TradingSymbol     string  `json:"tradingsymbol"`
	Exchange          string  `json:"exchange"`
	InstrumentToken   uint32  `json:"instrument_token"`
}

// Trades is a list of Trade entries.
type Trades []Trade

// GetOrders gets list of orders.
func (c *Client) GetOrders() (Orders, error) {
	var orders Orders
	err := c.doEnvelope(http.MethodGet, URIGetOrders, nil, nil, &orders)
	return orders, err
}

// GetTrades gets list of trades.
func (c *Client) GetTrades() (Trades, error) {
	var trades Trades
	err := c.doEnvelope(http.MethodGet, URIGetTrades, nil, nil, &trades)
	return trades, err
}

// GetOrderHistory gets history of individual order.
func (c *Client) GetOrderHistory(OrderID string) ([]Order, error) {
	var orderHistory []Order
	err := c.doEnvelope(http.MethodGet, fmt.Sprintf(URIGetOrderHistory, OrderID), nil, nil, &orderHistory)
	return orderHistory, err
}

// GetOrderTrades gets list of trades executed for a particular order.
func (c *Client) GetOrderTrades(OrderID string) ([]Trade, error) {
	var orderTrades []Trade
	err := c.doEnvelope(http.MethodGet, fmt.Sprintf(URIGetOrderTrades, OrderID), nil, nil, &orderTrades)
	return orderTrades, err
}

// PlaceOrder places an order.
func (c *Client) PlaceOrder(variety string, orderParams OrderParams) (OrderResponse, error) {
	var (
		orderResponse OrderResponse
		params        url.Values
		err           error
	)

	if params, err = query.Values(orderParams); err != nil {
		return orderResponse, NewError(InputError, fmt.Sprintf("Error decoding order params: %v", err), nil)
	}

	err = c.doEnvelope(http.MethodPost, fmt.Sprintf(URIPlaceOrder, variety), params, nil, &orderResponse)
	return orderResponse, err
}

// ModifyOrder modifies an order.
func (c *Client) ModifyOrder(variety string, orderID string, orderParams OrderParams) (OrderResponse, error) {
	var (
		orderResponse OrderResponse
		params        url.Values
		err           error
	)

	if params, err = query.Values(orderParams); err != nil {
		return orderResponse, NewError(InputError, fmt.Sprintf("Error decoding order params: %v", err), nil)
	}

	err = c.doEnvelope(http.MethodPut, fmt.Sprintf(URIModifyOrder, variety, orderID), params, nil, &orderResponse)
	return orderResponse, err
}

// CancelOrder cancels/exits an order.
func (c *Client) CancelOrder(variety string, orderID string, parentOrderID *string) (OrderResponse, error) {
	var (
		orderResponse OrderResponse
		params        url.Values
	)

	if parentOrderID != nil {
		params.Add("parent_order_id", *parentOrderID)
	}

	err := c.doEnvelope(http.MethodPut, fmt.Sprintf(URICancelOrder, variety, orderID), params, nil, &orderResponse)
	return orderResponse, err
}

// ExitOrder is an alias for CancelOrder which is used to cancel/exit an order.
func (c *Client) ExitOrder(variety string, orderID string, parentOrderID *string) (OrderResponse, error) {
	return c.CancelOrder(variety, orderID, parentOrderID)
}

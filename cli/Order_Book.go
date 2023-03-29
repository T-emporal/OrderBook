package cli

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/plot/plotter"
)

var _id, _id2 int = 1, 1
var _orders = make(map[int]Order)
var _bids = make(map[int]Bid)
var _xys plotter.XYs

//	var _desortedbid = make(map[int]Bid)

type Bid struct {
	bid      float64
	quantity float64
	ID       int
}

type Order struct {
	price    float64
	quantity float64
	duration int
}

type xy struct{ x, y float64 }

// type OrderBook struct {
// 	bids map[float64][]Order
// 	asks map[float64][]Order
// }

func InsertOrders() {

	var price, quantity float64
	var duration int

	fmt.Println("Enter Price of the order: ")
	fmt.Scanln(&price)

	fmt.Println("Enter Quantity: ")
	fmt.Scanln(&quantity)

	fmt.Println("Enter Duration: ")
	fmt.Scanln(&duration)

	order := Order{
		price:    price,
		quantity: quantity,
		duration: duration,
	}

	_orders[_id] = order
	_id += 1

}

func PlaceBid() {

	var bid, quantity float64
	var duration int

	fmt.Println("Enter the price of the bid: ")
	fmt.Scanln(&bid)

	fmt.Println("Enter the quantity you want to bid: ")
	fmt.Scanln(&quantity)

	fmt.Println("Enter the duration of the bid: ")
	fmt.Scanln(&duration)

	bids := Bid{
		bid:      bid,
		quantity: quantity,
		ID:       _id2,
	}

	_bids[duration] = bids
	_id2 += 1

}

func OrderMatchingMechanisum() {

	// var quantity, qxp float64 = 0, 0

	if len(_orders) == 0 {
		fmt.Println("Please Enter values in order book")
		return
	}

	if len(_bids) == 0 {
		fmt.Println("Please enter values to Bid")
		return
	}
	keys := make([]string, 0, len(_bids))

	//1. sort bid  orders in descending order of duration;
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	for duration, attrBid := range _bids {

		var quantity, qxp float64 = 0, 0

		//4. once the first bid is filled, move to the next one in the list with descending order of bid duration, and repeat steps 2 to 4
		for id, orders := range _orders {

			//2.check offer eligibility using following conditions:  ⁃
			//   i. price condition: (bid price) >= (offer price)
			//  ii. duration condition: (bid duration) <= (offer duration)
			if attrBid.bid >= orders.price && duration <= orders.duration {

				if _bids[duration].quantity == 0 {
					//fmt.Println("Order Quantity Empty for Order ID: ", id)
					continue
				}

				if _orders[id].quantity == 0 {
					//fmt.Println("Bid Quantity Empty")
					continue
				}

				//if attrBid.quantity >= orders.quantity {
				fmt.Println("\nBid ID: ", attrBid.ID)
				fmt.Println("Offer ID: ", id)

				minQuantity := math.Min(_bids[duration].quantity, _orders[id].quantity)
				minDuration := math.Min(float64(duration), float64(_orders[id].duration))
				minBidPrice := math.Min(attrBid.bid, orders.price)

				// fmt.Println("Bid Quantity1: ", _bids[duration].quantity)
				// fmt.Println("Order Quantity1: ", _orders[id].quantity)
				fmt.Println("Quantity: ", minQuantity)
				//fmt.Println("difference1 :", attrBid.quantity-orders.quantity)

				// fmt.Println("Bid Duration: ", _bids[duration].quantity)
				// fmt.Println("Order Duration: ", _orders[id].quantity)
				fmt.Println("Duration: ", minDuration)

				// fmt.Println("Bid Price: ", attrBid.bid)
				// fmt.Println("Order Price: ", orders.price)
				fmt.Println("Price:", minBidPrice)
				fmt.Println("Amount: ", minBidPrice*minQuantity)

				// fmt.Println("Bid Quantity2: ", _bids[duration].quantity)
				// fmt.Println("Difference: ", _bids[duration].quantity-minQuantity)
				// fmt.Println("_orders Quantity2: ", _orders[id].quantity)
				//fmt.Println("Difference: ", _orders[id].quantity-minQuantity)
				quantity += minQuantity

				qxp += minQuantity * minBidPrice

				_xys = append(_xys, struct{ X, Y float64 }{minDuration, minBidPrice})

				//3. when price condition and duration conditions are both met:
				//  ⁃ satisfy as much of bid order as possible with eligible offer having minimum duration
				//  - if any of the bid is left unfilled, move to the next lowest duration offer
				//  - for bids which are filled, i. price will be the lowest of bid and offer prices;
				//								ii. duration will be lowest of bid and offer durations
				if entry, ok := _bids[duration]; ok {
					entry.quantity = _bids[duration].quantity - minQuantity
					_bids[duration] = entry
				}

				if entry, ok := _orders[id]; ok {
					entry.quantity = _orders[id].quantity - minQuantity
					_orders[id] = entry
				}

			}
		}

		if quantity != 0 {
			fmt.Println("\nTotal Quantity :", quantity)
			fmt.Println("Price: ", qxp/quantity)
			fmt.Println("Amount: ", qxp)
		}

	}

	PlotGraph()

}

func ShowBids() {

	for keys, values := range _bids {
		fmt.Println("ID: ", values.ID, "Quantity: ", values.quantity, "Duration: ", keys, "Price: ", values.bid)
	}
}

func ShowOrderBook() {

	for keys, values := range _orders {
		fmt.Println("ID: ", keys, "Quantity: ", values.quantity, "Duration: ", values.duration, "Price: ", values.price)
	}
}

func CreateOrders() {
	var price, quantity float64
	var duration int

	price = 10
	quantity = 20
	duration = 15

	order := Order{
		price:    price,
		quantity: quantity,
		duration: duration,
	}

	_orders[_id] = order
	_id += 1

	price = 15
	quantity = 25
	duration = 20

	order = Order{
		price:    price,
		quantity: quantity,
		duration: duration,
	}

	_orders[_id] = order
	_id += 1

	price = 20
	quantity = 30
	duration = 25

	order = Order{
		price:    price,
		quantity: quantity,
		duration: duration,
	}

	_orders[_id] = order
	_id += 1

	price = 25
	quantity = 35
	duration = 30

	order = Order{
		price:    price,
		quantity: quantity,
		duration: duration,
	}

	_orders[_id] = order
	_id += 1

	price = 30
	quantity = 40
	duration = 35

	order = Order{
		price:    price,
		quantity: quantity,
		duration: duration,
	}

	_orders[_id] = order
	_id += 1
}

func CreateBid() {

	var bid, quantity float64

	bid = 10
	quantity = 15

	bids := Bid{
		bid:      bid,
		quantity: quantity,
		ID:       _id2,
	}

	_bids[10] = bids
	_id2 += 1

	bid = 12
	quantity = 20

	bids = Bid{
		bid:      bid,
		quantity: quantity,
		ID:       _id2,
	}

	_bids[13] = bids
	_id2 += 1

	bid = 15
	quantity = 22

	bids = Bid{
		bid:      bid,
		quantity: quantity,
		ID:       _id2,
	}

	_bids[15] = bids
	_id2 += 1

	bid = 20
	quantity = 35

	bids = Bid{
		bid:      bid,
		quantity: quantity,
		ID:       _id2,
	}

	_bids[16] = bids
	_id2 += 1

	bid = 30
	quantity = 40

	bids = Bid{
		bid:      bid,
		quantity: quantity,
		ID:       _id2,
	}

	_bids[20] = bids
	_id2 += 1

	bid = 21
	quantity = 45

	bids = Bid{
		bid:      bid,
		quantity: quantity,
		ID:       _id2,
	}

	_bids[18] = bids
	_id2 += 1

	bid = 24
	quantity = 35

	bids = Bid{
		bid:      bid,
		quantity: quantity,
		ID:       _id2,
	}

	_bids[24] = bids
	_id2 += 1
}

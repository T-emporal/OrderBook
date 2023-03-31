package cli

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/plot/plotter"
)

var _id, _id2 int = 1, 1
var _offers = make(map[int]Order)
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

	_offers[_id] = order
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

	if len(_offers) == 0 {
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
		var eligiableOrderDuration []int
		var eligiableIds []int

		//4. once the first bid is filled, move to the next one in the list with descending order of bid duration, and repeat steps 2 to 4
		for id, orders := range _offers {

			//2.check offer eligibility using following conditions:  ⁃
			//   i. price condition: (bid price) >= (offer price)
			//  ii. duration condition: (bid duration) <= (offer duration)
			if attrBid.bid >= orders.price && duration <= orders.duration {
				eligiableOrderDuration = append(eligiableOrderDuration, _offers[id].duration)
			}
		}
		sort.Slice(eligiableOrderDuration, func(i, j int) bool {
			return eligiableOrderDuration[i] < eligiableOrderDuration[j]
		})

		for _, orderDuration := range eligiableOrderDuration {

			for orderId, _ := range _offers {
				if orderDuration == _offers[orderId].duration {
					eligiableIds = append(eligiableIds, orderId)
				}
			}
		}

		//3. when price condition and duration conditions are both met:
		//  ⁃ satisfy as much of bid order as possible with eligible offer having minimum duration
		//  - if any of the bid is left unfilled, move to the next lowest duration offer
		//  - for bids which are filled:
		//         i. price will be the lowest of bid and offer prices;
		//		   ii. duration will be lowest of bid and offer durations
		for _, eligiableID := range eligiableIds {

			if _bids[duration].quantity == 0 {
				// fmt.Println("bids Quantity Empty ")
				continue
			}

			if _offers[eligiableID].quantity == 0 {
				// fmt.Println("Offers Quantity Empty")
				continue
			}
			fmt.Println("\nBid ID: ", attrBid.ID)
			fmt.Println("Offer ID: ", eligiableID)

			minQuantity := math.Min(_bids[duration].quantity, _offers[eligiableID].quantity)
			minDuration := math.Min(float64(duration), float64(_offers[eligiableID].duration))
			minBidPrice := math.Min(attrBid.bid, _offers[eligiableID].price)

			fmt.Println("Quantity: ", minQuantity)
			fmt.Println("Duration: ", minDuration)
			fmt.Println("Price:", minBidPrice)
			fmt.Println("Amount: ", minBidPrice*minQuantity)

			quantity += minQuantity

			qxp += minQuantity * minBidPrice

			_xys = append(_xys, struct{ X, Y float64 }{minDuration, minBidPrice})

			if entry, ok := _bids[duration]; ok {
				entry.quantity = _bids[duration].quantity - minQuantity
				_bids[duration] = entry
			}

			if entry, ok := _offers[eligiableID]; ok {
				entry.quantity = _offers[eligiableID].quantity - minQuantity
				_offers[eligiableID] = entry
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

	for keys, values := range _offers {
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

	_offers[_id] = order
	_id += 1

	price = 15
	quantity = 25
	duration = 20

	order = Order{
		price:    price,
		quantity: quantity,
		duration: duration,
	}

	_offers[_id] = order
	_id += 1

	price = 20
	quantity = 30
	duration = 25

	order = Order{
		price:    price,
		quantity: quantity,
		duration: duration,
	}

	_offers[_id] = order
	_id += 1

	price = 25
	quantity = 35
	duration = 30

	order = Order{
		price:    price,
		quantity: quantity,
		duration: duration,
	}

	_offers[_id] = order
	_id += 1

	price = 30
	quantity = 40
	duration = 35

	order = Order{
		price:    price,
		quantity: quantity,
		duration: duration,
	}

	_offers[_id] = order
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

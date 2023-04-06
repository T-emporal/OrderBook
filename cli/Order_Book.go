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
	deSortedDuration := make([]int, 0, len(_bids))

	//  Arrange The Bids In Descending Order Of Duration
	for durations, _ := range _bids {
		deSortedDuration = append(deSortedDuration, durations)
	}

	sort.SliceStable(deSortedDuration, func(i, j int) bool {
		return deSortedDuration[i] < deSortedDuration[j]
	})

	// 	If Any Portion Of Given Bid
	// 	Are Left Unfilled, Move To The Next Lowest Duration Offer.
	for _, duration := range deSortedDuration {
		// fmt.Println("----duration:", duration)
		var quantity, qxp float64 = 0, 0
		var eligiableOrderDuration []int
		var eligiableIds []int

		for id, orders := range _offers {

			// Condition Check:
			// 		1. Bid Price >= Offer Price
			// 		2. Bid Duration <= Offer Duration
			// fmt.Println("Offer ID: ", id)
			// fmt.Println("Bid ID: ", _bids[duration].ID)
			// fmt.Println("price-----", _bids[duration].bid >= orders.price)
			// fmt.Println("price-----", _bids[duration].bid, orders.price)
			// fmt.Println("duration-----", duration <= orders.duration)
			// fmt.Println("duration-----", duration, orders.duration)
			if _bids[duration].bid >= orders.price && duration <= orders.duration {
				eligiableOrderDuration = append(eligiableOrderDuration, _offers[id].duration)
			}
		}
		fmt.Println(eligiableOrderDuration)
		// Arrange Eligible Offers In Ascending Order Of Duration
		sort.Slice(eligiableOrderDuration, func(i, j int) bool {
			return eligiableOrderDuration[i] < eligiableOrderDuration[j]
		})
		// fmt.Println(eligiableOrderDuration)

		for _, orderDuration := range eligiableOrderDuration {
			// fmt.Println("--------orderDuration", orderDuration)
			for orderId, _ := range _offers {
				if orderDuration == _offers[orderId].duration {
					eligiableIds = append(eligiableIds, orderId)
				}
			}
		}
		// fmt.Println(eligiableIds)

		// Satisfy Bid:
		//  For Each Portion Of The Given Bid Filled By A Given Offer:
		//   - Price Will Be The Lowest Of Bid And Offer Prices
		//   - Duration Will Be Lowest Of Bid And Offer Durations
		for _, eligiableID := range eligiableIds {
			// fmt.Println("eligiableID", eligiableID)
			if _bids[duration].quantity == 0 {
				// fmt.Println("bids Quantity Empty ", _bids[duration].ID)
				continue
			}

			if _offers[eligiableID].quantity == 0 {
				// fmt.Println("Offers Quantity Empty for ID: ", eligiableID)
				continue
			}
			fmt.Println("\nBid ID: ", _bids[duration].ID)
			fmt.Println("Offer ID: ", eligiableID)

			minQuantity := math.Min(_bids[duration].quantity, _offers[eligiableID].quantity)
			minDuration := math.Min(float64(duration), float64(_offers[eligiableID].duration))
			minBidPrice := math.Min(_bids[duration].bid, _offers[eligiableID].price)

			fmt.Println("Quantity: ", minQuantity)
			fmt.Println("Duration: ", minDuration)
			fmt.Println("Price:", minBidPrice)
			fmt.Println("Amount: ", minBidPrice*minQuantity)

			quantity += minQuantity
			// Trade Executed
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

	fmt.Println("ID , Quantity , Duration , Price ")
	for keys, values := range _bids {
		fmt.Println(values.ID, ",", values.quantity, ",", keys, ",", values.bid)
	}
}

func ShowOrderBook() {

	fmt.Println("ID , Quantity , Duration , Price ")
	for keys, values := range _offers {
		fmt.Println(keys, ",", values.quantity, ",", values.duration, ",", values.price)
	}
}

func CreateOrders() {

	quantity := []float64{35, 30, 23, 33, 15, 12}
	duration := []int{20, 45, 26, 36, 10, 45}
	price := []float64{12, 9, 11, 10, 13, 14}

	for i := 0; i < len(quantity); i++ {
		_offers[i+1] = Order{price[i], quantity[i], duration[i]}
	}

}

func CreateBid() {

	bid := []float64{8, 11, 7, 15, 12, 6, 3, 5}
	quantity := []float64{20, 13, 23, 16, 21, 13, 10, 7}
	duration := []int{15, 12, 13, 19, 18, 22, 29, 11}

	for i := 0; i < len(quantity); i++ {
		_bids[duration[i]] = Bid{bid[i], quantity[i], i + 1}
	}

}

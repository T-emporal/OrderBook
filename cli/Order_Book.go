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

	keys := make([]int, 0, len(_bids))

	for k := range _bids {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	fmt.Println("ID , Quantity , Duration , Price ")
	for _, k := range keys {
		fmt.Println(_bids[k].ID, ",", _bids[k].quantity, ",", k, ",", _bids[k].bid)
	}
}

func ShowOrderBook() {

	keys := make([]int, 0, len(_offers))

	for k := range _offers {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	fmt.Println("ID , Quantity , Duration , Price ")
	for _, k := range keys {
		fmt.Println(k, ",", _offers[k].quantity, ",", _offers[k].duration, ",", _offers[k].price)
	}
}

func CreateOrders() {

	duration := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	quantity := []float64{90, 89, 88, 87, 86, 86, 85, 84, 83, 82, 81, 81, 80, 79, 78, 77, 77, 76, 75, 74, 74, 73, 72, 71, 71, 70, 69, 69, 68, 67, 67, 66, 65, 65, 64, 63, 63, 62, 61, 61, 60, 60, 59, 58, 58, 57, 57, 56, 56, 55, 54, 54, 53, 53, 52, 52, 51, 51, 50, 50, 49, 49, 48, 48, 47, 47, 46, 46, 45, 45, 45, 44, 44, 43, 43, 42, 42, 42, 41, 41, 40, 40, 39, 39, 39, 38, 38, 38, 37, 37}
	price := []float64{5.50, 5.51, 5.51, 5.52, 5.52, 5.53, 5.53, 5.54, 5.54, 5.55, 5.56, 5.56, 5.57, 5.57, 5.58, 5.58, 5.59, 5.59, 5.60, 5.61, 5.61, 5.62, 5.62, 5.63, 5.63, 5.64, 5.64, 5.65, 5.66, 5.66, 5.67, 5.67, 5.68, 5.68, 5.69, 5.70, 5.70, 5.71, 5.71, 5.72, 5.72, 5.73, 5.74, 5.74, 5.75, 5.75, 5.76, 5.76, 5.77, 5.78, 5.78, 5.79, 5.79, 5.80, 5.81, 5.81, 5.82, 5.82, 5.83, 5.83, 5.84, 5.85, 5.85, 5.86, 5.86, 5.87, 5.88, 5.88, 5.89, 5.89, 5.90, 5.90, 5.91, 5.92, 5.92, 5.93, 5.93, 5.94, 5.95, 5.95, 5.96, 5.96, 5.97, 5.98, 5.98, 5.99, 5.99, 6.00, 6.01, 6.01}

	for i := 0; i < len(quantity); i++ {
		_offers[i+1] = Order{price[i], quantity[i], duration[i]}
	}

}

func CreateBid() {

	duration := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	bid := []float64{5.00, 5.04, 5.07, 5.11, 5.14, 5.18, 5.21, 5.25, 5.29, 5.32, 5.36, 5.40, 5.44, 5.47, 5.51, 5.55, 5.59, 5.63, 5.67, 5.71, 5.75, 5.79, 5.83, 5.87, 5.91, 5.95, 5.99, 6.04, 6.08, 6.12, 6.16, 6.21, 6.25, 6.29, 6.34, 6.38, 6.43, 6.47, 6.52, 6.56, 6.61, 6.66, 6.70, 6.75, 6.80, 6.84, 6.89, 6.94, 6.99, 7.04, 7.09, 7.14, 7.19, 7.24, 7.29, 7.34, 7.39, 7.44, 7.49, 7.55, 7.60, 7.65, 7.71, 7.76, 7.81, 7.87, 7.92, 7.98, 8.03, 8.09, 8.15, 8.20, 8.26, 8.32, 8.38, 8.44, 8.50, 8.56, 8.62, 8.68, 8.74, 8.80, 8.86, 8.92, 8.98, 9.05, 9.11, 9.17, 9.24, 9.30}
	quantity := []float64{100, 99, 98, 97, 96, 95, 94, 93, 92, 91, 90, 90, 89, 88, 87, 86, 85, 84, 83, 83, 82, 81, 80, 79, 79, 78, 77, 76, 75, 75, 74, 73, 72, 72, 71, 70, 70, 69, 68, 68, 67, 66, 66, 65, 64, 64, 63, 62, 62, 61, 61, 60, 59, 59, 58, 58, 57, 56, 56, 55, 55, 54, 54, 53, 53, 52, 52, 51, 50, 50, 49, 49, 48, 48, 48, 47, 47, 46, 46, 45, 45, 44, 44, 43, 43, 43, 42, 42, 41, 41}

	for i := 0; i < len(quantity); i++ {
		_bids[duration[i]] = Bid{bid[i], quantity[i], i + 1}
	}

}

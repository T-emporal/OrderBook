package cli

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/plot/plotter"
)

var (
	_id, _id2 int = 1, 1
	_offers       = make(map[int]Order)
	_bids         = make(map[int]Bid)
	_xys      plotter.XYs
)

type (
	Bid struct {
		bid      float64
		quantity float64
		ID       int
	}

	Order struct {
		price    float64
		quantity float64
		duration int
	}

	xy struct{ x, y float64 }
)

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
	sortedDuration := make([]int, 0, len(_bids))

	//  Arrange The Bids In Descending Order Of Duration
	for durations, _ := range _bids {
		sortedDuration = append(sortedDuration, durations)
	}

	sort.SliceStable(sortedDuration, func(i, j int) bool {
		return sortedDuration[i] < sortedDuration[j]
	})

	// 	If Any Portion Of Given Bid
	// 	Are Left Unfilled, Move To The Next Lowest Duration Offer.
	for _, duration := range sortedDuration {
		// fmt.Println("----duration:", duration)
		var quantity, qxp float64 = 0, 0
		var eligiableOrderDuration []int
		var eligiableIds []int

		for id, orders := range _offers {

			// Condition Check:
			// 		1. Bid Price >= Offer Price
			// 		2. Bid Duration <= Offer Duration
			// fmt.Println("Offer ID: ", id)

			// fmt.Println("\nBids ID: ", _bids[duration].ID)
			// fmt.Println("Bids Duration: ", duration)
			// fmt.Println("Bids Quantity: ", _bids[duration].quantity)
			// fmt.Println("Bids Price : ", _bids[duration].bid)

			// fmt.Println("Offer ID: ", id)
			// fmt.Println("Offers Duration: ", orders.duration)
			// fmt.Println("Offers Quantity: ", orders.quantity)
			// fmt.Println("Offer Price : ", orders.price)

			// fmt.Println("Price Check : ", _bids[duration].bid >= orders.price)

			// fmt.Println("Duration Check: ", duration <= orders.duration)
			if _bids[duration].bid >= orders.price && duration <= orders.duration {
				eligiableOrderDuration = append(eligiableOrderDuration, _offers[id].duration)
			}
		}
		//fmt.Println(eligiableOrderDuration)
		// Arrange Eligible Offers In Ascending Order Of Duration
		sort.Slice(eligiableOrderDuration, func(i, j int) bool {
			return eligiableOrderDuration[i] < eligiableOrderDuration[j]
		})
		// fmt.Println(eligiableOrderDuration)

		for _, orderDuration := range eligiableOrderDuration {
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
	//Checking the Plot of the Graph
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

func CreateOffers() {

	duration := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	quantity := []float64{90, 87, 85, 82, 80, 77, 75, 73, 71, 68, 66, 64, 62, 61, 59, 57, 55, 54, 52, 50, 49, 47, 46, 45, 43, 42, 41, 40, 38, 37, 36, 35, 34, 33, 32, 31, 30, 29, 28, 27, 27, 26, 25, 24, 24, 23, 22, 22, 21, 20, 20, 19, 18, 18, 17, 17, 16, 16, 15, 15, 14, 14, 14, 13, 13, 12, 12, 12, 11, 11, 11, 10, 10, 10, 9, 9, 9, 9, 8, 8, 8, 8, 7, 7, 7, 7, 7, 6, 6, 6}
	price := []float64{5.50, 5.54, 5.58, 5.62, 5.66, 5.70, 5.74, 5.78, 5.82, 5.86, 5.90, 5.94, 5.98, 6.02, 6.06, 6.11, 6.15, 6.19, 6.24, 6.28, 6.32, 6.37, 6.41, 6.46, 6.50, 6.55, 6.59, 6.64, 6.69, 6.73, 6.78, 6.83, 6.88, 6.92, 6.97, 7.02, 7.07, 7.12, 7.17, 7.22, 7.27, 7.32, 7.37, 7.42, 7.48, 7.53, 7.58, 7.63, 7.69, 7.74, 7.80, 7.85, 7.90, 7.96, 8.02, 8.07, 8.13, 8.19, 8.24, 8.30, 8.36, 8.42, 8.48, 8.54, 8.60, 8.66, 8.72, 8.78, 8.84, 8.90, 8.96, 9.03, 9.09, 9.15, 9.22, 9.28, 9.35, 9.41, 9.48, 9.54, 9.61, 9.68, 9.74, 9.81, 9.88, 9.95, 10.02, 10.09, 10.16, 10.23}

	for i := 0; i < len(quantity); i++ {
		_offers[i+1] = Order{price[i], quantity[i], duration[i]}
	}

}

func CreateBid() {

	duration := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	bid := []float64{6.00, 6.03, 6.06, 6.09, 6.12, 6.15, 6.18, 6.21, 6.24, 6.28, 6.31, 6.34, 6.37, 6.40, 6.43, 6.47, 6.50, 6.53, 6.56, 6.60, 6.63, 6.66, 6.70, 6.73, 6.76, 6.80, 6.83, 6.86, 6.90, 6.93, 6.97, 7.00, 7.04, 7.07, 7.11, 7.14, 7.18, 7.22, 7.25, 7.29, 7.32, 7.36, 7.40, 7.44, 7.47, 7.51, 7.55, 7.59, 7.62, 7.66, 7.70, 7.74, 7.78, 7.82, 7.85, 7.89, 7.93, 7.97, 8.01, 8.05, 8.09, 8.13, 8.17, 8.22, 8.26, 8.30, 8.34, 8.38, 8.42, 8.46, 8.51, 8.55, 8.59, 8.64, 8.68, 8.72, 8.77, 8.81, 8.85, 8.90, 8.94, 8.99, 9.03, 9.08, 9.12, 9.17, 9.21, 9.26, 9.31, 9.35}
	quantity := []float64{95, 92, 89, 87, 84, 82, 79, 77, 74, 72, 70, 68, 66, 64, 62, 60, 58, 57, 55, 53, 52, 50, 49, 47, 46, 44, 43, 42, 40, 39, 38, 37, 36, 35, 34, 33, 32, 31, 30, 29, 28, 27, 26, 26, 25, 24, 23, 23, 22, 21, 21, 20, 19, 19, 18, 18, 17, 17, 16, 16, 15, 15, 14, 14, 14, 13, 13, 12, 12, 12, 11, 11, 11, 10, 10, 10, 9, 9, 9, 9, 8, 8, 8, 8, 7, 7, 7, 7, 7, 6}

	for i := 0; i < len(quantity); i++ {
		_bids[duration[i]] = Bid{bid[i], quantity[i], i + 1}
	}

}

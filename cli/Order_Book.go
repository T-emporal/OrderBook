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

func (order Order) MatchBidPerOffer() {

	fmt.Println("Entering ")

	var eligiableOrderDuration []int
	var eligiableIds []int

	for id, orders := range _offers {

		if order.price >= orders.price && order.duration <= orders.duration {
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

	for _, eligiableID := range eligiableIds {

		if order.quantity == 0 {
			break
		}

		if _offers[eligiableID].quantity == 0 {
			continue
		}

		minQuantity := math.Min(order.quantity, _offers[eligiableID].quantity)
		minDuration := math.Min(float64(order.duration), float64(_offers[eligiableID].duration))
		minBidPrice := math.Min(order.price, _offers[eligiableID].price)

		fmt.Println("Quantity: ", minQuantity)
		fmt.Println("Duration: ", minDuration)
		fmt.Println("Price:", minBidPrice)
		fmt.Println("Amount: ", minBidPrice*minQuantity)

		if entry, ok := _offers[eligiableID]; ok {
			entry.quantity = _offers[eligiableID].quantity - minQuantity
			_offers[eligiableID] = entry
		}

		order.quantity = order.quantity - minQuantity

	}

}

func MatchOfferPerBid() {

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

	order.MatchBidPerOffer()

}

func CreateOffers() {

	duration := []int{1, 2, 3, 4, 5, 6, 7}
	quantity := []float64{90, 87, 85, 82, 80, 77, 75}
	price := []float64{5.50, 5.54, 5.58, 5.62, 5.66, 5.70, 5.74}

	for i := 0; i < len(quantity); i++ {
		_offers[i+1] = Order{price[i], quantity[i], duration[i]}
	}
	_id = len(quantity) + 1

}

func CreateBid() {

	duration := []int{1, 2, 3, 4, 5, 6, 7}
	bid := []float64{6.00, 6.03, 6.06, 6.09, 6.12, 6.15, 6.18}
	quantity := []float64{95, 92, 89, 87, 84, 82, 79}

	for i := 0; i < len(quantity); i++ {
		_bids[duration[i]] = Bid{bid[i], quantity[i], i + 1}
	}
	_id2 = len(quantity) + 1
}

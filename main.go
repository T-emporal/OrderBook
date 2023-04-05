package main

import (
	"OrderBook/cli"
	"fmt"
)

func main() {

	var number int
	cli.CreateOrders()
	cli.CreateBid()

	for {
		fmt.Println("\nEnter 1 to Enter in Order Book ")
		fmt.Println("Enter 2 Show the offers ")
		fmt.Println("Enter 3 to place a bid ")
		fmt.Println("Enter 4 to show bids")
		fmt.Println("Enter 5 to Execute Trade ")
		fmt.Scanln(&number)

		switch number {
		case 1:
			cli.InsertOrders()
		case 2:
			cli.ShowOrderBook()
		case 3:
			cli.PlaceBid()
		case 4:
			cli.ShowBids()
		case 5:
			cli.OrderMatchingMechanisum()
		default:
			fmt.Println("Please Provide Input")
		}

	}

}

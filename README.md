# Temporal 
Temporal Orderbook Simulation

## Introduction

Temporal’s market mechanism unifies Temporally Discrete Markets such as Futures and
Lending & Borrowing into Realtime, Continuous, Forward Curves shaped purely by market
forces. This allows users to enter Market-Priced, Custom-Maturity Futures and Lending &
Borrowing contracts, creating an entirely new DeFi primitive; in contrast to typical
exchanges which offer a few standardized-maturity contracts with siloed liquidity.

This repository aims to provide Temporal's OrderBook Simulation/Prototype 

## Prerequisites

- Golang (1.19 or higher)

## Getting Started 
- Clone this Repository
- Run code through `go run main.go`

## Code Walk-Through
- Simulation has already been stored with initial Offers and Bids which could be displayed to execute Trade
- Initially, when we initiate Simulation it gives options with keys `2` and `4` to add more Offers and Bids, also we can see all Offers and Bids with keys `1` and `3` for both pre and post-trade
- When Executing Trade with key `5` we can see the Order Matching Mechanism. For every Bid ID, we could see which offer IDs are matched and which trade had taken place. Total Quantity, Price, and Amount are shown for every Bid ID
- A plot is drawn, against Price and Duration which shows the Continuous, Forward Curves shaped purely by market forces, which could be seen in the working directory as `out.png`

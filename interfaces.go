package main

type ClientHandler interface {
	SendOrder(customer Customer, order Order) *ExpressResponse
}
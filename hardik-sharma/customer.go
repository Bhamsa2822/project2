package main

type Customer struct {
	Id              string          `json:"id"`
	CustomerDetails CustomerDetails `json:"customerDetails" bun:"embed:customerdetails_"`
}

type CustomerDetails struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	ContactNo int    `json:"contactNo"`
}

package main

//Address Structure
type Address struct {
	StreetAddress string `json:"streetaddress,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	Country       string `json:"country,omitempty"`
	CountryCode   string `json:"countryCode,omitempty"`
}

type PaymentInfo struct {
	OrderId     string
	Qty         string
	PaymnetDate string
	Amount      string
}

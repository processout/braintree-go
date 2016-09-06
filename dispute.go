package braintree

import "encoding/xml"

type Dispute struct {
	XMLName            xml.Name           `xml:"dispute"`
	Id                 string             `xml:"id"`
	Amount             *Decimal           `xml:"amount"`
	Kind               string             `xml:"kind"`
	Reason             string             `xml:"reason"`
	TransactionDetails TransactionDetails `xml:"transaction-details"`
}

type TransactionDetails struct {
	Id     string   `xml:"id"`
	Amount *Decimal `xml:"amount"`
}

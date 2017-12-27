package braintree

import (
	"encoding/xml"
	"time"
)

type Dispute struct {
	XMLName            xml.Name           `xml:"dispute"`
	Id                 string             `xml:"id"`
	Amount             *Decimal           `xml:"amount"`
	AmountDisputed     *Decimal           `xml:"amount-disputed"`
	CurrencyISOCode    string             `xml:"currency-iso-code"`
	AmountWon          *Decimal           `xml:"amount-won"`
	CreatedAt          *time.Time         `xml:"created-at"`
	DateWon            *time.Time         `xml:"date-won"`
	Kind               string             `xml:"kind"`
	Reason             string             `xml:"reason"`
	TransactionDetails TransactionDetails `xml:"transaction-details"`

	Status        string                `xml:"status"`
	StatusHistory *DisputeStatusHistory `xml:"status-history"`
}

type TransactionDetails struct {
	Id     string   `xml:"id"`
	Amount *Decimal `xml:"amount"`
}

type DisputeStatusHistory struct {
	Events []DisputeStatusEvent `xml:"status-event"`
}

type DisputeStatusEvent struct {
	Timestamp string `xml:"timestamp"`
	Status    string `xml:"status"`
}

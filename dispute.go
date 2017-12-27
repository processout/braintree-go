package braintree

import (
	"encoding/xml"
	"time"
)

// Dispute object for a *Transaction*
type Dispute struct {
	XMLName xml.Name `xml:"dispute"`

	Id              string     `xml:"id"`
	Amount          *Decimal   `xml:"amount"`
	CurrencyISOCode string     `xml:"currency-iso-code"`
	AmountWon       *Decimal   `xml:"amount-won"`
	ReceivedDate    *time.Time `xml:"received-date"`
	DateWon         *time.Time `xml:"date-won"`
	Kind            string     `xml:"kind"`
	Reason          string     `xml:"reason"`

	Status string `xml:"status"`
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

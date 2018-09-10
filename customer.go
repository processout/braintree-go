package braintree

import (
	"github.com/processout/braintree-go/customfields"
	"github.com/processout/braintree-go/nullable"
)

type Customer struct {
	XMLName            string                    `xml:"customer"`
	Id                 string                    `xml:"id,omitempty"`
	FirstName          string                    `xml:"first-name,omitempty"`
	LastName           string                    `xml:"last-name,omitempty"`
	Company            string                    `xml:"company,omitempty"`
	Email              string                    `xml:"email,omitempty"`
	Phone              string                    `xml:"phone,omitempty"`
	Fax                string                    `xml:"fax,omitempty"`
	Website            string                    `xml:"website,omitempty"`
	CustomFields       customfields.CustomFields `xml:"custom-fields"`
	CreditCard         *CreditCard               `xml:"credit-card,omitempty"`
	CreditCards        *CreditCards              `xml:"credit-cards,omitempty"`
	PaymentMethodNonce string                    `xml:"payment-method-nonce,omitempty"`
	PaymentMethods     []struct {
		Token string `xml:"token,omitempty"`
	} `xml:"payment-methods,omitempty"`
}

/* :credit_card => {
   :payment_method_nonce => nonce_from_the_client,
   :options => {
     :verify_card => true
   }
*/

// DefaultCreditCard returns the default credit card, or nil
func (c *Customer) DefaultCreditCard() *CreditCard {
	for _, card := range c.CreditCards.CreditCard {
		if card.Default {
			return card
		}
	}
	return nil
}

type CustomerSearchResult struct {
	XMLName           string              `xml:"customers"`
	CurrentPageNumber *nullable.NullInt64 `xml:"current-page-number"`
	PageSize          *nullable.NullInt64 `xml:"page-size"`
	TotalItems        *nullable.NullInt64 `xml:"total-items"`
	Customers         []*Customer         `xml:"customer"`
}

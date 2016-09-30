package braintree

import (
	"encoding/xml"
	"fmt"
)

type TransactionGateway struct {
	*Braintree
}

// Create initiates a transaction.
func (g *TransactionGateway) Create(tx *Transaction) (*Transaction, error) {
	resp, err := g.execute("POST", "transactions", tx)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 201:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// SubmitForSettlement submits the transaction with the specified id for settlement.
// If the amount is omitted, the full amount is settled.
func (g *TransactionGateway) SubmitForSettlement(id string, amount ...*Decimal) (*Transaction, error) {
	var tx *Transaction
	if len(amount) > 0 {
		tx = &Transaction{
			Amount: amount[0],
		}
	}
	resp, err := g.execute("PUT", "transactions/"+id+"/submit_for_settlement", tx)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// Settle settles a transaction.
// This action is only available in the sandbox environment.
func (g *TransactionGateway) Settle(id string) (*Transaction, error) {
	if g.Environment != Production {
		resp, err := g.execute("PUT", "transactions/"+id+"/settle", nil)
		if err != nil {
			return nil, err
		}
		switch resp.StatusCode {
		case 200:
			return resp.transaction()
		}
		return nil, &invalidResponseError{resp}
	} else {
		return nil, &testOperationPerformedInProductionError{}
	}
}

// Void voids the transaction with the specified id if it has a status of authorized or
// submitted_for_settlement. When the transaction is voided Braintree will do an authorization
// reversal if possible so that the customer won’t have a pending charge on their card
func (g *TransactionGateway) Void(id string) (*Transaction, error) {
	resp, err := g.execute("PUT", "transactions/"+id+"/void", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// A transaction can be refunded if it is settled or settling.
// If the transaction has not yet begun settlement, use Void() instead.
// If you do not specify an amount to refund, the entire transaction amount will be refunded.
func (g *TransactionGateway) Refund(id string, amount ...*Decimal) (*Transaction, error) {
	var tx *Transaction
	if len(amount) > 0 {
		tx = &Transaction{
			Amount: amount[0],
		}
	}
	resp, err := g.execute("POST", "transactions/"+id+"/refund", tx)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	case 201:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// Find finds the transaction with the specified id.
func (g *TransactionGateway) Find(id string) (*Transaction, error) {
	resp, err := g.execute("GET", "transactions/"+id, nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// Search finds all transactions matching the search query.
func (g *TransactionGateway) Search(query *SearchQuery) (*TransactionSearchResult, error) {
	resp, err := g.execute("POST", "transactions/advanced_search_ids", query)
	if err != nil {
		return nil, err
	}
	var v TransactionSearchResult
	err = xml.Unmarshal(resp.Body, &v)
	if err != nil {
		return nil, err
	}
	return &v, err
}

type TransactionGetter struct {
	XMLName xml.Name `xml:"search"`
	Data    TrIDs    `xml:"ids"`
}

type TrIDs struct {
	Type string   `xml:"type,attr"`
	IDs  []string `xml:"item"`
}

type transactions50 struct {
	XMLName      xml.Name      `xml:"credit-card-transactions"`
	Transactions []Transaction `xml:"transaction"`
}

func (g *TransactionGateway) GetAll(ids []string) ([]Transaction, error) {
	var all []Transaction

	for i := 0; i < len(ids); i += 50 {
		top := i + 50
		if top > len(ids) {
			top = i + len(ids)%50
		}

		// Build xml
		tg := TransactionGetter{Data: TrIDs{Type: "array", IDs: ids[i:top]}}
		resp, err := g.execute("POST", "transactions/advanced_search", tg)
		if err != nil {
			return nil, err
		}

		var v transactions50
		err = xml.Unmarshal(resp.Body, &v)
		if err != nil {
			return nil, err
		}
		all = append(all, v.Transactions...)
	}
	return all, nil
}

type testOperationPerformedInProductionError struct {
	error
}

func (e *testOperationPerformedInProductionError) Error() string {
	return fmt.Sprint("Operation not allowed in production environment")
}

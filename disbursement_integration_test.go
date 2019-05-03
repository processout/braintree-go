package braintree

import (
	"testing"
)

// This test will fail unless you have a transaction with this ID on your sandbox.
func TestDisbursementTransactions(t *testing.T) {
	d := Disbursement{
		TransactionIds: []string{"dskdmb"},
	}

	result, err := d.Transactions(testGateway.Transaction())

	if err != nil {
		t.Fatal(err)
	}

	if len(result.TransactionIDs) != 1 {
		t.Fatal(result)
	}

	txn := result.TransactionIDs[0]
	if txn != "dskdmb" {
		t.Fatal(txn)
	}

}

package braintree

import (
	"math/rand"
	"testing"

	"github.com/processout/braintree-go/testhelpers"
)

func randomAmount() *Decimal {
	return NewDecimal(rand.Int63n(10000), 2)
}

func TestTransactionCreateSubmitForSettlementAndVoid(t *testing.T) {
	tx, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: NewDecimal(2000, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})

	t.Log(tx)

	if err != nil {
		t.Fatal(err)
	}
	if tx.Id == "" {
		t.Fatal("Received invalid ID on new transaction")
	}
	if tx.Status != "authorized" {
		t.Fatal(tx.Status)
	}

	// Submit for settlement
	ten := NewDecimal(1000, 2)
	tx2, err := testGateway.Transaction().SubmitForSettlement(tx.Id, ten)

	t.Log(tx2)

	if err != nil {
		t.Fatal(err)
	}
	if x := tx2.Status; x != "submitted_for_settlement" {
		t.Fatal(x)
	}
	if amount := tx2.Amount; amount.Cmp(ten) != 0 {
		t.Fatalf("transaction settlement amount (%s) did not equal amount requested (%s)", amount, ten)
	}

	// Void
	tx3, err := testGateway.Transaction().Void(tx2.Id)

	t.Log(tx3)

	if err != nil {
		t.Fatal(err)
	}
	if x := tx3.Status; x != "voided" {
		t.Fatal(x)
	}
}

func TestExternal3DSTransactionCreateSubmitForSettlementAndVoid(t *testing.T) {
	tx, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: NewDecimal(2000, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
		ThreeDSecurePassThru: &ThreeDSecurePassThrough{
			CAVV:                "200",
			DSTransactionID:     "",
			ECIFlag:             "02",
			ThreeDSecureVersion: "1.0.2",
			XID:                 "123",
		},
	})

	t.Log(tx)

	if err != nil {
		t.Fatal(err)
	}
	if tx.Id == "" {
		t.Fatal("Received invalid ID on new transaction")
	}
	if tx.Status != "authorized" {
		t.Fatal(tx.Status)
	}

	// Submit for settlement
	ten := NewDecimal(1000, 2)
	tx2, err := testGateway.Transaction().SubmitForSettlement(tx.Id, ten)

	t.Log(tx2)

	if err != nil {
		t.Fatal(err)
	}
	if x := tx2.Status; x != "submitted_for_settlement" {
		t.Fatal(x)
	}
	if amount := tx2.Amount; amount.Cmp(ten) != 0 {
		t.Fatalf("transaction settlement amount (%s) did not equal amount requested (%s)", amount, ten)
	}

	// Void
	tx3, err := testGateway.Transaction().Void(tx2.Id)

	t.Log(tx3)

	if err != nil {
		t.Fatal(err)
	}
	if x := tx3.Status; x != "voided" {
		t.Fatal(x)
	}
}

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestTransactionCreateSubmitForSettlementAndVoidWithNetworkToken(t *testing.T) {
	t.Skip("Requires specifying MerchantID which accepts Network Tokens")
	tx, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: NewDecimal(2000, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
			NetworkTokenAttributes: &NetworkTokenizationAttributes{
				Cryptogram:         "473dcfbc-06b8-431b-9885-ad1e904f8b1b",
				EcommerceIndicator: "5a4dacb5-03a4-4c8d-bf23-c11a5433c036",
				TokenRequestorID:   testGateway.MerchantId,
			},
		},
	})

	t.Log(tx)

	if err != nil {
		t.Fatal(err)
	}
	if tx.Id == "" {
		t.Fatal("Received invalid ID on new transaction")
	}
	if tx.Status != "authorized" {
		t.Fatal(tx.Status)
	}

	// Submit for settlement
	ten := NewDecimal(1000, 2)
	tx2, err := testGateway.Transaction().SubmitForSettlement(tx.Id, ten)

	t.Log(tx2)

	if err != nil {
		t.Fatal(err)
	}
	if x := tx2.Status; x != "submitted_for_settlement" {
		t.Fatal(x)
	}
	if amount := tx2.Amount; amount.Cmp(ten) != 0 {
		t.Fatalf("transaction settlement amount (%s) did not equal amount requested (%s)", amount, ten)
	}

	// Void
	tx3, err := testGateway.Transaction().Void(tx2.Id)

	t.Log(tx3)

	if err != nil {
		t.Fatal(err)
	}
	if x := tx3.Status; x != "voided" {
		t.Fatal(x)
	}
}

func TestTransactionSearch(t *testing.T) {
	txg := testGateway.Transaction()
	createTx := func(amount *Decimal, customerName string) error {
		_, err := txg.Create(&Transaction{
			Type:   "sale",
			Amount: amount,
			Customer: &Customer{
				FirstName: customerName,
			},
			CreditCard: &CreditCard{
				Number:         testCreditCards["visa"].Number,
				ExpirationDate: "05/14",
			},
		})
		return err
	}

	unique := testhelpers.RandomString()

	name0 := "Erik-" + unique
	if err := createTx(randomAmount(), name0); err != nil {
		t.Fatal(err)
	}

	name1 := "Lionel-" + unique
	if err := createTx(randomAmount(), name1); err != nil {
		t.Fatal(err)
	}

	query := new(SearchQuery)
	f := query.AddTextField("customer-first-name")
	f.Is = strp(name0)

	result, err := txg.Search(query)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.TransactionIDs) != 1 {
		t.Fatal(result.TransactionIDs)
	}
}

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestTransactionCreateWhenGatewayRejected(t *testing.T) {
	_, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: NewDecimal(201000, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err == nil {
		t.Fatal("Did not receive error when creating invalid transaction")
	}
	if err.Error() != "Card Issuer Declined CVV" {
		t.Fatal(err)
	}
}

func TestFindTransaction(t *testing.T) {
	createdTransaction, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["mastercard"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	foundTransaction, err := testGateway.Transaction().Find(createdTransaction.Id)
	if err != nil {
		t.Fatal(err)
	}

	if createdTransaction.Id != foundTransaction.Id {
		t.Fatal("transaction ids do not match")
	}
}

func TestFindNonExistantTransaction(t *testing.T) {
	_, err := testGateway.Transaction().Find("bad_transaction_id")
	if err == nil {
		t.Fatal("Did not receive error when finding an invalid tx ID")
	}
	if err.Error() != "404 Not Found" {
		t.Fatal(err)
	}
}

func TestAllTransactionFields(t *testing.T) {
	tx := &Transaction{
		Type:    "sale",
		Amount:  randomAmount(),
		OrderId: "my_custom_order",
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
			CVV:            "100",
		},
		Customer: &Customer{
			FirstName: "Lionel",
		},
		BillingAddress: &Address{
			StreetAddress: "1 E Main St",
			Locality:      "Chicago",
			Region:        "IL",
			PostalCode:    "60637",
		},
		ShippingAddress: &Address{
			StreetAddress: "1 E Main St",
			Locality:      "Chicago",
			Region:        "IL",
			PostalCode:    "60637",
		},
		Options: &TransactionOptions{
			SubmitForSettlement:              true,
			AddBillingAddressToPaymentMethod: true,
			StoreShippingAddressInVault:      true,
		},
	}

	tx2, err := testGateway.Transaction().Create(tx)
	if err != nil {
		t.Fatal(err)
	}

	if tx2.Type != tx.Type {
		t.Fatalf("expected Type to be equal, but %s was not %s", tx2.Type, tx.Type)
	}
	if tx2.Amount.Cmp(tx.Amount) != 0 {
		t.Fatalf("expected Amount to be equal, but %s was not %s", tx2.Amount, tx.Amount)
	}
	if tx2.OrderId != tx.OrderId {
		t.Fatalf("expected OrderId to be equal, but %s was not %s", tx2.OrderId, tx.OrderId)
	}
	if tx2.Customer.FirstName != tx.Customer.FirstName {
		t.Fatalf("expected Customer.FirstName to be equal, but %s was not %s", tx2.Customer.FirstName, tx.Customer.FirstName)
	}
	if tx2.BillingAddress.StreetAddress != tx.BillingAddress.StreetAddress {
		t.Fatalf("expected BillingAddress.StreetAddress to be equal, but %s was not %s", tx2.BillingAddress.StreetAddress, tx.BillingAddress.StreetAddress)
	}
	if tx2.BillingAddress.Locality != tx.BillingAddress.Locality {
		t.Fatalf("expected BillingAddress.Locality to be equal, but %s was not %s", tx2.BillingAddress.Locality, tx.BillingAddress.Locality)
	}
	if tx2.BillingAddress.Region != tx.BillingAddress.Region {
		t.Fatalf("expected BillingAddress.Region to be equal, but %s was not %s", tx2.BillingAddress.Region, tx.BillingAddress.Region)
	}
	if tx2.BillingAddress.PostalCode != tx.BillingAddress.PostalCode {
		t.Fatalf("expected BillingAddress.PostalCode to be equal, but %s was not %s", tx2.BillingAddress.PostalCode, tx.BillingAddress.PostalCode)
	}
	if tx2.ShippingAddress.StreetAddress != tx.ShippingAddress.StreetAddress {
		t.Fatalf("expected ShippingAddress.StreetAddress to be equal, but %s was not %s", tx2.ShippingAddress.StreetAddress, tx.ShippingAddress.StreetAddress)
	}
	if tx2.ShippingAddress.Locality != tx.ShippingAddress.Locality {
		t.Fatalf("expected ShippingAddress.Locality to be equal, but %s was not %s", tx2.ShippingAddress.Locality, tx.ShippingAddress.Locality)
	}
	if tx2.ShippingAddress.Region != tx.ShippingAddress.Region {
		t.Fatalf("expected ShippingAddress.Region to be equal, but %s was not %s", tx2.ShippingAddress.Region, tx.ShippingAddress.Region)
	}
	if tx2.ShippingAddress.PostalCode != tx.ShippingAddress.PostalCode {
		t.Fatalf("expected ShippingAddress.PostalCode to be equal, but %s was not %s", tx2.ShippingAddress.PostalCode, tx.ShippingAddress.PostalCode)
	}
	if tx2.CreditCard.Token == "" {
		t.Fatalf("expected CreditCard.Token to be equal, but %s was not %s", tx2.CreditCard.Token, tx.CreditCard.Token)
	}
	if tx2.Customer.Id == "" {
		t.Fatalf("expected Customer.Id to be equal, but %s was not %s", tx2.Customer.Id, tx.Customer.Id)
	}
	if tx2.Status != "submitted_for_settlement" {
		t.Fatalf("expected tx2.Status to be %s, but got %s", "submitted_for_settlement", tx2.Status)
	}
}

// This test will only pass on Travis. See TESTING.md for more details.
func TestTransactionDisbursementDetails(t *testing.T) {
	txn, err := testGateway.Transaction().Find("dskdmb")
	if err != nil {
		t.Fatal(err)
	}

	if txn.DisbursementDetails.DisbursementDate != "2013-06-27" {
		t.Fatalf("expected disbursement date to be %s, was %s", "2013-06-27", txn.DisbursementDetails.DisbursementDate)
	}
	if txn.DisbursementDetails.SettlementAmount.Cmp(NewDecimal(10000, 2)) != 0 {
		t.Fatalf("expected settlement amount to be %s, was %s", NewDecimal(10000, 2), txn.DisbursementDetails.SettlementAmount)
	}
	if txn.DisbursementDetails.SettlementCurrencyIsoCode != "USD" {
		t.Fatalf("expected settlement currency code to be %s, was %s", "USD", txn.DisbursementDetails.SettlementCurrencyIsoCode)
	}
	if txn.DisbursementDetails.SettlementCurrencyExchangeRate.Cmp(NewDecimal(100, 2)) != 0 {
		t.Fatalf("expected settlement currency exchange rate to be %s, was %s", NewDecimal(100, 2), txn.DisbursementDetails.SettlementCurrencyExchangeRate)
	}
	if !txn.DisbursementDetails.FundsHeld.Valid || txn.DisbursementDetails.FundsHeld.Bool {
		t.Error("funds held doesn't match")
	}
	if !txn.DisbursementDetails.Success.Valid || !txn.DisbursementDetails.Success.Bool {
		t.Error("success doesn't match")
	}
}

func TestTransactionCreateFromPaymentMethodCode(t *testing.T) {
	customer, err := testGateway.Customer().Create(&Customer{
		CreditCard: &CreditCard{
			Number:         testCreditCards["discover"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if customer.CreditCards.CreditCard[0].Token == "" {
		t.Fatal("invalid token")
	}

	tx, err := testGateway.Transaction().Create(&Transaction{
		Type:               "sale",
		CustomerID:         customer.Id,
		Amount:             randomAmount(),
		PaymentMethodToken: customer.CreditCards.CreditCard[0].Token,
	})
	if err != nil {
		t.Fatal(err)
	}
	if tx.Id == "" {
		t.Fatal("invalid tx id")
	}
}

func TestSettleTransaction(t *testing.T) {
	old_environment := testGateway.Environment

	txn, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().SubmitForSettlement(txn.Id, txn.Amount)
	if err != nil {
		t.Fatal(err)
	}

	testGateway.Environment = Production

	_, err = testGateway.Transaction().Settle(txn.Id)
	if err.Error() != "Operation not allowed in production environment" {
		t.Log(testGateway.Environment)
		t.Fatal(err)
	}

	testGateway.Environment = old_environment

	txn, err = testGateway.Transaction().Settle(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Status != "settled" {
		t.Fatal(txn.Status)
	}
}

func TestTrxPaymentMethodNonce(t *testing.T) {
	txn, err := testGateway.Transaction().Create(&Transaction{
		Type:               "sale",
		Amount:             randomAmount(),
		PaymentMethodNonce: "fake-apple-pay-mastercard-nonce",
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().SubmitForSettlement(txn.Id, txn.Amount)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransactionCreateSettleAndFullRefund(t *testing.T) {
	amount := NewDecimal(20000, 2)
	txn, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: amount,
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().SubmitForSettlement(txn.Id, txn.Amount)
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().Settle(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Status != "settled" {
		t.Fatal(txn.Status)
	}

	// Refund
	refundTxn, err := testGateway.Transaction().Refund(txn.Id)

	t.Log(refundTxn)

	if err != nil {
		t.Fatal(err)
	}
	if x := refundTxn.Status; x != "submitted_for_settlement" {
		t.Fatal(x)
	}

	refundTxn, err = testGateway.Transaction().Settle(refundTxn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if refundTxn.Status != "settled" {
		t.Fatal(txn.Status)
	}

	if *refundTxn.RefundedTransactionId != txn.Id {
		t.Fatal(*refundTxn.RefundedTransactionId)
	}

	// Check that the refund shows up in the original transaction
	txn, err = testGateway.Transaction().Find(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.RefundIds != nil && (*txn.RefundIds).IDs[0] != refundTxn.Id {
		t.Fatal(*txn.RefundIds)
	}

	// Second refund should fail
	refundTxn, err = testGateway.Transaction().Refund(txn.Id)
	t.Log(refundTxn)

	if err.Error() != "Transaction has already been fully refunded." {
		t.Fatal(err)
	}
}

func TestTransactionCreateSettleAndPartialRefund(t *testing.T) {
	amount := NewDecimal(10000, 2)
	refundAmt1 := NewDecimal(5000, 2)
	refundAmt2 := NewDecimal(5001, 2)
	txn, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: amount,
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().SubmitForSettlement(txn.Id, txn.Amount)
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().Settle(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Status != "settled" {
		t.Fatal(txn.Status)
	}

	// Refund
	refundTxn, err := testGateway.Transaction().Refund(txn.Id, refundAmt1)

	t.Log(refundTxn)

	if err != nil {
		t.Fatal(err)
	}
	if x := refundTxn.Status; x != "submitted_for_settlement" {
		t.Fatal(x)
	}

	refundTxn, err = testGateway.Transaction().Settle(refundTxn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if refundTxn.Status != "settled" {
		t.Fatal(txn.Status)
	}

	// Refund amount too large
	refundTxn, err = testGateway.Transaction().Refund(txn.Id, refundAmt2)

	t.Log(refundTxn)

	if err.Error() != "Refund amount is too large." {
		t.Fatal(err)
	}
}

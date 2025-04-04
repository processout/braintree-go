package braintree

import (
	"encoding/xml"
	"testing"

	"github.com/processout/braintree-go/testhelpers"
)

var acctId string

func TestMerchantAccountCreate(t *testing.T) {
	// TODO: Fix integration - Create ends with Business tax ID or Individual SSN is required ✓
	t.Skip()
	acctId = testhelpers.RandomString()
	acct := MerchantAccount{
		MasterMerchantAccountId: testMerchantAccountId,
		TOSAccepted:             true,
		Id:                      acctId,
		Individual: &MerchantAccountPerson{
			FirstName:   "Kayle",
			LastName:    "Gishen",
			Email:       "kayle.gishen@example.com",
			Phone:       "5556789012",
			DateOfBirth: "1-1-1989",
			Address: &Address{
				StreetAddress:   "1 E Main St",
				ExtendedAddress: "Suite 404",
				Locality:        "Chicago",
				Region:          "IL",
				PostalCode:      "60622",
			},
		},
		FundingOptions: &MerchantAccountFundingOptions{
			Destination: FUNDING_DEST_MOBILE_PHONE,
			MobilePhone: "5552344567",
		},
	}

	x, _ := xml.Marshal(&acct)
	t.Log(string(x))

	// TODO: Fix integration - Create ends with Business tax ID or Individual SSN is required ✓
	merchantAccount, err := testGateway.MerchantAccount().Create(&acct)

	t.Log(merchantAccount)

	if err != nil {
		t.Fatal(err)
	}

	if merchantAccount.Id == "" {
		t.Fatal("invalid merchant account id")
	}

	ma2, err := testGateway.MerchantAccount().Find(merchantAccount.Id)

	t.Log(ma2)

	if err != nil {
		t.Fatal(err)
	}

	if ma2.Id != merchantAccount.Id {
		t.Fatal("ids do not match")
	}
}

func TestMerchantAccountTransaction(t *testing.T) {
	// TODO: Fix integration - Create ends with 403 Forbidden ✓
	t.Skip()
	if acctId == "" {
		TestMerchantAccountCreate(t)
	}

	amount := NewDecimal(int64(randomAmount().Scale+500), 2)

	// TODO: Fix integration - Create ends with 403 Forbidden ✓
	tx, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: amount,
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
		ServiceFeeAmount:  NewDecimal(500, 2),
		MerchantAccountId: acctId,
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
}

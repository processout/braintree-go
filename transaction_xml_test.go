package braintree

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestTransactionMarshalsExternalVaultTLID(t *testing.T) {
	tx := &Transaction{
		Type: "sale",
		ExternalVault: &ExternalVault{
			PreviousNetworkTransactionID:        "123456789012345",
			PreviousMastercardTransactionLinkID: "JavmPovbDUHu8aIAyjexxG",
			Status:                              ExternalVaultStatusVaulted,
		},
	}

	b, err := xml.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}

	want := "<external-vault>" +
		"<previous-network-transaction-id>123456789012345</previous-network-transaction-id>" +
		"<previous-mastercard-transaction-link-id>JavmPovbDUHu8aIAyjexxG</previous-mastercard-transaction-link-id>" +
		"<status>vaulted</status>" +
		"</external-vault>"
	if !strings.Contains(string(b), want) {
		t.Fatalf("expected %s in %s", want, b)
	}
}

func TestTransactionMarshalsWillVaultWithoutPreviousIDs(t *testing.T) {
	tx := &Transaction{
		Type:          "sale",
		ExternalVault: &ExternalVault{Status: ExternalVaultStatusWillVault},
	}

	b, err := xml.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}

	want := "<external-vault><status>will_vault</status></external-vault>"
	if !strings.Contains(string(b), want) {
		t.Fatalf("expected %s in %s", want, b)
	}
}

func TestTransactionUnmarshalsTLID(t *testing.T) {
	payload := `<transaction>
		<network-transaction-id>020261102133130</network-transaction-id>
		<mastercard-transaction-link-id>JavmPovbDUHu8aIAyjexxG</mastercard-transaction-link-id>
	</transaction>`

	var tx Transaction
	if err := xml.Unmarshal([]byte(payload), &tx); err != nil {
		t.Fatal(err)
	}

	if tx.NetworkTransactionID != "020261102133130" {
		t.Errorf("unexpected network transaction id %q", tx.NetworkTransactionID)
	}
	if tx.MastercardTransactionLinkID != "JavmPovbDUHu8aIAyjexxG" {
		t.Errorf("unexpected mastercard transaction link id %q", tx.MastercardTransactionLinkID)
	}
}

func TestTransactionUnmarshalsNilTLID(t *testing.T) {
	payload := `<transaction><mastercard-transaction-link-id nil="true"/></transaction>`

	var tx Transaction
	if err := xml.Unmarshal([]byte(payload), &tx); err != nil {
		t.Fatal(err)
	}

	if tx.MastercardTransactionLinkID != "" {
		t.Errorf("expected empty mastercard transaction link id, got %q", tx.MastercardTransactionLinkID)
	}
}

func TestCreditCardVerificationUnmarshalsTLID(t *testing.T) {
	payload := `<verification><mastercard-transaction-link-id>JavmPovbDUHu8aIAyjexxG</mastercard-transaction-link-id></verification>`

	var v CreditCardVerification
	if err := xml.Unmarshal([]byte(payload), &v); err != nil {
		t.Fatal(err)
	}

	if v.MastercardTransactionLinkID != "JavmPovbDUHu8aIAyjexxG" {
		t.Errorf("unexpected mastercard transaction link id %q", v.MastercardTransactionLinkID)
	}
}

func TestTransactionUnmarshalsMerchantAdviceCode(t *testing.T) {
	payload := `<transaction>
		<processor-response-code>2043</processor-response-code>
		<merchant-advice-code>03</merchant-advice-code>
		<merchant-advice-code-text>Do Not Try Again</merchant-advice-code-text>
	</transaction>`

	var tx Transaction
	if err := xml.Unmarshal([]byte(payload), &tx); err != nil {
		t.Fatal(err)
	}

	if tx.MerchantAdviceCode != "03" {
		t.Fatalf("expected merchant advice code 03, got %q", tx.MerchantAdviceCode)
	}
	if tx.MerchantAdviceCodeText != "Do Not Try Again" {
		t.Fatalf("expected merchant advice code text, got %q", tx.MerchantAdviceCodeText)
	}
}

func TestBraintreeErrorUnmarshalsMerchantAdviceCode(t *testing.T) {
	payload := `<api-error-response>
		<message>Error - Do Not Retry, Call Issuer</message>
		<transaction>
			<processor-response-code>2043</processor-response-code>
			<merchant-advice-code>03</merchant-advice-code>
		</transaction>
	</api-error-response>`

	var btErr BraintreeError
	if err := xml.Unmarshal([]byte(payload), &btErr); err != nil {
		t.Fatal(err)
	}

	if btErr.Transaction.MerchantAdviceCode != "03" {
		t.Fatalf("expected merchant advice code 03, got %q", btErr.Transaction.MerchantAdviceCode)
	}
}

func TestTransactionUnmarshalsGatewayRejectionAndSettlementResponse(t *testing.T) {
	payload := `<transaction>
		<status>gateway_rejected</status>
		<gateway-rejection-reason>cvv</gateway-rejection-reason>
		<processor-settlement-response-code>4001</processor-settlement-response-code>
		<processor-settlement-response-text>Settlement Declined</processor-settlement-response-text>
	</transaction>`

	var tx Transaction
	if err := xml.Unmarshal([]byte(payload), &tx); err != nil {
		t.Fatal(err)
	}

	if tx.GatewayRejectionReason != "cvv" {
		t.Fatalf("expected gateway rejection reason cvv, got %q", tx.GatewayRejectionReason)
	}
	if tx.ProcessorSettlementResponseCode != "4001" {
		t.Fatalf("expected settlement response code 4001, got %q", tx.ProcessorSettlementResponseCode)
	}
	if tx.ProcessorSettlementResponseText != "Settlement Declined" {
		t.Fatalf("expected settlement response text, got %q", tx.ProcessorSettlementResponseText)
	}
}

func TestTransactionUnmarshalsNetworkTokenAndNetworkResponse(t *testing.T) {
	payload := `<transaction>
		<processor-response-type>hard_declined</processor-response-type>
		<additional-processor-response>05 : NOT AUTHORISED</additional-processor-response>
		<network-response-code>05</network-response-code>
		<network-response-text>Do not honor</network-response-text>
		<processed-with-network-token type="boolean">true</processed-with-network-token>
		<credit-card>
			<is-network-tokenized type="boolean">false</is-network-tokenized>
			<payment-account-reference>V0010013022298169667504231315</payment-account-reference>
		</credit-card>
		<network-token>
			<bin>411111</bin>
			<is-network-tokenized type="boolean">true</is-network-tokenized>
			<payment-account-reference>V0010013022298169667504231315</payment-account-reference>
		</network-token>
		<three-d-secure-info>
			<liability-shifted type="boolean">true</liability-shifted>
			<liability-shift-possible nil="true"></liability-shift-possible>
			<status>authenticate_successful</status>
			<eci-flag>05</eci-flag>
			<three-d-secure-version>2.2.0</three-d-secure-version>
			<ds-transaction-id>ds-trans-id</ds-transaction-id>
		</three-d-secure-info>
	</transaction>`

	var tx Transaction
	if err := xml.Unmarshal([]byte(payload), &tx); err != nil {
		t.Fatal(err)
	}

	if tx.ProcessorResponseType != "hard_declined" || tx.AdditionalProcessorResponse != "05 : NOT AUTHORISED" {
		t.Fatalf("unexpected processor response: %q %q", tx.ProcessorResponseType, tx.AdditionalProcessorResponse)
	}
	if tx.NetworkResponseCode != "05" || tx.NetworkResponseText != "Do not honor" {
		t.Fatalf("unexpected network response: %q %q", tx.NetworkResponseCode, tx.NetworkResponseText)
	}
	if tx.ProcessedWithNetworkToken == nil || !tx.ProcessedWithNetworkToken.Valid || !tx.ProcessedWithNetworkToken.Bool {
		t.Fatalf("expected processed with network token, got %+v", tx.ProcessedWithNetworkToken)
	}
	if tx.CreditCard.PaymentAccountReference != "V0010013022298169667504231315" || tx.CreditCard.IsNetworkTokenized.Bool {
		t.Fatalf("unexpected credit card details: %+v", tx.CreditCard)
	}
	if tx.NetworkToken == nil || !tx.NetworkToken.IsNetworkTokenized.Bool || tx.NetworkToken.Bin != "411111" {
		t.Fatalf("unexpected network token details: %+v", tx.NetworkToken)
	}
	info := tx.ThreeDSecureInfo
	if info == nil || !info.LiabilityShifted.Bool || info.LiabilityShiftPossible.Valid || info.ECIFlag != "05" ||
		info.ThreeDSecureVersion != "2.2.0" || info.DSTransactionID != "ds-trans-id" || info.Status != "authenticate_successful" {
		t.Fatalf("unexpected three d secure info: %+v", info)
	}
}

func TestTransactionRequestOmitsResponseOnlyFields(t *testing.T) {
	tx := &Transaction{Type: "sale", CreditCard: &CreditCard{Number: "4111111111111111"}}

	b, err := xml.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}

	for _, field := range []string{"network-token", "processed-with-network-token", "three-d-secure-info",
		"is-network-tokenized", "payment-account-reference", "network-response-code"} {
		if strings.Contains(string(b), "<"+field) {
			t.Fatalf("did not expect %s in %s", field, b)
		}
	}
}

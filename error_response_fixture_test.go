package braintree

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadErrorResponseFixture(t *testing.T, name string) *BraintreeError {
	t.Helper()

	data, err := os.ReadFile("testdata/" + name)
	require.NoError(t, err)

	var btErr BraintreeError
	require.NoError(t, xml.Unmarshal(data, &btErr))
	return &btErr
}

func TestErrorResponseFixture_ProcessorDeclinedNetworkToken(t *testing.T) {
	btErr := loadErrorResponseFixture(t, "processor_declined_network_token.xml")
	tx := btErr.Transaction

	assert.Equal(t, "Error - Do Not Retry, Call Issuer", btErr.ErrorMessage)
	assert.Equal(t, "processor_declined", tx.Status)
	assert.Equal(t, "2043", tx.ProcessorResponseCode)
	assert.Equal(t, "hard_declined", tx.ProcessorResponseType)
	assert.Equal(t, "31 : Network Decline - ErrorCode: 712 - ERROR", tx.AdditionalProcessorResponse)
	assert.Empty(t, tx.NetworkResponseCode)
	assert.Empty(t, tx.MerchantAdviceCode)
	assert.Empty(t, tx.GatewayRejectionReason)
	assert.Empty(t, tx.ProcessorSettlementResponseCode)

	require.NotNil(t, tx.ProcessedWithNetworkToken)
	assert.True(t, tx.ProcessedWithNetworkToken.Valid && tx.ProcessedWithNetworkToken.Bool)

	require.NotNil(t, tx.CreditCard)
	assert.False(t, tx.CreditCard.IsNetworkTokenized.Bool)
	assert.Empty(t, tx.CreditCard.PaymentAccountReference)

	require.NotNil(t, tx.NetworkToken)
	assert.True(t, tx.NetworkToken.IsNetworkTokenized.Bool)
	assert.Equal(t, "555555", tx.NetworkToken.Bin)
	assert.Equal(t, "MasterCard", tx.NetworkToken.CardType)

	require.NotNil(t, tx.ThreeDSecureInfo)
	assert.Nil(t, tx.ThreeDSecureInfo.LiabilityShifted)
	assert.Empty(t, tx.ThreeDSecureInfo.ECIFlag)

	assert.Equal(t, "TESTNTI000001  ", tx.NetworkTransactionID)
}

func TestErrorResponseFixture_GatewayRejectedCVV(t *testing.T) {
	btErr := loadErrorResponseFixture(t, "gateway_rejected_cvv.xml")
	tx := btErr.Transaction

	assert.Equal(t, "Gateway Rejected: cvv", btErr.ErrorMessage)
	assert.Equal(t, "gateway_rejected", tx.Status)
	assert.Equal(t, "cvv", tx.GatewayRejectionReason)
	assert.Equal(t, "1000", tx.ProcessorResponseCode)
	assert.Equal(t, "approved", tx.ProcessorResponseType)
	assert.Equal(t, "TEST01", tx.ProcessorAuthorizationCode)
	assert.Equal(t, "N", tx.CVVResponseCode)
	assert.Equal(t, "00", tx.NetworkResponseCode)
	assert.Equal(t, "Successful approval/completion or V.I.P. PIN verification is successful", tx.NetworkResponseText)
	assert.Empty(t, tx.AdditionalProcessorResponse)
	assert.Empty(t, tx.MerchantAdviceCode)
	assert.Nil(t, tx.NetworkToken)

	require.NotNil(t, tx.ProcessedWithNetworkToken)
	assert.True(t, tx.ProcessedWithNetworkToken.Valid)
	assert.False(t, tx.ProcessedWithNetworkToken.Bool)

	require.NotNil(t, tx.CreditCard)
	assert.Equal(t, "411111", tx.CreditCard.Bin)
	assert.False(t, tx.CreditCard.IsNetworkTokenized.Bool)

	assert.Equal(t, "000000000000002", tx.NetworkTransactionID)
}

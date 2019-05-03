package braintree

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime_UnmarshalXML(t *testing.T) {
	xmlNotif := `<?xml version="1.0" encoding="UTF-8"?>
	<notification>
		<kind>dispute_opened</kind>
		<timestamp type="datetime">2019-03-14T13:23:50Z</timestamp>
		<subject>
			<dispute>
				<id>asdf</id>
				<amount>110.60</amount>
				<amount-disputed>110.60</amount-disputed>
				<amount-won>0.00</amount-won>
				<case-number>bfagbafd</case-number>
				<created-at type="datetime">2019-03-14T13:23:40Z</created-at>
				<currency-iso-code>USD</currency-iso-code>
				<date-opened type="date">2019-03-14</date-opened>
				<date-won nil="true"/>
				<processor-comments nil="true"/>
				<kind>chargeback</kind>
				<merchant-account-id>asdfasdf</merchant-account-id>
				<reason>fraud</reason>
				<reason-code>7030</reason-code>
				<reason-description>UA02 - Fraudulent Transaction-Card Not Present</reason-description>
				<received-date type="date">2019-03-14</received-date>
				<reference-number>134534513451</reference-number>
				<reply-by-date type="date">2019-03-24</reply-by-date>
				<status>open</status>
				<updated-at type="datetime">2019-03-14T13:23:41Z</updated-at>
				<original-dispute-id nil="true"/>
				<evidence type="array"/>
				<status-history type="array">
					<status-history>
						<disbursement-date nil="true"/>
						<effective-date type="date">2019-03-14</effective-date>
						<status>open</status>
						<timestamp type="datetime">2019-03-14T13:23:41Z</timestamp>
					</status-history>
					<status-history>
						<disbursement-date type="date">2019-03-14</disbursement-date>
						<effective-date type="date">2019-03-14</effective-date>
						<status>pending</status>
						<timestamp type="datetime">2019-03-14T13:23:41Z</timestamp>
					</status-history>
				</status-history>
				<transaction>
					<id>adsgafg</id>
					<amount>110.60</amount>
					<created-at>2019-03-07T16:56:14Z</created-at>
					<order-id nil="true"/>
					<purchase-order-number nil="true"/>
					<payment-instrument-subtype>Discover</payment-instrument-subtype>
				</transaction>
			</dispute>
		</subject>
	</notification>
	`
	var n WebhookNotification
	err := xml.Unmarshal([]byte(xmlNotif), &n)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, time.Date(2019, 03, 14, 0, 0, 0, 0, time.UTC), time.Time(n.Subject.Dispute.DateOpened))
	assert.Equal(t, time.Time{}, time.Time(n.Subject.Dispute.DateWon))
}

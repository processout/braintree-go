package braintree

import (
	"encoding/xml"
	"time"
)

// Time is an lias to time.Time for custom xml unmarshaling
type Time time.Time

// UnmarshalXML is a custom unmarshaler that checks for two formats
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "2006-01-02" // yyyymmdd date format
	var v string
	d.DecodeElement(&v, &start)

	if v == "" {
		return nil
	}

	parse, err := time.Parse(time.RFC3339, v)
	if err != nil {
		parse, err = time.Parse(shortForm, v)
		if err != nil {
			return err
		}
	}

	*t = Time(parse)
	return nil
}

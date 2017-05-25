package sdt

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
)

type parcelRpost struct {
	NalScheme int        `xml:"nal_scheme,attr"`
	MailType  int        `xml:"mail_type,attr"`
	PostMark  int        `xml:"post_mark,attr"`
	Zip       string     `xml:"zip,attr"`
	SumVl     float64    `xml:"sum_vl,attr"`
	NalPlat   float64    `xml:"nal_plat,attr"`
	Weight    float64    `xml:"weight,attr"`
	PostDate  string     `xml:"post_date,attr"`
}

type rpost struct {
	Parcel struct {
			   NalPlat  string    `xml:"nal_plat,attr"`
			   AirRate  string    `xml:"air_rate,attr"`
			   SumVl    string    `xml:"sum_vl,attr"`
			   InstRate string    `xml:"inst_rate,attr"`
			   State    string    `xml:"state,attr"`
			   PostMark string    `xml:"post_mark,attr"`
			   MassRate string    `xml:"mass_rate,attr"`
			   MailType string    `xml:"mail_type,attr"`
			   MailCtg  string    `xml:"mail_ctg,attr"`
		   } `xml:"parcel"`
}

func (rp *rpost)decodeXml(r apiResult) error {
	decoder := xml.NewDecoder(r.Body)
	decoder.CharsetReader = charset.NewReaderLabel

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch element := token.(type) {
		case xml.StartElement:
			inElement := element.Name.Local
			if inElement == "response" {
				if err := decoder.DecodeElement(rp, &element); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
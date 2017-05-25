package sdt

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
)

type parcelCourier struct {
	CityFrom     string  `xml:"city_from,attr"`
	CityTo       string  `xml:"city_to,attr"`
	Weight       float64 `xml:"weight,attr"`
	Width        float64 `xml:"width,attr"`
	Height       float64 `xml:"height,attr"`
	Length       float64 `xml:"length,attr"`
	AssessedCost float64 `xml:"assessedCost,attr"`
	DeliveryType int     `xml:"delivery_type,attr"`
}

type couriers struct {
	Courier []courier
}

type courier struct {
	ServiceId int
	ToDoor    []courierToDoor  `xml:"deliveryToDoor>item"`
	ToPoint   []courierToPoint `xml:"deliveryToPoint>item"`
}

type courierToDoor struct {
	ProviderKey string `xml:"providerKey"`
	Tariffs     []struct {
		TariffName       string    `xml:"tariffName"`
		TariffId         string    `xml:"tariffId"`
		DeliveryCost     float64   `xml:"deliveryCost"`
		From             string    `xml:"from"`
		PickupTypes      string    `xml:"pickupTypes"`
		DaysMin          int       `xml:"daysMin"`
		TariffProviderId string    `xml:"tariffProviderId"`
		DaysMax          int       `xml:"daysMax"`
		DeliveryTypes    string    `xml:"deliveryTypes"`
	} `xml:"tariffs>item"`
}

type courierToPoint struct {
	ProviderKey string `xml:"providerKey"`
	Tariffs     []struct {
		TariffName       string    `xml:"tariffName"`
		TariffId         string    `xml:"tariffId"`
		DeliveryCost     float64   `xml:"deliveryCost"`
		From             string    `xml:"from"`
		PickupTypes      string    `xml:"pickupTypes"`
		DaysMin          int       `xml:"daysMin"`
		TariffProviderId string    `xml:"tariffProviderId"`
		DaysMax          int       `xml:"daysMax"`
		DeliveryTypes    string    `xml:"deliveryTypes"`
		PointIds         []int     `xml:"pointIds>text"`
	} `xml:"tariffs>item"`
}

func (cr *courier)decodeXml(r apiResult) error {
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
				if err := decoder.DecodeElement(cr, &element); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
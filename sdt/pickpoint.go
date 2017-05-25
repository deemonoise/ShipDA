package sdt

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
)

type Pickpoints struct {
	Pickpoint []pickpoint `xml:"rows>item" json:"result"`
	Total     int         `xml:"meta>total"`
}

type pickpoint struct {
	Url                string    `xml:"url" db:"url"`
	Area               string    `xml:"area" db:"area"`
	Street             string    `xml:"street" db:"street"`
	Description        string    `xml:"description" db:"description"`
	Office             string    `xml:"office" db:"office"`
	Email              string    `xml:"email" db:"email"`
	Cod                string    `xml:"cod" db:"cod"`
	AvailableOperation int       `xml:"availableOperation" db:"available_operation"`
	Id                 int       `xml:"id" db:"external_id"`
	PaymentCard        string    `xml:"paymentCard" db:"payment_card"`
	Code               string    `xml:"code" db:"code" json:"code"`
	CountryCode        string    `xml:"countryCode" db:"country_code"`
	Region             string    `xml:"region" db:"region"`
	Lat                float64   `xml:"lat" db:"lat" json:"coordinates>latitude"`
	Lng                float64   `xml:"lng" db:"lng" json:"coordinates>longitude"`
	City               string    `xml:"city" db:"city"`
	Timetable          string    `xml:"timetable" db:"timetable" json:"schedule"`
	Name               string    `xml:"name" db:"name" json:"name"`
	StreetType         string    `xml:"streetType" db:"street_type"`
	CityGuid           string    `xml:"cityGuid" db:"city_guid"`
	PostIndex          string    `xml:"postIndex" db:"post_index"`
	House              string    `xml:"house" db:"house"`
	ProviderKey        string    `xml:"providerKey" db:"provider_key"`
	Type               int       `xml:"type" db:"type"`
	Block              string    `xml:"block" db:"block"`
	Phone              string    `xml:"phone" db:"phone" json:"phone"`
	Address            string    `json:"address"`
}

type parcel struct {
	DeliveryType int    `xml:"delivery_type,attr"`
	City         string `xml:"city,attr"`
	Limit        int    `xml:"limit,attr"`
}


func (sdt *Sdt)GetPointsFromApi(city string, providerIds []int) (Pickpoints, error) {
	var points Pickpoints

	pr := parcel{
		City: city,
		Limit: 5000,
	}

	var parcels []parcel

	var rs []request

	if len(providerIds) > 0 {
		for _, p := range providerIds {
			pr.DeliveryType = p
			parcels = append(parcels, pr)
		}
	} else {
		parcels = append(parcels, pr)
	}

	for _, p := range parcels {
		rs = append(rs,
			request{
				Partner_id: sdt.PartnerId,
				Password: sdt.Password,
				RequestType: 57,
				ServiceId: p.DeliveryType,
				Parcel: p,
			})
	}

	req := requests{Request: rs}

	res, err := sdt.sendRequests(req)

	if err != nil {
		return points, err
	}

	for _, r := range res.Results {
		var p Pickpoints
		if r.Err == nil {
			err = p.decodeXml(r)
			if err != nil {
				return p, err
			}
		} else {
			return p, r.Err
		}
		for _, p := range p.Pickpoint {
			points.Pickpoint = append(points.Pickpoint, p)
		}
		points.Total += p.Total
	}
	return points, nil
}

func (points *Pickpoints)decodeXml(r apiResult) error {
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
				if err := decoder.DecodeElement(points, &element); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

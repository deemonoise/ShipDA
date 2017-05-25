package sdt

import "time"

type Parcel struct {
	CityFrom     string
	CityTo       string
	Zip          string
	Weight       float64
	Width        float64
	Height       float64
	Length       float64
	AssessedCost float64
}

func (sdt *Sdt)Calculate(parcel Parcel, providerIds []int) (couriers, rpost, error) {
	var crs couriers
	var rs []request
	var rp rpost

	if len(providerIds) > 0 {
		for _, p := range providerIds {
			pc := parcelCourier{
				DeliveryType: p,
				CityFrom: parcel.CityFrom,
				CityTo: parcel.CityTo,
				Weight: parcel.Weight,
				Width: parcel.Width,
				Height: parcel.Height,
				Length: parcel.Length,
				AssessedCost: parcel.AssessedCost,
			}
			rs = append(rs,
				request{
					Parcel: pc,
					Partner_id: sdt.PartnerId,
					Password: sdt.Password,
					RequestType: 54,
					ServiceId: p,
				},
			)
		}
	}

	if parcel.Zip != "" {
		tm := time.Now().Local()
		pc := parcelRpost{
			Zip: parcel.Zip,
			SumVl: parcel.AssessedCost,
			NalPlat: parcel.AssessedCost,
			Weight: float64(parcel.Weight) / 1000,
			PostDate: tm.Format("02.01.2006"),
			NalScheme: 0,
			MailType: 4,
			PostMark: 0,
		}
		rs = append(rs,
			request{
				Parcel: pc,
				Partner_id: sdt.PartnerId,
				Password: sdt.Password,
				RequestType: 53,
			},
		)
	}

	req := requests{Request: rs}

	res, err := sdt.sendRequests(req)
	if err != nil {
		return crs, rp, err
	}

	for _, r := range res.Results {
		switch r.Type {
		case 54:
			var cr courier
			err = cr.decodeXml(r)
			if err != nil {
				return crs, rp, err
			}
			cr.ServiceId = r.ServiceId
			crs.Courier = append(crs.Courier, cr)
		case 53:
			err = rp.decodeXml(r)
			if err != nil {
				return crs, rp, err
			}
		}
	}

	return crs, rp, nil
}
package main

import (
	"ShipDA/sdt"
	"log"
	"net/http"
)

type calcRequest struct {
	ClientId string    `json:"clientId"`
	Req      calculate `json:"calculate"`
}

type calculate struct {
	Packages        []pkg        `json:"packages"`        // Набор упаковок
	PayerType       string       `json:"payerType"`       // Плательщик за доставку (receiver или sender)
	DeclaredValue   float64      `json:"declaredValue"`   // Объявленная стоимость
	Cod             float64      `json:"cod"`             // Наложенный платеж
	DeliveryDate    string       `json:"deliveryDate"`    // Дата доставки
	DeliveryTime    deliveryTime `json:"deliveryTime"`    // Время доставки
	Currency        string       `json:"currency"`        // Код валюты
	DeliveryAddress address      `json:"deliveryAddress"` // Адрес доставки
	ShipmentAddress address      `json:"shipmentAddress"` // Адрес отгрузки
}

type address struct {
	Index        string `json:"index"`        // Индекс
	CountryIso   string `json:"countryIso"`   // ISO код страны (ISO 3166-1 alpha-2)
	Region       string `json:"region"`       // Регион
	RegionId     int    `json:"regionId"`     // Идентификатор региона в Geohelper
	City         string `json:"city"`         // Город
	CityId       int    `json:"cityId"`       // Идентификатор города в Geohelper
	CityType     string `json:"cityType"`     // Тип населенного пункта
	Street       string `json:"street"`       // Улица
	StreetId     int    `json:"streetId"`     // Идентификатор улицы в Geohelper
	StreetType   string `json:"streetType"`   // Тип улицы
	Building     string `json:"building"`     // Дом
	Flat         string `json:"flat"`         // Номер квартиры/офиса
	IntercomCode string `json:"intercomCode"` // Код домофона
	Floor        int    `json:"floor"`        // Этаж
	Block        int    `json:"block"`        // Подъезд
	House        string `json:"house"`        // Строение/корпус
	Metro        string `json:"metro"`        // Метро
	Notes        string `json:"notes"`        // Примечания к адресу
	Text         string `json:"text"`         // Адрес в текстовом виде
	Terminal     string `json:"terminal"`     // Код терминала отгрузки/доставки
}

type pkg struct {
	Weight float64 `json:"weight"` // вес в граммах
	Width  int     `json:"width"`  // ширина в мм
	Length int     `json:"length"` // длина в мм
	Height int     `json:"height"` // высота в мм
}

type deliveryTime struct {
	From   string `json:"from"`   // Время доставки "с"
	To     string `json:"to"`     // Время доставки "до"
	Custom string `json:"custom"` // Время доставки в свободной форме
}

type calcResponse struct {
	Success bool       `json:"success"`
	Result  calcResult `json:"result"`
}

type calcResult struct {
	Code            string      `json:"code"`            // Код тарифа
	Group           string      `json:"group"`           // Группа тарифов
	Name            string      `json:"name"`            // Наименование тарифа
	Type            string      `json:"type"`            // Тип тарифа (courier - курьерская доставка или selfDelivery - самовывоз)
	Description     string      `json:"description"`     // Описание
	Cost            float64     `json:"cost"`            // Стоимость доставки (Если не передана, то тариф будет выводиться, но не будет доступен для выбора)
	MinTerm         int         `json:"minTerm"`         // Минимальный срок доставки
	MaxTerm         int         `json:"maxTerm"`         // Максимальный срок доставки
	PickuppointList []pickpoint `json:"pickuppointList"` // Терминал отгрузки/получения
}

func (req *calcRequest) Calculate() (calcResponse, *appError) {
	var res calcResponse
	p := sdt.Parcel{
		CityFrom:     "Москва",
		CityTo:       req.Req.ShipmentAddress.City,
		Weight:       req.Req.Packages[0].Weight,
		Width:        float64(req.Req.Packages[0].Width / 10),
		Height:       float64(req.Req.Packages[0].Height / 10),
		Length:       float64(req.Req.Packages[0].Length / 10),
		AssessedCost: req.Req.Cod,
		Zip:          req.Req.ShipmentAddress.Index,
	}

	log.Printf("%#v", p)

	crs, rpost, err := Sdt.Calculate(p, []int{58, 62})
	if err != nil {
		e := NewAppError(http.StatusInternalServerError, err)
		return res, &e
	}

	log.Printf("%v", rpost)
	for _, c := range crs.Courier {
		log.Println("delivery_type:", c.ServiceId)
		//log.Printf("%v", c.ToDoor)
		//log.Printf("%v", c.ToPoint)

		for _, cr := range c.ToDoor {
			for _, t := range cr.Tariffs {
				cresult := calcResult{
					Name: t.TariffName,
					MinTerm: t.DaysMin,
					MaxTerm: t.DaysMax,
					Cost: t.DeliveryCost,
					Code: t.TariffId,
					Type: "courier",
					Group: cr.ProviderKey,
				}
				log.Printf("%#v", cresult)
			}
		}

		for _, cr := range c.ToPoint {
			for _, t := range cr.Tariffs {

				cresult := calcResult{
					Name: t.TariffName,
					MinTerm: t.DaysMin,
					MaxTerm: t.DaysMax,
					Cost: t.DeliveryCost,
					Code: t.TariffId,
					Type: "selfDelivery",
					Group: cr.ProviderKey,
				}

				for _, ptId := range t.PointIds {
					pt := pickpoint{
						Code: ptId,
					}
					cresult.PickuppointList = append(cresult.PickuppointList, pt)
				}
				log.Printf("%#v", cresult)
			}
		}
	}

	return res, nil
}

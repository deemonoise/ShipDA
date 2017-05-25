package main

import (
	"ShipDA/sdt"
	"fmt"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"errors"
	"net/http"
)

type pickpointCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type pickpointExtraData struct {
	Description string `json:"description"` // Описание
	PaymentCard bool   `json:"paymentCard"` // Имеется ли возможность оплаты банковской картой ,
	Cod         bool   `json:"cod"`         // Имеется ли возможность оплаты при доставке
}

type pickpoint struct {
	Url                string               `db:"url" json:"-"`
	Area               string               `db:"area" json:"-"`
	Street             string               `db:"street" json:"-"`
	Description        string               `db:"description" json:"-"`
	Office             string               `db:"office" json:"-"`
	Email              string               `db:"email" json:"-"`
	Cod                string               `db:"cod" json:"-"`
	AvailableOperation int                  `db:"available_operation" json:"-"`
	PaymentCard        string               `db:"payment_card" json:"-"`
	Code               string               `db:"code" json:"code"` // Код терминала
	CountryCode        string               `db:"country_code" json:"-"`
	Region             string               `db:"region" json:"-"`
	Lat                *float64             `db:"lat" json:"-"`
	Lng                *float64             `db:"lng" json:"-"`
	City               string               `db:"city" json:"-"`
	Timetable          string               `db:"timetable" json:"schedule"` // Режим работы
	Name               string               `db:"name" json:"name"`          // Наименование терминала
	StreetType         string               `db:"street_type" json:"-"`
	PostIndex          string               `db:"post_index" json:"-"`
	House              string               `db:"house" json:"-"`
	ProviderKey        string               `db:"provider_key" json:"-"`
	Type               int                  `db:"type" json:"-"`
	Block              string               `db:"block" json:"-"`
	Phone              string               `db:"phone" json:"phone"` // Телефон
	Address            string               `json:"address"`          // Адрес
	Coordinates        pickpointCoordinates `json:"coordinates"`      // Координаты
	ExtraData          pickpointExtraData   `json:"extraData"`        // Дополнительные данные
}

type pickpoints struct {
	Success   bool        `json:"success"`
	Pickpoint []pickpoint `json:"result"`
}

func savePoints(points sdt.Pickpoints) (int, int, error) {
	i := 0;
	u := 0;
	insert := `
	INSERT INTO shipda_pickpoints (external_id, provider_key, type, available_operation, cod, payment_card, name, lat, lng, code,
                               post_index, country_code, region, area, city, city_guid, street, street_type, house,
                               block, office, url, email, phone, timetable, description)
	VALUES (:external_id, :provider_key, :type, :available_operation, :cod, :payment_card, :name, :lat, :lng, :code,
                 :post_index, :country_code, :region, :area, :city, :city_guid, :street, :street_type, :house,
                              :block, :office, :url, :email, :phone, :timetable, :description);
	`
	update := `
	UPDATE shipda_pickpoints SET provider_key = :provider_key, type = :type, available_operation = :available_operation,
	  cod = :cod, payment_card = :payment_card, name = :name, lat = :lat, lng = :lng, code = :code, post_index = :post_index,
	  country_code = :country_code, region = :region, area = :area, city = :city, city_guid = :city_guid, street = :street,
	  street_type = :street_type, house = :house, block = :block, office = :office, url = :url, email = :email, phone = :phone,
	  timetable = :timetable, description = :description WHERE external_id = :external_id;
	`

	stmntIns, err := db.PrepareNamed(insert)
	if err != nil {
		return i, u, err
	}

	stmntUpd, err := db.PrepareNamed(update)
	if err != nil {
		return i, u, err
	}

	for _, point := range points.Pickpoint {
		var exist bool
		_ = db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM shipda_pickpoints WHERE external_id = ?", point.Id).Scan(&exist)

		if exist {
			_, err := stmntUpd.Exec(point)
			if err != nil {
				return i, u, err
			}
			u++
		} else {
			_, err := stmntIns.Exec(point)
			if err != nil {
				return i, u, err
			}
			i++
		}

	}
	return i, u, nil
}

func getPointsJson(city string, region string) ([]byte, *appError) {
	var result []byte
	var err error
	var rows *sqlx.Rows

	if city != "" && region != "" {
		sel := `SELECT
		  provider_key, type, available_operation, cod, payment_card, name, lat, lng, code, post_index, country_code,
		  region, area, city, street, street_type, house, block, office, url, email, phone, timetable, description
		FROM shipda_pickpoints WHERE city = ? AND region = ?`
		rows, err = db.Queryx(sel, city, region)
	} else if city != "" && region == "" {
		sel := `SELECT
		  provider_key, type, available_operation, cod, payment_card, name, lat, lng, code, post_index, country_code,
		  region, area, city, street, street_type, house, block, office, url, email, phone, timetable, description
		FROM shipda_pickpoints WHERE city = ?`
		rows, err = db.Queryx(sel, city)
	} else if city == "" && region != "" {
		sel := `SELECT
		  provider_key, type, available_operation, cod, payment_card, name, lat, lng, code, post_index, country_code,
		  region, area, city, street, street_type, house, block, office, url, email, phone, timetable, description
		FROM shipda_pickpoints WHERE region = ?`
		rows, err = db.Queryx(sel, region)
	} else {
		e := NewAppError(http.StatusBadRequest, errors.New("city and region fields is empty"))
		return result, &e
	}
	if err != nil {
		e := NewAppError(http.StatusInternalServerError, err)
		return result, &e
	}
	defer rows.Close()

	var points pickpoints

	ae := points.scan(rows)
	if ae != nil {
		return result, ae
	}

	if len(points.Pickpoint) == 0 {
		e := NewAppError(http.StatusNotFound, errors.New(fmt.Sprintf("pickpoints is not found for city = \"%s\" and region = \"%s\"", city, region)))
		return result, &e
	}

	points.Success = true
	result, err = json.Marshal(points)
	if err != nil {
		e := NewAppError(http.StatusInternalServerError, err)
		return result, &e
	}

	return result, nil
}

func (p *pickpoint) format() {
	if p.Lat != nil && p.Lng != nil {
		p.Coordinates.Latitude = *p.Lat
		p.Coordinates.Longitude = *p.Lng
	}

	p.Address = fmt.Sprintf("%s, %s, %s, %s, %s %s, д. %s", p.PostIndex, p.CountryCode, p.Region, p.City, p.StreetType, p.Street, p.House)
	if p.Block != "" {
		p.Address += "к" + p.Block
	}
	if p.Office != "" {
		p.Address += " офис " + p.Office
	}
	p.ExtraData.Description = p.Description
	p.ExtraData.PaymentCard = p.PaymentCard == "1"
	p.ExtraData.Cod = p.Cod == "1"
}

func (points *pickpoints) scan(rows *sqlx.Rows) *appError {
	for rows.Next() {
		var p pickpoint
		err := rows.StructScan(&p)
		if err != nil {
			e := NewAppError(http.StatusInternalServerError, err)
			return &e
		}
		p.format()
		points.Pickpoint = append(points.Pickpoint, p)
	}

	return nil
}

func getPointByExternalId(externalId string) (pickpoint, *appError) {
	var point pickpoint

	sel := `SELECT
		  provider_key, type, available_operation, cod, payment_card, name, lat, lng, code, post_index, country_code,
		  region, area, city, street, street_type, house, block, office, url, email, phone, timetable, description
		FROM shipda_pickpoints WHERE external_id = ?`
	err := db.Select(&point, sel, externalId)
	if err != nil {
		e := NewAppError(http.StatusInternalServerError, err)
		return point, &e
	}
	point.format()
	return point, nil
}

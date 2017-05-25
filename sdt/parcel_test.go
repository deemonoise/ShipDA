package sdt

import (
	"testing"
	"log"
)

func TestSdt_Calculate(t *testing.T) {
	client := Init("https://api.accordpost.ru/ff/v1/wsrv", "300", "300")

	p := Parcel{
		CityFrom: "Москва",
		CityTo: "Москва",
		Weight: 1000,
		Width: 10,
		Height: 10,
		Length: 10,
		AssessedCost: 1000,
		Zip: "353180",
	}

	crs, rp, err := client.Calculate(p, []int{58, 62})

	if err != nil {
		t.Error("Got", err)
	}

	log.Printf("%v", rp)
	log.Println("----------------------------------")
	for _, c := range crs.Courier {
		log.Println("delivery_type:", c.ServiceId)
		log.Printf("%#v", c.ToDoor)
		log.Printf("%#v", c.ToPoint)
	}
}

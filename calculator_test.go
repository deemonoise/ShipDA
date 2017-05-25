package main

import (
	"testing"
	"encoding/json"
	"log"
	"ShipDA/sdt"
)

var j = `{
  "clientId": "",
  "calculate": {
    "packages": [
      {
        "weight": 500,
        "width": 100,
        "length": 100,
        "height": 100
      }
    ],
    "payerType": "reciever",
    "declaredValue": 1500,
    "cod": 1500,
    "shipmentAddress": {
      "index": "353180",
      "city": "Самара"
    }
  }
}`

func TestCalculate(t *testing.T) {
	var req calcRequest

	Sdt =  sdt.Init("https://api.accordpost.ru/ff/v1/wsrv", "300", "300")

	Cfg = Config{
		Database: Database{
			Host: "172.20.100.214",
			Port: "3306",
			Db: "dpribylnov_common",
			User: "root",
			Password: "magento",
		},
	}

	dbInit()

	err := json.Unmarshal([]byte(j), &req)
	if err != nil {
		t.Error("Got", err)
	}

	res, aerr := req.Calculate()
	if aerr != nil {
		t.Error("Got", &aerr.Error)
	}

	jsn, err := json.Marshal(res)
	if err != nil {
		t.Error("Got", err)
	}

	log.Println(string(jsn))
}

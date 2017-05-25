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

	err := json.Unmarshal([]byte(j), &req)
	if err != nil {
		t.Error("Got", err)
	}

	res, _ := req.Calculate()
	if err != nil {
		//t.Error("Got", &aerr.Error)
	}

	log.Println(res)
}

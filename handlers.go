package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"strings"
)

/*
Index page
 */
func HandlerIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintln(writer, "Здесь нет ничего")
}

/*
Start Pickpoints update
 */
func HandlerPickpointsUpdate(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("Points update started")
	points, err := Sdt.GetPointsFromApi("", []int{62})
	if err != nil {
		e := NewAppError(http.StatusServiceUnavailable, err)
		e.Write(writer)
		return
	}

	log.Printf("Got %v points from SDT", points.Total)

	type pres struct {
		Message string `json:"message"`
	}

	inserted, updated, err := savePoints(points)
	if err != nil {
		e := NewAppError(http.StatusInternalServerError, err)
		e.Write(writer)
		return
	}

	log.Printf("Total %v point got from SDT. %v point updated, %v created", points.Total, updated, inserted)

	writer.WriteHeader(http.StatusOK)
	r := pres{
		Message: fmt.Sprintf("%v points updated, %v points created", updated, inserted),
	}
	err = json.NewEncoder(writer).Encode(r)
	if err != nil {
		e := NewAppError(http.StatusInternalServerError, err)
		e.Write(writer)
		return
	}
}

/*
Get Pickpoints list by city and region
 */
func HandlerPickpointsList(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := request.URL.Query();

	city := strings.Join(vars["city"], "")
	region := strings.Join(vars["region"], "")

	res, e := getPointsJson(city, region)
	if e != nil {
		e.Write(writer)
		return
	}

	_, err := writer.Write(res)
	if err != nil {
		e := NewAppError(http.StatusInternalServerError, err)
		e.Write(writer)
		return
	}
}

func HandlerCalculate(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	decoder := json.NewDecoder(request.Body)
	var calcReq calcRequest
	err := decoder.Decode(&calcReq)
	if err != nil {
		e := NewAppError(http.StatusInternalServerError, err)
		e.Write(writer)
		return
	}
}
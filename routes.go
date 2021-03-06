package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Comment     string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"",
		"GET",
		"/",
		HandlerIndex,
	},
	Route{
		"PickpointsUpdate",
		"",
		"GET",
		"/pickpoints/update",
		HandlerPickpointsUpdate,
	},
	Route{
		"Pickpoints list",
		"",
		"GET",
		"/pickpoints/list",
		HandlerPickpointsList,
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
	return router
}
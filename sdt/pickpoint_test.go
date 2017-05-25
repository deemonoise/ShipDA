package sdt

import (
	"testing"
	"fmt"
)

func TestSdt_GetPointsFromApi(t *testing.T) {
	client := Init("https://api.accordpost.ru/ff/v1/wsrv", "300", "300")
	points, err := client.GetPointsFromApi("", []int{})

	if err != nil {
		t.Error("Got", err)
	}

	fmt.Println("Total", points.Total, "points found")

	for _, point := range points.Pickpoint {
		fmt.Println("id:", point.Id,"providerKey:", point.ProviderKey, "| name:", point.Name, "| availableOperation:", point.AvailableOperation)
	}
}
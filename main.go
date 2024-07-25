package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Structs to represent the data
type Booking struct {
	BookingDate      string    `json:"bookingDate"`
	OfficeName       string    `json:"officeName"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	ListConsumption  []Consumption `json:"listConsumption"`
	Participants     int       `json:"participants"`
	RoomName         string    `json:"roomName"`
	ID               string    `json:"id"`
}

type Consumption struct {
	Name string `json:"name"`
}

type ConsumptionType struct {
	CreatedAt string `json:"createdAt"`
	Name      string `json:"name"`
	MaxPrice  int    `json:"maxPrice"`
	ID        string `json:"id"`
}

type Summary struct {
	TotalBookings    int
	TotalParticipants int
	ConsumptionCounts map[string]int
	TotalConsumptionCost int
}

func main() {
	bookings, err := fetchBookings()
	if err != nil {
		fmt.Println("Error fetching bookings:", err)
		return
	}

	consumptionTypes, err := fetchConsumptionTypes()
	if err != nil {
		fmt.Println("Error fetching consumption types:", err)
		return
	}

	summary := generateSummary(bookings, consumptionTypes)
	displaySummary(summary)
}

func fetchBookings() ([]Booking, error) {
	resp, err := http.Get("https://66876cc30bc7155dc017a662.mockapi.io/api/dummy-data/bookingList")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bookings []Booking
	err = json.Unmarshal(body, &bookings)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func fetchConsumptionTypes() ([]ConsumptionType, error) {
	resp, err := http.Get("https://6686cb5583c983911b03a7f3.mockapi.io/api/dummy-data/masterJenisKonsumsi")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var consumptionTypes []ConsumptionType
	err = json.Unmarshal(body, &consumptionTypes)
	if err != nil {
		return nil, err
	}

	return consumptionTypes, nil
}

func generateSummary(bookings []Booking, consumptionTypes []ConsumptionType) Summary {
	summary := Summary{
		ConsumptionCounts: make(map[string]int),
	}

	for _, booking := range bookings {
		summary.TotalBookings++
		summary.TotalParticipants += booking.Participants

		for _, consumption := range booking.ListConsumption {
			summary.ConsumptionCounts[consumption.Name]++
		}
	}

	for _, consumptionType := range consumptionTypes {
		count := summary.ConsumptionCounts[consumptionType.Name]
		summary.TotalConsumptionCost += count * consumptionType.MaxPrice
	}

	return summary
}

func displaySummary(summary Summary) {
	fmt.Println("Dashboard Summary:")
	fmt.Printf("Total Bookings: %d\n", summary.TotalBookings)
	fmt.Printf("Total Participants: %d\n", summary.TotalParticipants)
	fmt.Println("Consumption Counts:")
	for name, count := range summary.ConsumptionCounts {
		fmt.Printf("  %s: %d\n", name, count)
	}
	fmt.Printf("Total Consumption Cost: Rp %d\n", summary.TotalConsumptionCost)
}
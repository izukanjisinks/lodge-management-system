package utils

import "time"

// ZambianHoliday represents a Zambian public holiday
type ZambianHoliday struct {
	Month int
	Day   int
	Name  string
}

var ZambianHolidays = []ZambianHoliday{
	{Month: 1, Day: 1, Name: "New Year's Day"},
	{Month: 3, Day: 8, Name: "International Women's Day"},
	{Month: 3, Day: 12, Name: "Youth Day"},
	{Month: 4, Day: 28, Name: "Kenneth Kaunda Day"},
	{Month: 5, Day: 1, Name: "Labour Day"},
	{Month: 5, Day: 25, Name: "Africa Freedom Day"},
	{Month: 7, Day: 7, Name: "Heroes' Day"},
	{Month: 7, Day: 8, Name: "Unity Day"},
	{Month: 8, Day: 4, Name: "Farmers' Day"},
	{Month: 10, Day: 18, Name: "National Prayer Day"},
	{Month: 10, Day: 24, Name: "Independence Day"},
	{Month: 12, Day: 25, Name: "Christmas Day"},
}

// GetHolidaysForMonth returns all Zambian holidays for a given month
func GetHolidaysForMonth(month time.Month) []ZambianHoliday {
	var result []ZambianHoliday
	for _, h := range ZambianHolidays {
		if h.Month == int(month) {
			result = append(result, h)
		}
	}
	return result
}

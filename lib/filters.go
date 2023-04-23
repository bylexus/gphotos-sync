package lib

type MediaFilter struct {
	DateFilter *DateFilter `json:"dateFilter,omitempty"`
}

type DateFilter struct {
	Dates  []Date      `json:"dates,omitempty"`
	Ranges []DateRange `json:"ranges,omitempty"`
}

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type DateRange struct {
	StartDate Date `json:"startDate"`
	EndDate   Date `json:"endDate"`
}

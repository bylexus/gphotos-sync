package lib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type MediaFilter struct {
	DateFilter *DateFilter `json:"dateFilter,omitempty"`
}

func (m MediaFilter) String() string {
	var str = ""
	if m.DateFilter != nil {
		str += fmt.Sprintf("Date Filters:\n%s\n", *m.DateFilter)
	}
	return str
}

type DateFilter struct {
	Dates  []Date      `json:"dates,omitempty"`
	Ranges []DateRange `json:"ranges,omitempty"`
}

func (d DateFilter) String() string {
	var str = ""
	if len(d.Dates) > 0 {
		str += "Dates: "
		for _, d := range d.Dates {
			str += fmt.Sprintf("  Date: %s\n", d)
		}
	}
	if len(d.Ranges) > 0 {
		str += "DateRanges: "
		for _, d := range d.Ranges {
			str += fmt.Sprintf("  Range: %s\n", d)
		}
	}
	return str
}

func (m *MediaFilter) AppendDatesFromStrings(inputs []string) error {
	re := regexp.MustCompile(`^(\d{4})(-(\d{2})(-(\d{2}))?)?$`)
	for _, input := range inputs {
		groups := re.FindStringSubmatch(input)
		if len(groups) == 6 && len(groups[1]) == 4 {
			year, err := strconv.ParseInt(groups[1], 10, 32)
			if err != nil || year <= 0 {
				return errors.New("cannot parse Year in " + input)
			}
			date := Date{Year: int(year)}

			// month
			if len(groups[3]) == 2 {
				month, err := strconv.ParseInt(groups[3], 10, 32)
				if err != nil || month <= 1 || month >= 12 {
					return errors.New("cannot parse Month in " + input)
				}
				date.Month = int(month)
			}

			// day
			if len(groups[5]) == 2 {
				day, err := strconv.ParseInt(groups[5], 10, 32)
				if err != nil || day <= 1 || day >= 31 {
					return errors.New("cannot parse Day in " + input)
				}
				date.Day = int(day)
			}

			if m.DateFilter == nil {
				m.DateFilter = &DateFilter{}
			}
			m.DateFilter.Dates = append(m.DateFilter.Dates, date)
		}
	}
	return nil
}

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

func (d Date) String() string {
	var str = ""
	if d.Year > 0 {
		str += fmt.Sprintf("%d", d.Year)
	}
	if d.Month > 0 {
		str += fmt.Sprintf("-%d", d.Month)
	}
	if d.Day > 0 {
		str += fmt.Sprintf("-%d", d.Day)
	}
	return str
}

type DateRange struct {
	StartDate Date `json:"startDate"`
	EndDate   Date `json:"endDate"`
}

func (d DateRange) String() string {
	return "DateRange.String(): TBD"
}

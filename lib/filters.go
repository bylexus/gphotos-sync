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

func StringToDate(str string) (*Date, error) {
	re := regexp.MustCompile(`^(\d{4})(-(\d{2})(-(\d{2}))?)?$`)
	groups := re.FindStringSubmatch(str)
	if len(groups) == 6 && len(groups[1]) == 4 {
		year, err := strconv.ParseInt(groups[1], 10, 32)
		if err != nil || year <= 0 {
			return nil, errors.New("cannot parse Year in " + str)
		}
		date := Date{Year: int(year)}

		// month
		if len(groups[3]) == 2 {
			month, err := strconv.ParseInt(groups[3], 10, 32)
			if err != nil || month < 1 || month > 12 {
				return nil, errors.New("cannot parse Month in " + str)
			}
			date.Month = int(month)
		}

		// day
		if len(groups[5]) == 2 {
			day, err := strconv.ParseInt(groups[5], 10, 32)
			if err != nil || day < 1 || day > 31 {
				return nil, errors.New("cannot parse Day in " + str)
			}
			date.Day = int(day)
		}

		return &date, nil
	}
	return nil, errors.New("cannot parse date: " + str)
}

func (m *MediaFilter) AppendDatesFromStrings(inputs []string) error {
	for _, input := range inputs {
		date, err := StringToDate(input)
		if err != nil {
			return err
		}

		if m.DateFilter == nil {
			m.DateFilter = &DateFilter{}
		}
		m.DateFilter.Dates = append(m.DateFilter.Dates, *date)
	}
	return nil
}

func (m *MediaFilter) AppendDateRangesFromStrings(inputs []string) error {
	re := regexp.MustCompile(`^((\d{4})(-(\d{2})(-(\d{2}))?)?):((\d{4})(-(\d{2})(-(\d{2}))?)?)$`)
	for _, input := range inputs {
		groups := re.FindStringSubmatch(input)
		if len(groups) == 13 {
			startDateStr := groups[1]
			endDateStr := groups[7]

			startDate, err := StringToDate(startDateStr)
			if err != nil {
				return errors.New(fmt.Sprintf("cannot parse start date: %s", err))
			}
			endDate, err := StringToDate(endDateStr)
			if err != nil {
				return errors.New(fmt.Sprintf("cannot parse end date: %s", err))
			}

			// Both dates must be of the same format when using date ranges:
			if startDate.Day == 0 && endDate.Day > 0 || startDate.Month > 0 && endDate.Month == 0 {
				return errors.New("Both startDate and endDate must be of the same format")
			}

			if m.DateFilter == nil {
				m.DateFilter = &DateFilter{}
			}
			m.DateFilter.Ranges = append(m.DateFilter.Ranges, DateRange{StartDate: *startDate, EndDate: *endDate})
		}
	}
	return nil
}

type DateFilter struct {
	Dates  []Date      `json:"dates,omitempty"`
	Ranges []DateRange `json:"ranges,omitempty"`
}

func (d DateFilter) String() string {
	var str = ""
	if len(d.Dates) > 0 {
		str += "Dates: \n"
		for _, d := range d.Dates {
			str += fmt.Sprintf("  Date: %s\n", d)
		}
	}
	if len(d.Ranges) > 0 {
		str += "DateRanges: \n"
		for _, d := range d.Ranges {
			str += fmt.Sprintf("  Range: %s\n", d)
		}
	}
	return str
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
		str += fmt.Sprintf("-%02d", d.Month)
	}
	if d.Day > 0 {
		str += fmt.Sprintf("-%02d", d.Day)
	}
	return str
}

type DateRange struct {
	StartDate Date `json:"startDate"`
	EndDate   Date `json:"endDate"`
}

func (d DateRange) String() string {
	startDate := d.StartDate.String()
	endDate := d.EndDate.String()
	return fmt.Sprintf("%s : %s", startDate, endDate)
}

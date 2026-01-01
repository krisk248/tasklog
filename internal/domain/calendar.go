package domain

import (
	"fmt"
	"time"
)

type CalendarDate struct {
	Year  int
	Month time.Month
	Day   int
}

func NewCalendarDate(t time.Time) CalendarDate {
	return CalendarDate{
		Year:  t.Year(),
		Month: t.Month(),
		Day:   t.Day(),
	}
}

func Today() CalendarDate {
	return NewCalendarDate(time.Now())
}

func (d CalendarDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, int(d.Month), d.Day)
}

func (d CalendarDate) Time() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.Local)
}

func (d CalendarDate) Format(layout string) string {
	return d.Time().Format(layout)
}

func (d CalendarDate) AddDays(days int) CalendarDate {
	return NewCalendarDate(d.Time().AddDate(0, 0, days))
}

func (d CalendarDate) AddMonths(months int) CalendarDate {
	return NewCalendarDate(d.Time().AddDate(0, months, 0))
}

func (d CalendarDate) IsToday() bool {
	today := Today()
	return d.Year == today.Year && d.Month == today.Month && d.Day == today.Day
}

func (d CalendarDate) IsSameMonth(other CalendarDate) bool {
	return d.Year == other.Year && d.Month == other.Month
}

func (d CalendarDate) Equals(other CalendarDate) bool {
	return d.Year == other.Year && d.Month == other.Month && d.Day == other.Day
}

func (d CalendarDate) Weekday() time.Weekday {
	return d.Time().Weekday()
}

// MonthName returns the full month name
func (d CalendarDate) MonthName() string {
	return d.Month.String()
}

// ShortMonthName returns the abbreviated month name
func (d CalendarDate) ShortMonthName() string {
	return d.Month.String()[:3]
}

// DaysInMonth returns the number of days in the month
func (d CalendarDate) DaysInMonth() int {
	// Go to first day of next month, then subtract 1 day
	firstOfNextMonth := time.Date(d.Year, d.Month+1, 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastOfMonth.Day()
}

// FirstDayOfMonth returns the CalendarDate for the first day of the month
func (d CalendarDate) FirstDayOfMonth() CalendarDate {
	return CalendarDate{Year: d.Year, Month: d.Month, Day: 1}
}

// FirstDayWeekday returns the weekday of the first day of the month
func (d CalendarDate) FirstDayWeekday() time.Weekday {
	return d.FirstDayOfMonth().Weekday()
}

// WeekOfMonth returns which week of the month this date is in (0-indexed)
func (d CalendarDate) WeekOfMonth() int {
	firstDay := d.FirstDayOfMonth()
	dayOffset := int(firstDay.Weekday())
	return (d.Day + dayOffset - 1) / 7
}

// CalendarGrid represents a month view with weeks
type CalendarGrid struct {
	Year      int
	Month     time.Month
	Weeks     [][]CalendarDate // 6 weeks x 7 days
	MonthDays map[string]bool  // Quick lookup for days in current month
}

// GenerateCalendarGrid creates a calendar grid for display
func GenerateCalendarGrid(date CalendarDate) CalendarGrid {
	grid := CalendarGrid{
		Year:      date.Year,
		Month:     date.Month,
		Weeks:     make([][]CalendarDate, 6),
		MonthDays: make(map[string]bool),
	}

	// Find the first day to display (Sunday of the week containing the 1st)
	firstOfMonth := date.FirstDayOfMonth()
	startDate := firstOfMonth.AddDays(-int(firstOfMonth.Weekday()))

	// Fill in the grid
	currentDate := startDate
	for week := 0; week < 6; week++ {
		grid.Weeks[week] = make([]CalendarDate, 7)
		for day := 0; day < 7; day++ {
			grid.Weeks[week][day] = currentDate
			if currentDate.Month == date.Month {
				grid.MonthDays[currentDate.String()] = true
			}
			currentDate = currentDate.AddDays(1)
		}
	}

	return grid
}

// IsCurrentMonth checks if a date belongs to the grid's month
func (g CalendarGrid) IsCurrentMonth(date CalendarDate) bool {
	return date.Month == g.Month && date.Year == g.Year
}

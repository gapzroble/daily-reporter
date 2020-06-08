package schedule

var holidays = map[string]string{
	"2020-02-24": "Cebu City Charter Day",
	"2020-02-25": "EDSA People Power Revolution Anniversary",
	"2020-12-26": "leave",
}

// IsHolidayOrLeave check date
func IsHolidayOrLeave() bool {
	_, ok := holidays[Today]
	return ok
}

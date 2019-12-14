package types // import "golang.handcraftedbits.com/ezif"

import "C"

import (
	"fmt"
	"time"
)

//
// Public types
//

type IPTCDate interface {
	fmt.Stringer

	Day() int
	Month() time.Month
	ToGoTime(iptcTime IPTCTime) time.Time
	Year() int
}

type IPTCTime interface {
	fmt.Stringer

	Hour() int
	Minute() int
	Second() int
	Timezone() *time.Location
	ToGoTime(date IPTCDate) time.Time
}

//
// Public functions
//

func NewIPTCDate(year int, month time.Month, day int) IPTCDate {
	return &iptcDateImpl{
		day:   day,
		month: month,
		year:  year,
	}
}

func NewIPTCTime(hour, minute, second, offsetHours, offsetMinutes int) IPTCTime {
	return &iptcTimeImpl{
		hour:          hour,
		minute:        minute,
		offsetHours:   offsetHours,
		offsetMinutes: offsetMinutes,
		second:        second,
		timezone:      time.FixedZone("IPTC time", (offsetHours*60*60)+(offsetMinutes*60)),
	}
}

//
// Private types
//

// IPTCDate implementation
type iptcDateImpl struct {
	day   int
	month time.Month
	year  int
}

func (date *iptcDateImpl) Day() int {
	return date.day
}

func (date *iptcDateImpl) Month() time.Month {
	return date.month
}

func (date *iptcDateImpl) String() string {
	return fmt.Sprintf("%4d%2d%2d", date.year, date.month, date.day)
}

func (date *iptcDateImpl) ToGoTime(iptcTime IPTCTime) time.Time {
	return time.Date(date.year, date.month, date.day, iptcTime.Hour(), iptcTime.Minute(), iptcTime.Second(), 0,
		iptcTime.Timezone())
}

func (date *iptcDateImpl) Year() int {
	return date.year
}

// IPTCTime implementation
type iptcTimeImpl struct {
	hour          int
	minute        int
	offsetHours   int
	offsetMinutes int
	second        int
	timezone      *time.Location
}

func (iptcTime *iptcTimeImpl) Hour() int {
	return iptcTime.hour
}

func (iptcTime *iptcTimeImpl) Minute() int {
	return iptcTime.minute
}

func (iptcTime *iptcTimeImpl) Second() int {
	return iptcTime.second
}

func (iptcTime *iptcTimeImpl) String() string {
	return fmt.Sprintf("%2d%2d%2d:%2d%2d", iptcTime.hour, iptcTime.minute, iptcTime.second, iptcTime.offsetHours,
		iptcTime.offsetMinutes)
}

func (iptcTime *iptcTimeImpl) Timezone() *time.Location {
	return iptcTime.timezone
}

func (iptcTime *iptcTimeImpl) ToGoTime(date IPTCDate) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), iptcTime.hour, iptcTime.minute, iptcTime.second, 0,
		iptcTime.timezone)
}

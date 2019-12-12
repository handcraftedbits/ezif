package ezif // import "golang.handcraftedbits.com/ezif"

import "C"

import (
	"time"
)

//
// Public types
//

// TODO: move to iptc package?
type IPTCDate interface {
	Day() int
	Month() int
	Year() int
}

type IPTCTime interface {
	Hour() int
	Minute() int
	Second() int
	Timezone() *time.Location
}

//
// Private types
//

// IPTCDate implementation
type iptcDateImpl struct {
	day   int
	month int
	year  int
}

func (date *iptcDateImpl) Day() int {
	return date.day
}

func (date *iptcDateImpl) Month() int {
	return date.month
}

func (date *iptcDateImpl) Year() int {
	return date.year
}

// IPTCTime implementation
type iptcTimeImpl struct {
	hour     int
	minute   int
	second   int
	timezone *time.Location
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

func (iptcTime *iptcTimeImpl) Timezone() *time.Location {
	return iptcTime.timezone
}

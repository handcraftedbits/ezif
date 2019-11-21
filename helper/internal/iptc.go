package internal // import "golang.handcraftedbits.com/ezif/helper/internal"

import (
	"time"

	"golang.handcraftedbits.com/ezif"
)

//
// Public functions
//

func GetIPTCValueAsDate(metadatum ezif.IPTCDatum) time.Time {
	return metadatum.Values()[0].Date()
}

func GetIPTCValueAsDateSlice(metadatum ezif.IPTCDatum) []time.Time {
	var result = make([]time.Time, len(metadatum.Values()))

	for i := range metadatum.Values() {
		result[i] = metadatum.Values()[i].Date()
	}

	return result
}

func GetIPTCValueAsString(metadatum ezif.IPTCDatum) string {
	return metadatum.Values()[0].String()
}

func GetIPTCValueAsStringSlice(metadatum ezif.IPTCDatum) []string {
	var result = make([]string, len(metadatum.Values()))

	for i := range metadatum.Values() {
		result[i] = metadatum.Values()[i].String()
	}

	return result
}

func GetIPTCValueAsTime(metadatum ezif.IPTCDatum) time.Time {
	return metadatum.Values()[0].Time()
}

func GetIPTCValueAsUndefinedSlice(metadatum ezif.IPTCDatum) []byte {
	return metadatum.Values()[0].Undefined()
}

func GetIPTCValueAsUnsignedShort(metadatum ezif.IPTCDatum) uint16 {
	return metadatum.Values()[0].Short()
}

package internal // import "golang.handcraftedbits.com/ezif/helper/internal"

import (
	"time"

	"golang.handcraftedbits.com/ezif"
)

//
// Public functions
//

func GetIPTCValueAsDate(metadatum []ezif.IPTCDatum) time.Time {
	return metadatum[0].Value().Date()
}

func GetIPTCValueAsDateSlice(metadatum []ezif.IPTCDatum) []time.Time {
	var result = make([]time.Time, len(metadatum))

	for i := range metadatum {
		result[i] = metadatum[i].Value().Date()
	}

	return result
}

func GetIPTCValueAsString(metadatum []ezif.IPTCDatum) string {
	return metadatum[0].Value().String()
}

func GetIPTCValueAsStringSlice(metadatum []ezif.IPTCDatum) []string {
	var result = make([]string, len(metadatum))

	for i := range metadatum {
		result[i] = metadatum[i].Value().String()
	}

	return result
}

func GetIPTCValueAsTime(metadatum []ezif.IPTCDatum) time.Time {
	return metadatum[0].Value().Time()
}

func GetIPTCValueAsUndefinedSlice(metadatum []ezif.IPTCDatum) []byte {
	return metadatum[0].Value().Undefined()
}

func GetIPTCValueAsUnsignedShort(metadatum []ezif.IPTCDatum) uint16 {
	return metadatum[0].Value().Short()
}

package internal // import "golang.handcraftedbits.com/ezif/helper/internal"

import (
	"math/big"

	"golang.handcraftedbits.com/ezif"
)

//
// Public functions
//

func GetExifValueAsASCIIString(metadatum ezif.ExifDatum) string {
	return metadatum.Values()[0].ASCIIString()
}

func GetExifValueAsComment(metadatum ezif.ExifDatum) string {
	return metadatum.Values()[0].Comment()
}

func GetExifValueAsSignedByte(metadatum ezif.ExifDatum) int8 {
	return metadatum.Values()[0].SignedByte()
}

func GetExifValueAsSignedByteSlice(metadatum ezif.ExifDatum) []int8 {
	var result = make([]int8, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.SignedByte()
	}

	return result
}

func GetExifValueAsSignedLong(metadatum ezif.ExifDatum) int32 {
	return metadatum.Values()[0].SignedLong()
}

func GetExifValueAsSignedLongSlice(metadatum ezif.ExifDatum) []int32 {
	var result = make([]int32, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.SignedLong()
	}

	return result
}

func GetExifValueAsSignedRational(metadatum ezif.ExifDatum) *big.Rat {
	return metadatum.Values()[0].SignedRational()
}

func GetExifValueAsSignedRationalSlice(metadatum ezif.ExifDatum) []*big.Rat {
	var result = make([]*big.Rat, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.SignedRational()
	}

	return result
}

func GetExifValueAsSignedShort(metadatum ezif.ExifDatum) int16 {
	return metadatum.Values()[0].SignedShort()
}

func GetExifValueAsSignedShortSlice(metadatum ezif.ExifDatum) []int16 {
	var result = make([]int16, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.SignedShort()
	}

	return result
}

func GetExifValueAsTIFFDoubleSlice(metadatum ezif.ExifDatum) []float64 {
	var result = make([]float64, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.TIFFDouble()
	}

	return result
}

func GetExifValueAsTIFFFloatSlice(metadatum ezif.ExifDatum) []float32 {
	var result = make([]float32, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.TIFFFloat()
	}

	return result
}

func GetExifValueAsUndefined(metadatum ezif.ExifDatum) byte {
	return metadatum.Values()[0].Undefined()
}

func GetExifValueAsUndefinedSlice(metadatum ezif.ExifDatum) []byte {
	var result = make([]byte, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.Undefined()
	}

	return result
}

func GetExifValueAsUnsignedByte(metadatum ezif.ExifDatum) uint8 {
	return metadatum.Values()[0].UnsignedByte()
}

func GetExifValueAsUnsignedByteSlice(metadatum ezif.ExifDatum) []uint8 {
	var result = make([]uint8, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.UnsignedByte()
	}

	return result
}

func GetExifValueAsUnsignedLong(metadatum ezif.ExifDatum) uint32 {
	return metadatum.Values()[0].UnsignedLong()
}

func GetExifValueAsUnsignedLongSlice(metadatum ezif.ExifDatum) []uint32 {
	var result = make([]uint32, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.UnsignedLong()
	}

	return result
}

func GetExifValueAsUnsignedRational(metadatum ezif.ExifDatum) *big.Rat {
	return metadatum.Values()[0].UnsignedRational()
}

func GetExifValueAsUnsignedRationalSlice(metadatum ezif.ExifDatum) []*big.Rat {
	var result = make([]*big.Rat, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.UnsignedRational()
	}

	return result
}

func GetExifValueAsUnsignedShort(metadatum ezif.ExifDatum) uint16 {
	return metadatum.Values()[0].UnsignedShort()
}

func GetExifValueAsUnsignedShortSlice(metadatum ezif.ExifDatum) []uint16 {
	var result = make([]uint16, len(metadatum.Values()))

	for i, value := range metadatum.Values() {
		result[i] = value.UnsignedShort()
	}

	return result
}

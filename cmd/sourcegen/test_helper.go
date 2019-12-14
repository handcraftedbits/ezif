package main // import "golang.handcraftedbits.com/ezif/cmd/sourcegen"

import (
	"bytes"

	"golang.handcraftedbits.com/ezif/types"
)

//
// Private functions
//

func generateGroupTestSource(familyName string, f family, packageName string, gc groupConfig) (string, error) {
	return generateSource(familyName, f, packageName, gc, templateGroupTestSource)
}

func getStringOfFixedLength(length int) string {
	var result = make([]rune, length)

	for i := 0; i < len(result); i++ {
		result[i] = 'A'
	}

	return string(result)
}

func templateFuncMaxValue(familyName string, info functionInfo) string {
	var count = getAdjustedCount(familyName, info)

	// -1 and 0 are "don't know" values, so we'll just use 64 as a decently-sized value.  If nothing else, it can help
	// us identify keys that have hard limits which aren't documented.

	if count <= 0 {
		count = 64
	}

	if count > 1 {
		return testValueSlice(info, count, testMaxValueSingle)
	}

	return testMaxValueSingle(info)
}

func templateFuncMinValue(familyName string, info functionInfo) string {
	var count = getAdjustedCount(familyName, info)

	// -1 and 0 are "don't know" values, so we'll just use 64 as a decently-sized value.  If nothing else, it can help
	// us identify keys that have hard limits which aren't documented.

	if count <= 0 {
		count = 64
	}

	if count > 1 {
		return testValueSlice(info, count, testMinValueSingle)
	}

	return testMinValueSingle(info)
}

func testMaxValueSingle(info functionInfo) string {
	switch info.Tag.TypeID {
	case types.IDAsciiString, types.IDComment:
		// -1 and 0 are essentially "don't know" values, anything else is a fixed string length.

		if info.Tag.Count > 0 {
			return "\"" + getStringOfFixedLength(info.Tag.Count) + "\""
		}

		// Otherwise, try a decent length, like 64

		return "\"" + getStringOfFixedLength(64) + "\""

	case types.IDIPTCString:
		return "\"" + getStringOfFixedLength(info.Tag.MaxBytes) + "\""

	case types.IDIPTCDate, types.IDIPTCTime:
		return "time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)"

	case types.IDSignedByte:
		return "int8(math.MaxInt8)"

	case types.IDSignedLong:
		return "int32(math.MaxInt32)"

	case types.IDSignedRational:
		return "big.NewRat(math.MaxInt32, math.MaxInt32 - 1)"

	case types.IDSignedShort:
		return "int16(math.MaxInt16)"

	case types.IDTIFFDouble:
		// The Exiv2 command seemingly can't parse the MaxFloat64 value, so we'll go with something much bigger than
		// MaxFloat32 just to make sure we can handle a >32 bit value.

		return "float64(9.0e99)"

	case types.IDTIFFFloat:
		// The Exiv2 command will return a rounded value, so we'll round down MaxFloat32 to compensate and make a direct
		// comparison easier.

		return "float32(3.4e38)"

	case types.IDUndefined:
		return "uint8(math.MaxUint8)"

	case types.IDUnsignedByte:
		return "byte(math.MaxUint8)"

	case types.IDUnsignedLong:
		return "uint32(math.MaxUint32)"

	case types.IDUnsignedRational:
		return "big.NewRat(math.MaxUint32, math.MaxUint32 - 1)"

	case types.IDUnsignedShort:
		return "uint16(math.MaxUint16)"

	case types.IDXMPAlt, types.IDXMPBag, types.IDXMPSeq, types.IDXMPText:
		return "\"" + getStringOfFixedLength(64) + "\""

	case types.IDXMPLangAlt:
		return "nil"
	}

	return ""
}

func testMinValueSingle(info functionInfo) string {
	switch info.Tag.TypeID {
	case types.IDAsciiString, types.IDComment, types.IDIPTCString, types.IDXMPText:
		return "\"\""

	case types.IDIPTCDate, types.IDIPTCTime:
		return "time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)"

	case types.IDSignedByte:
		return "int8(math.MinInt8)"

	case types.IDSignedLong:
		return "int32(math.MinInt32)"

	case types.IDSignedRational:
		return "big.NewRat(math.MinInt32 + 1, math.MinInt32)"

	case types.IDSignedShort:
		return "int16(math.MinInt16)"

	case types.IDTIFFDouble:
		return "float64(0.0)"

	case types.IDTIFFFloat:
		return "float32(0.0)"

	case types.IDUndefined:
		return "byte(0)"

	case types.IDUnsignedByte:
		return "uint8(0)"

	case types.IDUnsignedLong:
		return "uint32(0)"

	case types.IDUnsignedRational:
		return "big.NewRat(0, 0)"

	case types.IDUnsignedShort:
		return "uint16(0)"
	}

	return ""
}

func testValueSlice(info functionInfo, count int, singleValueFunc func(info functionInfo) string) string {
	var buffer bytes.Buffer

	buffer.WriteString("[]")
	buffer.WriteString(getTypeIDMapping(info.Tag.TypeID).goType)
	buffer.WriteString("{")

	for i := 0; i < count; i++ {
		var value = singleValueFunc(info)

		if value != "" {
			buffer.WriteString(value)

			if i < count-1 {
				buffer.WriteString(", ")
			}
		}
	}

	buffer.WriteString("}")

	return buffer.String()
}

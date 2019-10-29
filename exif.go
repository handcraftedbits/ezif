package ezif // import "golang.handcraftedbits.com/ezif"

import "C"

import (
	"math/big"
)

//
// Private types
//

// ExifMetadata implementation
type exifMetadataImpl struct {
	currentDatum *exifDatumImpl
	keys         []string
	exifDatumMap map[string]ExifDatum
	valueCount   int
}

func (exifMetadata *exifMetadataImpl) HasKey(key string) bool {
	_, exists := exifMetadata.exifDatumMap[key]

	return exists
}

func (exifMetadata *exifMetadataImpl) Keys() []string {
	return exifMetadata.keys
}

func (exifMetadata *exifMetadataImpl) Get(key string) ExifDatum {
	return exifMetadata.exifDatumMap[key]
}

func (exifMetadata *exifMetadataImpl) add(exifDatum *exifDatumImpl) {
	var key = exifDatum.key()

	if !exifMetadata.HasKey(key) {
		exifMetadata.keys = append(exifMetadata.keys, key)
	}

	exifMetadata.currentDatum = exifDatum
	exifMetadata.exifDatumMap[key] = exifDatum
}

// ExifDatum implementation
type exifDatumImpl struct {
	*datumImpl

	values []ExifValue
}

func (exifDatum *exifDatumImpl) Values() []ExifValue {
	return exifDatum.values
}

// TODO: panic, return error, etc.
func (exifDatum *exifDatumImpl) populateValueFromValueHolder(index int, valueHolder *C.struct_valueHolder) {
	var value = &exifValueImpl{}

	switch exifDatum.TypeId() {
	case TypeIdAsciiString, TypeIdComment:
		value.stringValue = C.GoString(valueHolder.strValue)

	case TypeIdSignedByte, TypeIdSignedLong, TypeIdSignedShort:
		value.intValue = int32(valueHolder.longValue)

	case TypeIdSignedRational, TypeIdUnsignedRational:
		value.rationalValue = big.NewRat(int64(valueHolder.rationalValueN), int64(valueHolder.rationalValueD))

	case TypeIdTIFFDouble, TypeIdTIFFFloat:
		value.floatValue = float64(valueHolder.doubleValue)

	case TypeIdUndefined, TypeIdUnsignedByte, TypeIdUnsignedLong, TypeIdUnsignedShort:
		value.uintValue = uint32(valueHolder.longValue)
	}

	exifDatum.values[index] = value
}

// ExifValue implementation
type exifValueImpl struct {
	floatValue    float64
	intValue      int32
	rationalValue *big.Rat
	stringValue   string
	uintValue     uint32
}

func (exifValue *exifValueImpl) ASCIIString() string {
	return exifValue.stringValue
}

func (exifValue *exifValueImpl) Comment() string {
	return exifValue.stringValue
}

func (exifValue *exifValueImpl) Directory() string {
	return exifValue.stringValue
}

func (exifValue *exifValueImpl) SignedByte() int8 {
	return int8(exifValue.intValue)
}

func (exifValue *exifValueImpl) SignedLong() int32 {
	return int32(exifValue.intValue)
}

func (exifValue *exifValueImpl) SignedRational() *big.Rat {
	return big.NewRat(exifValue.rationalValue.Num().Int64(), exifValue.rationalValue.Denom().Int64())
}

func (exifValue *exifValueImpl) SignedShort() int16 {
	return int16(exifValue.intValue)
}

func (exifValue *exifValueImpl) TIFFDouble() float64 {
	return exifValue.floatValue
}

func (exifValue *exifValueImpl) TIFFFloat() float32 {
	return float32(exifValue.floatValue)
}

func (exifValue *exifValueImpl) Undefined() byte {
	return byte(exifValue.uintValue)
}

func (exifValue *exifValueImpl) UnsignedByte() uint8 {
	return uint8(exifValue.uintValue)
}

func (exifValue *exifValueImpl) UnsignedLong() uint32 {
	return uint32(exifValue.uintValue)
}

func (exifValue *exifValueImpl) UnsignedRational() *big.Rat {
	return big.NewRat(exifValue.rationalValue.Num().Int64(), exifValue.rationalValue.Denom().Int64())
}

func (exifValue *exifValueImpl) UnsignedShort() uint16 {
	return uint16(exifValue.uintValue)
}

//
// Private functions
//

func newExifMetadata() *exifMetadataImpl {
	return &exifMetadataImpl{
		exifDatumMap: make(map[string]ExifDatum),
	}
}

func newExifDatum(familyName, groupName, tagName string, typeId int, label, interpretedValue string,
	numValues int) *exifDatumImpl {
	return &exifDatumImpl{
		datumImpl: &datumImpl{
			familyName:       familyName,
			groupName:        groupName,
			interpretedValue: interpretedValue,
			label:            label,
			tagName:          tagName,
			typeId:           TypeId(typeId),
		},
		values: make([]ExifValue, numValues),
	}
}

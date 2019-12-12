package ezif // import "golang.handcraftedbits.com/ezif"

import "C"

import (
	"fmt"
	"math/big"
	"sort"
)

//
// Public types
//

type Datum interface {
	FamilyName() string
	GroupName() string
	InterpretedValue() string
	Label() string
	TagName() string
	TypeID() TypeID
	Value() interface{}
}

type ImageMetadata interface {
	Exif() Metadata
	IPTC() Metadata
	XMP() Metadata
}

type Metadata interface {
	Get(key string) Datum
	Keys() []string
}

type TypeID int

// TODO: move into xmp package
type XMPLangAlt interface {
	Language() string
	Value() string
}

//
// Public constants
//

const (
	TypeIDUnsignedByte     TypeID = 1
	TypeIDAsciiString      TypeID = 2
	TypeIDUnsignedShort    TypeID = 3
	TypeIDUnsignedLong     TypeID = 4
	TypeIDUnsignedRational TypeID = 5
	TypeIDSignedByte       TypeID = 6
	TypeIDUndefined        TypeID = 7
	TypeIDSignedShort      TypeID = 8
	TypeIDSignedLong       TypeID = 9
	TypeIDSignedRational   TypeID = 10
	TypeIDTIFFFloat        TypeID = 11
	TypeIDTIFFDouble       TypeID = 12
	TypeIDIPTCString       TypeID = 0x10000
	TypeIDIPTCDate         TypeID = 0x10001
	TypeIDIPTCTime         TypeID = 0x10002
	TypeIDComment          TypeID = 0x10003
	TypeIDXMPText          TypeID = 0x10005
	TypeIDXMPAlt           TypeID = 0x10006
	TypeIDXMPBag           TypeID = 0x10007
	TypeIDXMPSeq           TypeID = 0x10008
	TypeIDXMPLangAlt       TypeID = 0x10009
	TypeIDInvalid          TypeID = 0x1FFFE
)

//
// Private types
//

// Datum implementation
type datumImpl struct {
	familyName       string
	groupName        string
	interpretedValue string
	label            string
	tagName          string
	typeId           TypeID
	value            interface{}
}

func (datum *datumImpl) FamilyName() string {
	return datum.familyName
}

func (datum *datumImpl) GroupName() string {
	return datum.groupName
}

func (datum *datumImpl) InterpretedValue() string {
	return datum.interpretedValue
}

func (datum *datumImpl) Label() string {
	return datum.label
}

func (datum *datumImpl) TagName() string {
	return datum.tagName
}

func (datum *datumImpl) TypeID() TypeID {
	return datum.typeId
}

func (datum *datumImpl) Value() interface{} {
	return datum.value
}

func (datum *datumImpl) key() string {
	return datum.familyName + "." + datum.groupName + "." + datum.tagName
}

// ImageMetadata implementation
type imageMetadataImpl struct {
	exifMetadata *metadataImpl
	iptcMetadata *metadataImpl
	xmpMetadata  *metadataImpl
}

func (imageMetadata *imageMetadataImpl) Exif() Metadata {
	return imageMetadata.exifMetadata
}

func (imageMetadata *imageMetadataImpl) IPTC() Metadata {
	return imageMetadata.iptcMetadata
}

func (imageMetadata *imageMetadataImpl) XMP() Metadata {
	return imageMetadata.xmpMetadata
}

// Metadata implementation
type metadataImpl struct {
	datumMap map[string]Datum
	keys     []string
}

func (metadata *metadataImpl) Get(key string) Datum {
	return metadata.datumMap[key]
}

func (metadata *metadataImpl) Keys() []string {
	return metadata.keys
}

func (metadata *metadataImpl) add(datum *datumImpl, values []interface{}) {
	if datum.familyName == "Exif" && datum.groupName == "Photo" && datum.tagName == "UserComment" {
		fmt.Printf("*** add... typeId=%d\n", datum.typeId)
	}
	metadata.datumMap[datum.key()] = datum

	if len(values) == 1 {
		datum.value = values[0]

		return
	}

	// For slice values, we need to do a manual conversion from an interface{} slice to a concrete-typed slice.  This
	// keeps us from having to do a similar conversion later on.

	switch datum.TypeID() {
	case TypeIDAsciiString, TypeIDComment, TypeIDIPTCString, TypeIDXMPAlt, TypeIDXMPBag, TypeIDXMPSeq, TypeIDXMPText:
		var newSlice = make([]string, len(values))

		for i, value := range values {
			newSlice[i] = value.(string)
		}

		datum.value = newSlice

	case TypeIDIPTCDate:
		var newSlice = make([]IPTCDate, len(values))

		for i, value := range values {
			newSlice[i] = value.(IPTCDate)
		}

		datum.value = newSlice

	case TypeIDIPTCTime:
		var newSlice = make([]IPTCTime, len(values))

		for i, value := range values {
			newSlice[i] = value.(IPTCTime)
		}

		datum.value = newSlice

	case TypeIDSignedByte:
		var newSlice = make([]int8, len(values))

		for i, value := range values {
			newSlice[i] = value.(int8)
		}

		datum.value = newSlice

	case TypeIDSignedLong:
		var newSlice = make([]int32, len(values))

		for i, value := range values {
			newSlice[i] = value.(int32)
		}

		datum.value = newSlice

	case TypeIDSignedShort:
		var newSlice = make([]int16, len(values))

		for i, value := range values {
			newSlice[i] = value.(int16)
		}

		datum.value = newSlice

	case TypeIDSignedRational, TypeIDUnsignedRational:
		var newSlice = make([]*big.Rat, len(values))

		for i, value := range values {
			newSlice[i] = value.(*big.Rat)
		}

		datum.value = newSlice

	case TypeIDTIFFDouble:
		var newSlice = make([]float64, len(values))

		for i, value := range values {
			newSlice[i] = value.(float64)
		}

		datum.value = newSlice

	case TypeIDTIFFFloat:
		var newSlice = make([]float32, len(values))

		for i, value := range values {
			newSlice[i] = value.(float32)
		}

		datum.value = newSlice

	case TypeIDUndefined:
		var newSlice = make([]byte, len(values))

		for i, value := range values {
			newSlice[i] = value.(byte)
		}

		datum.value = newSlice

	case TypeIDUnsignedByte:
		var newSlice = make([]uint8, len(values))

		for i, value := range values {
			newSlice[i] = value.(uint8)
		}

		datum.value = newSlice

	case TypeIDUnsignedLong:
		var newSlice = make([]uint32, len(values))

		for i, value := range values {
			newSlice[i] = value.(uint32)
		}

		datum.value = newSlice

	case TypeIDUnsignedShort:
		var newSlice = make([]uint16, len(values))

		for i, value := range values {
			newSlice[i] = value.(uint16)
		}

		datum.value = newSlice

	case TypeIDXMPLangAlt:
		var newSlice = make([]XMPLangAlt, len(values))

		for i, value := range values {
			newSlice[i] = value.(XMPLangAlt)
		}

		datum.value = newSlice
	}
}

func (metadata *metadataImpl) finish() {
	var i = 0

	metadata.keys = make([]string, len(metadata.datumMap))

	for key := range metadata.datumMap {
		metadata.keys[i] = key

		i++
	}

	sort.Strings(metadata.keys)
}

//
// Private functions
//

func newDatum(familyName, groupName, tagName string, typeId int, label, interpretedValue string) *datumImpl {
	return &datumImpl{
		familyName:       familyName,
		groupName:        groupName,
		interpretedValue: interpretedValue,
		label:            label,
		tagName:          tagName,
		typeId:           TypeID(typeId),
	}
}

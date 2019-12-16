package ezif // import "golang.handcraftedbits.com/ezif"

import "C"

import (
	"math/big"
	"sort"

	"golang.handcraftedbits.com/ezif/types"
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
	TypeID() types.ID
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

//
// Private types
//

// Datum implementation
type datumImpl struct {
	familyName       string
	groupName        string
	interpretedValue string
	label            string
	repeatable       bool
	tagName          string
	typeId           types.ID
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

func (datum *datumImpl) TypeID() types.ID {
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
	metadata.datumMap[datum.key()] = datum

	// TODO: handle repeatable appending into an existing value.

	// We need to do a manual conversion from an interface{} slice to a concrete-typed slice.  This keeps us from having
	// to do a similar conversion later on.

	switch datum.TypeID() {
	case types.IDAsciiString, types.IDComment, types.IDIPTCString, types.IDXMPAlt, types.IDXMPBag, types.IDXMPSeq,
		types.IDXMPText:
		var newSlice = make([]string, len(values))

		for i, value := range values {
			newSlice[i] = value.(string)
		}

		datum.value = newSlice

	case types.IDIPTCDate:
		var newSlice = make([]types.IPTCDate, len(values))

		for i, value := range values {
			newSlice[i] = value.(types.IPTCDate)
		}

		datum.value = newSlice

	case types.IDIPTCTime:
		var newSlice = make([]types.IPTCTime, len(values))

		for i, value := range values {
			newSlice[i] = value.(types.IPTCTime)
		}

		datum.value = newSlice

	case types.IDSignedByte:
		var newSlice = make([]int8, len(values))

		for i, value := range values {
			newSlice[i] = value.(int8)
		}

		datum.value = newSlice

	case types.IDSignedLong:
		var newSlice = make([]int32, len(values))

		for i, value := range values {
			newSlice[i] = value.(int32)
		}

		datum.value = newSlice

	case types.IDSignedShort:
		var newSlice = make([]int16, len(values))

		for i, value := range values {
			newSlice[i] = value.(int16)
		}

		datum.value = newSlice

	case types.IDSignedRational, types.IDUnsignedRational:
		var newSlice = make([]*big.Rat, len(values))

		for i, value := range values {
			newSlice[i] = value.(*big.Rat)
		}

		datum.value = newSlice

	case types.IDTIFFDouble:
		var newSlice = make([]float64, len(values))

		for i, value := range values {
			newSlice[i] = value.(float64)
		}

		datum.value = newSlice

	case types.IDTIFFFloat:
		var newSlice = make([]float32, len(values))

		for i, value := range values {
			newSlice[i] = value.(float32)
		}

		datum.value = newSlice

	case types.IDUndefined:
		var newSlice = make([]byte, len(values))

		for i, value := range values {
			newSlice[i] = value.(byte)
		}

		datum.value = newSlice

	case types.IDUnsignedByte:
		var newSlice = make([]uint8, len(values))

		for i, value := range values {
			newSlice[i] = value.(uint8)
		}

		datum.value = newSlice

	case types.IDUnsignedLong:
		var newSlice = make([]uint32, len(values))

		for i, value := range values {
			newSlice[i] = value.(uint32)
		}

		datum.value = newSlice

	case types.IDUnsignedShort:
		var newSlice = make([]uint16, len(values))

		for i, value := range values {
			newSlice[i] = value.(uint16)
		}

		datum.value = newSlice

	case types.IDXMPLangAlt:
		var newSlice = make([]types.XMPLangAlt, len(values))

		for i, value := range values {
			newSlice[i] = value.(types.XMPLangAlt)
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

func newDatum(familyName, groupName, tagName string, typeId int, label, interpretedValue string,
	repeatable bool) *datumImpl {
	return &datumImpl{
		familyName:       familyName,
		groupName:        groupName,
		interpretedValue: interpretedValue,
		label:            label,
		repeatable:       repeatable,
		tagName:          tagName,
		typeId:           types.ID(typeId),
	}
}

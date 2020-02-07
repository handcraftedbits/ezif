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
	Family() types.Family
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
	HasKey(key string) bool
	Keys() []string
}

//
// Private types
//

// Datum implementation
type datumImpl struct {
	family           types.Family
	groupName        string
	interpretedValue string
	label            string
	repeatable       bool
	tagName          string
	typeId           types.ID
	value            interface{}
}

func (datum *datumImpl) Family() types.Family {
	return datum.family
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
	return string(datum.family) + "." + datum.groupName + "." + datum.tagName
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
	datumMap map[string]*datumImpl
	keys     []string
}

func (metadata *metadataImpl) Get(key string) Datum {
	return metadata.datumMap[key]
}

func (metadata *metadataImpl) HasKey(key string) bool {
	if _, ok := metadata.datumMap[key]; ok {
		return true
	}

	return false
}

func (metadata *metadataImpl) Keys() []string {
	return metadata.keys
}

func (metadata *metadataImpl) add(datum *datumImpl, values []interface{}) {
	var oldDatum = metadata.datumMap[datum.key()]
	var valuesLength = len(values)

	// IPTC metadata properties can be "repeatable" (at this time, this only applies to dates and strings), meaning that
	// the property can be defined multiple times and the values still need to be preserved.  Exif and XMP properties
	// can be repeated multiple times but we only preserve the last value.  Therefore, if the metadata property is
	// repeatable and it already exists we won't do anything here.  Later, we'll append the new value to the existing
	// array value.

	if datum.repeatable {
		if oldDatum == nil {
			metadata.datumMap[datum.key()] = datum
		} else {
			datum = oldDatum
		}

		// Decrement valuesLength because IPTC repeatable metadata property values come in one at a time and in the
		// event that this is the first value we want to make a slice of length 0 to append into (otherwise the first
		// element will be a zero value).

		valuesLength -= 1
	} else {
		metadata.datumMap[datum.key()] = datum
	}

	// We need to do a manual conversion from an interface{} slice to a concrete-typed slice.  This keeps us from having
	// to do a similar conversion later on.

	switch datum.TypeID() {
	case types.IDAsciiString, types.IDComment, types.IDXMPAlt, types.IDXMPBag, types.IDXMPSeq, types.IDXMPText:
		var slice = make([]string, valuesLength)

		for i, value := range values {
			slice[i] = value.(string)
		}

		datum.value = slice

	case types.IDIPTCDate:
		var slice []types.IPTCDate

		if !datum.repeatable || oldDatum == nil {
			slice = make([]types.IPTCDate, valuesLength)
		} else {
			slice = datum.value.([]types.IPTCDate)
		}

		for i, value := range values {
			if !datum.repeatable {
				slice[i] = value.(types.IPTCDate)
			} else {
				slice = append(slice, value.(types.IPTCDate))
			}
		}

		datum.value = slice

	case types.IDIPTCString:
		var slice []string

		if !datum.repeatable || oldDatum == nil {
			slice = make([]string, valuesLength)
		} else {
			slice = datum.value.([]string)
		}

		for i, value := range values {
			if !datum.repeatable {
				slice[i] = value.(string)
			} else {
				slice = append(slice, value.(string))
			}
		}

		datum.value = slice

	case types.IDIPTCTime:
		var slice = make([]types.IPTCTime, valuesLength)

		for i, value := range values {
			slice[i] = value.(types.IPTCTime)
		}

		datum.value = slice

	case types.IDSignedByte:
		var slice = make([]int8, valuesLength)

		for i, value := range values {
			slice[i] = value.(int8)
		}

		datum.value = slice

	case types.IDSignedLong:
		var slice = make([]int32, valuesLength)

		for i, value := range values {
			slice[i] = value.(int32)
		}

		datum.value = slice

	case types.IDSignedShort:
		var slice = make([]int16, valuesLength)

		for i, value := range values {
			slice[i] = value.(int16)
		}

		datum.value = slice

	case types.IDSignedRational, types.IDUnsignedRational:
		var slice = make([]*big.Rat, valuesLength)

		for i, value := range values {
			slice[i] = value.(*big.Rat)
		}

		datum.value = slice

	case types.IDTIFFDouble:
		var slice = make([]float64, valuesLength)

		for i, value := range values {
			slice[i] = value.(float64)
		}

		datum.value = slice

	case types.IDTIFFFloat:
		var slice = make([]float32, valuesLength)

		for i, value := range values {
			slice[i] = value.(float32)
		}

		datum.value = slice

	case types.IDUndefined:
		var slice = make([]byte, valuesLength)

		for i, value := range values {
			slice[i] = value.(byte)
		}

		datum.value = slice

	case types.IDUnsignedByte:
		var slice = make([]uint8, valuesLength)

		for i, value := range values {
			slice[i] = value.(uint8)
		}

		datum.value = slice

	case types.IDUnsignedLong:
		var slice = make([]uint32, valuesLength)

		for i, value := range values {
			slice[i] = value.(uint32)
		}

		datum.value = slice

	case types.IDUnsignedShort:
		var slice = make([]uint16, valuesLength)

		for i, value := range values {
			slice[i] = value.(uint16)
		}

		datum.value = slice

	// XMPLangAlt is a special case, there's really only a single "value", which is a map.

	case types.IDXMPLangAlt:
		var langAlt = make(map[string]string)

		for _, value := range values {
			curValue := value.(*xmpLangAltEntry)

			langAlt[curValue.language] = curValue.value
		}

		// TODO: can this be done as a single value instead of a slice?
		datum.value = []map[string]string{langAlt}
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

func newDatum(family types.Family, groupName, tagName string, typeId types.ID, label, interpretedValue string,
	repeatable bool) *datumImpl {
	return &datumImpl{
		family:           family,
		groupName:        groupName,
		interpretedValue: interpretedValue,
		label:            label,
		repeatable:       repeatable,
		tagName:          tagName,
		typeId:           typeId,
	}
}

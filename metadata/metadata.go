package metadata // import "golang.handcraftedbits.com/ezif/metadata"

import (
	"math/big"
	"sort"

	"golang.handcraftedbits.com/ezif/types"
)

//
// Public types
//

type Family string

type Collection interface {
	Exif() Properties
	IPTC() Properties
	XMP() Properties
}

type Properties interface {
	Get(key string) Property
	HasKey(key string) bool
	Keys() []string
}

type Property interface {
	Family() Family
	GroupName() string
	InterpretedValue() string
	Label() string
	TagName() string
	TypeID() types.ID
	Value() interface{}
}

//
// Public constants
//

const (
	FamilyExif Family = "Exif"
	FamilyIPTC Family = "Iptc"
	FamilyXMP  Family = "Xmp"
)

//
// Private types
//

// Collection implementation
type collectionImpl struct {
	exifProperties *propertiesImpl
	iptcProperties *propertiesImpl
	xmpProperties  *propertiesImpl
}

func (collection *collectionImpl) Exif() Properties {
	return collection.exifProperties
}

func (collection *collectionImpl) IPTC() Properties {
	return collection.iptcProperties
}

func (collection *collectionImpl) XMP() Properties {
	return collection.xmpProperties
}

// Properties implementation
type propertiesImpl struct {
	propertyMap map[string]*propertyImpl
	keys        []string
}

func (properties *propertiesImpl) Get(key string) Property {
	return properties.propertyMap[key]
}

func (properties *propertiesImpl) HasKey(key string) bool {
	if _, ok := properties.propertyMap[key]; ok {
		return true
	}

	return false
}

func (properties *propertiesImpl) Keys() []string {
	return properties.keys
}

func (properties *propertiesImpl) add(property *propertyImpl, values []interface{}) {
	var oldProperty = properties.propertyMap[property.key()]
	var valuesLength = len(values)

	// IPTC metadata properties can be "repeatable" (at this time, this only applies to dates and strings), meaning that
	// the property can be defined multiple times and the values still need to be preserved.  Exif and XMP properties
	// can be repeated multiple times but we only preserve the last value.  Therefore, if the metadata property is
	// repeatable and it already exists we won't do anything here.  Later, we'll append the new value to the existing
	// array value.

	if property.repeatable {
		if oldProperty == nil {
			properties.propertyMap[property.key()] = property
		} else {
			property = oldProperty
		}

		// Decrement valuesLength because IPTC repeatable metadata property values come in one at a time and in the
		// event that this is the first value we want to make a slice of length 0 to append into (otherwise the first
		// element will be a zero value).

		valuesLength -= 1
	} else {
		properties.propertyMap[property.key()] = property
	}

	// We need to do a manual conversion from an interface{} slice to a concrete-typed slice.  This keeps us from having
	// to do a similar conversion later on.

	switch property.TypeID() {
	case types.IDAsciiString, types.IDComment, types.IDXMPAlt, types.IDXMPBag, types.IDXMPSeq, types.IDXMPText:
		var slice = make([]string, valuesLength)

		for i, value := range values {
			slice[i] = value.(string)
		}

		property.value = slice

	case types.IDIPTCDate:
		var slice []types.IPTCDate

		if !property.repeatable || oldProperty == nil {
			slice = make([]types.IPTCDate, valuesLength)
		} else {
			slice = property.value.([]types.IPTCDate)
		}

		for i, value := range values {
			if !property.repeatable {
				slice[i] = value.(types.IPTCDate)
			} else {
				slice = append(slice, value.(types.IPTCDate))
			}
		}

		property.value = slice

	case types.IDIPTCString:
		var slice []string

		if !property.repeatable || oldProperty == nil {
			slice = make([]string, valuesLength)
		} else {
			slice = property.value.([]string)
		}

		for i, value := range values {
			if !property.repeatable {
				slice[i] = value.(string)
			} else {
				slice = append(slice, value.(string))
			}
		}

		property.value = slice

	case types.IDIPTCTime:
		var slice = make([]types.IPTCTime, valuesLength)

		for i, value := range values {
			slice[i] = value.(types.IPTCTime)
		}

		property.value = slice

	case types.IDSignedByte:
		var slice = make([]int8, valuesLength)

		for i, value := range values {
			slice[i] = value.(int8)
		}

		property.value = slice

	case types.IDSignedLong:
		var slice = make([]int32, valuesLength)

		for i, value := range values {
			slice[i] = value.(int32)
		}

		property.value = slice

	case types.IDSignedShort:
		var slice = make([]int16, valuesLength)

		for i, value := range values {
			slice[i] = value.(int16)
		}

		property.value = slice

	case types.IDSignedRational, types.IDUnsignedRational:
		var slice = make([]*big.Rat, valuesLength)

		for i, value := range values {
			slice[i] = value.(*big.Rat)
		}

		property.value = slice

	case types.IDTIFFDouble:
		var slice = make([]float64, valuesLength)

		for i, value := range values {
			slice[i] = value.(float64)
		}

		property.value = slice

	case types.IDTIFFFloat:
		var slice = make([]float32, valuesLength)

		for i, value := range values {
			slice[i] = value.(float32)
		}

		property.value = slice

	case types.IDUndefined:
		var slice = make([]byte, valuesLength)

		for i, value := range values {
			slice[i] = value.(byte)
		}

		property.value = slice

	case types.IDUnsignedByte:
		var slice = make([]uint8, valuesLength)

		for i, value := range values {
			slice[i] = value.(uint8)
		}

		property.value = slice

	case types.IDUnsignedLong:
		var slice = make([]uint32, valuesLength)

		for i, value := range values {
			slice[i] = value.(uint32)
		}

		property.value = slice

	case types.IDUnsignedShort:
		var slice = make([]uint16, valuesLength)

		for i, value := range values {
			slice[i] = value.(uint16)
		}

		property.value = slice

	// XMPLangAlt is a special case, there's really only a single "value", which is a map.

	case types.IDXMPLangAlt:
		var langAlt = make(map[string]string)

		for _, value := range values {
			curValue := value.(*xmpLangAltEntry)

			langAlt[curValue.language] = curValue.value
		}

		// TODO: can this be done as a single value instead of a slice?
		property.value = []map[string]string{langAlt}
	}
}

func (properties *propertiesImpl) finish() {
	var i = 0

	properties.keys = make([]string, len(properties.propertyMap))

	for key := range properties.propertyMap {
		properties.keys[i] = key

		i++
	}

	sort.Strings(properties.keys)
}

// Property implementation
type propertyImpl struct {
	family           Family
	groupName        string
	interpretedValue string
	label            string
	repeatable       bool
	tagName          string
	typeId           types.ID
	value            interface{}
}

func (property *propertyImpl) Family() Family {
	return property.family
}

func (property *propertyImpl) GroupName() string {
	return property.groupName
}

func (property *propertyImpl) InterpretedValue() string {
	return property.interpretedValue
}

func (property *propertyImpl) Label() string {
	return property.label
}

func (property *propertyImpl) TagName() string {
	return property.tagName
}

func (property *propertyImpl) TypeID() types.ID {
	return property.typeId
}

func (property *propertyImpl) Value() interface{} {
	return property.value
}

func (property *propertyImpl) key() string {
	return string(property.family) + "." + property.groupName + "." + property.tagName
}

type xmpLangAltEntry struct {
	language string
	value    string
}

//
// Private functions
//

func newProperty(family Family, groupName, tagName string, typeId types.ID, label, interpretedValue string,
	repeatable bool) *propertyImpl {
	return &propertyImpl{
		family:           family,
		groupName:        groupName,
		interpretedValue: interpretedValue,
		label:            label,
		repeatable:       repeatable,
		tagName:          tagName,
		typeId:           typeId,
	}
}

package imageio // import "golang.handcraftedbits.com/ezif/imageio"

import (
	"math/big"
	"sort"

	"golang.handcraftedbits.com/ezif"
)

//
// Private types
//

// ezif.Metadata implementation
type metadataImpl struct {
	exifProperties *propertiesImpl
	iptcProperties *propertiesImpl
	xmpProperties  *propertiesImpl
}

func (metadata *metadataImpl) Exif() ezif.Properties {
	return metadata.exifProperties
}

func (metadata *metadataImpl) IPTC() ezif.Properties {
	return metadata.iptcProperties
}

func (metadata *metadataImpl) XMP() ezif.Properties {
	return metadata.xmpProperties
}

// ezif.Properties implementation
type propertiesImpl struct {
	propertyMap map[string]*propertyImpl
	keys        []string
}

func (properties *propertiesImpl) Get(key string) ezif.Property {
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
	case ezif.IDAsciiString, ezif.IDComment, ezif.IDXMPAlt, ezif.IDXMPBag, ezif.IDXMPSeq, ezif.IDXMPText:
		var slice = make([]string, valuesLength)

		for i, value := range values {
			slice[i] = value.(string)
		}

		property.value = slice

	case ezif.IDIPTCDate:
		var slice []ezif.IPTCDate

		if !property.repeatable || oldProperty == nil {
			slice = make([]ezif.IPTCDate, valuesLength)
		} else {
			slice = property.value.([]ezif.IPTCDate)
		}

		for i, value := range values {
			if !property.repeatable {
				slice[i] = value.(ezif.IPTCDate)
			} else {
				slice = append(slice, value.(ezif.IPTCDate))
			}
		}

		property.value = slice

	case ezif.IDIPTCString:
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

	case ezif.IDIPTCTime:
		var slice = make([]ezif.IPTCTime, valuesLength)

		for i, value := range values {
			slice[i] = value.(ezif.IPTCTime)
		}

		property.value = slice

	case ezif.IDSignedByte:
		var slice = make([]int8, valuesLength)

		for i, value := range values {
			slice[i] = value.(int8)
		}

		property.value = slice

	case ezif.IDSignedLong:
		var slice = make([]int32, valuesLength)

		for i, value := range values {
			slice[i] = value.(int32)
		}

		property.value = slice

	case ezif.IDSignedShort:
		var slice = make([]int16, valuesLength)

		for i, value := range values {
			slice[i] = value.(int16)
		}

		property.value = slice

	case ezif.IDSignedRational, ezif.IDUnsignedRational:
		var slice = make([]*big.Rat, valuesLength)

		for i, value := range values {
			slice[i] = value.(*big.Rat)
		}

		property.value = slice

	case ezif.IDTIFFDouble:
		var slice = make([]float64, valuesLength)

		for i, value := range values {
			slice[i] = value.(float64)
		}

		property.value = slice

	case ezif.IDTIFFFloat:
		var slice = make([]float32, valuesLength)

		for i, value := range values {
			slice[i] = value.(float32)
		}

		property.value = slice

	case ezif.IDUndefined:
		var slice = make([]byte, valuesLength)

		for i, value := range values {
			slice[i] = value.(byte)
		}

		property.value = slice

	case ezif.IDUnsignedByte:
		var slice = make([]uint8, valuesLength)

		for i, value := range values {
			slice[i] = value.(uint8)
		}

		property.value = slice

	case ezif.IDUnsignedLong:
		var slice = make([]uint32, valuesLength)

		for i, value := range values {
			slice[i] = value.(uint32)
		}

		property.value = slice

	case ezif.IDUnsignedShort:
		var slice = make([]uint16, valuesLength)

		for i, value := range values {
			slice[i] = value.(uint16)
		}

		property.value = slice

	// XMPLangAlt is a special case, there's really only a single "value", which is a map.

	case ezif.IDXMPLangAlt:
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

// ezif.Property implementation
type propertyImpl struct {
	family           ezif.Family
	groupName        string
	interpretedValue string
	label            string
	repeatable       bool
	tagName          string
	typeId           ezif.ID
	value            interface{}
}

func (property *propertyImpl) Family() ezif.Family {
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

func (property *propertyImpl) TypeID() ezif.ID {
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

func newProperty(family ezif.Family, groupName, tagName string, typeId ezif.ID, label, interpretedValue string,
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

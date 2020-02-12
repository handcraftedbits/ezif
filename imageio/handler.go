package imageio // import "golang.handcraftedbits.com/ezif/imageio"

/*
#include "exiv2.h"
*/
import "C"

import (
	"math/big"
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	log "github.com/sirupsen/logrus"

	"golang.handcraftedbits.com/ezif"
	"golang.handcraftedbits.com/ezif/internal"
)

//
// Private types
//

type readHandler struct {
	index    int
	metadata *metadataImpl
	property *propertyImpl
	values   []interface{}
}

func (handler *readHandler) finish() {
	handler.metadata.exifProperties.finish()
	handler.metadata.iptcProperties.finish()
	handler.metadata.xmpProperties.finish()
}

func (handler *readHandler) onPropertyEnd(familyName string) {
	if internal.Log.IsLevelEnabled(log.DebugLevel) {
		internal.Log.WithFields(log.Fields{
			"name": string(handler.property.family) + "." + handler.property.groupName + "." + handler.property.tagName,
		}).Debug("property end")
	}

	switch ezif.Family(familyName) {
	case ezif.FamilyExif:
		handler.metadata.exifProperties.add(handler.property, handler.values)

	case ezif.FamilyIPTC:
		handler.metadata.iptcProperties.add(handler.property, handler.values)

	case ezif.FamilyXMP:
		handler.metadata.xmpProperties.add(handler.property, handler.values)
	}
}

func (handler *readHandler) onPropertyStart(familyName, groupName, tagName string, typeId int, label,
	interpretedValue string, numValues int, repeatable bool) {
	if internal.Log.IsLevelEnabled(log.DebugLevel) {
		internal.Log.WithFields(log.Fields{
			"interpretedValue": interpretedValue,
			"label":            label,
			"name":             familyName + "." + groupName + "." + tagName,
			"numValues":        numValues,
			"repeatable":       repeatable,
			"typeId":           ezif.ID(typeId),
		}).Debug("property start")
	}

	handler.index = 0
	handler.property = newProperty(ezif.Family(familyName), groupName, tagName, ezif.ID(typeId), label,
		interpretedValue, repeatable)
	handler.values = make([]interface{}, numValues)
}

func (handler *readHandler) onValue(valueHolder *C.struct_valueHolder) {
	handler.values[handler.index] = convertValueFromValueHolder(handler.property.TypeID(), valueHolder)

	if internal.Log.IsLevelEnabled(log.DebugLevel) {
		internal.Log.WithFields(log.Fields{
			"name": string(handler.property.family) + "." + handler.property.groupName + "." +
				handler.property.tagName,
			"value": handler.values[handler.index],
		}).Debug("property value encountered")
	}

	handler.index++
}

//
// Private functions
//

func convertValueFromValueHolder(typeId ezif.ID, valueHolder *C.struct_valueHolder) interface{} {
	switch typeId {
	case ezif.IDAsciiString, ezif.IDComment, ezif.IDIPTCString, ezif.IDXMPAlt, ezif.IDXMPBag, ezif.IDXMPSeq,
		ezif.IDXMPText:
		return C.GoString(valueHolder.strValue)

	case ezif.IDInvalid:
		// TODO: handle somehow?  Ignore?
		return nil

	case ezif.IDIPTCDate:
		return ezif.NewIPTCDate(int(valueHolder.yearValue), time.Month(int(valueHolder.monthValue)),
			int(valueHolder.dayValue))

	case ezif.IDIPTCTime:
		return ezif.NewIPTCTime(int(valueHolder.hourValue), int(valueHolder.minuteValue), int(valueHolder.secondValue),
			int(valueHolder.timezoneHourOffset), int(valueHolder.timezoneMinuteOffset))

	case ezif.IDSignedByte:
		return int8(valueHolder.longValue)

	case ezif.IDSignedLong:
		return int32(valueHolder.longValue)

	case ezif.IDSignedShort:
		return int16(valueHolder.longValue)

	case ezif.IDSignedRational, ezif.IDUnsignedRational:
		return big.NewRat(int64(valueHolder.rationalValueN), int64(valueHolder.rationalValueD))

	case ezif.IDTIFFDouble:
		return float64(valueHolder.doubleValue)

	case ezif.IDTIFFFloat:
		return float32(valueHolder.doubleValue)

	case ezif.IDUndefined:
		return byte(valueHolder.longValue)

	case ezif.IDUnsignedByte:
		return uint8(valueHolder.longValue)

	case ezif.IDUnsignedLong:
		return uint32(valueHolder.longValue)

	case ezif.IDUnsignedShort:
		return uint16(valueHolder.longValue)

	case ezif.IDXMPLangAlt:
		return &xmpLangAltEntry{
			language: C.GoString(valueHolder.langValue),
			value:    C.GoString(valueHolder.strValue),
		}

	default:
		return nil
	}
}

func newReadHandler() *readHandler {
	return &readHandler{
		metadata: &metadataImpl{
			exifProperties: &propertiesImpl{propertyMap: make(map[string]*propertyImpl)},
			iptcProperties: &propertiesImpl{propertyMap: make(map[string]*propertyImpl)},
			xmpProperties:  &propertiesImpl{propertyMap: make(map[string]*propertyImpl)},
		},
	}
}

//export onPropertyEndGo
func onPropertyEndGo(rhPointer unsafe.Pointer, familyName *C.char) {
	var handlers = gopointer.Restore(rhPointer).(*readHandler)

	handlers.onPropertyEnd(C.GoString(familyName))
}

//export onPropertyStartGo
func onPropertyStartGo(rhPointer unsafe.Pointer, familyName, groupName, tagName *C.char, typeId C.int,
	label, interpretedValue *C.char, numValues C.int, repeatable C.int) {
	var canRepeat bool
	var handlers = gopointer.Restore(rhPointer).(*readHandler)

	if int(repeatable) == 1 {
		canRepeat = true
	} else {
		canRepeat = false
	}

	handlers.onPropertyStart(C.GoString(familyName), C.GoString(groupName), C.GoString(tagName), int(typeId),
		C.GoString(label), C.GoString(interpretedValue), int(numValues), canRepeat)
}

//export onValueGo
func onValueGo(rhPointer unsafe.Pointer, valueHolder *C.struct_valueHolder) {
	var handlers = gopointer.Restore(rhPointer).(*readHandler)

	handlers.onValue(valueHolder)
}

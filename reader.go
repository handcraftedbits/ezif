package ezif // import "golang.handcraftedbits.com/ezif"

/*
#cgo LDFLAGS: -lexiv2 -lexiv2-xmp -lexpat -lz

#include <stdlib.h>

#include "exiv2_bridge.h"
*/
import "C"

import (
	"fmt"
	"math/big"
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	log "github.com/sirupsen/logrus"

	"golang.handcraftedbits.com/ezif/internal"
	"golang.handcraftedbits.com/ezif/types"
)

//
// Public functions
//

func ReadImageMetadata(filename string) (ImageMetadata, error) {
	var datum *datumImpl
	var err error
	var imageMetadata = &imageMetadataImpl{
		exifMetadata: &metadataImpl{datumMap: make(map[string]*datumImpl)},
		iptcMetadata: &metadataImpl{datumMap: make(map[string]*datumImpl)},
		xmpMetadata:  &metadataImpl{datumMap: make(map[string]*datumImpl)},
	}
	var index int
	var values []interface{}

	err = cReadImageMetadata(filename, &readHandlers{
		onDatumEnd: func(familyName string) {
			if internal.Log.IsLevelEnabled(log.DebugLevel) {
				internal.Log.WithFields(log.Fields{
					"name": string(datum.family) + "." + datum.groupName + "." + datum.tagName,
				}).Debug("datum end")
			}

			switch types.Family(familyName) {
			case types.FamilyExif:
				imageMetadata.exifMetadata.add(datum, values)

			case types.FamilyIPTC:
				imageMetadata.iptcMetadata.add(datum, values)

			case types.FamilyXMP:
				imageMetadata.xmpMetadata.add(datum, values)
			}
		},

		onDatumStart: func(familyName, groupName, tagName string, typeId int, label, interpretedValue string,
			numValues int, repeatable bool) {
			if internal.Log.IsLevelEnabled(log.DebugLevel) {
				internal.Log.WithFields(log.Fields{
					"interpretedValue": interpretedValue,
					"label":            label,
					"name":             familyName + "." + groupName + "." + tagName,
					"numValues":        numValues,
					"repeatable":       repeatable,
					"typeId":           types.ID(typeId),
				}).Debug("datum start")
			}

			datum = newDatum(types.Family(familyName), groupName, tagName, types.ID(typeId), label, interpretedValue,
				repeatable)
			index = 0
			values = make([]interface{}, numValues)
		},

		onValue: func(valueHolder *C.struct_valueHolder) {
			values[index] = convertValueFromValueHolder(datum.TypeID(), valueHolder)

			if internal.Log.IsLevelEnabled(log.DebugLevel) {
				internal.Log.WithFields(log.Fields{
					"name":  string(datum.family) + "." + datum.groupName + "." + datum.tagName,
					"value": values[index],
				}).Debug("datum value encountered")
			}

			index++
		},
	})

	if err != nil {
		return nil, err
	}

	imageMetadata.exifMetadata.finish()
	imageMetadata.iptcMetadata.finish()
	imageMetadata.xmpMetadata.finish()

	return imageMetadata, nil
}

//
// Private types
//

type readHandlers struct {
	onDatumEnd   func(familyName string)
	onDatumStart func(familyName, groupName, tagName string, typeId int, label, interpretedValue string, numValues int,
		repeatable bool)
	onValue func(valueHolder *C.struct_valueHolder)
}

type xmpLangAltEntry struct {
	language string
	value    string
}

//
// Private functions
//

func convertValueFromValueHolder(typeId types.ID, valueHolder *C.struct_valueHolder) interface{} {
	switch typeId {
	case types.IDAsciiString, types.IDComment, types.IDIPTCString, types.IDXMPAlt, types.IDXMPBag, types.IDXMPSeq,
		types.IDXMPText:
		return C.GoString(valueHolder.strValue)

	case types.IDInvalid:
		// TODO: handle somehow?  Ignore?
		return nil

	case types.IDIPTCDate:
		return types.NewIPTCDate(int(valueHolder.yearValue), time.Month(int(valueHolder.monthValue)),
			int(valueHolder.dayValue))

	case types.IDIPTCTime:
		return types.NewIPTCTime(int(valueHolder.hourValue), int(valueHolder.minuteValue), int(valueHolder.secondValue),
			int(valueHolder.timezoneHourOffset), int(valueHolder.timezoneMinuteOffset))

	case types.IDSignedByte:
		return int8(valueHolder.longValue)

	case types.IDSignedLong:
		return int32(valueHolder.longValue)

	case types.IDSignedShort:
		return int16(valueHolder.longValue)

	case types.IDSignedRational, types.IDUnsignedRational:
		return big.NewRat(int64(valueHolder.rationalValueN), int64(valueHolder.rationalValueD))

	case types.IDTIFFDouble:
		return float64(valueHolder.doubleValue)

	case types.IDTIFFFloat:
		return float32(valueHolder.doubleValue)

	case types.IDUndefined:
		return byte(valueHolder.longValue)

	case types.IDUnsignedByte:
		return uint8(valueHolder.longValue)

	case types.IDUnsignedLong:
		return uint32(valueHolder.longValue)

	case types.IDUnsignedShort:
		return uint16(valueHolder.longValue)

	case types.IDXMPLangAlt:
		return &xmpLangAltEntry{
			language: C.GoString(valueHolder.langValue),
			value:    C.GoString(valueHolder.strValue),
		}

	default:
		return nil
	}
}

func cReadImageMetadata(filename string, handlers *readHandlers) error {
	var cExiv2Error = C.struct_exiv2Error{
		code: C.int(-999),
	}
	var cFilename = C.CString(filename)
	var cReadHandlers = C.struct_readHandlers{
		doec: C.datumOnEndCallback(C.onDatumEnd),
		dosc: C.datumOnStartCallback(C.onDatumStart),
		vc:   C.valueCallback(C.onValue),
	}
	var cValueHolder = C.struct_valueHolder{}
	var rhPointer = gopointer.Save(handlers)

	defer C.free(unsafe.Pointer(cFilename))
	defer gopointer.Unref(rhPointer)

	if log.IsLevelEnabled(log.InfoLevel) {
		internal.Log.WithFields(log.Fields{
			"filename": filename,
		}).Info("reading image metadata from file")
	}

	C.readImageMetadata(cFilename, &cExiv2Error, &cValueHolder, &cReadHandlers, rhPointer)

	if cExiv2Error.code != C.int(-999) {
		defer C.free(unsafe.Pointer(cExiv2Error.message))

		// Not using %s because it creates a spurious warning about the argument not being a string.

		return fmt.Errorf(C.GoString(cExiv2Error.message)+" (Exiv2 error code %d)", int(cExiv2Error.code))
	}

	return nil
}

//export onDatumEndGo
func onDatumEndGo(rhPointer unsafe.Pointer, familyName *C.char) {
	var handlers = gopointer.Restore(rhPointer).(*readHandlers)

	handlers.onDatumEnd(C.GoString(familyName))
}

//export onDatumStartGo
func onDatumStartGo(rhPointer unsafe.Pointer, familyName, groupName, tagName *C.char, typeId C.int,
	label, interpretedValue *C.char, numValues C.int, repeatable C.int) {
	var canRepeat bool
	var handlers = gopointer.Restore(rhPointer).(*readHandlers)

	if int(repeatable) == 1 {
		canRepeat = true
	} else {
		canRepeat = false
	}

	handlers.onDatumStart(C.GoString(familyName), C.GoString(groupName), C.GoString(tagName), int(typeId),
		C.GoString(label), C.GoString(interpretedValue), int(numValues), canRepeat)
}

//export onValueGo
func onValueGo(rhPointer unsafe.Pointer, valueHolder *C.struct_valueHolder) {
	var handlers = gopointer.Restore(rhPointer).(*readHandlers)

	handlers.onValue(valueHolder)
}

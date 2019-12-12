package ezif // import "golang.handcraftedbits.com/ezif"

/*
#cgo LDFLAGS: -lexiv2 -lexiv2-xmp -lexpat -lintl -lz

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
)

//
// Public functions
//

func ReadImageMetadata(filename string) (ImageMetadata, error) {
	var datum *datumImpl
	var err error
	var imageMetadata = &imageMetadataImpl{
		exifMetadata: &metadataImpl{datumMap: make(map[string]Datum)},
		iptcMetadata: &metadataImpl{datumMap: make(map[string]Datum)},
		xmpMetadata:  &metadataImpl{datumMap: make(map[string]Datum)},
	}
	var index int
	var values []interface{}

	err = cReadImageMetadata(filename, &readHandlers{
		onDatumEnd: func(familyName string) {
			switch familyName {
			case familyNameExif:
				imageMetadata.exifMetadata.add(datum, values)

			case familyNameIPTC:
				imageMetadata.iptcMetadata.add(datum, values)

			case familyNameXMP:
				imageMetadata.xmpMetadata.add(datum, values)
			}
		},

		onDatumStart: func(familyName, groupName, tagName string, typeId int, label, interpretedValue string,
			numValues int) {
			if familyName == "Exif" && groupName == "Photo" && tagName == "UserComment" {
				fmt.Printf("*** datum start... typeId=%d\n", typeId)
			}
			switch familyName {
			case familyNameExif:
				datum = newDatum(familyName, groupName, tagName, typeId, label, interpretedValue)
				index = 0
				values = make([]interface{}, numValues)
			}
		},

		onValue: func(valueHolder *C.struct_valueHolder) {
			if datum.familyName == "Exif" && datum.groupName == "Photo" && datum.tagName == "UserComment" {
				fmt.Printf("*** convert... typeId=%d\n", datum.typeId)
			}
			values[index] = convertValueFromValueHolder(datum.TypeID(), valueHolder)

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
	onDatumStart func(familyName, groupName, tagName string, typeId int, label, interpretedValue string, numValues int)
	onValue      func(valueHolder *C.struct_valueHolder)
}

//
// Private constants
//

const (
	familyNameExif = "Exif"
	familyNameIPTC = "IPTC"
	familyNameXMP  = "XMP"
)

//
// Private functions
//

func convertValueFromValueHolder(typeId TypeID, valueHolder *C.struct_valueHolder) interface{} {
	switch typeId {
	case TypeIDAsciiString, TypeIDComment, TypeIDIPTCString, TypeIDXMPAlt, TypeIDXMPBag, TypeIDXMPSeq, TypeIDXMPText:
		return C.GoString(valueHolder.strValue)

	case TypeIDInvalid:
		// TODO: handle somehow?  Ignore?
		return nil

	case TypeIDIPTCDate:
		return &iptcDateImpl{
			day:   int(valueHolder.dayValue),
			month: int(valueHolder.monthValue),
			year:  int(valueHolder.yearValue),
		}

	case TypeIDIPTCTime:
		return &iptcTimeImpl{
			hour:   int(valueHolder.hourValue),
			minute: int(valueHolder.minuteValue),
			second: int(valueHolder.secondValue),
			timezone: time.FixedZone("IPTC time",
				(int(valueHolder.timezoneHourOffset)*60*60)+(int(valueHolder.timezoneMinuteOffset)*60)),
		}

	case TypeIDSignedByte:
		return int8(valueHolder.longValue)

	case TypeIDSignedLong:
		return int32(valueHolder.longValue)

	case TypeIDSignedShort:
		return int16(valueHolder.longValue)

	case TypeIDSignedRational, TypeIDUnsignedRational:
		return big.NewRat(int64(valueHolder.rationalValueN), int64(valueHolder.rationalValueD))

	case TypeIDTIFFDouble:
		return float64(valueHolder.doubleValue)

	case TypeIDTIFFFloat:
		return float32(valueHolder.doubleValue)

	case TypeIDUndefined:
		return byte(valueHolder.longValue)

	case TypeIDUnsignedByte:
		return uint8(valueHolder.longValue)

	case TypeIDUnsignedLong:
		return uint32(valueHolder.longValue)

	case TypeIDUnsignedShort:
		return uint16(valueHolder.longValue)

	case TypeIDXMPLangAlt:
		// TODO: fix
		return nil

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
	label, interpretedValue *C.char, numValues C.int) {
	var handlers = gopointer.Restore(rhPointer).(*readHandlers)

	handlers.onDatumStart(C.GoString(familyName), C.GoString(groupName), C.GoString(tagName), int(typeId),
		C.GoString(label), C.GoString(interpretedValue), int(numValues))
}

//export onValueGo
func onValueGo(rhPointer unsafe.Pointer, valueHolder *C.struct_valueHolder) {
	var handlers = gopointer.Restore(rhPointer).(*readHandlers)

	handlers.onValue(valueHolder)
}

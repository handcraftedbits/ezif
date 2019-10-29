package ezif // import "golang.handcraftedbits.com/ezif"

/*
#cgo LDFLAGS: -lexiv2 -lexiv2-xmp -lexpat -lintl -lz

#include <stdlib.h>

#include "exiv2_bridge.h"
*/
import "C"

import (
	"fmt"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

//
// Public functions
//

func ReadImageMetadata(filename string) (ImageMetadata, error) {
	var err error
	var exif = newExifMetadata()

	err = readImageMetadata(filename, &readHandlers{
		onDatumStart: func(familyName, groupName, tagName string, typeId int, label, interpretedValue string,
			numValues int) {
			if familyName == "Exif" {
				exif.add(newExifDatum(familyName, groupName, tagName, typeId, label, interpretedValue, numValues))

				exif.valueCount = 0
			}
		},

		onValue: func(familyName string, valueHolder *C.struct_valueHolder) {
			switch familyName {
			case familyNameExif:
				exif.currentDatum.populateValueFromValueHolder(exif.valueCount, valueHolder)

				exif.valueCount++
			}
		},
	})

	if err != nil {
		return nil, err
	}

	return &imageMetadataImpl{
		exifMetadata: exif,
	}, nil
}

//
// Private types
//

type readHandlers struct {
	onDatumStart func(familyName, groupName, tagName string, typeId int, label, interpretedValue string, numValues int)
	onValue      func(familyName string, valueHolder *C.struct_valueHolder)
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

//export onDatumStartGo
func onDatumStartGo(rhPointer unsafe.Pointer, familyName, groupName, tagName *C.char, typeId C.int,
	label, interpretedValue *C.char, numValues C.int) {
	var handlers = gopointer.Restore(rhPointer).(*readHandlers)

	handlers.onDatumStart(C.GoString(familyName), C.GoString(groupName), C.GoString(tagName), int(typeId),
		C.GoString(label), C.GoString(interpretedValue), int(numValues))
}

//export onValueGo
func onValueGo(rhPointer unsafe.Pointer, family *C.char, valueHolder *C.struct_valueHolder) {
	var handlers = gopointer.Restore(rhPointer).(*readHandlers)

	handlers.onValue(C.GoString(family), valueHolder)
}

func readImageMetadata(filename string, handlers *readHandlers) error {
	var cExiv2Error = C.struct_exiv2Error{
		code: C.int(-999),
	}
	var cFilename = C.CString(filename)
	var cReadHandlers = C.struct_readHandlers{
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

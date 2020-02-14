package metadata // import "golang.handcraftedbits.com/ezif/metadata"

/*
#include <stdlib.h>

#include "exiv2.h"
*/
import "C"

import (
	"fmt"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	log "github.com/sirupsen/logrus"

	"golang.handcraftedbits.com/ezif/internal"
)

//
// Public functions
//

func FromFile(filename string) (Collection, error) {
	return readCollection(func(cExiv2Error *C.struct_exiv2Error, cValueHolder *C.struct_valueHolder,
		cReadHandler *C.struct_readHandler, rhPointer unsafe.Pointer) {
		var cFilename = C.CString(filename)

		defer C.free(unsafe.Pointer(cFilename))

		if log.IsLevelEnabled(log.InfoLevel) {
			internal.Log.WithFields(log.Fields{
				"filename": filename,
			}).Info("reading image metadata from file")
		}

		C.readCollectionFromFile(cFilename, cExiv2Error, cValueHolder, cReadHandler, rhPointer)
	})
}

func FromURL(url string) (Collection, error) {
	return readCollection(func(cExiv2Error *C.struct_exiv2Error, cValueHolder *C.struct_valueHolder,
		cReadHandler *C.struct_readHandler, rhPointer unsafe.Pointer) {
		var cURL = C.CString(url)

		defer C.free(unsafe.Pointer(cURL))

		if log.IsLevelEnabled(log.InfoLevel) {
			internal.Log.WithFields(log.Fields{
				"url": url,
			}).Info("reading image metadata from URL")
		}

		C.readCollectionFromURL(cURL, cExiv2Error, cValueHolder, cReadHandler, rhPointer)
	})
}

//
// Private types
//

type readCollectionInvoker func(cExiv2Error *C.struct_exiv2Error, cValueHolder *C.struct_valueHolder,
	cReadHandler *C.struct_readHandler, rhPointer unsafe.Pointer)

//
// Private functions
//

func cReadCollection(invoker readCollectionInvoker, handler *readHandler) error {
	var cExiv2Error = C.struct_exiv2Error{
		code: C.int(-999),
	}
	var cReadHandler = C.struct_readHandler{
		poec: C.propertyOnEndCallback(C.onPropertyEnd),
		posc: C.propertyOnStartCallback(C.onPropertyStart),
		vc:   C.valueCallback(C.onValue),
	}
	var cValueHolder = C.struct_valueHolder{}
	var rhPointer = gopointer.Save(handler)

	defer gopointer.Unref(rhPointer)

	invoker(&cExiv2Error, &cValueHolder, &cReadHandler, rhPointer)

	if cExiv2Error.code != C.int(-999) {
		defer C.free(unsafe.Pointer(cExiv2Error.message))

		// Not using %s because it creates a spurious warning about the argument not being a string.

		return fmt.Errorf(C.GoString(cExiv2Error.message)+" (Exiv2 error code %d)", int(cExiv2Error.code))
	}

	return nil
}

func readCollection(invoker readCollectionInvoker) (Collection, error) {
	var handler = newReadHandler()

	if err := cReadCollection(invoker, handler); err != nil {
		return nil, err
	}

	handler.finish()

	return handler.metadata, nil
}

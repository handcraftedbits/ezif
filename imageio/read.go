package imageio // import "golang.handcraftedbits.com/ezif/imageio"

/*
#include <stdlib.h>

#include "exiv2.h"
*/
import "C"

import (
	"fmt"
	"io"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	log "github.com/sirupsen/logrus"

	"golang.handcraftedbits.com/ezif"
	"golang.handcraftedbits.com/ezif/internal"
)

//
// Public types
//

type MetadataSource interface {
	FromFile(filename string) (ezif.Metadata, error)
	FromReader(reader io.ReadCloser) (ezif.Metadata, error)
	FromURL(url string) (ezif.Metadata, error)
}

//
// Public functions
//

func ReadMetadata() MetadataSource {
	return &metadataSourceImpl{
		handler: newReadHandler(),
	}
}

//
// Private types
//

// MetadataSource implementation
type metadataSourceImpl struct {
	handler *readHandler
}

func (metadataSource *metadataSourceImpl) FromFile(filename string) (ezif.Metadata, error) {
	if err := cReadImageMetadata(filename, metadataSource.handler); err != nil {
		return nil, err
	}

	metadataSource.handler.finish()

	return metadataSource.handler.metadata, nil
}

func (metadataSource *metadataSourceImpl) FromReader(reader io.ReadCloser) (ezif.Metadata, error) {
	return nil, nil
}

func (metadataSource *metadataSourceImpl) FromURL(url string) (ezif.Metadata, error) {
	return nil, nil
}

//
// Private functions
//

func cReadImageMetadata(filename string, handler *readHandler) error {
	var cExiv2Error = C.struct_exiv2Error{
		code: C.int(-999),
	}
	var cFilename = C.CString(filename)
	var cReadHandler = C.struct_readHandler{
		poec: C.propertyOnEndCallback(C.onPropertyEnd),
		posc: C.propertyOnStartCallback(C.onPropertyStart),
		vc:   C.valueCallback(C.onValue),
	}
	var cValueHolder = C.struct_valueHolder{}
	var rhPointer = gopointer.Save(handler)

	defer C.free(unsafe.Pointer(cFilename))
	defer gopointer.Unref(rhPointer)

	if log.IsLevelEnabled(log.InfoLevel) {
		internal.Log.WithFields(log.Fields{
			"filename": filename,
		}).Info("reading image metadata from file")
	}

	C.readImageMetadata(cFilename, &cExiv2Error, &cValueHolder, &cReadHandler, rhPointer)

	if cExiv2Error.code != C.int(-999) {
		defer C.free(unsafe.Pointer(cExiv2Error.message))

		// Not using %s because it creates a spurious warning about the argument not being a string.

		return fmt.Errorf(C.GoString(cExiv2Error.message)+" (Exiv2 error code %d)", int(cExiv2Error.code))
	}

	return nil
}

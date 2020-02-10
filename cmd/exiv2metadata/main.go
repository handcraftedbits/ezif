package main // import "golang.handcraftedbits.com/ezif/cmd/exiv2metadata"

/*
#cgo LDFLAGS: -lexiv2 -lexiv2-xmp -lexpat -ljansson -lz

#include "exiv2metadata.h"
*/
import "C"

//
// Private functions
//

func main() {
	C.dumpMetadataToJSON()
}

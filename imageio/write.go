package imageio // import "golang.handcraftedbits.com/ezif/imageio"

import (
	"io"

	"golang.handcraftedbits.com/ezif"
)

//
// Public types
//

type DestinationContents interface {
	AsImage() error
	AsJSON() error
}

type MetadataDestination interface {
	ToFile(filename string) DestinationContents
	ToWriter(writer io.WriteCloser) DestinationContents
}

//
// Public functions
//

func WriteMetadata(metadata ezif.Metadata) MetadataDestination {
	return nil
}

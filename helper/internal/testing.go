package internal // import "golang.handcraftedbits.com/ezif/helper/internal"

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"golang.handcraftedbits.com/ezif"
	"golang.handcraftedbits.com/ezif/helper"
)

//
// Public types
//

type GeneratedTestContext struct {
	AccessorFunc func(ezif.ImageMetadata) helper.Accessor
	MaxValues    interface{}
	Name         string
}

//
// Public functions
//

func GeneratedTests(t *testing.T, context *GeneratedTestContext) {
	var exiv2 = newExternalExiv2(testPNG)

	t.Run("MaxValue", func(t *testing.T) {
		testGetValueFromHelper(t, exiv2, context.Name, context.AccessorFunc, context.MaxValues)
	})

	t.Run("MissingValue", func(t *testing.T) {
		testGetMissingValueFromHelper(t, context.Name, context.AccessorFunc)
	})
}

//
// Private variables
//

var (
	// Simple 1x1 PNG image for metadata testing.
	testPNG = []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 1, 0, 0, 0, 1, 1, 3, 0, 0,
		0, 37, 219, 86, 202, 0, 0, 0, 6, 80, 76, 84, 69, 0, 0, 0, 255, 255, 255, 165, 217, 159, 221, 0, 0, 0, 9, 112,
		72, 89, 115, 0, 0, 14, 196, 0, 0, 14, 196, 1, 149, 43, 14, 27, 0, 0, 0, 10, 73, 68, 65, 84, 8, 153, 99, 96, 0,
		0, 0, 2, 0, 1, 244, 113, 100, 166, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130}
)

//
// Private functions
//

func getMethodFromAccessor(accessor helper.Accessor, name string) reflect.Value {
	return reflect.ValueOf(accessor).MethodByName(name)
}

func getRawValueFromAccessor(accessor helper.Accessor) interface{} {
	var method = getMethodFromAccessor(accessor, "Raw")
	var result = method.Call([]reflect.Value{})

	return result[0].Interface()
}

func normalizeValuesAsSlice(values interface{}) []interface{} {
	var kind reflect.Kind
	var value = reflect.ValueOf(values)

	kind = value.Kind()

	if kind == reflect.Array || kind == reflect.Slice {
		var result = make([]interface{}, value.Len())

		// Already an array/slice, so just convert to []interface{}.

		for i := 0; i < value.Len(); i++ {
			result[i] = value.Index(i).Interface()
		}

		return result
	}

	return []interface{}{values}
}

func testGetMissingValueFromHelper(t *testing.T, name string, accessorFunc func(ezif.ImageMetadata) helper.Accessor) {
	var err error
	var imageFilename string
	var metadata ezif.ImageMetadata

	imageFilename, err = saveImage(testPNG)

	require.Nil(t, err, "could not save temporary dummy image")

	defer func() {
		_ = os.Remove(imageFilename)
	}()

	// Make sure we cannot find the metadata in the temporary image.

	metadata, err = ezif.ReadImageMetadata(imageFilename)

	require.Nil(t, err)
	require.Nil(t, accessorFunc(metadata), "expected not to find metadata with name '%s' in test image", name)
}

func testGetValueFromHelper(t *testing.T, exiv2 *externalExiv2Impl, name string,
	accessorFunc func(ezif.ImageMetadata) helper.Accessor, valuesToSet interface{}) {
	var err error
	var metadata ezif.ImageMetadata
	var result interface{}
	var stdErr string
	var stdOut string

	// Write the metadata using an external copy of Exiv2 that's known to produce good results...

	exiv2.Set(name, normalizeValuesAsSlice(valuesToSet))

	err, stdOut, stdErr = exiv2.execute()

	require.Nil(t, err, "could not save metadata with name '%s' via external Exiv2 command\nstdout: %s\nstderr: %s\n",
		name, stdOut, stdErr)

	// ...and make sure we can read back the exact same values that we provided.

	metadata, err = ezif.ReadImageMetadata(exiv2.tempFilename)

	require.Nil(t, err)

	result = getRawValueFromAccessor(accessorFunc(metadata))

	require.NotNil(t, result, "couldn't find metadata with name '%s' in test image", name)
	require.Exactly(t, valuesToSet, result, "expected value(s) and results do not match")
}

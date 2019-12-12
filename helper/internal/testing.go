package internal // import "golang.handcraftedbits.com/ezif/helper/internal"

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"golang.handcraftedbits.com/ezif"
	"golang.handcraftedbits.com/ezif/helper"
)

//
// Public constants
//

const (
	SubTestMaxValue     = "MaxValue"
	SubTestMissingValue = "MissingValue"
)

//
// Public functions
//

func TestGetMissingValueFromHelper(t *testing.T, key string, accessorFunc func(ezif.ImageMetadata) helper.Accessor) {
	var err error
	var metadata ezif.ImageMetadata
	var tempFile *os.File

	tempFile, err = saveTempDummyImage()

	require.Nil(t, err, "could not save temporary dummy image")

	defer os.Remove(tempFile.Name())

	// Make sure we cannot find the key in the temporary image.

	metadata, err = ezif.ReadImageMetadata(tempFile.Name())

	require.Nil(t, err)

	require.Nil(t, accessorFunc(metadata), "expected not to find metadata with key '%s' in test image", key)
}

func TestGetValueFromHelper(t *testing.T, key string, getFunc func(ezif.ImageMetadata) interface{},
	valuesToSet interface{}) interface{} {
	var err error
	var metadata ezif.ImageMetadata
	var normalizedValues = normalizeValuesAsSlice(valuesToSet)
	var result interface{}
	var stdErr string
	var stdOut string
	var tempFile *os.File

	tempFile, err = saveTempDummyImage()

	require.Nil(t, err, "could not save temporary dummy image")

	//defer os.Remove(tempFile.Name())

	// Write the metadata using an external copy of Exiv2 that's known to produce good results...

	err, stdOut, stdErr = setMetadataViaExternalExiv2(key, tempFile.Name(), normalizedValues)

	require.Nil(t, err, "could not save metadata with key '%s' via external Exiv2 command\nstdout: %s\nstderr: %s\n",
		key, stdOut, stdErr)

	// ...and make sure we can read back the exact same values that we provided.

	metadata, err = ezif.ReadImageMetadata(tempFile.Name())

	require.Nil(t, err)

	result = getFunc(metadata)

	require.NotNil(t, result, "couldn't find metadata with key '%s' in test image", key)

	if len(normalizedValues) == 1 {
		// In order for require.Exactly() to work we need to compare against the single value in the array (since this
		// must mean that the result is a single value, not an array).

		require.Exactly(t, valuesToSet, result, "expected value and result do not match")
	} else {
		require.Exactly(t, valuesToSet, result, "expected values and results do not match")
	}

	return result
}

//
// Private variables
//

var (
	// Simple 1x1 PNG image for metadata testing.
	dummyImage = []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 1, 0, 0, 0, 1, 1, 3, 0,
		0, 0, 37, 219, 86, 202, 0, 0, 0, 6, 80, 76, 84, 69, 0, 0, 0, 255, 255, 255, 165, 217, 159, 221, 0, 0, 0, 9, 112,
		72, 89, 115, 0, 0, 14, 196, 0, 0, 14, 196, 1, 149, 43, 14, 27, 0, 0, 0, 10, 73, 68, 65, 84, 8, 153, 99, 96, 0,
		0, 0, 2, 0, 1, 244, 113, 100, 166, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130}
)

//
// Private functions
//

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

func saveTempDummyImage() (*os.File, error) {
	var err error
	var tempFile *os.File

	tempFile, err = ioutil.TempFile("", "ezif-test")

	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(tempFile.Name(), dummyImage, os.ModePerm)

	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

func setMetadataViaExternalExiv2(key, filename string, values []interface{}) (error, string, string) {
	var cmd = exec.Command("exiv2", "-M", fmt.Sprintf("set %s %s", key, valuesToExiv2Format(values)), filename)
	var stdErr bytes.Buffer
	var stdOut bytes.Buffer

	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	fmt.Printf("*** cmd: %v\n", cmd)

	if err := cmd.Run(); err != nil {
		return err, stdOut.String(), stdErr.String()
	}

	return nil, "", ""
}

func valuesToExiv2Format(values []interface{}) string {
	var buffer bytes.Buffer

	for i, value := range values {
		switch v := value.(type) {
		case *big.Rat:
			buffer.WriteString(fmt.Sprintf("%v/%v", v.Num(), v.Denom()))
		default:
			buffer.WriteString(fmt.Sprintf("%v", v))
		}

		if i < len(values)-1 {
			buffer.WriteRune(' ')
		}
	}

	return buffer.String()
}

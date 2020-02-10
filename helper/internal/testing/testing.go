package testing // import "golang.handcraftedbits.com/ezif/internal/testing"

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"golang.handcraftedbits.com/ezif"
	"golang.handcraftedbits.com/ezif/helper"
	"golang.handcraftedbits.com/ezif/types"
)

//
// Public types
//

type GeneratedTestContext struct {
	AccessorFunc func(ezif.ImageMetadata) helper.Accessor
	Family       types.Family
	IsSlice      bool
	Name         string
	TypeID       types.ID
}

//
// Public functions
//

func GeneratedTests(t *testing.T, context *GeneratedTestContext) {
	var exiv2 = newExternalExiv2(testPNG)
	var sliceLength int

	if context.IsSlice {
		sliceLength = defaultSliceLength
	} else {
		sliceLength = 1
	}

	t.Run("MaxValue", func(t *testing.T) {
		testGetValueFromHelper(t, exiv2, context, makeSlice(context.TypeID, sliceLength, maxValue))
	})

	t.Run("MissingValue", func(t *testing.T) {
		testGetMissingValueFromHelper(t, context)
	})
}

//
// Private types
//

type typeInfo struct {
	maxValue   interface{}
	minValue   interface{}
	emptyValue interface{}
}

type xmpLangAltEntry struct {
	language string
	value    string
}

//
// Private variables
//

var (
	alphabet = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's',
		't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	defaultSliceLength = 8
	defaultValueLength = 32

	// Simple 1x1 PNG image for metadata testing.
	testPNG = []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 1, 0, 0, 0, 1, 1, 3, 0, 0,
		0, 37, 219, 86, 202, 0, 0, 0, 6, 80, 76, 84, 69, 0, 0, 0, 255, 255, 255, 165, 217, 159, 221, 0, 0, 0, 9, 112,
		72, 89, 115, 0, 0, 14, 196, 0, 0, 14, 196, 1, 149, 43, 14, 27, 0, 0, 0, 10, 73, 68, 65, 84, 8, 153, 99, 96, 0,
		0, 0, 2, 0, 1, 244, 113, 100, 166, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130}

	// TODO: 0 for empty numbers, right or wrong?  Maybe see if there's a better way to handle empty values.
	// TODO: pretty sure max/min for rationals isn't right.
	// TODO: for randomStringOfLength(), should probably make a function that pre-generates a bunch of long random
	//   strings, and have each invocation cycle through them.
	typeInfos = map[types.ID]typeInfo{
		types.IDAsciiString:      {randomStringOfLength(defaultValueLength), randomStringOfLength(1), ""},
		types.IDComment:          {randomStringOfLength(defaultValueLength), randomStringOfLength(1), ""},
		types.IDIPTCDate:         {types.NewIPTCDate(9999, 12, 31), types.NewIPTCDate(1, 1, 1), nil},
		types.IDIPTCString:       {randomStringOfLength(defaultValueLength), randomStringOfLength(1), ""},
		types.IDIPTCTime:         {types.NewIPTCTime(23, 59, 59, 0, 0), types.NewIPTCTime(0, 0, 0, 0, 0), nil},
		types.IDSignedByte:       {int8(math.MaxInt8), int8(math.MinInt8), int8(0)},
		types.IDSignedLong:       {int32(math.MaxInt32), int32(math.MinInt32), int32(0)},
		types.IDSignedRational:   {big.NewRat(math.MaxInt32, 1), big.NewRat(1, math.MaxInt32), nil},
		types.IDSignedShort:      {int16(math.MaxInt16), int16(math.MinInt16), int16(0)},
		types.IDTIFFDouble:       {9.0e99, -9.0e99, float64(0)},
		types.IDTIFFFloat:        {float32(3.4e38), float32(-3.4e38), float32(0)},
		types.IDUndefined:        {byte(math.MaxUint8), byte(0), byte(0)},
		types.IDUnsignedByte:     {uint8(math.MaxUint8), uint8(0), uint8(0)},
		types.IDUnsignedLong:     {uint32(math.MaxUint32), uint32(0), uint32(0)},
		types.IDUnsignedRational: {big.NewRat(math.MaxUint32, 1), big.NewRat(1, math.MaxUint32), nil},
		types.IDUnsignedShort:    {uint16(math.MaxUint16), uint16(0), uint16(0)},
		types.IDXMPAlt:           {randomStringOfLength(defaultValueLength), randomStringOfLength(1), ""},
		types.IDXMPBag:           {randomStringOfLength(defaultValueLength), randomStringOfLength(1), ""},
		types.IDXMPSeq:           {randomStringOfLength(defaultValueLength), randomStringOfLength(1), ""},
		types.IDXMPText:          {randomStringOfLength(defaultValueLength), randomStringOfLength(1), ""},
		types.IDXMPLangAlt: {xmpLangAltEntry{"en", randomStringOfLength(defaultValueLength)},
			xmpLangAltEntry{"en", randomStringOfLength(1)}, nil},
	}
)

//
// Private functions
//

func emptyValue(typeID types.ID) interface{} {
	return typeInfos[typeID].emptyValue
}

func expectEqualValues(t *testing.T, typeID types.ID, expected []interface{}, actual []interface{}) {
	require.NotNil(t, actual)
	require.Equal(t, len(expected), len(actual))

	for i := 0; i < len(expected); i++ {
		switch typeID {
		case types.IDSignedRational, types.IDUnsignedRational:
			// Unfortunately we can't use require.Equal() directly on the two lists because big.Rat values can't
			// necessarily be compared on a field-by-field basis -- there seems to be some sort of constant value
			// replacement for the denominator that occurs (seemingly at random!) that throws everything off.  So we
			// have to use big.Rat.Cmp() in that case.

			require.True(t, expected[i].(*big.Rat).Cmp(actual[i].(*big.Rat)) == 0, fmt.Sprintf("value at index %d "+
				"does not equal expected value", i))

		case types.IDXMPLangAlt:
			var entry = expected[i].(xmpLangAltEntry)
			var resultMap = actual[0].(map[string]string)

			// The simple lang alt entry we use in testing is not the same as types.XMPLangAlt, so a manual conversion
			// is necessary.

			require.NotNil(t, resultMap[entry.language], fmt.Sprintf("XMP lang alt does not contain value for "+
				"language '%s'", entry.language))
			require.Equal(t, entry.value, resultMap[entry.language], fmt.Sprintf("XMP lang alt value for language "+
				"'%s' does not match expected value", entry.language))

		default:
			require.Equal(t, expected[i], actual[i], fmt.Sprintf("value at index %d does not equal expected value", i))
		}
	}
}

func getMethodFromAccessor(accessor helper.Accessor, name string) reflect.Value {
	var value = reflect.ValueOf(accessor)

	if value == reflect.ValueOf(nil) {
		return value
	}

	return value.MethodByName(name)
}

func getRawValueFromAccessor(accessor helper.Accessor) []interface{} {
	var kind reflect.Kind
	var method = getMethodFromAccessor(accessor, "Raw")

	if method == reflect.ValueOf(nil) {
		return nil
	}

	var result = method.Call([]reflect.Value{})
	var value = result[0]

	kind = value.Kind()

	if kind == reflect.Array || kind == reflect.Slice {
		var intfSlice = make([]interface{}, value.Len())

		// Value is already an array/slice, so just convert to []interface{}.

		for i := 0; i < value.Len(); i++ {
			intfSlice[i] = value.Index(i).Interface()
		}

		return intfSlice
	}

	return []interface{}{value.Interface()}
}

func makeSlice(typeID types.ID, length int, valueFunc func(types.ID) interface{}) []interface{} {
	var result = make([]interface{}, length)

	for i := 0; i < length; i++ {
		result[i] = valueFunc(typeID)
	}

	return result
}

func maxValue(typeID types.ID) interface{} {
	return typeInfos[typeID].maxValue
}

func minValue(typeID types.ID) interface{} {
	return typeInfos[typeID].minValue
}

func randomStringOfLength(length int) string {
	var result = make([]rune, length)

	for i := 0; i < length; i++ {
		result[i] = alphabet[rand.Int()%len(alphabet)]
	}

	return string(result)
}

func testGetMissingValueFromHelper(t *testing.T, context *GeneratedTestContext) {
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
	require.Nil(t, context.AccessorFunc(metadata), "expected not to find metadata with name '%s' in test image",
		context.Name)
}

// TODO: no maxBytes/minBytes support!
func testGetValueFromHelper(t *testing.T, exiv2 *externalExiv2Impl, context *GeneratedTestContext,
	valuesToSet []interface{}) {
	var err error
	var metadata ezif.ImageMetadata
	var result []interface{}
	var stdErr string
	var stdOut string

	// Write the metadata using an external copy of Exiv2 that's known to produce good results...

	if context.IsSlice && (context.Family == types.FamilyIPTC) {
		if context.TypeID == types.IDUndefined {
			// An undefined IPTC type needs to be set with a single "set" command.

			exiv2.Set(context.Name, valuesToSet)
		} else {
			// An IPTC value that's marked as "repeatable" is a special case.  We can't just provide all the values in a
			// single "set" command, we have to "add" the metadata property with a single value repeatedly.

			for _, value := range valuesToSet {
				exiv2.Add(context.Name, []interface{}{value})
			}
		}
	} else {
		if context.Family == types.FamilyXMP {
			// For XMP values, we need to use multiple "set" commands.

			for _, value := range valuesToSet {
				exiv2.Set(context.Name, []interface{}{value})
			}
		} else {
			// For Exif values, we have to do everything in a single "set" command.

			exiv2.Set(context.Name, valuesToSet)
		}
	}

	err, stdOut, stdErr = exiv2.execute()

	require.Nil(t, err, "could not save metadata with name '%s' via external Exiv2 command\nstdout: %s\nstderr: %s\n",
		context.Name, stdOut, stdErr)

	// ...and make sure we can read back the exact same values that we provided.

	metadata, err = ezif.ReadImageMetadata(exiv2.tempFilename)

	require.Nil(t, err)

	result = getRawValueFromAccessor(context.AccessorFunc(metadata))

	require.NotNil(t, result, "couldn't find metadata with name '%s' in test image", context.Name)

	expectEqualValues(t, context.TypeID, valuesToSet, result)
}

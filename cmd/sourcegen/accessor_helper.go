package main // import "golang.handcraftedbits.com/ezif/cmd/sourcegen"

import (
	"bytes"
	"go/format"
	"sort"
	"strings"
	"text/template"
)

//
// Private types
//

type accessorInfo struct {
	ImplName string
	IsSlice  bool
	Name     string
	Type     string
}

type accessorTemplateContext struct {
	AccessorInfos []accessorInfo
	PackageName   string
}

//
// Private variables
//

var accessors = map[string]string{
	"Date":                  "types.IPTCDate",
	"DateSlice":             "[]types.IPTCDate",
	"Double":                "float64",
	"DoubleSlice":           "[]float64",
	"Float":                 "float32",
	"FloatSlice":            "[]float32",
	"SignedByte":            "int8",
	"SignedByteSlice":       "[]int8",
	"SignedLong":            "int32",
	"SignedLongSlice":       "[]int32",
	"SignedRational":        "*big.Rat",
	"SignedRationalSlice":   "[]*big.Rat",
	"SignedShort":           "int16",
	"SignedShortSlice":      "[]int16",
	"String":                "string",
	"StringSlice":           "[]string",
	"Time":                  "types.IPTCTime",
	"Undefined":             "byte",
	"UndefinedSlice":        "[]byte",
	"UnsignedByte":          "uint8",
	"UnsignedByteSlice":     "[]uint8",
	"UnsignedLong":          "uint32",
	"UnsignedLongSlice":     "[]uint32",
	"UnsignedRational":      "*big.Rat",
	"UnsignedRationalSlice": "[]*big.Rat",
	"UnsignedShort":         "uint16",
	"UnsignedShortSlice":    "[]uint16",
	"XMPLangAlt":            "map[string]string",
}

//
// Private functions
//

func generateAccessorSource(packageName string, templateBody string) (string, error) {
	var buffer bytes.Buffer
	var err error
	var formattedSource []byte
	var templateRoot *template.Template

	templateRoot, err = initTemplate("root", templateBody)

	if err != nil {
		return "", err
	}

	err = templateRoot.Execute(&buffer, &accessorTemplateContext{
		AccessorInfos: getAccessorInfos(accessors),
		PackageName:   packageName,
	})

	if err != nil {
		return "", err
	}

	// Do a gofmt pass on the generated source.

	formattedSource, err = format.Source(buffer.Bytes())

	if err != nil {
		return "", err
	}

	return string(formattedSource), nil
}

func getAccessorInfos(accessors map[string]string) []accessorInfo {
	var result []accessorInfo

	for name, goType := range accessors {
		var implName string
		var nameRunes = []rune(name)

		implName = strings.ToLower(string(nameRunes[0])) + string(nameRunes[1:])

		if strings.HasPrefix(implName, "xMP") {
			implName = "xmp" + string(nameRunes[3:])
		}

		result = append(result, accessorInfo{
			ImplName: implName,
			IsSlice:  strings.HasPrefix(goType, "[]"),
			Name:     name,
			Type:     goType,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i].Name, result[j].Name) < 0
	})

	return result
}

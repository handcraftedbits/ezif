package main // import "golang.handcraftedbits.com/ezif/cmd/sourcegen"

import (
	"bytes"
	"fmt"
	"go/format"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

//
// Private types
//

type families map[string]family
type family map[string]group

type functionInfo struct {
	FullTagName string
	Tag         tag
}
type group map[string]tag

type groupConfig struct {
	Description string `yaml:"description"`
	Family      string `yaml:"family"`
	Reference   string `yaml:"reference"`
	Regexp      string `yaml:"regexp"`

	disabledHelpers map[string]bool
	disabledTests   map[string]bool
	regexp          *regexp.Regexp
}

type helperTemplateContext struct {
	DisabledHelpers    map[string]bool
	DisabledTests      map[string]bool
	FamilyName         string
	FunctionMappings   map[string]functionInfo
	FunctionNames      []string
	PackageDescription string
	PackageName        string
	Reference          string
}

type tag struct {
	Count       int    `json:"count"`
	Description string `json:"description"`
	Label       string `json:"label"`
	MaxBytes    int    `json:"maxBytes"`
	MinBytes    int    `json:"minBytes"`
	Repeatable  bool   `json:"repeatable"`
	TypeID      int    `json:"typeId"`
}

type typeIdMapping struct {
	goType              string
	requiredTestImports []string
	returnType          string
}

//
// Private constants
//

// Redefining these since we don't want this build tool to depend on the ezif package.
const (
	typeIdUnsignedByte     int = 1
	typeIdAsciiString      int = 2
	typeIdUnsignedShort    int = 3
	typeIdUnsignedLong     int = 4
	typeIdUnsignedRational int = 5
	typeIdSignedByte       int = 6
	typeIdUndefined        int = 7
	typeIdSignedShort      int = 8
	typeIdSignedLong       int = 9
	typeIdSignedRational   int = 10
	typeIdTIFFFloat        int = 11
	typeIdTIFFDouble       int = 12
	typeIdIPTCString       int = 0x10000
	typeIdIPTCDate         int = 0x10001
	typeIdIPTCTime         int = 0x10002
	typeIdComment          int = 0x10003
	typeIdXMPText          int = 0x10005
	typeIdXMPAlt           int = 0x10006
	typeIdXMPBag           int = 0x10007
	typeIdXMPSeq           int = 0x10008
	typeIdXMPLangAlt       int = 0x10009
)

//
// Private variables
//

var funcMap = template.FuncMap{
	"FixDescription":  templateFuncFixDescription,
	"IsHelperEnabled": templateFuncIsHelperEnabled,
	"IsTestEnabled":   templateFuncIsTestEnabled,
	"LastPackage":     templateFuncLastPackage,
	"MaxValue":        templateFuncMaxValue,
	"MinValue":        templateFuncMinValue,
	"RequiredImports": templateFuncRequiredImports,
	"ReturnType":      templateFuncReturnType,
}

var typeIdMappings = map[int]typeIdMapping{
	typeIdAsciiString:      {"string", nil, "String"},
	typeIdComment:          {"string", nil, "String"},
	typeIdIPTCDate:         {"time.Time", nil, "Date"},
	typeIdIPTCString:       {"string", nil, "String"},
	typeIdIPTCTime:         {"time.Time", nil, "Time"},
	typeIdSignedByte:       {"int8", []string{"math"}, "SignedByte"},
	typeIdSignedLong:       {"int32", []string{"math"}, "SignedLong"},
	typeIdSignedRational:   {"*big.Rat", []string{"math", "math/big"}, "SignedRational"},
	typeIdSignedShort:      {"int16", []string{"math"}, "SignedShort"},
	typeIdTIFFDouble:       {"float64", nil, "Double"},
	typeIdTIFFFloat:        {"float32", nil, "Float"},
	typeIdUndefined:        {"byte", []string{"math"}, "Undefined"},
	typeIdUnsignedByte:     {"uint8", []string{"math"}, "UnsignedByte"},
	typeIdUnsignedLong:     {"uint32", []string{"math"}, "UnsignedLong"},
	typeIdUnsignedRational: {"*big.Rat", []string{"math", "math/big"}, "UnsignedRational"},
	typeIdUnsignedShort:    {"uint16", []string{"math"}, "UnsignedShort"},
	typeIdXMPAlt:           {"[]string", nil, "StringSlice"},
	typeIdXMPBag:           {"[]string", nil, "StringSlice"},
	typeIdXMPLangAlt:       {"[]ezif.XMPLangAlt", nil, "XMPLangAlt"},
	typeIdXMPSeq:           {"[]string", nil, "StringSlice"},
	typeIdXMPText:          {"string", nil, "String"},
}

//
// Private functions
//

func generateGroupSource(familyName string, f family, packageName string, gc groupConfig) (string, error) {
	return generateSource(familyName, f, packageName, gc, templateGroupSource)
}

func generateSource(familyName string, f family, packageName string, gc groupConfig,
	templateBody string) (string, error) {
	var buffer bytes.Buffer
	var err error
	var formattedSource []byte
	var functionMappings map[string]functionInfo
	var functionNames []string
	var groupNames []string
	var templateRoot *template.Template

	// Find matching groups and make sure we iterate over them in order later.

	for group := range f {
		if gc.regexp.MatchString(group) {
			groupNames = append(groupNames, group)
		}
	}

	templateRoot, err = initTemplate("root", templateBody)

	if err != nil {
		return "", err
	}

	functionNames, functionMappings = getFunctionMappings(familyName, f, groupNames)

	err = templateRoot.Execute(&buffer, &helperTemplateContext{
		DisabledHelpers:    gc.disabledHelpers,
		DisabledTests:      gc.disabledTests,
		FamilyName:         familyName,
		FunctionMappings:   functionMappings,
		FunctionNames:      functionNames,
		PackageDescription: gc.Description,
		PackageName:        strings.ToLower(packageName),
		Reference:          gc.Reference,
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

func getAdjustedCount(familyName string, info functionInfo) int {
	switch familyName {
	case "Exif":
		var typeId = info.Tag.TypeID

		// Exiv2 seems to treat the ASCII string and comment types as characters instead of full strings because the
		// count values range from "unknown" and "any" to fixed values that don't make sense in context (e.g.,
		// Exif.Image.DateTime has a count of 20, are there really supposed to be 20 dates?  Probably not, it's more
		// likely that it's referring to a string with a maximum of 20 characters).  Therefore, we'll disregard count
		// values for these types.

		if typeId == typeIdAsciiString || typeId == typeIdComment {
			return 1
		}

		return info.Tag.Count

	case "IPTC":
		// IPTC datasets don't have a count value, they're marked as "repeatable" or not.  We can just adjust the count
		// to 1 or 2 depending on the value of that boolean -- the important thing is whether or not downstream callers
		// realize they're dealing with a single value or a slice.

		if info.Tag.Repeatable {
			return 2
		}

		// An exception to the above rule is that when the "undefined" type (i.e., raw bytes) is used the intention is
		// for us to generate a byte array as the type, so we'll adjust the count accordingly.

		if info.Tag.TypeID == typeIdUndefined {
			return 2
		}

		return 1
	}

	// Anything else is assumed to be XMP metadata, which we consider to be single-valued.

	return 1
}

func getDuplicateTagNames(f family, groupNames []string) map[string]bool {
	var foundTagNames = map[string]bool{}
	var result = map[string]bool{}

	for _, groupName := range groupNames {
		for tagName := range f[groupName] {
			if foundTagNames[tagName] {
				result[tagName] = true
			} else {
				foundTagNames[tagName] = true
			}
		}
	}

	return result
}

func getFixedTagName(tagName string) string {
	// Standard set of characters that some tags include that we can't use for a function name.

	tagName = strings.ReplaceAll(tagName, " ", "")
	tagName = strings.ReplaceAll(tagName, "-", "")

	// Some XMP tags start with lowercase letters.

	return strings.Title(tagName)
}

func getFunctionMappings(familyName string, f family, groupNames []string) ([]string, map[string]functionInfo) {
	var duplicateTagNames = getDuplicateTagNames(f, groupNames)
	var functionMappings = make(map[string]functionInfo)
	var sortedFunctionNames []string

	for _, groupName := range groupNames {
		for tagName := range f[groupName] {
			var functionName string

			if duplicateTagNames[tagName] {
				functionName = groupName + getFixedTagName(tagName)
			} else {
				functionName = getFixedTagName(tagName)
			}

			sortedFunctionNames = append(sortedFunctionNames, functionName)

			functionMappings[functionName] = functionInfo{
				FullTagName: familyName + "." + groupName + "." + tagName,
				Tag:         f[groupName][tagName],
			}
		}
	}

	sort.Strings(sortedFunctionNames)

	return sortedFunctionNames, functionMappings
}

func getTypeIDMapping(typeId int) typeIdMapping {
	if mapping, ok := typeIdMappings[typeId]; ok {
		return mapping
	}

	panic(fmt.Sprintf("invalid type id %d encountered", typeId))
}

func initTemplate(name, content string) (*template.Template, error) {
	var err error
	var tmpl = template.New(name).Funcs(funcMap)

	tmpl, err = tmpl.Parse(content)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func templateFuncFixDescription(description string) string {
	description = strings.TrimSpace(description)

	description = strings.ToLower(string(description[0])) + description[1:]

	// Fix double quotes to single quotes since this description will appear within double quotes.

	description = strings.ReplaceAll(description, "\"", "'")

	// Remove any trailing "." since we'll be inserting one.

	if strings.HasSuffix(description, ".") {
		description = description[:len(description)-1]
	}

	return description
}

func templateFuncIsHelperEnabled(info functionInfo, disabledHelpers map[string]bool) bool {
	return !disabledHelpers[info.FullTagName]
}

func templateFuncIsTestEnabled(info functionInfo, disabledTests map[string]bool) bool {
	return !disabledTests[info.FullTagName]
}

func templateFuncLastPackage(packageName string) string {
	var index = strings.Index(packageName, "/")

	if index == -1 {
		return packageName
	}

	return packageName[index+1:]
}

func templateFuncRequiredImports(functionMappings map[string]functionInfo, testing bool) []string {
	var i = 0
	var importMap = map[string]bool{
		"golang.handcraftedbits.com/ezif":                 true,
		"golang.handcraftedbits.com/ezif/helper":          true,
		"golang.handcraftedbits.com/ezif/helper/internal": true,
	}
	var result []string

	if testing {
		importMap["testing"] = true
		importMap["golang.handcraftedbits.com/ezif/helper/internal"] = true
	}

	for _, info := range functionMappings {
		var mapping = getTypeIDMapping(info.Tag.TypeID)
		var requiredTestImports = mapping.requiredTestImports

		if testing && requiredTestImports != nil {
			for _, requiredImport := range requiredTestImports {
				importMap[requiredImport] = true
			}
		}
	}

	result = make([]string, len(importMap))

	for key := range importMap {
		result[i] = key

		i++
	}

	sort.Strings(result)

	return result
}

func templateFuncReturnType(familyName string, info functionInfo) string {
	var count = getAdjustedCount(familyName, info)
	var result = getTypeIDMapping(info.Tag.TypeID).returnType

	if count != 1 {
		result += "Slice"
	}

	return result
}

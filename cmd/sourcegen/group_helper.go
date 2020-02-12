package main // import "golang.handcraftedbits.com/ezif/cmd/sourcegen"

import (
	"bytes"
	"fmt"
	"go/format"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"golang.handcraftedbits.com/ezif"
)

//
// Private types
//

type families map[string]family
type family map[string]group

type functionInfo struct {
	Family      ezif.Family
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
	FunctionMappings   map[string]functionInfo
	FunctionNames      []string
	PackageDescription string
	PackageName        string
	Reference          string
}

type tag struct {
	Count       int     `json:"count"`
	Description string  `json:"description"`
	Label       string  `json:"label"`
	MaxBytes    int     `json:"maxBytes"`
	MinBytes    int     `json:"minBytes"`
	Repeatable  bool    `json:"repeatable"`
	TypeID      ezif.ID `json:"typeID"`
}

type typeIDMapping struct {
	returnType string
}

//
// Private variables
//

var funcMap = template.FuncMap{
	"FixDescription":  templateFuncFixDescription,
	"IsHelperEnabled": templateFuncIsHelperEnabled,
	"IsSlice":         templateFuncIsSlice,
	"IsTestEnabled":   templateFuncIsTestEnabled,
	"LastPackage":     templateFuncLastPackage,
	"PropertyName":    templateFuncPropertyName,
	"ReturnType":      templateFuncReturnType,
}

var typeIDMappings = map[ezif.ID]typeIDMapping{
	ezif.IDAsciiString:      {"String"},
	ezif.IDComment:          {"String"},
	ezif.IDIPTCDate:         {"Date"},
	ezif.IDIPTCString:       {"String"},
	ezif.IDIPTCTime:         {"Time"},
	ezif.IDSignedByte:       {"SignedByte"},
	ezif.IDSignedLong:       {"SignedLong"},
	ezif.IDSignedRational:   {"SignedRational"},
	ezif.IDSignedShort:      {"SignedShort"},
	ezif.IDTIFFDouble:       {"Double"},
	ezif.IDTIFFFloat:        {"Float"},
	ezif.IDUndefined:        {"Undefined"},
	ezif.IDUnsignedByte:     {"UnsignedByte"},
	ezif.IDUnsignedLong:     {"UnsignedLong"},
	ezif.IDUnsignedRational: {"UnsignedRational"},
	ezif.IDUnsignedShort:    {"UnsignedShort"},
	ezif.IDXMPAlt:           {"String"},
	ezif.IDXMPBag:           {"String"},
	ezif.IDXMPLangAlt:       {"StringMap"},
	ezif.IDXMPSeq:           {"String"},
	ezif.IDXMPText:          {"String"},
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

func getAdjustedCount(info functionInfo) int {
	switch info.Family {
	case ezif.FamilyExif:
		var typeID = info.Tag.TypeID

		// Exiv2 seems to treat the ASCII string and comment ezif as characters instead of full strings because the
		// count values range from "unknown" and "any" to fixed values that don't make sense in context (e.g.,
		// Exif.Image.DateTime has a count of 20, are there really supposed to be 20 dates?  Probably not, it's more
		// likely that it's referring to a string with a maximum of 20 characters).  Therefore, we'll disregard count
		// values for these ezif.

		if typeID == ezif.IDAsciiString || typeID == ezif.IDComment {
			return 1
		}

		return info.Tag.Count

	case ezif.FamilyIPTC:
		// IPTC datasets don't have a count value, they're marked as "repeatable" or not.  We can just adjust the count
		// to 1 or 2 depending on the value of that boolean -- the important thing is whether or not downstream callers
		// realize they're dealing with a single value or a slice.

		if info.Tag.Repeatable {
			return 2
		}

		// An exception to the above rule is that when the "undefined" type (i.e., raw bytes) is used the intention is
		// for us to generate a byte array as the type, so we'll adjust the count accordingly.

		if info.Tag.TypeID == ezif.IDUndefined {
			return 2
		}

		return 1
	}

	// Anything else is assumed to be XMP metadata, which is multi-valued (we'll just use a count of "2" in that case)
	// unless the type is ezif.IDXMPText or ezif.IDXMPLangAlt (technically multi-valued, but exposed through a map).

	if info.Tag.TypeID == ezif.IDXMPText || info.Tag.TypeID == ezif.IDXMPLangAlt {
		return 1
	}

	return 2
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
				Family:      ezif.Family(familyName),
				FullTagName: familyName + "." + groupName + "." + tagName,
				Tag:         f[groupName][tagName],
			}
		}
	}

	sort.Strings(sortedFunctionNames)

	return sortedFunctionNames, functionMappings
}

func getTypeIDMapping(typeID ezif.ID) typeIDMapping {
	if mapping, ok := typeIDMappings[typeID]; ok {
		return mapping
	}

	panic(fmt.Sprintf("invalid type id %d encountered", typeID))
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

func templateFuncLastPackage(packageName string) string {
	var index = strings.Index(packageName, "/")

	if index == -1 {
		return packageName
	}

	return packageName[index+1:]
}

func templateFuncPropertyName(info functionInfo) string {
	if info.Family == ezif.FamilyExif {
		return string(ezif.FamilyExif)
	}

	return strings.ToUpper(string(info.Family))
}

func templateFuncReturnType(info functionInfo) string {
	var count = getAdjustedCount(info)
	var result = getTypeIDMapping(info.Tag.TypeID).returnType

	if count != 1 {
		result += "Slice"
	}

	return result
}

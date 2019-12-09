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

type accessorConfig struct {
	TreatSliceAsSingleValue bool   `yaml:"treatSliceAsSingleValue"`
	Type                    string `yaml:"type"`
	ValueMethod             string `yaml:"valueMethod"`
	ValueMethodAlias        string `yaml:"valueMethodAlias"`
}

type accessorImplementation struct {
	Family           string
	InterfaceName    string
	IsSlice          bool
	Name             string
	Type             string
	ValueMethod      string
	ValueMethodAlias string
}

type accessorInterface struct {
	Name string
	Type string
}

type accessorsConfig map[string]map[string][]accessorConfig

type accessorTemplateContext struct {
	AccessorImplementations []accessorImplementation
	AccessorInterfaces      []accessorInterface
	PackageName             string
}

//
// Private functions
//

func generateAccessorSource(packageName string, accessors accessorsConfig, templateBody string) (string, error) {
	var buffer bytes.Buffer
	var err error
	var formattedSource []byte
	var templateRoot *template.Template

	templateRoot, err = initTemplate("root", templateBody)

	if err != nil {
		return "", err
	}

	err = templateRoot.Execute(&buffer, &accessorTemplateContext{
		AccessorImplementations: getAccessorImplementations(accessors),
		AccessorInterfaces:      getAccessorInterfaces(accessors),
		PackageName:             packageName,
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

func getAccessorImplementations(accessors accessorsConfig) []accessorImplementation {
	var result []accessorImplementation

	for name := range accessors {
		for family := range accessors[name] {
			for _, value := range accessors[name][family] {
				var isSlice = strings.HasPrefix(value.Type, "[]") && !value.TreatSliceAsSingleValue
				var structName = strings.ToLower(family) + value.ValueMethod

				if isSlice {
					structName += "Slice"
				}

				result = append(result, accessorImplementation{
					Family:           family,
					InterfaceName:    name,
					IsSlice:          isSlice,
					Name:             structName,
					Type:             value.Type,
					ValueMethod:      value.ValueMethod,
					ValueMethodAlias: value.ValueMethodAlias,
				})
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		var firstName = result[i].Family + result[i].ValueMethod
		var lastName = result[j].Family + result[j].ValueMethod

		return strings.Compare(firstName, lastName) == -1
	})

	return result
}

func getAccessorInterfaces(accessors accessorsConfig) []accessorInterface {
	var result []accessorInterface
	var sortedNames []string

	for name := range accessors {
		sortedNames = append(sortedNames, name)
	}

	sort.Strings(sortedNames)

	for _, name := range sortedNames {
		var valueType string

		// The type for the interface getters and setters will be the same for all implementations, so just use the
		// first one.

		for key := range accessors[name] {
			valueType = accessors[name][key][0].Type

			break
		}

		result = append(result, accessorInterface{
			Name: name,
			Type: valueType,
		})
	}

	return result
}

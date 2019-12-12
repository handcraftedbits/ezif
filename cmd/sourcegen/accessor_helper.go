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
	Name     string
	Type     string
}

type accessorTemplateContext struct {
	AccessorInfos []accessorInfo
	PackageName   string
}

//
// Private functions
//

func generateAccessorSource(packageName string, accessors map[string]string, templateBody string) (string, error) {
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
			Name:     name,
			Type:     goType,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i].Name, result[j].Name) < 0
	})

	return result
}

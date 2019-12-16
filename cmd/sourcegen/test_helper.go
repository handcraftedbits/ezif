package main // import "golang.handcraftedbits.com/ezif/cmd/sourcegen"

//
// Private functions
//

func generateGroupTestSource(familyName string, f family, packageName string, gc groupConfig) (string, error) {
	return generateSource(familyName, f, packageName, gc, templateGroupTestSource)
}

func templateFuncIsSlice(familyName string, info functionInfo) bool {
	return getAdjustedCount(familyName, info) != 1
}

func templateFuncIsTestEnabled(info functionInfo, disabledTests map[string]bool) bool {
	return !disabledTests[info.FullTagName]
}

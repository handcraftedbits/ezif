package main // import "golang.handcraftedbits.com/ezif/cmd/sourcegen"

//
// Private functions
//

func generateGroupTestSource(familyName string, f family, packageName string, gc groupConfig) (string, error) {
	return generateSource(familyName, f, packageName, gc, templateGroupTestSource)
}

func templateFuncIsSlice(info functionInfo) bool {
	return getAdjustedCount(info) != 1
}

func templateFuncIsTestEnabled(info functionInfo, disabledTests map[string]bool) bool {
	return !disabledTests[info.FullTagName]
}

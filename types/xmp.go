package types // import "golang.handcraftedbits.com/ezif/types"
import (
	"fmt"
)

//
// Public types
//

type XMPLangAlt interface {
	fmt.Stringer

	Language() string
	Value() string
}

//
// Public functions
//

func NewXMPLangAlt(language, value string) XMPLangAlt {
	return &xmpLangAltImpl{
		language: language,
		value:    value,
	}
}

//
// Private types
//

// XMPLangAlt implementation
type xmpLangAltImpl struct {
	language string
	value    string
}

func (xmpLangAlt *xmpLangAltImpl) Language() string {
	return xmpLangAlt.language
}

func (xmpLangAlt *xmpLangAltImpl) String() string {
	return fmt.Sprintf("lang=\"%s\" %s", xmpLangAlt.language, xmpLangAlt.value)
}

func (xmpLangAlt *xmpLangAltImpl) Value() string {
	return xmpLangAlt.value
}

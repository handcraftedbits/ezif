package ezif // import "golang.handcraftedbits.com/ezif"

//
// Public types
//

type Family string

type Metadata interface {
	Exif() Properties
	IPTC() Properties
	XMP() Properties
}

type Properties interface {
	Get(key string) Property
	HasKey(key string) bool
	Keys() []string
}

type Property interface {
	Family() Family
	GroupName() string
	InterpretedValue() string
	Label() string
	TagName() string
	TypeID() ID
	Value() interface{}
}

//
// Public constants
//

const (
	FamilyExif Family = "Exif"
	FamilyIPTC Family = "Iptc"
	FamilyXMP  Family = "Xmp"
)

package ezif // import "golang.handcraftedbits.com/ezif"

import (
	"math/big"
	"time"
)

//
// Public types
//

type ImageMetadata interface {
	Exif() ExifMetadata
	IPTC() IPTCMetadata
	XMP() XMPMetadata
}

type Metadata interface {
	HasKey(key string) bool
	Keys() []string
}

type Datum interface {
	FamilyName() string
	GroupName() string
	InterpretedValue() string
	Label() string
	TagName() string
	TypeID() TypeID
}

type ExifMetadata interface {
	Metadata

	Get(key string) ExifDatum
}

type ExifDatum interface {
	Datum

	Values() []ExifValue
}

type ExifValue interface {
	ASCIIString() string
	Comment() string
	SignedByte() int8
	SignedLong() int32
	SignedRational() *big.Rat
	SignedShort() int16
	TIFFDouble() float64
	TIFFFloat() float32
	Undefined() byte
	UnsignedByte() uint8
	UnsignedLong() uint32
	UnsignedRational() *big.Rat
	UnsignedShort() uint16
}

type IPTCMetadata interface {
	Metadata

	Get(key string) IPTCDatum
}

type IPTCDatum interface {
	Datum

	Values() []IPTCValue
}

type IPTCValue interface {
	Date() time.Time
	Short() uint16
	String() string
	Time() time.Time
	Undefined() []byte
}

type TypeID int

type XMPLangAlt interface {
	Language() string
	Value() string
}
type XMPMetadata interface {
	Metadata

	Get(key string) XMPDatum
}

type XMPDatum interface {
	Datum

	Value() XMPValue
}

type XMPValue interface {
	Alt() []string
	Bag() []string
	LangAlt() []XMPLangAlt
	Seq() []string
	Text() string
}

//
// Public constants
//

const (
	TypeIDUnsignedByte     TypeID = 1
	TypeIDAsciiString      TypeID = 2
	TypeIDUnsignedShort    TypeID = 3
	TypeIDUnsignedLong     TypeID = 4
	TypeIDUnsignedRational TypeID = 5
	TypeIDSignedByte       TypeID = 6
	TypeIDUndefined        TypeID = 7
	TypeIDSignedShort      TypeID = 8
	TypeIDSignedLong       TypeID = 9
	TypeIDSignedRational   TypeID = 10
	TypeIDTIFFFloat        TypeID = 11
	TypeIDTIFFDouble       TypeID = 12
	TypeIDIPTCString       TypeID = 0x10000
	TypeIDIPTCDate         TypeID = 0x10001
	TypeIDIPTCTime         TypeID = 0x10002
	TypeIDComment          TypeID = 0x10003
	TypeIDXMPText          TypeID = 0x10005
	TypeIDXMPAlt           TypeID = 0x10006
	TypeIDXMPBag           TypeID = 0x10007
	TypeIDXMPSeq           TypeID = 0x10008
	TypeIDXMPLangAlt       TypeID = 0x10009
	TypeIDInvalid          TypeID = 0x1FFFE
)

//
// Private types
//

// ImageMetadata implementation
type imageMetadataImpl struct {
	exifMetadata ExifMetadata
}

func (imageMetadata *imageMetadataImpl) Exif() ExifMetadata {
	return imageMetadata.exifMetadata
}

func (imageMetadata *imageMetadataImpl) IPTC() IPTCMetadata {
	return nil
}

func (imageMetadata *imageMetadataImpl) XMP() XMPMetadata {
	return nil
}

// Datum implementation
type datumImpl struct {
	familyName       string
	groupName        string
	interpretedValue string
	label            string
	tagName          string
	typeId           TypeID
}

func (datum *datumImpl) FamilyName() string {
	return datum.familyName
}

func (datum *datumImpl) GroupName() string {
	return datum.groupName
}

func (datum *datumImpl) InterpretedValue() string {
	return datum.interpretedValue
}

func (datum *datumImpl) Label() string {
	return datum.label
}

func (datum *datumImpl) TagName() string {
	return datum.tagName
}

func (datum *datumImpl) TypeID() TypeID {
	return datum.typeId
}

func (datum *datumImpl) key() string {
	return datum.familyName + "." + datum.groupName + "." + datum.tagName
}

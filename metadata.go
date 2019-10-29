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
	TypeId() TypeId
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

	Get(key string) []IPTCDatum
}

type IPTCDatum interface {
	Datum

	Value() IPTCValue
}

type IPTCValue interface {
	Date() time.Time
	Short() uint16
	String() string
	Time() time.Time
	Undefined() []byte
}

type TypeId int

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
	TypeIdUnsignedByte     TypeId = 1
	TypeIdAsciiString      TypeId = 2
	TypeIdUnsignedShort    TypeId = 3
	TypeIdUnsignedLong     TypeId = 4
	TypeIdUnsignedRational TypeId = 5
	TypeIdSignedByte       TypeId = 6
	TypeIdUndefined        TypeId = 7
	TypeIdSignedShort      TypeId = 8
	TypeIdSignedLong       TypeId = 9
	TypeIdSignedRational   TypeId = 10
	TypeIdTIFFFloat        TypeId = 11
	TypeIdTIFFDouble       TypeId = 12
	TypeIdIPTCString       TypeId = 0x10000
	TypeIdIPTCDate         TypeId = 0x10001
	TypeIdIPTCTime         TypeId = 0x10002
	TypeIdComment          TypeId = 0x10003
	TypeIdXMPText          TypeId = 0x10005
	TypeIdXMPAlt           TypeId = 0x10006
	TypeIdXMPBag           TypeId = 0x10007
	TypeIdXMPSeq           TypeId = 0x10008
	TypeIdXMPLangAlt       TypeId = 0x10009
	TypeIdInvalid          TypeId = 0x1FFFE
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
	typeId           TypeId
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

func (datum *datumImpl) TypeId() TypeId {
	return datum.typeId
}

func (datum *datumImpl) key() string {
	return datum.familyName + "." + datum.groupName + "." + datum.tagName
}

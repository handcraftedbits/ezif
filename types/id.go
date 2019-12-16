package types // import "golang.handcraftedbits.com/ezif/types"

//
// Public types
//

type ID int

func (id ID) String() string {
	switch id {
	case IDAsciiString:
		return "IDAsciiString"

	case IDComment:
		return "IDComment"

	case IDIPTCDate:
		return "IDIPTCDate"

	case IDIPTCString:
		return "IDIPTCString"

	case IDIPTCTime:
		return "IDIPTCTime"

	case IDSignedByte:
		return "IDSignedByte"

	case IDSignedLong:
		return "IDSignedLong"

	case IDSignedRational:
		return "IDSignedRational"

	case IDSignedShort:
		return "IDSignedShort"

	case IDTIFFDouble:
		return "IDTIFFDouble"

	case IDTIFFFloat:
		return "IDTIFFFloat"

	case IDUndefined:
		return "IDUndefined"

	case IDUnsignedByte:
		return "IDUnsignedByte"

	case IDUnsignedLong:
		return "IDUnsignedLong"

	case IDUnsignedRational:
		return "IDUnsignedRational"

	case IDUnsignedShort:
		return "IDUnsignedShort"

	case IDXMPAlt:
		return "IDXMPAlt"

	case IDXMPBag:
		return "IDXMPBag"

	case IDXMPLangAlt:
		return "IDXMPLangAlt"

	case IDXMPSeq:
		return "IDXMPSeq"

	case IDXMPText:
		return "IDXMPText"
	}

	return "IDInvalid"
}

//
// Public constants
//

const (
	IDUnsignedByte     ID = 1
	IDAsciiString      ID = 2
	IDUnsignedShort    ID = 3
	IDUnsignedLong     ID = 4
	IDUnsignedRational ID = 5
	IDSignedByte       ID = 6
	IDUndefined        ID = 7
	IDSignedShort      ID = 8
	IDSignedLong       ID = 9
	IDSignedRational   ID = 10
	IDTIFFFloat        ID = 11
	IDTIFFDouble       ID = 12
	IDIPTCString       ID = 0x10000
	IDIPTCDate         ID = 0x10001
	IDIPTCTime         ID = 0x10002
	IDComment          ID = 0x10003
	IDXMPText          ID = 0x10005
	IDXMPAlt           ID = 0x10006
	IDXMPBag           ID = 0x10007
	IDXMPSeq           ID = 0x10008
	IDXMPLangAlt       ID = 0x10009
	IDInvalid          ID = 0x1FFFE
)

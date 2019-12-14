package types // import "golang.handcraftedbits.com/ezif/types"

//
// Public types
//

type ID int

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

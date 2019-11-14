package internal // import "golang.handcraftedbits.com/ezif/helper/internal"

import (
	"golang.handcraftedbits.com/ezif"
)

//
// Public functions
//

func GetXMPValueAsAlt(metadatum ezif.XMPDatum) []string {
	return metadatum.Value().Alt()
}

func GetXMPValueAsBag(metadatum ezif.XMPDatum) []string {
	return metadatum.Value().Bag()
}

func GetXMPValueAsLangAlt(metadatum ezif.XMPDatum) []ezif.XMPLangAlt {
	return metadatum.Value().LangAlt()
}

func GetXMPValueAsSeq(metadatum ezif.XMPDatum) []string {
	return metadatum.Value().Seq()
}

func GetXMPValueAsText(metadatum ezif.XMPDatum) string {
	return metadatum.Value().Text()
}

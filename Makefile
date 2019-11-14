# Variables

CMD_DOCKER_RUN=docker run -it --rm -v $(DIR_BASE):/ezif $(DOCKER_IMAGE)
CMD_GENHELPER_RUN=go run $(DIR_CMD_GENHELPER) -m $(FILE_EXIV2_METADATA)
CMD_EXIV2METADATA_RUN=$(CMD_DOCKER_RUN) go run ./cmd/exiv2metadata

DIR_BASE=$(dir $(realpath $(lastword $(MAKEFILE_LIST))))
DIR_CMD=$(DIR_BASE)cmd
DIR_CMD_EXIV2METADATA=$(DIR_CMD)/exiv2metadata
DIR_CMD_GENHELPER=$(DIR_CMD)/genhelper
DIR_HELPER=$(DIR_BASE)helper
DIR_HELPER_EXIF=$(DIR_HELPER)/exif
DIR_HELPER_IPTC=$(DIR_HELPER)/iptc
DIR_HELPER_XMP=$(DIR_HELPER)/xmp

DOCKER_IMAGE=handcraftedbits/ezif-build:$(VERSION)

FILE_DOCKERFILE=$(DIR_BASE)Dockerfile
FILE_DOCKER_BUILT=$(DIR_BASE).docker-built
FILE_EXIV2_METADATA=$(DIR_BASE).exiv2metadata.json
FILE_GENHELPER_CONFIG=$(DIR_BASE).genhelper.yaml

# A couple libpthread symbols seem to be marked as weak, causing a segfault when run in a non-musl environment.
LDFLAGS=-ldflags "-linkmode external -extldflags '-Wl,-u,pthread_mutexattr_init -Wl,-u,pthread_mutexattr_destroy -Wl,-u,pthread_mutexattr_settype -static'"

VERSION=0.9.0

# Phony/special targets

.DELETE_ON_ERROR: $(DIR_HELPER)/%.go $(DIR_HELPER)/%_test.go $(FILE_EXIV2_METADATA)
.PHONY: all helpers helpers_test

all: helpers helpers_test

helpers: $(DIR_HELPER_EXIF)/exif.go \
	$(DIR_HELPER_IPTC)/iptc.go \
	$(DIR_HELPER_XMP)/xmp.go \
	$(DIR_HELPER_XMP)/acdsee/acdsee.go \
	$(DIR_HELPER_XMP)/aux/aux.go \
	$(DIR_HELPER_XMP)/bj/bj.go \
	$(DIR_HELPER_XMP)/crs/crs.go \
	$(DIR_HELPER_XMP)/crss/crss.go \
	$(DIR_HELPER_XMP)/dc/dc.go \
	$(DIR_HELPER_XMP)/dcterms/dcterms.go \
	$(DIR_HELPER_XMP)/digikam/digikam.go \
	$(DIR_HELPER_XMP)/dm/dm.go \
	$(DIR_HELPER_XMP)/dwc/dwc.go \
	$(DIR_HELPER_XMP)/exif/exif.go \
	$(DIR_HELPER_XMP)/exifex/exifex.go \
	$(DIR_HELPER_XMP)/expressionmedia/expressionmedia.go \
	$(DIR_HELPER_XMP)/gpano/gpano.go \
	$(DIR_HELPER_XMP)/iptccore/iptccore.go \
	$(DIR_HELPER_XMP)/iptcext/iptcext.go \
	$(DIR_HELPER_XMP)/kipi/kipi.go \
	$(DIR_HELPER_XMP)/lr/lr.go \
	$(DIR_HELPER_XMP)/mediapro/mediapro.go \
	$(DIR_HELPER_XMP)/microsoftphoto/microsoftphoto.go \
	$(DIR_HELPER_XMP)/mm/mm.go \
	$(DIR_HELPER_XMP)/mp/mp.go \
	$(DIR_HELPER_XMP)/mwgkw/mwgkw.go \
	$(DIR_HELPER_XMP)/mwgrs/mwgrs.go \
	$(DIR_HELPER_XMP)/pdf/pdf_test.go \
	$(DIR_HELPER_XMP)/photoshop/photoshop_test.go \
	$(DIR_HELPER_XMP)/plus/plus.go \
	$(DIR_HELPER_XMP)/rights/rights.go \
	$(DIR_HELPER_XMP)/tiff/tiff.go \
	$(DIR_HELPER_XMP)/tpg/tpg.go
helpers_test: $(DIR_HELPER_EXIF)/exif_test.go \
	$(DIR_HELPER_IPTC)/iptc_test.go \
	$(DIR_HELPER_XMP)/xmp_test.go \
	$(DIR_HELPER_XMP)/acdsee/acdsee_test.go \
	$(DIR_HELPER_XMP)/aux/aux_test.go \
	$(DIR_HELPER_XMP)/bj/bj_test.go \
	$(DIR_HELPER_XMP)/crs/crs_test.go \
	$(DIR_HELPER_XMP)/crss/crss_test.go \
	$(DIR_HELPER_XMP)/dc/dc_test.go \
	$(DIR_HELPER_XMP)/dcterms/dcterms_test.go \
	$(DIR_HELPER_XMP)/digikam/digikam_test.go \
	$(DIR_HELPER_XMP)/dm/dm_test.go \
	$(DIR_HELPER_XMP)/dwc/dwc_test.go \
	$(DIR_HELPER_XMP)/exif/exif_test.go \
	$(DIR_HELPER_XMP)/exifex/exifex_test.go \
	$(DIR_HELPER_XMP)/expressionmedia/expressionmedia_test.go \
	$(DIR_HELPER_XMP)/gpano/gpano_test.go \
	$(DIR_HELPER_XMP)/iptccore/iptccore_test.go \
	$(DIR_HELPER_XMP)/iptcext/iptcext_test.go \
	$(DIR_HELPER_XMP)/kipi/kipi_test.go \
	$(DIR_HELPER_XMP)/lr/lr_test.go \
	$(DIR_HELPER_XMP)/mediapro/mediapro_test.go \
	$(DIR_HELPER_XMP)/microsoftphoto/microsoftphoto_test.go \
	$(DIR_HELPER_XMP)/mm/mm_test.go \
	$(DIR_HELPER_XMP)/mp/mp_test.go \
	$(DIR_HELPER_XMP)/mwgkw/mwgkw_test.go \
	$(DIR_HELPER_XMP)/mwgrs/mwgrs_test.go \
	$(DIR_HELPER_XMP)/pdf/pdf.go \
	$(DIR_HELPER_XMP)/photoshop/photoshop.go \
	$(DIR_HELPER_XMP)/plus/plus_test.go \
	$(DIR_HELPER_XMP)/rights/rights_test.go \
	$(DIR_HELPER_XMP)/tiff/tiff_test.go \
	$(DIR_HELPER_XMP)/tpg/tpg_test.go

# File targets

$(FILE_DOCKER_BUILT): $(FILE_DOCKERFILE)
	docker build -t $(DOCKER_IMAGE) -f $(FILE_DOCKERFILE) $(DIR_BASE)
	docker images -f "reference=$(DOCKER_IMAGE)" --format="{{ .ID }}" > $@

$(DIR_HELPER)/%.go: $(FILE_EXIV2_METADATA) $(FILE_GENHELPER_CONFIG) $(wildcard $(DIR_CMD_GENHELPER)/*)
	mkdir -p $(dir $@)
	$(CMD_GENHELPER_RUN) -c $(FILE_GENHELPER_CONFIG) -p $(patsubst %/,%,$(dir $*)) > $@
$(DIR_HELPER)/%_test.go: $(FILE_EXIV2_METADATA) $(FILE_GENHELPER_CONFIG) $(wildcard $(DIR_CMD_GENHELPER)/*)
	mkdir -p $(dir $@)
	$(CMD_GENHELPER_RUN) -c $(FILE_GENHELPER_CONFIG) -p $(patsubst %/,%,$(dir $*)) -t > $@

$(FILE_EXIV2_METADATA): $(wildcard $(DIR_CMD_EXIV2METADATA)/*) $(FILE_DOCKER_BUILT)
	$(CMD_EXIV2METADATA_RUN) > $@

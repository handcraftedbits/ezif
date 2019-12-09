# Variables

CMD_DOCKER_RUN=docker run -it --rm -v $(DIR_BASE):/ezif $(DOCKER_IMAGE)
CMD_EXIV2METADATA_RUN=$(CMD_DOCKER_RUN) go run ./cmd/exiv2metadata
CMD_SOURCEGEN_HELPER_RUN=$(CMD_SOURCEGEN_RUN) helper -m $(FILE_EXIV2_METADATA)
CMD_SOURCEGEN_RUN=go run $(DIR_CMD_SOURCEGEN) -c $(FILE_SOURCEGEN_CONFIG)

DIR_BASE=$(dir $(realpath $(lastword $(MAKEFILE_LIST))))
DIR_CMD=$(DIR_BASE)cmd
DIR_CMD_EXIV2METADATA=$(DIR_CMD)/exiv2metadata
DIR_CMD_SOURCEGEN=$(DIR_CMD)/sourcegen
DIR_HELPER=$(DIR_BASE)helper
DIR_HELPER_EXIF=$(DIR_HELPER)/exif
DIR_HELPER_IPTC=$(DIR_HELPER)/iptc
DIR_HELPER_XMP=$(DIR_HELPER)/xmp

DOCKER_IMAGE=handcraftedbits/ezif-build:$(VERSION)

FILE_ACCESSOR=$(DIR_HELPER)/accessor.go
FILE_DOCKERFILE=$(DIR_BASE)Dockerfile
FILE_DOCKER_BUILT=$(DIR_BASE).docker-built
FILE_EXIV2_METADATA=$(DIR_BASE).exiv2metadata.json
FILE_SOURCEGEN_CONFIG=$(DIR_BASE).sourcegen.yaml

# A couple libpthread symbols seem to be marked as weak, causing a segfault when run in a non-musl environment.
LDFLAGS=-ldflags "-linkmode external -extldflags '-Wl,-u,pthread_mutexattr_init -Wl,-u,pthread_mutexattr_destroy -Wl,-u,pthread_mutexattr_settype -static'"

VERSION=0.9.0

# Phony/special targets

.DELETE_ON_ERROR: $(DIR_HELPER)/%.go $(DIR_HELPER)/%_test.go $(FILE_EXIV2_METADATA)
.PHONY: all clean helpers helpers_test test

all: helpers_test

build: helpers
	$(CMD_DOCKER_RUN) go build ./...

clean:
	rm -rf $(DIR_HELPER_EXIF) $(DIR_HELPER_IPTC) $(DIR_HELPER_XMP) $(FILE_ACCESSOR) $(FILE_DOCKER_BUILT) \
		$(FILE_EXIV2_METADATA)

helpers: $(FILE_ACCESSOR) \
	$(DIR_HELPER_EXIF)/exif.go \
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
	$(DIR_HELPER_XMP)/exifcore/exifcore.go \
	$(DIR_HELPER_XMP)/exifext/exifext.go \
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
helpers_test: helpers $(DIR_HELPER_EXIF)/exif_test.go \
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
	$(DIR_HELPER_XMP)/exifcore/exifcore_test.go \
	$(DIR_HELPER_XMP)/exifext/exifext_test.go \
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

test: helpers_test
	$(CMD_DOCKER_RUN) go test ./helper/exif

# File targets

$(FILE_DOCKER_BUILT): $(FILE_DOCKERFILE)
	docker build -t $(DOCKER_IMAGE) -f $(FILE_DOCKERFILE) $(DIR_BASE)
	docker images -f "reference=$(DOCKER_IMAGE)" --format="{{ .ID }}" > $@

$(DIR_HELPER)/%.go: $(FILE_EXIV2_METADATA) $(FILE_SOURCEGEN_CONFIG) $(wildcard $(DIR_CMD_SOURCEGEN)/*)
	mkdir -p $(dir $@)
	$(CMD_SOURCEGEN_HELPER_RUN) -g $(patsubst %/,%,$(dir $*)) > $@
$(DIR_HELPER)/%_test.go: $(FILE_EXIV2_METADATA) $(FILE_SOURCEGEN_CONFIG) $(wildcard $(DIR_CMD_SOURCEGEN)/*)
	mkdir -p $(dir $@)
	$(CMD_SOURCEGEN_HELPER_RUN) -g $(patsubst %/,%,$(dir $*)) -t > $@
$(FILE_ACCESSOR): $(FILE_SOURCEGEN_CONFIG) $(wildcard $(DIR_CMD_SOURCEGEN)/*)
	mkdir -p $(dir $@)
	$(CMD_SOURCEGEN_RUN) accessor -p helper > $@

$(FILE_EXIV2_METADATA): $(wildcard $(DIR_CMD_EXIV2METADATA)/*) $(FILE_DOCKER_BUILT)
	$(CMD_EXIV2METADATA_RUN) > $@

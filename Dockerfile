FROM alpine:3.10 AS cdeps
MAINTAINER HandcraftedBits <opensource@handcraftedbits.com>

ENV DIR_EXIV2=/tmp/exiv2
ENV DIR_EXPAT=/tmp/expat
ENV DIR_JANSSON=/tmp/jansson
ENV DIR_PCRE=/tmp/pcre
ENV DIR_SWIG=/tmp/swig
ENV DIR_ZLIB=/tmp/zlib

ENV VERSION_EXIV2=0.27.2
ENV VERSION_EXPAT=R_2_2_7/expat-2.2.7
ENV VERSION_JANSSON=2.12
ENV VERSION_PCRE=8.43
ENV VERSION_SWIG=4.0.1
ENV VERSION_ZLIB=1.2.11

RUN apk update && \
  apk add autoconf automake cmake curl file g++ gcc libtool linux-headers make musl-dev pkgconfig

RUN mkdir -p ${DIR_EXPAT} && \
  curl -L https://github.com/libexpat/libexpat/releases/download/${VERSION_EXPAT}.tar.gz | tar -xzvf - \
    -C ${DIR_EXPAT} --strip-components 1 && \
  cd ${DIR_EXPAT} && \
  ./configure --disable-shared && \
  make install
RUN mkdir -p ${DIR_PCRE} && \
  curl -L https://ftp.pcre.org/pub/pcre/pcre-${VERSION_PCRE}.tar.gz | tar -xzvf - -C ${DIR_PCRE} \
    --strip-components 1 && \
  cd ${DIR_PCRE} && \
  ./configure --enable-shared=no && \
  make install
RUN mkdir -p ${DIR_ZLIB} && \
  curl -L https://zlib.net/zlib-${VERSION_ZLIB}.tar.gz | tar -xzvf - -C ${DIR_ZLIB} --strip-components 1 && \
  cd ${DIR_ZLIB} && \
  ./configure --static --64 && \
  make install
RUN mkdir -p ${DIR_EXIV2} && \
  curl -L https://www.exiv2.org/builds/exiv2-${VERSION_EXIV2}-Source.tar.gz | tar -xzvf - -C ${DIR_EXIV2} \
    --strip-components 1 && \
  cd ${DIR_EXIV2} &&\
  cmake . -DBUILD_SHARED_LIBS=OFF -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_LIBDIR=lib \
    -DEXIV2_BUILD_EXIV2_COMMAND=ON -DEXIV2_BUILD_SAMPLES=OFF -DEXIV2_ENABLE_NLS=OFF -DEXIV2_ENABLE_VIDEO=ON && \
  cmake --build . && \
  make install
RUN mkdir -p ${DIR_JANSSON} && \
  curl -L http://www.digip.org/jansson/releases/jansson-${VERSION_JANSSON}.tar.gz | tar -xzvf - -C ${DIR_JANSSON} \
    --strip-components 1 && \
  cd ${DIR_JANSSON} && \
  ./configure --disable-shared && \
  make install

FROM golang:alpine
MAINTAINER HandcraftedBits <opensource@handcraftedbits.com>

RUN apk update && \
  apk add g++ gcc git make musl-dev

COPY --from=cdeps /usr/local/bin/exiv2 /usr/local/bin
COPY --from=cdeps /usr/local/include /usr/local/include
COPY --from=cdeps /usr/local/lib /usr/local/lib
COPY go.mod /
COPY go.sum /

RUN cd / && \
  go mod download && \
  rm /go.*

VOLUME /ezif

WORKDIR /ezif

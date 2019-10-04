FROM golang:1.12.10-alpine3.10

ENV GO111MODULE off
RUN go version \
    && echo $GOPATH \
    && apk update \
    && apk add --no-cache git wget unzip fontconfig alpine-sdk bash \
    && wget https://github.com/tomokuni/Myrica/raw/master/product/MyricaM.zip -q -O /tmp/MyricaM.zip \
    && (cd /tmp && unzip MyricaM.zip) \
    && git clone https://github.com/googlefonts/noto-emoji /usr/local/src/noto-emoji \
    && wget https://www.wfonts.com/download/data/2016/04/23/symbola/symbola.zip -q -O /tmp/symbola.zip \
    && (cd /tmp && unzip symbola.zip)
COPY . /go/src/github.com/jiro4989/textimg
WORKDIR /go/src/github.com/jiro4989/textimg
RUN go install

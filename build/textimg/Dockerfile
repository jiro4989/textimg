FROM golang:1.12.5-alpine3.9 AS build-stage

ENV GO111MODULE off
RUN go version \
    && echo $GOPATH \
    && apk update \
    && apk add --no-cache git wget unzip fontconfig \
    && go get github.com/jiro4989/textimg \
    && wget https://github.com/tomokuni/Myrica/raw/master/product/MyricaM.zip -q -O /tmp/MyricaM.zip \
    && (cd /tmp && unzip MyricaM.zip) \
    && git clone https://github.com/googlefonts/noto-emoji /usr/local/src/noto-emoji \
    && wget https://www.wfonts.com/download/data/2016/04/23/symbola/symbola.zip -q -O /tmp/symbola.zip \
    && (cd /tmp && unzip symbola.zip)

FROM alpine:latest AS exec-stage
COPY --from=build-stage /go/bin/textimg /usr/local/bin/
COPY --from=build-stage /tmp/MyricaM.TTC /usr/share/fonts/truetype/myrica/MyricaM.TTC
COPY --from=build-stage /usr/local/src/noto-emoji/png/128 /usr/share/emoji-image
COPY --from=build-stage /tmp/Symbola_hint.ttf /usr/share/fonts/truetype/symbola/
COPY --from=build-stage /tmp/Symbola_hint.ttf /usr/share/fonts/truetype/ancient-scripts/

ENV TEXTIMG_FONT_FILE /usr/share/fonts/truetype/myrica/MyricaM.TTC
ENV TEXTIMG_EMOJI_DIR /usr/share/emoji-image
ENV TEXTIMG_EMOJI_FONT_FILE /usr/share/fonts/truetype/symbola/Symbola_hint.ttf
RUN mkdir /images

ENTRYPOINT ["/usr/local/bin/textimg"]

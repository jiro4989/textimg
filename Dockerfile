FROM golang:1.12.5-alpine3.9 AS build-stage

ENV GO111MODULE on
RUN go version \
    && echo $GOPATH \
    && apk update \
    && apk add --no-cache git wget unzip fontconfig \
    && wget https://noto-website.storage.googleapis.com/pkgs/NotoSansCJKjp-hinted.zip \
    && mkdir /tmp/NotoSansCJKjp \
    && unzip NotoSansCJKjp-hinted.zip -d /tmp/NotoSansCJKjp \
    && go get github.com/jiro4989/textimg \
    && git clone https://github.com/googlefonts/noto-emoji /usr/local/src/noto-emoji \
    && wget https://www.wfonts.com/download/data/2016/04/23/symbola/symbola.zip \
    && unzip symbola.zip

FROM alpine:latest AS exec-stage
COPY --from=build-stage /root/go/bin/textimg /usr/local/bin/
COPY --from=build-stage /usr/local/src/noto-emoji /usr/local/src/
COPY --from=build-stage /Symbola_hint.ttf /usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf
COPY --from=build-stage /tmp/NotoSansCJKjp /usr/share/fonts/truetype/

ENV TEXTIMG_FONT_FILE /usr/share/fonts/truetype/
ENV TEXTIMG_EMOJI_FONT_FILE /usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf

ENTRYPOINT ["/usr/local/bin/textimg"]

FROM golang:1.18.3-alpine3.15 AS base

RUN go version \
    && echo $GOPATH \
    && apk update \
    && apk add --no-cache git wget unzip fontconfig alpine-sdk bash \
    && wget https://github.com/tomokuni/Myrica/raw/master/product/MyricaM.zip -q -O /tmp/MyricaM.zip \
    && (cd /tmp && unzip MyricaM.zip) \
    && git clone https://github.com/googlefonts/noto-emoji /usr/local/src/noto-emoji \
    && wget https://www.wfonts.com/download/data/2016/04/23/symbola/symbola.zip -q -O /tmp/symbola.zip \
    && (cd /tmp && unzip symbola.zip)

################################################################################

FROM base AS builder

COPY . /app
WORKDIR /app
RUN go install

################################################################################

FROM alpine:3.16.0 AS runtime
COPY --from=builder /go/bin/textimg /usr/local/bin/
COPY --from=builder /tmp/MyricaM.TTC /usr/share/fonts/truetype/myrica/MyricaM.TTC
COPY --from=builder /usr/local/src/noto-emoji/png/128 /usr/share/emoji-image
COPY --from=builder /tmp/Symbola_hint.ttf /usr/share/fonts/truetype/symbola/
COPY --from=builder /tmp/Symbola_hint.ttf /usr/share/fonts/truetype/ancient-scripts/

ENV TEXTIMG_OUTPUT_DIR /images
ENV TEXTIMG_FONT_FILE /usr/share/fonts/truetype/myrica/MyricaM.TTC
ENV TEXTIMG_EMOJI_DIR /usr/share/emoji-image
ENV TEXTIMG_EMOJI_FONT_FILE /usr/share/fonts/truetype/symbola/Symbola_hint.ttf
ENV LANG ja_JP.UTF-8
RUN mkdir /images

ENTRYPOINT ["/usr/local/bin/textimg"]

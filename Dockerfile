FROM golang:1.12.5-alpine3.9 as build-stage

ENV GO111MODULE on
RUN go version \
    && echo $GOPATH \
    && apk add --no-cache git wget unzip \
    && go get github.com/jiro4989/textimg \
    && git clone https://github.com/googlefonts/noto-emoji /usr/local/src/noto-emoji \
    && wget https://www.wfonts.com/download/data/2016/04/23/symbola/symbola.zip \
    && unzip symbola.zip

FROM alpine:latest as exec-stage
COPY --from=build-stage /root/go/bin/textimg /usr/local/bin/
COPY --from=build-stage /usr/local/src/noto-emoji /usr/local/src/
COPY --from=build-stage /Symbola_hint.ttf /usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf

ENV TEXTIMG_FONT_FILE /usr/share/fonts/TTF/HackGen-Regular.ttf
ENV TEXTIMG_EMOJI_FONT_FILE /usr/share/fonts/TTF/Symbola.ttf

ENTRYPOINT ["/usr/local/bin/textimg"]

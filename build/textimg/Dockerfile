FROM jiro4989/textimg-runtime AS runtime

FROM alpine:latest AS exec
COPY --from=runtime /go/bin/textimg /usr/local/bin/
COPY --from=runtime /tmp/MyricaM.TTC /usr/share/fonts/truetype/myrica/MyricaM.TTC
COPY --from=runtime /usr/local/src/noto-emoji/png/128 /usr/share/emoji-image
COPY --from=runtime /tmp/Symbola_hint.ttf /usr/share/fonts/truetype/symbola/
COPY --from=runtime /tmp/Symbola_hint.ttf /usr/share/fonts/truetype/ancient-scripts/

ENV TEXTIMG_FONT_FILE /usr/share/fonts/truetype/myrica/MyricaM.TTC
ENV TEXTIMG_EMOJI_DIR /usr/share/emoji-image
ENV TEXTIMG_EMOJI_FONT_FILE /usr/share/fonts/truetype/symbola/Symbola_hint.ttf
RUN mkdir /images

ENTRYPOINT ["/usr/local/bin/textimg"]

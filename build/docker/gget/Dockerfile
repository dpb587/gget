FROM alpine as base
RUN apk add --no-cache ca-certificates wget

ARG GGET_VERSION
ARG GGET_REF
ARG GITHUB_TOKEN
ENV gget /usr/local/bin/gget
RUN true \
  && wget -qO "$gget" https://github.com/dpb587/gget/releases/download/v0.1.0/gget-0.1.0-linux-amd64 \
  && echo "418b68caeb302898322465abc31c316a1dbdb7872cdad0fe4a9661c5118860c4  $gget" | sha256sum -c \
  && chmod +x "$gget"
RUN gget --exec ${GGET_REF:-github.com/dpb587/gget}@${GGET_VERSION:-} "$gget=gget-*-linux-amd64"
RUN gget --version -v

FROM alpine
RUN apk add --no-cache ca-certificates xz
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /usr/local/bin/gget /bin/gget
RUN mkdir /result
WORKDIR /result
ENTRYPOINT ["/bin/gget"]
CMD ["--help"]

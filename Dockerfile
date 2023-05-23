FROM golang:1.20.4-alpine AS build

LABEL maintainer="Iori Mizutani <iori.mizutani@gmail.com>"

# build the app
RUN mkdir -p /build
COPY go.mod /build/
COPY go.sum /build/
COPY *.go /build/
WORKDIR /build
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o radiko-auto-downloader .

# export to a single layer image
FROM alpine:latest

# install some required binaries
RUN apk add --no-cache ca-certificates \
  ffmpeg \
  rtmpdump \
  tzdata

COPY --from=build /build/radiko-auto-downloader /app/radiko-auto-downloader

# set timezone
ENV TZ "Asia/Tokyo"
# set the default download dir
ENV RADIGO_HOME "/downloads"
VOLUME ["/downloads"]

ENTRYPOINT ["/app/radiko-auto-downloader"]
CMD ["-c", "/app/config.toml"]

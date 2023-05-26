FROM golang:1.20.4-alpine AS build

LABEL maintainer="Iori Mizutani <iori.mizutani@gmail.com>"

# build the app
RUN mkdir -p /build
COPY go.mod /build/
COPY go.sum /build/
COPY *.go /build/
WORKDIR /build
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o radicron .

# export to a single layer image
FROM alpine:latest

# install some required binaries
RUN apk add --no-cache ca-certificates \
    ffmpeg \
    rtmpdump \
    tzdata

COPY --from=build /build/radicron /app/radicron

# set timezone
ENV TZ "Asia/Tokyo"
# set the default download dir
ENV RADIGO_HOME "/downloads"
VOLUME ["/downloads"]

ENTRYPOINT ["/app/radicron"]
CMD ["-c", "/app/config.yml"]

FROM 1.20.4-bullseye

LABEL maintainer="Iori Mizutani <iori.mizutani@gmail.com>"

RUN mkdir -p /build
COPY vendor /build/
COPY go.mod /build/
COPY go.sum /build/
COPY main.go /build/
WORKDIR /build
RUN go mod vendor
RUN go build -mod=vendor -o /app/radiko-auto-recorder .

ENTRYPOINT ["/app/radiko-auto-recorder"]
CMD ["-c", "/app/config.toml"]

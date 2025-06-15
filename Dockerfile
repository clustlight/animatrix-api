FROM golang:1.24-bullseye AS builder
WORKDIR /src
COPY . /src

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -o ./build/app .

FROM scratch

COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /usr/share/zoneinfo/Asia/Tokyo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /src/build/app /app
ENTRYPOINT ["/app"]
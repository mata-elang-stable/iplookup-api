FROM golang:1.14-stretch AS builder
WORKDIR /go/src/app/
COPY . /go/src/app/
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/main /go/src/app/

FROM scratch
COPY --from=builder /go/bin/main /app/
COPY --from=builder /go/src/app/assets/GeoLite2-City.mmdb /app/assets/
EXPOSE 80
CMD ["/app/main", "/app/assets/GeoLite2-City.mmdb"]
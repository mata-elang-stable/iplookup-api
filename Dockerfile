# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.23 AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /go/src/app/

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

COPY . /go/src/app/
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/main /go/src/app/cmd/iplookup-go/

FROM scratch
COPY --from=builder /go/bin/main /app/main
EXPOSE 3000
CMD ["/app/main"]
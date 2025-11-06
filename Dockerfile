FROM golang:1.25.3 as builder
COPY go.mod go.sum /go/src/oybek/io/sigma/
WORKDIR /go/src/oybek/io/sigma
RUN go mod download
COPY . /go/src/oybek/io/sigma
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/sigma oybek.io/sigma

FROM alpine/curl
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/oybek/io/sigma/build/sigma /usr/bin/sigma
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/sigma"]

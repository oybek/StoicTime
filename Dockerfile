FROM golang:1.24.2 as builder
COPY go.mod go.sum /go/src/oybek/io/kerege/
WORKDIR /go/src/oybek/io/kerege
RUN go mod download
COPY . /go/src/oybek/io/kerege
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/kerege oybek.io/kerege

FROM alpine/curl
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/oybek/io/kerege/build/kerege /usr/bin/kerege
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/kerege"]

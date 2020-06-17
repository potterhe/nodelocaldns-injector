FROM golang:1.14 as builder

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
#RUN go install -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -o nodelocaldns-injector

FROM alpine:3.12

COPY --from=builder /go/src/app/nodelocaldns-injector /nodelocaldns-injector
CMD ["/nodelocaldns-injector", "webhook"]

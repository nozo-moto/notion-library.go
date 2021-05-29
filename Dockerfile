FROM golang:1.16-buster as builder

WORKDIR /go/src/server
ADD . /go/src/server

RUN go get -d -v ./...
RUN go build -o /go/bin/server ./cmd/server 


FROM gcr.io/distroless/base-debian10
COPY --from=builder /go/bin/server /
CMD ["/server"]
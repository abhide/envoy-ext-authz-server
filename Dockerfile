FROM golang:alpine3.12
WORKDIR /go/src/github.com/abhide/envoy-ext-authz-server/
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go build -o envoy-ext-authz-server ./main.go

FROM alpine:3.12
WORKDIR /root/
COPY --from=0 /go/src/github.com/abhide/envoy-ext-authz-server/envoy-ext-authz-server .
CMD ["./envoy-ext-authz-server"]
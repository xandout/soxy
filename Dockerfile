FROM golang as builder

WORKDIR /go/src/github.com/xandout/soxy
COPY . ./

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/soxy
RUN chmod +x /go/bin/soxy


FROM alpine
COPY --from=builder /go/bin/soxy /
RUN apk add --no-cache \
        libc6-compat && \
        chmod +x /soxy
CMD ["/soxy"]


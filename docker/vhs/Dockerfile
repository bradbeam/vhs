FROM golang:1.15.2-alpine3.12 as builder

RUN go install -v -a std

RUN apk add --no-cache --virtual .build-deps \
    ca-certificates \
    g++ \
    libpcap-dev

WORKDIR /src

COPY go.* .

RUN go mod download

COPY . .

RUN GODEBUG=netdns=go go build \
	--ldflags '-extldflags "-static"' \
	./cmd/vhs

########################################################

FROM scratch

COPY --from=builder /src/vhs /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/vhs"]

########################################################

FROM builer as test

RUN go test -v ./...


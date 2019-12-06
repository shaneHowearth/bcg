FROM golang:1.13 as builder

# modules
WORKDIR $GOPATH/src/notify
ADD . $GOPATH/src/notify

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN apt-get install ca-certificates -y
# build time
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/notify cmd/main.go

# stage 2
FROM scratch
WORKDIR /root/
COPY --from=0 /go/bin/notify .
COPY --from=0 /go/src/notify/internal/emailprovider/gmail/credentials.json ./
# This will fail if the token file has not been created
# If this fails please `cd $PROJECT/notify/internal/emailprovider/gmail/getToken`
# and run `go run main.go` to get a token file from Google
COPY --from=0 /go/src/notify/internal/emailprovider/gmail/token.json ./

COPY --from=0 /etc/ssl/certs/* /etc/ssl/certs/

ENTRYPOINT ["./notify"]

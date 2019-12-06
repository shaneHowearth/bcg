FROM golang:1.13 as builder

# modules
WORKDIR $GOPATH/src/github.com/shanehowearth/bcg/restserver
ADD . $GOPATH/src/github.com/shanehowearth/bcg/restserver

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# build time
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/restserver cmd/main.go cmd/routes.go

FROM scratch
WORKDIR /root/
COPY --from=0 /go/bin/restserver .

# run the rest server
ENTRYPOINT ["./restserver"]

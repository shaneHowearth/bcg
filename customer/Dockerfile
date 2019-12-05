FROM golang:1.13 as builder

# modules
WORKDIR $GOPATH/src/customer
ADD . $GOPATH/src/customer

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

# build time
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/customer cmd/main.go

# stage 2
FROM redis
WORKDIR /root/
COPY --from=0 /go/bin/customer .
COPY . .
COPY customer.sh .
RUN chmod 777 customer.sh

# run the customer script
ENTRYPOINT ["./customer.sh"]

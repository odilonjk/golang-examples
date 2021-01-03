FROM golang:1.15.6 AS builder
WORKDIR /go/src/github.com/odilonjk/golang-examples/consistencia-dados 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/odilonjk/golang-examples/consistencia-dados/app .
CMD ["./app"]  

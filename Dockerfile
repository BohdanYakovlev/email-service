FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .
ADD go.sum .

COPY . .

RUN go mod tidy

RUN go build -o email-service email-service/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/email-service /build/email-service

CMD ["/build/email-service/main"]
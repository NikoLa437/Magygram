FROM golang:1.15-alpine3.13
WORKDIR /app
RUN apk update \
    && apk add sqlite \
    && apk add socat \
    && apk add build-base
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o message-service
EXPOSE 467
CMD ["./message-service"]
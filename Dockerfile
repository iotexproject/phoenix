FROM golang:1.14-alpine as build

WORKDIR apps/phoenix

RUN apk add --no-cache make gcc musl-dev linux-headers git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o phoenix -v ./

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=build /go/apps/phoenix/phoenix /usr/local/bin/iotex-phoenix

CMD [ "iotex-phoenix"]

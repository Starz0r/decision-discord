FROM golang:1.12.17-alpine3.10 as builder

RUN apk update
RUN apk upgrade
RUN apk add git
RUN mkdir /build
ADD . /build/
WORKDIR /build
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -i -v -a -installsuffix cgo -ldflags '-extldflags "-static"' -o decision ./src/

FROM alpine:3.12.0

RUN apk update && apk upgrade && apk add ca-certificates
COPY --from=builder /build/decision /app/
WORKDIR /app
CMD ["./decision"]
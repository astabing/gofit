FROM golang:1.18-alpine as builder

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

COPY /go.* /build/
COPY /cmd /build/cmd
COPY /vendor /build/vendor
WORKDIR /build 

RUN apk add --no-cache make git tzdata ca-certificates
#	&& adduser -D -H -g "" -s "" -u 1001 gofit

RUN CGO_ENABLED=0 go build -a -o gofit ./cmd/gofit

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /build/gofit /app/

#USER gofit:gofit
WORKDIR /app

ENTRYPOINT ["/app/gofit"]

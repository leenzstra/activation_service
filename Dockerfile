# builder image
FROM golang:1.19-alpine as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o license_server ./cmd/main.go


# generate clean, final image for end users
FROM alpine:latest
COPY --from=builder /build/license_server .

# executable
ENTRYPOINT [ "./license_server" ]
# arguments that can be overridden
# CMD [ "3", "300" ]
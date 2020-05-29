FROM golang:1.14.2-alpine AS builder
COPY . /app/
WORKDIR /app
RUN go build -o app

# --

FROM maxisme/megatools-alpine
RUN apk add --update rsync openssh

WORKDIR /app
COPY . /app/
COPY --from=builder /app/app /app/app
CMD ["/app/app"]
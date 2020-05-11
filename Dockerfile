FROM golang:1.14.2-alpine AS builder
COPY . /app/
WORKDIR /app
RUN go build -o app

FROM alpine
RUN apk add --update build-base libcurl curl-dev openssl-dev glib-dev glib libtool automake autoconf rsync

# install mega
ARG mega_version=1.10.3
RUN wget https://megatools.megous.com/builds/megatools-$mega_version.tar.gz
RUN tar -xzf megatools-$mega_version.tar.gz
RUN bash megatools-$mega_version/configure --disable-docs
RUN make -j4
RUN make install
# cleanup
RUN rm -rf megatools*

WORKDIR /app
COPY . /app/
COPY --from=builder /app/app /app/app
CMD ["/app/app"]